package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// LoyaltyProgram represents a loyalty program defined by a business
type LoyaltyProgram struct {
	ID          uuid.UUID  `json:"id"`
	BusinessID  uuid.UUID  `json:"business_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ProgramType string     `json:"program_type"` // visit, spend, service, tier
	Rules       Rules      `json:"rules"`        // Configuration rules based on program type
	RewardType  string     `json:"reward_type"`  // percentage, fixed, free_service, upgrade, product
	RewardValue RewardInfo `json:"reward_value"` // Details of the reward
	IsActive    bool       `json:"is_active"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	DeletedBy   *uuid.UUID `json:"deleted_by,omitempty"`

	// Expanded relationships (populated by service when needed)
	Business *Business `json:"business,omitempty"`
}

// Rules represents the configuration rules for a loyalty program
type Rules struct {
	PointsPerSpend      float64     `json:"points_per_spend,omitempty"`      // For spend-based programs
	PointsPerVisit      int         `json:"points_per_visit,omitempty"`      // For visit-based programs
	RequiredVisits      int         `json:"required_visits,omitempty"`       // For visit-based programs
	RequiredSpend       float64     `json:"required_spend,omitempty"`        // For spend-based programs
	ApplicableServices  []uuid.UUID `json:"applicable_services,omitempty"`   // For service-specific programs
	MinimumSpend        float64     `json:"minimum_spend,omitempty"`         // Minimum spend per transaction
	PointsExpiry        int         `json:"points_expiry,omitempty"`         // Days until points expire (0 = no expiry)
	AllowCombineOffers  bool        `json:"allow_combine_offers"`            // Can combine with other offers
	TierThresholds      []int       `json:"tier_thresholds,omitempty"`       // For tiered programs
	BirthdayBonus       int         `json:"birthday_bonus,omitempty"`        // Birthday bonus points
	ReferralBonus       int         `json:"referral_bonus,omitempty"`        // Referral bonus points
	EnrollmentBonus     int         `json:"enrollment_bonus,omitempty"`      // Welcome bonus points
}

// RewardInfo represents the details of a reward
type RewardInfo struct {
	DiscountPercentage float64 `json:"discount_percentage,omitempty"` // For percentage discounts
	DiscountAmount     float64 `json:"discount_amount,omitempty"`     // For fixed amount discounts
	FreeServiceID      *uuid.UUID `json:"free_service_id,omitempty"`  // For free service rewards
	UpgradeServiceID   *uuid.UUID `json:"upgrade_service_id,omitempty"` // For upgrade rewards
	ProductID          *uuid.UUID `json:"product_id,omitempty"`       // For product rewards
	Description        string  `json:"description,omitempty"`        // Human readable description
}

// ClientLoyaltyMembership represents a client's membership in a loyalty program
type ClientLoyaltyMembership struct {
	ID           uuid.UUID  `json:"id"`
	ProgramID    uuid.UUID  `json:"program_id"`
	ClientID     uuid.UUID  `json:"client_id"`
	CurrentPoints int        `json:"current_points"`
	VisitsCount   int        `json:"visits_count"`
	TotalSpent    float64    `json:"total_spent"`
	TierLevel     string     `json:"tier_level,omitempty"`
	Progress      Progress   `json:"progress,omitempty"` // Progress data specific to program type
	JoinDate      time.Time  `json:"join_date"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`

	// Expanded relationships (populated by service when needed)
	Program *LoyaltyProgram `json:"program,omitempty"`
	Client  *Client         `json:"client,omitempty"`
}

// Progress represents progress data specific to program type
type Progress struct {
	CurrentTierProgress  int     `json:"current_tier_progress,omitempty"`  // For tiered programs
	NextTierThreshold    int     `json:"next_tier_threshold,omitempty"`    // For tiered programs
	PointsToNextReward   int     `json:"points_to_next_reward,omitempty"`  // Points needed for next reward
	VisitsToNextReward   int     `json:"visits_to_next_reward,omitempty"`  // Visits needed for next reward
	SpendToNextReward    float64 `json:"spend_to_next_reward,omitempty"`   // Spend needed for next reward
	LastRewardEarnedAt   *time.Time `json:"last_reward_earned_at,omitempty"` // When last reward was earned
	CompletionPercentage float64 `json:"completion_percentage,omitempty"`  // % towards next reward
}

