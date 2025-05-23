package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Campaign represents a marketing campaign or promotion
type Campaign struct {
	ID              uuid.UUID      `json:"id"`
	BusinessID      uuid.UUID      `json:"business_id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	CampaignType    string         `json:"campaign_type"`    // promotion, seasonal, reactivation, birthday
	TargetAudience  TargetAudience `json:"target_audience"`  // Criteria for selecting clients
	OfferType       string         `json:"offer_type"`       // discount, free_service, bundle, gift
	OfferDetails    OfferDetails   `json:"offer_details"`    // Details of the offer
	StartDate       time.Time      `json:"start_date"`
	EndDate         time.Time      `json:"end_date"`
	IsActive        bool           `json:"is_active"`
	MessageTemplate string         `json:"message_template"`
	CreatedAt       time.Time      `json:"created_at"`
	CreatedBy       uuid.UUID      `json:"created_by"`
	UpdatedAt       *time.Time     `json:"updated_at,omitempty"`
	UpdatedBy       *uuid.UUID     `json:"updated_by,omitempty"`
	DeletedAt       *time.Time     `json:"deleted_at,omitempty"`
	DeletedBy       *uuid.UUID     `json:"deleted_by,omitempty"`

	// Expanded relationships (populated by service when needed)
	Business *Business        `json:"business,omitempty"`
	Clients  []*CampaignClient `json:"clients,omitempty"`
	Messages []*CampaignMessage `json:"messages,omitempty"`
}

// TargetAudience represents criteria for selecting clients for a campaign
type TargetAudience struct {
	AgeMin            *int        `json:"age_min,omitempty"`
	AgeMax            *int        `json:"age_max,omitempty"`
	Gender            *string     `json:"gender,omitempty"`           // male, female, other, any
	LocationRadius    *float64    `json:"location_radius,omitempty"`  // km from business
	LastVisitDaysAgo  *int        `json:"last_visit_days_ago,omitempty"` // Days since last visit
	TotalSpentMin     *float64    `json:"total_spent_min,omitempty"`
	TotalSpentMax     *float64    `json:"total_spent_max,omitempty"`
	VisitCountMin     *int        `json:"visit_count_min,omitempty"`
	VisitCountMax     *int        `json:"visit_count_max,omitempty"`
	ServiceIDs        []uuid.UUID `json:"service_ids,omitempty"`      // Clients who used specific services
	LoyaltyTier       *string     `json:"loyalty_tier,omitempty"`     // Specific loyalty tier
	BirthdayMonth     *int        `json:"birthday_month,omitempty"`   // For birthday campaigns
	ClientTags        []string    `json:"client_tags,omitempty"`      // Custom client tags
	HasLoyaltyPoints  *bool       `json:"has_loyalty_points,omitempty"` // Clients with/without points
	IsNewClient       *bool       `json:"is_new_client,omitempty"`    // Clients registered recently
	ClientIDs         []uuid.UUID `json:"client_ids,omitempty"`       // Specific client IDs
}

// OfferDetails represents the details of a campaign offer
type OfferDetails struct {
	DiscountPercentage   *float64    `json:"discount_percentage,omitempty"`   // For percentage discounts
	DiscountAmount       *float64    `json:"discount_amount,omitempty"`       // For fixed amount discounts
	FreeServiceID        *uuid.UUID  `json:"free_service_id,omitempty"`       // For free service offers
	BundleServiceIDs     []uuid.UUID `json:"bundle_service_ids,omitempty"`    // For bundle offers
	BundlePrice          *float64    `json:"bundle_price,omitempty"`          // Special bundle price
	GiftValue            *float64    `json:"gift_value,omitempty"`            // For gift offers
	MinimumSpend         *float64    `json:"minimum_spend,omitempty"`         // Minimum spend requirement
	MaxUsage             *int        `json:"max_usage,omitempty"`             // Max uses per client
	CombineWithLoyalty   bool        `json:"combine_with_loyalty"`            // Can combine with loyalty discounts
	ValidServiceIDs      []uuid.UUID `json:"valid_service_ids,omitempty"`     // Services offer applies to
	ExcludedServiceIDs   []uuid.UUID `json:"excluded_service_ids,omitempty"`  // Services excluded from offer
	RequiredAdvanceBooking *int      `json:"required_advance_booking,omitempty"` // Days in advance booking required
	ValidDaysOfWeek      []int       `json:"valid_days_of_week,omitempty"`    // 0=Sunday, 1=Monday, etc.
	ValidTimeStart       *string     `json:"valid_time_start,omitempty"`      // HH:MM format
	ValidTimeEnd         *string     `json:"valid_time_end,omitempty"`        // HH:MM format
	Description          string      `json:"description,omitempty"`           // Human readable description
}

// CampaignClient represents a client targeted by a campaign and their response
type CampaignClient struct {
	ID              uuid.UUID  `json:"id"`
	CampaignID      uuid.UUID  `json:"campaign_id"`
	ClientID        uuid.UUID  `json:"client_id"`
	Status          string     `json:"status"`          // pending, sent, opened, clicked, converted, unsubscribed
	SentAt          *time.Time `json:"sent_at,omitempty"`
	OpenedAt        *time.Time `json:"opened_at,omitempty"`
	ClickedAt       *time.Time `json:"clicked_at,omitempty"`
	ConvertedAt     *time.Time `json:"converted_at,omitempty"`
	ConversionValue *float64   `json:"conversion_value,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`

	// Expanded relationships (populated by service when needed)
	Campaign *Campaign `json:"campaign,omitempty"`
	Client   *Client   `json:"client,omitempty"`
}

