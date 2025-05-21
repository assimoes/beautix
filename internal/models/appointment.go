package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	// AppointmentStatusPending represents a pending appointment
	AppointmentStatusPending AppointmentStatus = "pending"
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
	// AppointmentStatusRescheduled represents a rescheduled appointment
	AppointmentStatusRescheduled AppointmentStatus = "rescheduled"
)

// PaymentStatus represents the payment status of an appointment
type PaymentStatus string

const (
	// PaymentStatusPending represents a pending payment
	PaymentStatusPending PaymentStatus = "pending"
	// PaymentStatusPaid represents a completed payment
	PaymentStatusPaid PaymentStatus = "paid"
	// PaymentStatusPartiallyPaid represents a partially paid payment
	PaymentStatusPartiallyPaid PaymentStatus = "partially_paid"
	// PaymentStatusRefunded represents a refunded payment
	PaymentStatusRefunded PaymentStatus = "refunded"
	// PaymentStatusFailed represents a failed payment
	PaymentStatusFailed PaymentStatus = "failed"
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
	// PaymentMethodOnline represents online payment
	PaymentMethodOnline PaymentMethod = "online"
)

// BookingSource represents the source of the booking
type BookingSource string

const (
	// BookingSourceWebsite represents booking from the website
	BookingSourceWebsite BookingSource = "website"
	// BookingSourceApp represents booking from the mobile app
	BookingSourceApp BookingSource = "app"
	// BookingSourceInPerson represents in-person booking
	BookingSourceInPerson BookingSource = "in_person"
	// BookingSourcePhone represents booking by phone
	BookingSourcePhone BookingSource = "phone"
	// BookingSourceThirdParty represents booking through a third-party platform
	BookingSourceThirdParty BookingSource = "third_party"
)

// Appointment represents a booking for a service
type Appointment struct {
	BaseModel
	BusinessID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"business_id"`
	Business         Business         `gorm:"foreignKey:BusinessID" json:"business"`
	ClientID         uuid.UUID        `gorm:"type:uuid;not null;index" json:"client_id"`
	Client           Client           `gorm:"foreignKey:ClientID" json:"client"`
	LocationID       *uuid.UUID       `gorm:"type:uuid;index" json:"location_id"`
	BusinessLocation *BusinessLocation `gorm:"foreignKey:LocationID" json:"location"`
	StartTime        time.Time        `gorm:"not null;index" json:"start_time"`
	EndTime          time.Time        `gorm:"not null" json:"end_time"`
	Status           AppointmentStatus `gorm:"type:text;not null;default:'pending'" json:"status"`
	TotalDuration    int              `gorm:"not null" json:"total_duration"` // In minutes
	TotalPrice       float64          `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Discount         float64          `gorm:"type:decimal(10,2);not null;default:0" json:"discount"`
	Tax              float64          `gorm:"type:decimal(10,2);not null;default:0" json:"tax"`
	FinalPrice       float64          `gorm:"type:decimal(10,2);not null" json:"final_price"`
	PaymentStatus    PaymentStatus    `gorm:"type:text;not null;default:'pending'" json:"payment_status"`
	PaymentMethod    *PaymentMethod   `gorm:"type:text" json:"payment_method"`
	Notes            string           `json:"notes"`
	CancellationReason string         `json:"cancellation_reason"`
	CancellationFee  float64          `gorm:"type:decimal(10,2);not null;default:0" json:"cancellation_fee"`
	IsRescheduled    bool             `gorm:"not null;default:false" json:"is_rescheduled"`
	PreviousAppointmentID *uuid.UUID  `gorm:"type:uuid;index" json:"previous_appointment_id"`
	BookingSource    BookingSource    `gorm:"type:text;not null;default:'in_person'" json:"booking_source"`
	BookingReference string           `json:"booking_reference"`
	Notifications    NotificationSettings `gorm:"type:jsonb" json:"notifications"`
	InternalNotes    string           `json:"internal_notes"`
	CheckInTime      *time.Time       `json:"check_in_time"`
	CheckOutTime     *time.Time       `json:"check_out_time"`
	NoShowFee        float64          `gorm:"type:decimal(10,2);not null;default:0" json:"no_show_fee"`
	DepositAmount    float64          `gorm:"type:decimal(10,2);not null;default:0" json:"deposit_amount"`
	DepositPaid      bool             `gorm:"not null;default:false" json:"deposit_paid"`
}

// TableName overrides the table name
func (Appointment) TableName() string {
	return "appointments"
}

