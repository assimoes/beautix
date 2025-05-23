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
	// Connect to the test database using simple approach
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	t.Run("User CRUD Operations", func(t *testing.T) {
		// Clean up all tables comprehensively to avoid foreign key issues
		database.CleanupAllTables(t, testDB.DB)

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
		err = testDB.DB.Create(&user).Error
		assert.NoError(t, err, "Failed to create user")

		// Verify user was created with ID
		var savedUser models.User
		err = testDB.DB.First(&savedUser, "id = ?", userID).Error
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
		err = testDB.DB.Create(&connectedAccount).Error
		assert.NoError(t, err, "Failed to create connected account")

		// Verify connected account was created with proper associations
		var savedConnectedAccounts []models.UserConnectedAccount
		err = testDB.DB.Preload("User").Where("user_id = ?", userID).Find(&savedConnectedAccounts).Error
		assert.NoError(t, err, "Failed to find connected accounts")
		assert.Len(t, savedConnectedAccounts, 1, "Should have one connected account")
		assert.Equal(t, "google", savedConnectedAccounts[0].ProviderType)
		assert.Equal(t, "google_12345", savedConnectedAccounts[0].ProviderID)
		assert.Equal(t, userID, savedConnectedAccounts[0].User.ID)

		// Test soft delete
		err = testDB.DB.Delete(&user).Error
		assert.NoError(t, err, "Failed to soft delete user")

		// Verify user is soft deleted
		var deletedUser models.User
		err = testDB.DB.Unscoped().First(&deletedUser, "id = ?", userID).Error
		assert.NoError(t, err, "Failed to find soft deleted user")
		assert.False(t, deletedUser.DeletedAt.Time.IsZero(), "DeletedAt should be set")

		// Verify we can't find the user with normal queries
		err = testDB.DB.First(&models.User{}, "id = ?", userID).Error
		assert.Error(t, err, "Should not find soft deleted user")
	})

	t.Run("Unique Constraints", func(t *testing.T) {
		// Test unique constraints in separate transactions to avoid transaction abort

		// Test duplicate ClerkID
		t.Run("Duplicate ClerkID", func(t *testing.T) {
			// Create first user
			user1 := models.User{
				ClerkID:   "clerk_unique_test",
				Email:     "first@example.com",
				FirstName: "First",
				LastName:  "User",
				Role:      models.UserRoleUser,
			}
			err = testDB.DB.Create(&user1).Error
			assert.NoError(t, err, "Failed to create first user")

			// Try to create user with duplicate ClerkID
			user2 := models.User{
				ClerkID:   "clerk_unique_test", // Duplicate ClerkID
				Email:     "second@example.com",
				FirstName: "Second",
				LastName:  "User",
				Role:      models.UserRoleUser,
			}
			err = testDB.DB.Create(&user2).Error
			assert.Error(t, err, "Should not allow duplicate ClerkID")
		})

		// Test duplicate Email
		t.Run("Duplicate Email", func(t *testing.T) {
			// Create first user
			user1 := models.User{
				ClerkID:   "clerk_email_test1",
				Email:     "duplicate@example.com",
				FirstName: "First",
				LastName:  "User",
				Role:      models.UserRoleUser,
			}
			err = testDB.DB.Create(&user1).Error
			assert.NoError(t, err, "Failed to create first user")

			// Try to create user with duplicate email
			user2 := models.User{
				ClerkID:   "clerk_email_test2",
				Email:     "duplicate@example.com", // Duplicate email
				FirstName: "Second",
				LastName:  "User",
				Role:      models.UserRoleUser,
			}
			err = testDB.DB.Create(&user2).Error
			assert.Error(t, err, "Should not allow duplicate email")
		})
	})
}
