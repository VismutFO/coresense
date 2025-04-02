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

type BusinessProcessor struct {
	config config.Server
	db     *gorm.DB
	logger zerolog.Logger
}

func NewBusinessProcessor(config config.Server, db *gorm.DB, logger zerolog.Logger) *BusinessProcessor {
	return &BusinessProcessor{
		config: config,
		db:     db,
		logger: logger,
	}
}

func (p *BusinessProcessor) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/business/register", p.RegisterBusiness)
	mux.HandleFunc("/business/login", p.LoginBusiness)

	// Protect the dashboard route with middleware
	mux.Handle("/business/dashboard", authMiddleware("business", http.HandlerFunc(p.Dashboard)))
	mux.Handle("/business/addServiceTemplate",
		authMiddleware("business", http.HandlerFunc(p.AddServiceTemplate)))
	mux.Handle("/business/editServiceTemplate",
		authMiddleware("business", http.HandlerFunc(p.EditServiceTemplate)))
	mux.Handle("/business/getFilledServices",
		authMiddleware("business", http.HandlerFunc(p.GetFilledServices)))
	mux.Handle("/business/addScript",
		authMiddleware("business", http.HandlerFunc(p.AddScript)))
	mux.Handle("/business/editScript",
		authMiddleware("business", http.HandlerFunc(p.EditScript)))
	mux.Handle("/business/getServiceTemplates",
		authMiddleware("business", http.HandlerFunc(p.GetServiceTemplates)))
	mux.Handle("/business/getFilledServiceWithDetails",
		authMiddleware("business", http.HandlerFunc(p.GetFilledServiceWithDetails)))
}

func (p *BusinessProcessor) RegisterBusiness(w http.ResponseWriter, r *http.Request) {
	var businessCustomer dto.BusinessCustomer
	if err := json.NewDecoder(r.Body).Decode(&businessCustomer); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(businessCustomer.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	businessCustomer.Password = string(hashedPassword)
	businessCustomer.CreatedAt = time.Now()
	businessCustomer.ID = uuid.New()

	if err := p.db.Create(&businessCustomer).Error; err != nil {
		http.Error(w, "Error saving filled service", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).
		Encode(map[string]string{"message": "Customer registered successfully"})
	if err != nil {
		http.Error(w, "Invalid request", http.StatusInternalServerError)
		return
	}
}

func (p *BusinessProcessor) LoginBusiness(w http.ResponseWriter, r *http.Request) {
	var businessCustomer dto.BusinessCustomer
	if err := json.NewDecoder(r.Body).Decode(&businessCustomer); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var foundBusinessCustomer dto.BusinessCustomer
	result := p.db.Where("name = ?", businessCustomer.Name).First(&foundBusinessCustomer)
	if result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(foundBusinessCustomer.Password), []byte(businessCustomer.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWT(foundBusinessCustomer.Name, map[string]any{
		"role": "business",
		"name": foundBusinessCustomer.Name,
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

func (p *BusinessProcessor) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Implement the dashboard logic here
	err := json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the business dashboard"})
	if err != nil {
		http.Error(w, "Internal json error", http.StatusInternalServerError)
		return
	}
}
