package repository

import (
	"context"
	"strings"

	"github.com/assimoes/beautix/internal/domain"
	"gorm.io/gorm"
)

// businessRepositoryImpl implements the BusinessRepository interface
type businessRepositoryImpl struct {
	*BaseRepositoryImpl[domain.Business]
}

// NewBusinessRepository creates a new business repository
func NewBusinessRepository(db *gorm.DB) domain.BusinessRepository {
	return &businessRepositoryImpl{
		BaseRepositoryImpl: &BaseRepositoryImpl[domain.Business]{db: db},
	}
}

// FindByUserID finds businesses owned by a user
func (r *businessRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*domain.Business, error) {
	var businesses []*domain.Business
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&businesses).Error
	return businesses, err
}

// FindByName finds businesses by name (case-insensitive)
func (r *businessRepositoryImpl) FindByName(ctx context.Context, name string) ([]*domain.Business, error) {
	var businesses []*domain.Business
	err := r.db.WithContext(ctx).
		Where("LOWER(name) LIKE ? AND deleted_at IS NULL", "%"+strings.ToLower(name)+"%").
		Find(&businesses).Error
	return businesses, err
}

// ExistsByName checks if a business with the given name exists
func (r *businessRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Business{}).
		Where("LOWER(name) = ? AND deleted_at IS NULL", strings.ToLower(name)).
		Count(&count).Error
	return count > 0, err
}

