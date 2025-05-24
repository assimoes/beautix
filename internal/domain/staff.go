package domain

import (
	"context"
	"time"
)

// BusinessRole represents the role a user has within a specific business
type BusinessRole string

const (
	BusinessRoleOwner     BusinessRole = "owner"
	BusinessRoleManager   BusinessRole = "manager"
	BusinessRoleEmployee  BusinessRole = "employee"
	BusinessRoleAssistant BusinessRole = "assistant"
)

// Staff represents a user's role and permissions within a specific business
type Staff struct {
	BaseModel
	BusinessID   string       `gorm:"not null;type:uuid;index" json:"business_id"`
	UserID       string       `gorm:"not null;type:uuid;index" json:"user_id"`
	Role         BusinessRole `gorm:"not null;size:20;check:role IN ('owner','manager','employee','assistant')" json:"role"`
	IsActive     bool         `gorm:"not null;default:true" json:"is_active"`
	Permissions  *string      `gorm:"type:jsonb;default:'{}'" json:"permissions,omitempty"` // JSON object with specific permissions
	StartDate    *time.Time   `gorm:"" json:"start_date,omitempty"`
	EndDate      *time.Time   `gorm:"" json:"end_date,omitempty"`

	// Relationships
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
	User     User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

// TableName returns the table name for Staff
func (Staff) TableName() string { return "staff" }

// IsOwner returns true if the staff member has owner role
func (s Staff) IsOwner() bool {
	return s.Role == BusinessRoleOwner
}

// IsManager returns true if the staff member has manager role
func (s Staff) IsManager() bool {
	return s.Role == BusinessRoleManager
}

// CanManage returns true if the staff member can manage the business
func (s Staff) CanManage() bool {
	return s.Role == BusinessRoleOwner || s.Role == BusinessRoleManager
}

// Validate validates the staff model
func (s *Staff) Validate() error {
	if s.BusinessID == "" {
		return ErrValidation
	}
	if s.UserID == "" {
		return ErrValidation
	}
	if s.Role != BusinessRoleOwner && s.Role != BusinessRoleManager && 
		s.Role != BusinessRoleEmployee && s.Role != BusinessRoleAssistant {
		return ErrValidation
	}
	return nil
}

// StaffRepository defines the repository interface for Staff
type StaffRepository interface {
	BaseRepository[Staff]
	FindByBusinessID(ctx context.Context, businessID string) ([]*Staff, error)
	FindByUserID(ctx context.Context, userID string) ([]*Staff, error)
	FindByBusinessAndUser(ctx context.Context, businessID, userID string) (*Staff, error)
	FindActiveByBusinessID(ctx context.Context, businessID string) ([]*Staff, error)
	FindByRole(ctx context.Context, businessID string, role BusinessRole) ([]*Staff, error)
	UpdateRole(ctx context.Context, businessID, userID string, role BusinessRole) error
	DeactivateStaff(ctx context.Context, businessID, userID string) error
}