// LoyaltyTransaction represents a loyalty point/visit transaction
type LoyaltyTransaction struct {
	ID             uuid.UUID  `json:"id"`
	MembershipID   uuid.UUID  `json:"membership_id"`
	AppointmentID  *uuid.UUID `json:"appointment_id,omitempty"`
	TransactionType string    `json:"transaction_type"` // earn, redeem, adjust, expire
	Points         int        `json:"points"`            // Positive for earn, negative for redeem
	Description    string     `json:"description"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      uuid.UUID  `json:"created_by"`

	// Expanded relationships (populated by service when needed)
	Membership  *ClientLoyaltyMembership `json:"membership,omitempty"`
	Appointment *Appointment             `json:"appointment,omitempty"`
}

// Input types for creating and updating loyalty programs

// CreateLoyaltyProgramInput is the input for creating a loyalty program
type CreateLoyaltyProgramInput struct {
	BusinessID  uuid.UUID  `json:"business_id" validate:"required"`
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description"`
	ProgramType string     `json:"program_type" validate:"required,oneof=visit spend service tier"`
	Rules       Rules      `json:"rules" validate:"required"`
	RewardType  string     `json:"reward_type" validate:"required,oneof=percentage fixed free_service upgrade product"`
	RewardValue RewardInfo `json:"reward_value" validate:"required"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

// UpdateLoyaltyProgramInput is the input for updating a loyalty program
type UpdateLoyaltyProgramInput struct {
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	ProgramType *string     `json:"program_type" validate:"omitempty,oneof=visit spend service tier"`
	Rules       *Rules      `json:"rules"`
	RewardType  *string     `json:"reward_type" validate:"omitempty,oneof=percentage fixed free_service upgrade product"`
	RewardValue *RewardInfo `json:"reward_value"`
	IsActive    *bool       `json:"is_active"`
	StartDate   *time.Time  `json:"start_date"`
	EndDate     *time.Time  `json:"end_date"`
}

// CreateClientLoyaltyMembershipInput is the input for creating a client loyalty membership
type CreateClientLoyaltyMembershipInput struct {
	ProgramID uuid.UUID `json:"program_id" validate:"required"`
	ClientID  uuid.UUID `json:"client_id" validate:"required"`
}

// UpdateClientLoyaltyMembershipInput is the input for updating a client loyalty membership
type UpdateClientLoyaltyMembershipInput struct {
	CurrentPoints *int       `json:"current_points"`
	VisitsCount   *int       `json:"visits_count"`
	TotalSpent    *float64   `json:"total_spent"`
	TierLevel     *string    `json:"tier_level"`
	Progress      *Progress  `json:"progress"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	IsActive      *bool      `json:"is_active"`
}

// CreateLoyaltyTransactionInput is the input for creating a loyalty transaction
type CreateLoyaltyTransactionInput struct {
	MembershipID    uuid.UUID  `json:"membership_id" validate:"required"`
	AppointmentID   *uuid.UUID `json:"appointment_id"`
	TransactionType string     `json:"transaction_type" validate:"required,oneof=earn redeem adjust expire"`
	Points          int        `json:"points" validate:"required"`
	Description     string     `json:"description"`
}

// Repository interfaces

// LoyaltyProgramRepository defines methods for loyalty program data store
type LoyaltyProgramRepository interface {
	Create(ctx context.Context, program *LoyaltyProgram) error
	GetByID(ctx context.Context, id uuid.UUID) (*LoyaltyProgram, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateLoyaltyProgramInput, updatedBy uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*LoyaltyProgram, error)
	ListActiveBybusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*LoyaltyProgram, error)
	Count(ctx context.Context) (int64, error)
	CountByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
	CountActiveByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// ClientLoyaltyMembershipRepository defines methods for client loyalty membership data store
