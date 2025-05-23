package models

import (
	"github.com/google/uuid"
	"time"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	// AppointmentStatusScheduled represents a scheduled appointment
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	// AppointmentStatusConfirmed represents a confirmed appointment
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	// AppointmentStatusInProgress represents an appointment in progress
	AppointmentStatusInProgress AppointmentStatus = "in_progress"
	// AppointmentStatusCompleted represents a completed appointment
	AppointmentStatusCompleted AppointmentStatus = "completed"
	// AppointmentStatusCancelled represents a cancelled appointment
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
	// AppointmentStatusNoShow represents a no-show appointment
	AppointmentStatusNoShow AppointmentStatus = "no_show"
)

// PaymentStatus represents the payment status of an appointment
type PaymentStatus string

const (
	// PaymentStatusPending represents a pending payment
	PaymentStatusPending PaymentStatus = "pending"
	// PaymentStatusPaid represents a completed payment
	PaymentStatusPaid PaymentStatus = "paid"
	// PaymentStatusPartial represents a partially paid payment
	PaymentStatusPartial PaymentStatus = "partial"
	// PaymentStatusRefunded represents a refunded payment
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// PaymentMethod represents the method of payment
type PaymentMethod string

const (
	// PaymentMethodCash represents cash payment
	PaymentMethodCash PaymentMethod = "cash"
	// PaymentMethodCard represents card payment
	PaymentMethodCard PaymentMethod = "card"
	// PaymentMethodTransfer represents bank transfer payment
	PaymentMethodTransfer PaymentMethod = "transfer"
	// PaymentMethodOther represents other payment methods
	PaymentMethodOther PaymentMethod = "other"
)

// Appointment represents a simplified booking for a service that matches the database schema
type Appointment struct {
	BaseModel
	BusinessID         uuid.UUID         `gorm:"type:uuid;not null;index" json:"business_id"`
	Business           *Business         `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	ClientID           uuid.UUID         `gorm:"type:uuid;not null;index" json:"client_id"`
	Client             *Client           `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	StaffID            uuid.UUID         `gorm:"type:uuid;not null;index" json:"staff_id"`
	Staff              *Staff            `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	ServiceID          uuid.UUID         `gorm:"type:uuid;not null;index" json:"service_id"`
	Service            *Service          `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	StartTime          time.Time         `gorm:"not null;index" json:"start_time"`
	EndTime            time.Time         `gorm:"not null" json:"end_time"`
	Status             AppointmentStatus `gorm:"type:varchar(20);not null;default:'scheduled'" json:"status"`
	Notes              string            `gorm:"type:text" json:"notes"`
	EstimatedPrice     *float64          `gorm:"type:decimal(10,2)" json:"estimated_price"`
	ActualPrice        *float64          `gorm:"type:decimal(10,2)" json:"actual_price"`
	PaymentMethod      *PaymentMethod    `gorm:"type:varchar(20)" json:"payment_method"`
	PaymentStatus      PaymentStatus     `gorm:"type:varchar(20);default:'pending'" json:"payment_status"`
	ClientConfirmed    bool              `gorm:"default:false" json:"client_confirmed"`
	StaffConfirmed     bool              `gorm:"default:false" json:"staff_confirmed"`
	CancellationReason string            `gorm:"type:text" json:"cancellation_reason"`
}

// TableName overrides the table name
func (Appointment) TableName() string {
	return "appointments"
}

// IsCompleted returns true if the appointment is completed
func (a *Appointment) IsCompleted() bool {
	return a.Status == AppointmentStatusCompleted
}

// IsCancelled returns true if the appointment is cancelled
func (a *Appointment) IsCancelled() bool {
	return a.Status == AppointmentStatusCancelled
}

// IsConfirmed returns true if both client and staff have confirmed
func (a *Appointment) IsConfirmed() bool {
	return a.ClientConfirmed && a.StaffConfirmed
}

// Duration returns the duration of the appointment in minutes
func (a *Appointment) Duration() int {
	return int(a.EndTime.Sub(a.StartTime).Minutes())
}
