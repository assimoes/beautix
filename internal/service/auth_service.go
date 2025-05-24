package service

import (
	"context"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/dto"
	"github.com/assimoes/beautix/internal/infrastructure/auth"
	"github.com/assimoes/beautix/pkg/errors"
	"github.com/assimoes/beautix/pkg/utils"
	"gorm.io/gorm"
)

// AuthService provides authentication-related operations
type AuthService interface {
	RegisterWithClerk(ctx context.Context, clerkUserData dto.ClerkUserDTO) (*dto.UserResponseDTO, error)
	SyncClerkUser(ctx context.Context, clerkID string, userData dto.ClerkUserDTO) (*dto.UserResponseDTO, error)
	GetCurrentUser(ctx context.Context, clerkID string) (*dto.UserResponseDTO, error)
	VerifyToken(ctx context.Context, token string) (*dto.UserResponseDTO, error)
	GetUserFromContext(ctx context.Context) (*dto.UserResponseDTO, error)
}

// authServiceImpl implements the AuthService interface
type authServiceImpl struct {
	userRepo    domain.UserRepository
	clerkClient *auth.ClerkClient
	db          *gorm.DB
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo domain.UserRepository,
	clerkClient *auth.ClerkClient,
	db *gorm.DB,
) AuthService {
	return &authServiceImpl{
		userRepo:    userRepo,
		clerkClient: clerkClient,
		db:          db,
	}
}

// RegisterWithClerk registers a new user using Clerk data
func (s *authServiceImpl) RegisterWithClerk(ctx context.Context, clerkUserData dto.ClerkUserDTO) (*dto.UserResponseDTO, error) {
	// Validate the clerk user data
	if err := utils.ValidateStruct(clerkUserData); err != nil {
		return nil, err
	}

	// Check if user already exists by email
	existingUser, err := s.userRepo.FindByEmail(ctx, clerkUserData.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.NewInternalError("failed to check existing user", err)
	}

	if existingUser != nil {
		// User exists, update Clerk ID if not set
		if existingUser.ClerkID == nil || *existingUser.ClerkID != clerkUserData.ClerkID {
			if err := s.userRepo.UpdateClerkID(ctx, existingUser.ID, clerkUserData.ClerkID); err != nil {
				return nil, errors.NewInternalError("failed to update clerk ID", err)
			}
			existingUser.ClerkID = &clerkUserData.ClerkID
		}
		return dto.ToUserResponseDTO(existingUser), nil
	}

	// Create new user
	createDTO := clerkUserData.ToCreateUserDTO()
	user := &domain.User{
		Email:     utils.NormalizeEmail(createDTO.Email),
		ClerkID:   createDTO.ClerkID,
		FirstName: utils.NormalizeName(createDTO.FirstName),
		LastName:  utils.NormalizeName(createDTO.LastName),
		Phone:     createDTO.Phone,
		IsActive:  true,
	}

	// Set audit fields
	user.SetAuditFields(nil) // No user context for registration

	// Validate the user
	if err := user.Validate(); err != nil {
		return nil, errors.NewValidationError("invalid user data", err)
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewInternalError("failed to create user", err)
	}

	return dto.ToUserResponseDTO(user), nil
}

// SyncClerkUser synchronizes user data from Clerk
func (s *authServiceImpl) SyncClerkUser(ctx context.Context, clerkID string, userData dto.ClerkUserDTO) (*dto.UserResponseDTO, error) {
	// Find user by Clerk ID
	user, err := s.userRepo.FindByClerkID(ctx, clerkID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create new one
			return s.RegisterWithClerk(ctx, userData)
		}
		return nil, errors.NewInternalError("failed to find user", err)
	}

	// Update user data
	hasChanges := false

	normalizedEmail := utils.NormalizeEmail(userData.Email)
	if user.Email != normalizedEmail {
		user.Email = normalizedEmail
		hasChanges = true
	}

	normalizedFirstName := utils.NormalizeName(userData.FirstName)
	if user.FirstName != normalizedFirstName {
		user.FirstName = normalizedFirstName
		hasChanges = true
	}

	normalizedLastName := utils.NormalizeName(userData.LastName)
	if user.LastName != normalizedLastName {
		user.LastName = normalizedLastName
		hasChanges = true
	}

	phoneMatches := (user.Phone == nil && userData.Phone == "") ||
		(user.Phone != nil && userData.Phone != "" && *user.Phone == userData.Phone)
	if !phoneMatches {
		if userData.Phone == "" {
			user.Phone = nil
		} else {
			user.Phone = &userData.Phone
		}
		hasChanges = true
	}

	// Save changes if any
	if hasChanges {
		userID := GetUserIDFromContext(ctx)
		user.SetAuditFields(userID)

		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, errors.NewInternalError("failed to update user", err)
		}
	}

	return dto.ToUserResponseDTO(user), nil
}

// GetCurrentUser retrieves the current user by Clerk ID
func (s *authServiceImpl) GetCurrentUser(ctx context.Context, clerkID string) (*dto.UserResponseDTO, error) {
	user, err := s.userRepo.FindByClerkID(ctx, clerkID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("user")
		}
		return nil, errors.NewInternalError("failed to get user", err)
	}

	if !user.IsActive {
		return nil, errors.NewUnauthorizedError("user account is inactive")
	}

	return dto.ToUserResponseDTO(user), nil
}

// VerifyToken verifies a Clerk token and returns the user
func (s *authServiceImpl) VerifyToken(ctx context.Context, token string) (*dto.UserResponseDTO, error) {
	// Verify token with Clerk
	clerkUser, err := s.clerkClient.VerifyToken(ctx, token)
	if err != nil {
		return nil, err // Error already wrapped by Clerk client
	}

	// Get or sync user from our database
	user, err := s.userRepo.FindByClerkID(ctx, clerkUser.ClerkID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create from Clerk data
			clerkUserDTO := dto.ClerkUserDTO{
				ClerkID:   clerkUser.ClerkID,
				Email:     clerkUser.Email,
				FirstName: clerkUser.FirstName,
				LastName:  clerkUser.LastName,
				Phone:     clerkUser.Phone,
			}
			return s.RegisterWithClerk(ctx, clerkUserDTO)
		}
		return nil, errors.NewInternalError("failed to get user", err)
	}

	if !user.IsActive {
		return nil, errors.NewUnauthorizedError("user account is inactive")
	}

	return dto.ToUserResponseDTO(user), nil
}

// GetUserFromContext extracts user information from context
func (s *authServiceImpl) GetUserFromContext(ctx context.Context) (*dto.UserResponseDTO, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == nil {
		return nil, errors.NewUnauthorizedError("no user in context")
	}

	user, err := s.userRepo.GetByID(ctx, *userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("user")
		}
		return nil, errors.NewInternalError("failed to get user", err)
	}

	return dto.ToUserResponseDTO(user), nil
}
