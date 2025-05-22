package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/assimoes/beautix/internal/domain"
)

func TestServiceCategoryRepositoryIntegration_Create(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()

	category := &domain.ServiceCategory{
		Name:        "Hair Services",
		Description: "All hair-related services",
	}

	err := repos.ServiceCategoryRepo.Create(context.Background(), category)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, category.ID)
	assert.NotZero(t, category.CreatedAt)
}

func TestServiceCategoryRepositoryIntegration_GetByID(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Use the category from test data
	retrievedCategory, err := repos.ServiceCategoryRepo.GetByID(context.Background(), testData.ServiceCategory.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedCategory)

	assert.Equal(t, testData.ServiceCategory.ID, retrievedCategory.ID)
	assert.Equal(t, testData.ServiceCategory.Name, retrievedCategory.Name)
	assert.Equal(t, testData.ServiceCategory.Description, retrievedCategory.Description)
}

func TestServiceCategoryRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()

	category, err := repos.ServiceCategoryRepo.GetByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Nil(t, category)
}

func TestServiceCategoryRepositoryIntegration_Update(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Update the category
	updatedName := "Updated Category Name"
	updatedDescription := "Updated description"

	err := repos.ServiceCategoryRepo.Update(context.Background(), testData.ServiceCategory.ID, updatedName, updatedDescription, testData.User.ID)
	require.NoError(t, err)

	// Verify the update
	updatedCategory, err := repos.ServiceCategoryRepo.GetByID(context.Background(), testData.ServiceCategory.ID)
	require.NoError(t, err)
	assert.Equal(t, updatedName, updatedCategory.Name)
	assert.Equal(t, updatedDescription, updatedCategory.Description)
}

func TestServiceCategoryRepositoryIntegration_Delete(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Delete the category
	err := repos.ServiceCategoryRepo.Delete(context.Background(), testData.ServiceCategory.ID, testData.User.ID)
	require.NoError(t, err)

	// Verify the category is deleted (soft delete)
	deletedCategory, err := repos.ServiceCategoryRepo.GetByID(context.Background(), testData.ServiceCategory.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedCategory)
}

func TestServiceCategoryRepositoryIntegration_List(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Create additional categories
	createdBy := testData.User.ID
	createTestServiceCategoryTx(t, suite.Tx, &createdBy)
	createTestServiceCategoryTx(t, suite.Tx, &createdBy)

	// List categories
	categories, err := repos.ServiceCategoryRepo.List(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.Len(t, categories, 3) // 2 created + 1 from test data
}

func TestServiceCategoryRepositoryIntegration_Count(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Create additional categories
	createdBy := testData.User.ID
	createTestServiceCategoryTx(t, suite.Tx, &createdBy)
	createTestServiceCategoryTx(t, suite.Tx, &createdBy)

	// Count categories
	count, err := repos.ServiceCategoryRepo.Count(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(3), count) // 2 created + 1 from test data
}

func TestServiceCategoryRepositoryIntegration_Pagination(t *testing.T) {
	suite := NewTransactionTestSuite(t)
	repos := suite.CreateTestRepositories()
	testData := suite.CreateTestData()

	// Create multiple categories
	createdBy := testData.User.ID
	for i := 0; i < 5; i++ {
		createTestServiceCategoryTx(t, suite.Tx, &createdBy)
	}

	// Test pagination
	page1Categories, err := repos.ServiceCategoryRepo.List(context.Background(), 1, 3)
	require.NoError(t, err)
	assert.Len(t, page1Categories, 3)

	page2Categories, err := repos.ServiceCategoryRepo.List(context.Background(), 2, 3)
	require.NoError(t, err)
	assert.Len(t, page2Categories, 3) // 5 created + 1 from test data = 6 total, so page 2 has 3

	// Verify no overlap between pages
	page1IDs := make(map[uuid.UUID]bool)
	for _, category := range page1Categories {
		page1IDs[category.ID] = true
	}

	for _, category := range page2Categories {
		assert.False(t, page1IDs[category.ID], "Category should not appear on both pages")
	}
}