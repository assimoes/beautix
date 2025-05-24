package repository

import (
	"context"
	"strings"

	"github.com/assimoes/beautix/internal/domain"
	"gorm.io/gorm"
)

// userRepositoryImpl implements the UserRepository interface
type userRepositoryImpl struct {
	*BaseRepositoryImpl[domain.User]
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepositoryImpl{
		BaseRepositoryImpl: &BaseRepositoryImpl[domain.User]{db: db},
	}
}

// FindByEmail finds a user by email address
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("email = ?", strings.ToLower(email)).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByClerkID finds a user by Clerk ID
func (r *userRepositoryImpl) FindByClerkID(ctx context.Context, clerkID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("clerk_id = ?", clerkID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateClerkID updates the Clerk ID for a user
func (r *userRepositoryImpl) UpdateClerkID(ctx context.Context, userID, clerkID string) error {
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Update("clerk_id", clerkID).Error
}


// GetWithBusinesses retrieves a user with their businesses preloaded
func (r *userRepositoryImpl) GetWithBusinesses(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("Businesses").
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SearchUsers searches for users by name or email
func (r *userRepositoryImpl) SearchUsers(ctx context.Context, query string, limit int) ([]*domain.User, error) {
	var users []*domain.User
	
	searchTerm := "%" + strings.ToLower(query) + "%"
	
	err := r.db.WithContext(ctx).
		Where("(LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email) LIKE ?) AND deleted_at IS NULL",
			searchTerm, searchTerm, searchTerm).
		Limit(limit).
		Find(&users).Error
		
	return users, err
}

// WithTx returns a new repository instance with the given transaction
func (r *userRepositoryImpl) WithTx(tx *gorm.DB) domain.BaseRepository[domain.User] {
	return &BaseRepositoryImpl[domain.User]{db: tx}
}