package service

import (
	"context"
	"errors"
	"strings"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/dto"
	"github.com/assimoes/beautix/internal/infrastructure/validation"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// UserService defines the service interface for User operations
type UserService interface {
	domain.BaseService[dto.CreateUserDTO, dto.UpdateUserDTO, dto.UserResponseDTO]
	GetByEmail(ctx context.Context, email string) (*dto.UserResponseDTO, error)
	GetByClerkID(ctx context.Context, clerkID string) (*dto.UserResponseDTO, error)
	UpdateClerkID(ctx context.Context, userID, clerkID string) error
	GetWithBusinesses(ctx context.Context, userID string) (*dto.UserWithBusinessesDTO, error)
	SearchUsers(ctx context.Context, query string, limit int) ([]*dto.UserResponseDTO, error)
	ActivateUser(ctx context.Context, userID string) error
	DeactivateUser(ctx context.Context, userID string) error
}

// userServiceImpl implements the UserService interface
type userServiceImpl struct {
	*BaseServiceImpl[domain.User, dto.CreateUserDTO, dto.UpdateUserDTO, dto.UserResponseDTO]
	userRepo     domain.UserRepository
	businessRepo domain.BusinessRepository
	staffRepo    domain.StaffRepository
	validator    *validator.Validate
}

// NewUserService creates a new user service
func NewUserService(userRepo domain.UserRepository, businessRepo domain.BusinessRepository, staffRepo domain.StaffRepository, validator *validator.Validate) UserService {
	return &userServiceImpl{
		BaseServiceImpl: NewBaseService(
			userRepo,
			func(createDTO dto.CreateUserDTO) (*domain.User, error) {
				user := &domain.User{
					Email:     strings.ToLower(createDTO.Email),
					ClerkID:   createDTO.ClerkID,
					FirstName: strings.TrimSpace(createDTO.FirstName),
					LastName:  strings.TrimSpace(createDTO.LastName),
					Phone:     createDTO.Phone,
					IsActive:  true,
				}
				return user, user.Validate()
			},
			func(entity *domain.User, updateDTO dto.UpdateUserDTO) error {
				if updateDTO.FirstName != nil {
					entity.FirstName = strings.TrimSpace(*updateDTO.FirstName)
				}
				if updateDTO.LastName != nil {
					entity.LastName = strings.TrimSpace(*updateDTO.LastName)
				}
				if updateDTO.Phone != nil {
					entity.Phone = updateDTO.Phone
				}
				if updateDTO.IsActive != nil {
					entity.IsActive = *updateDTO.IsActive
				}
				return entity.Validate()
			},
			func(entity *domain.User) *dto.UserResponseDTO {
				return dto.ToUserResponseDTO(entity)
			},
		),
		userRepo:     userRepo,
		businessRepo: businessRepo,
		staffRepo:    staffRepo,
		validator:    validator,
	}
}

// Create creates a new user and automatically creates a default business with the user as owner
func (s *userServiceImpl) Create(ctx context.Context, createDTO dto.CreateUserDTO) (*dto.UserResponseDTO, error) {
	// Validate the input
	if err := s.validator.Struct(createDTO); err != nil {
		return nil, validation.NewValidationError(err.Error())
	}

	// Create the user first using the base service
	userResponse, err := s.BaseServiceImpl.Create(ctx, createDTO)
	if err != nil {
		return nil, err
	}

	// TODO: Temporarily disabled automatic business/staff creation due to schema mismatch
	// The staff table schema doesn't match the domain model (missing position, employment_type, join_date fields)
	// This needs to be fixed either by updating the domain model or the database schema
	/*
	// Create a default business for the user
	defaultBusinessName := userResponse.FullName + "'s Business"
	defaultBusiness := &domain.Business{
		UserID:           userResponse.ID,
		Name:             defaultBusinessName,
		Email:            userResponse.Email,
		Currency:         "EUR",
		TimeZone:         "Europe/Lisbon",
		IsActive:         true,
		SubscriptionTier: "free",
	}

	// Save the business
	err = s.businessRepo.Create(ctx, defaultBusiness)
	if err != nil {
		// If business creation fails, we should ideally rollback the user creation
		// For now, we'll log the error but still return the user
		// In a production system, this should be handled with a transaction
		return userResponse, NewServiceError("user created but failed to create default business", err)
	}

	// Create a staff entry for the user as owner of their business
	ownerStaff := &domain.Staff{
		BusinessID: defaultBusiness.ID,
		UserID:     userResponse.ID,
		Role:       domain.BusinessRoleOwner,
		IsActive:   true,
	}

	err = s.staffRepo.Create(ctx, ownerStaff)
	if err != nil {
		// If staff creation fails, log the error but still return the user
		return userResponse, NewServiceError("user and business created but failed to create owner staff record", err)
	}
	*/

	return userResponse, nil
}

// GetByEmail retrieves a user by email address
func (s *userServiceImpl) GetByEmail(ctx context.Context, email string) (*dto.UserResponseDTO, error) {
	if email == "" {
		return nil, validation.NewValidationError("email is required")
	}

	user, err := s.userRepo.FindByEmail(ctx, strings.ToLower(email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewNotFoundError("user", "email", email)
		}
		return nil, NewServiceError("failed to retrieve user by email", err)
	}

	return dto.ToUserResponseDTO(user), nil
}

// GetByClerkID retrieves a user by Clerk ID
func (s *userServiceImpl) GetByClerkID(ctx context.Context, clerkID string) (*dto.UserResponseDTO, error) {
	if clerkID == "" {
		return nil, validation.NewValidationError("clerk_id is required")
	}

	user, err := s.userRepo.FindByClerkID(ctx, clerkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewNotFoundError("user", "clerk_id", clerkID)
		}
		return nil, NewServiceError("failed to retrieve user by clerk_id", err)
	}

	return dto.ToUserResponseDTO(user), nil
}

