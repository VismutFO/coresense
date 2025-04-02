package dto

import (
	"time"

	"github.com/google/uuid"
)

// BusinessCustomer represents the business_customers table
type BusinessCustomerGridRecord struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name string    `gorm:"size:255;not null"`
}

// BusinessCustomer represents the business_customers table
type BusinessCustomer struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string    `gorm:"size:255;not null" json:"username"`
	Email     string    `gorm:"size:255;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt *time.Time
}
