package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormStaffRepository implements the domain.StaffRepository interface using GORM
type GormStaffRepository struct {
	db *database.DB
}

// NewStaffRepository creates a new instance of GormStaffRepository
func NewStaffRepository(db *database.DB) domain.StaffRepository {
	return &GormStaffRepository{
		db: db,
	}
}

// Create creates a new staff member
func (r *GormStaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	staffModel := mapDomainToModel(staff)
	
	if err := r.db.WithContext(ctx).Create(&staffModel).Error; err != nil {
		return fmt.Errorf("failed to create staff: %w", err)
	}
	
	// Update the domain entity with any generated fields
	staff.StaffID = staffModel.ID
	staff.CreatedAt = staffModel.CreatedAt
	
	return nil
}

// GetByID retrieves a staff member by ID
func (r *GormStaffRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Staff, error) {
	var staffModel models.Staff
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Business").
		First(&staffModel, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("staff not found with ID %s", id)
		}
		return nil, fmt.Errorf("failed to get staff by ID: %w", err)
	}
	
	return mapModelToDomain(&staffModel), nil
}

// GetByUserID retrieves staff members by user ID
func (r *GormStaffRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Staff, error) {
	var staffModels []models.Staff
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("user_id = ?", userID).
		Find(&staffModels).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get staff by user ID: %w", err)
	}
	
	// Map staff models to domain entities
	return mapModelsToDomainSlice(staffModels), nil
}

// GetByBusinessID retrieves staff members by business ID
func (r *GormStaffRepository) GetByBusinessID(ctx context.Context, businessID uuid.UUID) ([]*domain.Staff, error) {
	var staffModels []models.Staff
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("business_id = ?", businessID).
		Find(&staffModels).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get staff by business ID: %w", err)
	}
	
	// Map staff models to domain entities
	return mapModelsToDomainSlice(staffModels), nil
}

// Update updates a staff member
func (r *GormStaffRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateStaffInput, updatedBy uuid.UUID) error {
	// First find the staff to ensure it exists
	var staffModel models.Staff
	err := r.db.WithContext(ctx).First(&staffModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("staff not found with ID %s", id)
		}
		return fmt.Errorf("failed to find staff for update: %w", err)
	}
	
	// Apply updates from the input
	updates := map[string]interface{}{
		"updated_at": time.Now(),
		"updated_by": updatedBy,
	}
	
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
	
	// Perform the update
	err = r.db.WithContext(ctx).Model(&staffModel).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update staff: %w", err)
	}
	
	return nil
}

// Delete soft deletes a staff member
func (r *GormStaffRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// First find the staff to ensure it exists
	var staffModel models.Staff
	err := r.db.WithContext(ctx).First(&staffModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("staff not found with ID %s", id)
		}
		return fmt.Errorf("failed to find staff for deletion: %w", err)
	}
	
	// Set deleted by
	updates := map[string]interface{}{
		"deleted_by": deletedBy,
	}
	
	err = r.db.WithContext(ctx).Model(&staffModel).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update staff before deletion: %w", err)
	}
	
	// Perform soft delete
	err = r.db.WithContext(ctx).Delete(&staffModel).Error
	if err != nil {
		return fmt.Errorf("failed to delete staff: %w", err)
	}
	
	return nil
}

// List retrieves a paginated list of staff members
func (r *GormStaffRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Staff, error) {
	var staffModels []models.Staff
	
	// Apply pagination
	offset := (page - 1) * pageSize
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&staffModels).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to list staff: %w", err)
	}
	
	return mapModelsToDomainSlice(staffModels), nil
}

// ListByBusiness retrieves a paginated list of staff members by business ID
func (r *GormStaffRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*domain.Staff, error) {
	var staffModels []models.Staff
	
	// Apply pagination
	offset := (page - 1) * pageSize
	
	err := r.db.WithContext(ctx).
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
	
	return mapModelsToDomainSlice(staffModels), nil
}

// Search searches for staff members by query within a business
func (r *GormStaffRepository) Search(ctx context.Context, businessID uuid.UUID, query string, page, pageSize int) ([]*domain.Staff, error) {
	var staffModels []models.Staff
	
	// Apply pagination
	offset := (page - 1) * pageSize
	
	// Build the search query
	searchQuery := r.db.WithContext(ctx).
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
	
	return mapModelsToDomainSlice(staffModels), nil
}

// Count counts all staff members
func (r *GormStaffRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	
	err := r.db.WithContext(ctx).Model(&models.Staff{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count staff: %w", err)
	}
	
	return count, nil
}

// CountByBusiness counts staff members by business ID
func (r *GormStaffRepository) CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	var count int64
	
	err := r.db.WithContext(ctx).Model(&models.Staff{}).Where("business_id = ?", businessID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count staff by business: %w", err)
	}
	
	return count, nil
}

// Helper functions to map between domain entities and models

// mapDomainToModel converts a domain Staff entity to a model Staff
func mapDomainToModel(s *domain.Staff) *models.Staff {
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
	
	if s.DeletedAt != nil {
		deletedAt := gorm.DeletedAt{Time: *s.DeletedAt, Valid: true}
		staffModel.DeletedAt = deletedAt
	}
	
	if s.DeletedBy != nil {
		staffModel.DeletedBy = s.DeletedBy
	}
	
	return staffModel
}

// mapModelToDomain converts a model Staff to a domain Staff entity
func mapModelToDomain(s *models.Staff) *domain.Staff {
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

// mapModelsToDomainSlice converts a slice of model Staff to a slice of domain Staff entities
func mapModelsToDomainSlice(staffModels []models.Staff) []*domain.Staff {
	result := make([]*domain.Staff, len(staffModels))
	for i, model := range staffModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapModelToDomain(&modelCopy)
	}
	return result
}

// mapUserModelToDomain converts a model User to a domain User entity
func mapUserModelToDomain(u *models.User) *domain.User {
	if u == nil || u.ID == uuid.Nil {
		return nil
	}
	
	return &domain.User{
		UserID:    u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		IsActive:  u.IsActive,
		Role:      string(u.Role),
	}
}

// mapBusinessModelToDomain converts a model Business to a domain Business entity
func mapBusinessModelToDomain(b *models.Business) *domain.Business {
	if b == nil || b.ID == uuid.Nil {
		return nil
	}
	
	return &domain.Business{
		BusinessID:   b.ID,
		OwnerID:      b.UserID,
		BusinessName: b.Name,
		City:         b.City,
		Country:      b.Country,
		Phone:        b.Phone,
		Email:        b.Email,
		IsActive:     b.IsActive,
	}
}