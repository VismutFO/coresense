package dto

import (
	"time"

	"github.com/google/uuid"
)

type QuestionAnswered struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	QuestionID      uuid.UUID `gorm:"type:uuid;not null"`
	FilledServiceID uuid.UUID `gorm:"type:uuid;not null"`
	Answer          string
	Number          int       `gorm:"not null"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       *time.Time
}

type QuestionAnsweredGridRecord struct {
	Question string
	Answer   string
}

// TableName overrides the table name used by QuestionAnswered to `questions_answered`
func (QuestionAnswered) TableName() string {
	return "questions_answered"
}
