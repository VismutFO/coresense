package dto

import (
	"github.com/google/uuid"
	"time"
)

// Script represents the scripts table
type Script struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name       string    `gorm:"size:255;not null"`
	ScriptCode string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  *time.Time
}
