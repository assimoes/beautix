package models

import (
	"time"

	"github.com/google/uuid"
)

// Client represents a client of a beauty business
type Client struct {
	BaseModel
	BusinessID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"business_id"`
	Business         Business   `gorm:"foreignKey:BusinessID" json:"business"`
	UserID           *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	User             *User      `gorm:"foreignKey:UserID" json:"user"`
	FirstName        string     `gorm:"not null" json:"first_name"`
	LastName         string     `gorm:"not null" json:"last_name"`
	Email            string     `gorm:"index" json:"email"`
	Phone            string     `gorm:"index" json:"phone"`
	DateOfBirth      *time.Time `gorm:"type:date" json:"date_of_birth"`
	AddressLine1     string     `gorm:"column:address_line1" json:"address_line1"`
	City             string     `json:"city"`
	PostalCode       string     `json:"postal_code"`
	Country          string     `json:"country"`
	Notes            string     `json:"notes"`
	Allergies        string     `json:"allergies"`
	HealthConditions string     `json:"health_conditions"`
	ReferralSource   string     `gorm:"column:referral_source" json:"referral_source"`
	AcceptsMarketing bool       `gorm:"not null;default:false" json:"accepts_marketing"`
	IsActive         bool       `gorm:"not null;default:true" json:"is_active"`
}

// TableName overrides the table name
func (Client) TableName() string {
	return "clients"
}

// TODO: The following types are defined but not yet implemented in the database schema:
// - ClientNote and client_notes table
// - ClientDocument and client_documents table
// - JSONB types like ClientTags, ClientPreferences, HealthInfo etc.
// These should be added via migrations when needed.