// CampaignMessage represents a message sent to a client as part of a campaign
type CampaignMessage struct {
	ID            uuid.UUID  `json:"id"`
	CampaignID    uuid.UUID  `json:"campaign_id"`
	ClientID      uuid.UUID  `json:"client_id"`
	MessageType   string     `json:"message_type"`   // email, sms, push, whatsapp
	MessageContent string    `json:"message_content"`
	ScheduledTime time.Time  `json:"scheduled_time"`
	SentTime      *time.Time `json:"sent_time,omitempty"`
	Status        string     `json:"status"`        // pending, sent, failed, cancelled
	ErrorMessage  string     `json:"error_message,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`

	// Expanded relationships (populated by service when needed)
	Campaign *Campaign `json:"campaign,omitempty"`
	Client   *Client   `json:"client,omitempty"`
}

// Input types for creating and updating campaigns

// CreateCampaignInput is the input for creating a campaign
type CreateCampaignInput struct {
	BusinessID      uuid.UUID      `json:"business_id" validate:"required"`
	Name            string         `json:"name" validate:"required"`
	Description     string         `json:"description"`
	CampaignType    string         `json:"campaign_type" validate:"required,oneof=promotion seasonal reactivation birthday"`
	TargetAudience  TargetAudience `json:"target_audience"`
	OfferType       string         `json:"offer_type" validate:"required,oneof=discount free_service bundle gift"`
	OfferDetails    OfferDetails   `json:"offer_details" validate:"required"`
	StartDate       time.Time      `json:"start_date" validate:"required"`
	EndDate         time.Time      `json:"end_date" validate:"required,gtfield=StartDate"`
	MessageTemplate string         `json:"message_template"`
}

// UpdateCampaignInput is the input for updating a campaign
type UpdateCampaignInput struct {
	Name            *string         `json:"name"`
	Description     *string         `json:"description"`
	CampaignType    *string         `json:"campaign_type" validate:"omitempty,oneof=promotion seasonal reactivation birthday"`
	TargetAudience  *TargetAudience `json:"target_audience"`
	OfferType       *string         `json:"offer_type" validate:"omitempty,oneof=discount free_service bundle gift"`
	OfferDetails    *OfferDetails   `json:"offer_details"`
	StartDate       *time.Time      `json:"start_date"`
	EndDate         *time.Time      `json:"end_date"`
	IsActive        *bool           `json:"is_active"`
	MessageTemplate *string         `json:"message_template"`
}

// CreateCampaignClientInput is the input for adding a client to a campaign
type CreateCampaignClientInput struct {
	CampaignID uuid.UUID `json:"campaign_id" validate:"required"`
	ClientID   uuid.UUID `json:"client_id" validate:"required"`
}

// UpdateCampaignClientInput is the input for updating campaign client status
type UpdateCampaignClientInput struct {
	Status          *string    `json:"status" validate:"omitempty,oneof=pending sent opened clicked converted unsubscribed"`
	SentAt          *time.Time `json:"sent_at"`
	OpenedAt        *time.Time `json:"opened_at"`
	ClickedAt       *time.Time `json:"clicked_at"`
	ConvertedAt     *time.Time `json:"converted_at"`
	ConversionValue *float64   `json:"conversion_value"`
}

// CreateCampaignMessageInput is the input for creating a campaign message
type CreateCampaignMessageInput struct {
	CampaignID     uuid.UUID `json:"campaign_id" validate:"required"`
	ClientID       uuid.UUID `json:"client_id" validate:"required"`
	MessageType    string    `json:"message_type" validate:"required,oneof=email sms push whatsapp"`
	MessageContent string    `json:"message_content" validate:"required"`
	ScheduledTime  time.Time `json:"scheduled_time" validate:"required"`
}