type ClientLoyaltyMembershipRepository interface {
	Create(ctx context.Context, membership *ClientLoyaltyMembership) error
	GetByID(ctx context.Context, id uuid.UUID) (*ClientLoyaltyMembership, error)
	GetByClientAndProgram(ctx context.Context, clientID, programID uuid.UUID) (*ClientLoyaltyMembership, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateClientLoyaltyMembershipInput) error
	AddPoints(ctx context.Context, membershipID uuid.UUID, points int, description string, createdBy uuid.UUID) error
	RedeemPoints(ctx context.Context, membershipID uuid.UUID, points int, description string, createdBy uuid.UUID) error
	ListByClient(ctx context.Context, clientID uuid.UUID) ([]*ClientLoyaltyMembership, error)
	ListByProgram(ctx context.Context, programID uuid.UUID, page, pageSize int) ([]*ClientLoyaltyMembership, error)
	Count(ctx context.Context) (int64, error)
	CountByProgram(ctx context.Context, programID uuid.UUID) (int64, error)
}

// LoyaltyTransactionRepository defines methods for loyalty transaction data store
type LoyaltyTransactionRepository interface {
	Create(ctx context.Context, transaction *LoyaltyTransaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*LoyaltyTransaction, error)
	ListByMembership(ctx context.Context, membershipID uuid.UUID, page, pageSize int) ([]*LoyaltyTransaction, error)
	ListByClient(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*LoyaltyTransaction, error)
	Count(ctx context.Context) (int64, error)
	CountByMembership(ctx context.Context, membershipID uuid.UUID) (int64, error)
}

// Service interfaces

// LoyaltyProgramService defines business logic for loyalty program operations
type LoyaltyProgramService interface {
	CreateLoyaltyProgram(ctx context.Context, input *CreateLoyaltyProgramInput) (*LoyaltyProgram, error)
	GetLoyaltyProgram(ctx context.Context, id uuid.UUID) (*LoyaltyProgram, error)
	UpdateLoyaltyProgram(ctx context.Context, id uuid.UUID, input *UpdateLoyaltyProgramInput, updatedBy uuid.UUID) error
	DeleteLoyaltyProgram(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	ListLoyaltyProgramsByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*LoyaltyProgram, error)
	ListActiveLoyaltyProgramsByBusiness(ctx context.Context, businessID uuid.UUID, page, pageSize int) ([]*LoyaltyProgram, error)
	CountLoyaltyPrograms(ctx context.Context) (int64, error)
	CountLoyaltyProgramsByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
	CountActiveLoyaltyProgramsByBusiness(ctx context.Context, businessID uuid.UUID) (int64, error)
}

// ClientLoyaltyService defines business logic for client loyalty operations
type ClientLoyaltyService interface {
	EnrollClient(ctx context.Context, input *CreateClientLoyaltyMembershipInput) (*ClientLoyaltyMembership, error)
	GetClientMembership(ctx context.Context, clientID, programID uuid.UUID) (*ClientLoyaltyMembership, error)
	UpdateMembership(ctx context.Context, id uuid.UUID, input *UpdateClientLoyaltyMembershipInput) error
	EarnPoints(ctx context.Context, membershipID uuid.UUID, points int, appointmentID *uuid.UUID, description string, createdBy uuid.UUID) error
	RedeemReward(ctx context.Context, membershipID uuid.UUID, points int, description string, createdBy uuid.UUID) error
	GetMembershipProgress(ctx context.Context, membershipID uuid.UUID) (*Progress, error)
	ListClientMemberships(ctx context.Context, clientID uuid.UUID) ([]*ClientLoyaltyMembership, error)
	ListProgramMembers(ctx context.Context, programID uuid.UUID, page, pageSize int) ([]*ClientLoyaltyMembership, error)
	GetMembershipTransactions(ctx context.Context, membershipID uuid.UUID, page, pageSize int) ([]*LoyaltyTransaction, error)
}

// JSON marshaling helpers for JSONB fields

func (r Rules) Value() (interface{}, error) {
	return json.Marshal(r)
}

func (r *Rules) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, r)
}

func (ri RewardInfo) Value() (interface{}, error) {
	return json.Marshal(ri)
}

func (ri *RewardInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, ri)
}

func (p Progress) Value() (interface{}, error) {
	return json.Marshal(p)
}

func (p *Progress) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, p)
}