package repository

import (
	"context"

	"github.com/assimoes/beautix/internal/domain"
	"gorm.io/gorm"
)

// BaseRepositoryImpl provides a base implementation for repositories
type BaseRepositoryImpl[T any] struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository[T any](db *gorm.DB) domain.BaseRepository[T] {
	return &BaseRepositoryImpl[T]{
		db: db,
	}
}

// Create creates a new entity
func (r *BaseRepositoryImpl[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID retrieves an entity by ID
func (r *BaseRepositoryImpl[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update updates an existing entity
func (r *BaseRepositoryImpl[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete soft deletes an entity by ID
func (r *BaseRepositoryImpl[T]) Delete(ctx context.Context, id string) error {
	var entity T
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
}

// List retrieves entities with pagination
func (r *BaseRepositoryImpl[T]) List(ctx context.Context, page, pageSize int) ([]*T, int64, error) {
	var entities []*T
	var total int64
	
	// Count total records
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Calculate offset
	offset := (page - 1) * pageSize
	
	// Retrieve paginated results
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Find(&entities).Error
	
	return entities, total, err
}

// FindBy finds entities by criteria
func (r *BaseRepositoryImpl[T]) FindBy(ctx context.Context, criteria map[string]any) ([]*T, error) {
	var entities []*T
	query := r.db.WithContext(ctx)
	
	for key, value := range criteria {
		query = query.Where(key+" = ?", value)
	}
	
	err := query.Find(&entities).Error
	return entities, err
}

// ExistsByID checks if an entity exists by ID
func (r *BaseRepositoryImpl[T]) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// WithTx returns a new repository instance with the given transaction
func (r *BaseRepositoryImpl[T]) WithTx(tx *gorm.DB) domain.BaseRepository[T] {
	return &BaseRepositoryImpl[T]{db: tx}
}

// GetDB returns the database instance
func (r *BaseRepositoryImpl[T]) GetDB() *gorm.DB {
	return r.db
}