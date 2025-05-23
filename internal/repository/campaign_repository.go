package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CampaignRepository implements domain.CampaignRepository
type CampaignRepository struct {
	*BaseRepository
}

// NewCampaignRepository creates a new campaign repository
func NewCampaignRepository(db *gorm.DB) domain.CampaignRepository {
	return &CampaignRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// NewCampaignRepositoryTx creates a new campaign repository with transaction support
func NewCampaignRepositoryTx(base *BaseRepository) domain.CampaignRepository {
	return &CampaignRepository{
		BaseRepository: base,
	}
}

// Create creates a new campaign
func (r *CampaignRepository) Create(ctx context.Context, campaign *domain.Campaign) error {
	db := r.WithContext(ctx)
	
	query := `
		INSERT INTO campaigns (
			id, business_id, name, description, campaign_type, target_audience,
			offer_type, offer_details, start_date, end_date, is_active, 
			message_template, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result := db.Exec(query,
		campaign.ID, campaign.BusinessID, campaign.Name, campaign.Description,
		campaign.CampaignType, campaign.TargetAudience, campaign.OfferType, campaign.OfferDetails,
		campaign.StartDate, campaign.EndDate, campaign.IsActive, campaign.MessageTemplate,
		campaign.CreatedAt, campaign.CreatedBy, campaign.UpdatedAt, campaign.UpdatedBy,
		campaign.DeletedAt, campaign.DeletedBy,
	)
	
	if result.Error != nil {
		return fmt.Errorf("failed to create campaign: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a campaign by ID
func (r *CampaignRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Campaign, error) {
	db := r.WithContext(ctx)
	
	query := `
		SELECT id, business_id, name, description, campaign_type, target_audience,
			   offer_type, offer_details, start_date, end_date, is_active,
			   message_template, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM campaigns
		WHERE id = ? AND deleted_at IS NULL
	`
	
	campaign := &domain.Campaign{}
	err := db.Raw(query, id).Scan(campaign).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get campaign by ID: %w", err)
	}

	return campaign, nil
}

// Update updates a campaign
func (r *CampaignRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateCampaignInput, updatedBy uuid.UUID) error {
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
	if input.CampaignType != nil {
		setParts = append(setParts, "campaign_type = ?")
		args = append(args, *input.CampaignType)
	}
	if input.TargetAudience != nil {
		setParts = append(setParts, "target_audience = ?")
		args = append(args, *input.TargetAudience)
	}
	if input.OfferType != nil {
		setParts = append(setParts, "offer_type = ?")
		args = append(args, *input.OfferType)
	}
	if input.OfferDetails != nil {
		setParts = append(setParts, "offer_details = ?")
		args = append(args, *input.OfferDetails)
	}
	if input.StartDate != nil {
		setParts = append(setParts, "start_date = ?")
		args = append(args, input.StartDate)
	}
	if input.EndDate != nil {
		setParts = append(setParts, "end_date = ?")
		args = append(args, input.EndDate)
	}
	if input.IsActive != nil {
		setParts = append(setParts, "is_active = ?")
		args = append(args, *input.IsActive)
	}
	if input.MessageTemplate != nil {
		setParts = append(setParts, "message_template = ?")
		args = append(args, *input.MessageTemplate)
	}
	
	if len(setParts) == 0 {
		return nil // Nothing to update
	}
	
	setParts = append(setParts, "updated_by = ?", "updated_at = NOW()")
	args = append(args, updatedBy, id)
	
	query := fmt.Sprintf(`
		UPDATE campaigns 
		SET %s
		WHERE id = ? AND deleted_at IS NULL
	`, fmt.Sprintf("%s", setParts))
	
	result := db.Exec(query, args...)
	if result.Error != nil {
		return fmt.Errorf("failed to update campaign: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete soft deletes a campaign
func (r *CampaignRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	db := r.WithContext(ctx)
	
	query := `
		UPDATE campaigns 
		SET deleted_by = ?, deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	
	result := db.Exec(query, deletedBy, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete campaign: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ListByBusiness lists campaigns by business with pagination
func (r *CampaignRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*domain.Campaign, error) {
	db := r.WithContext(ctx)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT id, business_id, name, description, campaign_type, target_audience,
			   offer_type, offer_details, start_date, end_date, is_active,
			   message_template, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM campaigns
		WHERE business_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	var campaigns []*domain.Campaign
	err := db.Raw(query, businessID, pageSize, offset).Scan(&campaigns).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns by business: %w", err)
	}

	return campaigns, nil
}

// ListActiveByBusiness lists active campaigns by business with pagination
func (r *CampaignRepository) ListActiveByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*domain.Campaign, error) {
	db := r.WithContext(ctx)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT id, business_id, name, description, campaign_type, target_audience,
			   offer_type, offer_details, start_date, end_date, is_active,
			   message_template, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM campaigns
		WHERE business_id = ? AND is_active = true AND deleted_at IS NULL
			AND start_date <= NOW() AND end_date >= NOW()
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	var campaigns []*domain.Campaign
	err := db.Raw(query, businessID, pageSize, offset).Scan(&campaigns).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list active campaigns by business: %w", err)
	}

	return campaigns, nil
}

// Count counts total campaigns
func (r *CampaignRepository) Count(ctx context.Context) (int64, error) {
	db := r.WithContext(ctx)
	
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM campaigns WHERE deleted_at IS NULL").Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count campaigns: %w", err)
	}

	return count, nil
}

// CountByBusiness counts campaigns by business
func (r *CampaignRepository) CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	db := r.WithContext(ctx)
	
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM campaigns WHERE business_id = ? AND deleted_at IS NULL", businessID).Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count campaigns by business: %w", err)
	}

	return count, nil
}

// CountActiveByBusiness counts active campaigns by business
func (r *CampaignRepository) CountActiveByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error) {
	db := r.WithContext(ctx)
	
	query := `
		SELECT COUNT(*) FROM campaigns 
		WHERE business_id = ? AND is_active = true AND deleted_at IS NULL
			AND start_date <= NOW() AND end_date >= NOW()
	`
	
	var count int64
	err := db.Raw(query, businessID).Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count active campaigns by business: %w", err)
	}

	return count, nil
}