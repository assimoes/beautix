package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormServiceAssignmentRepository implements the domain.ServiceAssignmentRepository interface using GORM
type GormServiceAssignmentRepository struct {
	db *database.DB
}

// NewServiceAssignmentRepository creates a new instance of GormServiceAssignmentRepository
func NewServiceAssignmentRepository(db *database.DB) domain.ServiceAssignmentRepository {
	return &GormServiceAssignmentRepository{
		db: db,
	}
}

// Create creates a new service assignment
func (r *GormServiceAssignmentRepository) Create(ctx context.Context, assignment *domain.ServiceAssignment) error {
	assignmentModel := mapAssignmentDomainToModel(assignment)
	
	if err := r.db.WithContext(ctx).Create(&assignmentModel).Error; err != nil {
		return fmt.Errorf("failed to create service assignment: %w", err)
	}
	
	// Update the domain entity with any generated fields
	assignment.AssignmentID = assignmentModel.ID
	assignment.CreatedAt = assignmentModel.CreatedAt
	
	return nil
}

// GetByID retrieves a service assignment by ID
func (r *GormServiceAssignmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ServiceAssignment, error) {
	var assignmentModel models.ServiceAssignment
	
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		First(&assignmentModel, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("service assignment not found with ID %s", id)
		}
		return nil, fmt.Errorf("failed to get service assignment by ID: %w", err)
	}
	
	return mapAssignmentModelToDomain(&assignmentModel), nil
}

// GetByStaffAndService retrieves a service assignment by staff ID and service ID
func (r *GormServiceAssignmentRepository) GetByStaffAndService(ctx context.Context, staffID, serviceID uuid.UUID) (*domain.ServiceAssignment, error) {
	var assignmentModel models.ServiceAssignment
	
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("staff_id = ? AND service_id = ?", staffID, serviceID).
		First(&assignmentModel).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("service assignment not found for staff %s and service %s", staffID, serviceID)
		}
		return nil, fmt.Errorf("failed to get service assignment: %w", err)
	}
	
	return mapAssignmentModelToDomain(&assignmentModel), nil
}

// GetByStaff retrieves service assignments by staff ID
func (r *GormServiceAssignmentRepository) GetByStaff(ctx context.Context, staffID uuid.UUID) ([]*domain.ServiceAssignment, error) {
	var assignmentModels []models.ServiceAssignment
	
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("staff_id = ?", staffID).
		Find(&assignmentModels).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get service assignments by staff ID: %w", err)
	}
	
	return mapAssignmentModelsToDomainSlice(assignmentModels), nil
}

// GetByService retrieves service assignments by service ID
func (r *GormServiceAssignmentRepository) GetByService(ctx context.Context, serviceID uuid.UUID) ([]*domain.ServiceAssignment, error) {
	var assignmentModels []models.ServiceAssignment
	
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("service_id = ?", serviceID).
		Find(&assignmentModels).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get service assignments by service ID: %w", err)
	}
	
	return mapAssignmentModelsToDomainSlice(assignmentModels), nil
}

// Update updates a service assignment
func (r *GormServiceAssignmentRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateServiceAssignmentInput, updatedBy uuid.UUID) error {
	// First find the assignment to ensure it exists
	var assignmentModel models.ServiceAssignment
	err := r.db.WithContext(ctx).First(&assignmentModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("service assignment not found with ID %s", id)
		}
		return fmt.Errorf("failed to find service assignment for update: %w", err)
	}
	
	// Apply updates from the input
	updates := map[string]interface{}{
		"updated_at": time.Now(),
		"updated_by": updatedBy,
	}
	
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	
	// Perform the update
	err = r.db.WithContext(ctx).Model(&assignmentModel).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update service assignment: %w", err)
	}
	
	return nil
}

// Delete soft deletes a service assignment
func (r *GormServiceAssignmentRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// First find the assignment to ensure it exists
	var assignmentModel models.ServiceAssignment
	err := r.db.WithContext(ctx).First(&assignmentModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("service assignment not found with ID %s", id)
		}
		return fmt.Errorf("failed to find service assignment for deletion: %w", err)
	}
	
	// Set deleted by
	updates := map[string]interface{}{
		"deleted_by": deletedBy,
	}
	
	err = r.db.WithContext(ctx).Model(&assignmentModel).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update service assignment before deletion: %w", err)
	}
	
	// Perform soft delete
	err = r.db.WithContext(ctx).Delete(&assignmentModel).Error
	if err != nil {
		return fmt.Errorf("failed to delete service assignment: %w", err)
	}
	
	return nil
}

