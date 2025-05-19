package mock

import (
	"context"
	"errors"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// UserService is a mock implementation of domain.UserService
type UserService struct {
	users []*domain.User
}

// NewUserService creates a new mock user service with some test data
func NewUserService() *UserService {
	// Create some mock users for testing
	users := []*domain.User{
		{
			ID:           uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			Email:        "admin@example.com",
			PasswordHash: "$2a$10$zY5MdNIX8vxMUiZjrcHGn.jgRmrTyQA9QZG3HJURcfxjr9QfDfqyW", // "password"
			FirstName:    "Admin",
			LastName:     "User",
			Phone:        "+351910000000",
			Role:         "admin",
			CreatedAt:    time.Now().Add(-24 * time.Hour),
		},
		{
			ID:           uuid.MustParse("5f6e8b3a-2d98-4d80-8f4e-6c79f236a3d2"),
			Email:        "provider@example.com",
			PasswordHash: "$2a$10$zY5MdNIX8vxMUiZjrcHGn.jgRmrTyQA9QZG3HJURcfxjr9QfDfqyW", // "password"
			FirstName:    "Provider",
			LastName:     "User",
			Phone:        "+351910000001",
			Role:         "provider",
			CreatedAt:    time.Now().Add(-12 * time.Hour),
		},
		{
			ID:           uuid.MustParse("9a8f7e6d-5c4b-3a2d-1e0f-9b8a7c6d5e4f"),
			Email:        "client@example.com",
			PasswordHash: "$2a$10$zY5MdNIX8vxMUiZjrcHGn.jgRmrTyQA9QZG3HJURcfxjr9QfDfqyW", // "password"
			FirstName:    "Client",
			LastName:     "User",
			Phone:        "+351910000002",
			Role:         "user",
			CreatedAt:    time.Now().Add(-6 * time.Hour),
		},
	}

	return &UserService{
		users: users,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, input *domain.CreateUserInput) (*domain.User, error) {
	// Check if email already exists
	for _, user := range s.users {
		if user.Email == input.Email {
			return nil, errors.New("email already exists")
		}
	}

	// Create a new user
	user := &domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: input.Password, // In a real application, this would be hashed
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Phone:        input.Phone,
		Role:         input.Role,
		CreatedAt:    time.Now(),
	}

	// Add the user to the mock database
	s.users = append(s.users, user)

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, input *domain.UpdateUserInput, updatedBy uuid.UUID) error {
	for i, user := range s.users {
		if user.ID == id {
			if input.Email != nil {
				user.Email = *input.Email
			}
			if input.Password != nil {
				user.PasswordHash = *input.Password // In a real application, this would be hashed
			}
			if input.FirstName != nil {
				user.FirstName = *input.FirstName
			}
			if input.LastName != nil {
				user.LastName = *input.LastName
			}
			if input.Phone != nil {
				user.Phone = *input.Phone
			}

			now := time.Now()
			user.UpdatedAt = &now
			user.UpdatedBy = &updatedBy

			s.users[i] = user
			return nil
		}
	}

	return errors.New("user not found")
}

// DeleteUser marks a user as deleted
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	for i, user := range s.users {
		if user.ID == id {
			now := time.Now()
			user.DeletedAt = &now
			user.DeletedBy = &deletedBy

			s.users[i] = user
			return nil
		}
	}

	return errors.New("user not found")
}

// ListUsers retrieves a list of users with pagination
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) ([]*domain.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(s.users) {
		return []*domain.User{}, nil
	}

	if end > len(s.users) {
		end = len(s.users)
	}

	return s.users[start:end], nil
}

// CountUsers returns the total number of users
func (s *UserService) CountUsers(ctx context.Context) (int64, error) {
	return int64(len(s.users)), nil
}

// Authenticate authenticates a user with email and password
func (s *UserService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	// In a real application, we would verify the password hash
	for _, user := range s.users {
		if user.Email == email {
			// Mock implementation just checks if password is "password"
			if password == "password" {
				return user, nil
			}
			return nil, errors.New("invalid password")
		}
	}

	return nil, errors.New("user not found")
}

// GenerateToken generates a JWT token for a user
func (s *UserService) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	// Create claims with user info
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	// In a real application, this would be a secure key from configuration
	tokenString, err := token.SignedString([]byte("mock-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *UserService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		// Return secret key
		return []byte("mock-secret-key"), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Extract user ID
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user ID")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	return s.GetUser(ctx, userID)
}