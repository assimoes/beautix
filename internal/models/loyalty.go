package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LoyaltyProgramType represents the type of loyalty program
type LoyaltyProgramType string

const (
	// LoyaltyProgramTypePoints represents a points-based loyalty program
	LoyaltyProgramTypePoints LoyaltyProgramType = "points"
	// LoyaltyProgramTypeVisits represents a visits-based loyalty program
	LoyaltyProgramTypeVisits LoyaltyProgramType = "visits"
	// LoyaltyProgramTypeSpend represents a spend-based loyalty program
	LoyaltyProgramTypeSpend LoyaltyProgramType = "spend"
	// LoyaltyProgramTypeTiered represents a tiered loyalty program
	LoyaltyProgramTypeTiered LoyaltyProgramType = "tiered"
)

// RewardType represents the type of reward
type RewardType string

const (
	// RewardTypeDiscount represents a discount reward
	RewardTypeDiscount RewardType = "discount"
	// RewardTypeService represents a free service reward
	RewardTypeService RewardType = "service"
	// RewardTypeProduct represents a free product reward
	RewardTypeProduct RewardType = "product"
	// RewardTypeGiftCard represents a gift card reward
	RewardTypeGiftCard RewardType = "gift_card"
	// RewardTypeUpgrade represents a service upgrade reward
	RewardTypeUpgrade RewardType = "upgrade"
)

// LoyaltyProgram represents a loyalty program for a business
type LoyaltyProgram struct {
	BaseModel
	BusinessID      uuid.UUID           `gorm:"type:uuid;not null;index" json:"business_id"`
	Business        Business            `gorm:"foreignKey:BusinessID" json:"business"`
	Name            string              `gorm:"not null" json:"name"`
	Description     string              `json:"description"`
	ProgramType     LoyaltyProgramType  `gorm:"type:text;not null" json:"program_type"`
	PointsPerSpend  float64             `gorm:"type:decimal(10,2)" json:"points_per_spend"`
	PointsPerVisit  int                 `json:"points_per_visit"`
	VisitsRequired  int                 `json:"visits_required"`
	SpendRequired   float64             `gorm:"type:decimal(10,2)" json:"spend_required"`
	EnrollmentBonus int                 `json:"enrollment_bonus"`
	BirthdayBonus   int                 `json:"birthday_bonus"`
	ReferralBonus   int                 `json:"referral_bonus"`
	ExpiryDays      int                 `json:"expiry_days"` // 0 means no expiry
	IsActive        bool                `gorm:"not null;default:true" json:"is_active"`
	StartDate       *time.Time          `json:"start_date"`
	EndDate         *time.Time          `json:"end_date"`
	Rules           LoyaltyProgramRules `gorm:"type:jsonb" json:"rules"`
	TierDefinitions TierDefinitions     `gorm:"type:jsonb" json:"tier_definitions"`
}

// TableName overrides the table name
func (LoyaltyProgram) TableName() string {
	return "loyalty_programs"
}

// LoyaltyProgramRules stores the rules for a loyalty program
type LoyaltyProgramRules struct {
	ApplicableServices     []uuid.UUID `json:"applicable_services,omitempty"`
	ExcludedServices       []uuid.UUID `json:"excluded_services,omitempty"`
	ApplicableProducts     []uuid.UUID `json:"applicable_products,omitempty"`
	ExcludedProducts       []uuid.UUID `json:"excluded_products,omitempty"`
	MinimumSpend           float64     `json:"minimum_spend,omitempty"`
	PointsRoundingMethod   string      `json:"points_rounding_method,omitempty"` // "up", "down", "nearest"
	AllowPointsExpiry      bool        `json:"allow_points_expiry"`
	PointsExpiryMonths     int         `json:"points_expiry_months,omitempty"`
	AllowCombineWithOffers bool        `json:"allow_combine_with_offers"`
	BlackoutDates          []string    `json:"blackout_dates,omitempty"`
	RequireOptIn           bool        `json:"require_opt_in"`
}

