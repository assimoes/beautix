package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Client represents a client of a service provider
type Client struct {
	ID         uuid.UUID  `json:"id"`
	UserID     *uuid.UUID `json:"user_id,omitempty"` // Can be null if client doesn't have account
	BusinessID uuid.UUID  `json:"business_id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Notes      string     `json:"notes"`
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	UpdatedBy  *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	DeletedBy  *uuid.UUID `json:"deleted_by,omitempty"`

	// Expanded relationships (populated by service when needed)
	User     *User     `json:"user,omitempty"`
	Business *Business `json:"business,omitempty"`
}

// CreateClientInput is the input for creating a client
type CreateClientInput struct {
	UserID     *uuid.UUID `json:"user_id"`
	BusinessID uuid.UUID  `json:"business_id" validate:"required"`
	FirstName  string     `json:"first_name" validate:"required"`
	LastName   string     `json:"last_name" validate:"required"`
	Email      string     `json:"email" validate:"omitempty,email"`
	Phone      string     `json:"phone"`
	Notes      string     `json:"notes"`
}

// UpdateClientInput is the input for updating a client
type UpdateClientInput struct {
	UserID    *uuid.UUID `json:"user_id"`
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	Email     *string    `json:"email" validate:"omitempty,email"`
	Phone     *string    `json:"phone"`
	Notes     *string    `json:"notes"`
}

// ClientRepository defines methods for client data store
type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	GetByID(ctx context.Context, id uuid.UUID) (*Client, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Client, error)
	GetByBusinessAndEmail(ctx context.Context, businessID uuid.UUID, email string) (*Client, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateClientInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Client, error)
	Search(ctx context.Context, businessID uuid.UUID, query string, page, pageSize int) ([]*Client, error)
	Count(ctx context.Context) (int64, error)
	CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// ClientService defines business logic for client operations
type ClientService interface {
	CreateClient(ctx context.Context, input *CreateClientInput) (*Client, error)
	GetClient(ctx context.Context, id uuid.UUID) (*Client, error)
	GetClientsByUserID(ctx context.Context, userID uuid.UUID) ([]*Client, error)
	GetClientByBusinessAndEmail(ctx context.Context, businessID uuid.UUID, email string) (*Client, error)
	UpdateClient(ctx context.Context, id uuid.UUID, input *UpdateClientInput, updatedBy uuid.UUID) error
	DeleteClient(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListClientsByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Client, error)
	SearchClients(ctx context.Context, businessID uuid.UUID, query string, page, pageSize int) ([]*Client, error)
	CountClients(ctx context.Context) (int64, error)
	CountClientsByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}
