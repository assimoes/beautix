package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServiceCompletionRepository implements the domain.ServiceCompletionRepository interface using GORM
type ServiceCompletionRepository struct {
	*BaseRepository
}

// NewServiceCompletionRepository creates a new instance of ServiceCompletionRepository
func NewServiceCompletionRepository(db DBAdapter) domain.ServiceCompletionRepository {
	return &ServiceCompletionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new service completion
func (r *ServiceCompletionRepository) Create(ctx context.Context, completion *domain.ServiceCompletion) error {
	completionModel := mapDomainToServiceCompletion(completion)

	err := r.WithContext(ctx).Create(&completionModel).Error
	if err != nil {
		return fmt.Errorf("failed to create service completion: %w", err)
	}

	// Update the domain entity with any generated fields
	completion.ID = completionModel.ID
	completion.CreatedAt = completionModel.CreatedAt

	return nil
}

// GetByID retrieves a service completion by ID
func (r *ServiceCompletionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ServiceCompletion, error) {
	var completionModel models.ServiceCompletion

	err := r.WithContext(ctx).
		Preload("Appointment").
		Where("deleted_at IS NULL").
		First(&completionModel, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("service completion not found")
		}
		return nil, fmt.Errorf("failed to get service completion: %w", err)
	}

	return mapServiceCompletionToDomain(&completionModel), nil
}

// GetByAppointmentID retrieves a service completion by appointment ID
func (r *ServiceCompletionRepository) GetByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*domain.ServiceCompletion, error) {
	var completionModel models.ServiceCompletion

	err := r.WithContext(ctx).
		Preload("Appointment").
		Where("deleted_at IS NULL").
		First(&completionModel, "appointment_id = ?", appointmentID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("service completion not found for appointment")
		}
		return nil, fmt.Errorf("failed to get service completion by appointment: %w", err)
	}

	return mapServiceCompletionToDomain(&completionModel), nil
}

