package tests

import (
	"context"
	"testing"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestUserForBusiness creates a test user for business tests
func createTestUserForBusiness(t *testing.T, userRepo *repository.UserRepository, email string) *domain.User {
	user := &domain.User{
		ClerkID:   uuid.New(),
		Email:     email,
		FirstName: "Test",
		LastName:  "Owner",
		Phone:     "+351123456789",
		Role:      "owner",
		IsActive:  true,
	}

	err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)
	return user
}

// TestBusinessRepositoryIntegration_CreateBusiness tests creating a business
func TestBusinessRepositoryIntegration_CreateBusiness(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create owner user first
	owner := createTestUserForBusiness(t, userRepo, "owner@testbusiness.com")

	business := &domain.Business{
		OwnerID:          owner.UserID,
		BusinessName:     "Test Beauty Salon",
		BusinessType:     "beauty_salon",
		AddressLine1:     "123 Test Street",
		City:             "Lisbon",
		Region:           "Lisbon",
		PostalCode:       "1000-001",
		Country:          "Portugal",
		Phone:            "+351123456789",
		Email:            "info@testbeautysalon.com",
		TimeZone:         "Europe/Lisbon",
		SubscriptionPlan: "basic",
		IsActive:         true,
	}

	// Test creation
	err = businessRepo.Create(ctx, business)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, business.BusinessID)
	assert.False(t, business.CreatedAt.IsZero())
	assert.False(t, business.UpdatedAt.IsZero())
}

// TestBusinessRepositoryIntegration_GetBusinessByID tests retrieving a business by ID
func TestBusinessRepositoryIntegration_GetBusinessByID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create owner and business
	owner := createTestUserForBusiness(t, userRepo, "owner_getbyid@test.com")
	
	originalBusiness := &domain.Business{
		OwnerID:          owner.UserID,
		BusinessName:     "get-by-id-salon",
		BusinessType:     "beauty_salon",
		City:             "Lisbon",
		Country:          "Portugal",
		Phone:            "+351987654321",
		Email:            "getbyid@salon.com",
		TimeZone:         "Europe/Lisbon",
		SubscriptionPlan: "pro",
		IsActive:         true,
	}

	err = businessRepo.Create(ctx, originalBusiness)
	require.NoError(t, err)

	// Test retrieval
	retrievedBusiness, err := businessRepo.GetByID(ctx, originalBusiness.BusinessID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedBusiness)
	
	// Verify all fields match
	assert.Equal(t, originalBusiness.BusinessID, retrievedBusiness.BusinessID)
	assert.Equal(t, originalBusiness.OwnerID, retrievedBusiness.OwnerID)
	assert.Equal(t, originalBusiness.BusinessName, retrievedBusiness.BusinessName)
	assert.Equal(t, originalBusiness.Phone, retrievedBusiness.Phone)
	assert.Equal(t, originalBusiness.Email, retrievedBusiness.Email)

	// Test non-existent ID
	nonExistentID := uuid.New()
	retrievedBusiness, err = businessRepo.GetByID(ctx, nonExistentID)
	assert.Error(t, err)
	assert.Nil(t, retrievedBusiness)
}

