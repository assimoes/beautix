package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Business represents a business entity, which serves as a tenant
type Business struct {
	BusinessID       uuid.UUID  `json:"business_id"`
	OwnerID          uuid.UUID  `json:"owner_id"`
	BusinessName     string     `json:"business_name"`
	BusinessType     string     `json:"business_type"`
	TaxID            string     `json:"tax_id,omitempty"`
	Phone            string     `json:"phone"`
	Email            string     `json:"email"`
	AddressLine1     string     `json:"address_line1,omitempty"`
	City             string     `json:"city,omitempty"`
	Region           string     `json:"region,omitempty"`
	PostalCode       string     `json:"postal_code,omitempty"`
	Country          string     `json:"country"`
	TimeZone         string     `json:"time_zone"`
	BusinessHours    []byte     `json:"business_hours,omitempty"` // JSONB in database
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	IsActive         bool       `json:"is_active"`
	SubscriptionPlan string     `json:"subscription_plan"`
	TrialEndsAt      *time.Time `json:"trial_ends_at,omitempty"`

	// Expanded relationships (populated by service when needed)
	Owner *User `json:"owner,omitempty"`
}

// CreateBusinessInput is the input for creating a business
type CreateBusinessInput struct {
	OwnerID          uuid.UUID `json:"owner_id" validate:"required"`
	BusinessName     string    `json:"business_name" validate:"required"`
	BusinessType     string    `json:"business_type" validate:"required"`
	TaxID            string    `json:"tax_id,omitempty"`
	Phone            string    `json:"phone" validate:"required"`
	Email            string    `json:"email" validate:"required,email"`
	AddressLine1     string    `json:"address_line1,omitempty"`
	City             string    `json:"city,omitempty"`
	Region           string    `json:"region,omitempty"`
	PostalCode       string    `json:"postal_code,omitempty"`
	Country          string    `json:"country,omitempty"`
	TimeZone         string    `json:"time_zone,omitempty"`
	BusinessHours    []byte    `json:"business_hours,omitempty"`
	SubscriptionPlan string    `json:"subscription_plan,omitempty"`
}

// UpdateBusinessInput is the input for updating a business
type UpdateBusinessInput struct {
	BusinessName     *string    `json:"business_name,omitempty"`
	BusinessType     *string    `json:"business_type,omitempty"`
	TaxID            *string    `json:"tax_id,omitempty"`
	Phone            *string    `json:"phone,omitempty"`
	Email            *string    `json:"email,omitempty" validate:"omitempty,email"`
	AddressLine1     *string    `json:"address_line1,omitempty"`
	City             *string    `json:"city,omitempty"`
	Region           *string    `json:"region,omitempty"`
	PostalCode       *string    `json:"postal_code,omitempty"`
	Country          *string    `json:"country,omitempty"`
	TimeZone         *string    `json:"time_zone,omitempty"`
	BusinessHours    *[]byte    `json:"business_hours,omitempty"`
	IsActive         *bool      `json:"is_active,omitempty"`
	SubscriptionPlan *string    `json:"subscription_plan,omitempty"`
	TrialEndsAt      *time.Time `json:"trial_ends_at,omitempty"`
}

// BusinessRepository defines methods for business data store
type BusinessRepository interface {
	Create(ctx context.Context, business *Business) error
	GetByID(ctx context.Context, id uuid.UUID) (*Business, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*Business, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateBusinessInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int) ([]*Business, error)
	Count(ctx context.Context) (int64, error)
}

// BusinessService defines business logic for business operations
type BusinessService interface {
	CreateBusiness(ctx context.Context, input *CreateBusinessInput) (*Business, error)
	GetBusiness(ctx context.Context, id uuid.UUID) (*Business, error)
	GetBusinessesByOwner(ctx context.Context, ownerID uuid.UUID) ([]*Business, error)
	UpdateBusiness(ctx context.Context, id uuid.UUID, input *UpdateBusinessInput) error
	DeleteBusiness(ctx context.Context, id uuid.UUID) error
	ListBusinesses(ctx context.Context, page, pageSize int) ([]*Business, error)
	CountBusinesses(ctx context.Context) (int64, error)
}

