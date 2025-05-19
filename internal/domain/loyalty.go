package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// LoyaltyProgram represents a loyalty program defined by a provider
type LoyaltyProgram struct {
	ID          uuid.UUID  `json:"id"`
	ProviderID  uuid.UUID  `json:"provider_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ProgramType string     `json:"program_type"` // visit-based, spending-based, service-specific, tiered
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	DeletedBy   *uuid.UUID `json:"deleted_by,omitempty"`
	
	// Expanded relationships (populated by service when needed)
	Provider *Provider       `json:"provider,omitempty"`
	Rewards  []*LoyaltyReward `json:"rewards,omitempty"`
}

// LoyaltyReward represents a reward in a loyalty program
type LoyaltyReward struct {
	ID             uuid.UUID  `json:"id"`
	ProgramID      uuid.UUID  `json:"program_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	RequiredPoints int        `json:"required_points"`
	RewardType     string     `json:"reward_type"` // percentage, fixed_amount, free_service, service_upgrade, product
	RewardValue    string     `json:"reward_value"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	DeletedBy      *uuid.UUID `json:"deleted_by,omitempty"`
	
	// Expanded relationships (populated by service when needed)
	Program *LoyaltyProgram `json:"program,omitempty"`
}

// ClientLoyalty represents a client's points in a loyalty program
type ClientLoyalty struct {
	ID        uuid.UUID  `json:"id"`
	ClientID  uuid.UUID  `json:"client_id"`
	ProgramID uuid.UUID  `json:"program_id"`
	Points    int        `json:"points"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uuid.UUID `json:"deleted_by,omitempty"`
	
	// Expanded relationships (populated by service when needed)
	Client  *Client        `json:"client,omitempty"`
	Program *LoyaltyProgram `json:"program,omitempty"`
}

