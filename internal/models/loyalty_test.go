package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLoyaltyProgramModel(t *testing.T) {
	// Test LoyaltyProgram Struct
	rules := LoyaltyProgramRules{
		ApplicableServices:     []uuid.UUID{uuid.New(), uuid.New()},
		MinimumSpend:           50.00,
		PointsRoundingMethod:   "up",
		AllowPointsExpiry:      true,
		PointsExpiryMonths:     12,
		AllowCombineWithOffers: false,
		RequireOptIn:           true,
	}

	tiers := TierDefinitions{
		{
			Name:             "Bronze",
			Level:            1,
			RequiredPoints:   0,
			PointsMultiplier: 1.0,
			Benefits:         []string{"Standard Earning Rate"},
		},
		{
			Name:             "Silver",
			Level:            2,
			RequiredPoints:   500,
			PointsMultiplier: 1.25,
			Benefits:         []string{"25% Bonus Points", "Birthday Gift"},
		},
		{
			Name:             "Gold",
			Level:            3,
			RequiredPoints:   1000,
			PointsMultiplier: 1.5,
			Benefits:         []string{"50% Bonus Points", "Birthday Gift", "Priority Booking"},
		},
	}

	program := LoyaltyProgram{
		BusinessID:      uuid.New(),
		Name:            "Beauty Rewards",
		Description:     "Earn points with every visit and service",
		ProgramType:     LoyaltyProgramTypePoints,
		PointsPerSpend:  10.00,
		PointsPerVisit:  50,
		EnrollmentBonus: 100,
		BirthdayBonus:   250,
		ReferralBonus:   150,
		ExpiryDays:      365,
		IsActive:        true,
		Rules:           rules,
		TierDefinitions: tiers,
	}

	assert.Equal(t, "Beauty Rewards", program.Name)
	assert.Equal(t, LoyaltyProgramTypePoints, program.ProgramType)
	assert.Equal(t, 10.00, program.PointsPerSpend)
	assert.Equal(t, 50, program.PointsPerVisit)
	assert.Equal(t, 365, program.ExpiryDays)
	assert.Equal(t, 2, len(program.Rules.ApplicableServices))
	assert.Equal(t, 3, len(program.TierDefinitions))
	assert.Equal(t, "Bronze", program.TierDefinitions[0].Name)
	assert.Equal(t, 3, program.TierDefinitions[2].Level)
}

func TestLoyaltyRewardModel(t *testing.T) {
	// Test LoyaltyReward struct
	reward := LoyaltyReward{
		BusinessID:     uuid.New(),
		ProgramID:      uuid.New(),
		Name:           "Free Haircut",
		Description:    "Redeem points for a free haircut",
		RewardType:     RewardTypeService,
		PointsRequired: 1000,
		ServiceID:      func() *uuid.UUID { id := uuid.New(); return &id }(),
		IsActive:       true,
		MinTierLevel:   2,
		MaxRedemptions: 1,
	}

	assert.Equal(t, "Free Haircut", reward.Name)
	assert.Equal(t, RewardTypeService, reward.RewardType)
	assert.Equal(t, 1000, reward.PointsRequired)
	assert.Equal(t, 2, reward.MinTierLevel)
	assert.Equal(t, 1, reward.MaxRedemptions)
	assert.NotNil(t, reward.ServiceID)
}

func TestClientLoyaltyModel(t *testing.T) {
	// Test ClientLoyalty struct
	now := time.Now()
	clientLoyalty := ClientLoyalty{
		BusinessID:       uuid.New(),
		ClientID:         uuid.New(),
		ProgramID:        uuid.New(),
		Points:           450,
		Visits:           8,
		TotalSpend:       560.75,
		CurrentTier:      2,
		EnrollmentDate:   now.AddDate(0, -3, 0), // 3 months ago
		LastActivityDate: &now,
		IsActive:         true,
		CardNumber:       "LYT123456789",
		MembershipStatus: "active",
	}

	assert.Equal(t, 450, clientLoyalty.Points)
	assert.Equal(t, 8, clientLoyalty.Visits)
	assert.Equal(t, 560.75, clientLoyalty.TotalSpend)
	assert.Equal(t, 2, clientLoyalty.CurrentTier)
	assert.Equal(t, "active", clientLoyalty.MembershipStatus)
	assert.Equal(t, now.Day(), clientLoyalty.LastActivityDate.Day())
	assert.Equal(t, "LYT123456789", clientLoyalty.CardNumber)
}

