package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

// SubscriptionTier represents a business subscription tier
type SubscriptionTier string

const (
	// SubscriptionTierFree represents the free tier
	SubscriptionTierFree SubscriptionTier = "free"
	// SubscriptionTierBasic represents the basic paid tier
	SubscriptionTierBasic SubscriptionTier = "basic"
	// SubscriptionTierPro represents the professional tier
	SubscriptionTierPro SubscriptionTier = "pro"
	// SubscriptionTierEnterprise represents the enterprise tier
	SubscriptionTierEnterprise SubscriptionTier = "enterprise"
)

// Business represents a business entity that provides beauty services
type Business struct {
	BaseModel
	UserID           uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	User             User            `gorm:"foreignKey:UserID" json:"user"`
	Name             string          `gorm:"not null" json:"name"`
	DisplayName      string          `gorm:"not null" json:"display_name"`
	Description      string          `json:"description"`
	Address          string          `json:"address"`
	City             string          `json:"city"`
	State            string          `json:"state"`
	PostalCode       string          `json:"postal_code"`
	Country          string          `gorm:"not null;default:'Portugal'" json:"country"`
	Phone            string          `json:"phone"`
	Email            string          `json:"email"`
	Website          string          `json:"website"`
	LogoURL          string          `json:"logo_url"`
	CoverPhotoURL    string          `json:"cover_photo_url"`
	TaxID            string          `json:"tax_id"`
	SubscriptionTier SubscriptionTier `gorm:"type:text;not null;default:'free'" json:"subscription_tier"`
	IsActive         bool            `gorm:"not null;default:true" json:"is_active"`
	IsVerified       bool            `gorm:"not null;default:false" json:"is_verified"`
	Timezone         string          `gorm:"not null;default:'Europe/Lisbon'" json:"timezone"`
	Currency         string          `gorm:"not null;default:'EUR'" json:"currency"`
	SocialLinks      SocialLinks     `gorm:"type:jsonb" json:"social_links"`
	Settings         BusinessSettings `gorm:"type:jsonb" json:"settings"`
}

// TableName overrides the table name
func (Business) TableName() string {
	return "businesses"
}

// SocialLinks stores URLs to social media profiles
type SocialLinks struct {
	Facebook  string `json:"facebook,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	LinkedIn  string `json:"linkedin,omitempty"`
	TikTok    string `json:"tiktok,omitempty"`
	YouTube   string `json:"youtube,omitempty"`
}

// Scan implements the sql.Scanner interface for SocialLinks
func (s *SocialLinks) Scan(value interface{}) error {
	if value == nil {
		*s = SocialLinks{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp SocialLinks
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*s = temp

	return nil
}

// Value implements the driver.Valuer interface for SocialLinks
func (s SocialLinks) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// BusinessSettings stores business-specific settings
type BusinessSettings struct {
	AllowOnlineBooking       bool   `json:"allow_online_booking"`
	RequireDeposit           bool   `json:"require_deposit"`
	DepositAmount            float64 `json:"deposit_amount,omitempty"`
	CancellationPolicyHours  int    `json:"cancellation_policy_hours"`
	CancellationFeePercentage int    `json:"cancellation_fee_percentage"`
	WorkingHours             WorkingHours `json:"working_hours"`
	BookingNotificationsEnabled bool  `json:"booking_notifications_enabled"`
	MarketingEnabled         bool   `json:"marketing_enabled"`
	CustomTheme              Theme  `json:"custom_theme,omitempty"`
}

// Scan implements the sql.Scanner interface for BusinessSettings
func (bs *BusinessSettings) Scan(value interface{}) error {
	if value == nil {
		*bs = BusinessSettings{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp BusinessSettings
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*bs = temp

	return nil
}

// Value implements the driver.Valuer interface for BusinessSettings
func (bs BusinessSettings) Value() (driver.Value, error) {
	return json.Marshal(bs)
}

// Theme defines the branding colors for a business
type Theme struct {
	PrimaryColor   string `json:"primary_color,omitempty"`
	SecondaryColor string `json:"secondary_color,omitempty"`
	AccentColor    string `json:"accent_color,omitempty"`
	TextColor      string `json:"text_color,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
}

// Scan implements the sql.Scanner interface for Theme
func (t *Theme) Scan(value interface{}) error {
	if value == nil {
		*t = Theme{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp Theme
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*t = temp

	return nil
}

// Value implements the driver.Valuer interface for Theme
func (t Theme) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// WorkingHours defines the business hours for each day of the week
type WorkingHours struct {
	Monday    DayHours `json:"monday"`
	Tuesday   DayHours `json:"tuesday"`
	Wednesday DayHours `json:"wednesday"`
	Thursday  DayHours `json:"thursday"`
	Friday    DayHours `json:"friday"`
	Saturday  DayHours `json:"saturday"`
	Sunday    DayHours `json:"sunday"`
}

// Scan implements the sql.Scanner interface for WorkingHours
func (wh *WorkingHours) Scan(value interface{}) error {
	if value == nil {
		*wh = WorkingHours{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp WorkingHours
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*wh = temp

	return nil
}

// Value implements the driver.Valuer interface for WorkingHours
func (wh WorkingHours) Value() (driver.Value, error) {
	return json.Marshal(wh)
}

// DayHours defines working hours for a specific day
type DayHours struct {
	IsOpen     bool   `json:"is_open"`
	OpenTime   string `json:"open_time,omitempty"`  // Format: "09:00"
	CloseTime  string `json:"close_time,omitempty"` // Format: "18:00"
	BreakStart string `json:"break_start,omitempty"`
	BreakEnd   string `json:"break_end,omitempty"`
}

// Scan implements the sql.Scanner interface for DayHours
func (dh *DayHours) Scan(value interface{}) error {
	if value == nil {
		*dh = DayHours{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp DayHours
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*dh = temp

	return nil
}

// Value implements the driver.Valuer interface for DayHours
func (dh DayHours) Value() (driver.Value, error) {
	return json.Marshal(dh)
}

// BusinessLocation represents a physical location for a business (for businesses with multiple locations)
type BusinessLocation struct {
	BaseModel
	BusinessID uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	Name       string    `gorm:"not null" json:"name"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Country    string    `gorm:"not null;default:'Portugal'" json:"country"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	IsMain     bool      `gorm:"not null;default:false" json:"is_main"`
	Timezone   string    `gorm:"not null;default:'Europe/Lisbon'" json:"timezone"`
	Settings   LocationSettings `gorm:"type:jsonb" json:"settings"`
}

// LocationSettings stores location-specific settings
type LocationSettings struct {
	WorkingHours WorkingHours `json:"working_hours"`
	Capacity     int          `json:"capacity"`
}

// Scan implements the sql.Scanner interface for LocationSettings
func (ls *LocationSettings) Scan(value interface{}) error {
	if value == nil {
		*ls = LocationSettings{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp LocationSettings
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*ls = temp

	return nil
}

// Value implements the driver.Valuer interface for LocationSettings
func (ls LocationSettings) Value() (driver.Value, error) {
	return json.Marshal(ls)
}

// TableName overrides the table name
func (BusinessLocation) TableName() string {
	return "business_locations"
}