// CreateLoyaltyProgramInput is the input for creating a loyalty program
type CreateLoyaltyProgramInput struct {
	ProviderID  uuid.UUID `json:"provider_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	ProgramType string    `json:"program_type" validate:"required,oneof=visit-based spending-based service-specific tiered"`
}

// UpdateLoyaltyProgramInput is the input for updating a loyalty program
type UpdateLoyaltyProgramInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ProgramType *string `json:"program_type" validate:"omitempty,oneof=visit-based spending-based service-specific tiered"`
}

// CreateLoyaltyRewardInput is the input for creating a loyalty reward
type CreateLoyaltyRewardInput struct {
	ProgramID      uuid.UUID `json:"program_id" validate:"required"`
	Name           string    `json:"name" validate:"required"`
	Description    string    `json:"description"`
	RequiredPoints int       `json:"required_points" validate:"required,min=1"`
	RewardType     string    `json:"reward_type" validate:"required,oneof=percentage fixed_amount free_service service_upgrade product"`
	RewardValue    string    `json:"reward_value" validate:"required"`
}

// UpdateLoyaltyRewardInput is the input for updating a loyalty reward
type UpdateLoyaltyRewardInput struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	RequiredPoints *int    `json:"required_points" validate:"omitempty,min=1"`
	RewardType     *string `json:"reward_type" validate:"omitempty,oneof=percentage fixed_amount free_service service_upgrade product"`
	RewardValue    *string `json:"reward_value"`
}

// LoyaltyProgramRepository defines methods for loyalty program data store
type LoyaltyProgramRepository interface {
	Create(ctx context.Context, program *LoyaltyProgram) error
	GetByID(ctx context.Context, id uuid.UUID) (*LoyaltyProgram, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateLoyaltyProgramInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) ([]*LoyaltyProgram, error)
	Count(ctx context.Context) (int64, error)
	CountByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
}

// LoyaltyRewardRepository defines methods for loyalty reward data store
type LoyaltyRewardRepository interface {
	Create(ctx context.Context, reward *LoyaltyReward) error
	GetByID(ctx context.Context, id uuid.UUID) (*LoyaltyReward, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateLoyaltyRewardInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByProgram(ctx context.Context, programID uuid.UUID) ([]*LoyaltyReward, error)
	Count(ctx context.Context) (int64, error)
	CountByProgram(ctx context.Context, programID uuid.UUID) (int64, error)
}

// ClientLoyaltyRepository defines methods for client loyalty data store
type ClientLoyaltyRepository interface {
	Create(ctx context.Context, clientLoyalty *ClientLoyalty) error
	GetByClientAndProgram(ctx context.Context, clientID, programID uuid.UUID) (*ClientLoyalty, error)
	UpdatePoints(ctx context.Context, clientID, programID uuid.UUID, points int, updatedBy uuid.UUID) error
	ListByClient(ctx context.Context, clientID uuid.UUID) ([]*ClientLoyalty, error)
	ListByProgram(ctx context.Context, programID uuid.UUID, page, pageSize int) ([]*ClientLoyalty, error)
	Delete(ctx context.Context, clientID, programID uuid.UUID, deletedBy uuid.UUID) error
}

// LoyaltyProgramService defines business logic for loyalty program operations
type LoyaltyProgramService interface {
	CreateLoyaltyProgram(ctx context.Context, input *CreateLoyaltyProgramInput) (*LoyaltyProgram, error)
	GetLoyaltyProgram(ctx context.Context, id uuid.UUID) (*LoyaltyProgram, error)
	UpdateLoyaltyProgram(ctx context.Context, id uuid.UUID, input *UpdateLoyaltyProgramInput, updatedBy uuid.UUID) error
	DeleteLoyaltyProgram(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListLoyaltyProgramsByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) ([]*LoyaltyProgram, error)
	CountLoyaltyPrograms(ctx context.Context) (int64, error)
	CountLoyaltyProgramsByProvider(ctx context.Context, providerID uuid.UUID) (int64, error)
}

// LoyaltyRewardService defines business logic for loyalty reward operations
type LoyaltyRewardService interface {
	CreateLoyaltyReward(ctx context.Context, input *CreateLoyaltyRewardInput) (*LoyaltyReward, error)
	GetLoyaltyReward(ctx context.Context, id uuid.UUID) (*LoyaltyReward, error)
	UpdateLoyaltyReward(ctx context.Context, id uuid.UUID, input *UpdateLoyaltyRewardInput, updatedBy uuid.UUID) error
	DeleteLoyaltyReward(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListLoyaltyRewardsByProgram(ctx context.Context, programID uuid.UUID) ([]*LoyaltyReward, error)
	CountLoyaltyRewards(ctx context.Context) (int64, error)
	CountLoyaltyRewardsByProgram(ctx context.Context, programID uuid.UUID) (int64, error)
}

// ClientLoyaltyService defines business logic for client loyalty operations
type ClientLoyaltyService interface {
	CreateClientLoyalty(ctx context.Context, clientID, programID uuid.UUID, initialPoints int) (*ClientLoyalty, error)
	GetClientLoyalty(ctx context.Context, clientID, programID uuid.UUID) (*ClientLoyalty, error)
	UpdateClientLoyaltyPoints(ctx context.Context, clientID, programID uuid.UUID, points int, updatedBy uuid.UUID) error
	AddClientLoyaltyPoints(ctx context.Context, clientID, programID uuid.UUID, pointsToAdd int, updatedBy uuid.UUID) error
	RedeemReward(ctx context.Context, clientID, programID, rewardID uuid.UUID, updatedBy uuid.UUID) error
	ListClientLoyalties(ctx context.Context, clientID uuid.UUID) ([]*ClientLoyalty, error)
	ListClientsByProgram(ctx context.Context, programID uuid.UUID, page, pageSize int) ([]*ClientLoyalty, error)
	DeleteClientLoyalty(ctx context.Context, clientID, programID uuid.UUID, deletedBy uuid.UUID) error
}