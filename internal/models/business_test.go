package models_test

import (
	"encoding/json"
	"testing"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBusinessModel(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate the models
	err = testDB.AutoMigrate(&models.User{}, &models.Business{}, &models.BusinessLocation{})
	require.NoError(t, err, "Failed to migrate models")

	// Create a user first (since business requires a user)
	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		ClerkID:   "clerk_business_test",
		Email:     "business_test@example.com",
		FirstName: "Business",
		LastName:  "Owner",
		Phone:     "+1234567890",
		Role:      models.UserRoleProvider,
		IsActive:  true,
	}

	// Save the user
	err = testDB.Create(&user).Error
	assert.NoError(t, err, "Failed to create user")

	// Create business
	businessID := uuid.New()
	business := models.Business{
		BaseModel: models.BaseModel{
			ID: businessID,
		},
		UserID:           userID,
		Name:             "beauty-salon-1",
		DisplayName:      "Beauty Salon Example",
		Description:      "A premium beauty salon offering a variety of services",
		Address:          "123 Main St",
		City:             "Lisbon",
		Country:          "Portugal",
		PostalCode:       "1000-100",
		Phone:            "+351123456789",
		Email:            "contact@beautysalon.com",
		Website:          "https://beautysalon.com",
		SubscriptionTier: models.SubscriptionTierPro,
		IsActive:         true,
		SocialLinks: models.SocialLinks{
			Instagram: "beautysalon",
			Facebook:  "beautysalon",
		},
		Settings: models.BusinessSettings{
			AllowOnlineBooking:       true,
			RequireDeposit:           true,
			DepositAmount:            15.00,
			CancellationPolicyHours:  24,
			CancellationFeePercentage: 50,
			WorkingHours: models.WorkingHours{
				Monday: models.DayHours{
					IsOpen:    true,
					OpenTime:  "09:00",
					CloseTime: "18:00",
				},
				Tuesday: models.DayHours{
					IsOpen:    true,
					OpenTime:  "09:00",
					CloseTime: "18:00",
				},
			},
			BookingNotificationsEnabled: true,
		},
	}

	// Save the business
	err = testDB.Create(&business).Error
	assert.NoError(t, err, "Failed to create business")

	// Verify business was created with ID
	var savedBusiness models.Business
	err = testDB.First(&savedBusiness, "id = ?", businessID).Error
	assert.NoError(t, err, "Failed to find business")
	assert.Equal(t, businessID, savedBusiness.ID)
	assert.Equal(t, "Beauty Salon Example", savedBusiness.DisplayName)
	assert.Equal(t, models.SubscriptionTierPro, savedBusiness.SubscriptionTier)

	// Test JSON fields
	assert.Equal(t, "beautysalon", savedBusiness.SocialLinks.Instagram)
	assert.True(t, savedBusiness.Settings.AllowOnlineBooking)
	assert.Equal(t, 24, savedBusiness.Settings.CancellationPolicyHours)

	// Test loaded relationship
	err = testDB.Preload("User").First(&savedBusiness, "id = ?", businessID).Error
	assert.NoError(t, err, "Failed to find business with user")
	assert.Equal(t, userID, savedBusiness.User.ID)
	assert.Equal(t, "Business", savedBusiness.User.FirstName)

	// Create business location
	locationID := uuid.New()
	location := models.BusinessLocation{
		BaseModel: models.BaseModel{
			ID: locationID,
		},
		BusinessID: businessID,
		Name:       "Downtown Branch",
		Address:    "456 Center St",
		City:       "Lisbon",
		Country:    "Portugal",
		PostalCode: "1000-200",
		Phone:      "+351987654321",
		IsMain:     true,
		Settings: models.LocationSettings{
			WorkingHours: models.WorkingHours{
				Monday: models.DayHours{
					IsOpen:    true,
					OpenTime:  "10:00",
					CloseTime: "19:00",
				},
			},
			Capacity: 5,
		},
	}

	// Save the location
	err = testDB.Create(&location).Error
	assert.NoError(t, err, "Failed to create business location")

	// Verify location was created and correctly associated with business
	var locations []models.BusinessLocation
	err = testDB.Preload("Business").Where("business_id = ?", businessID).Find(&locations).Error
	assert.NoError(t, err, "Failed to find business locations")
	assert.Len(t, locations, 1, "Should have one business location")
	assert.Equal(t, "Downtown Branch", locations[0].Name)
	assert.Equal(t, businessID, locations[0].Business.ID)

	// Test serialization and deserialization of JSON fields
	jsonBytes, err := json.Marshal(business)
	assert.NoError(t, err, "Failed to marshal business to JSON")

	var unmarshaledBusiness models.Business
	err = json.Unmarshal(jsonBytes, &unmarshaledBusiness)
	assert.NoError(t, err, "Failed to unmarshal business from JSON")
	assert.Equal(t, business.DisplayName, unmarshaledBusiness.DisplayName)
	assert.Equal(t, business.SocialLinks.Instagram, unmarshaledBusiness.SocialLinks.Instagram)
	assert.Equal(t, business.Settings.CancellationPolicyHours, unmarshaledBusiness.Settings.CancellationPolicyHours)

	// Test multiple locations for the same business
	location2 := models.BusinessLocation{
		BusinessID: businessID,
		Name:       "Suburbs Branch",
		Address:    "789 Outer St",
		City:       "Lisbon",
		Country:    "Portugal",
		IsMain:     false,
	}

	// Save the second location
	err = testDB.Create(&location2).Error
	assert.NoError(t, err, "Failed to create second business location")

	// Verify that we now have two locations
	err = testDB.Preload("Business").Where("business_id = ?", businessID).Find(&locations).Error
	assert.NoError(t, err, "Failed to find business locations")
	assert.Len(t, locations, 2, "Should have two business locations")

	// Test soft delete
	err = testDB.Delete(&business).Error
	assert.NoError(t, err, "Failed to soft delete business")

	// Verify business is soft deleted
	var deletedBusiness models.Business
	err = testDB.Unscoped().First(&deletedBusiness, "id = ?", businessID).Error
	assert.NoError(t, err, "Failed to find soft deleted business")
	assert.False(t, deletedBusiness.DeletedAt.Time.IsZero(), "DeletedAt should be set")

	// Verify we can't find the business with normal queries
	err = testDB.First(&models.Business{}, "id = ?", businessID).Error
	assert.Error(t, err, "Should not find soft deleted business")
}