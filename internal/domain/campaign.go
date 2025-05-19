package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Campaign represents a marketing campaign or promotion
type Campaign struct {
	ID           uuid.UUID  `json:"id"`
	ProviderID   uuid.UUID  `json:"provider_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      time.Time  `json:"end_date"`
	DiscountType string     `json:"discount_type"` // percentage, fixed_amount
	DiscountValue float64    `json:"discount_value"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	DeletedBy    *uuid.UUID `json:"deleted_by,omitempty"`
	
	// Expanded relationships (populated by service when needed)
	Provider *Provider   `json:"provider,omitempty"`
	Services []*Service  `json:"services,omitempty"`
}

// CampaignServiceMapping represents a service included in a campaign
type CampaignServiceMapping struct {
	CampaignID uuid.UUID  `json:"campaign_id"`
	ServiceID  uuid.UUID  `json:"service_id"`
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  *uuid.UUID `json:"created_by,omitempty"`
	
	// Expanded relationships (populated by service when needed)
	Campaign *Campaign `json:"campaign,omitempty"`
	Service  *Service  `json:"service,omitempty"`
}

// CreateCampaignInput is the input for creating a campaign
type CreateCampaignInput struct {
	ProviderID    uuid.UUID  `json:"provider_id" validate:"required"`
	Name          string     `json:"name" validate:"required"`
	Description   string     `json:"description"`
	StartDate     time.Time  `json:"start_date" validate:"required"`
	EndDate       time.Time  `json:"end_date" validate:"required,gtfield=StartDate"`
	DiscountType  string     `json:"discount_type" validate:"required,oneof=percentage fixed_amount"`
	DiscountValue float64    `json:"discount_value" validate:"required,min=0"`
	ServiceIDs    []uuid.UUID `json:"service_ids" validate:"required,min=1"`
}

// UpdateCampaignInput is the input for updating a campaign
type UpdateCampaignInput struct {
	Name          *string    `json:"name"`
	Description   *string    `json:"description"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	DiscountType  *string    `json:"discount_type" validate:"omitempty,oneof=percentage fixed_amount"`
	DiscountValue *float64   `json:"discount_value" validate:"omitempty,min=0"`
	ServiceIDs    []uuid.UUID `json:"service_ids"`
}

// CampaignRepository defines methods for campaign data store
type CampaignRepository interface {
	Create(ctx context.Context, campaign *Campaign) error
	AddServices(ctx context.Context, campaignID uuid.UUID, serviceIDs []uuid.UUID, createdBy uuid.UUID) error
	RemoveServices(ctx context.Context, campaignID uuid.UUID, serviceIDs []uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Campaign, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateCampaignInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) ([]*Campaign, error)
	ListActiveByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) ([]*Campaign, error)
	ListServiceMappingsByID(ctx context.Context, campaignID uuid.UUID) ([]*CampaignServiceMapping, error)
	ListServicesByID(ctx context.Context, campaignID uuid.UUID) ([]*Service, error)
	Count(ctx context.Context) (int64, error)
	CountByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
	CountActiveByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
}

// CampaignService defines business logic for campaign operations
type CampaignService interface {
	CreateCampaign(ctx context.Context, input *CreateCampaignInput) (*Campaign, error)
	GetCampaign(ctx context.Context, id uuid.UUID) (*Campaign, error)
	UpdateCampaign(ctx context.Context, id uuid.UUID, input *UpdateCampaignInput, updatedBy uuid.UUID) error
	DeleteCampaign(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListCampaignsByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) ([]*Campaign, error)
	ListActiveCampaignsByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) ([]*Campaign, error)
	ListCampaignServices(ctx context.Context, campaignID uuid.UUID) ([]*Service, error)
	AddServicesToCompaign(ctx context.Context, campaignID uuid.UUID, serviceIDs []uuid.UUID, updatedBy uuid.UUID) error
	RemoveServicesFromCampaign(ctx context.Context, campaignID uuid.UUID, serviceIDs []uuid.UUID) error
	CountCampaigns(ctx context.Context) (int64, error)
	CountCampaignsByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
	CountActiveCampaignsByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
}