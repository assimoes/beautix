package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
)

// StaffRepository implements the domain.StaffRepository interface using GORM
type StaffRepository struct {
	*BaseRepository
}

// NewStaffRepository creates a new instance of StaffRepository
func NewStaffRepository(db DBAdapter) domain.StaffRepository {
	return &StaffRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new staff member
func (r *StaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	staffModel := mapStaffDomainToModel(staff)

	// Option 1: Use CreateWithAudit helper method
	// This automatically sets CreatedBy and CreatedAt fields
	if err := r.CreateWithAudit(ctx, &staffModel, &staff.CreatedBy); err != nil {
		return fmt.Errorf("failed to create staff: %w", err)
	}

	// Update the domain entity with any generated fields
	staff.StaffID = staffModel.ID
	staff.CreatedAt = staffModel.CreatedAt

	return nil
}

// GetByID retrieves a staff member by ID
func (r *StaffRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Staff, error) {
	var staffModel models.Staff

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		First(&staffModel, "id = ?", id).Error

	if err != nil {
		return nil, r.HandleNotFound(err, "staff", id)
	}

	return mapStaffModelToDomain(&staffModel), nil
}

// GetByUserID retrieves staff members by user ID
func (r *StaffRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Staff, error) {
	var staffModels []models.Staff

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("user_id = ?", userID).
		Find(&staffModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get staff by user ID: %w", err)
	}

	return mapStaffModelsToDomainSlice(staffModels), nil
}

// GetByBusinessID retrieves staff members by business ID
func (r *StaffRepository) GetByBusinessID(ctx context.Context, businessID uuid.UUID) ([]*domain.Staff, error) {
	var staffModels []models.Staff

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("business_id = ?", businessID).
		Find(&staffModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get staff by business ID: %w", err)
	}

	return mapStaffModelsToDomainSlice(staffModels), nil
}

// Update updates a staff member
func (r *StaffRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateStaffInput, updatedBy uuid.UUID) error {
	// First find the staff to ensure it exists
	var staffModel models.Staff
	err := r.WithContext(ctx).First(&staffModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "staff", id)
	}

	// Apply updates from the input
	// Option 2: Manual audit field handling in updates map
	// The UpdateWithAudit method will automatically add updated_by and updated_at
	updates := map[string]interface{}{}

	if input.Position != nil {
		updates["position"] = *input.Position
	}

	if input.Bio != nil {
		updates["bio"] = *input.Bio
	}

	if input.SpecialtyAreas != nil {
		updates["specialty_areas"] = models.SpecialtyAreas(*input.SpecialtyAreas)
	}

	if input.ProfileImageURL != nil {
		updates["profile_image_url"] = *input.ProfileImageURL
	}

	if input.WorkingHours != nil {
		updates["working_hours"] = *input.WorkingHours
	}

	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if input.EmploymentType != nil {
		updates["employment_type"] = *input.EmploymentType
	}

	if input.JoinDate != nil {
		updates["join_date"] = *input.JoinDate
	}

	if input.EndDate != nil {
		updates["end_date"] = input.EndDate
	}

	if input.CommissionRate != nil {
		updates["commission_rate"] = *input.CommissionRate
	}

	// Option 2: Use UpdateWithAudit helper method
	// This automatically adds updated_by and updated_at to the updates map
	err = r.UpdateWithAudit(ctx, &staffModel, updates, updatedBy)
	if err != nil {
		return fmt.Errorf("failed to update staff: %w", err)
	}

	return nil
}

// Delete soft deletes a staff member
func (r *StaffRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// First find the staff to ensure it exists
	var staffModel models.Staff
	err := r.WithContext(ctx).First(&staffModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "staff", id)
	}

	// Option 3: Use SoftDeleteWithAudit helper method
	// This automatically sets deleted_by and deleted_at fields
	err = r.SoftDeleteWithAudit(ctx, &staffModel, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete staff: %w", err)
	}

	return nil
}

// List retrieves a paginated list of staff members
func (r *StaffRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Staff, error) {
	var staffModels []models.Staff

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&staffModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list staff: %w", err)
	}

	return mapStaffModelsToDomainSlice(staffModels), nil
}

// ListByBusiness retrieves a paginated list of staff members by business ID
func (r *StaffRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*domain.Staff, error) {
	var staffModels []models.Staff

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("business_id = ?", businessID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&staffModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list staff by business: %w", err)
	}

	return mapStaffModelsToDomainSlice(staffModels), nil
}

