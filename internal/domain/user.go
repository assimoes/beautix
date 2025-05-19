package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Phone        string     `json:"phone"`
	Role         string     `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	DeletedBy    *uuid.UUID `json:"deleted_by,omitempty"`
}

// CreateUserInput is the input for creating a user
type CreateUserInput struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone"`
	Role      string `json:"role" validate:"required,oneof=user admin provider"`
}

// UpdateUserInput is the input for updating a user
type UpdateUserInput struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=8"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
}

// UserRepository defines the methods to interact with the user data store
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateUserInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	List(ctx context.Context, page, pageSize int) ([]*User, error)
	Count(ctx context.Context) (int64, error)
}

// UserService defines the business logic for user operations
type UserService interface {
	CreateUser(ctx context.Context, input *CreateUserInput) (*User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, input *UpdateUserInput, updatedBy uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*User, error)
	CountUsers(ctx context.Context) (int64, error)
	Authenticate(ctx context.Context, email, password string) (*User, error)
	GenerateToken(ctx context.Context, user *User) (string, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
}