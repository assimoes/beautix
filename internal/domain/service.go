package domain

import (
	"context"
	"github.com/shopspring/decimal"
)

// Service represents a service offered by a business
type Service struct {
	BaseModel
	BusinessID   string          `gorm:"not null;type:uuid;index" json:"business_id"`
	CategoryID   *string         `gorm:"type:uuid;index" json:"category_id,omitempty"`
	Name         string          `gorm:"not null;size:200" json:"name"`
	Description  *string         `gorm:"type:text" json:"description,omitempty"`
	Duration     int             `gorm:"not null;default:30" json:"duration"` // Duration in minutes
	Price        decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	IsActive     bool            `gorm:"not null;default:true" json:"is_active"`
	DisplayOrder int             `gorm:"not null;default:0" json:"display_order"`
	PreparationTime *int         `gorm:"default:0" json:"preparation_time,omitempty"` // Buffer time before in minutes
	CleanupTime     *int         `gorm:"default:0" json:"cleanup_time,omitempty"`     // Buffer time after in minutes
	MaxAdvanceBooking *int       `gorm:"" json:"max_advance_booking,omitempty"`       // Days in advance
	MinAdvanceBooking *int       `gorm:"default:0" json:"min_advance_booking,omitempty"` // Hours in advance
	RequiresDeposit   bool       `gorm:"not null;default:false" json:"requires_deposit"`
	DepositAmount     *decimal.Decimal `gorm:"type:decimal(10,2)" json:"deposit_amount,omitempty"`

	// Relationships
	Business Business        `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
	Category *ServiceCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
}

// TableName returns the table name for Service
func (Service) TableName() string { return "services" }

// Validate validates the service model
func (s *Service) Validate() error {
	if s.BusinessID == "" {
		return ErrValidation
	}
	if s.Name == "" {
		return ErrValidation
	}
	if s.Duration <= 0 {
		return ErrValidation
	}
	if s.Price.IsNegative() {
		return ErrValidation
	}
	if s.DepositAmount != nil && s.DepositAmount.IsNegative() {
		return ErrValidation
	}
	return nil
}

// GetFullName returns the full display name including category
func (s *Service) GetFullName() string {
	if s.Category != nil {
		return s.Category.Name + " - " + s.Name
	}
	return s.Name
}

// GetTotalDuration returns the total duration including preparation and cleanup
func (s *Service) GetTotalDuration() int {
	total := s.Duration
	if s.PreparationTime != nil {
		total += *s.PreparationTime
	}
	if s.CleanupTime != nil {
		total += *s.CleanupTime
	}
	return total
}

// ServiceCategory represents a category of services
type ServiceCategory struct {
	BaseModel
	BusinessID   string  `gorm:"not null;type:uuid;index" json:"business_id"`
	Name         string  `gorm:"not null;size:100" json:"name"`
	Description  *string `gorm:"type:text" json:"description,omitempty"`
	DisplayOrder int     `gorm:"not null;default:0" json:"display_order"`
	IsActive     bool    `gorm:"not null;default:true" json:"is_active"`
	ColorCode    *string `gorm:"size:7" json:"color_code,omitempty"` // Hex color code for UI

	// Relationships
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
	Services []Service `gorm:"foreignKey:CategoryID" json:"services,omitempty"`
}

// TableName returns the table name for ServiceCategory
func (ServiceCategory) TableName() string { return "service_categories" }

// Validate validates the service category model
func (sc *ServiceCategory) Validate() error {
	if sc.BusinessID == "" {
		return ErrValidation
	}
	if sc.Name == "" {
		return ErrValidation
	}
	return nil
}

// ServiceRepository defines the repository interface for Service
type ServiceRepository interface {
	BaseRepository[Service]
	FindByBusinessID(ctx context.Context, businessID string) ([]*Service, error)
	FindByCategory(ctx context.Context, businessID, categoryID string) ([]*Service, error)
	FindActiveByBusiness(ctx context.Context, businessID string) ([]*Service, error)
	UpdatePricing(ctx context.Context, serviceID string, price decimal.Decimal) error
	ReorderServices(ctx context.Context, businessID string, serviceOrders []ServiceOrder) error
	ExistsByNameAndBusiness(ctx context.Context, name, businessID string) (bool, error)
}

// ServiceCategoryRepository defines the repository interface for ServiceCategory
type ServiceCategoryRepository interface {
	BaseRepository[ServiceCategory]
	FindByBusinessID(ctx context.Context, businessID string) ([]*ServiceCategory, error)
	GetByDisplayOrder(ctx context.Context, businessID string) ([]*ServiceCategory, error)
	ExistsByNameAndBusiness(ctx context.Context, name, businessID string) (bool, error)
	ReorderCategories(ctx context.Context, businessID string, categoryOrders []CategoryOrder) error
}

// Helper types for repository methods
type ServiceOrder struct {
	ServiceID    string `json:"service_id"`
	DisplayOrder int    `json:"display_order"`
}

type CategoryOrder struct {
	CategoryID   string `json:"category_id"`
	DisplayOrder int    `json:"display_order"`
}