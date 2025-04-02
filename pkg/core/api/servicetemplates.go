package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"coresense/pkg/common/utils/jwtutils"
	"coresense/pkg/core/model/dto"
)

func (p *BusinessProcessor) AddServiceTemplate(w http.ResponseWriter, r *http.Request) {
	p.logger.Debug().Msg("In AddServiceTemplate")
	p.logger.Debug().Msg("Out AddServiceTemplate")
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
		p.logger.Err(err).Send()
		http.Error(w, "Invalid token, missing claims", http.StatusUnauthorized)
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
	var serviceTemplateRequest dto.ServiceTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&serviceTemplateRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Set BusinessCustomerID and other necessary fields
	var serviceTemplate = dto.ServiceTemplate{
		ID:                 uuid.New(),
		BusinessCustomerID: businessCustomer.ID,
		Name:               serviceTemplateRequest.Name,
		Description:        serviceTemplateRequest.Description,
		CreatedAt:          time.Now(),
	}

	for i := range len(serviceTemplateRequest.FieldsFormat) {
		serviceTemplate.FieldsFormat = append(serviceTemplate.FieldsFormat, dto.Question{
			ID:                uuid.New(),
			ServiceTemplateID: serviceTemplate.ID,
			Type:              serviceTemplateRequest.FieldsFormat[i].Type,
			Description:       serviceTemplateRequest.FieldsFormat[i].Description,
			Number:            i,
			CreatedAt:         time.Now(),
		})
		if serviceTemplateRequest.FieldsFormat[i].ScriptID.Valid {
			var script dto.Script
			if err := p.db.Model(&dto.Script{}).
				Where("id = ?", serviceTemplateRequest.FieldsFormat[i].ScriptID.UUID).
				First(&script).Error; err != nil {
				http.Error(w, "Error updating script", http.StatusInternalServerError)
				return
			}
			serviceTemplate.FieldsFormat[i].ScriptID = serviceTemplateRequest.FieldsFormat[i].ScriptID
		}
	}

	// Save to database
	if err := p.db.Create(&serviceTemplate).Error; err != nil {
		http.Error(w, "Error saving service template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Service template created successfully"})
	if err != nil {
		http.Error(w, "Internal json error", http.StatusInternalServerError)
		return
	}
}

// EditServiceTemplate allows a business customer to edit their ServiceTemplate.
func (p *BusinessProcessor) EditServiceTemplate(w http.ResponseWriter, r *http.Request) {
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
	var updatedServiceTemplate dto.ServiceTemplate
	if err := json.NewDecoder(r.Body).Decode(&updatedServiceTemplate); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve existing service template
	var existingServiceTemplate dto.ServiceTemplate
	if err := p.db.Where("id = ? AND business_customer_id = ?", updatedServiceTemplate.ID, businessCustomer.ID).First(&existingServiceTemplate).Error; err != nil {
		http.Error(w, "Service template not found or unauthorized", http.StatusForbidden)
		return
	}

	// Update fields
	if updatedServiceTemplate.Name != "" {
		existingServiceTemplate.Name = updatedServiceTemplate.Name
	}
	if updatedServiceTemplate.Description != "" {
		existingServiceTemplate.Description = updatedServiceTemplate.Description
	}
	if len(updatedServiceTemplate.FieldsFormat) != 0 {
		existingServiceTemplate.FieldsFormat = updatedServiceTemplate.FieldsFormat
	}
	existingServiceTemplate.Description = updatedServiceTemplate.Description
	existingServiceTemplate.FieldsFormat = updatedServiceTemplate.FieldsFormat
	updatedTime := time.Now()
	existingServiceTemplate.UpdatedAt = &updatedTime

	for i := range len(existingServiceTemplate.FieldsFormat) {
		if existingServiceTemplate.FieldsFormat[i].ScriptID.Valid {
			var script dto.Script
			if err := p.db.Model(&dto.Script{}).
				Where("id = ?", existingServiceTemplate.FieldsFormat[i].ScriptID.UUID).
				First(&script).Error; err != nil {
				http.Error(w, "Error updating script", http.StatusInternalServerError)
				return
			}
		}
		existingServiceTemplate.FieldsFormat[i].Number = i
		existingServiceTemplate.FieldsFormat[i].UpdatedAt = &updatedTime
	}

	// Save changes
	if err := p.db.Save(&existingServiceTemplate).Error; err != nil {
		http.Error(w, "Error updating service template", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "Service template updated successfully"})
	if err != nil {
		http.Error(w, "Internal json error", http.StatusInternalServerError)
		return
	}
}

