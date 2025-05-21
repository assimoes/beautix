package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientModel(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate the models
	err = testDB.AutoMigrate(
		&models.User{},
		&models.Business{},
		&models.Client{},
		&models.ClientNote{},
		&models.ClientDocument{},
	)
	require.NoError(t, err, "Failed to migrate models")

	// Create a user first
	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		ClerkID:   "clerk_client_test",
		Email:     "client_test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+1234567890",
		Role:      models.UserRoleUser,
		IsActive:  true,
	}

	// Save the user
	err = testDB.Create(&user).Error
	assert.NoError(t, err, "Failed to create user")

	// Create a business
	businessID := uuid.New()
	business := models.Business{
		BaseModel: models.BaseModel{
			ID: businessID,
		},
		UserID:           userID,
		Name:             "client-test-business",
		DisplayName:      "Client Test Business",
		Description:      "A business for testing client models",
		Address:          "123 Main St",
		City:             "Lisbon",
		Country:          "Portugal",
		PostalCode:       "1000-100",
		Phone:            "+351123456789",
		Email:            "contact@clienttestbusiness.com",
		SubscriptionTier: models.SubscriptionTierPro,
		IsActive:         true,
	}

	// Save the business
	err = testDB.Create(&business).Error
	assert.NoError(t, err, "Failed to create business")

	// Create a client
	clientID := uuid.New()
	dob := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	client := models.Client{
		BaseModel: models.BaseModel{
			ID: clientID,
		},
		BusinessID:      businessID,
		UserID:          &userID,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Phone:           "+351987654321",
		DateOfBirth:     &dob,
		Address:         "456 Client St",
		City:            "Lisbon",
		Country:         "Portugal",
		PostalCode:      "1000-200",
		Notes:           "Regular client with monthly appointments",
		ProfileImageURL: "https://example.com/client.jpg",
		Tags:            models.ClientTags{"VIP", "Regular"},
		Preferences: models.ClientPreferences{
			PreferredDays:      []string{"monday", "wednesday"},
			PreferredTimeStart: "14:00",
			PreferredTimeEnd:   "18:00",
			CommunicationPrefs: models.CommunicationPreferences{
				AllowEmail:          true,
				AllowSMS:            true,
				AppointmentReminders: true,
			},
			LanguagePreference: "en",
		},
		HealthInfo: models.HealthInfo{
			Allergies: []string{"Latex"},
			Notes:     "No other health concerns",
		},
		Source:           "Referral",
		AcceptsMarketing: true,
		IsActive:         true,
	}

	// Save the client
	err = testDB.Create(&client).Error
	assert.NoError(t, err, "Failed to create client")

	// Verify client was created with ID
	var savedClient models.Client
	err = testDB.First(&savedClient, "id = ?", clientID).Error
	assert.NoError(t, err, "Failed to find client")
	assert.Equal(t, clientID, savedClient.ID)
	assert.Equal(t, "John", savedClient.FirstName)
	assert.Equal(t, "Doe", savedClient.LastName)
	assert.Equal(t, "+351987654321", savedClient.Phone)

	// Test JSONB fields
	assert.Len(t, savedClient.Tags, 2)
	assert.Contains(t, savedClient.Tags, "VIP")
	assert.Equal(t, "en", savedClient.Preferences.LanguagePreference)
	assert.True(t, savedClient.Preferences.CommunicationPrefs.AllowEmail)
	assert.Contains(t, savedClient.HealthInfo.Allergies, "Latex")

	// Test loaded relationships
	err = testDB.Preload("User").Preload("Business").First(&savedClient, "id = ?", clientID).Error
	assert.NoError(t, err, "Failed to find client with relationships")
	assert.Equal(t, userID, *savedClient.UserID)
	assert.Equal(t, "Test", savedClient.User.FirstName)
	assert.Equal(t, businessID, savedClient.BusinessID)
	assert.Equal(t, "Client Test Business", savedClient.Business.DisplayName)

	// Create a client note
	noteID := uuid.New()
	note := models.ClientNote{
		BaseModel: models.BaseModel{
			ID: noteID,
		},
		BusinessID: businessID,
		ClientID:   clientID,
		Title:      "Important Note",
		Content:    "This client prefers a specific hair treatment",
		IsPrivate:  true,
		Pinned:     true,
	}

	// Save the note
	err = testDB.Create(&note).Error
	assert.NoError(t, err, "Failed to create client note")

	// Verify note was created
	var savedNote models.ClientNote
	err = testDB.Preload("Client").First(&savedNote, "id = ?", noteID).Error
	assert.NoError(t, err, "Failed to find client note")
	assert.Equal(t, noteID, savedNote.ID)
	assert.Equal(t, "Important Note", savedNote.Title)
	assert.Equal(t, clientID, savedNote.Client.ID)
	assert.True(t, savedNote.IsPrivate)
	assert.True(t, savedNote.Pinned)

	// Create a client document
	documentID := uuid.New()
	now := time.Now()
	expiresAt := now.AddDate(1, 0, 0)
	document := models.ClientDocument{
		BaseModel: models.BaseModel{
			ID: documentID,
		},
		BusinessID:          businessID,
		ClientID:            clientID,
		DocumentName:        "Consent Form",
		DocumentType:        "consent_form",
		FileURL:             "https://example.com/documents/consent_form.pdf",
		ContentType:         "application/pdf",
		FileSize:            1024 * 100, // 100KB
		IsSignatureRequired: true,
		SignedAt:            &now,
		ExpiresAt:           &expiresAt,
		IsPrivate:           true,
	}

	// Save the document
	err = testDB.Create(&document).Error
	assert.NoError(t, err, "Failed to create client document")

	// Verify document was created
	var savedDocument models.ClientDocument
	err = testDB.Preload("Client").First(&savedDocument, "id = ?", documentID).Error
	assert.NoError(t, err, "Failed to find client document")
	assert.Equal(t, documentID, savedDocument.ID)
	assert.Equal(t, "Consent Form", savedDocument.DocumentName)
	assert.Equal(t, clientID, savedDocument.Client.ID)
	assert.True(t, savedDocument.IsSignatureRequired)
	assert.NotNil(t, savedDocument.SignedAt)
	assert.NotNil(t, savedDocument.ExpiresAt)

	// Test serialization and deserialization of JSON fields
	jsonBytes, err := json.Marshal(client)
	assert.NoError(t, err, "Failed to marshal client to JSON")

	var unmarshaledClient models.Client
	err = json.Unmarshal(jsonBytes, &unmarshaledClient)
	assert.NoError(t, err, "Failed to unmarshal client from JSON")
	assert.Equal(t, client.FirstName, unmarshaledClient.FirstName)
	assert.Equal(t, client.Tags, unmarshaledClient.Tags)
	assert.Equal(t, client.Preferences.LanguagePreference, unmarshaledClient.Preferences.LanguagePreference)
	assert.Equal(t, client.HealthInfo.Allergies, unmarshaledClient.HealthInfo.Allergies)

	// Test soft delete
	err = testDB.Delete(&client).Error
	assert.NoError(t, err, "Failed to soft delete client")

	// Verify client is soft deleted
	var deletedClient models.Client
	err = testDB.Unscoped().First(&deletedClient, "id = ?", clientID).Error
	assert.NoError(t, err, "Failed to find soft deleted client")
	assert.False(t, deletedClient.DeletedAt.Time.IsZero(), "DeletedAt should be set")

	// Verify we can't find the client with normal queries
	err = testDB.First(&models.Client{}, "id = ?", clientID).Error
	assert.Error(t, err, "Should not find soft deleted client")

	// Test cascade effect on notes and documents - they should still exist
	// but won't be returned in normal queries due to the client being soft deleted
	var noteCount int64
	err = testDB.Model(&models.ClientNote{}).Where("client_id = ?", clientID).Count(&noteCount).Error
	assert.NoError(t, err, "Failed to count notes")
	assert.Equal(t, int64(1), noteCount, "Note should still exist")

	var documentCount int64
	err = testDB.Model(&models.ClientDocument{}).Where("client_id = ?", clientID).Count(&documentCount).Error
	assert.NoError(t, err, "Failed to count documents")
	assert.Equal(t, int64(1), documentCount, "Document should still exist")
}