package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoyaltyTransactionRepository implements domain.LoyaltyTransactionRepository
type LoyaltyTransactionRepository struct {
	*BaseRepository
}

// NewLoyaltyTransactionRepository creates a new loyalty transaction repository
func NewLoyaltyTransactionRepository(db *gorm.DB) domain.LoyaltyTransactionRepository {
	return &LoyaltyTransactionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// NewLoyaltyTransactionRepositoryTx creates a new loyalty transaction repository with transaction support
func NewLoyaltyTransactionRepositoryTx(base *BaseRepository) domain.LoyaltyTransactionRepository {
	return &LoyaltyTransactionRepository{
		BaseRepository: base,
	}
}

// Create creates a new loyalty transaction
func (r *LoyaltyTransactionRepository) Create(ctx context.Context, transaction *domain.LoyaltyTransaction) error {
	db := r.WithContext(ctx)
	
	query := `
		INSERT INTO loyalty_transactions (
			id, membership_id, appointment_id, transaction_type, points, description, created_at, created_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result := db.Exec(query,
		transaction.ID, transaction.MembershipID, transaction.AppointmentID,
		transaction.TransactionType, transaction.Points, transaction.Description,
		transaction.CreatedAt, transaction.CreatedBy,
	)
	
	if result.Error != nil {
		return fmt.Errorf("failed to create loyalty transaction: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a loyalty transaction by ID
func (r *LoyaltyTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.LoyaltyTransaction, error) {
	db := r.WithContext(ctx)
	
	query := `
		SELECT id, membership_id, appointment_id, transaction_type, points, description, created_at, created_by
		FROM loyalty_transactions
		WHERE id = ?
	`
	
	transaction := &domain.LoyaltyTransaction{}
	err := db.Raw(query, id).Scan(transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get loyalty transaction by ID: %w", err)
	}

	return transaction, nil
}

// ListByMembership lists all transactions for a loyalty membership with pagination
func (r *LoyaltyTransactionRepository) ListByMembership(ctx context.Context, membershipID uuid.UUID, page, pageSize int) ([]*domain.LoyaltyTransaction, error) {
	db := r.WithContext(ctx)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT id, membership_id, appointment_id, transaction_type, points, description, created_at, created_by
		FROM loyalty_transactions
		WHERE membership_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	var transactions []*domain.LoyaltyTransaction
	err := db.Raw(query, membershipID, pageSize, offset).Scan(&transactions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list loyalty transactions by membership: %w", err)
	}

	return transactions, nil
}

// ListByClient lists all transactions for a client across all their memberships with pagination
func (r *LoyaltyTransactionRepository) ListByClient(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*domain.LoyaltyTransaction, error) {
	db := r.WithContext(ctx)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT lt.id, lt.membership_id, lt.appointment_id, lt.transaction_type, 
			   lt.points, lt.description, lt.created_at, lt.created_by
		FROM loyalty_transactions lt
		INNER JOIN client_loyalty_memberships clm ON lt.membership_id = clm.id
		WHERE clm.client_id = ?
		ORDER BY lt.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	var transactions []*domain.LoyaltyTransaction
	err := db.Raw(query, clientID, pageSize, offset).Scan(&transactions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list loyalty transactions by client: %w", err)
	}

	return transactions, nil
}

// Count counts total loyalty transactions
func (r *LoyaltyTransactionRepository) Count(ctx context.Context) (int64, error) {
	db := r.WithContext(ctx)
	
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM loyalty_transactions").Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count loyalty transactions: %w", err)
	}

	return count, nil
}

// CountByMembership counts loyalty transactions by membership
func (r *LoyaltyTransactionRepository) CountByMembership(ctx context.Context, membershipID uuid.UUID) (int64, error) {
	db := r.WithContext(ctx)
	
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM loyalty_transactions WHERE membership_id = ?", membershipID).Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count loyalty transactions by membership: %w", err)
	}

	return count, nil
}