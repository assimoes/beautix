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

func TestServiceModels(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate the models
	err = testDB.AutoMigrate(
		&models.User{},
		&models.Business{},
		&models.ServiceCategory{},
		&models.Service{},
		&models.ServiceVariant{},
		&models.ServiceOption{},
		&models.ServiceOptionChoice{},
		&models.ServiceBundle{},
		&models.ServiceBundleItem{},
	)
	require.NoError(t, err, "Failed to migrate models")

	// Create a user first (since business requires a user)
	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		ClerkID:   "clerk_service_test",
		Email:     "service_test@example.com",
		FirstName: "Service",
		LastName:  "Provider",
		Phone:     "+1234567890",
		Role:      models.UserRoleProvider,
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
		Name:             "beauty-salon-service",
		DisplayName:      "Beauty Salon Services",
		Description:      "A salon for testing service models",
		Address:          "123 Main St",
		City:             "Lisbon",
		Country:          "Portugal",
		PostalCode:       "1000-100",
		Phone:            "+351123456789",
		Email:            "contact@beautysalonservice.com",
		SubscriptionTier: models.SubscriptionTierPro,
		IsActive:         true,
	}

	// Save the business
	err = testDB.Create(&business).Error
	assert.NoError(t, err, "Failed to create business")

	// Create a service category
	categoryID := uuid.New()
	category := models.ServiceCategory{
		BaseModel: models.BaseModel{
			ID: categoryID,
		},
		Name:        "Hair Services",
		Description: "All hair-related services",
	}

	// Save the category
	err = testDB.Create(&category).Error
	assert.NoError(t, err, "Failed to create service category")

	// Verify category was created with ID
	var savedCategory models.ServiceCategory
	err = testDB.First(&savedCategory, "id = ?", categoryID).Error
	assert.NoError(t, err, "Failed to find service category")
	assert.Equal(t, categoryID, savedCategory.ID)
	assert.Equal(t, "Hair Services", savedCategory.Name)

	// Create a service
	serviceID := uuid.New()
	service := models.Service{
		BaseModel: models.BaseModel{
			ID: serviceID,
		},
		BusinessID:  businessID,
		CategoryID:  &categoryID,
		Name:        "Haircut",
		Description: "Basic haircut service",
		Duration:    45,
		Price:       35.50,
		ImageURL:    "https://example.com/haircut.jpg",
		IsActive:    true,
		Tags:        models.ServiceTags{"Hair", "Cut", "Styling"},
		Settings: models.ServiceSettings{
			AllowOnlineBooking:   true,
			MinAdvanceTimeHours:  2,
			MaxAdvanceTimeDays:   30,
			RequireDeposit:       true,
			DepositAmount:        10.00,
			CanBeBooked:          true,
			BufferTimeAfterMin:   15,
			CancellationPolicyHours: 24,
		},
	}

	// Save the service
	err = testDB.Create(&service).Error
	assert.NoError(t, err, "Failed to create service")

	// Verify service was created with ID
	var savedService models.Service
	err = testDB.First(&savedService, "id = ?", serviceID).Error
	assert.NoError(t, err, "Failed to find service")
	assert.Equal(t, serviceID, savedService.ID)
	assert.Equal(t, "Haircut", savedService.Name)
	assert.Equal(t, 45, savedService.Duration)
	assert.Equal(t, 35.50, savedService.Price)

	// Test JSONB fields
	assert.Len(t, savedService.Tags, 3)
	assert.Contains(t, savedService.Tags, "Hair")
	assert.True(t, savedService.Settings.AllowOnlineBooking)
	assert.Equal(t, 24, savedService.Settings.CancellationPolicyHours)

	// Test loaded relationships
	err = testDB.Preload("Business").Preload("Category").First(&savedService, "id = ?", serviceID).Error
	assert.NoError(t, err, "Failed to find service with relationships")
	assert.Equal(t, businessID, savedService.Business.ID)
	assert.Equal(t, "Beauty Salon Services", savedService.Business.DisplayName)
	assert.Equal(t, categoryID, *savedService.CategoryID)
	assert.Equal(t, "Hair Services", savedService.Category.Name)

	// Create a service variant
	variantID := uuid.New()
	variant := models.ServiceVariant{
		BaseModel: models.BaseModel{
			ID: variantID,
		},
		ServiceID:   serviceID,
		Name:        "Premium Haircut",
		Description: "Haircut with premium experience",
		Duration:    60,
		Price:       50.00,
		IsActive:    true,
	}

	// Save the variant
	err = testDB.Create(&variant).Error
	assert.NoError(t, err, "Failed to create service variant")

	// Verify variant was created
	var savedVariant models.ServiceVariant
	err = testDB.Preload("Service").First(&savedVariant, "id = ?", variantID).Error
	assert.NoError(t, err, "Failed to find service variant")
	assert.Equal(t, variantID, savedVariant.ID)
	assert.Equal(t, "Premium Haircut", savedVariant.Name)
	assert.Equal(t, 60, savedVariant.Duration)
	assert.Equal(t, serviceID, savedVariant.Service.ID)

	// Create a service option
	optionID := uuid.New()
	option := models.ServiceOption{
		BaseModel: models.BaseModel{
			ID: optionID,
		},
		ServiceID:     serviceID,
		Name:          "Hair Length",
		Description:   "Select your hair length",
		IsRequired:    true,
		IsMultiple:    false,
		MinSelections: 1,
		MaxSelections: 1,
		IsActive:      true,
	}

	// Save the option
	err = testDB.Create(&option).Error
	assert.NoError(t, err, "Failed to create service option")

	// Create option choices
	choiceID := uuid.New()
	choice := models.ServiceOptionChoice{
		BaseModel: models.BaseModel{
			ID: choiceID,
		},
		OptionID:        optionID,
		Name:            "Long Hair",
		Description:     "For hair longer than shoulder length",
		PriceAdjustment: 10.00,
		TimeAdjustment:  15,
		IsDefault:       false,
		IsActive:        true,
	}

	// Save the choice
	err = testDB.Create(&choice).Error
	assert.NoError(t, err, "Failed to create service option choice")

	// Verify option and choice were created
	var savedOption models.ServiceOption
	err = testDB.Preload("Service").First(&savedOption, "id = ?", optionID).Error
	assert.NoError(t, err, "Failed to find service option")
	assert.Equal(t, optionID, savedOption.ID)
	assert.Equal(t, "Hair Length", savedOption.Name)
	assert.True(t, savedOption.IsRequired)

	var savedChoice models.ServiceOptionChoice
	err = testDB.Preload("Option").First(&savedChoice, "id = ?", choiceID).Error
	assert.NoError(t, err, "Failed to find service option choice")
	assert.Equal(t, choiceID, savedChoice.ID)
	assert.Equal(t, "Long Hair", savedChoice.Name)
	assert.Equal(t, 10.00, savedChoice.PriceAdjustment)
	assert.Equal(t, optionID, savedChoice.Option.ID)

	// Create a service bundle
	bundleID := uuid.New()
	tomorrow := time.Now().Add(24 * time.Hour)
	nextMonth := time.Now().AddDate(0, 1, 0)
	bundle := models.ServiceBundle{
		BaseModel: models.BaseModel{
			ID: bundleID,
		},
		BusinessID:        businessID,
		Name:              "Hair Package",
		Description:       "Complete hair care package",
		Price:             80.00,
		DiscountPercentage: 15,
		ImageURL:          "https://example.com/hair-package.jpg",
		IsActive:          true,
		StartDate:         &tomorrow,
		EndDate:           &nextMonth,
	}

	// Save the bundle
	err = testDB.Create(&bundle).Error
	assert.NoError(t, err, "Failed to create service bundle")

	// Create bundle item
	bundleItemID := uuid.New()
	bundleItem := models.ServiceBundleItem{
		BaseModel: models.BaseModel{
			ID: bundleItemID,
		},
		BundleID:   bundleID,
		ServiceID:  serviceID,
		Quantity:   1,
		IsRequired: true,
	}

	// Save the bundle item
	err = testDB.Create(&bundleItem).Error
	assert.NoError(t, err, "Failed to create service bundle item")

	// Verify bundle and bundle item were created
	var savedBundle models.ServiceBundle
	err = testDB.Preload("Business").First(&savedBundle, "id = ?", bundleID).Error
	assert.NoError(t, err, "Failed to find service bundle")
	assert.Equal(t, bundleID, savedBundle.ID)
	assert.Equal(t, "Hair Package", savedBundle.Name)
	assert.Equal(t, 80.00, savedBundle.Price)
	assert.Equal(t, 15, savedBundle.DiscountPercentage)
	assert.Equal(t, businessID, savedBundle.Business.ID)

	var savedBundleItem models.ServiceBundleItem
	err = testDB.Preload("Bundle").Preload("Service").First(&savedBundleItem, "id = ?", bundleItemID).Error
	assert.NoError(t, err, "Failed to find service bundle item")
	assert.Equal(t, bundleItemID, savedBundleItem.ID)
	assert.Equal(t, bundleID, savedBundleItem.Bundle.ID)
	assert.Equal(t, serviceID, savedBundleItem.Service.ID)
	assert.Equal(t, 1, savedBundleItem.Quantity)
	assert.True(t, savedBundleItem.IsRequired)

	// Test serialization and deserialization of JSON fields
	jsonBytes, err := json.Marshal(service)
	assert.NoError(t, err, "Failed to marshal service to JSON")

	var unmarshaledService models.Service
	err = json.Unmarshal(jsonBytes, &unmarshaledService)
	assert.NoError(t, err, "Failed to unmarshal service from JSON")
	assert.Equal(t, service.Name, unmarshaledService.Name)
	assert.Equal(t, service.Tags, unmarshaledService.Tags)
	assert.Equal(t, service.Settings.RequireDeposit, unmarshaledService.Settings.RequireDeposit)
	assert.Equal(t, service.Settings.CancellationPolicyHours, unmarshaledService.Settings.CancellationPolicyHours)

	// Test soft delete
	err = testDB.Delete(&service).Error
	assert.NoError(t, err, "Failed to soft delete service")

	// Verify service is soft deleted
	var deletedService models.Service
	err = testDB.Unscoped().First(&deletedService, "id = ?", serviceID).Error
	assert.NoError(t, err, "Failed to find soft deleted service")
	assert.False(t, deletedService.DeletedAt.Time.IsZero(), "DeletedAt should be set")

	// Verify we can't find the service with normal queries
	err = testDB.First(&models.Service{}, "id = ?", serviceID).Error
	assert.Error(t, err, "Should not find soft deleted service")

	// Verify cascade effect on variants and options - they should remain in the db
	// but should be filtered out when loading services normally due to the relationship
	var variantCount int64
	err = testDB.Model(&models.ServiceVariant{}).Where("service_id = ?", serviceID).Count(&variantCount).Error
	assert.NoError(t, err, "Failed to count variants")
	assert.Equal(t, int64(1), variantCount, "Variant should still exist")

	var optionCount int64
	err = testDB.Model(&models.ServiceOption{}).Where("service_id = ?", serviceID).Count(&optionCount).Error
	assert.NoError(t, err, "Failed to count options")
	assert.Equal(t, int64(1), optionCount, "Option should still exist")

	// Service category should be unaffected by service deletion
	err = testDB.First(&models.ServiceCategory{}, "id = ?", categoryID).Error
	assert.NoError(t, err, "Service category should still be found after service deletion")
}