// Subscription represents a business subscription
type Subscription struct {
	SubscriptionID uuid.UUID  `json:"subscription_id"`
	BusinessID     uuid.UUID  `json:"business_id"`
	PlanType       string     `json:"plan_type"`
	Status         string     `json:"status"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	TrialEndDate   *time.Time `json:"trial_end_date,omitempty"`
	BillingCycle   string     `json:"billing_cycle"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Expanded relationships
	Business *Business `json:"business,omitempty"`
}

// CreateSubscriptionInput is the input for creating a subscription
type CreateSubscriptionInput struct {
	BusinessID   uuid.UUID  `json:"business_id" validate:"required"`
	PlanType     string     `json:"plan_type" validate:"required"`
	Status       string     `json:"status,omitempty"`
	StartDate    time.Time  `json:"start_date" validate:"required"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	TrialEndDate *time.Time `json:"trial_end_date,omitempty"`
	BillingCycle string     `json:"billing_cycle,omitempty"`
}

// UpdateSubscriptionInput is the input for updating a subscription
type UpdateSubscriptionInput struct {
	PlanType     *string    `json:"plan_type,omitempty"`
	Status       *string    `json:"status,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	TrialEndDate *time.Time `json:"trial_end_date,omitempty"`
	BillingCycle *string    `json:"billing_cycle,omitempty"`
}

// SubscriptionRepository defines methods for subscription data store
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	GetByBusinessID(ctx context.Context, businessID uuid.UUID) (*Subscription, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateSubscriptionInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int) ([]*Subscription, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
}

