package repository

import (
	"context"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
)

// BusinessRepository implements the domain.BusinessRepository interface using GORM
type BusinessRepository struct {
	*BaseRepository
}

// NewBusinessRepository creates a new instance of BusinessRepository
func NewBusinessRepository(db DBAdapter) domain.BusinessRepository {
	return &BusinessRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new business
func (r *BusinessRepository) Create(ctx context.Context, business *domain.Business) error {
	businessModel := mapBusinessDomainToModel(business)

	// Set created_by if available (from user context)
	// For now, we'll set it to nil since business creation might not have an existing user context
	if err := r.CreateWithAudit(ctx, businessModel, nil); err != nil {
		return fmt.Errorf("failed to create business: %w", err)
	}

	// Update the domain entity with generated fields
	business.BusinessID = businessModel.ID
	business.CreatedAt = businessModel.CreatedAt
	business.UpdatedAt = businessModel.UpdatedAt

	return nil
}

// GetByID retrieves a business by ID
func (r *BusinessRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Business, error) {
	var businessModel models.Business

	err := r.WithContext(ctx).First(&businessModel, "id = ?", id).Error
	if err != nil {
		return nil, r.HandleNotFound(err, "business", id)
	}

	return mapBusinessModelToDomain(&businessModel), nil
}

// GetByOwnerID retrieves all businesses owned by a specific user
func (r *BusinessRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Business, error) {
	var businessModels []models.Business

	err := r.WithContext(ctx).Where("user_id = ?", ownerID).Find(&businessModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find businesses for owner %s: %w", ownerID, err)
	}

	return mapBusinessModelsToDomainSlice(businessModels), nil
}

// Update updates a business
func (r *BusinessRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateBusinessInput) error {
	// First find the business to ensure it exists
	var businessModel models.Business
	err := r.WithContext(ctx).First(&businessModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "business", id)
	}

	// Build updates map from input
	updates := map[string]interface{}{}

	if input.BusinessName != nil {
		updates["name"] = *input.BusinessName
		updates["display_name"] = *input.BusinessName
	}

	if input.BusinessType != nil {
		// Note: BusinessType is not in models.Business
		// This would need to be added to the model or handled separately
	}

	if input.TaxID != nil {
		updates["tax_id"] = *input.TaxID
	}

	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}

	if input.Email != nil {
		updates["email"] = *input.Email
	}

	if input.AddressLine1 != nil {
		updates["address"] = *input.AddressLine1
	}

	if input.City != nil {
		updates["city"] = *input.City
	}

	if input.Region != nil {
		updates["state"] = *input.Region
	}

	if input.PostalCode != nil {
		updates["postal_code"] = *input.PostalCode
	}

	if input.Country != nil {
		updates["country"] = *input.Country
	}

	if input.TimeZone != nil {
		updates["timezone"] = *input.TimeZone
	}

	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if input.SubscriptionPlan != nil {
		updates["subscription_tier"] = models.SubscriptionTier(*input.SubscriptionPlan)
	}

	if input.BusinessHours != nil {
		// Handle business hours update - need to update the settings JSONB field
		var workingHours models.WorkingHours
		err := workingHours.Scan(*input.BusinessHours)
		if err == nil {
			// Update the settings with new working hours
			settings := businessModel.Settings
			settings.WorkingHours = workingHours
			updates["settings"] = settings
		}
	}

	// For now, we'll assume updatedBy comes from context or is passed separately
	// In a real implementation, this would be extracted from the authentication context
	updatedBy := uuid.New() // Placeholder - should come from authenticated user context

	err = r.UpdateWithAudit(ctx, &businessModel, updates, updatedBy)
	if err != nil {
		return fmt.Errorf("failed to update business: %w", err)
	}

	return nil
}

// Delete soft deletes a business
func (r *BusinessRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// First find the business to ensure it exists
	var businessModel models.Business
	err := r.WithContext(ctx).First(&businessModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "business", id)
	}

	// For now, we'll assume deletedBy comes from context
	// In a real implementation, this would be extracted from the authentication context
	deletedBy := uuid.New() // Placeholder - should come from authenticated user context

	err = r.SoftDeleteWithAudit(ctx, &businessModel, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete business: %w", err)
	}

	return nil
}

// List retrieves a paginated list of businesses
func (r *BusinessRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Business, error) {
	var businessModels []models.Business

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&businessModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list businesses: %w", err)
	}

	return mapBusinessModelsToDomainSlice(businessModels), nil
}

// Count counts all businesses
func (r *BusinessRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.WithContext(ctx).Model(&models.Business{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count businesses: %w", err)
	}

	return count, nil
}

// Helper function to map slice of models to slice of domain entities
func mapBusinessModelsToDomainSlice(businessModels []models.Business) []*domain.Business {
	result := make([]*domain.Business, len(businessModels))
	for i, model := range businessModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapBusinessModelToDomain(&modelCopy)
	}
	return result
}
