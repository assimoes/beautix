package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServiceAssignmentRepository implements the domain.ServiceAssignmentRepository interface using GORM
type ServiceAssignmentRepository struct {
	*BaseRepository
}

// NewServiceAssignmentRepository creates a new instance of ServiceAssignmentRepository
func NewServiceAssignmentRepository(db DBAdapter) domain.ServiceAssignmentRepository {
	return &ServiceAssignmentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new service assignment
func (r *ServiceAssignmentRepository) Create(ctx context.Context, assignment *domain.ServiceAssignment) error {
	assignmentModel := mapAssignmentDomainToModel(assignment)
	
	if err := r.CreateWithAudit(ctx, &assignmentModel, &assignment.CreatedBy); err != nil {
		return fmt.Errorf("failed to create service assignment: %w", err)
	}
	
	// Update the domain entity with any generated fields
	assignment.AssignmentID = assignmentModel.ID
	assignment.CreatedAt = assignmentModel.CreatedAt
	
	return nil
}

// GetByID retrieves a service assignment by ID
func (r *ServiceAssignmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ServiceAssignment, error) {
	var assignmentModel models.ServiceAssignment
	
	err := r.WithContext(ctx).
		Preload("Staff").
		Preload("Staff.User").
		First(&assignmentModel, "id = ?", id).Error
	
	if err != nil {
		return nil, r.HandleNotFound(err, "service assignment", id)
	}
	
	return mapAssignmentModelToDomain(&assignmentModel), nil
}

// GetByStaffAndService retrieves a service assignment by staff ID and service ID
func (r *ServiceAssignmentRepository) GetByStaffAndService(ctx context.Context, staffID, serviceID uuid.UUID) (*domain.ServiceAssignment, error) {
	var assignmentModel models.ServiceAssignment
	
	err := r.WithContext(ctx).
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
func (r *ServiceAssignmentRepository) GetByStaff(ctx context.Context, staffID uuid.UUID) ([]*domain.ServiceAssignment, error) {
	var assignmentModels []models.ServiceAssignment
	
	err := r.WithContext(ctx).
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
func (r *ServiceAssignmentRepository) GetByService(ctx context.Context, serviceID uuid.UUID) ([]*domain.ServiceAssignment, error) {
	var assignmentModels []models.ServiceAssignment
	
	err := r.WithContext(ctx).
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
func (r *ServiceAssignmentRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateServiceAssignmentInput, updatedBy uuid.UUID) error {
	// First find the assignment to ensure it exists
	var assignmentModel models.ServiceAssignment
	err := r.WithContext(ctx).First(&assignmentModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "service assignment", id)
	}
	
	// Apply updates from the input
	updates := map[string]interface{}{}
	
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	
	// Perform the update with audit
	err = r.UpdateWithAudit(ctx, &assignmentModel, updates, updatedBy)
	if err != nil {
		return fmt.Errorf("failed to update service assignment: %w", err)
	}
	
	return nil
}

// Delete soft deletes a service assignment
func (r *ServiceAssignmentRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// First find the assignment to ensure it exists
	var assignmentModel models.ServiceAssignment
	err := r.WithContext(ctx).First(&assignmentModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "service assignment", id)
	}
	
	// Perform soft delete with audit
	err = r.SoftDeleteWithAudit(ctx, &assignmentModel, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete service assignment: %w", err)
	}
	
	return nil
}

// ListByBusiness retrieves a paginated list of service assignments by business ID
func (r *ServiceAssignmentRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*domain.ServiceAssignment, error) {
	var assignmentModels []models.ServiceAssignment
	
	// Apply pagination
	offset := r.CalculateOffset(page, pageSize)
	
	err := r.WithContext(ctx).
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
func (r *ServiceAssignmentRepository) CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	var count int64
	
	err := r.WithContext(ctx).Model(&models.ServiceAssignment{}).Where("business_id = ?", businessID).Count(&count).Error
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
		assignment.Staff = mapStaffModelToDomain(&a.Staff)
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