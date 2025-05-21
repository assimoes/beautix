package test

import (
	"testing"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseSchema(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create a user with clerk_id
	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		ClerkID:   "test_clerk_id",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+1234567890",
		Role:      models.UserRoleProvider,
		IsActive:  true,
	}

	// Save the user
	err = testDB.Create(&user).Error
	assert.NoError(t, err, "Failed to create user")

	// Retrieve the user
	var savedUser models.User
	err = testDB.First(&savedUser, "clerk_id = ?", "test_clerk_id").Error
	assert.NoError(t, err, "Failed to find user by clerk_id")
	assert.Equal(t, userID, savedUser.ID)

	// Create a business
	businessID := uuid.New()
	business := models.Business{
		BaseModel: models.BaseModel{
			ID: businessID,
		},
		UserID:           userID,
		Name:             "test-business",
		DisplayName:      "Test Business",
		Description:      "A test business",
		Address:          "123 Test St",
		City:             "Test City",
		Country:          "Portugal",
		PostalCode:       "1000-100",
		Phone:            "+351123456789",
		Email:            "contact@testbusiness.com",
		SubscriptionTier: models.SubscriptionTierPro,
		IsActive:         true,
		SocialLinks: models.SocialLinks{
			Instagram: "testbusiness",
		},
		Settings: models.BusinessSettings{
			AllowOnlineBooking: true,
		},
	}

	// Save the business
	err = testDB.Create(&business).Error
	assert.NoError(t, err, "Failed to create business")

	// Retrieve the business
	var savedBusiness models.Business
	err = testDB.First(&savedBusiness, "id = ?", businessID).Error
	assert.NoError(t, err, "Failed to find business")
	assert.Equal(t, businessID, savedBusiness.ID)
	assert.Equal(t, "Test Business", savedBusiness.DisplayName)
	assert.Equal(t, "testbusiness", savedBusiness.SocialLinks.Instagram)
	assert.True(t, savedBusiness.Settings.AllowOnlineBooking)
}