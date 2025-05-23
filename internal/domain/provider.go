package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Provider represents a beauty service provider
type Provider struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	BusinessName     string     `json:"business_name"`
	Description      string     `json:"description"`
	Address          string     `json:"address"`
	City             string     `json:"city"`
	PostalCode       string     `json:"postal_code"`
	Country          string     `json:"country"`
	Website          string     `json:"website"`
	LogoURL          string     `json:"logo_url"`
	SubscriptionTier string     `json:"subscription_tier"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	DeletedBy        *uuid.UUID `json:"deleted_by,omitempty"`

	// Expanded relationships (populated by service when needed)
	User *User `json:"user,omitempty"`
}

// CreateProviderInput is the input for creating a provider
type CreateProviderInput struct {
	UserID           uuid.UUID `json:"user_id" validate:"required"`
	BusinessName     string    `json:"business_name" validate:"required"`
	Description      string    `json:"description"`
	Address          string    `json:"address"`
	City             string    `json:"city"`
	PostalCode       string    `json:"postal_code"`
	Country          string    `json:"country" validate:"required"`
	Website          string    `json:"website"`
	LogoURL          string    `json:"logo_url"`
	SubscriptionTier string    `json:"subscription_tier"`
}

// UpdateProviderInput is the input for updating a provider
type UpdateProviderInput struct {
	BusinessName     *string `json:"business_name"`
	Description      *string `json:"description"`
	Address          *string `json:"address"`
	City             *string `json:"city"`
	PostalCode       *string `json:"postal_code"`
	Country          *string `json:"country"`
	Website          *string `json:"website"`
	LogoURL          *string `json:"logo_url"`
	SubscriptionTier *string `json:"subscription_tier"`
}

// ProviderRepository defines the methods to interact with the provider data store
type ProviderRepository interface {
	Create(ctx context.Context, provider *Provider) error
	GetByID(ctx context.Context, id uuid.UUID) (*Provider, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Provider, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateProviderInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	List(ctx context.Context, page, pageSize int) ([]*Provider, error)
	Search(ctx context.Context, query string, page, pageSize int) ([]*Provider, error)
	Count(ctx context.Context) (int64, error)
}

// ProviderService defines the business logic for provider operations
type ProviderService interface {
	CreateProvider(ctx context.Context, input *CreateProviderInput) (*Provider, error)
	GetProvider(ctx context.Context, id uuid.UUID) (*Provider, error)
	GetProviderByUserID(ctx context.Context, userID uuid.UUID) (*Provider, error)
	UpdateProvider(ctx context.Context, id uuid.UUID, input *UpdateProviderInput, updatedBy uuid.UUID) error
	DeleteProvider(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListProviders(ctx context.Context, page, pageSize int) ([]*Provider, error)
	SearchProviders(ctx context.Context, query string, page, pageSize int) ([]*Provider, error)
	CountProviders(ctx context.Context) (int64, error)
}
