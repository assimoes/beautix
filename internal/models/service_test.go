package models_test

import (
	"testing"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceModels(t *testing.T) {
	// Connect to the test database using simple approach
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Clean up all tables comprehensively to avoid foreign key issues
	database.CleanupAllTables(t, testDB.DB)

	// Models are already migrated by the database migration system

	// Create a user first
	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{
			ID:        userID,
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		ClerkID:   "clerk_service_test_" + userID.String()[:8], // Unique ClerkID
		Email:     "service@example.com",
		FirstName: "Service",
		LastName:  "User",
		Role:      models.UserRoleOwner,
	}

	err = testDB.DB.Create(&user).Error
	require.NoError(t, err, "Failed to create user")

	// Create a business first
	businessID := uuid.New()
	business := models.Business{
		BaseModel: models.BaseModel{
			ID:        businessID,
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		UserID:       userID,
		Name:         "Service Test Business",
		BusinessType: "salon",
		Phone:        "+351123456789",
		Email:        "business@example.com",
		Country:      "Portugal",
		TimeZone:     "Europe/Lisbon",
	}

	err = testDB.DB.Create(&business).Error
	require.NoError(t, err, "Failed to create business")

	// Create a service category
	categoryID := uuid.New()
	category := models.ServiceCategory{
		ID:          categoryID,
		BusinessID:  businessID,
		Name:        "Hair Services",
		Description: "All hair-related services",
	}

	err = testDB.DB.Create(&category).Error
	require.NoError(t, err, "Failed to create service category")

	// Create a service
	serviceID := uuid.New()
	service := models.Service{
		BaseModel: models.BaseModel{
			ID:        serviceID,
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		BusinessID:      businessID,
		Name:            "Haircut",
		Description:     "Basic haircut service",
		Duration:        45,
		Price:           35.50,
		Category:        "hair",
		Color:           "#FF5733",
		IsActive:        true,
		PreparationTime: 5,
		CleanupTime:     10,
	}

	// Save the service
	err = testDB.DB.Create(&service).Error
	assert.NoError(t, err, "Failed to create service")

	// Verify service was created
	var savedService models.Service
	err = testDB.DB.First(&savedService, "id = ?", serviceID).Error
	assert.NoError(t, err, "Failed to find service")
	assert.Equal(t, serviceID, savedService.ID)
	assert.Equal(t, "Haircut", savedService.Name)
	assert.Equal(t, "Basic haircut service", savedService.Description)
	assert.Equal(t, 45, savedService.Duration)
	assert.Equal(t, 35.50, savedService.Price)
	assert.Equal(t, "hair", savedService.Category)
	assert.Equal(t, "#FF5733", savedService.Color)
	assert.True(t, savedService.IsActive)
	assert.Equal(t, 5, savedService.PreparationTime)
	assert.Equal(t, 10, savedService.CleanupTime)

	// Test loaded relationships
	err = testDB.DB.Preload("Business").First(&savedService, "id = ?", serviceID).Error
	assert.NoError(t, err, "Failed to find service with relationships")
	assert.Equal(t, businessID, savedService.BusinessID)
	assert.Equal(t, "Service Test Business", savedService.Business.Name)

	// Test update service
	savedService.Price = 40.00
	savedService.UpdatedBy = &userID
	err = testDB.DB.Save(&savedService).Error
	assert.NoError(t, err, "Failed to update service")

	// Verify update
	var updatedService models.Service
	err = testDB.DB.First(&updatedService, "id = ?", serviceID).Error
	assert.NoError(t, err, "Failed to find updated service")
	assert.Equal(t, 40.00, updatedService.Price)

	// Test soft delete
	err = testDB.DB.Delete(&savedService).Error
	assert.NoError(t, err, "Failed to soft delete service")

	// Verify soft delete
	var deletedService models.Service
	err = testDB.DB.Unscoped().First(&deletedService, "id = ?", serviceID).Error
	assert.NoError(t, err, "Failed to find soft deleted service")
	assert.False(t, deletedService.DeletedAt.Time.IsZero(), "DeletedAt should be set")

	// Verify service is not found in normal queries
	err = testDB.DB.First(&deletedService, "id = ?", serviceID).Error
	assert.Error(t, err, "Should not find soft deleted service")
}