// Scan implements the sql.Scanner interface for LoyaltyProgramRules
func (lpr *LoyaltyProgramRules) Scan(value interface{}) error {
	if value == nil {
		*lpr = LoyaltyProgramRules{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp LoyaltyProgramRules
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*lpr = temp

	return nil
}

// Value implements the driver.Valuer interface for LoyaltyProgramRules
func (lpr LoyaltyProgramRules) Value() (driver.Value, error) {
	return json.Marshal(lpr)
}

// TierDefinition represents a single tier in a tiered loyalty program
type TierDefinition struct {
	Name             string   `json:"name"`
	Level            int      `json:"level"`
	RequiredPoints   int      `json:"required_points,omitempty"`
	RequiredSpend    float64  `json:"required_spend,omitempty"`
	RequiredVisits   int      `json:"required_visits,omitempty"`
	PointsMultiplier float64  `json:"points_multiplier,omitempty"`
	Benefits         []string `json:"benefits,omitempty"`
}

// TierDefinitions is a slice of TierDefinition that can be stored as JSONB
type TierDefinitions []TierDefinition

// Scan implements the sql.Scanner interface for TierDefinitions
func (td *TierDefinitions) Scan(value interface{}) error {
	if value == nil {
		*td = TierDefinitions{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var temp TierDefinitions
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*td = temp

	return nil
}

// Value implements the driver.Valuer interface for TierDefinitions
func (td TierDefinitions) Value() (driver.Value, error) {
	return json.Marshal(td)
}

// LoyaltyReward represents a reward that can be redeemed with loyalty points or visits
type LoyaltyReward struct {
	BaseModel
	BusinessID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"business_id"`
	Business        Business       `gorm:"foreignKey:BusinessID" json:"business"`
	ProgramID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"program_id"`
	LoyaltyProgram  LoyaltyProgram `gorm:"foreignKey:ProgramID" json:"loyalty_program"`
	Name            string         `gorm:"not null" json:"name"`
	Description     string         `json:"description"`
	RewardType      RewardType     `gorm:"type:text;not null" json:"reward_type"`
	PointsRequired  int            `json:"points_required"`
	VisitsRequired  int            `json:"visits_required"`
	DiscountAmount  float64        `gorm:"type:decimal(10,2)" json:"discount_amount"`
	DiscountPercent int            `json:"discount_percent"`
	ServiceID       *uuid.UUID     `gorm:"type:uuid;index" json:"service_id"`
	ProductID       *uuid.UUID     `gorm:"type:uuid;index" json:"product_id"`
	GiftCardAmount  float64        `gorm:"type:decimal(10,2)" json:"gift_card_amount"`
	IsActive        bool           `gorm:"not null;default:true" json:"is_active"`
	StartDate       *time.Time     `json:"start_date"`
	EndDate         *time.Time     `json:"end_date"`
	MinTierLevel    int            `gorm:"not null;default:0" json:"min_tier_level"`  // 0 means available to all tiers
	MaxRedemptions  int            `gorm:"not null;default:0" json:"max_redemptions"` // 0 means unlimited
	ImageURL        string         `json:"image_url"`
}

// TableName overrides the table name
func (LoyaltyReward) TableName() string {
	return "loyalty_rewards"
}

// ClientLoyalty represents a client's loyalty program membership
type ClientLoyalty struct {
	BaseModel
	BusinessID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"business_id"`
	Business         Business       `gorm:"foreignKey:BusinessID" json:"business"`
	ClientID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"client_id"`
	Client           Client         `gorm:"foreignKey:ClientID" json:"client"`
	ProgramID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"program_id"`
	LoyaltyProgram   LoyaltyProgram `gorm:"foreignKey:ProgramID" json:"loyalty_program"`
	Points           int            `gorm:"not null;default:0" json:"points"`
	Visits           int            `gorm:"not null;default:0" json:"visits"`
	TotalSpend       float64        `gorm:"type:decimal(10,2);not null;default:0" json:"total_spend"`
	CurrentTier      int            `gorm:"not null;default:0" json:"current_tier"`
	EnrollmentDate   time.Time      `gorm:"not null" json:"enrollment_date"`
	LastActivityDate *time.Time     `json:"last_activity_date"`
	IsActive         bool           `gorm:"not null;default:true" json:"is_active"`
	CardNumber       string         `json:"card_number"`
	ExpiryDate       *time.Time     `json:"expiry_date"`
	MembershipStatus string         `gorm:"not null;default:'active'" json:"membership_status"` // "active", "inactive", "suspended"
}

// TableName overrides the table name
func (ClientLoyalty) TableName() string {
	return "client_loyalty"
}

// LoyaltyTransaction represents a loyalty point/visit transaction
type LoyaltyTransaction struct {
	BaseModel
	BusinessID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"business_id"`
	Business        Business       `gorm:"foreignKey:BusinessID" json:"business"`
	ClientID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"client_id"`
	Client          Client         `gorm:"foreignKey:ClientID" json:"client"`
	ProgramID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"program_id"`
	LoyaltyProgram  LoyaltyProgram `gorm:"foreignKey:ProgramID" json:"loyalty_program"`
	ClientLoyaltyID uuid.UUID      `gorm:"type:uuid;not null;index" json:"client_loyalty_id"`
	ClientLoyalty   ClientLoyalty  `gorm:"foreignKey:ClientLoyaltyID" json:"client_loyalty"`
	AppointmentID   *uuid.UUID     `gorm:"type:uuid;index" json:"appointment_id"`
	Appointment     *Appointment   `gorm:"foreignKey:AppointmentID" json:"appointment"`
	PointsEarned    int            `json:"points_earned"`
	PointsRedeemed  int            `json:"points_redeemed"`
	VisitCounted    bool           `gorm:"not null;default:false" json:"visit_counted"`
	Amount          float64        `gorm:"type:decimal(10,2)" json:"amount"`
	TransactionType string         `gorm:"not null" json:"transaction_type"` // "earn", "redeem", "adjust", "expire", "bonus"
	Description     string         `json:"description"`
	RewardID        *uuid.UUID     `gorm:"type:uuid;index" json:"reward_id"`
	ExpiryDate      *time.Time     `json:"expiry_date"`
	ReferralSource  *uuid.UUID     `gorm:"type:uuid;index" json:"referral_source"` // ClientID of referring client
}

// TableName overrides the table name
func (LoyaltyTransaction) TableName() string {
	return "loyalty_transactions"
}

// RewardRedemption represents the redemption of a loyalty reward
type RewardRedemption struct {
	BaseModel
	BusinessID         uuid.UUID          `gorm:"type:uuid;not null;index" json:"business_id"`
	Business           Business           `gorm:"foreignKey:BusinessID" json:"business"`
	ClientID           uuid.UUID          `gorm:"type:uuid;not null;index" json:"client_id"`
	Client             Client             `gorm:"foreignKey:ClientID" json:"client"`
	ProgramID          uuid.UUID          `gorm:"type:uuid;not null;index" json:"program_id"`
	LoyaltyProgram     LoyaltyProgram     `gorm:"foreignKey:ProgramID" json:"loyalty_program"`
	RewardID           uuid.UUID          `gorm:"type:uuid;not null;index" json:"reward_id"`
	LoyaltyReward      LoyaltyReward      `gorm:"foreignKey:RewardID" json:"loyalty_reward"`
	TransactionID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"transaction_id"`
	LoyaltyTransaction LoyaltyTransaction `gorm:"foreignKey:TransactionID" json:"loyalty_transaction"`
	AppointmentID      *uuid.UUID         `gorm:"type:uuid;index" json:"appointment_id"`
	Appointment        *Appointment       `gorm:"foreignKey:AppointmentID" json:"appointment"`
	PointsRedeemed     int                `json:"points_redeemed"`
	RedemptionStatus   string             `gorm:"not null;default:'pending'" json:"redemption_status"` // "pending", "redeemed", "cancelled", "expired"
	RedemptionDate     *time.Time         `json:"redemption_date"`
	ExpiryDate         *time.Time         `json:"expiry_date"`
	CancellationReason string             `json:"cancellation_reason"`
	RedemptionCode     string             `gorm:"not null" json:"redemption_code"`
	IsDigital          bool               `gorm:"not null;default:false" json:"is_digital"`
}

// TableName overrides the table name
func (RewardRedemption) TableName() string {
	return "reward_redemptions"
}