// UpdateCampaignMessageInput is the input for updating a campaign message
type UpdateCampaignMessageInput struct {
	MessageContent *string    `json:"message_content"`
	ScheduledTime  *time.Time `json:"scheduled_time"`
	SentTime       *time.Time `json:"sent_time"`
	Status         *string    `json:"status" validate:"omitempty,oneof=pending sent failed cancelled"`
	ErrorMessage   *string    `json:"error_message"`
}

// Repository interfaces

// CampaignRepository defines methods for campaign data store
type CampaignRepository interface {
	Create(ctx context.Context, campaign *Campaign) error
	GetByID(ctx context.Context, id uuid.UUID) (*Campaign, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateCampaignInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Campaign, error)
	ListActiveByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Campaign, error)
	Count(ctx context.Context) (int64, error)
	CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
	CountActiveByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// CampaignClientRepository defines methods for campaign client data store
type CampaignClientRepository interface {
	Create(ctx context.Context, campaignClient *CampaignClient) error
	GetByID(ctx context.Context, id uuid.UUID) (*CampaignClient, error)
	GetByCampaignAndClient(ctx context.Context, campaignID, clientID uuid.UUID) (*CampaignClient, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateCampaignClientInput) error
	ListByCampaign(ctx context.Context, campaignID uuid.UUID, page, pageSize int) ([]*CampaignClient, error)
	ListByClient(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*CampaignClient, error)
	CountByCampaign(ctx context.Context, campaignID uuid.UUID) (int64, error)
	CountByCampaignAndStatus(ctx context.Context, campaignID uuid.UUID, status string) (int64, error)
}

// CampaignMessageRepository defines methods for campaign message data store
type CampaignMessageRepository interface {
	Create(ctx context.Context, message *CampaignMessage) error
	GetByID(ctx context.Context, id uuid.UUID) (*CampaignMessage, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateCampaignMessageInput) error
	ListByCampaign(ctx context.Context, campaignID uuid.UUID, page, pageSize int) ([]*CampaignMessage, error)
	ListByClient(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*CampaignMessage, error)
	ListPendingMessages(ctx context.Context, beforeTime time.Time, page, pageSize int) ([]*CampaignMessage, error)
	CountByCampaign(ctx context.Context, campaignID uuid.UUID) (int64, error)
	CountByCampaignAndStatus(ctx context.Context, campaignID uuid.UUID, status string) (int64, error)
}

// Service interfaces

// CampaignService defines business logic for campaign operations
type CampaignService interface {
	CreateCampaign(ctx context.Context, input *CreateCampaignInput) (*Campaign, error)
	GetCampaign(ctx context.Context, id uuid.UUID) (*Campaign, error)
	UpdateCampaign(ctx context.Context, id uuid.UUID, input *UpdateCampaignInput, updatedBy uuid.UUID) error
	DeleteCampaign(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListCampaignsByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Campaign, error)
	ListActiveCampaignsByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Campaign, error)
	CountCampaigns(ctx context.Context) (int64, error)
	CountCampaignsByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
	CountActiveCampaignsByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
	
	// Campaign targeting and execution
	TargetClients(ctx context.Context, campaignID uuid.UUID) error
	SendCampaign(ctx context.Context, campaignID uuid.UUID) error
	GetCampaignAnalytics(ctx context.Context, campaignID uuid.UUID) (*CampaignAnalytics, error)
}

// CampaignAnalytics represents analytics data for a campaign
type CampaignAnalytics struct {
	CampaignID       uuid.UUID `json:"campaign_id"`
	TotalTargeted    int64     `json:"total_targeted"`
	TotalSent        int64     `json:"total_sent"`
	TotalOpened      int64     `json:"total_opened"`
	TotalClicked     int64     `json:"total_clicked"`
	TotalConverted   int64     `json:"total_converted"`
	TotalRevenue     float64   `json:"total_revenue"`
	OpenRate         float64   `json:"open_rate"`          // Opened / Sent
	ClickRate        float64   `json:"click_rate"`         // Clicked / Opened
	ConversionRate   float64   `json:"conversion_rate"`    // Converted / Clicked
	RevenuePerClient float64   `json:"revenue_per_client"` // Total Revenue / Converted
}

// JSON marshaling helpers for JSONB fields

func (ta TargetAudience) Value() (interface{}, error) {
	return json.Marshal(ta)
}

func (ta *TargetAudience) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, ta)
}

func (od OfferDetails) Value() (interface{}, error) {
	return json.Marshal(od)
}

func (od *OfferDetails) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, od)
}