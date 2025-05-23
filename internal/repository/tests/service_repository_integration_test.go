package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/repository"
)

// createTestUserForService creates a test user for service tests
func createTestUserForService(t *testing.T, userRepo *repository.UserRepository, email string) *domain.User {
	user := &domain.User{
		ClerkID:   uuid.New(),
		Email:     email,
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+351123456789",
		Role:      "staff",
		IsActive:  true,
	}

	err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)
	return user
}

// createTestBusinessForService creates a test business for service tests
func createTestBusinessForService(t *testing.T, businessRepo *repository.BusinessRepository, ownerID uuid.UUID, name string) *domain.Business {
	business := &domain.Business{
		OwnerID:          ownerID,
		BusinessName:     name,
		BusinessType:     "beauty_salon",
		AddressLine1:     "123 Test Street",
		City:             "Lisbon",
		Country:          "Portugal",
		Phone:            "+351987654321",
		Email:            name + "@salon.com",
		TimeZone:         "Europe/Lisbon",
		SubscriptionPlan: "basic",
		IsActive:         true,
	}

	err := businessRepo.Create(context.Background(), business)
	require.NoError(t, err)
	return business
}

// createTestServiceCategory creates a test service category
func createTestServiceCategory(t *testing.T, categoryRepo *repository.ServiceCategoryRepository, businessID uuid.UUID, name string) *domain.ServiceCategory {
	category := &domain.ServiceCategory{
		BusinessID:  businessID,
		Name:        name,
		Description: "Test category description",
	}

	err := categoryRepo.Create(context.Background(), category)
	require.NoError(t, err)
	return category
}

func TestServiceRepositoryIntegration_Create(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForService(t, userRepo, "service_create@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "create-service-salon")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Facial Treatments")

	service := &domain.Service{
		BusinessID:  business.BusinessID,
		CategoryID:  &category.ID,
		Name:        "Facial Treatment",
		Description: "Deep cleansing facial",
		Duration:    60,
		Price:       120.00,
		
		CreatedBy:   &user.UserID,
	}

	err = serviceRepo.Create(ctx, service)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, service.ID)
	assert.NotZero(t, service.CreatedAt)
}

func TestServiceRepositoryIntegration_GetByID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test data
	user := createTestUserForService(t, userRepo, "service_getbyid@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbyid-service-salon")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")

	// Create service
	service := &domain.Service{
		BusinessID:  business.BusinessID,
		CategoryID:  &category.ID,
		Name:        "Hair Cut",
		Description: "Professional hair cut",
		Duration:    45,
		Price:       50.00,
		
		CreatedBy:   &user.UserID,
	}

	err = serviceRepo.Create(ctx, service)
	require.NoError(t, err)

	// Test retrieval
	retrievedService, err := serviceRepo.GetByID(ctx, service.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedService)

	assert.Equal(t, service.ID, retrievedService.ID)
	assert.Equal(t, service.BusinessID, retrievedService.BusinessID)
	assert.Equal(t, service.Name, retrievedService.Name)
	assert.Equal(t, service.Description, retrievedService.Description)
	assert.Equal(t, service.Duration, retrievedService.Duration)
	assert.Equal(t, service.Price, retrievedService.Price)

	// Verify related entities are populated
	assert.NotNil(t, retrievedService.Business)
	assert.Equal(t, business.BusinessID, retrievedService.Business.BusinessID)
	
	// Verify category ID is set correctly
	assert.NotNil(t, retrievedService.CategoryID)
	assert.Equal(t, category.ID, *retrievedService.CategoryID)
}

func TestServiceRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	service, err := serviceRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "not found")
}

func TestServiceRepositoryIntegration_Update(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test data
	user := createTestUserForService(t, userRepo, "service_update@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "update-service-salon")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Nail Services")

	// Create service
	service := &domain.Service{
		BusinessID:  business.BusinessID,
		CategoryID:  &category.ID,
		Name:        "Basic Manicure",
		Description: "Standard manicure service",
		Duration:    30,
		Price:       35.00,
		
		CreatedBy:   &user.UserID,
	}

	err = serviceRepo.Create(ctx, service)
	require.NoError(t, err)

	// Update the service
	updatedName := "Deluxe Manicure"
	updatedDescription := "Premium manicure with spa treatment"
	updatedPrice := 55.00
	updatedDuration := 45

	updateInput := &domain.UpdateServiceInput{
		Name:        &updatedName,
		Description: &updatedDescription,
		Price:       &updatedPrice,
		Duration:    &updatedDuration,
	}

	err = serviceRepo.Update(ctx, service.ID, updateInput, user.UserID)
	require.NoError(t, err)

	// Verify the update
	updatedService, err := serviceRepo.GetByID(ctx, service.ID)
	require.NoError(t, err)
	assert.Equal(t, updatedName, updatedService.Name)
	assert.Equal(t, updatedDescription, updatedService.Description)
	assert.Equal(t, updatedPrice, updatedService.Price)
	assert.Equal(t, updatedDuration, updatedService.Duration)
}

