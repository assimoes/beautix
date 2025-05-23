package models

import (
	"time"

	"github.com/google/uuid"
)

// ServiceCategory represents a category for services
type ServiceCategory struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BusinessID   uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
	Business     Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	DisplayOrder int       `gorm:"not null;default:0" json:"display_order"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName overrides the table name
func (ServiceCategory) TableName() string {
	return "service_categories"
}

// Service represents a beauty service offered by a provider
type Service struct {
	BaseModel
	BusinessID      uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
	Business        Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	Duration        int       `gorm:"not null" json:"duration"` // in minutes
	Price           float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Category        string    `json:"category"`
	Color           string    `json:"color"`
	IsActive        bool      `gorm:"not null;default:true" json:"is_active"`
	PreparationTime int       `gorm:"not null;default:0" json:"preparation_time"`
	CleanupTime     int       `gorm:"not null;default:0" json:"cleanup_time"`
}

// TableName overrides the table name
func (Service) TableName() string {
	return "services"
}

// TODO: The following types are not yet implemented in the database schema:
// - ServiceVariant and service_variants table
// - ServiceOption and service_options table
// - ServiceOptionChoice and service_option_choices table
// - ServiceBundle and service_bundles table
// - ServiceBundleItem and service_bundle_items table
// - JSONB types like ServiceTags and ServiceSettings
// These should be added via migrations when needed.

// Temporary placeholders for types referenced in base.go
type ServiceVariant struct{ BaseModel }
type ServiceOption struct{ BaseModel }
type ServiceOptionChoice struct{ BaseModel }
type ServiceBundle struct{ BaseModel }
type ServiceBundleItem struct{ BaseModel }