// FindActiveBusinesses finds active businesses with pagination
func (r *businessRepositoryImpl) FindActiveBusinesses(ctx context.Context, page, pageSize int) ([]*domain.Business, int64, error) {
	var businesses []*domain.Business
	var total int64
	
	// Count total active businesses
	if err := r.db.WithContext(ctx).
		Model(&domain.Business{}).
		Where("is_active = true AND deleted_at IS NULL").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("is_active = true AND deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Find(&businesses).Error
		
	return businesses, total, err
}

// SearchByLocation finds businesses by city and country
func (r *businessRepositoryImpl) SearchByLocation(ctx context.Context, city, country string) ([]*domain.Business, error) {
	var businesses []*domain.Business
	
	query := r.db.WithContext(ctx).
		Joins("JOIN business_locations bl ON businesses.id = bl.business_id").
		Where("businesses.is_active = true AND businesses.deleted_at IS NULL AND bl.deleted_at IS NULL")
	
	if city != "" {
		query = query.Where("LOWER(bl.city) = ?", strings.ToLower(city))
	}
	
	if country != "" {
		query = query.Where("LOWER(bl.country) = ?", strings.ToLower(country))
	}
	
	err := query.Distinct().Find(&businesses).Error
	return businesses, err
}

// SearchByService finds businesses that offer a specific service
func (r *businessRepositoryImpl) SearchByService(ctx context.Context, serviceName string) ([]*domain.Business, error) {
	var businesses []*domain.Business
	
	err := r.db.WithContext(ctx).
		Joins("JOIN services s ON businesses.id = s.business_id").
		Where("businesses.is_active = true AND businesses.deleted_at IS NULL AND s.deleted_at IS NULL").
		Where("LOWER(s.name) LIKE ?", "%"+strings.ToLower(serviceName)+"%").
		Distinct().
		Find(&businesses).Error
		
	return businesses, err
}

// GetBusinessWithDetails retrieves a business with all related data
func (r *businessRepositoryImpl) GetBusinessWithDetails(ctx context.Context, businessID string) (*domain.Business, error) {
	var business domain.Business
	err := r.db.WithContext(ctx).
		Preload("Locations").
		Preload("Settings_").
		Preload("User").
		Where("id = ?", businessID).
		First(&business).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

// GetWithLocations retrieves a business with its locations
func (r *businessRepositoryImpl) GetWithLocations(ctx context.Context, businessID string) (*domain.Business, error) {
	var business domain.Business
	err := r.db.WithContext(ctx).
		Preload("Locations", "deleted_at IS NULL").
		Where("id = ?", businessID).
		First(&business).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

// WithTx returns a new repository instance with the given transaction
func (r *businessRepositoryImpl) WithTx(tx *gorm.DB) domain.BaseRepository[domain.Business] {
	return &BaseRepositoryImpl[domain.Business]{db: tx}
}

// businessLocationRepositoryImpl implements the BusinessLocationRepository interface
type businessLocationRepositoryImpl struct {
	*BaseRepositoryImpl[domain.BusinessLocation]
}

// NewBusinessLocationRepository creates a new business location repository
func NewBusinessLocationRepository(db *gorm.DB) domain.BusinessLocationRepository {
	return &businessLocationRepositoryImpl{
		BaseRepositoryImpl: &BaseRepositoryImpl[domain.BusinessLocation]{db: db},
	}
}

// FindByBusinessID finds all locations for a business
func (r *businessLocationRepositoryImpl) FindByBusinessID(ctx context.Context, businessID string) ([]*domain.BusinessLocation, error) {
	var locations []*domain.BusinessLocation
	err := r.db.WithContext(ctx).
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Order("is_main DESC, name ASC").
		Find(&locations).Error
	return locations, err
}

// GetMainLocation retrieves the main location for a business
func (r *businessLocationRepositoryImpl) GetMainLocation(ctx context.Context, businessID string) (*domain.BusinessLocation, error) {
	var location domain.BusinessLocation
	err := r.db.WithContext(ctx).
		Where("business_id = ? AND is_main = true AND deleted_at IS NULL", businessID).
		First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

// SetMainLocation sets a location as the main location for a business
func (r *businessLocationRepositoryImpl) SetMainLocation(ctx context.Context, businessID, locationID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First, unset all main locations for this business
		if err := tx.Model(&domain.BusinessLocation{}).
			Where("business_id = ?", businessID).
			Update("is_main", false).Error; err != nil {
			return err
		}
		
		// Then set the specified location as main
		return tx.Model(&domain.BusinessLocation{}).
			Where("id = ? AND business_id = ?", locationID, businessID).
			Update("is_main", true).Error
	})
}

// WithTx returns a new repository instance with the given transaction
func (r *businessLocationRepositoryImpl) WithTx(tx *gorm.DB) domain.BaseRepository[domain.BusinessLocation] {
	return &BaseRepositoryImpl[domain.BusinessLocation]{db: tx}
}

// businessSettingsRepositoryImpl implements the BusinessSettingsRepository interface
type businessSettingsRepositoryImpl struct {
	*BaseRepositoryImpl[domain.BusinessSettings]
}

// NewBusinessSettingsRepository creates a new business settings repository
func NewBusinessSettingsRepository(db *gorm.DB) domain.BusinessSettingsRepository {
	return &businessSettingsRepositoryImpl{
		BaseRepositoryImpl: &BaseRepositoryImpl[domain.BusinessSettings]{db: db},
	}
}

// GetByBusinessID retrieves settings for a business
func (r *businessSettingsRepositoryImpl) GetByBusinessID(ctx context.Context, businessID string) (*domain.BusinessSettings, error) {
	var settings domain.BusinessSettings
	err := r.db.WithContext(ctx).
		Where("business_id = ?", businessID).
		First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// UpdateByBusinessID updates settings for a business
func (r *businessSettingsRepositoryImpl) UpdateByBusinessID(ctx context.Context, businessID string, settings *domain.BusinessSettings) error {
	return r.db.WithContext(ctx).
		Model(&domain.BusinessSettings{}).
		Where("business_id = ?", businessID).
		Updates(settings).Error
}

// WithTx returns a new repository instance with the given transaction
func (r *businessSettingsRepositoryImpl) WithTx(tx *gorm.DB) domain.BaseRepository[domain.BusinessSettings] {
	return &BaseRepositoryImpl[domain.BusinessSettings]{db: tx}
}