package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"coresense/pkg/common/utils/jwtutils"
	"coresense/pkg/common/utils/validation"
	"coresense/pkg/core/model/dto"
)

func (p *UserProcessor) AddFilledService(w http.ResponseWriter, r *http.Request) {
	// Extract claims from the context
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unable to retrieve claims from context", http.StatusInternalServerError)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	var user dto.User
	result := p.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Decode request body
	var serviceRequest dto.FilledServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&serviceRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Fetch ServiceTemplate
	var serviceTemplate dto.ServiceTemplate
	if err := p.db.Where("id = ?", serviceRequest.ServiceTemplateID).First(&serviceTemplate).Error; err != nil {
		http.Error(w, "Service template not found", http.StatusBadRequest)
		return
	}

	result = p.db.Where("service_template_id = ?", serviceTemplate.ID).
		Find(&serviceTemplate.FieldsFormat)
	if result.Error != nil {
		http.Error(w, "Questions for service template not found", http.StatusInternalServerError)
		return
	}

	// sort by number
	sort.Slice(serviceTemplate.FieldsFormat, func(i, j int) bool {
		return serviceTemplate.FieldsFormat[i].Number < serviceTemplate.FieldsFormat[j].Number
	})

	// Validate ServiceData length
	if len(serviceRequest.ServiceData) != len(serviceTemplate.FieldsFormat) {
		http.Error(w, "Service data length mismatch", http.StatusBadRequest)
		return
	}

	// Set UserID and timestamps
	answer := dto.FilledService{
		ID:                uuid.New(),
		ServiceTemplateID: serviceTemplate.ID,
		UserID:            user.ID,
		CreatedAt:         time.Now(),
	}

	// Validate data types
	for i, question := range serviceTemplate.FieldsFormat {
		if !validation.ValidateFieldType(serviceRequest.ServiceData[i].Answer, question.Type) {
			http.Error(w, "Invalid data type for field", http.StatusBadRequest)
			return
		}

		// Check script validation if ScriptID is present
		if question.ScriptID.Valid {
			var script dto.Script
			if err := p.db.Where("id = ?", question.ScriptID).First(&script).Error; err != nil {
				http.Error(w, "Script not found", http.StatusInternalServerError)
				return
			}

			// Run the script with serviceData[i]
			valid, err := validation.RunValidationScript(script.ScriptCode, serviceRequest.ServiceData[i].Answer)
			if err != nil {
				http.Error(w, "Error executing validation script", http.StatusInternalServerError)
				return
			}

			if !valid {
				http.Error(w, "Validation script failed for field", http.StatusBadRequest)
				return
			}
		}

		answer.ServiceData = append(answer.ServiceData, dto.QuestionAnswered{
			ID:              uuid.New(),
			FilledServiceID: answer.ID,
			Answer:          serviceRequest.ServiceData[i].Answer,
			QuestionID:      question.ID,
			Number:          question.Number,
			CreatedAt:       time.Now(),
		})
	}

	// Save to database
	if err := p.db.Create(&answer).Error; err != nil {
		http.Error(w, "Error saving filled service", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(map[string]string{"message": "Filled service created successfully"})
	if err != nil {
		http.Error(w, "Error saving filled service", http.StatusInternalServerError)
		return
	}
}

// GetFilledServices retrieves FilledServices based on optional UserID and ServiceTemplateID.
func (p *BusinessProcessor) GetFilledServices(w http.ResponseWriter, r *http.Request) {
	// Extract JWT from request header
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString, err := jwtutils.GetTokenFromHeader(tokenHeader)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Validate token and get BusinessCustomerID
	claims, err := jwtutils.GetClaims(tokenString, jwtSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	businessCustomerName, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	var businessCustomer dto.BusinessCustomer
	result := p.db.Where("name = ?", businessCustomerName).First(&businessCustomer)
	if result.Error != nil {
		http.Error(w, "Business customer not found", http.StatusUnauthorized)
		return
	}

	// Decode request body
	var request struct {
		Count             int        `json:"count"`
		Skip              int        `json:"skip"`
		ServiceTemplateID *uuid.UUID `json:"service_template_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate count and skip parameters
	if request.Count <= 0 {
		http.Error(w, "Count must be greater than 0", http.StatusBadRequest)
		return
	}
	if request.Skip < 0 {
		http.Error(w, "Skip cannot be negative", http.StatusBadRequest)
		return
	}
	if request.ServiceTemplateID == nil {
		http.Error(w, "ServiceTemplateID cannot be null", http.StatusBadRequest)
		return
	}

	var total int64
	if err := p.db.Model(&dto.FilledService{}).
		Where("service_template_id = ?", *request.ServiceTemplateID).Count(&total).Error; err != nil {
		http.Error(w, "Error retrieving filled services", http.StatusInternalServerError)
		return
	}

	var filledServices []dto.FilledService
	query := p.db.Model(&dto.FilledService{}).Where("service_template_id = ?", *request.ServiceTemplateID)
	if err := query.Limit(request.Count).Offset(request.Skip).Find(&filledServices).Error; err != nil {
		http.Error(w, "Error retrieving filled services", http.StatusInternalServerError)
		return
	}

	var filledServiceGridRecords []dto.FilledServiceGridRecord

	for k := range len(filledServices) {
		var serviceTemplate dto.ServiceTemplate
		result := p.db.Where("id = ?", filledServices[k].ServiceTemplateID).First(&serviceTemplate)
		if result.Error != nil {
			http.Error(w, "Service template not found", http.StatusInternalServerError)
			return
		}
		// sort by number
		filledServiceGridRecords = append(filledServiceGridRecords, dto.FilledServiceGridRecord{
			ID:                         filledServices[k].ID,
			ServiceTemplateName:        serviceTemplate.Name,
			ServiceTemplateDescription: serviceTemplate.Description,
			CreatedAt:                  filledServices[k].CreatedAt,
		})
	}

	// Prepare the response with `data` and `total` fields
	response := struct {
		Data  []dto.FilledServiceGridRecord `json:"data"`
		Total int64                         `json:"total"`
	}{
		Data:  filledServiceGridRecords,
		Total: total,
	}

	// Respond with filled services
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetFilledServices retrieves FilledServices based on optional UserID and ServiceTemplateID.
func (p *UserProcessor) GetFilledServices(w http.ResponseWriter, r *http.Request) {
	// Extract JWT from request header
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString, err := jwtutils.GetTokenFromHeader(tokenHeader)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Validate token and get UserID
	claims, err := jwtutils.GetClaims(tokenString, jwtSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userName, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	var user dto.User
	result := p.db.Where("username = ?", userName).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Decode request body
	var request struct {
		Count int `json:"count"`
		Skip  int `json:"skip"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate count and skip parameters
	if request.Count <= 0 {
		http.Error(w, "Count must be greater than 0", http.StatusBadRequest)
		return
	}
	if request.Skip < 0 {
		http.Error(w, "Skip cannot be negative", http.StatusBadRequest)
		return
	}

	// Build query
	var total int64
	if err := p.db.Model(&dto.FilledService{}).
		Where("user_id = ?", user.ID).
		Count(&total).
		Error; err != nil {
		http.Error(w, "Error retrieving filled services", http.StatusInternalServerError)
		return
	}

	var filledServices []dto.FilledService
	query := p.db.Model(&dto.FilledService{}).Where("user_id = ?", user.ID)
	if err := query.Limit(request.Count).Offset(request.Skip).Find(&filledServices).Error; err != nil {
		http.Error(w, "Error retrieving filled services", http.StatusInternalServerError)
		return
	}

	var filledServiceGridRecords []dto.FilledServiceGridRecord

	for k := range len(filledServices) {
		var serviceTemplate dto.ServiceTemplate
		result := p.db.Where("id = ?", filledServices[k].ServiceTemplateID).First(&serviceTemplate)
		if result.Error != nil {
			http.Error(w, "Service template not found", http.StatusInternalServerError)
			return
		}
		// sort by number
		filledServiceGridRecords = append(filledServiceGridRecords, dto.FilledServiceGridRecord{
			ID:                         filledServices[k].ID,
			ServiceTemplateName:        serviceTemplate.Name,
			ServiceTemplateDescription: serviceTemplate.Description,
			CreatedAt:                  filledServices[k].CreatedAt,
		})
	}

	// Prepare the response with `data` and `total` fields
	response := struct {
		Data  []dto.FilledServiceGridRecord `json:"data"`
		Total int64                         `json:"total"`
	}{
		Data:  filledServiceGridRecords,
		Total: total,
	}

	// Respond with filled services
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (p *UserProcessor) GetFilledServiceWithDetails(w http.ResponseWriter, r *http.Request) {
	// Extract claims from the context
	_, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unable to retrieve claims from context", http.StatusInternalServerError)
		return
	}

	// Decode request body
	var request struct {
		FilledServiceID *uuid.UUID `json:"filled_service_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate
	if request.FilledServiceID == nil {
		http.Error(w, "Invalid request, missing FilledServiceID", http.StatusBadRequest)
		return
	}

	var filledService dto.FilledService
	result := p.db.Where("id = ?", request.FilledServiceID).First(&filledService)
	if result.Error != nil {
		http.Error(w, "Filled service not found", http.StatusInternalServerError)
		return
	}

	result = p.db.Where("filled_service_id = ?", request.FilledServiceID).
		Find(&filledService.ServiceData)
	if result.Error != nil {
		http.Error(w, "Questions answered for filled service not found", http.StatusInternalServerError)
		return
	}

	// sort by number
	if len(filledService.ServiceData) == 0 {
		http.Error(w, "Filled service (service data) not found", http.StatusInternalServerError)
		return
	}

	sort.Slice(filledService.ServiceData, func(i, j int) bool {
		return filledService.ServiceData[i].Number < filledService.ServiceData[j].Number
	})

	var serviceTemplate dto.ServiceTemplate
	result = p.db.Where("id = ?", filledService.ServiceTemplateID).First(&serviceTemplate)
	if result.Error != nil {
		http.Error(w, "Service template not found", http.StatusInternalServerError)
		return
	}

	result = p.db.Where("service_template_id = ?", serviceTemplate.ID).
		Find(&serviceTemplate.FieldsFormat)
	if result.Error != nil {
		http.Error(w, "Questions for service template not found", http.StatusInternalServerError)
		return
	}

	if len(serviceTemplate.FieldsFormat) == 0 {
		http.Error(w, "Service template (fields format) not found", http.StatusInternalServerError)
		return
	}

	// sort by number
	sort.Slice(serviceTemplate.FieldsFormat, func(i, j int) bool {
		return serviceTemplate.FieldsFormat[i].Number < serviceTemplate.FieldsFormat[j].Number
	})

	var response = dto.FilledServiceWithDetails{
		ServiceTemplateName:        serviceTemplate.Name,
		ServiceTemplateDescription: serviceTemplate.Description,
		CreatedAt:                  filledService.CreatedAt,
	}

	for i := range len(filledService.ServiceData) {
		response.ServiceData = append(response.ServiceData, dto.QuestionAnsweredGridRecord{
			Question: serviceTemplate.FieldsFormat[i].Description,
			Answer:   filledService.ServiceData[i].Answer,
		})
	}

	// Respond with service templates
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (p *BusinessProcessor) GetFilledServiceWithDetails(w http.ResponseWriter, r *http.Request) {
	// Extract claims from the context
	_, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unable to retrieve claims from context", http.StatusInternalServerError)
		return
	}

	// Decode request body
	var request struct {
		FilledServiceID *uuid.UUID `json:"filled_service_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate
	if request.FilledServiceID == nil {
		http.Error(w, "Invalid request, missing FilledServiceID", http.StatusBadRequest)
		return
	}

	var filledService dto.FilledService
	result := p.db.Where("id = ?", request.FilledServiceID).First(&filledService)
	if result.Error != nil {
		http.Error(w, "Filled service not found", http.StatusInternalServerError)
		return
	}

	result = p.db.Where("filled_service_id = ?", request.FilledServiceID).
		Find(&filledService.ServiceData)
	if result.Error != nil {
		http.Error(w, "Questions answered for filled service not found", http.StatusInternalServerError)
		return
	}

	// sort by number
	if len(filledService.ServiceData) == 0 {
		http.Error(w, "Filled service (service data) not found", http.StatusInternalServerError)
		return
	}

	sort.Slice(filledService.ServiceData, func(i, j int) bool {
		return filledService.ServiceData[i].Number < filledService.ServiceData[j].Number
	})

	var serviceTemplate dto.ServiceTemplate
	result = p.db.Where("id = ?", filledService.ServiceTemplateID).First(&serviceTemplate)
	if result.Error != nil {
		http.Error(w, "Service template not found", http.StatusInternalServerError)
		return
	}

	result = p.db.Where("service_template_id = ?", serviceTemplate.ID).
		Find(&serviceTemplate.FieldsFormat)
	if result.Error != nil {
		http.Error(w, "Questions for service template not found", http.StatusInternalServerError)
		return
	}

	if len(serviceTemplate.FieldsFormat) == 0 {
		http.Error(w, "Service template (fields format) not found", http.StatusInternalServerError)
		return
	}

	// sort by number
	sort.Slice(serviceTemplate.FieldsFormat, func(i, j int) bool {
		return serviceTemplate.FieldsFormat[i].Number < serviceTemplate.FieldsFormat[j].Number
	})

	var response = dto.FilledServiceWithDetails{
		ServiceTemplateName:        serviceTemplate.Name,
		ServiceTemplateDescription: serviceTemplate.Description,
		CreatedAt:                  filledService.CreatedAt,
		User:                       filledService.UserID.String(), // todo: add user name instead
	}

	for i := range len(filledService.ServiceData) {
		response.ServiceData = append(response.ServiceData, dto.QuestionAnsweredGridRecord{
			Question: serviceTemplate.FieldsFormat[i].Description,
			Answer:   filledService.ServiceData[i].Answer,
		})
	}

	// Respond with service templates
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetStatistics retrieves and processes statistics based on filled services
func (p *BusinessProcessor) GetStatistics(w http.ResponseWriter, r *http.Request) {
	p.logger.Debug().Msg("In")
	defer p.logger.Debug().Msg("Out")

	// Decode request body
	var request struct {
		UserID            *uuid.UUID `json:"user_id"`
		ServiceTemplateID *uuid.UUID `json:"service_template_id"`
		ScriptID          uuid.UUID  `json:"script_id"`
		QuestionNumber    int        `json:"question_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Query filled services
	query := p.db.Model(&dto.FilledService{})
	if request.UserID != nil {
		query = query.Where("user_id = ?", *request.UserID)
	}
	if request.ServiceTemplateID != nil {
		query = query.Where("service_template_id = ?", *request.ServiceTemplateID)
	}

	var filledServices []dto.FilledService
	if err := query.Find(&filledServices).Error; err != nil {
		http.Error(w, "Error retrieving filled services", http.StatusInternalServerError)
		return
	}

	// Extract answers for the specified question number
	var answers []string
	for _, filledService := range filledServices {
		for _, questionAnswered := range filledService.ServiceData {
			if questionAnswered.Number == request.QuestionNumber {
				answers = append(answers, questionAnswered.Answer)
				break
			}
		}
	}

	// Fetch the script
	var script dto.Script
	if err := p.db.Where("id = ?", request.ScriptID).First(&script).Error; err != nil {
		http.Error(w, "Script not found", http.StatusInternalServerError)
		return
	}

	// Run the script with the collected answers
	result, err := validation.RunStatisticsScript(script.ScriptCode, answers)
	if err != nil {
		http.Error(w, "Error executing validation script", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"result": result})
	if err != nil {
		return
	}
}
