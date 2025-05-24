package domain

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Common domain errors
var (
	ErrValidation = errors.New("validation error")
)

// User represents a user in the system
type User struct {
	BaseModel
	Email     string   `gorm:"not null;uniqueIndex" json:"email"`
	ClerkID   *string  `gorm:"uniqueIndex" json:"clerk_id,omitempty"`
	FirstName string   `gorm:"not null;size:100" json:"first_name"`
	LastName  string   `gorm:"not null;size:100" json:"last_name"`
	Phone     *string  `gorm:"size:50" json:"phone,omitempty"`
	IsActive  bool     `gorm:"not null;default:true" json:"is_active"`

	// Relationships
	Businesses           []Business           `gorm:"foreignKey:UserID" json:"businesses,omitempty"`
	StaffPositions       []Staff              `gorm:"foreignKey:UserID" json:"staff_positions,omitempty"` // Staff positions this user holds at businesses
	ConnectedAccounts    []UserConnectedAccount `gorm:"foreignKey:UserID" json:"connected_accounts,omitempty"`
}

// UserConnectedAccount represents external account connections (OAuth providers)
type UserConnectedAccount struct {
	BaseModel
	UserID       string `gorm:"not null;type:uuid" json:"user_id"`
	ProviderType string `gorm:"not null;size:50" json:"provider_type"`
	ProviderID   string `gorm:"not null;size:255" json:"provider_id"`
	IsActive     bool   `gorm:"not null;default:true" json:"is_active"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// TableName returns the table name for UserConnectedAccount
func (UserConnectedAccount) TableName() string {
	return "user_connected_accounts"
}

// GetFullName returns the user's full name
func (u User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsOwnerOfBusiness returns true if the user is the owner of the given business
func (u User) IsOwnerOfBusiness(businessID string) bool {
	for _, business := range u.Businesses {
		if business.ID == businessID {
			return true
		}
	}
	return false
}

// GetBusinessRole returns the user's role in the given business
func (u User) GetBusinessRole(businessID string) *BusinessRole {
	// Check if user is the business owner
	for _, business := range u.Businesses {
		if business.ID == businessID {
			role := BusinessRoleOwner
			return &role
		}
	}

	// Check if user has a staff position in the business
	for _, staffPosition := range u.StaffPositions {
		if staffPosition.BusinessID == businessID && staffPosition.IsActive {
			return &staffPosition.Role
		}
	}

	return nil
}

// CanAccessBusiness returns true if the user can access the given business
func (u User) CanAccessBusiness(businessID string) bool {
	// Check if user owns the business
	if u.IsOwnerOfBusiness(businessID) {
		return true
	}

	// Check if user has active staff position in the business
	for _, staffPosition := range u.StaffPositions {
		if staffPosition.BusinessID == businessID && staffPosition.IsActive {
			return true
		}
	}

	return false
}

// HasPermissionInBusiness checks if user has specific permission in a business
func (u User) HasPermissionInBusiness(businessID string, permission string) bool {
	role := u.GetBusinessRole(businessID)
	if role == nil {
		return false
	}

	// Owners have all permissions
	if *role == BusinessRoleOwner {
		return true
	}

	// For other roles, implement permission checking logic
	// This can be extended based on specific permission requirements
	return false
}

// Validate validates the user model
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrValidation
	}
	if u.FirstName == "" {
		return ErrValidation
	}
	if u.LastName == "" {
		return ErrValidation
	}
	return nil
}

// BeforeCreate is called before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// No role-related setup needed anymore
	return nil
}

// UserRepository defines the repository interface for User
type UserRepository interface {
	BaseRepository[User]
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByClerkID(ctx context.Context, clerkID string) (*User, error)
	UpdateClerkID(ctx context.Context, userID, clerkID string) error
	GetWithBusinesses(ctx context.Context, userID string) (*User, error)
	SearchUsers(ctx context.Context, query string, limit int) ([]*User, error)
}

// Note: Service interface is defined in the service layer to avoid circular dependencies