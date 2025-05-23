package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ClientLoyaltyMembershipRepository implements domain.ClientLoyaltyMembershipRepository
type ClientLoyaltyMembershipRepository struct {
	*BaseRepository
}

// NewClientLoyaltyMembershipRepository creates a new client loyalty membership repository
func NewClientLoyaltyMembershipRepository(db *gorm.DB) domain.ClientLoyaltyMembershipRepository {
	return &ClientLoyaltyMembershipRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// NewClientLoyaltyMembershipRepositoryTx creates a new client loyalty membership repository with transaction support
func NewClientLoyaltyMembershipRepositoryTx(base *BaseRepository) domain.ClientLoyaltyMembershipRepository {
	return &ClientLoyaltyMembershipRepository{
		BaseRepository: base,
	}
}

// Create creates a new client loyalty membership
func (r *ClientLoyaltyMembershipRepository) Create(ctx context.Context, membership *domain.ClientLoyaltyMembership) error {
	db := r.WithContext(ctx)

	query := `
		INSERT INTO client_loyalty_memberships (
			id, program_id, client_id, current_points, visits_count, total_spent,
			tier_level, progress, join_date, expiry_date, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result := db.Exec(query,
		membership.ID, membership.ProgramID, membership.ClientID, membership.CurrentPoints,
		membership.VisitsCount, membership.TotalSpent, membership.TierLevel, membership.Progress,
		membership.JoinDate, membership.ExpiryDate, membership.IsActive, membership.CreatedAt, membership.UpdatedAt,
	)

	if result.Error != nil {
		return fmt.Errorf("failed to create client loyalty membership: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a client loyalty membership by ID
func (r *ClientLoyaltyMembershipRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ClientLoyaltyMembership, error) {
	db := r.WithContext(ctx)

	query := `
		SELECT id, program_id, client_id, current_points, visits_count, total_spent,
			   tier_level, progress, join_date, expiry_date, is_active, created_at, updated_at
		FROM client_loyalty_memberships
		WHERE id = ? AND is_active = true
	`

	membership := &domain.ClientLoyaltyMembership{}
	err := db.Raw(query, id).Scan(membership).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get client loyalty membership by ID: %w", err)
	}

	return membership, nil
}

// GetByClientAndProgram retrieves a client loyalty membership by client and program IDs
func (r *ClientLoyaltyMembershipRepository) GetByClientAndProgram(ctx context.Context, clientID, programID uuid.UUID) (*domain.ClientLoyaltyMembership, error) {
	db := r.WithContext(ctx)

	query := `
		SELECT id, program_id, client_id, current_points, visits_count, total_spent,
			   tier_level, progress, join_date, expiry_date, is_active, created_at, updated_at
		FROM client_loyalty_memberships
		WHERE client_id = ? AND program_id = ? AND is_active = true
	`

	membership := &domain.ClientLoyaltyMembership{}
	err := db.Raw(query, clientID, programID).Scan(membership).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get client loyalty membership: %w", err)
	}

	return membership, nil
}

// Update updates a client loyalty membership
func (r *ClientLoyaltyMembershipRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateClientLoyaltyMembershipInput) error {
	db := r.WithContext(ctx)

	// Build dynamic update query
	setParts := []string{}
	args := []any{}

	if input.CurrentPoints != nil {
		setParts = append(setParts, "current_points = ?")
		args = append(args, *input.CurrentPoints)
	}
	if input.VisitsCount != nil {
		setParts = append(setParts, "visits_count = ?")
		args = append(args, *input.VisitsCount)
	}
	if input.TotalSpent != nil {
		setParts = append(setParts, "total_spent = ?")
		args = append(args, *input.TotalSpent)
	}
	if input.TierLevel != nil {
		setParts = append(setParts, "tier_level = ?")
		args = append(args, *input.TierLevel)
	}
	if input.Progress != nil {
		setParts = append(setParts, "progress = ?")
		args = append(args, *input.Progress)
	}
	if input.ExpiryDate != nil {
		setParts = append(setParts, "expiry_date = ?")
		args = append(args, input.ExpiryDate)
	}
	if input.IsActive != nil {
		setParts = append(setParts, "is_active = ?")
		args = append(args, *input.IsActive)
	}

	if len(setParts) == 0 {
		return nil // Nothing to update
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE client_loyalty_memberships 
		SET %s
		WHERE id = ?
	`, fmt.Sprintf("%s", setParts))

	result := db.Exec(query, args...)
	if result.Error != nil {
		return fmt.Errorf("failed to update client loyalty membership: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// AddPoints adds points to a membership and creates a transaction record
func (r *ClientLoyaltyMembershipRepository) AddPoints(ctx context.Context, membershipID uuid.UUID, points int, description string, createdBy uuid.UUID) error {
	db := r.WithContext(ctx)

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update membership points
	updateQuery := `
		UPDATE client_loyalty_memberships 
		SET current_points = current_points + ?, updated_at = NOW()
		WHERE id = ?
	`

	result := tx.Exec(updateQuery, points, membershipID)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add points to membership: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	// Create transaction record
	transactionQuery := `
		INSERT INTO loyalty_transactions (id, membership_id, transaction_type, points, description, created_at, created_by)
		VALUES (?, ?, 'earn', ?, ?, NOW(), ?)
	`

	transactionID := uuid.New()
	result = tx.Exec(transactionQuery, transactionID, membershipID, points, description, createdBy)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create loyalty transaction: %w", result.Error)
	}

	return tx.Commit().Error
}

// RedeemPoints redeems points from a membership and creates a transaction record
func (r *ClientLoyaltyMembershipRepository) RedeemPoints(ctx context.Context, membershipID uuid.UUID, points int, description string, createdBy uuid.UUID) error {
	db := r.WithContext(ctx)

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if membership has enough points
	var currentPoints int
	checkQuery := `SELECT current_points FROM client_loyalty_memberships WHERE id = ?`
	err := tx.Raw(checkQuery, membershipID).Scan(&currentPoints).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check current points: %w", err)
	}

	if currentPoints < points {
		tx.Rollback()
		return fmt.Errorf("insufficient points: have %d, need %d", currentPoints, points)
	}

	// Update membership points
	updateQuery := `
		UPDATE client_loyalty_memberships 
		SET current_points = current_points - ?, updated_at = NOW()
		WHERE id = ?
	`

	result := tx.Exec(updateQuery, points, membershipID)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to redeem points from membership: %w", result.Error)
	}

	// Create transaction record (negative points for redemption)
	transactionQuery := `
		INSERT INTO loyalty_transactions (id, membership_id, transaction_type, points, description, created_at, created_by)
		VALUES (?, ?, 'redeem', ?, ?, NOW(), ?)
	`

	transactionID := uuid.New()
	result = tx.Exec(transactionQuery, transactionID, membershipID, -points, description, createdBy)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create loyalty transaction: %w", result.Error)
	}

	return tx.Commit().Error
}