// TestBusinessRepositoryIntegration_UpdateBusiness tests updating a business
func TestBusinessRepositoryIntegration_UpdateBusiness(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create owner and business
	owner := createTestUserForBusiness(t, userRepo, "owner_update@test.com")
	
	business := &domain.Business{
		OwnerID:          owner.UserID,
		BusinessName:     "update-salon",
		BusinessType:     "beauty_salon",
		Phone:            "+351111111111",
		Email:            "original@salon.com",
		IsActive:         true,
		TimeZone:         "Europe/Lisbon",
		SubscriptionPlan: "basic",
	}

	err = businessRepo.Create(ctx, business)
	require.NoError(t, err)

	// Update the business
	businessName := "Updated Salon Name"
	phone := "+351999999999"
	email := "updated@salon.com"
	city := "Porto"
	isActive := false
	
	updateInput := &domain.UpdateBusinessInput{
		BusinessName: &businessName,
		Phone:        &phone,
		Email:        &email,
		City:         &city,
		IsActive:     &isActive,
	}

	err = businessRepo.Update(ctx, business.BusinessID, updateInput)
	require.NoError(t, err)

	// Retrieve and verify
	updatedBusiness, err := businessRepo.GetByID(ctx, business.BusinessID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Salon Name", updatedBusiness.BusinessName)
	assert.Equal(t, "+351999999999", updatedBusiness.Phone)
	assert.Equal(t, "updated@salon.com", updatedBusiness.Email)
	assert.Equal(t, "Porto", updatedBusiness.City)
	assert.False(t, updatedBusiness.IsActive)
	assert.True(t, updatedBusiness.UpdatedAt.After(business.CreatedAt))
}

// TestBusinessRepositoryIntegration_DeleteBusiness tests deleting a business
func TestBusinessRepositoryIntegration_DeleteBusiness(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create owner and business
	owner := createTestUserForBusiness(t, userRepo, "owner_delete@test.com")
	
	business := &domain.Business{
		OwnerID:          owner.UserID,
		BusinessName:     "delete-salon",
		BusinessType:     "beauty_salon",
		IsActive:         true,
		TimeZone:         "Europe/Lisbon",
		SubscriptionPlan: "basic",
	}

	err = businessRepo.Create(ctx, business)
	require.NoError(t, err)

	// Delete the business
	err = businessRepo.Delete(ctx, business.BusinessID)
	require.NoError(t, err)

	// Try to retrieve - should fail
	retrievedBusiness, err := businessRepo.GetByID(ctx, business.BusinessID)
	assert.Error(t, err)
	assert.Nil(t, retrievedBusiness)
}

// TestBusinessRepositoryIntegration_ListBusinessesByOwner tests listing businesses by owner
func TestBusinessRepositoryIntegration_ListBusinessesByOwner(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create two owners
	owner1 := createTestUserForBusiness(t, userRepo, "owner1_list@test.com")
	owner2 := createTestUserForBusiness(t, userRepo, "owner2_list@test.com")

	// Create businesses for owner1
	businesses := []*domain.Business{
		{
			OwnerID:          owner1.UserID,
			BusinessName:     "salon1",
			BusinessType:     "beauty_salon",
			IsActive:         true,
			TimeZone:         "Europe/Lisbon",
			SubscriptionPlan: "basic",
		},
		{
			OwnerID:          owner1.UserID,
			BusinessName:     "salon2",
			BusinessType:     "beauty_salon",
			IsActive:         true,
			TimeZone:         "Europe/Lisbon",
			SubscriptionPlan: "pro",
		},
		{
			OwnerID:          owner2.UserID,
			BusinessName:     "salon3",
			BusinessType:     "beauty_salon",
			IsActive:         true,
			TimeZone:         "Europe/Lisbon",
			SubscriptionPlan: "basic",
		},
	}

	for _, business := range businesses {
		err := businessRepo.Create(ctx, business)
		require.NoError(t, err)
	}

	// Test listing businesses for owner1
	owner1Businesses, err := businessRepo.GetByOwnerID(ctx, owner1.UserID)
	require.NoError(t, err)
	assert.Len(t, owner1Businesses, 2)

	// Verify all returned businesses belong to owner1
	for _, b := range owner1Businesses {
		assert.Equal(t, owner1.UserID, b.OwnerID)
	}

	// Test listing businesses for owner2
	owner2Businesses, err := businessRepo.GetByOwnerID(ctx, owner2.UserID)
	require.NoError(t, err)
	assert.Len(t, owner2Businesses, 1)
	assert.Equal(t, owner2.UserID, owner2Businesses[0].OwnerID)
}

// TestBusinessRepositoryIntegration_ListBusinesses tests listing businesses with pagination
func TestBusinessRepositoryIntegration_ListBusinesses(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create owner
	owner := createTestUserForBusiness(t, userRepo, "owner_list@test.com")

	// Create multiple businesses
	for i := 0; i < 5; i++ {
		business := &domain.Business{
			OwnerID:          owner.UserID,
			BusinessName:     uuid.New().String()[:8] + "-salon",
			BusinessType:     "beauty_salon",
			IsActive:         true,
			TimeZone:         "Europe/Lisbon",
			SubscriptionPlan: "basic",
		}
		err := businessRepo.Create(ctx, business)
		require.NoError(t, err)
	}

	// Test listing with pagination
	page1Businesses, err := businessRepo.List(ctx, 1, 2)
	require.NoError(t, err)
	assert.Len(t, page1Businesses, 2)

	page2Businesses, err := businessRepo.List(ctx, 2, 2)
	require.NoError(t, err)
	assert.Len(t, page2Businesses, 2)

	page3Businesses, err := businessRepo.List(ctx, 3, 2)
	require.NoError(t, err)
	assert.Len(t, page3Businesses, 1)

	// Test count
	count, err := businessRepo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}