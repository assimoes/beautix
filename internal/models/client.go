package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Client represents a client of a beauty business
type Client struct {
	BaseModel
	BusinessID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"business_id"`
	Business        Business        `gorm:"foreignKey:BusinessID" json:"business"`
	UserID          *uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	User            *User           `gorm:"foreignKey:UserID" json:"user"`
	FirstName       string          `gorm:"not null" json:"first_name"`
	LastName        string          `gorm:"not null" json:"last_name"`
	Email           string          `gorm:"index" json:"email"`
	Phone           string          `gorm:"index" json:"phone"`
	DateOfBirth     *time.Time      `json:"date_of_birth"`
	Address         string          `json:"address"`
	City            string          `json:"city"`
	State           string          `json:"state"`
	PostalCode      string          `json:"postal_code"`
	Country         string          `json:"country"`
	Notes           string          `json:"notes"`
	ProfileImageURL string          `json:"profile_image_url"`
	Tags            ClientTags      `gorm:"type:jsonb" json:"tags"`
	Preferences     ClientPreferences `gorm:"type:jsonb" json:"preferences"`
	HealthInfo      HealthInfo      `gorm:"type:jsonb" json:"health_info"`
	Source          string          `json:"source"`
	AcceptsMarketing bool           `gorm:"not null;default:false" json:"accepts_marketing"`
	IsActive        bool            `gorm:"not null;default:true" json:"is_active"`
}

// TableName overrides the table name
func (Client) TableName() string {
	return "clients"
}

// ClientTags is a string array that can be stored as JSONB
type ClientTags []string

// Scan implements the sql.Scanner interface for ClientTags
func (ct *ClientTags) Scan(value interface{}) error {
	if value == nil {
		*ct = ClientTags{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp ClientTags
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*ct = temp

	return nil
}

// Value implements the driver.Valuer interface for ClientTags
func (ct ClientTags) Value() (driver.Value, error) {
	return json.Marshal(ct)
}

// ClientPreferences stores client preferences as JSONB
type ClientPreferences struct {
	PreferredStaffIDs  []uuid.UUID `json:"preferred_staff_ids,omitempty"`
	PreferredDays      []string    `json:"preferred_days,omitempty"`      // "monday", "tuesday", etc.
	PreferredTimeStart string      `json:"preferred_time_start,omitempty"` // "09:00"
	PreferredTimeEnd   string      `json:"preferred_time_end,omitempty"`   // "18:00"
	CommunicationPrefs CommunicationPreferences `json:"communication_prefs"`
	ReminderSettings   ReminderSettings `json:"reminder_settings"`
	LanguagePreference string      `json:"language_preference,omitempty"`
}

// Scan implements the sql.Scanner interface for ClientPreferences
func (cp *ClientPreferences) Scan(value interface{}) error {
	if value == nil {
		*cp = ClientPreferences{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp ClientPreferences
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*cp = temp

	return nil
}

// Value implements the driver.Valuer interface for ClientPreferences
func (cp ClientPreferences) Value() (driver.Value, error) {
	return json.Marshal(cp)
}

// CommunicationPreferences stores how a client prefers to be contacted
type CommunicationPreferences struct {
	AllowEmail        bool `json:"allow_email"`
	AllowSMS          bool `json:"allow_sms"`
	AllowPush         bool `json:"allow_push"`
	AllowWhatsApp     bool `json:"allow_whatsapp"`
	MarketingEmails   bool `json:"marketing_emails"`
	AppointmentReminders bool `json:"appointment_reminders"`
	PromotionalOffers bool `json:"promotional_offers"`
}

// ReminderSettings stores when and how appointment reminders should be sent
type ReminderSettings struct {
	DaysBefore        int    `json:"days_before"`
	HoursBefore       int    `json:"hours_before"`
	MinutesBefore     int    `json:"minutes_before"`
	PreferredChannel  string `json:"preferred_channel"`  // "email", "sms", "push", "whatsapp"
	SecondaryChannel  string `json:"secondary_channel"`
}

// HealthInfo stores health-related information about the client
type HealthInfo struct {
	Allergies       []string `json:"allergies,omitempty"`
	HealthConditions []string `json:"health_conditions,omitempty"`
	Medications     []string `json:"medications,omitempty"`
	EmergencyContact *EmergencyContact `json:"emergency_contact,omitempty"`
	Notes           string   `json:"notes,omitempty"`
}

// Scan implements the sql.Scanner interface for HealthInfo
func (hi *HealthInfo) Scan(value interface{}) error {
	if value == nil {
		*hi = HealthInfo{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp HealthInfo
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*hi = temp

	return nil
}

// Value implements the driver.Valuer interface for HealthInfo
func (hi HealthInfo) Value() (driver.Value, error) {
	return json.Marshal(hi)
}

// EmergencyContact stores emergency contact information
type EmergencyContact struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
	Relationship string `json:"relationship,omitempty"`
}

// ClientNote represents a note about a client
type ClientNote struct {
	BaseModel
	BusinessID uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	ClientID   uuid.UUID `gorm:"type:uuid;not null;index" json:"client_id"`
	Client     Client    `gorm:"foreignKey:ClientID" json:"client"`
	Title      string    `gorm:"not null" json:"title"`
	Content    string    `gorm:"not null" json:"content"`
	IsPrivate  bool      `gorm:"not null;default:false" json:"is_private"`
	Pinned     bool      `gorm:"not null;default:false" json:"pinned"`
}

// TableName overrides the table name
func (ClientNote) TableName() string {
	return "client_notes"
}

// ClientDocument represents a document attached to a client record
type ClientDocument struct {
	BaseModel
	BusinessID    uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
	Business      Business  `gorm:"foreignKey:BusinessID" json:"business"`
	ClientID      uuid.UUID `gorm:"type:uuid;not null;index" json:"client_id"`
	Client        Client    `gorm:"foreignKey:ClientID" json:"client"`
	DocumentName  string    `gorm:"not null" json:"document_name"`
	DocumentType  string    `gorm:"not null" json:"document_type"`  // "consent_form", "medical_form", etc.
	FileURL       string    `gorm:"not null" json:"file_url"`
	ContentType   string    `json:"content_type"`
	FileSize      int64     `json:"file_size"`
	IsSignatureRequired bool  `gorm:"not null;default:false" json:"is_signature_required"`
	SignedAt      *time.Time `json:"signed_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	IsPrivate     bool      `gorm:"not null;default:false" json:"is_private"`
}

// TableName overrides the table name
func (ClientDocument) TableName() string {
	return "client_documents"
}