// NotificationSettings stores settings for appointment notifications
type NotificationSettings struct {
	ClientReminded    bool      `json:"client_reminded"`
	ClientReminderSent time.Time `json:"client_reminder_sent,omitempty"`
	ConfirmationSent  bool      `json:"confirmation_sent"`
	ConfirmationTime  time.Time `json:"confirmation_time,omitempty"`
	FollowUpSent      bool      `json:"follow_up_sent"`
	FollowUpTime      time.Time `json:"follow_up_time,omitempty"`
}

// Scan implements the sql.Scanner interface for NotificationSettings
func (ns *NotificationSettings) Scan(value interface{}) error {
	if value == nil {
		*ns = NotificationSettings{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp NotificationSettings
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*ns = temp

	return nil
}

// Value implements the driver.Valuer interface for NotificationSettings
func (ns NotificationSettings) Value() (driver.Value, error) {
	return json.Marshal(ns)
}

// AppointmentService represents a service within an appointment
type AppointmentService struct {
	BaseModel
	AppointmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"appointment_id"`
	Appointment   Appointment `gorm:"foreignKey:AppointmentID" json:"appointment"`
	ServiceID     uuid.UUID `gorm:"type:uuid;not null;index" json:"service_id"`
	Service       Service   `gorm:"foreignKey:ServiceID" json:"service"`
	StaffID       uuid.UUID `gorm:"type:uuid;not null;index" json:"staff_id"`
	Staff         Staff     `gorm:"foreignKey:StaffID" json:"staff"`
	VariantID     *uuid.UUID `gorm:"type:uuid;index" json:"variant_id"`
	ServiceVariant *ServiceVariant `gorm:"foreignKey:VariantID" json:"service_variant"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Discount      float64   `gorm:"type:decimal(10,2);not null;default:0" json:"discount"`
	Duration      int       `gorm:"not null" json:"duration"` // In minutes
	StartTime     time.Time `gorm:"not null" json:"start_time"`
	EndTime       time.Time `gorm:"not null" json:"end_time"`
	Notes         string    `json:"notes"`
	Status        AppointmentStatus `gorm:"type:text;not null;default:'pending'" json:"status"`
	SelectedOptions AppointmentServiceOptions `gorm:"type:jsonb" json:"selected_options"`
}

// TableName overrides the table name
func (AppointmentService) TableName() string {
	return "appointment_services"
}

// AppointmentServiceOptions stores the options selected for a service
type AppointmentServiceOptions []AppointmentServiceOption

// AppointmentServiceOption represents a single option chosen for a service
type AppointmentServiceOption struct {
	OptionID       uuid.UUID `json:"option_id"`
	OptionName     string    `json:"option_name"`
	ChoiceID       uuid.UUID `json:"choice_id"`
	ChoiceName     string    `json:"choice_name"`
	PriceAdjustment float64   `json:"price_adjustment"`
	TimeAdjustment  int       `json:"time_adjustment"`
}

// Scan implements the sql.Scanner interface for AppointmentServiceOptions
func (aso *AppointmentServiceOptions) Scan(value interface{}) error {
	if value == nil {
		*aso = AppointmentServiceOptions{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp AppointmentServiceOptions
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*aso = temp

	return nil
}

// Value implements the driver.Valuer interface for AppointmentServiceOptions
func (aso AppointmentServiceOptions) Value() (driver.Value, error) {
	return json.Marshal(aso)
}

// AppointmentPayment represents a payment for an appointment
type AppointmentPayment struct {
	BaseModel
	AppointmentID  uuid.UUID     `gorm:"type:uuid;not null;index" json:"appointment_id"`
	Appointment    Appointment   `gorm:"foreignKey:AppointmentID" json:"appointment"`
	Amount         float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentMethod  PaymentMethod `gorm:"type:text;not null" json:"payment_method"`
	PaymentStatus  PaymentStatus `gorm:"type:text;not null" json:"payment_status"`
	TransactionID  string        `json:"transaction_id"`
	Notes          string        `json:"notes"`
	RefundAmount   float64       `gorm:"type:decimal(10,2);not null;default:0" json:"refund_amount"`
	IsDeposit      bool          `gorm:"not null;default:false" json:"is_deposit"`
	PaymentDate    time.Time     `gorm:"not null" json:"payment_date"`
	ReceiptURL     string        `json:"receipt_url"`
	FailureReason  string        `json:"failure_reason"`
}

// TableName overrides the table name
func (AppointmentPayment) TableName() string {
	return "appointment_payments"
}

// AppointmentForm represents a form filled for an appointment
type AppointmentForm struct {
	BaseModel
	AppointmentID uuid.UUID    `gorm:"type:uuid;not null;index" json:"appointment_id"`
	Appointment   Appointment  `gorm:"foreignKey:AppointmentID" json:"appointment"`
	FormType      string       `gorm:"not null" json:"form_type"` // "consent", "intake", "medical", etc.
	FormData      FormData     `gorm:"type:jsonb" json:"form_data"`
	IsCompleted   bool         `gorm:"not null;default:false" json:"is_completed"`
	CompletedAt   *time.Time   `json:"completed_at"`
	SignatureURL  string       `json:"signature_url"`
}

// TableName overrides the table name
func (AppointmentForm) TableName() string {
	return "appointment_forms"
}

// FormData stores the data for a form
type FormData map[string]interface{}

// Scan implements the sql.Scanner interface for FormData
func (fd *FormData) Scan(value interface{}) error {
	if value == nil {
		*fd = FormData{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp FormData
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*fd = temp

	return nil
}

// Value implements the driver.Valuer interface for FormData
func (fd FormData) Value() (driver.Value, error) {
	return json.Marshal(fd)
}

// AppointmentNote represents a note for an appointment
type AppointmentNote struct {
	BaseModel
	AppointmentID  uuid.UUID   `gorm:"type:uuid;not null;index" json:"appointment_id"`
	Appointment    Appointment `gorm:"foreignKey:AppointmentID" json:"appointment"`
	Title          string      `gorm:"not null" json:"title"`
	Content        string      `gorm:"not null" json:"content"`
	IsPrivate      bool        `gorm:"not null;default:false" json:"is_private"`
	AttachmentURL  string      `json:"attachment_url"`
}

// TableName overrides the table name
func (AppointmentNote) TableName() string {
	return "appointment_notes"
}

// AppointmentReminder represents a reminder for an appointment
type AppointmentReminder struct {
	BaseModel
	AppointmentID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"appointment_id"`
	Appointment    Appointment `gorm:"foreignKey:AppointmentID" json:"appointment"`
	ReminderType   string     `gorm:"not null" json:"reminder_type"` // "sms", "email", "push", "whatsapp"
	ReminderTime   time.Time  `gorm:"not null" json:"reminder_time"`
	IsSent         bool       `gorm:"not null;default:false" json:"is_sent"`
	SentAt         *time.Time `json:"sent_at"`
	IsSuccess      bool       `gorm:"not null;default:false" json:"is_success"`
	ErrorMessage   string     `json:"error_message"`
}

// TableName overrides the table name
func (AppointmentReminder) TableName() string {
	return "appointment_reminders"
}

// AppointmentFeedback represents client feedback for an appointment
type AppointmentFeedback struct {
	BaseModel
	AppointmentID uuid.UUID   `gorm:"type:uuid;not null;index" json:"appointment_id"`
	Appointment   Appointment `gorm:"foreignKey:AppointmentID" json:"appointment"`
	Rating        int         `gorm:"not null" json:"rating"` // 1-5 stars
	Comments      string      `json:"comments"`
	IsPublic      bool        `gorm:"not null;default:false" json:"is_public"`
	FeedbackItems FeedbackItems `gorm:"type:jsonb" json:"feedback_items"`
}

// TableName overrides the table name
func (AppointmentFeedback) TableName() string {
	return "appointment_feedback"
}

// FeedbackItems stores specific feedback categories and ratings
type FeedbackItems struct {
	ServiceQuality    int `json:"service_quality,omitempty"`     // 1-5
	StaffProfessionalism int `json:"staff_professionalism,omitempty"` // 1-5
	Cleanliness       int `json:"cleanliness,omitempty"`        // 1-5
	ValueForMoney     int `json:"value_for_money,omitempty"`    // 1-5
	Atmosphere        int `json:"atmosphere,omitempty"`         // 1-5
	WouldRecommend    bool `json:"would_recommend,omitempty"`
}

// Scan implements the sql.Scanner interface for FeedbackItems
func (fi *FeedbackItems) Scan(value interface{}) error {
	if value == nil {
		*fi = FeedbackItems{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp FeedbackItems
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*fi = temp

	return nil
}

// Value implements the driver.Valuer interface for FeedbackItems
func (fi FeedbackItems) Value() (driver.Value, error) {
	return json.Marshal(fi)
}