// ListByBusiness retrieves a paginated list of service assignments by business ID
func (r *GormServiceAssignmentRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*domain.ServiceAssignment, error) {
	var assignmentModels []models.ServiceAssignment
	
	// Apply pagination
	offset := (page - 1) * pageSize
	
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		Where("business_id = ?", businessID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&assignmentModels).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to list service assignments by business: %w", err)
	}
	
	return mapAssignmentModelsToDomainSlice(assignmentModels), nil
}

// CountByBusiness counts service assignments by business ID
func (r *GormServiceAssignmentRepository) CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	var count int64
	
	err := r.db.WithContext(ctx).Model(&models.ServiceAssignment{}).Where("business_id = ?", businessID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count service assignments by business: %w", err)
	}
	
	return count, nil
}

// Helper functions to map between domain entities and models

// mapAssignmentDomainToModel converts a domain ServiceAssignment entity to a model ServiceAssignment
func mapAssignmentDomainToModel(a *domain.ServiceAssignment) *models.ServiceAssignment {
	if a == nil {
		return nil
	}
	
	assignmentModel := &models.ServiceAssignment{
		BaseModel: models.BaseModel{
			ID:        a.AssignmentID,
			CreatedAt: a.CreatedAt,
		},
		BusinessID: a.BusinessID,
		StaffID:    a.StaffID,
		ServiceID:  a.ServiceID,
		IsActive:   a.IsActive,
	}
	
	if a.CreatedBy != uuid.Nil {
		createdBy := a.CreatedBy
		assignmentModel.CreatedBy = &createdBy
	}
	
	if a.UpdatedAt != nil {
		assignmentModel.UpdatedAt = *a.UpdatedAt
	}
	
	if a.UpdatedBy != nil {
		assignmentModel.UpdatedBy = a.UpdatedBy
	}
	
	if a.DeletedAt != nil {
		deletedAt := gorm.DeletedAt{Time: *a.DeletedAt, Valid: true}
		assignmentModel.DeletedAt = deletedAt
	}
	
	if a.DeletedBy != nil {
		assignmentModel.DeletedBy = a.DeletedBy
	}
	
	return assignmentModel
}

// mapAssignmentModelToDomain converts a model ServiceAssignment to a domain ServiceAssignment entity
func mapAssignmentModelToDomain(a *models.ServiceAssignment) *domain.ServiceAssignment {
	if a == nil {
		return nil
	}
	
	assignment := &domain.ServiceAssignment{
		AssignmentID: a.ID,
		BusinessID:   a.BusinessID,
		StaffID:      a.StaffID,
		ServiceID:    a.ServiceID,
		IsActive:     a.IsActive,
		CreatedAt:    a.CreatedAt,
	}
	
	if a.CreatedBy != nil {
		assignment.CreatedBy = *a.CreatedBy
	}
	
	if !a.UpdatedAt.IsZero() {
		assignment.UpdatedAt = &a.UpdatedAt
	}
	
	if a.UpdatedBy != nil {
		assignment.UpdatedBy = a.UpdatedBy
	}
	
	if a.DeletedAt.Valid {
		deletedAt := a.DeletedAt.Time
		assignment.DeletedAt = &deletedAt
	}
	
	if a.DeletedBy != nil {
		assignment.DeletedBy = a.DeletedBy
	}
	
	// Map related entities if loaded
	if a.Staff.ID != uuid.Nil {
		assignment.Staff = mapModelToDomain(&a.Staff)
	}
	
	return assignment
}

// mapAssignmentModelsToDomainSlice converts a slice of model ServiceAssignment to a slice of domain ServiceAssignment entities
func mapAssignmentModelsToDomainSlice(assignmentModels []models.ServiceAssignment) []*domain.ServiceAssignment {
	result := make([]*domain.ServiceAssignment, len(assignmentModels))
	for i, model := range assignmentModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapAssignmentModelToDomain(&modelCopy)
	}
	return result
}