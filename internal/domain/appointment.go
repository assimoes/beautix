package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Appointment represents a scheduled appointment between a client and provider
type Appointment struct {
	ID          uuid.UUID  `json:"id"`
	ProviderID  uuid.UUID  `json:"provider_id"`
	ClientID    uuid.UUID  `json:"client_id"`
	ServiceID   uuid.UUID  `json:"service_id"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	Status      string     `json:"status"` // scheduled, confirmed, completed, cancelled, no-show
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	DeletedBy   *uuid.UUID `json:"deleted_by,omitempty"`
	
	// Expanded relationships (populated by service when needed)
	Provider *Provider `json:"provider,omitempty"`
	Client   *Client   `json:"client,omitempty"`
	Service  *Service  `json:"service,omitempty"`
}

// ServiceCompletion represents a completed service with financial tracking
type ServiceCompletion struct {
	ID               uuid.UUID  `json:"id"`
	AppointmentID    uuid.UUID  `json:"appointment_id"`
	PriceCharged     float64    `json:"price_charged"`
	PaymentMethod    string     `json:"payment_method"`
	ProviderConfirmed bool       `json:"provider_confirmed"`
	ClientConfirmed  bool       `json:"client_confirmed"`
	CompletionDate   *time.Time `json:"completion_date,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	DeletedBy        *uuid.UUID `json:"deleted_by,omitempty"`
	
	// Expanded relationships (populated by service when needed)
	Appointment *Appointment `json:"appointment,omitempty"`
}

// CreateAppointmentInput is the input for creating an appointment
type CreateAppointmentInput struct {
	ProviderID uuid.UUID `json:"provider_id" validate:"required"`
	ClientID   uuid.UUID `json:"client_id" validate:"required"`
	ServiceID  uuid.UUID `json:"service_id" validate:"required"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Status     string    `json:"status" validate:"required,oneof=scheduled confirmed completed cancelled no-show"`
	Notes      string    `json:"notes"`
}

// UpdateAppointmentInput is the input for updating an appointment
type UpdateAppointmentInput struct {
	ServiceID *uuid.UUID `json:"service_id"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Status    *string    `json:"status" validate:"omitempty,oneof=scheduled confirmed completed cancelled no-show"`
	Notes     *string    `json:"notes"`
}

// CreateServiceCompletionInput is the input for creating a service completion
type CreateServiceCompletionInput struct {
	AppointmentID    uuid.UUID `json:"appointment_id" validate:"required"`
	PriceCharged     float64   `json:"price_charged" validate:"required,min=0"`
	PaymentMethod    string    `json:"payment_method" validate:"required"`
	ProviderConfirmed bool      `json:"provider_confirmed"`
	ClientConfirmed  bool      `json:"client_confirmed"`
	CompletionDate   *time.Time `json:"completion_date"`
}

// UpdateServiceCompletionInput is the input for updating a service completion
type UpdateServiceCompletionInput struct {
	PriceCharged     *float64   `json:"price_charged" validate:"omitempty,min=0"`
	PaymentMethod    *string    `json:"payment_method"`
	ProviderConfirmed *bool      `json:"provider_confirmed"`
	ClientConfirmed  *bool      `json:"client_confirmed"`
	CompletionDate   *time.Time `json:"completion_date"`
}

// AppointmentRepository defines methods for appointment data store
type AppointmentRepository interface {
	Create(ctx context.Context, appointment *Appointment) error
	GetByID(ctx context.Context, id uuid.UUID) (*Appointment, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateAppointmentInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByProvider(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*Appointment, error)
	ListByClient(ctx context.Context, clientID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*Appointment, error)
	Count(ctx context.Context) (int64, error)
	CountByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
	CountByProviderAndDateRange(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time) (int64, error)
}

// ServiceCompletionRepository defines methods for service completion data store
type ServiceCompletionRepository interface {
	Create(ctx context.Context, completion *ServiceCompletion) error
	GetByID(ctx context.Context, id uuid.UUID) (*ServiceCompletion, error)
	GetByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*ServiceCompletion, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateServiceCompletionInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByProvider(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*ServiceCompletion, error)
	Count(ctx context.Context) (int64, error)
	CountByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
	GetProviderRevenue(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time) (float64, error)
}

// AppointmentService defines business logic for appointment operations
type AppointmentService interface {
	CreateAppointment(ctx context.Context, input *CreateAppointmentInput) (*Appointment, error)
	GetAppointment(ctx context.Context, id uuid.UUID) (*Appointment, error)
	UpdateAppointment(ctx context.Context, id uuid.UUID, input *UpdateAppointmentInput, updatedBy uuid.UUID) error
	DeleteAppointment(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListAppointmentsByProvider(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*Appointment, error)
	ListAppointmentsByClient(ctx context.Context, clientID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*Appointment, error)
	CountAppointments(ctx context.Context) (int64, error)
	CountAppointmentsByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
	CountAppointmentsByProviderAndDateRange(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time) (int64, error)
	CheckAvailability(ctx context.Context, providerID uuid.UUID, startTime, endTime time.Time) (bool, error)
}

// ServiceCompletionService defines business logic for service completion operations
type ServiceCompletionService interface {
	CreateServiceCompletion(ctx context.Context, input *CreateServiceCompletionInput) (*ServiceCompletion, error)
	GetServiceCompletion(ctx context.Context, id uuid.UUID) (*ServiceCompletion, error)
	GetServiceCompletionByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*ServiceCompletion, error)
	UpdateServiceCompletion(ctx context.Context, id uuid.UUID, input *UpdateServiceCompletionInput, updatedBy uuid.UUID) error
	DeleteServiceCompletion(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListServiceCompletionsByProvider(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*ServiceCompletion, error)
	CountServiceCompletions(ctx context.Context) (int64, error)
	CountServiceCompletionsByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
	GetProviderRevenue(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time) (float64, error)
}