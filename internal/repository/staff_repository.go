package repository

import (
	"context"

	"github.com/assimoes/beautix/internal/domain"
	"gorm.io/gorm"
)

// staffRepositoryImpl implements the StaffRepository interface
type staffRepositoryImpl struct {
	*BaseRepositoryImpl[domain.Staff]
}

// NewStaffRepository creates a new staff repository
func NewStaffRepository(db *gorm.DB) domain.StaffRepository {
	return &staffRepositoryImpl{
		BaseRepositoryImpl: &BaseRepositoryImpl[domain.Staff]{db: db},
	}
}

// FindByBusinessID finds all staff members for a business
func (r *staffRepositoryImpl) FindByBusinessID(ctx context.Context, businessID string) ([]*domain.Staff, error) {
	var staff []*domain.Staff
	err := r.db.WithContext(ctx).
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Find(&staff).Error
	return staff, err
}

// FindByUserID finds all staff positions for a user
func (r *staffRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*domain.Staff, error) {
	var staff []*domain.Staff
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&staff).Error
	return staff, err
}

// FindByBusinessAndUser finds a specific staff record
func (r *staffRepositoryImpl) FindByBusinessAndUser(ctx context.Context, businessID, userID string) (*domain.Staff, error) {
	var staff domain.Staff
	err := r.db.WithContext(ctx).
		Where("business_id = ? AND user_id = ? AND deleted_at IS NULL", businessID, userID).
		First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}

// FindActiveByBusinessID finds all active staff members for a business
func (r *staffRepositoryImpl) FindActiveByBusinessID(ctx context.Context, businessID string) ([]*domain.Staff, error) {
	var staff []*domain.Staff
	err := r.db.WithContext(ctx).
		Where("business_id = ? AND is_active = true AND deleted_at IS NULL", businessID).
		Find(&staff).Error
	return staff, err
}

// FindByRole finds staff members by role in a business
func (r *staffRepositoryImpl) FindByRole(ctx context.Context, businessID string, role domain.BusinessRole) ([]*domain.Staff, error) {
	var staff []*domain.Staff
	err := r.db.WithContext(ctx).
		Where("business_id = ? AND role = ? AND deleted_at IS NULL", businessID, role).
		Find(&staff).Error
	return staff, err
}

// UpdateRole updates the role of a staff member
func (r *staffRepositoryImpl) UpdateRole(ctx context.Context, businessID, userID string, role domain.BusinessRole) error {
	return r.db.WithContext(ctx).
		Model(&domain.Staff{}).
		Where("business_id = ? AND user_id = ?", businessID, userID).
		Update("role", role).Error
}

// DeactivateStaff deactivates a staff member
func (r *staffRepositoryImpl) DeactivateStaff(ctx context.Context, businessID, userID string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Staff{}).
		Where("business_id = ? AND user_id = ?", businessID, userID).
		Update("is_active", false).Error
}