package tests

import (
	"context"
	"testing"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClientRepositoryIntegration_Create(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()

	// Create a new client
	ctx := context.Background()
	client := &domain.Client{
		ProviderID: testData.Business.ID,
		UserID:     &testData.User.ID,
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		Phone:      "+1234567890",
		Notes:      "Test client notes",
		CreatedBy:  &testData.User.ID,
	}

	// Test creation
	err := repos.ClientRepo.Create(ctx, client)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, client.ID, "Client ID should be generated")
	assert.NotZero(t, client.CreatedAt, "Created at timestamp should be set")

	// Verify the client was created using the repository rather than direct DB access
	result, err := repos.ClientRepo.GetByID(ctx, client.ID)
	assert.NoError(t, err)
	assert.Equal(t, testData.Business.ID, result.ProviderID)
	assert.Equal(t, testData.User.ID, *result.UserID)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, "john.doe@example.com", result.Email)
	assert.Equal(t, "+1234567890", result.Phone)
	assert.Equal(t, "Test client notes", result.Notes)
	assert.Equal(t, testData.User.ID, *result.CreatedBy)
}

func TestClientRepositoryIntegration_GetByID(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()

	// Test GetByID
	ctx := context.Background()
	result, err := repos.ClientRepo.GetByID(ctx, testData.Client.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testData.Client.ID, result.ID)
	assert.Equal(t, testData.Business.ID, result.ProviderID)
	assert.Equal(t, testData.User.ID, *result.UserID)
	
	// Verify related entities are populated
	assert.NotNil(t, result.User)
	assert.Equal(t, testData.User.ID, result.User.UserID)
	assert.NotEmpty(t, result.User.Email)
	
	assert.NotNil(t, result.Provider)
	assert.Equal(t, testData.Business.ID, result.Provider.ID)
}

func TestClientRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := repos.ClientRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestClientRepositoryIntegration_GetByUserID(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create an additional business
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	
	// Create clients for both businesses but with the same user ID
	client2 := createTestClientTx(t, suite.Tx, business2.ID, &testData.User.ID, &testData.User.ID)

	// Test GetByUserID
	ctx := context.Background()
	results, err := repos.ClientRepo.GetByUserID(ctx, testData.User.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	
	// Check that both clients are returned
	clientIDs := []uuid.UUID{results[0].ID, results[1].ID}
	assert.Contains(t, clientIDs, testData.Client.ID)
	assert.Contains(t, clientIDs, client2.ID)
}

func TestClientRepositoryIntegration_Update(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()

	// Create update input
	ctx := context.Background()
	updatedFirstName := "UpdatedFirst"
	updatedLastName := "UpdatedLast"
	updatedEmail := "updated.email@example.com"
	updatedPhone := "+9876543210"
	updatedNotes := "Updated notes"
	
	updateInput := &domain.UpdateClientInput{
		FirstName: &updatedFirstName,
		LastName:  &updatedLastName,
		Email:     &updatedEmail,
		Phone:     &updatedPhone,
		Notes:     &updatedNotes,
	}

	// Test Update
	err := repos.ClientRepo.Update(ctx, testData.Client.ID, updateInput, testData.User.ID)
	assert.NoError(t, err)

	// Verify the client was updated
	updated, err := repos.ClientRepo.GetByID(ctx, testData.Client.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedFirstName, updated.FirstName)
	assert.Equal(t, updatedLastName, updated.LastName)
	assert.Equal(t, updatedEmail, updated.Email)
	assert.Equal(t, updatedPhone, updated.Phone)
	assert.Equal(t, updatedNotes, updated.Notes)
	assert.Equal(t, testData.User.ID, *updated.UpdatedBy)
	assert.NotNil(t, updated.UpdatedAt)
}

func TestClientRepositoryIntegration_Delete(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()

	// Test Delete
	ctx := context.Background()
	err := repos.ClientRepo.Delete(ctx, testData.Client.ID, testData.User.ID)
	assert.NoError(t, err)

	// Verify the client was soft deleted by trying to retrieve it
	_, err = repos.ClientRepo.GetByID(ctx, testData.Client.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestClientRepositoryIntegration_ListByProvider(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	
	// Create 2 more clients for business1 (total 3)
	for i := 0; i < 2; i++ {
		createTestClientTx(t, suite.Tx, testData.Business.ID, &testData.User.ID, &testData.User.ID)
	}
	
	// Create 2 clients for business2
	for i := 0; i < 2; i++ {
		createTestClientTx(t, suite.Tx, business2.ID, &testData.User.ID, &testData.User.ID)
	}

	// Test ListByProvider
	ctx := context.Background()
	clientsBusiness1, err := repos.ClientRepo.ListByProvider(ctx, testData.Business.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, clientsBusiness1, 3)
	
	clientsBusiness2, err := repos.ClientRepo.ListByProvider(ctx, business2.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, clientsBusiness2, 2)
	
	// Verify all clients are for the correct provider
	for _, c := range clientsBusiness1 {
		assert.Equal(t, testData.Business.ID, c.ProviderID)
	}
	
	for _, c := range clientsBusiness2 {
		assert.Equal(t, business2.ID, c.ProviderID)
	}
}

func TestClientRepositoryIntegration_Count(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	
	// Create 2 more clients for business1 (total 3)
	for i := 0; i < 2; i++ {
		createTestClientTx(t, suite.Tx, testData.Business.ID, &testData.User.ID, &testData.User.ID)
	}
	
	// Create 2 clients for business2
	for i := 0; i < 2; i++ {
		createTestClientTx(t, suite.Tx, business2.ID, &testData.User.ID, &testData.User.ID)
	}

	// Test Count
	ctx := context.Background()
	count, err := repos.ClientRepo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count) // 1 + 2 + 2 = 5

	// Test CountByProvider
	count1, err := repos.ClientRepo.CountByProvider(ctx, testData.Business.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count1)
	
	count2, err := repos.ClientRepo.CountByProvider(ctx, business2.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count2)
}