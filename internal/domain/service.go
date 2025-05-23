package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ServiceCategory represents a category for services
type ServiceCategory struct {
	ID          uuid.UUID  `json:"id"`
	BusinessID  uuid.UUID  `json:"business_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	DeletedBy   *uuid.UUID `json:"deleted_by,omitempty"`
}

// Service represents a beauty service offered by a provider
type Service struct {
	ID          uuid.UUID  `json:"id"`
	BusinessID  uuid.UUID  `json:"business_id"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Duration    int        `json:"duration"` // in minutes
	Price       float64    `json:"price"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	DeletedBy   *uuid.UUID `json:"deleted_by,omitempty"`

	// Expanded relationships (populated by service when needed)
	Business *Business        `json:"business,omitempty"`
	Category *ServiceCategory `json:"category,omitempty"`
}

// CreateServiceInput is the input for creating a service
type CreateServiceInput struct {
	BusinessID  uuid.UUID  `json:"business_id" validate:"required"`
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description"`
	Duration    int        `json:"duration" validate:"required,min=1"`
	Price       float64    `json:"price" validate:"required,min=0"`
}

// UpdateServiceInput is the input for updating a service
type UpdateServiceInput struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Duration    *int       `json:"duration" validate:"omitempty,min=1"`
	Price       *float64   `json:"price" validate:"omitempty,min=0"`
}

// ServiceCategoryRepository defines methods for service category data store
type ServiceCategoryRepository interface {
	Create(ctx context.Context, category *ServiceCategory) error
	GetByID(ctx context.Context, id uuid.UUID) (*ServiceCategory, error)
	Update(ctx context.Context, id uuid.UUID, name, description string, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	List(ctx context.Context, page, pageSize int) ([]*ServiceCategory, error)
	Count(ctx context.Context) (int64, error)
}

// ServiceRepository defines methods for service data store
type ServiceRepository interface {
	Create(ctx context.Context, service *Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*Service, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateServiceInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Service, error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int) ([]*Service, error)
	Count(ctx context.Context) (int64, error)
	CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// ServiceCategoryService defines business logic for service category operations
type ServiceCategoryService interface {
	CreateCategory(ctx context.Context, name, description string) (*ServiceCategory, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*ServiceCategory, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, name, description string, updatedBy uuid.UUID) error
	DeleteCategory(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListCategories(ctx context.Context, page, pageSize int) ([]*ServiceCategory, error)
	CountCategories(ctx context.Context) (int64, error)
}

// ServiceService defines business logic for service operations
type ServiceService interface {
	CreateService(ctx context.Context, input *CreateServiceInput) (*Service, error)
	GetService(ctx context.Context, id uuid.UUID) (*Service, error)
	UpdateService(ctx context.Context, id uuid.UUID, input *UpdateServiceInput, updatedBy uuid.UUID) error
	DeleteService(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListServicesByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*Service, error)
	ListServicesByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int) ([]*Service, error)
	CountServices(ctx context.Context) (int64, error)
	CountServicesByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}
