package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Staff represents a staff member in a business
type Staff struct {
	StaffID         uuid.UUID  `json:"staff_id"`
	BusinessID      uuid.UUID  `json:"business_id"`
	UserID          uuid.UUID  `json:"user_id"`
	Position        string     `json:"position"`
	Bio             string     `json:"bio,omitempty"`
	SpecialtyAreas  []string   `json:"specialty_areas,omitempty"`
	ProfileImageURL string     `json:"profile_image_url,omitempty"`
	WorkingHours    []byte     `json:"working_hours,omitempty"` // JSONB in database
	IsActive        bool       `json:"is_active"`
	EmploymentType  string     `json:"employment_type"`
	JoinDate        time.Time  `json:"join_date"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	CommissionRate  float64    `json:"commission_rate,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	UpdatedBy       *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	DeletedBy       *uuid.UUID `json:"deleted_by,omitempty"`

	// Expanded relationships (populated by service when needed)
	User     *User     `json:"user,omitempty"`
	Business *Business `json:"business,omitempty"`
}

// CreateStaffInput is the input for creating a staff member
type CreateStaffInput struct {
	BusinessID      uuid.UUID  `json:"business_id" validate:"required"`
	UserID          uuid.UUID  `json:"user_id" validate:"required"`
	Position        string     `json:"position" validate:"required"`
	Bio             string     `json:"bio,omitempty"`
	SpecialtyAreas  []string   `json:"specialty_areas,omitempty"`
	ProfileImageURL string     `json:"profile_image_url,omitempty"`
	WorkingHours    []byte     `json:"working_hours,omitempty"`
	IsActive        *bool      `json:"is_active,omitempty"`
	EmploymentType  string     `json:"employment_type" validate:"required"`
	JoinDate        time.Time  `json:"join_date" validate:"required"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	CommissionRate  *float64   `json:"commission_rate,omitempty"`
}

// UpdateStaffInput is the input for updating a staff member
type UpdateStaffInput struct {
	Position        *string    `json:"position,omitempty"`
	Bio             *string    `json:"bio,omitempty"`
	SpecialtyAreas  *[]string  `json:"specialty_areas,omitempty"`
	ProfileImageURL *string    `json:"profile_image_url,omitempty"`
	WorkingHours    *[]byte    `json:"working_hours,omitempty"`
	IsActive        *bool      `json:"is_active,omitempty"`
	EmploymentType  *string    `json:"employment_type,omitempty"`
	JoinDate        *time.Time `json:"join_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	CommissionRate  *float64   `json:"commission_rate,omitempty"`
}

