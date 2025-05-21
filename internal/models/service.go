package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ServiceCategory represents a category for services
type ServiceCategory struct {
	BaseModel
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
}

// TableName overrides the table name
func (ServiceCategory) TableName() string {
	return "service_categories"
}

// Service represents a beauty service offered by a provider
type Service struct {
	BaseModel
	BusinessID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"business_id"`
	Business    Business       `gorm:"foreignKey:BusinessID" json:"business"`
	CategoryID  *uuid.UUID     `gorm:"type:uuid;index" json:"category_id"`
	Category    ServiceCategory `gorm:"foreignKey:CategoryID" json:"category"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Duration    int            `gorm:"not null" json:"duration"` // in minutes
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	ImageURL    string         `json:"image_url"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
	Tags        ServiceTags    `gorm:"type:jsonb" json:"tags"`
	Settings    ServiceSettings `gorm:"type:jsonb" json:"settings"`
}

// TableName overrides the table name
func (Service) TableName() string {
	return "services"
}

// ServiceTags is a string array that can be stored as JSONB
type ServiceTags []string

// Scan implements the sql.Scanner interface for ServiceTags
func (st *ServiceTags) Scan(value interface{}) error {
	if value == nil {
		*st = ServiceTags{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp ServiceTags
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*st = temp

	return nil
}

// Value implements the driver.Valuer interface for ServiceTags
func (st ServiceTags) Value() (driver.Value, error) {
	return json.Marshal(st)
}

// ServiceSettings contains configurable settings for a service
type ServiceSettings struct {
	AllowOnlineBooking    bool  `json:"allow_online_booking"`
	MinAdvanceTimeHours   int   `json:"min_advance_time_hours"`
	MaxAdvanceTimeDays    int   `json:"max_advance_time_days"`
	RequireDeposit        bool  `json:"require_deposit"`
	DepositAmount         float64 `json:"deposit_amount"`
	DepositPercentage     int   `json:"deposit_percentage"`
	CanBeBooked           bool  `json:"can_be_booked"`
	RequiresConsultation  bool  `json:"requires_consultation"`
	BufferTimeBeforeMin   int   `json:"buffer_time_before_min"`
	BufferTimeAfterMin    int   `json:"buffer_time_after_min"`
	CancellationPolicyHours int  `json:"cancellation_policy_hours"`
}

// Scan implements the sql.Scanner interface for ServiceSettings
func (ss *ServiceSettings) Scan(value interface{}) error {
	if value == nil {
		*ss = ServiceSettings{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp ServiceSettings
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*ss = temp

	return nil
}

// Value implements the driver.Valuer interface for ServiceSettings
func (ss ServiceSettings) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

// ServiceVariant represents a variant of a service with different options
type ServiceVariant struct {
	BaseModel
	ServiceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"service_id"`
	Service     Service   `gorm:"foreignKey:ServiceID" json:"service"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Duration    int       `gorm:"not null" json:"duration"` // in minutes
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
}

// TableName overrides the table name
func (ServiceVariant) TableName() string {
	return "service_variants"
}

// ServiceOption represents a customizable option for a service
type ServiceOption struct {
	BaseModel
	ServiceID    uuid.UUID `gorm:"type:uuid;not null;index" json:"service_id"`
	Service      Service   `gorm:"foreignKey:ServiceID" json:"service"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	IsRequired   bool      `gorm:"not null;default:false" json:"is_required"`
	IsMultiple   bool      `gorm:"not null;default:false" json:"is_multiple"`
	MinSelections int       `gorm:"not null;default:0" json:"min_selections"`
	MaxSelections int       `gorm:"not null;default:1" json:"max_selections"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
}

// TableName overrides the table name
func (ServiceOption) TableName() string {
	return "service_options"
}

// ServiceOptionChoice represents a choice within a service option
type ServiceOptionChoice struct {
	BaseModel
	OptionID    uuid.UUID `gorm:"type:uuid;not null;index" json:"option_id"`
	Option      ServiceOption `gorm:"foreignKey:OptionID" json:"option"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	PriceAdjustment float64 `gorm:"type:decimal(10,2);not null;default:0" json:"price_adjustment"`
	TimeAdjustment int     `gorm:"not null;default:0" json:"time_adjustment"` // in minutes
	IsDefault     bool     `gorm:"not null;default:false" json:"is_default"`
	IsActive      bool     `gorm:"not null;default:true" json:"is_active"`
}

// TableName overrides the table name
func (ServiceOptionChoice) TableName() string {
	return "service_option_choices"
}

// ServiceBundle represents a bundle of services offered together
type ServiceBundle struct {
	BaseModel
	BusinessID  uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
	Business    Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	DiscountPercentage int `gorm:"not null;default:0" json:"discount_percentage"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

// TableName overrides the table name
func (ServiceBundle) TableName() string {
	return "service_bundles"
}

// ServiceBundleItem represents a service included in a bundle
type ServiceBundleItem struct {
	BaseModel
	BundleID   uuid.UUID `gorm:"type:uuid;not null;index" json:"bundle_id"`
	Bundle     ServiceBundle `gorm:"foreignKey:BundleID" json:"bundle"`
	ServiceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"service_id"`
	Service    Service   `gorm:"foreignKey:ServiceID" json:"service"`
	Quantity   int       `gorm:"not null;default:1" json:"quantity"`
	IsRequired bool      `gorm:"not null;default:true" json:"is_required"`
}

// TableName overrides the table name
func (ServiceBundleItem) TableName() string {
	return "service_bundle_items"
}