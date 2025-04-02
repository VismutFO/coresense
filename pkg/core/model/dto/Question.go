package dto

import (
	"github.com/google/uuid"
	"time"
)

// Question represents the service_templates table
type Question struct {
	ID                uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ServiceTemplateID uuid.UUID `gorm:"type:uuid;not null"`
	ScriptID          uuid.NullUUID
	Type              string
	Description       string
	Number            int       `gorm:"not null"`
	CreatedAt         time.Time `gorm:"not null"`
	UpdatedAt         *time.Time
}
