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

// AvailabilityExceptionRepository implements the domain.AvailabilityExceptionRepository interface using GORM
type AvailabilityExceptionRepository struct {
	*BaseRepository
}

// NewAvailabilityExceptionRepository creates a new instance of AvailabilityExceptionRepository
func NewAvailabilityExceptionRepository(db DBAdapter) domain.AvailabilityExceptionRepository {
	return &AvailabilityExceptionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new availability exception
func (r *AvailabilityExceptionRepository) Create(ctx context.Context, exception *domain.AvailabilityException) error {
	exceptionModel := mapExceptionDomainToModel(exception)

	if err := r.CreateWithAudit(ctx, &exceptionModel, &exception.CreatedBy); err != nil {
		return fmt.Errorf("failed to create availability exception: %w", err)
	}

	// Update the domain entity with any generated fields
	exception.ExceptionID = exceptionModel.ID
	exception.CreatedAt = exceptionModel.CreatedAt

	return nil
}

// GetByID retrieves an availability exception by ID
func (r *AvailabilityExceptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AvailabilityException, error) {
	var exceptionModel models.AvailabilityException

	err := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		First(&exceptionModel, "id = ?", id).Error

	if err != nil {
		return nil, r.HandleNotFound(err, "availability exception", id)
	}

	return mapExceptionModelToDomain(&exceptionModel), nil
}

// GetByStaff retrieves availability exceptions by staff ID
func (r *AvailabilityExceptionRepository) GetByStaff(ctx context.Context, staffID uuid.UUID) ([]*domain.AvailabilityException, error) {
	var exceptionModels []models.AvailabilityException

	err := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("staff_id = ?", staffID).
		Order("start_time ASC").
		Find(&exceptionModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get availability exceptions by staff ID: %w", err)
	}

	return mapExceptionModelsToDomainSlice(exceptionModels), nil
}

// GetByStaffAndDateRange retrieves availability exceptions by staff ID and date range
func (r *AvailabilityExceptionRepository) GetByStaffAndDateRange(ctx context.Context, staffID uuid.UUID, start, end time.Time) ([]*domain.AvailabilityException, error) {
	var exceptionModels []models.AvailabilityException

	query := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("staff_id = ?", staffID).
		Where("start_time BETWEEN ? AND ? OR end_time BETWEEN ? AND ? OR (start_time <= ? AND end_time >= ?)", 
			start, end, start, end, start, end)

	// Also include recurring exceptions
	query = query.Or("is_recurring = ? AND staff_id = ?", true, staffID)

	err := query.Order("start_time ASC").Find(&exceptionModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get availability exceptions by date range: %w", err)
	}

	return mapExceptionModelsToDomainSlice(exceptionModels), nil
}

// Update updates an availability exception
func (r *AvailabilityExceptionRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateAvailabilityExceptionInput, updatedBy uuid.UUID) error {
	// First find the exception to ensure it exists
	var exceptionModel models.AvailabilityException
	err := r.WithContext(ctx).First(&exceptionModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "availability exception", id)
	}

	// Apply updates from the input
	updates := map[string]interface{}{}

	if input.ExceptionType != nil {
		updates["exception_type"] = models.ExceptionType(*input.ExceptionType)
	}

	if input.StartTime != nil {
		updates["start_time"] = *input.StartTime
	}

	if input.EndTime != nil {
		updates["end_time"] = *input.EndTime
	}

	if input.IsFullDay != nil {
		updates["is_full_day"] = *input.IsFullDay
	}

	if input.IsRecurring != nil {
		updates["is_recurring"] = *input.IsRecurring
	}

	if input.RecurrenceRule != nil {
		updates["recurrence_rule"] = *input.RecurrenceRule
	}

	if input.Notes != nil {
		updates["notes"] = *input.Notes
	}

	// Perform the update with audit
	err = r.UpdateWithAudit(ctx, &exceptionModel, updates, updatedBy)
	if err != nil {
		return fmt.Errorf("failed to update availability exception: %w", err)
	}

	return nil
}

// Delete soft deletes an availability exception
func (r *AvailabilityExceptionRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// First find the exception to ensure it exists
	var exceptionModel models.AvailabilityException
	err := r.WithContext(ctx).First(&exceptionModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "availability exception", id)
	}

	// Perform soft delete with audit
	err = r.SoftDeleteWithAudit(ctx, &exceptionModel, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete availability exception: %w", err)
	}

	return nil
}

// ListByBusiness retrieves a paginated list of availability exceptions by business ID
func (r *AvailabilityExceptionRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*domain.AvailabilityException, error) {
	var exceptionModels []models.AvailabilityException

	// Apply pagination
	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("business_id = ?", businessID).
		Offset(offset).
		Limit(pageSize).
		Order("start_time ASC").
		Find(&exceptionModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list availability exceptions by business: %w", err)
	}

	return mapExceptionModelsToDomainSlice(exceptionModels), nil
}

