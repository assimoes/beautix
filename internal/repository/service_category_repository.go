package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
)

// ServiceCategoryRepository implements the domain.ServiceCategoryRepository interface using GORM
type ServiceCategoryRepository struct {
	*BaseRepository
}

// NewServiceCategoryRepository creates a new instance of ServiceCategoryRepository
func NewServiceCategoryRepository(db DBAdapter) domain.ServiceCategoryRepository {
	return &ServiceCategoryRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new service category
func (r *ServiceCategoryRepository) Create(ctx context.Context, category *domain.ServiceCategory) error {
	categoryModel := mapServiceCategoryDomainToModel(category)

	// Set creation time (ServiceCategory table doesn't have audit fields)
	if categoryModel.CreatedAt.IsZero() {
		categoryModel.CreatedAt = time.Now()
	}

	if err := r.WithContext(ctx).Create(&categoryModel).Error; err != nil {
		return fmt.Errorf("failed to create service category: %w", err)
	}

	// Update the domain entity with any generated fields
	category.ID = categoryModel.ID
	category.CreatedAt = categoryModel.CreatedAt

	return nil
}

// GetByID retrieves a service category by ID
func (r *ServiceCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ServiceCategory, error) {
	var categoryModel models.ServiceCategory

	err := r.WithContext(ctx).First(&categoryModel, "id = ?", id).Error

	if err != nil {
		return nil, r.HandleNotFound(err, "service category", id)
	}

	return mapServiceCategoryModelToDomain(&categoryModel), nil
}

// Update updates a service category
func (r *ServiceCategoryRepository) Update(ctx context.Context, id uuid.UUID, name, description string, updatedBy uuid.UUID) error {
	// First find the category to ensure it exists
	var categoryModel models.ServiceCategory
	err := r.WithContext(ctx).First(&categoryModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "service category", id)
	}

	// Apply updates (ServiceCategory table doesn't have audit fields)
	updates := map[string]interface{}{
		"name":        name,
		"description": description,
		"updated_at":  time.Now(),
	}

	// Perform the update without audit (ServiceCategory doesn't have audit fields)
	err = r.WithContext(ctx).Model(&categoryModel).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update service category: %w", err)
	}

	return nil
}

// Delete hard deletes a service category (no soft delete fields in table)
func (r *ServiceCategoryRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// Perform hard delete (ServiceCategory table doesn't have soft delete fields)
	err := r.WithContext(ctx).Delete(&models.ServiceCategory{}, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("failed to delete service category: %w", err)
	}

	return nil
}

// List retrieves a paginated list of service categories
func (r *ServiceCategoryRepository) List(ctx context.Context, page, pageSize int) ([]*domain.ServiceCategory, error) {
	var categoryModels []models.ServiceCategory

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("name ASC").
		Find(&categoryModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list service categories: %w", err)
	}

	return mapServiceCategoryModelsToDomainSlice(categoryModels), nil
}

// Count counts all service categories
func (r *ServiceCategoryRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.WithContext(ctx).Model(&models.ServiceCategory{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count service categories: %w", err)
	}

	return count, nil
}

// Helper functions to map between domain entities and models

// mapServiceCategoryDomainToModel converts a domain ServiceCategory entity to a model ServiceCategory
func mapServiceCategoryDomainToModel(sc *domain.ServiceCategory) *models.ServiceCategory {
	if sc == nil {
		return nil
	}

	categoryModel := &models.ServiceCategory{
		ID:          sc.ID,
		BusinessID:  sc.BusinessID,
		Name:        sc.Name,
		Description: sc.Description,
		CreatedAt:   sc.CreatedAt,
	}

	// Handle optional/pointer fields
	if sc.UpdatedAt != nil {
		categoryModel.UpdatedAt = *sc.UpdatedAt
	}

	return categoryModel
}

// mapServiceCategoryModelsToDomainSlice converts a slice of model ServiceCategory to a slice of domain ServiceCategory entities
func mapServiceCategoryModelsToDomainSlice(categoryModels []models.ServiceCategory) []*domain.ServiceCategory {
	result := make([]*domain.ServiceCategory, len(categoryModels))
	for i, model := range categoryModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapServiceCategoryModelToDomain(&modelCopy)
	}
	return result
}
