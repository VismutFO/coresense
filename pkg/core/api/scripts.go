package api

import (
	"coresense/pkg/core/model/dto"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// AddScript allows a BusinessCustomer to add a new script.
func (p *BusinessProcessor) AddScript(w http.ResponseWriter, r *http.Request) {
	var script dto.Script
	if err := json.NewDecoder(r.Body).Decode(&script); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	script.ID = uuid.New()
	script.CreatedAt = time.Now()

	if err := p.db.Create(&script).Error; err != nil {
		http.Error(w, "Error saving script", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(script)
	if err != nil {
		http.Error(w, "Error saving script", http.StatusInternalServerError)
		return
	}
}

// EditScript allows a BusinessCustomer to edit an existing script.
func (p *BusinessProcessor) EditScript(w http.ResponseWriter, r *http.Request) {
	var script dto.Script
	if err := json.NewDecoder(r.Body).Decode(&script); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	script.UpdatedAt = new(time.Time)
	*script.UpdatedAt = time.Now()

	if err := p.db.Model(&dto.Script{}).Where("id = ?", script.ID).Updates(script).Error; err != nil {
		http.Error(w, "Error updating script", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(script)
	if err != nil {
		http.Error(w, "Error saving script", http.StatusInternalServerError)
		return
	}
}
