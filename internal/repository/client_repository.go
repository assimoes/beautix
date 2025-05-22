package repository

import (
	"context"
	"fmt"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
)

// ClientRepository implements the domain.ClientRepository interface using GORM
type ClientRepository struct {
	*BaseRepository
}

// NewClientRepository creates a new instance of ClientRepository
func NewClientRepository(db DBAdapter) domain.ClientRepository {
	return &ClientRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new client
func (r *ClientRepository) Create(ctx context.Context, client *domain.Client) error {
	clientModel := mapClientDomainToModel(client)

	if err := r.CreateWithAudit(ctx, &clientModel, client.CreatedBy); err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Update the domain entity with any generated fields
	client.ID = clientModel.ID
	client.CreatedAt = clientModel.CreatedAt

	return nil
}

// GetByID retrieves a client by ID
func (r *ClientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	var clientModel models.Client

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		First(&clientModel, "id = ?", id).Error

	if err != nil {
		return nil, r.HandleNotFound(err, "client", id)
	}

	return mapClientModelToDomain(&clientModel), nil
}

// GetByUserID retrieves clients by user ID
func (r *ClientRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Client, error) {
	var clientModels []models.Client

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("user_id = ?", userID).
		Find(&clientModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get clients by user ID: %w", err)
	}

	return mapClientModelsToDomainSlice(clientModels), nil
}

// GetByProviderAndEmail retrieves a client by provider ID and email
func (r *ClientRepository) GetByProviderAndEmail(ctx context.Context, providerID uuid.UUID, email string) (*domain.Client, error) {
	var clientModel models.Client

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("business_id = ? AND email = ?", providerID, email).
		First(&clientModel).Error

	if err != nil {
		return nil, r.HandleNotFound(err, "client with provider ID "+providerID.String()+" and email "+email, uuid.Nil)
	}

	return mapClientModelToDomain(&clientModel), nil
}

// Update updates a client
func (r *ClientRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateClientInput, updatedBy uuid.UUID) error {
	// First find the client to ensure it exists
	var clientModel models.Client
	err := r.WithContext(ctx).First(&clientModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "client", id)
	}

	// Apply updates from the input
	updates := map[string]interface{}{}

	if input.FirstName != nil {
		updates["first_name"] = *input.FirstName
	}

	if input.LastName != nil {
		updates["last_name"] = *input.LastName
	}

	if input.Email != nil {
		updates["email"] = *input.Email
	}

	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}

	if input.Notes != nil {
		updates["notes"] = *input.Notes
	}

	// Perform the update with audit
	err = r.UpdateWithAudit(ctx, &clientModel, updates, updatedBy)
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}

	return nil
}

// Delete soft deletes a client
func (r *ClientRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// First find the client to ensure it exists
	var clientModel models.Client
	err := r.WithContext(ctx).First(&clientModel, "id = ?", id).Error
	if err != nil {
		return r.HandleNotFound(err, "client", id)
	}

	// Perform soft delete
	err = r.SoftDeleteWithAudit(ctx, &clientModel, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	return nil
}

// Search searches for clients by query within a provider
func (r *ClientRepository) Search(ctx context.Context, providerID uuid.UUID, query string, page, pageSize int) ([]*domain.Client, error) {
	var clientModels []models.Client

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("business_id = ?", providerID).
		Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR CONCAT(first_name, ' ', last_name) ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&clientModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search clients: %w", err)
	}

	return mapClientModelsToDomainSlice(clientModels), nil
}

// ListByProvider retrieves a paginated list of clients by provider ID
func (r *ClientRepository) ListByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) ([]*domain.Client, error) {
	var clientModels []models.Client

	offset := r.CalculateOffset(page, pageSize)

	err := r.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Where("business_id = ?", providerID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&clientModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list clients by provider: %w", err)
	}

	return mapClientModelsToDomainSlice(clientModels), nil
}

// Count counts all clients
func (r *ClientRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.WithContext(ctx).Model(&models.Client{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count clients: %w", err)
	}

	return count, nil
}

// CountByProvider counts clients by provider ID
func (r *ClientRepository) CountByProvider(ctx context.Context, providerID uuid.UUID) (int64, error) {
	var count int64

	err := r.WithContext(ctx).Model(&models.Client{}).Where("business_id = ?", providerID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count clients by provider: %w", err)
	}

	return count, nil
}

// Helper functions to map between domain entities and models

// mapClientDomainToModel converts a domain Client entity to a model Client
func mapClientDomainToModel(c *domain.Client) *models.Client {
	if c == nil {
		return nil
	}

	clientModel := &models.Client{
		BaseModel: models.BaseModel{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
		},
		BusinessID: c.ProviderID,
		UserID:     c.UserID,
		FirstName:  c.FirstName,
		LastName:   c.LastName,
		Email:      c.Email,
		Phone:      c.Phone,
		Notes:      c.Notes,
	}

	// Handle optional/pointer fields
	if c.CreatedBy != nil {
		clientModel.CreatedBy = c.CreatedBy
	}

	if c.UpdatedAt != nil {
		clientModel.UpdatedAt = *c.UpdatedAt
	}

	if c.UpdatedBy != nil {
		clientModel.UpdatedBy = c.UpdatedBy
	}

	if c.DeletedAt != nil {
		clientModel.DeletedAt.Time = *c.DeletedAt
		clientModel.DeletedAt.Valid = true
	}

	if c.DeletedBy != nil {
		clientModel.DeletedBy = c.DeletedBy
	}

	return clientModel
}

// mapClientModelToDomain converts a model Client to a domain Client entity
func mapClientModelToDomain(c *models.Client) *domain.Client {
	if c == nil {
		return nil
	}

	client := &domain.Client{
		ID:         c.ID,
		ProviderID: c.BusinessID,
		UserID:     c.UserID,
		FirstName:  c.FirstName,
		LastName:   c.LastName,
		Email:      c.Email,
		Phone:      c.Phone,
		Notes:      c.Notes,
		CreatedAt:  c.CreatedAt,
	}

	// Handle optional/pointer fields
	if c.CreatedBy != nil {
		client.CreatedBy = c.CreatedBy
	}

	if !c.UpdatedAt.IsZero() {
		client.UpdatedAt = &c.UpdatedAt
	}

	if c.UpdatedBy != nil {
		client.UpdatedBy = c.UpdatedBy
	}

	if c.DeletedAt.Valid {
		deletedAt := c.DeletedAt.Time
		client.DeletedAt = &deletedAt
	}

	if c.DeletedBy != nil {
		client.DeletedBy = c.DeletedBy
	}

	// Map related entities if loaded
	if c.User != nil && c.User.ID != uuid.Nil {
		client.User = mapUserModelToDomain(c.User)
	}

	if c.Business.ID != uuid.Nil {
		client.Provider = mapProviderModelToDomain(&c.Business)
	}

	return client
}

// mapClientModelsToDomainSlice converts a slice of model Client to a slice of domain Client entities
func mapClientModelsToDomainSlice(clientModels []models.Client) []*domain.Client {
	result := make([]*domain.Client, len(clientModels))
	for i, model := range clientModels {
		modelCopy := model // create a copy to avoid pointer issues
		result[i] = mapClientModelToDomain(&modelCopy)
	}
	return result
}