// ListByClient lists all loyalty memberships for a client
func (r *ClientLoyaltyMembershipRepository) ListByClient(ctx context.Context, clientID uuid.UUID) ([]*domain.ClientLoyaltyMembership, error) {
	db := r.WithContext(ctx)

	query := `
		SELECT id, program_id, client_id, current_points, visits_count, total_spent,
			   tier_level, progress, join_date, expiry_date, is_active, created_at, updated_at
		FROM client_loyalty_memberships
		WHERE client_id = ? AND is_active = true
		ORDER BY created_at DESC
	`

	var memberships []*domain.ClientLoyaltyMembership
	err := db.Raw(query, clientID).Scan(&memberships).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list client loyalty memberships: %w", err)
	}

	return memberships, nil
}

// ListByProgram lists all clients in a loyalty program with pagination
func (r *ClientLoyaltyMembershipRepository) ListByProgram(ctx context.Context, programID uuid.UUID, page, pageSize int) ([]*domain.ClientLoyaltyMembership, error) {
	db := r.WithContext(ctx)

	offset := (page - 1) * pageSize
	query := `
		SELECT id, program_id, client_id, current_points, visits_count, total_spent,
			   tier_level, progress, join_date, expiry_date, is_active, created_at, updated_at
		FROM client_loyalty_memberships
		WHERE program_id = ? AND is_active = true
		ORDER BY current_points DESC, created_at DESC
		LIMIT ? OFFSET ?
	`

	var memberships []*domain.ClientLoyaltyMembership
	err := db.Raw(query, programID, pageSize, offset).Scan(&memberships).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list loyalty memberships by program: %w", err)
	}

	return memberships, nil
}

// Count counts total client loyalty memberships
func (r *ClientLoyaltyMembershipRepository) Count(ctx context.Context) (int64, error) {
	db := r.WithContext(ctx)

	var count int64
	err := db.Raw("SELECT COUNT(*) FROM client_loyalty_memberships WHERE is_active = true").Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count client loyalty memberships: %w", err)
	}

	return count, nil
}

// CountByProgram counts client loyalty memberships by program
func (r *ClientLoyaltyMembershipRepository) CountByProgram(ctx context.Context, programID uuid.UUID) (int64, error) {
	db := r.WithContext(ctx)

	var count int64
	err := db.Raw("SELECT COUNT(*) FROM client_loyalty_memberships WHERE program_id = ? AND is_active = true", programID).Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count client loyalty memberships by program: %w", err)
	}

	return count, nil
}