func TestLoyaltyTransactionModel(t *testing.T) {
	// Test LoyaltyTransaction struct
	transaction := LoyaltyTransaction{
		BusinessID:      uuid.New(),
		ClientID:        uuid.New(),
		ProgramID:       uuid.New(),
		ClientLoyaltyID: uuid.New(),
		PointsEarned:    85,
		VisitCounted:    true,
		Amount:          85.00,
		TransactionType: "earn",
		Description:     "Points earned for haircut and coloring",
	}

	assert.Equal(t, 85, transaction.PointsEarned)
	assert.Equal(t, true, transaction.VisitCounted)
	assert.Equal(t, 85.00, transaction.Amount)
	assert.Equal(t, "earn", transaction.TransactionType)
}

func TestRewardRedemptionModel(t *testing.T) {
	// Test RewardRedemption struct
	redemptionDate := time.Now()
	redemption := RewardRedemption{
		BusinessID:       uuid.New(),
		ClientID:         uuid.New(),
		ProgramID:        uuid.New(),
		RewardID:         uuid.New(),
		TransactionID:    uuid.New(),
		PointsRedeemed:   1000,
		RedemptionStatus: "redeemed",
		RedemptionDate:   &redemptionDate,
		RedemptionCode:   "RED123456789",
		IsDigital:        true,
	}

	assert.Equal(t, 1000, redemption.PointsRedeemed)
	assert.Equal(t, "redeemed", redemption.RedemptionStatus)
	assert.Equal(t, "RED123456789", redemption.RedemptionCode)
	assert.Equal(t, true, redemption.IsDigital)
	assert.Equal(t, redemptionDate.Day(), redemption.RedemptionDate.Day())
}

func TestLoyaltyProgramRulesSerialization(t *testing.T) {
	// Test JSON serialization and deserialization for LoyaltyProgramRules
	rules := LoyaltyProgramRules{
		ApplicableServices:     []uuid.UUID{uuid.New(), uuid.New()},
		ExcludedServices:       []uuid.UUID{uuid.New()},
		MinimumSpend:           25.00,
		PointsRoundingMethod:   "nearest",
		AllowPointsExpiry:      true,
		PointsExpiryMonths:     6,
		AllowCombineWithOffers: true,
		BlackoutDates:          []string{"2023-12-24", "2023-12-25", "2023-12-31"},
		RequireOptIn:           false,
	}

	// Serialize to JSON
	jsonBytes, err := json.Marshal(rules)
	assert.NoError(t, err)

	// Deserialize back
	var deserialized LoyaltyProgramRules
	err = json.Unmarshal(jsonBytes, &deserialized)
	assert.NoError(t, err)

	assert.Equal(t, len(rules.ApplicableServices), len(deserialized.ApplicableServices))
	assert.Equal(t, len(rules.ExcludedServices), len(deserialized.ExcludedServices))
	assert.Equal(t, rules.MinimumSpend, deserialized.MinimumSpend)
	assert.Equal(t, rules.PointsRoundingMethod, deserialized.PointsRoundingMethod)
	assert.Equal(t, rules.AllowPointsExpiry, deserialized.AllowPointsExpiry)
	assert.Equal(t, rules.PointsExpiryMonths, deserialized.PointsExpiryMonths)
	assert.Equal(t, len(rules.BlackoutDates), len(deserialized.BlackoutDates))
}

func TestTierDefinitionsSerialization(t *testing.T) {
	// Test JSON serialization and deserialization for TierDefinitions
	tiers := TierDefinitions{
		{
			Name:             "Standard",
			Level:            1,
			RequiredPoints:   0,
			PointsMultiplier: 1.0,
			Benefits:         []string{"Basic Services"},
		},
		{
			Name:             "Premium",
			Level:            2,
			RequiredPoints:   500,
			RequiredSpend:    1000.00,
			PointsMultiplier: 1.5,
			Benefits:         []string{"Priority Booking", "Complimentary Drink"},
		},
	}

	// Serialize to JSON
	jsonBytes, err := json.Marshal(tiers)
	assert.NoError(t, err)

	// Deserialize back
	var deserialized TierDefinitions
	err = json.Unmarshal(jsonBytes, &deserialized)
	assert.NoError(t, err)

	assert.Equal(t, len(tiers), len(deserialized))
	assert.Equal(t, tiers[0].Name, deserialized[0].Name)
	assert.Equal(t, tiers[0].Level, deserialized[0].Level)
	assert.Equal(t, tiers[1].RequiredPoints, deserialized[1].RequiredPoints)
	assert.Equal(t, tiers[1].RequiredSpend, deserialized[1].RequiredSpend)
	assert.Equal(t, len(tiers[1].Benefits), len(deserialized[1].Benefits))
}
