package models_test

import (
	"testing"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserModel(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate the models
	err = testDB.AutoMigrate(&models.User{}, &models.UserConnectedAccount{})
	require.NoError(t, err, "Failed to migrate models")

	// Create a user
	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		ClerkID:   "clerk_123456789",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+1234567890",
		Role:      models.UserRoleUser,
		IsActive:  true,
	}

	// Save the user
	err = testDB.Create(&user).Error
	assert.NoError(t, err, "Failed to create user")

	// Verify user was created with ID
	var savedUser models.User
	err = testDB.First(&savedUser, "id = ?", userID).Error
	assert.NoError(t, err, "Failed to find user")
	assert.Equal(t, userID, savedUser.ID)
	assert.Equal(t, "clerk_123456789", savedUser.ClerkID)
	assert.Equal(t, "test@example.com", savedUser.Email)

	// Create a connected account
	connectedAccount := models.UserConnectedAccount{
		UserID:       userID,
		ProviderType: "google",
		ProviderID:   "google_12345",
		IsActive:     true,
	}

	// Save the connected account
	err = testDB.Create(&connectedAccount).Error
	assert.NoError(t, err, "Failed to create connected account")

	// Verify connected account was created with proper associations
	var savedConnectedAccounts []models.UserConnectedAccount
	err = testDB.Preload("User").Where("user_id = ?", userID).Find(&savedConnectedAccounts).Error
	assert.NoError(t, err, "Failed to find connected accounts")
	assert.Len(t, savedConnectedAccounts, 1, "Should have one connected account")
	assert.Equal(t, "google", savedConnectedAccounts[0].ProviderType)
	assert.Equal(t, "google_12345", savedConnectedAccounts[0].ProviderID)
	assert.Equal(t, userID, savedConnectedAccounts[0].User.ID)

	// Test unique constraints
	duplicateUser := models.User{
		ClerkID:   "clerk_123456789", // Duplicate ClerkID
		Email:     "another@example.com",
		FirstName: "Another",
		LastName:  "User",
		Role:      models.UserRoleUser,
	}
	err = testDB.Create(&duplicateUser).Error
	assert.Error(t, err, "Should not allow duplicate ClerkID")

	duplicateEmail := models.User{
		ClerkID:   "clerk_different",
		Email:     "test@example.com", // Duplicate email
		FirstName: "Another",
		LastName:  "User",
		Role:      models.UserRoleUser,
	}
	err = testDB.Create(&duplicateEmail).Error
	assert.Error(t, err, "Should not allow duplicate email")

	// Note: In a real application, indexes would be added through migrations
	// We're relying on GORM's built-in index creation for tests

	// Test soft delete
	err = testDB.Delete(&user).Error
	assert.NoError(t, err, "Failed to soft delete user")

	// Verify user is soft deleted
	var deletedUser models.User
	err = testDB.Unscoped().First(&deletedUser, "id = ?", userID).Error
	assert.NoError(t, err, "Failed to find soft deleted user")
	assert.False(t, deletedUser.DeletedAt.Time.IsZero(), "DeletedAt should be set")

	// Verify we can't find the user with normal queries
	err = testDB.First(&models.User{}, "id = ?", userID).Error
	assert.Error(t, err, "Should not find soft deleted user")
}