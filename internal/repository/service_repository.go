package repository

import (
	"context"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
)

// ServiceRepository implements the domain.ServiceRepository interface using GORM
type ServiceRepository struct {
	*BaseRepository
}

// NewServiceRepository creates a new instance of ServiceRepository
func NewServiceRepository(db DBAdapter) domain.ServiceRepository {
	return &ServiceRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new service
func (r *ServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	serviceModel := mapServiceDomainToModel(service)
	
	if err := r.CreateWithAudit(ctx, &serviceModel, service.CreatedBy); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	
	// Update the domain entity with any generated fields
	service.ID = serviceModel.ID
	service.CreatedAt = serviceModel.CreatedAt
	
	return nil
}

// GetByID retrieves a service by ID
func (r *ServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	var serviceModel models.Service
	
	err := r.WithContext(ctx).
		Preload("Business").
		Preload("Category").
		First(&serviceModel, "id = ?", id).Error
	
	if err != nil {
		return nil, r.HandleNotFound(err, "service", id)
	}
	
	return mapServiceModelToDomain(&serviceModel), nil
}

// Update updates a service
func (r *ServiceRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateServiceInput, updatedBy uuid.UUID) error {
	// First find the service to ensure it exists
	var serviceModel models.Service
	err := r.WithContext(ctx).First(&serviceModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "service", id)
	}
	
	// Apply updates from the input
	updates := map[string]interface{}{}
	
	if input.CategoryID != nil {
		updates["category_id"] = *input.CategoryID
	}
	
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	
	if input.Duration != nil {
		updates["duration"] = *input.Duration
	}
	
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	
	// Perform the update with audit
	err = r.UpdateWithAudit(ctx, &serviceModel, updates, updatedBy)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}
	
	return nil
}

// Delete soft deletes a service
func (r *ServiceRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// First find the service to ensure it exists
	var serviceModel models.Service
	err := r.WithContext(ctx).First(&serviceModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "service", id)
	}
	
	// Perform soft delete with audit
	err = r.SoftDeleteWithAudit(ctx, &serviceModel, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	
	return nil
}

// ListByProvider retrieves services by provider ID
func (r *ServiceRepository) ListByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) ([]*domain.Service, error) {
	var serviceModels []models.Service
	
	offset := r.CalculateOffset(page, pageSize)
	
	err := r.WithContext(ctx).
		Preload("Business").
		Preload("Category").
		Where("business_id = ? AND is_active = ?", providerID, true).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&serviceModels).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to list services by provider: %w", err)
	}
	
	return mapServiceModelsToDomainSlice(serviceModels), nil
}

// ListByCategory retrieves services by category ID
func (r *ServiceRepository) ListByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int) ([]*domain.Service, error) {
	var serviceModels []models.Service
	
	offset := r.CalculateOffset(page, pageSize)
	
	err := r.WithContext(ctx).
		Preload("Business").
		Preload("Category").
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&serviceModels).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to list services by category: %w", err)
	}
	
	return mapServiceModelsToDomainSlice(serviceModels), nil
}

// Count counts all services
func (r *ServiceRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	
	err := r.WithContext(ctx).Model(&models.Service{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count services: %w", err)
	}
	
	return count, nil
}

// CountByProvider counts services by provider ID
func (r *ServiceRepository) CountByProvider(ctx context.Context, providerID uuid.UUID) (int64, error) {
	var count int64
	
	err := r.WithContext(ctx).Model(&models.Service{}).Where("business_id = ?", providerID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count services by provider: %w", err)
	}
	
	return count, nil
}

// Helper functions to map between domain entities and models

// mapServiceDomainToModel converts a domain Service entity to a model Service
func mapServiceDomainToModel(s *domain.Service) *models.Service {
	if s == nil {
		return nil
	}
	
	serviceModel := &models.Service{
		BaseModel: models.BaseModel{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
		},
		BusinessID:  s.ProviderID,
		CategoryID:  s.CategoryID,
		Name:        s.Name,
		Description: s.Description,
		Duration:    s.Duration,
		Price:       s.Price,
		IsActive:    true,
	}
	
	// Handle optional/pointer fields
	if s.CreatedBy != nil {
		serviceModel.CreatedBy = s.CreatedBy
	}
	
	if s.UpdatedAt != nil {
		serviceModel.UpdatedAt = *s.UpdatedAt
	}
	
	if s.UpdatedBy != nil {
		serviceModel.UpdatedBy = s.UpdatedBy
	}
	
	if s.DeletedAt != nil {
		serviceModel.DeletedAt.Time = *s.DeletedAt
		serviceModel.DeletedAt.Valid = true
	}
	
	if s.DeletedBy != nil {
		serviceModel.DeletedBy = s.DeletedBy
	}
	
	return serviceModel
}

// mapServiceModelToDomain converts a model Service to a domain Service entity
func mapServiceModelToDomain(s *models.Service) *domain.Service {
	if s == nil {
		return nil
	}
	
	service := &domain.Service{
		ID:          s.ID,
		ProviderID:  s.BusinessID,
		CategoryID:  s.CategoryID,
		Name:        s.Name,
		Description: s.Description,
		Duration:    s.Duration,
		Price:       s.Price,
		CreatedAt:   s.CreatedAt,
	}
	
	// Handle optional/pointer fields
	if s.CreatedBy != nil {
		service.CreatedBy = s.CreatedBy
	}
	
	if !s.UpdatedAt.IsZero() {
		service.UpdatedAt = &s.UpdatedAt
	}
	
	if s.UpdatedBy != nil {
		service.UpdatedBy = s.UpdatedBy
	}
	
	if s.DeletedAt.Valid {
		deletedAt := s.DeletedAt.Time
		service.DeletedAt = &deletedAt
	}
	
	if s.DeletedBy != nil {
		service.DeletedBy = s.DeletedBy
	}
	
	// Map related entities if loaded
	if s.Business.ID != uuid.Nil {
		service.Provider = mapProviderModelToDomain(&s.Business)
	}
	
	if s.Category.ID != uuid.Nil {
		service.Category = mapServiceCategoryModelToDomain(&s.Category)
	}
	
	return service
}

// mapServiceModelsToDomainSlice converts a slice of model Service to a slice of domain Service entities
func mapServiceModelsToDomainSlice(serviceModels []models.Service) []*domain.Service {
	result := make([]*domain.Service, len(serviceModels))
	for i, model := range serviceModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapServiceModelToDomain(&modelCopy)
	}
	return result
}