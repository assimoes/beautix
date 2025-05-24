package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Business represents a business entity in the system
type Business struct {
	BaseModel
	UserID              string  `gorm:"not null;type:uuid;index" json:"user_id"`
	Name                string  `gorm:"not null;size:100;index" json:"name"`
	DisplayName         *string `gorm:"size:100" json:"display_name,omitempty"`
	BusinessType        *string `gorm:"size:50" json:"business_type,omitempty"`
	TaxID               *string `gorm:"size:50" json:"tax_id,omitempty"`
	Email               string  `gorm:"not null;size:255" json:"email"`
	Website             *string `gorm:"size:255" json:"website,omitempty"`
	LogoURL             *string `gorm:"size:255" json:"logo_url,omitempty"`
	CoverPhotoURL       *string `gorm:"size:255" json:"cover_photo_url,omitempty"`
	IsVerified          bool    `gorm:"default:false" json:"is_verified"`
	SocialLinks         *string `gorm:"type:jsonb;default:'{}'" json:"social_links,omitempty"`
	Settings            *string `gorm:"type:jsonb;default:'{}'" json:"settings,omitempty"`
	Currency            string  `gorm:"size:3;default:'EUR'" json:"currency"`
	TimeZone            string  `gorm:"size:50;default:'Europe/Lisbon'" json:"time_zone"`
	BusinessHours       *string `gorm:"type:jsonb" json:"business_hours,omitempty"`
	IsActive            bool    `gorm:"not null;default:true" json:"is_active"`
	SubscriptionTier    string  `gorm:"size:50;default:'free'" json:"subscription_tier"`
	TrialEndsAt         *time.Time `gorm:"" json:"trial_ends_at,omitempty"`

	// Relationships
	User             User               `gorm:"foreignKey:UserID" json:"user"`
	Locations        []BusinessLocation `gorm:"foreignKey:BusinessID" json:"locations"`
	Settings_        *BusinessSettings  `gorm:"foreignKey:BusinessID" json:"business_settings"`
	Services         []Service          `gorm:"foreignKey:BusinessID" json:"services"`
	ServiceCategories []ServiceCategory `gorm:"foreignKey:BusinessID" json:"service_categories"`
	Staff            []Staff            `gorm:"foreignKey:BusinessID" json:"staff"`
	Clients          []Client           `gorm:"foreignKey:BusinessID" json:"clients"`
	Appointments     []Appointment      `gorm:"foreignKey:BusinessID" json:"appointments,omitempty"`
}

// BusinessLocation represents a physical location for a business
type BusinessLocation struct {
	BaseModel
	BusinessID  string  `gorm:"not null;type:uuid;index" json:"business_id"`
	Name        string  `gorm:"not null;size:100" json:"name"`
	Address     *string `gorm:"size:255" json:"address,omitempty"`
	City        *string `gorm:"size:100" json:"city,omitempty"`
	State       *string `gorm:"size:100" json:"state,omitempty"`
	PostalCode  *string `gorm:"size:20" json:"postal_code,omitempty"`
	Country     string  `gorm:"not null;size:50;default:'Portugal'" json:"country"`
	Phone       *string `gorm:"size:50" json:"phone,omitempty"`
	Email       *string `gorm:"size:255" json:"email,omitempty"`
	IsActive    bool    `gorm:"not null;default:true" json:"is_active"`
	IsMain      bool    `gorm:"not null;default:false" json:"is_main"`
	Timezone    string  `gorm:"size:50;default:'Europe/Lisbon'" json:"timezone"`
	Settings    *string `gorm:"type:jsonb;default:'{}'" json:"settings,omitempty"`

	// Relationships
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
}

// BusinessSettings represents business configuration settings
type BusinessSettings struct {
	BaseModel
	BusinessID                   string  `gorm:"not null;type:uuid;uniqueIndex" json:"business_id"`
	CalendarStartHour            int     `gorm:"not null;default:9" json:"calendar_start_hour"`
	CalendarEndHour              int     `gorm:"not null;default:18" json:"calendar_end_hour"`
	AppointmentBufferMinutes     int     `gorm:"not null;default:0" json:"appointment_buffer_minutes"`
	AllowOnlineBooking           bool    `gorm:"not null;default:true" json:"allow_online_booking"`
	DefaultAppointmentDuration   int     `gorm:"not null;default:60" json:"default_appointment_duration"`
	Currency                     string  `gorm:"not null;size:3;default:'EUR'" json:"currency"`
	DateFormat                   string  `gorm:"not null;size:20;default:'DD-MM-YYYY'" json:"date_format"`
	TimeFormat                   string  `gorm:"not null;size:10;default:'24h'" json:"time_format"`

	// Relationships
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
}