// StaffRepository defines the methods to interact with the staff data store
type StaffRepository interface {
	Create(ctx context.Context, staff *Staff) error
	GetByID(ctx context.Context, id uuid.UUID) (*Staff, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Staff, error)
	GetByBusinessID(ctx context.Context, businessID uuid.UUID) ([]*Staff, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateStaffInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	List(ctx context.Context, page, pageSize int) ([]*Staff, error)
	ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Staff, error)
	Search(ctx context.Context, businessID uuid.UUID, query string, page, pageSize int) ([]*Staff, error)
	Count(ctx context.Context) (int64, error)
	CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// StaffService defines the business logic for staff operations
type StaffService interface {
	CreateStaff(ctx context.Context, input *CreateStaffInput) (*Staff, error)
	GetStaff(ctx context.Context, id uuid.UUID) (*Staff, error)
	GetStaffByUser(ctx context.Context, userID uuid.UUID) ([]*Staff, error)
	GetBusinessStaff(ctx context.Context, businessID uuid.UUID) ([]*Staff, error)
	UpdateStaff(ctx context.Context, id uuid.UUID, input *UpdateStaffInput, updatedBy uuid.UUID) error
	DeleteStaff(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListStaff(ctx context.Context, page, pageSize int) ([]*Staff, error)
	ListBusinessStaff(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Staff, error)
	SearchStaff(ctx context.Context, businessID uuid.UUID, query string, page, pageSize int) ([]*Staff, error)
	CountStaff(ctx context.Context) (int64, error)
	CountBusinessStaff(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// ServiceAssignment represents the assignment of services to staff members
type ServiceAssignment struct {
	AssignmentID uuid.UUID  `json:"assignment_id"`
	BusinessID   uuid.UUID  `json:"business_id"`
	StaffID      uuid.UUID  `json:"staff_id"`
	ServiceID    uuid.UUID  `json:"service_id"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    uuid.UUID  `json:"created_by"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	DeletedBy    *uuid.UUID `json:"deleted_by,omitempty"`

	// Expanded relationships (populated by service when needed)
	Staff   *Staff   `json:"staff,omitempty"`
	Service *Service `json:"service,omitempty"`
}

// CreateServiceAssignmentInput is the input for creating a service assignment
type CreateServiceAssignmentInput struct {
	BusinessID uuid.UUID `json:"business_id" validate:"required"`
	StaffID    uuid.UUID `json:"staff_id" validate:"required"`
	ServiceID  uuid.UUID `json:"service_id" validate:"required"`
	IsActive   *bool     `json:"is_active,omitempty"`
}

// UpdateServiceAssignmentInput is the input for updating a service assignment
type UpdateServiceAssignmentInput struct {
	IsActive *bool `json:"is_active,omitempty"`
}

// ServiceAssignmentRepository defines the methods to interact with the service assignment data store
type ServiceAssignmentRepository interface {
	Create(ctx context.Context, assignment *ServiceAssignment) error
	GetByID(ctx context.Context, id uuid.UUID) (*ServiceAssignment, error)
	GetByStaffAndService(ctx context.Context, staffID, serviceID uuid.UUID) (*ServiceAssignment, error)
	GetByStaff(ctx context.Context, staffID uuid.UUID) ([]*ServiceAssignment, error)
	GetByService(ctx context.Context, serviceID uuid.UUID) ([]*ServiceAssignment, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateServiceAssignmentInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*ServiceAssignment, error)
	CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// ServiceAssignmentService defines the business logic for service assignment operations
type ServiceAssignmentService interface {
	CreateServiceAssignment(ctx context.Context, input *CreateServiceAssignmentInput) (*ServiceAssignment, error)
	GetServiceAssignment(ctx context.Context, id uuid.UUID) (*ServiceAssignment, error)
	GetServiceAssignmentByStaffAndService(ctx context.Context, staffID, serviceID uuid.UUID) (*ServiceAssignment, error)
	GetStaffServiceAssignments(ctx context.Context, staffID uuid.UUID) ([]*ServiceAssignment, error)
	GetServiceStaffAssignments(ctx context.Context, serviceID uuid.UUID) ([]*ServiceAssignment, error)
	UpdateServiceAssignment(ctx context.Context, id uuid.UUID, input *UpdateServiceAssignmentInput, updatedBy uuid.UUID) error
	DeleteServiceAssignment(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListBusinessServiceAssignments(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*ServiceAssignment, error)
	CountBusinessServiceAssignments(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// AvailabilityException represents exceptions to a staff member's regular working hours
type AvailabilityException struct {
	ExceptionID    uuid.UUID  `json:"exception_id"`
	BusinessID     uuid.UUID  `json:"business_id"`
	StaffID        uuid.UUID  `json:"staff_id"`
	ExceptionType  string     `json:"exception_type"` // "time_off", "holiday", "custom_hours"
	StartTime      time.Time  `json:"start_time"`
	EndTime        time.Time  `json:"end_time"`
	IsFullDay      bool       `json:"is_full_day"`
	IsRecurring    bool       `json:"is_recurring"`
	RecurrenceRule string     `json:"recurrence_rule,omitempty"` // iCalendar RRULE format
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	DeletedBy      *uuid.UUID `json:"deleted_by,omitempty"`

	// Expanded relationships (populated by service when needed)
	Staff *Staff `json:"staff,omitempty"`
}

// CreateAvailabilityExceptionInput is the input for creating an availability exception
type CreateAvailabilityExceptionInput struct {
	BusinessID     uuid.UUID `json:"business_id" validate:"required"`
	StaffID        uuid.UUID `json:"staff_id" validate:"required"`
	ExceptionType  string    `json:"exception_type" validate:"required"`
	StartTime      time.Time `json:"start_time" validate:"required"`
	EndTime        time.Time `json:"end_time" validate:"required"`
	IsFullDay      *bool     `json:"is_full_day,omitempty"`
	IsRecurring    *bool     `json:"is_recurring,omitempty"`
	RecurrenceRule string    `json:"recurrence_rule,omitempty"`
	Notes          string    `json:"notes,omitempty"`
}

// UpdateAvailabilityExceptionInput is the input for updating an availability exception
type UpdateAvailabilityExceptionInput struct {
	ExceptionType  *string    `json:"exception_type,omitempty"`
	StartTime      *time.Time `json:"start_time,omitempty"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	IsFullDay      *bool      `json:"is_full_day,omitempty"`
	IsRecurring    *bool      `json:"is_recurring,omitempty"`
	RecurrenceRule *string    `json:"recurrence_rule,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
}

// AvailabilityExceptionRepository defines the methods to interact with the availability exception data store
type AvailabilityExceptionRepository interface {
	Create(ctx context.Context, exception *AvailabilityException) error
	GetByID(ctx context.Context, id uuid.UUID) (*AvailabilityException, error)
	GetByStaff(ctx context.Context, staffID uuid.UUID) ([]*AvailabilityException, error)
	GetByStaffAndDateRange(ctx context.Context, staffID uuid.UUID, start, end time.Time) ([]*AvailabilityException, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateAvailabilityExceptionInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*AvailabilityException, error)
	ListByBusinessAndDateRange(ctx context.Context, businessID uuid.UUID, start, end time.Time, page, pageSize int) ([]*AvailabilityException, error)
	CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// AvailabilityExceptionService defines the business logic for availability exception operations
type AvailabilityExceptionService interface {
	CreateAvailabilityException(ctx context.Context, input *CreateAvailabilityExceptionInput) (*AvailabilityException, error)
	GetAvailabilityException(ctx context.Context, id uuid.UUID) (*AvailabilityException, error)
	GetStaffAvailabilityExceptions(ctx context.Context, staffID uuid.UUID) ([]*AvailabilityException, error)
	GetStaffAvailabilityExceptionsByDateRange(ctx context.Context, staffID uuid.UUID, start, end time.Time) ([]*AvailabilityException, error)
	UpdateAvailabilityException(ctx context.Context, id uuid.UUID, input *UpdateAvailabilityExceptionInput, updatedBy uuid.UUID) error
	DeleteAvailabilityException(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListBusinessAvailabilityExceptions(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*AvailabilityException, error)
	ListBusinessAvailabilityExceptionsByDateRange(ctx context.Context, businessID uuid.UUID, start, end time.Time, page, pageSize int) ([]*AvailabilityException, error)
	CountBusinessAvailabilityExceptions(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// StaffPerformance represents performance metrics for a staff member
type StaffPerformance struct {
	PerformanceID         uuid.UUID `json:"performance_id"`
	BusinessID            uuid.UUID `json:"business_id"`
	StaffID               uuid.UUID `json:"staff_id"`
	Period                string    `json:"period"` // "daily", "weekly", "monthly", "yearly"
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
	TotalAppointments     int       `json:"total_appointments"`
	CompletedAppointments int       `json:"completed_appointments"`
	CanceledAppointments  int       `json:"canceled_appointments"`
	NoShowAppointments    int       `json:"no_show_appointments"`
	TotalRevenue          float64   `json:"total_revenue"`
	AverageRating         float64   `json:"average_rating,omitempty"`
	ClientRetentionRate   float64   `json:"client_retention_rate,omitempty"`
	NewClients            int       `json:"new_clients"`
	ReturnClients         int       `json:"return_clients"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	// Expanded relationships (populated by service when needed)
	Staff *Staff `json:"staff,omitempty"`
}

// StaffPerformanceRepository defines the methods to interact with the staff performance data store
type StaffPerformanceRepository interface {
	Create(ctx context.Context, performance *StaffPerformance) error
	GetByID(ctx context.Context, id uuid.UUID) (*StaffPerformance, error)
	GetByStaffAndPeriod(ctx context.Context, staffID uuid.UUID, period string, startDate time.Time) (*StaffPerformance, error)
	GetByStaffAndDateRange(ctx context.Context, staffID uuid.UUID, startDate, endDate time.Time) ([]*StaffPerformance, error)
	Update(ctx context.Context, id uuid.UUID, performance *StaffPerformance) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByBusiness(ctx context.Context, businessID uuid.UUID, period string, page, pageSize int) ([]*StaffPerformance, error)
}

// StaffPerformanceService defines the business logic for staff performance operations
type StaffPerformanceService interface {
	CreateStaffPerformance(ctx context.Context, performance *StaffPerformance) error
	GetStaffPerformance(ctx context.Context, id uuid.UUID) (*StaffPerformance, error)
	GetStaffPerformanceByPeriod(ctx context.Context, staffID uuid.UUID, period string, startDate time.Time) (*StaffPerformance, error)
	GetStaffPerformanceByDateRange(ctx context.Context, staffID uuid.UUID, startDate, endDate time.Time) ([]*StaffPerformance, error)
	UpdateStaffPerformance(ctx context.Context, id uuid.UUID, performance *StaffPerformance) error
	DeleteStaffPerformance(ctx context.Context, id uuid.UUID) error
	ListBusinessStaffPerformance(ctx context.Context, businessID uuid.UUID, period string, page, pageSize int) ([]*StaffPerformance, error)
	CalculateStaffPerformanceMetrics(ctx context.Context, staffID uuid.UUID, startDate, endDate time.Time) (*StaffPerformance, error)
}
