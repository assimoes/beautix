package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/assimoes/beautix/internal/domain"
)

func TestServiceRepositoryIntegration_Create(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	service := &domain.Service{
		ProviderID:  testData.Business.ID,
		CategoryID:  &testData.ServiceCategory.ID,
		Name:        "Facial Treatment",
		Description: "Deep cleansing facial",
		Duration:    60,
		Price:       120.00,
	}

	err := repos.ServiceRepo.Create(context.Background(), service)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, service.ID)
	assert.NotZero(t, service.CreatedAt)
}

func TestServiceRepositoryIntegration_GetByID(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Use the service from test data
	retrievedService, err := repos.ServiceRepo.GetByID(context.Background(), testData.Service.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedService)

	assert.Equal(t, testData.Service.ID, retrievedService.ID)
	assert.Equal(t, testData.Service.BusinessID, retrievedService.ProviderID)
	assert.Equal(t, testData.Service.Name, retrievedService.Name)
	assert.Equal(t, testData.Service.Description, retrievedService.Description)
	assert.Equal(t, testData.Service.Duration, retrievedService.Duration)
	assert.Equal(t, testData.Service.Price, retrievedService.Price)
}

func TestServiceRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()

	service, err := repos.ServiceRepo.GetByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Nil(t, service)
}

func TestServiceRepositoryIntegration_Update(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Update the service
	updatedName := "Updated Service Name"
	updatedPrice := 150.00
	updatedDuration := 90

	updateInput := &domain.UpdateServiceInput{
		Name:     &updatedName,
		Price:    &updatedPrice,
		Duration: &updatedDuration,
	}

	err := repos.ServiceRepo.Update(context.Background(), testData.Service.ID, updateInput, testData.User.ID)
	require.NoError(t, err)

	// Verify the update
	updatedService, err := repos.ServiceRepo.GetByID(context.Background(), testData.Service.ID)
	require.NoError(t, err)
	assert.Equal(t, updatedName, updatedService.Name)
	assert.Equal(t, updatedPrice, updatedService.Price)
	assert.Equal(t, updatedDuration, updatedService.Duration)
}

func TestServiceRepositoryIntegration_Delete(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Delete the service
	err := repos.ServiceRepo.Delete(context.Background(), testData.Service.ID, testData.User.ID)
	require.NoError(t, err)

	// Verify the service is deleted (soft delete)
	deletedService, err := repos.ServiceRepo.GetByID(context.Background(), testData.Service.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedService)
}

func TestServiceRepositoryIntegration_ListByProvider(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Create additional services for the provider
	createdBy := testData.User.ID
	categoryID := testData.ServiceCategory.ID
	createTestServiceTx(t, suite.Tx, testData.Business.ID, &categoryID, &createdBy)
	createTestServiceTx(t, suite.Tx, testData.Business.ID, &categoryID, &createdBy)

	// Create a service for a different provider
	differentProvider := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	createTestServiceTx(t, suite.Tx, differentProvider.ID, &categoryID, &createdBy)

	// List services for the original provider
	services, err := repos.ServiceRepo.ListByProvider(context.Background(), testData.Business.ID, 1, 10)
	require.NoError(t, err)
	assert.Len(t, services, 3) // 2 created + 1 from test data

	// Verify services belong to the correct provider
	for _, service := range services {
		assert.Equal(t, testData.Business.ID, service.ProviderID)
	}
}

func TestServiceRepositoryIntegration_ListByCategory(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Create another category
	createdBy := testData.User.ID
	category2 := createTestServiceCategoryTx(t, suite.Tx, &createdBy)

	// Create services in different categories
	categoryID1 := testData.ServiceCategory.ID
	categoryID2 := category2.ID
	createTestServiceTx(t, suite.Tx, testData.Business.ID, &categoryID1, &createdBy)
	createTestServiceTx(t, suite.Tx, testData.Business.ID, &categoryID2, &createdBy)

	// List services by the first category
	services, err := repos.ServiceRepo.ListByCategory(context.Background(), testData.ServiceCategory.ID, 1, 10)
	require.NoError(t, err)
	assert.Len(t, services, 2) // 1 created + 1 from test data

	// Verify services belong to the correct category
	for _, service := range services {
		assert.Equal(t, testData.ServiceCategory.ID, *service.CategoryID)
	}

	// List services by the second category
	services2, err := repos.ServiceRepo.ListByCategory(context.Background(), category2.ID, 1, 10)
	require.NoError(t, err)
	assert.Len(t, services2, 1)
}

func TestServiceRepositoryIntegration_Pagination(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Create multiple services
	createdBy := testData.User.ID
	categoryID := testData.ServiceCategory.ID
	for i := 0; i < 5; i++ {
		createTestServiceTx(t, suite.Tx, testData.Business.ID, &categoryID, &createdBy)
	}

	// Test pagination
	page1Services, err := repos.ServiceRepo.ListByProvider(context.Background(), testData.Business.ID, 1, 3)
	require.NoError(t, err)
	assert.Len(t, page1Services, 3)

	page2Services, err := repos.ServiceRepo.ListByProvider(context.Background(), testData.Business.ID, 2, 3)
	require.NoError(t, err)
	assert.Len(t, page2Services, 3) // 5 created + 1 from test data = 6 total, so page 2 has 3

	// Verify no overlap between pages
	page1IDs := make(map[uuid.UUID]bool)
	for _, service := range page1Services {
		page1IDs[service.ID] = true
	}

	for _, service := range page2Services {
		assert.False(t, page1IDs[service.ID], "Service should not appear on both pages")
	}
}