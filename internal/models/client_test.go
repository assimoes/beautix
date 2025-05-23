package models_test

import (
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientModel(t *testing.T) {
	// Connect to the test database using simple approach
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Clean up all tables comprehensively to avoid foreign key issues
	database.CleanupAllTables(t, testDB.DB)

	// Create a user first
	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{
			ID:        userID,
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		ClerkID:   "clerk_client_" + userID.String()[:8], // Unique ClerkID
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      models.UserRoleOwner,
	}

	err = testDB.DB.Create(&user).Error
	assert.NoError(t, err, "Failed to create user")

	// Create a business first
	businessID := uuid.New()
	business := models.Business{
		BaseModel: models.BaseModel{
			ID:        businessID,
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		UserID:       userID,
		Name:         "Client Test Business",
		BusinessType: "salon",
		Phone:        "+351123456789",
		Email:        "business@example.com",
		Country:      "Portugal",
		TimeZone:     "Europe/Lisbon",
	}

	err = testDB.DB.Create(&business).Error
	assert.NoError(t, err, "Failed to create business")

	// Create a client
	clientID := uuid.New()
	dob := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	client := models.Client{
		BaseModel: models.BaseModel{
			ID:        clientID,
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		BusinessID:       businessID,
		UserID:           &userID,
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "john.doe@example.com",
		Phone:            "+351987654321",
		DateOfBirth:      &dob,
		AddressLine1:     "456 Client St",
		City:             "Lisbon",
		Country:          "Portugal",
		PostalCode:       "1000-200",
		Notes:            "Regular client with monthly appointments",
		Allergies:        "Latex, Nuts",
		HealthConditions: "None",
		ReferralSource:   "Website",
		AcceptsMarketing: true,
		IsActive:         true,
	}

	// Save the client
	err = testDB.DB.Create(&client).Error
	assert.NoError(t, err, "Failed to create client")

	// Verify client was created
	var savedClient models.Client
	err = testDB.DB.First(&savedClient, "id = ?", clientID).Error
	assert.NoError(t, err, "Failed to find client")
	assert.Equal(t, clientID, savedClient.ID)
	assert.Equal(t, "John", savedClient.FirstName)
	assert.Equal(t, "Doe", savedClient.LastName)
	assert.Equal(t, "john.doe@example.com", savedClient.Email)
	assert.Equal(t, "+351987654321", savedClient.Phone)
	assert.Equal(t, dob, *savedClient.DateOfBirth)
	assert.Equal(t, "456 Client St", savedClient.AddressLine1)
	assert.Equal(t, "Lisbon", savedClient.City)
	assert.Equal(t, "Portugal", savedClient.Country)
	assert.Equal(t, "1000-200", savedClient.PostalCode)
	assert.Equal(t, "Regular client with monthly appointments", savedClient.Notes)
	assert.Equal(t, "Latex, Nuts", savedClient.Allergies)
	assert.Equal(t, "None", savedClient.HealthConditions)
	assert.Equal(t, "Website", savedClient.ReferralSource)
	assert.True(t, savedClient.AcceptsMarketing)
	assert.True(t, savedClient.IsActive)

	// Test loaded relationships
	err = testDB.DB.Preload("User").Preload("Business").First(&savedClient, "id = ?", clientID).Error
	assert.NoError(t, err, "Failed to find client with relationships")
	assert.NotNil(t, savedClient.UserID, "UserID should not be nil")
	if savedClient.UserID != nil {
		assert.Equal(t, userID, *savedClient.UserID)
	}
	assert.NotNil(t, savedClient.User, "User relationship should be loaded")
	if savedClient.User != nil {
		assert.Equal(t, "Test", savedClient.User.FirstName)
	}
	assert.Equal(t, businessID, savedClient.BusinessID)
	assert.Equal(t, "Client Test Business", savedClient.Business.Name)

	// Test update client
	savedClient.Phone = "+351999888777"
	savedClient.UpdatedBy = &userID
	err = testDB.DB.Save(&savedClient).Error
	assert.NoError(t, err, "Failed to update client")

	// Verify update
	var updatedClient models.Client
	err = testDB.DB.First(&updatedClient, "id = ?", clientID).Error
	assert.NoError(t, err, "Failed to find updated client")
	assert.Equal(t, "+351999888777", updatedClient.Phone)

	// Test soft delete
	err = testDB.DB.Delete(&savedClient).Error
	assert.NoError(t, err, "Failed to soft delete client")

	// Verify soft delete
	var deletedClient models.Client
	err = testDB.DB.Unscoped().First(&deletedClient, "id = ?", clientID).Error
	assert.NoError(t, err, "Failed to find soft deleted client")
	assert.False(t, deletedClient.DeletedAt.Time.IsZero(), "DeletedAt should be set")

	// Verify client is not found in normal queries
	err = testDB.DB.First(&deletedClient, "id = ?", clientID).Error
	assert.Error(t, err, "Should not find soft deleted client")
}
