package domain

import (
	"context"
	"time"
	"github.com/shopspring/decimal"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusInProgress AppointmentStatus = "in_progress"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
	AppointmentStatusNoShow    AppointmentStatus = "no_show"
	AppointmentStatusRescheduled AppointmentStatus = "rescheduled"
)

// Appointment represents a scheduled appointment
type Appointment struct {
	BaseModel
	BusinessID      string            `gorm:"not null;type:uuid;index" json:"business_id"`
	ClientID        string            `gorm:"not null;type:uuid;index" json:"client_id"`
	StaffID         string            `gorm:"not null;type:uuid;index" json:"staff_id"`
	StartTime       time.Time         `gorm:"not null;index" json:"start_time"`
	EndTime         time.Time         `gorm:"not null;index" json:"end_time"`
	Status          AppointmentStatus `gorm:"not null;size:20;default:'scheduled';check:status IN ('scheduled','confirmed','in_progress','completed','cancelled','no_show','rescheduled')" json:"status"`
	Title           *string           `gorm:"size:200" json:"title,omitempty"`
	Notes           *string           `gorm:"type:text" json:"notes,omitempty"`
	InternalNotes   *string           `gorm:"type:text" json:"internal_notes,omitempty"`
	CancellationReason *string        `gorm:"type:text" json:"cancellation_reason,omitempty"`
	TotalPrice      decimal.Decimal   `gorm:"type:decimal(10,2);not null;default:0" json:"total_price"`
	DepositPaid     decimal.Decimal   `gorm:"type:decimal(10,2);not null;default:0" json:"deposit_paid"`
	ReminderSent    bool              `gorm:"not null;default:false" json:"reminder_sent"`
	ConfirmedAt     *time.Time        `gorm:"" json:"confirmed_at,omitempty"`
	CompletedAt     *time.Time        `gorm:"" json:"completed_at,omitempty"`
	CancelledAt     *time.Time        `gorm:"" json:"cancelled_at,omitempty"`

	// Relationships
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
	Client   Client   `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE" json:"client"`
	Staff    Staff    `gorm:"foreignKey:StaffID;constraint:OnDelete:CASCADE" json:"staff"`
}

// TableName returns the table name for Appointment
func (Appointment) TableName() string { return "appointments" }

// Validate validates the appointment model
func (a *Appointment) Validate() error {
	if a.BusinessID == "" {
		return ErrValidation
	}
	if a.ClientID == "" {
		return ErrValidation
	}
	if a.StaffID == "" {
		return ErrValidation
	}
	if a.StartTime.IsZero() {
		return ErrValidation
	}
	if a.EndTime.IsZero() {
		return ErrValidation
	}
	if a.EndTime.Before(a.StartTime) {
		return ErrValidation
	}
	if a.TotalPrice.IsNegative() {
		return ErrValidation
	}
	if a.DepositPaid.IsNegative() {
		return ErrValidation
	}
	return nil
}

// GetDuration returns the appointment duration in minutes
func (a *Appointment) GetDuration() int {
	return int(a.EndTime.Sub(a.StartTime).Minutes())
}

// IsUpcoming returns true if the appointment is in the future
func (a *Appointment) IsUpcoming() bool {
	return a.StartTime.After(time.Now())
}

// IsPast returns true if the appointment is in the past
func (a *Appointment) IsPast() bool {
	return a.EndTime.Before(time.Now())
}

// CanBeCancelled returns true if the appointment can be cancelled
func (a *Appointment) CanBeCancelled() bool {
	return a.Status == AppointmentStatusScheduled || a.Status == AppointmentStatusConfirmed
}

// CanBeCompleted returns true if the appointment can be marked as completed
func (a *Appointment) CanBeCompleted() bool {
	return a.Status == AppointmentStatusConfirmed || a.Status == AppointmentStatusInProgress
}

// MarkCompleted marks the appointment as completed
func (a *Appointment) MarkCompleted() {
	a.Status = AppointmentStatusCompleted
	now := time.Now()
	a.CompletedAt = &now
}

// MarkCancelled marks the appointment as cancelled
func (a *Appointment) MarkCancelled(reason string) {
	a.Status = AppointmentStatusCancelled
	a.CancellationReason = &reason
	now := time.Now()
	a.CancelledAt = &now
}

// AppointmentRepository defines the repository interface for Appointment
type AppointmentRepository interface {
	BaseRepository[Appointment]
	FindByBusinessID(ctx context.Context, businessID string, filters AppointmentFilters) ([]*Appointment, error)
	FindByClientID(ctx context.Context, clientID string) ([]*Appointment, error)
	FindByStaffID(ctx context.Context, staffID string, dateRange DateRange) ([]*Appointment, error)
	FindByDateRange(ctx context.Context, staffID string, start, end time.Time) ([]*Appointment, error)
	CheckOverlap(ctx context.Context, staffID string, start, end time.Time, excludeID *string) (bool, error)
	GetUpcomingByStaff(ctx context.Context, staffID string, limit int) ([]*Appointment, error)
	GetDashboardData(ctx context.Context, businessID string, date time.Time) (*DashboardData, error)
	GetCalendarView(ctx context.Context, businessID string, start, end time.Time) ([]*CalendarAppointment, error)
	GetByStatus(ctx context.Context, businessID string, status AppointmentStatus) ([]*Appointment, error)
}

// Helper types for repository methods
type AppointmentFilters struct {
	Status    *AppointmentStatus `json:"status,omitempty"`
	StaffID   *string           `json:"staff_id,omitempty"`
	ServiceID *string           `json:"service_id,omitempty"`
	DateRange *DateRange        `json:"date_range,omitempty"`
}

type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type DashboardData struct {
	// Dashboard data structure to be defined based on requirements
}

type CalendarAppointment struct {
	// Calendar appointment structure to be defined based on requirements
}