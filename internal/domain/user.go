package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	UserID             uuid.UUID  `json:"user_id"`
	Email              string     `json:"email"`
	ClerkID            uuid.UUID  `json:"clerk_id"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	Phone              string     `json:"phone,omitempty"`
	Role               string     `json:"role"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	LastLogin          *time.Time `json:"last_login,omitempty"`
	IsActive           bool       `json:"is_active"`
	EmailVerified      bool       `json:"email_verified"`
	LanguagePreference string     `json:"language_preference"`
}

// CreateUserInput is the input for creating a user
type CreateUserInput struct {
	Email              string `json:"email" validate:"required,email"`
	Password           string `json:"password" validate:"required,min=8"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	Phone              string `json:"phone,omitempty"`
	Role               string `json:"role" validate:"required,oneof=admin owner staff user"`
	LanguagePreference string `json:"language_preference,omitempty"`
}

// UpdateUserInput is the input for updating a user
type UpdateUserInput struct {
	Email              *string `json:"email,omitempty" validate:"omitempty,email"`
	Password           *string `json:"password,omitempty" validate:"omitempty,min=8"`
	FirstName          *string `json:"first_name,omitempty"`
	LastName           *string `json:"last_name,omitempty"`
	Phone              *string `json:"phone,omitempty"`
	IsActive           *bool   `json:"is_active,omitempty"`
	EmailVerified      *bool   `json:"email_verified,omitempty"`
	LanguagePreference *string `json:"language_preference,omitempty"`
}

// UserRepository defines the methods to interact with the user data store
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateUserInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int) ([]*User, error)
	Count(ctx context.Context) (int64, error)
}

// UserService defines the business logic for user operations
type UserService interface {
	CreateUser(ctx context.Context, input *CreateUserInput) (*User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, input *UpdateUserInput) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*User, error)
	CountUsers(ctx context.Context) (int64, error)
	Authenticate(ctx context.Context, email, password string) (*User, error)
	GenerateToken(ctx context.Context, user *User, businessID *uuid.UUID) (string, error)
	ValidateToken(ctx context.Context, token string) (*User, *uuid.UUID, error)
}

// TenantContext holds the tenant context for multi-tenancy
type TenantContext struct {
	BusinessID uuid.UUID
}

// Session represents a user authentication session
type Session struct {
	SessionID    uuid.UUID  `json:"session_id"`
	UserID       uuid.UUID  `json:"user_id"`
	BusinessID   *uuid.UUID `json:"business_id,omitempty"`
	Token        string     `json:"token"`
	IPAddress    string     `json:"ip_address,omitempty"`
	UserAgent    string     `json:"user_agent,omitempty"`
	ExpiresAt    time.Time  `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at"`
	LastActivity time.Time  `json:"last_activity"`
}

// SessionRepository defines methods for session management
type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByToken(ctx context.Context, token string) (*Session, error)
	Update(ctx context.Context, id uuid.UUID, lastActivity time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// PermissionScope represents a permission for a user in a business
type PermissionScope struct {
	ScopeID    uuid.UUID `json:"scope_id"`
	UserID     uuid.UUID `json:"user_id"`
	BusinessID uuid.UUID `json:"business_id"`
	Resource   string    `json:"resource"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"created_at"`
}

// PermissionScopeRepository defines methods for permission scope management
type PermissionScopeRepository interface {
	Create(ctx context.Context, scope *PermissionScope) error
	GetByUserAndBusiness(ctx context.Context, userID, businessID uuid.UUID) ([]*PermissionScope, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserAndBusiness(ctx context.Context, userID, businessID uuid.UUID) error
}

// PermissionScopeService defines business logic for permission scope operations
type PermissionScopeService interface {
	CreatePermissionScope(ctx context.Context, userID, businessID uuid.UUID, resource, action string) (*PermissionScope, error)
	GetUserPermissions(ctx context.Context, userID, businessID uuid.UUID) ([]*PermissionScope, error)
	HasPermission(ctx context.Context, userID, businessID uuid.UUID, resource, action string) (bool, error)
	RevokePermission(ctx context.Context, userID, businessID uuid.UUID, resource, action string) error
	RevokeAllPermissions(ctx context.Context, userID, businessID uuid.UUID) error
}