// TableName returns the table name for Business
func (Business) TableName() string {
	return "businesses"
}

// TableName returns the table name for BusinessLocation
func (BusinessLocation) TableName() string {
	return "business_locations"
}

// TableName returns the table name for BusinessSettings
func (BusinessSettings) TableName() string {
	return "business_settings"
}

// GetDisplayName returns the display name or name if display name is not set
func (b Business) GetDisplayName() string {
	if b.DisplayName != nil && *b.DisplayName != "" {
		return *b.DisplayName
	}
	return b.Name
}

// IsOwner checks if the given user ID is the owner of this business
func (b Business) IsOwner(userID string) bool {
	return b.UserID == userID
}

// GetMainLocation returns the main location for the business
func (b Business) GetMainLocation() *BusinessLocation {
	for _, location := range b.Locations {
		if location.IsMain {
			return &location
		}
	}
	return nil
}

// Validate validates the business model
func (b *Business) Validate() error {
	if b.UserID == "" {
		return ErrValidation
	}
	if b.Name == "" {
		return ErrValidation
	}
	if b.Email == "" {
		return ErrValidation
	}
	if b.Currency == "" || len(b.Currency) != 3 {
		return ErrValidation
	}
	if b.TimeZone == "" {
		return ErrValidation
	}
	return nil
}

// BeforeCreate is called before creating a business
func (b *Business) BeforeCreate(tx *gorm.DB) error {
	if b.Currency == "" {
		b.Currency = "EUR"
	}
	if b.TimeZone == "" {
		b.TimeZone = "Europe/Lisbon"
	}
	if b.SubscriptionTier == "" {
		b.SubscriptionTier = "free"
	}
	return nil
}

// Validate validates the business location model
func (bl *BusinessLocation) Validate() error {
	if bl.BusinessID == "" {
		return ErrValidation
	}
	if bl.Name == "" {
		return ErrValidation
	}
	if bl.Country == "" {
		return ErrValidation
	}
	return nil
}

// Validate validates the business settings model
func (bs *BusinessSettings) Validate() error {
	if bs.BusinessID == "" {
		return ErrValidation
	}
	if bs.CalendarStartHour < 0 || bs.CalendarStartHour >= 24 {
		return ErrValidation
	}
	if bs.CalendarEndHour <= 0 || bs.CalendarEndHour > 24 {
		return ErrValidation
	}
	if bs.CalendarStartHour >= bs.CalendarEndHour {
		return ErrValidation
	}
	return nil
}

// BusinessRepository defines the repository interface for Business
type BusinessRepository interface {
	BaseRepository[Business]
	FindByUserID(ctx context.Context, userID string) ([]*Business, error)
	FindByName(ctx context.Context, name string) ([]*Business, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	FindActiveBusinesses(ctx context.Context, page, pageSize int) ([]*Business, int64, error)
	SearchByLocation(ctx context.Context, city, country string) ([]*Business, error)
	SearchByService(ctx context.Context, serviceName string) ([]*Business, error)
	GetBusinessWithDetails(ctx context.Context, businessID string) (*Business, error)
	GetWithLocations(ctx context.Context, businessID string) (*Business, error)
}

// BusinessLocationRepository defines the repository interface for BusinessLocation
type BusinessLocationRepository interface {
	BaseRepository[BusinessLocation]
	FindByBusinessID(ctx context.Context, businessID string) ([]*BusinessLocation, error)
	GetMainLocation(ctx context.Context, businessID string) (*BusinessLocation, error)
	SetMainLocation(ctx context.Context, businessID, locationID string) error
}

// BusinessSettingsRepository defines the repository interface for BusinessSettings
type BusinessSettingsRepository interface {
	BaseRepository[BusinessSettings]
	GetByBusinessID(ctx context.Context, businessID string) (*BusinessSettings, error)
	UpdateByBusinessID(ctx context.Context, businessID string, settings *BusinessSettings) error
}

// Note: Service interfaces will be defined in the service layer to avoid circular dependencies