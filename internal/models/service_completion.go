package models

import (
	"time"

	"github.com/google/uuid"
)

// ServiceCompletion represents a completed service with financial tracking
type ServiceCompletion struct {
	BaseModel
	AppointmentID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"appointment_id"`
	PriceCharged      float64    `gorm:"type:decimal(10,2);not null" json:"price_charged"`
	PaymentMethod     string     `gorm:"type:varchar(50);not null" json:"payment_method"`
	ProviderConfirmed bool       `gorm:"not null;default:false" json:"provider_confirmed"`
	ClientConfirmed   bool       `gorm:"not null;default:false" json:"client_confirmed"`
	CompletionDate    *time.Time `gorm:"index" json:"completion_date,omitempty"`

	// Relationships
	Appointment *Appointment `gorm:"foreignKey:AppointmentID;references:ID" json:"appointment,omitempty"`
}

// TableName returns the table name for the ServiceCompletion model
func (ServiceCompletion) TableName() string {
	return "service_completions"
}
