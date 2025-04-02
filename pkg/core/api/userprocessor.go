package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"coresense/app/server/config"
	"coresense/pkg/core/model/dto"
)

type UserProcessor struct {
	config config.Server
	db     *gorm.DB
	logger zerolog.Logger
}

func NewUserProcessor(config config.Server, db *gorm.DB, logger zerolog.Logger) *UserProcessor {
	return &UserProcessor{
		config: config,
		db:     db,
		logger: logger,
	}
}

func (p *UserProcessor) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/register", p.RegisterUser)
	mux.HandleFunc("/login", p.AuthorizeUser)

	// Protect the dashboard route with middleware
	mux.Handle("/dashboard", authMiddleware("user", http.HandlerFunc(p.Dashboard)))
	mux.Handle("/addFilledService", authMiddleware("user", http.HandlerFunc(p.AddFilledService)))
	mux.Handle("/getServiceTemplates",
		authMiddleware("user", http.HandlerFunc(p.GetServiceTemplates)))
	mux.Handle("/getServiceTemplateWithQuestions",
		authMiddleware("user", http.HandlerFunc(p.GetServiceTemplateWithQuestions)))
	mux.Handle("/getBusinessUsers", authMiddleware("user", http.HandlerFunc(p.GetBusinessUsers)))
	mux.Handle("/getFilledServices",
		authMiddleware("user", http.HandlerFunc(p.GetFilledServices)))
	mux.Handle("/getFilledServiceWithDetails",
		authMiddleware("user", http.HandlerFunc(p.GetFilledServiceWithDetails)))
}

func (p *UserProcessor) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.ID = uuid.New()

	if err := p.db.Create(&user).Error; err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
	if err != nil {
		http.Error(w, "Invalid request", http.StatusInternalServerError)
		return
	}
}

// AuthorizeUser handles user authentication.
func (p *UserProcessor) AuthorizeUser(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var foundUser dto.User
	result := p.db.Where("username = ?", user.Username).First(&foundUser)
	if result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWT(foundUser.Username, map[string]any{
		"role":     "user",
		"username": foundUser.Username,
	}) // second parameter is to add permissions later
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "token": token})
	if err != nil {
		http.Error(w, "Internal json error", http.StatusInternalServerError)
		return
	}
}

func (p *UserProcessor) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Implement the dashboard logic here
	err := json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the user dashboard"})
	if err != nil {
		http.Error(w, "Internal json error", http.StatusInternalServerError)
		return
	}
}

// GetBusinessUsers retrieves a paginated list of business users.
func (p *UserProcessor) GetBusinessUsers(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a struct
	var request struct {
		Count int `json:"count"`
		Skip  int `json:"skip"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
	if err := p.db.Model(&dto.BusinessCustomer{}).Count(&total).Error; err != nil {
		http.Error(w, "Error counting business users", http.StatusInternalServerError)
		return
	}

	// Query the database to fetch the business users with pagination
	var businessUsers []dto.BusinessCustomerGridRecord
	err := p.db.Model(&dto.BusinessCustomer{}).
		Limit(request.Count).
		Offset(request.Skip).
		Find(&businessUsers).
		Error
	if err != nil {
		http.Error(w, "Error retrieving business users", http.StatusInternalServerError)
		return
	}

	// Prepare the response with `data` and `total` fields
	response := struct {
		Data  []dto.BusinessCustomerGridRecord `json:"data"`
		Total int64                            `json:"total"`
	}{
		Data:  businessUsers,
		Total: total,
	}

	// Encode the result as JSON and send it back to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
