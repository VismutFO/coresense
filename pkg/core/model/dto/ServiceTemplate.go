package dto

import (
	"time"

	"github.com/google/uuid"
)

// ServiceTemplateRequest represents the service_templates table
type ServiceTemplateRequest struct {
	Name         string `gorm:"size:255;not null"`
	Description  string `gorm:"type:text;not null"`
	FieldsFormat []QuestionInput
}

type QuestionInput struct {
	ScriptID    uuid.NullUUID
	Type        string
	Description string
	Number      int `gorm:"not null"`
}

// ServiceTemplateGridRecord represents the service_templates table
type ServiceTemplateGridRecord struct {
	ID                 uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BusinessCustomerID uuid.UUID `gorm:"not null"`
	Name               string    `gorm:"size:255;not null"`
	Description        string    `gorm:"type:text;not null"`
}

// ServiceTemplate represents the service_templates table
type ServiceTemplate struct {
	ID                 uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BusinessCustomerID uuid.UUID `gorm:"not null"`
	Name               string    `gorm:"size:255;not null"`
	Description        string    `gorm:"type:text;not null"`
	FieldsFormat       []Question
	CreatedAt          time.Time `gorm:"not null"`
	UpdatedAt          *time.Time
}
