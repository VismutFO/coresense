package dto

import (
	"github.com/google/uuid"
	"time"
)

// User represents the users table
type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username  string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt *time.Time
}
