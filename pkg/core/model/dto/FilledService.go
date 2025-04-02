package dto

import (
	"time"

	"github.com/google/uuid"
)

// FilledServiceRequest represents the filled_services table
type FilledServiceRequest struct {
	ServiceTemplateID uuid.UUID               `json:"service_template_id"`
	ServiceData       []QuestionAnsweredInput `json:"service_data"`
}

type QuestionAnsweredInput struct {
	QuestionID uuid.UUID `json:"question_id"`
	Answer     string    `json:"answer"`
}

// FilledService represents the filled_services table
type FilledService struct {
	ID                uuid.UUID          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID            uuid.UUID          `gorm:"not null"`
	ServiceTemplateID uuid.UUID          `gorm:"not null"`
	ServiceData       []QuestionAnswered `gorm:"not null"`
	CreatedAt         time.Time          `gorm:"not null"`
	UpdatedAt         *time.Time
}

type FilledServiceGridRecord struct {
	ID                         uuid.UUID
	ServiceTemplateName        string
	ServiceTemplateDescription string
	CreatedAt                  time.Time
}

type FilledServiceWithDetails struct {
	ServiceTemplateName        string
	ServiceTemplateDescription string
	ServiceData                []QuestionAnsweredGridRecord
	CreatedAt                  time.Time
	User                       string `json:",omitempty"`
}