// GetServiceTemplates retrieves ServiceTemplates based on required BusinessCustomerID.
func (p *UserProcessor) GetServiceTemplates(w http.ResponseWriter, r *http.Request) {
	// Extract claims from the context
	_, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unable to retrieve claims from context", http.StatusInternalServerError)
		return
	}

	// Decode request body
	var request struct {
		BusinessCustomerID *uuid.UUID `json:"user_id"`
		Count              int        `json:"count"`
		Skip               int        `json:"skip"`
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
	if request.BusinessCustomerID == nil {
		http.Error(w, "Invalid request, missing BusinessCustomerID", http.StatusBadRequest)
		return
	}

	// Get the complete total number of business users
	var total int64
	if err := p.db.Model(&dto.ServiceTemplate{}).
		Where("business_customer_id = ?", *request.BusinessCustomerID).
		Count(&total).Error; err != nil {
		http.Error(w, "Error counting business users", http.StatusInternalServerError)
		return
	}

	// Build query
	query := p.db.Model(&dto.ServiceTemplate{})
	query = query.Where("business_customer_id = ?", *request.BusinessCustomerID)

	var serviceTemplates []dto.ServiceTemplateGridRecord
	if err := query.Limit(request.Count).Offset(request.Skip).Find(&serviceTemplates).Error; err != nil {
		http.Error(w, "Error retrieving service templates", http.StatusInternalServerError)
		return
	}

	// Prepare the response with `data` and `total` fields
	response := struct {
		Data  []dto.ServiceTemplateGridRecord `json:"data"`
		Total int64                           `json:"total"`
	}{
		Data:  serviceTemplates,
		Total: total,
	}

	// Encode the result as JSON and send it back to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GetServiceTemplates retrieves ServiceTemplates based on required BusinessCustomerID.
func (p *BusinessProcessor) GetServiceTemplates(w http.ResponseWriter, r *http.Request) {
	// Extract claims from the context
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unable to retrieve claims from context", http.StatusInternalServerError)
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

	// Get the complete total number of business users
	var total int64
	if err := p.db.Model(&dto.ServiceTemplate{}).
		Where("business_customer_id = ?", businessCustomer.ID).
		Count(&total).Error; err != nil {
		http.Error(w, "Error counting business users", http.StatusInternalServerError)
		return
	}

	// Build query
	query := p.db.Model(&dto.ServiceTemplate{})
	query = query.Where("business_customer_id = ?", businessCustomer.ID)

	var serviceTemplates []dto.ServiceTemplateGridRecord
	if err := query.Limit(request.Count).Offset(request.Skip).Find(&serviceTemplates).Error; err != nil {
		http.Error(w, "Error retrieving service templates", http.StatusInternalServerError)
		return
	}

	// Prepare the response with `data` and `total` fields
	response := struct {
		Data  []dto.ServiceTemplateGridRecord `json:"data"`
		Total int64                           `json:"total"`
	}{
		Data:  serviceTemplates,
		Total: total,
	}

	// Encode the result as JSON and send it back to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GetServiceTemplateWithQuestions retrieves ServiceTemplate based on its ID.
func (p *UserProcessor) GetServiceTemplateWithQuestions(w http.ResponseWriter, r *http.Request) {
	// Extract claims from the context
	_, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unable to retrieve claims from context", http.StatusInternalServerError)
		return
	}

	// Decode request body
	var request struct {
		ServiceTemplateID *uuid.UUID `json:"service_template_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate
	if request.ServiceTemplateID == nil {
		http.Error(w, "Invalid request, missing ServiceTemplateID", http.StatusBadRequest)
		return
	}

	var serviceTemplate dto.ServiceTemplate
	result := p.db.Where("id = ?", request.ServiceTemplateID).First(&serviceTemplate)
	if result.Error != nil {
		http.Error(w, "Service template not found", http.StatusInternalServerError)
		return
	}

	result = p.db.Where("service_template_id = ?", request.ServiceTemplateID).
		Find(&serviceTemplate.FieldsFormat)
	if result.Error != nil {
		http.Error(w, "Questions for service template not found", http.StatusInternalServerError)
		return
	}
	// sort by number
	sort.Slice(serviceTemplate.FieldsFormat, func(i, j int) bool {
		return serviceTemplate.FieldsFormat[i].Number < serviceTemplate.FieldsFormat[j].Number
	})

	// Respond with service templates
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(serviceTemplate); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
