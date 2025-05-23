package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoyaltyProgramRepository implements domain.LoyaltyProgramRepository
type LoyaltyProgramRepository struct {
	*BaseRepository
}

// NewLoyaltyProgramRepository creates a new loyalty program repository
func NewLoyaltyProgramRepository(db *gorm.DB) domain.LoyaltyProgramRepository {
	return &LoyaltyProgramRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// NewLoyaltyProgramRepositoryTx creates a new loyalty program repository with transaction support
func NewLoyaltyProgramRepositoryTx(base *BaseRepository) domain.LoyaltyProgramRepository {
	return &LoyaltyProgramRepository{
		BaseRepository: base,
	}
}

// Create creates a new loyalty program
func (r *LoyaltyProgramRepository) Create(ctx context.Context, program *domain.LoyaltyProgram) error {
	db := r.WithContext(ctx)
	
	// Use raw SQL to match the database schema exactly
	query := `
		INSERT INTO loyalty_programs (
			id, business_id, name, description, program_type, rules, 
			reward_type, reward_value, is_active, start_date, end_date, 
			created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result := db.Exec(query,
		program.ID, program.BusinessID, program.Name, program.Description,
		program.ProgramType, program.Rules, program.RewardType, program.RewardValue,
		program.IsActive, program.StartDate, program.EndDate,
		program.CreatedAt, program.CreatedBy, program.UpdatedAt, program.UpdatedBy,
		program.DeletedAt, program.DeletedBy,
	)
	
	if result.Error != nil {
		return fmt.Errorf("failed to create loyalty program: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a loyalty program by ID
func (r *LoyaltyProgramRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.LoyaltyProgram, error) {
	db := r.WithContext(ctx)
	
	query := `
		SELECT id, business_id, name, description, program_type, rules, 
			   reward_type, reward_value, is_active, start_date, end_date,
			   created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM loyalty_programs 
		WHERE id = ? AND deleted_at IS NULL
	`
	
	program := &domain.LoyaltyProgram{}
	err := db.Raw(query, id).Scan(program).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get loyalty program by ID: %w", err)
	}

	return program, nil
}

// Update updates a loyalty program
func (r *LoyaltyProgramRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateLoyaltyProgramInput, updatedBy uuid.UUID) error {
	db := r.WithContext(ctx)
	
	// Build dynamic update query
	setParts := []string{}
	args := []any{}
	
	if input.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *input.Name)
	}
	if input.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *input.Description)
	}
	if input.ProgramType != nil {
		setParts = append(setParts, "program_type = ?")
		args = append(args, *input.ProgramType)
	}
	if input.Rules != nil {
		setParts = append(setParts, "rules = ?")
		args = append(args, *input.Rules)
	}
	if input.RewardType != nil {
		setParts = append(setParts, "reward_type = ?")
		args = append(args, *input.RewardType)
	}
	if input.RewardValue != nil {
		setParts = append(setParts, "reward_value = ?")
		args = append(args, *input.RewardValue)
	}
	if input.IsActive != nil {
		setParts = append(setParts, "is_active = ?")
		args = append(args, *input.IsActive)
	}
	if input.StartDate != nil {
		setParts = append(setParts, "start_date = ?")
		args = append(args, input.StartDate)
	}
	if input.EndDate != nil {
		setParts = append(setParts, "end_date = ?")
		args = append(args, input.EndDate)
	}
	
	if len(setParts) == 0 {
		return nil // Nothing to update
	}
	
	setParts = append(setParts, "updated_by = ?", "updated_at = NOW()")
	args = append(args, updatedBy, id)
	
	query := fmt.Sprintf(`
		UPDATE loyalty_programs 
		SET %s
		WHERE id = ? AND deleted_at IS NULL
	`, fmt.Sprintf("%s", setParts))
	
	result := db.Exec(query, args...)
	if result.Error != nil {
		return fmt.Errorf("failed to update loyalty program: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete soft deletes a loyalty program
func (r *LoyaltyProgramRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	db := r.WithContext(ctx)
	
	query := `
		UPDATE loyalty_programs 
		SET deleted_by = ?, deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	
	result := db.Exec(query, deletedBy, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete loyalty program: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ListByBusiness lists loyalty programs by business with pagination
func (r *LoyaltyProgramRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*domain.LoyaltyProgram, error) {
	db := r.WithContext(ctx)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT id, business_id, name, description, program_type, rules, 
			   reward_type, reward_value, is_active, start_date, end_date,
			   created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM loyalty_programs 
		WHERE business_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	var programs []*domain.LoyaltyProgram
	err := db.Raw(query, businessID, pageSize, offset).Scan(&programs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list loyalty programs by business: %w", err)
	}

	return programs, nil
}

// ListActiveBybusiness lists active loyalty programs by business with pagination
func (r *LoyaltyProgramRepository) ListActiveBybusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*domain.LoyaltyProgram, error) {
	db := r.WithContext(ctx)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT id, business_id, name, description, program_type, rules, 
			   reward_type, reward_value, is_active, start_date, end_date,
			   created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM loyalty_programs 
		WHERE business_id = ? AND is_active = true AND deleted_at IS NULL
			AND (start_date IS NULL OR start_date <= NOW())
			AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	var programs []*domain.LoyaltyProgram
	err := db.Raw(query, businessID, pageSize, offset).Scan(&programs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list active loyalty programs by business: %w", err)
	}

	return programs, nil
}

// Count counts total loyalty programs
func (r *LoyaltyProgramRepository) Count(ctx context.Context) (int64, error) {
	db := r.WithContext(ctx)
	
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM loyalty_programs WHERE deleted_at IS NULL").Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count loyalty programs: %w", err)
	}

	return count, nil
}

// CountByBusiness counts loyalty programs by business
func (r *LoyaltyProgramRepository) CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	db := r.WithContext(ctx)
	
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM loyalty_programs WHERE business_id = ? AND deleted_at IS NULL", businessID).Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count loyalty programs by business: %w", err)
	}

	return count, nil
}

// CountActiveByBusiness counts active loyalty programs by business
func (r *LoyaltyProgramRepository) CountActiveByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	db := r.WithContext(ctx)
	
	query := `
		SELECT COUNT(*) FROM loyalty_programs 
		WHERE business_id = ? AND is_active = true AND deleted_at IS NULL
			AND (start_date IS NULL OR start_date <= NOW())
			AND (end_date IS NULL OR end_date >= NOW())
	`
	
	var count int64
	err := db.Raw(query, businessID).Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count active loyalty programs by business: %w", err)
	}

	return count, nil
}