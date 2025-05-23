package repository

import (
	"context"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
)

// UserRepository implements the domain.UserRepository interface using GORM
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db DBAdapter) domain.UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	userModel := mapUserDomainToModel(user)

	// Set created_by if available (from user context)
	// For now, we'll set it to nil since user creation might not have an existing user context
	if err := r.CreateWithAudit(ctx, userModel, nil); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Update the domain entity with generated fields
	user.UserID = userModel.ID
	user.CreatedAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var userModel models.User

	err := r.WithContext(ctx).First(&userModel, "id = ?", id).Error
	if err != nil {
		return nil, r.HandleNotFound(err, "user", id)
	}

	return mapUserModelToDomain(&userModel), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userModel models.User

	err := r.WithContext(ctx).First(&userModel, "email = ?", email).Error
	if err != nil {
		return nil, fmt.Errorf("user with email %s not found: %w", email, err)
	}

	return mapUserModelToDomain(&userModel), nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateUserInput) error {
	// First find the user to ensure it exists
	var userModel models.User
	err := r.WithContext(ctx).First(&userModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "user", id)
	}

	// Build updates map from input
	updates := map[string]interface{}{}

	if input.Email != nil {
		updates["email"] = *input.Email
	}

	if input.Password != nil {
		// Note: Password hashing should be handled at the service layer
		// Here we assume it's already hashed
		updates["password_hash"] = *input.Password
	}

	if input.FirstName != nil {
		updates["first_name"] = *input.FirstName
	}

	if input.LastName != nil {
		updates["last_name"] = *input.LastName
	}

	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}

	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if input.EmailVerified != nil {
		// Note: EmailVerified is not in models.User
		// This would need to be handled in a separate table or added to the model
	}

	if input.LanguagePreference != nil {
		// Note: LanguagePreference is not in models.User
		// This would need to be handled in a separate table or added to the model
	}

	// For now, we'll assume updatedBy comes from context or is passed separately
	// In a real implementation, this would be extracted from the authentication context
	updatedBy := uuid.New() // Placeholder - should come from authenticated user context

	err = r.UpdateWithAudit(ctx, &userModel, updates, updatedBy)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// First find the user to ensure it exists
	var userModel models.User
	err := r.WithContext(ctx).First(&userModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "user", id)
	}

	// For now, we'll assume deletedBy comes from context
	// In a real implementation, this would be extracted from the authentication context
	deletedBy := uuid.New() // Placeholder - should come from authenticated user context

	err = r.SoftDeleteWithAudit(ctx, &userModel, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List retrieves a paginated list of users
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]*domain.User, error) {
	var userModels []models.User

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&userModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return mapUserModelsToDomainSlice(userModels), nil
}

// Count counts all users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// Helper function to map slice of models to slice of domain entities
func mapUserModelsToDomainSlice(userModels []models.User) []*domain.User {
	result := make([]*domain.User, len(userModels))
	for i, model := range userModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapUserModelToDomain(&modelCopy)
	}
	return result
}