// ListByBusinessAndDateRange retrieves a paginated list of availability exceptions by business ID and date range
func (r *AvailabilityExceptionRepository) ListByBusinessAndDateRange(ctx context.Context, businessID uuid.UUID, start, end time.Time, page, pageSize int) ([]*domain.AvailabilityException, error) {
	var exceptionModels []models.AvailabilityException

	// Apply pagination
	offset := r.CalculateOffset(page, pageSize)

	query := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("business_id = ?", businessID).
		Where("start_time BETWEEN ? AND ? OR end_time BETWEEN ? AND ? OR (start_time <= ? AND end_time >= ?)", 
			start, end, start, end, start, end)

	// Also include recurring exceptions
	query = query.Or("is_recurring = ? AND business_id = ?", true, businessID)

	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("start_time ASC").
		Find(&exceptionModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list availability exceptions by business and date range: %w", err)
	}

	return mapExceptionModelsToDomainSlice(exceptionModels), nil
}

// CountByBusiness counts availability exceptions by business ID
func (r *AvailabilityExceptionRepository) CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	var count int64

	err := r.WithContext(ctx).Model(&models.AvailabilityException{}).Where("business_id = ?", businessID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count availability exceptions by business: %w", err)
	}

	return count, nil
}

// Helper functions to map between domain entities and models

// mapExceptionDomainToModel converts a domain AvailabilityException entity to a model AvailabilityException
func mapExceptionDomainToModel(e *domain.AvailabilityException) *models.AvailabilityException {
	if e == nil {
		return nil
	}

	exceptionModel := &models.AvailabilityException{
		BaseModel: models.BaseModel{
			ID:        e.ExceptionID,
			CreatedAt: e.CreatedAt,
		},
		BusinessID:     e.BusinessID,
		StaffID:        e.StaffID,
		ExceptionType:  models.ExceptionType(e.ExceptionType),
		StartTime:      e.StartTime,
		EndTime:        e.EndTime,
		IsFullDay:      e.IsFullDay,
		IsRecurring:    e.IsRecurring,
		RecurrenceRule: e.RecurrenceRule,
		Notes:          e.Notes,
	}

	if e.CreatedBy != uuid.Nil {
		createdBy := e.CreatedBy
		exceptionModel.CreatedBy = &createdBy
	}

	if e.UpdatedAt != nil {
		exceptionModel.UpdatedAt = *e.UpdatedAt
	}

	if e.UpdatedBy != nil {
		exceptionModel.UpdatedBy = e.UpdatedBy
	}

	if e.DeletedAt != nil {
		deletedAt := gorm.DeletedAt{Time: *e.DeletedAt, Valid: true}
		exceptionModel.DeletedAt = deletedAt
	}

	if e.DeletedBy != nil {
		exceptionModel.DeletedBy = e.DeletedBy
	}

	return exceptionModel
}

// mapExceptionModelToDomain converts a model AvailabilityException to a domain AvailabilityException entity
func mapExceptionModelToDomain(e *models.AvailabilityException) *domain.AvailabilityException {
	if e == nil {
		return nil
	}

	exception := &domain.AvailabilityException{
		ExceptionID:    e.ID,
		BusinessID:     e.BusinessID,
		StaffID:        e.StaffID,
		ExceptionType:  string(e.ExceptionType),
		StartTime:      e.StartTime,
		EndTime:        e.EndTime,
		IsFullDay:      e.IsFullDay,
		IsRecurring:    e.IsRecurring,
		RecurrenceRule: e.RecurrenceRule,
		Notes:          e.Notes,
		CreatedAt:      e.CreatedAt,
	}

	if e.CreatedBy != nil {
		exception.CreatedBy = *e.CreatedBy
	}

	if !e.UpdatedAt.IsZero() {
		exception.UpdatedAt = &e.UpdatedAt
	}

	if e.UpdatedBy != nil {
		exception.UpdatedBy = e.UpdatedBy
	}

	if e.DeletedAt.Valid {
		deletedAt := e.DeletedAt.Time
		exception.DeletedAt = &deletedAt
	}

	if e.DeletedBy != nil {
		exception.DeletedBy = e.DeletedBy
	}

	// Map related entities if loaded
	if e.Staff.ID != uuid.Nil {
		exception.Staff = mapStaffModelToDomain(&e.Staff)
	}

	return exception
}

// mapExceptionModelsToDomainSlice converts a slice of model AvailabilityException to a slice of domain AvailabilityException entities
func mapExceptionModelsToDomainSlice(exceptionModels []models.AvailabilityException) []*domain.AvailabilityException {
	result := make([]*domain.AvailabilityException, len(exceptionModels))
	for i, model := range exceptionModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapExceptionModelToDomain(&modelCopy)
	}
	return result
}