// SubscriptionService defines business logic for subscription operations
type SubscriptionService interface {
	CreateSubscription(ctx context.Context, input *CreateSubscriptionInput) (*Subscription, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (*Subscription, error)
	GetBusinessSubscription(ctx context.Context, businessID uuid.UUID) (*Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, input *UpdateSubscriptionInput) error
	CancelSubscription(ctx context.Context, id uuid.UUID) error
	ListSubscriptions(ctx context.Context, page, pageSize int) ([]*Subscription, error)
	CountSubscriptionsByStatus(ctx context.Context, status string) (int64, error)
}

// BusinessSettings represents settings for a business
type BusinessSettings struct {
	SettingID                 uuid.UUID `json:"setting_id"`
	BusinessID                uuid.UUID `json:"business_id"`
	CalendarStartHour         int       `json:"calendar_start_hour"`
	CalendarEndHour           int       `json:"calendar_end_hour"`
	AppointmentBufferMinutes  int       `json:"appointment_buffer_minutes"`
	AllowOnlineBooking        bool      `json:"allow_online_booking"`
	DefaultAppointmentDuration int       `json:"default_appointment_duration"`
	Currency                  string    `json:"currency"`
	DateFormat                string    `json:"date_format"`
	TimeFormat                string    `json:"time_format"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`

	// Expanded relationships
	Business *Business `json:"business,omitempty"`
}

// CreateBusinessSettingsInput is the input for creating business settings
type CreateBusinessSettingsInput struct {
	BusinessID                uuid.UUID `json:"business_id" validate:"required"`
	CalendarStartHour         *int      `json:"calendar_start_hour,omitempty"`
	CalendarEndHour           *int      `json:"calendar_end_hour,omitempty"`
	AppointmentBufferMinutes  *int      `json:"appointment_buffer_minutes,omitempty"`
	AllowOnlineBooking        *bool     `json:"allow_online_booking,omitempty"`
	DefaultAppointmentDuration *int      `json:"default_appointment_duration,omitempty"`
	Currency                  *string   `json:"currency,omitempty"`
	DateFormat                *string   `json:"date_format,omitempty"`
	TimeFormat                *string   `json:"time_format,omitempty"`
}

// UpdateBusinessSettingsInput is the input for updating business settings
type UpdateBusinessSettingsInput struct {
	CalendarStartHour         *int    `json:"calendar_start_hour,omitempty"`
	CalendarEndHour           *int    `json:"calendar_end_hour,omitempty"`
	AppointmentBufferMinutes  *int    `json:"appointment_buffer_minutes,omitempty"`
	AllowOnlineBooking        *bool   `json:"allow_online_booking,omitempty"`
	DefaultAppointmentDuration *int    `json:"default_appointment_duration,omitempty"`
	Currency                  *string `json:"currency,omitempty"`
	DateFormat                *string `json:"date_format,omitempty"`
	TimeFormat                *string `json:"time_format,omitempty"`
}

// BusinessSettingsRepository defines methods for business settings data store
type BusinessSettingsRepository interface {
	Create(ctx context.Context, settings *BusinessSettings) error
	GetByBusinessID(ctx context.Context, businessID uuid.UUID) (*BusinessSettings, error)
	Update(ctx context.Context, businessID uuid.UUID, input *UpdateBusinessSettingsInput) error
	Delete(ctx context.Context, businessID uuid.UUID) error
}

// BusinessSettingsService defines business logic for business settings operations
type BusinessSettingsService interface {
	CreateBusinessSettings(ctx context.Context, input *CreateBusinessSettingsInput) (*BusinessSettings, error)
	GetBusinessSettings(ctx context.Context, businessID uuid.UUID) (*BusinessSettings, error)
	UpdateBusinessSettings(ctx context.Context, businessID uuid.UUID, input *UpdateBusinessSettingsInput) error
	DeleteBusinessSettings(ctx context.Context, businessID uuid.UUID) error
}

// NotificationSettings represents notification settings for a business
type NotificationSettings struct {
	SettingID        uuid.UUID  `json:"setting_id"`
	BusinessID       uuid.UUID  `json:"business_id"`
	NotificationType string     `json:"notification_type"`
	IsEnabled        bool       `json:"is_enabled"`
	TemplateSubject  string     `json:"template_subject,omitempty"`
	TemplateContent  string     `json:"template_content,omitempty"`
	SendTimeBefore   *int       `json:"send_time_before,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Expanded relationships
	Business *Business `json:"business,omitempty"`
}

// CreateNotificationSettingsInput is the input for creating notification settings
type CreateNotificationSettingsInput struct {
	BusinessID       uuid.UUID `json:"business_id" validate:"required"`
	NotificationType string    `json:"notification_type" validate:"required"`
	IsEnabled        *bool     `json:"is_enabled,omitempty"`
	TemplateSubject  string    `json:"template_subject,omitempty"`
	TemplateContent  string    `json:"template_content,omitempty"`
	SendTimeBefore   *int      `json:"send_time_before,omitempty"`
}

// UpdateNotificationSettingsInput is the input for updating notification settings
type UpdateNotificationSettingsInput struct {
	IsEnabled       *bool   `json:"is_enabled,omitempty"`
	TemplateSubject *string `json:"template_subject,omitempty"`
	TemplateContent *string `json:"template_content,omitempty"`
	SendTimeBefore  *int    `json:"send_time_before,omitempty"`
}

// NotificationSettingsRepository defines methods for notification settings data store
type NotificationSettingsRepository interface {
	Create(ctx context.Context, settings *NotificationSettings) error
	GetByID(ctx context.Context, id uuid.UUID) (*NotificationSettings, error)
	GetByBusinessAndType(ctx context.Context, businessID uuid.UUID, notificationType string) (*NotificationSettings, error)
	ListByBusiness(ctx context.Context, businessID uuid.UUID) ([]*NotificationSettings, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateNotificationSettingsInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// NotificationSettingsService defines business logic for notification settings operations
type NotificationSettingsService interface {
	CreateNotificationSettings(ctx context.Context, input *CreateNotificationSettingsInput) (*NotificationSettings, error)
	GetNotificationSettings(ctx context.Context, id uuid.UUID) (*NotificationSettings, error)
	GetBusinessNotificationSettings(ctx context.Context, businessID uuid.UUID, notificationType string) (*NotificationSettings, error)
	ListBusinessNotificationSettings(ctx context.Context, businessID uuid.UUID) ([]*NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, id uuid.UUID, input *UpdateNotificationSettingsInput) error
	DeleteNotificationSettings(ctx context.Context, id uuid.UUID) error
}