// Update updates an existing service completion
func (r *ServiceCompletionRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateServiceCompletionInput, updatedBy uuid.UUID) error {
	updateData := map[string]interface{}{
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}

	if input.PriceCharged != nil {
		updateData["price_charged"] = *input.PriceCharged
	}
	if input.PaymentMethod != nil {
		updateData["payment_method"] = *input.PaymentMethod
	}
	if input.ProviderConfirmed != nil {
		updateData["provider_confirmed"] = *input.ProviderConfirmed
	}
	if input.ClientConfirmed != nil {
		updateData["client_confirmed"] = *input.ClientConfirmed
	}
	if input.CompletionDate != nil {
		updateData["completion_date"] = *input.CompletionDate
	}

	result := r.WithContext(ctx).
		Model(&models.ServiceCompletion{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updateData)

	if result.Error != nil {
		return fmt.Errorf("failed to update service completion: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("service completion not found")
	}

	return nil
}

// Delete soft deletes a service completion
func (r *ServiceCompletionRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	updateData := map[string]interface{}{
		"deleted_by": deletedBy,
		"deleted_at": time.Now(),
	}

	result := r.WithContext(ctx).
		Model(&models.ServiceCompletion{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updateData)

	if result.Error != nil {
		return fmt.Errorf("failed to delete service completion: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("service completion not found")
	}

	return nil
}

// ListByProvider lists service completions for a provider within a date range
func (r *ServiceCompletionRepository) ListByProvider(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time, page, pageSize int) ([]*domain.ServiceCompletion, error) {
	var completionModels []models.ServiceCompletion

	offset := (page - 1) * pageSize

	err := r.WithContext(ctx).
		Preload("Appointment").
		Joins("JOIN appointments a ON service_completions.appointment_id = a.id").
		Where("a.business_id = ?", providerID).
		Where("service_completions.completion_date >= ? AND service_completions.completion_date <= ?", startDate, endDate).
		Where("service_completions.deleted_at IS NULL").
		Order("service_completions.completion_date DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&completionModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list service completions by provider: %w", err)
	}

	completions := make([]*domain.ServiceCompletion, 0, len(completionModels))
	for _, model := range completionModels {
		completions = append(completions, mapServiceCompletionToDomain(&model))
	}

	return completions, nil
}

// Count returns the total number of service completions
func (r *ServiceCompletionRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.WithContext(ctx).
		Model(&models.ServiceCompletion{}).
		Where("deleted_at IS NULL").
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count service completions: %w", err)
	}

	return count, nil
}

// CountByProvider returns the number of service completions for a provider
func (r *ServiceCompletionRepository) CountByProvider(ctx context.Context, providerID uuid.UUID) (int64, error) {
	var count int64

	err := r.WithContext(ctx).
		Model(&models.ServiceCompletion{}).
		Joins("JOIN appointments a ON service_completions.appointment_id = a.id").
		Where("a.business_id = ?", providerID).
		Where("service_completions.deleted_at IS NULL").
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count service completions by provider: %w", err)
	}

	return count, nil
}

// GetProviderRevenue calculates total revenue for a provider within a date range
func (r *ServiceCompletionRepository) GetProviderRevenue(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time) (float64, error) {
	var revenue float64

	err := r.WithContext(ctx).
		Model(&models.ServiceCompletion{}).
		Select("COALESCE(SUM(price_charged), 0)").
		Joins("JOIN appointments a ON service_completions.appointment_id = a.id").
		Where("a.business_id = ?", providerID).
		Where("service_completions.completion_date >= ? AND service_completions.completion_date <= ?", startDate, endDate).
		Where("service_completions.deleted_at IS NULL").
		Scan(&revenue).Error

	if err != nil {
		return 0, fmt.Errorf("failed to calculate provider revenue: %w", err)
	}

	return revenue, nil
}

// Helper functions to map between domain entities and database models

// mapDomainToServiceCompletion converts a domain ServiceCompletion entity to a database ServiceCompletion model
func mapDomainToServiceCompletion(sc *domain.ServiceCompletion) *models.ServiceCompletion {
	if sc == nil {
		return nil
	}

	completionModel := &models.ServiceCompletion{
		BaseModel: models.BaseModel{
			ID:        sc.ID,
			CreatedAt: sc.CreatedAt,
			CreatedBy: sc.CreatedBy,
		},
		AppointmentID:     sc.AppointmentID,
		PriceCharged:      sc.PriceCharged,
		PaymentMethod:     sc.PaymentMethod,
		ProviderConfirmed: sc.ProviderConfirmed,
		ClientConfirmed:   sc.ClientConfirmed,
		CompletionDate:    sc.CompletionDate,
	}

	// Handle optional/pointer fields
	if sc.UpdatedAt != nil {
		completionModel.UpdatedAt = *sc.UpdatedAt
		completionModel.UpdatedBy = sc.UpdatedBy
	}

	if sc.DeletedAt != nil {
		completionModel.DeletedAt = gorm.DeletedAt{Time: *sc.DeletedAt, Valid: true}
		completionModel.DeletedBy = sc.DeletedBy
	}

	return completionModel
}

// mapServiceCompletionToDomain converts a database ServiceCompletion model to a domain ServiceCompletion entity
func mapServiceCompletionToDomain(scModel *models.ServiceCompletion) *domain.ServiceCompletion {
	if scModel == nil {
		return nil
	}

	completion := &domain.ServiceCompletion{
		ID:                scModel.ID,
		AppointmentID:     scModel.AppointmentID,
		PriceCharged:      scModel.PriceCharged,
		PaymentMethod:     scModel.PaymentMethod,
		ProviderConfirmed: scModel.ProviderConfirmed,
		ClientConfirmed:   scModel.ClientConfirmed,
		CompletionDate:    scModel.CompletionDate,
		CreatedAt:         scModel.CreatedAt,
		CreatedBy:         scModel.CreatedBy,
	}

	// Handle optional/pointer fields
	if !scModel.UpdatedAt.IsZero() {
		completion.UpdatedAt = &scModel.UpdatedAt
	}

	if scModel.UpdatedBy != nil {
		completion.UpdatedBy = scModel.UpdatedBy
	}

	if scModel.DeletedAt.Valid {
		completion.DeletedAt = &scModel.DeletedAt.Time
	}

	if scModel.DeletedBy != nil {
		completion.DeletedBy = scModel.DeletedBy
	}

	// Map expanded relationships if loaded
	if scModel.Appointment != nil {
		completion.Appointment = mapAppointmentToDomain(scModel.Appointment)
	}

	return completion
}