func TestServiceRepositoryIntegration_Delete(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test data
	user := createTestUserForService(t, userRepo, "service_delete@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "delete-service-salon")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Spa Services")

	// Create service
	service := &domain.Service{
		BusinessID:  business.BusinessID,
		CategoryID:  &category.ID,
		Name:        "Relaxation Massage",
		Description: "60 minute full body massage",
		Duration:    60,
		Price:       80.00,
		
		CreatedBy:   &user.UserID,
	}

	err = serviceRepo.Create(ctx, service)
	require.NoError(t, err)

	// Delete the service
	err = serviceRepo.Delete(ctx, service.ID, user.UserID)
	require.NoError(t, err)

	// Verify deletion - should not be found
	result, err := serviceRepo.GetByID(ctx, service.ID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestServiceRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test data
	user := createTestUserForService(t, userRepo, "service_getbybusiness@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbybusiness-service-salon")
	category1 := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Category 1")
	category2 := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Category 2")

	// Create multiple services
	services := []*domain.Service{
		{
			BusinessID:  business.BusinessID,
			CategoryID:  &category1.ID,
			Name:        "Service 1",
			Description: "Description 1",
			Duration:    30,
			Price:       40.00,
			
			CreatedBy:   &user.UserID,
		},
		{
			BusinessID:  business.BusinessID,
			CategoryID:  &category1.ID,
			Name:        "Service 2",
			Description: "Description 2",
			Duration:    60,
			Price:       80.00,
			
			CreatedBy:   &user.UserID,
		},
		{
			BusinessID:  business.BusinessID,
			CategoryID:  &category2.ID,
			Name:        "Service 3",
			Description: "Description 3",
			Duration:    90,
			Price:       120.00,
			CreatedBy:   &user.UserID,
		},
	}

	for _, svc := range services {
		err := serviceRepo.Create(ctx, svc)
		require.NoError(t, err)
	}

	// Test ListByBusiness
	results, err := serviceRepo.ListByBusiness(ctx, business.BusinessID, 1, 100)
	require.NoError(t, err)
	assert.Len(t, results, 3) // Should return all services including inactive

	// Verify all services belong to the business
	for _, svc := range results {
		assert.Equal(t, business.BusinessID, svc.BusinessID)
	}
}

func TestServiceRepositoryIntegration_ListByCategory(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test data
	user := createTestUserForService(t, userRepo, "service_getbycategory@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbycategory-service-salon")
	category1 := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	category2 := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Nail Services")

	// Create services in different categories
	services := []*domain.Service{
		{
			BusinessID:  business.BusinessID,
			CategoryID:  &category1.ID,
			Name:        "Hair Cut",
			Duration:    30,
			Price:       40.00,
			
			CreatedBy:   &user.UserID,
		},
		{
			BusinessID:  business.BusinessID,
			CategoryID:  &category1.ID,
			Name:        "Hair Color",
			Duration:    90,
			Price:       120.00,
			
			CreatedBy:   &user.UserID,
		},
		{
			BusinessID:  business.BusinessID,
			CategoryID:  &category2.ID,
			Name:        "Manicure",
			Duration:    45,
			Price:       35.00,
			
			CreatedBy:   &user.UserID,
		},
	}

	for _, svc := range services {
		err := serviceRepo.Create(ctx, svc)
		require.NoError(t, err)
	}

	// Test ListByCategory for category1
	results, err := serviceRepo.ListByCategory(ctx, category1.ID, 1, 10)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify all services belong to category1
	for _, svc := range results {
		assert.NotNil(t, svc.ID)
		assert.Equal(t, category1.ID, *svc.CategoryID)
	}

	// Test ListByCategory for category2
	results, err = serviceRepo.ListByCategory(ctx, category2.ID, 1, 10)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, category2.ID, *results[0].CategoryID)
}

func TestServiceRepositoryIntegration_ListWithPagination(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test data
	user := createTestUserForService(t, userRepo, "service_pagination@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "pagination-service-salon")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "All Services")

	// Create multiple services
	for i := 0; i < 5; i++ {
		service := &domain.Service{
			BusinessID:  business.BusinessID,
			CategoryID:  &category.ID,
			Name:        "Service " + uuid.New().String()[:8],
			Duration:    30 + i*15,
			Price:       float64(40 + i*10),
			
			CreatedBy:   &user.UserID,
		}
		err := serviceRepo.Create(ctx, service)
		require.NoError(t, err)
	}

	// Test pagination
	page1, err := serviceRepo.ListByBusiness(ctx, business.BusinessID, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, err := serviceRepo.ListByBusiness(ctx, business.BusinessID, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, page2, 2)

	page3, err := serviceRepo.ListByBusiness(ctx, business.BusinessID, 3, 2)
	assert.NoError(t, err)
	assert.Len(t, page3, 1)

	// Test count
	count, err := serviceRepo.CountByBusiness(ctx, business.BusinessID)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}