// UpdateClerkID updates the Clerk ID for a user
func (s *userServiceImpl) UpdateClerkID(ctx context.Context, userID, clerkID string) error {
	if userID == "" {
		return validation.NewValidationError("user_id is required")
	}
	if clerkID == "" {
		return validation.NewValidationError("clerk_id is required")
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewNotFoundError("user", "id", userID)
		}
		return NewServiceError("failed to retrieve user", err)
	}

	// Check if Clerk ID is already in use
	existingUser, err := s.userRepo.FindByClerkID(ctx, clerkID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return NewServiceError("failed to check existing clerk_id", err)
	}
	if existingUser != nil && existingUser.ID != userID {
		return validation.NewValidationError("clerk_id already in use")
	}

	err = s.userRepo.UpdateClerkID(ctx, userID, clerkID)
	if err != nil {
		return NewServiceError("failed to update clerk_id", err)
	}

	return nil
}

// GetWithBusinesses retrieves a user with their businesses
func (s *userServiceImpl) GetWithBusinesses(ctx context.Context, userID string) (*dto.UserWithBusinessesDTO, error) {
	if userID == "" {
		return nil, validation.NewValidationError("user_id is required")
	}

	user, err := s.userRepo.GetWithBusinesses(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewNotFoundError("user", "id", userID)
		}
		return nil, NewServiceError("failed to retrieve user with businesses", err)
	}

	// Convert businesses to DTOs
	businessDTOs := make([]*dto.BusinessResponseDTO, len(user.Businesses))
	for i, business := range user.Businesses {
		businessDTOs[i] = dto.ToBusinessResponseDTO(&business)
	}

	return dto.ToUserWithBusinessesDTO(user, businessDTOs), nil
}

// SearchUsers searches for users by name or email
func (s *userServiceImpl) SearchUsers(ctx context.Context, query string, limit int) ([]*dto.UserResponseDTO, error) {
	if query == "" {
		return nil, validation.NewValidationError("search query is required")
	}

	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}

	users, err := s.userRepo.SearchUsers(ctx, query, limit)
	if err != nil {
		return nil, NewServiceError("failed to search users", err)
	}

	return dto.ToUserResponseDTOs(users), nil
}

// ActivateUser activates a user account
func (s *userServiceImpl) ActivateUser(ctx context.Context, userID string) error {
	if userID == "" {
		return validation.NewValidationError("user_id is required")
	}

	updateDTO := dto.UpdateUserDTO{
		IsActive: &[]bool{true}[0],
	}

	_, err := s.Update(ctx, userID, updateDTO)
	if err != nil {
		return NewServiceError("failed to activate user", err)
	}

	return nil
}

// DeactivateUser deactivates a user account
func (s *userServiceImpl) DeactivateUser(ctx context.Context, userID string) error {
	if userID == "" {
		return validation.NewValidationError("user_id is required")
	}

	updateDTO := dto.UpdateUserDTO{
		IsActive: &[]bool{false}[0],
	}

	_, err := s.Update(ctx, userID, updateDTO)
	if err != nil {
		return NewServiceError("failed to deactivate user", err)
	}

	return nil
}