// Search searches for staff members by query within a business
func (r *StaffRepository) Search(ctx context.Context, businessID uuid.UUID, query string, page, pageSize int) ([]*domain.Staff, error) {
	var staffModels []models.Staff

	offset := r.CalculateOffset(page, pageSize)

	// Build the search query
	searchQuery := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("business_id = ?", businessID).
		Where("position ILIKE ? OR bio ILIKE ? OR specialty_areas::text ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC")

	// Join with user table for searching by user's name
	searchQuery = searchQuery.Joins("LEFT JOIN users ON users.id = staff.user_id").
		Where("users.first_name ILIKE ? OR users.last_name ILIKE ? OR CONCAT(users.first_name, ' ', users.last_name) ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")

	err := searchQuery.Find(&staffModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to search staff: %w", err)
	}

	return mapStaffModelsToDomainSlice(staffModels), nil
}

// Count counts all staff members
func (r *StaffRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.WithContext(ctx).Model(&models.Staff{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count staff: %w", err)
	}

	return count, nil
}

// CountByBusiness counts staff members by business ID
func (r *StaffRepository) CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	var count int64

	err := r.WithContext(ctx).Model(&models.Staff{}).Where("business_id = ?", businessID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count staff by business: %w", err)
	}

	return count, nil
}

// Helper functions to map between domain entities and models

// mapStaffDomainToModel converts a domain Staff entity to a model Staff
func mapStaffDomainToModel(s *domain.Staff) *models.Staff {
	if s == nil {
		return nil
	}

	staffModel := &models.Staff{
		BaseModel: models.BaseModel{
			ID:        s.StaffID,
			CreatedAt: s.CreatedAt,
		},
		BusinessID:      s.BusinessID,
		UserID:          s.UserID,
		Position:        s.Position,
		Bio:             s.Bio,
		SpecialtyAreas:  models.SpecialtyAreas(s.SpecialtyAreas),
		ProfileImageURL: s.ProfileImageURL,
		IsActive:        s.IsActive,
		EmploymentType:  models.StaffEmploymentType(s.EmploymentType),
		JoinDate:        s.JoinDate,
		CommissionRate:  s.CommissionRate,
	}

	// Handle optional/pointer fields
	if s.EndDate != nil {
		staffModel.EndDate = s.EndDate
	}

	if s.WorkingHours != nil {
		// Convert JSON bytes to WorkingHours
		var workingHours models.WorkingHours
		if err := json.Unmarshal(s.WorkingHours, &workingHours); err == nil {
			staffModel.WorkingHours = workingHours
		}
	}

	if s.CreatedBy != uuid.Nil {
		createdBy := s.CreatedBy
		staffModel.CreatedBy = &createdBy
	}

	if s.UpdatedAt != nil {
		staffModel.UpdatedAt = *s.UpdatedAt
	}

	if s.UpdatedBy != nil {
		staffModel.UpdatedBy = s.UpdatedBy
	}

	if s.UpdatedBy != nil {
		staffModel.UpdatedBy = s.UpdatedBy
	}

	return staffModel
}

// mapStaffModelToDomain converts a model Staff to a domain Staff entity
func mapStaffModelToDomain(s *models.Staff) *domain.Staff {
	if s == nil {
		return nil
	}

	staff := &domain.Staff{
		StaffID:         s.ID,
		BusinessID:      s.BusinessID,
		UserID:          s.UserID,
		Position:        s.Position,
		Bio:             s.Bio,
		SpecialtyAreas:  []string(s.SpecialtyAreas),
		ProfileImageURL: s.ProfileImageURL,
		IsActive:        s.IsActive,
		EmploymentType:  string(s.EmploymentType),
		JoinDate:        s.JoinDate,
		CommissionRate:  s.CommissionRate,
		CreatedAt:       s.CreatedAt,
	}

	// Convert WorkingHours to JSON bytes
	if workingHoursBytes, err := json.Marshal(s.WorkingHours); err == nil {
		staff.WorkingHours = workingHoursBytes
	}

	// Handle optional/pointer fields
	if s.EndDate != nil {
		staff.EndDate = s.EndDate
	}

	if s.CreatedBy != nil {
		staff.CreatedBy = *s.CreatedBy
	}

	if !s.UpdatedAt.IsZero() {
		staff.UpdatedAt = &s.UpdatedAt
	}

	if s.UpdatedBy != nil {
		staff.UpdatedBy = s.UpdatedBy
	}

	if s.DeletedAt.Valid {
		deletedAt := s.DeletedAt.Time
		staff.DeletedAt = &deletedAt
	}

	if s.DeletedBy != nil {
		staff.DeletedBy = s.DeletedBy
	}

	// Map related entities if loaded
	if s.User.ID != uuid.Nil {
		staff.User = mapUserModelToDomain(&s.User)
	}

	if s.Business.ID != uuid.Nil {
		staff.Business = mapBusinessModelToDomain(&s.Business)
	}

	return staff
}

// mapStaffModelsToDomainSlice converts a slice of model Staff to a slice of domain Staff entities
func mapStaffModelsToDomainSlice(staffModels []models.Staff) []*domain.Staff {
	result := make([]*domain.Staff, len(staffModels))
	for i, model := range staffModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapStaffModelToDomain(&modelCopy)
	}
	return result
}
