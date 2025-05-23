package models

import (
	"github.com/google/uuid"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	// UserRoleUser represents a regular user
	UserRoleUser UserRole = "user"
	// UserRoleAdmin represents an admin user
	UserRoleAdmin UserRole = "admin"
	// UserRoleOwner represents a provider (business owner)
	UserRoleOwner UserRole = "owner"
	// UserRoleStaff represents a staff member
	UserRoleStaff UserRole = "staff"
)

// User represents a user in the system
type User struct {
	BaseModel
	ClerkID   string   `gorm:"uniqueIndex;not null" json:"clerk_id"`
	Email     string   `gorm:"uniqueIndex;not null" json:"email"`
	FirstName string   `gorm:"not null" json:"first_name"`
	LastName  string   `gorm:"not null" json:"last_name"`
	Phone     string   `json:"phone"`
	Role      UserRole `gorm:"type:text;not null;default:'user'" json:"role"`
	IsActive  bool     `gorm:"not null;default:true" json:"is_active"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// UserConnectedAccount represents an authentication provider connected to a user
type UserConnectedAccount struct {
	BaseModel
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	ProviderType string    `gorm:"not null" json:"provider_type"` // e.g., "email", "google", "facebook", "apple"
	ProviderID   string    `gorm:"not null" json:"provider_id"`   // ID from the provider
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
}

// TableName overrides the table name
func (UserConnectedAccount) TableName() string {
	return "user_connected_accounts"
}
