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

// createTestUserForClient creates a test user for client tests
func createTestUserForClient(t *testing.T, userRepo *repository.UserRepository, email string) *domain.User {
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

// createTestBusinessForClient creates a test business for client tests
func createTestBusinessForClient(t *testing.T, businessRepo *repository.BusinessRepository, ownerID uuid.UUID, name string) *domain.Business {
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

func TestClientRepositoryIntegration_Create(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForClient(t, userRepo, "client_create@test.com")
	business := createTestBusinessForClient(t, businessRepo, user.UserID, "create-client-salon")

	// Create a new client
	client := &domain.Client{
		BusinessID: business.BusinessID,
		UserID:     &user.UserID,
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		Phone:      "+1234567890",
		Notes:      "Test client notes",
		CreatedBy:  &user.UserID,
	}

	// Test creation
	err = clientRepo.Create(ctx, client)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, client.ID, "Client ID should be generated")
	assert.NotZero(t, client.CreatedAt, "Created at timestamp should be set")

	// Verify the client was created
	result, err := clientRepo.GetByID(ctx, client.ID)
	assert.NoError(t, err)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, user.UserID, *result.UserID)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, "john.doe@example.com", result.Email)
	assert.Equal(t, "+1234567890", result.Phone)
	assert.Equal(t, "Test client notes", result.Notes)
}

func TestClientRepositoryIntegration_GetByID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForClient(t, userRepo, "client_getbyid@test.com")
	business := createTestBusinessForClient(t, businessRepo, user.UserID, "getbyid-client-salon")

	// Create client
	client := &domain.Client{
		BusinessID: business.BusinessID,
		UserID:     &user.UserID,
		FirstName:  "Jane",
		LastName:   "Smith",
		Email:      "jane.smith@example.com",
		Phone:      "+9876543210",
		CreatedBy:  &user.UserID,
	}

	err = clientRepo.Create(ctx, client)
	require.NoError(t, err)

	// Test GetByID
	result, err := clientRepo.GetByID(ctx, client.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, client.ID, result.ID)
	assert.Equal(t, "Jane", result.FirstName)
	assert.Equal(t, "Smith", result.LastName)

	// Verify related entities are populated
	assert.NotNil(t, result.User)
	assert.Equal(t, user.UserID, result.User.UserID)
	
	assert.NotNil(t, result.Business)
	assert.Equal(t, business.BusinessID, result.Business.BusinessID)
}

func TestClientRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Test GetByID with non-existent ID
	result, err := clientRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestClientRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test users and business
	user := createTestUserForClient(t, userRepo, "client_getbybusiness@test.com")
	business := createTestBusinessForClient(t, businessRepo, user.UserID, "getbybusiness-client-salon")

	// Create multiple clients for the business
	clients := []*domain.Client{
		{
			BusinessID: business.BusinessID,
			FirstName:  "Client",
			LastName:   "One",
			Email:      "client1@example.com",
			Phone:      "+1111111111",
			
			CreatedBy:  &user.UserID,
		},
		{
			BusinessID: business.BusinessID,
			FirstName:  "Client",
			LastName:   "Two",
			Email:      "client2@example.com",
			Phone:      "+2222222222",
			
			CreatedBy:  &user.UserID,
		},
		{
			BusinessID: business.BusinessID,
			FirstName:  "Client",
			LastName:   "Three",
			Email:      "client3@example.com",
			Phone:      "+3333333333",
			
			CreatedBy:  &user.UserID,
		},
	}

	for _, c := range clients {
		err := clientRepo.Create(ctx, c)
		require.NoError(t, err)
	}

	// Test ListByBusiness
	results, err := clientRepo.ListByBusiness(ctx, business.BusinessID, 1, 100)
	assert.NoError(t, err)
	assert.Len(t, results, 3, "Should return all clients for the business")

	// Verify all clients belong to the business
	for _, c := range results {
		assert.Equal(t, business.BusinessID, c.BusinessID)
	}
}

func TestClientRepositoryIntegration_GetByBusinessAndEmail(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForClient(t, userRepo, "client_getbyemail@test.com")
	business := createTestBusinessForClient(t, businessRepo, user.UserID, "getbyemail-client-salon")

	// Create client
	client := &domain.Client{
		BusinessID: business.BusinessID,
		FirstName:  "Email",
		LastName:   "Test",
		Email:      "email.test@example.com",
		Phone:      "+4444444444",
		
		CreatedBy:  &user.UserID,
	}

	err = clientRepo.Create(ctx, client)
	require.NoError(t, err)

	// Test GetByBusinessAndEmail
	result, err := clientRepo.GetByBusinessAndEmail(ctx, business.BusinessID, "email.test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, client.ID, result.ID)
	assert.Equal(t, "email.test@example.com", result.Email)
}

func TestClientRepositoryIntegration_Update(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForClient(t, userRepo, "client_update@test.com")
	business := createTestBusinessForClient(t, businessRepo, user.UserID, "update-client-salon")

	// Create client
	client := &domain.Client{
		BusinessID:       business.BusinessID,
		FirstName:        "Original",
		LastName:         "Name",
		Email:            "original@example.com",
		Phone:            "+5555555555",
		Notes:            "Original notes",
		CreatedBy:        &user.UserID,
	}

	err = clientRepo.Create(ctx, client)
	require.NoError(t, err)

	// Update the client
	firstName := "Updated"
	lastName := "Client"
	email := "updated@example.com"
	phone := "+6666666666"
	notes := "Updated notes"

	updateInput := &domain.UpdateClientInput{
		FirstName:        &firstName,
		LastName:         &lastName,
		Email:            &email,
		Phone:            &phone,
		Notes:            &notes,
		
	}

	err = clientRepo.Update(ctx, client.ID, updateInput, user.UserID)
	assert.NoError(t, err)

	// Verify the update
	updated, err := clientRepo.GetByID(ctx, client.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.FirstName)
	assert.Equal(t, "Client", updated.LastName)
	assert.Equal(t, "updated@example.com", updated.Email)
	assert.Equal(t, "+6666666666", updated.Phone)
	assert.Equal(t, "Updated notes", updated.Notes)
	
}

func TestClientRepositoryIntegration_Delete(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForClient(t, userRepo, "client_delete@test.com")
	business := createTestBusinessForClient(t, businessRepo, user.UserID, "delete-client-salon")

	// Create client
	client := &domain.Client{
		BusinessID: business.BusinessID,
		FirstName:  "Delete",
		LastName:   "Me",
		Email:      "delete.me@example.com",
		Phone:      "+7777777777",
		
		CreatedBy:  &user.UserID,
	}

	err = clientRepo.Create(ctx, client)
	require.NoError(t, err)

	// Delete the client
	err = clientRepo.Delete(ctx, client.ID, user.UserID)
	assert.NoError(t, err)

	// Verify deletion
	result, err := clientRepo.GetByID(ctx, client.ID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestClientRepositoryIntegration_ListWithPagination(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForClient(t, userRepo, "client_pagination@test.com")
	business := createTestBusinessForClient(t, businessRepo, user.UserID, "pagination-client-salon")

	// Create multiple clients
	for i := 0; i < 5; i++ {
		client := &domain.Client{
			BusinessID: business.BusinessID,
			FirstName:  "Client",
			LastName:   uuid.New().String()[:8],
			Email:      uuid.New().String()[:8] + "@example.com",
			Phone:      "+123456789" + string(rune('0'+i)),
			
			CreatedBy:  &user.UserID,
		}
		err := clientRepo.Create(ctx, client)
		require.NoError(t, err)
	}

	// Test pagination
	page1, err := clientRepo.ListByBusiness(ctx, business.BusinessID, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, err := clientRepo.ListByBusiness(ctx, business.BusinessID, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, page2, 2)

	page3, err := clientRepo.ListByBusiness(ctx, business.BusinessID, 3, 2)
	assert.NoError(t, err)
	assert.Len(t, page3, 1)

	// Test count
	count, err := clientRepo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestClientRepositoryIntegration_Search(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForClient(t, userRepo, "client_search@test.com")
	business := createTestBusinessForClient(t, businessRepo, user.UserID, "search-client-salon")

	// Create clients with different names
	clients := []*domain.Client{
		{
			BusinessID: business.BusinessID,
			FirstName:  "John",
			LastName:   "Smith",
			Email:      "john.smith@example.com",
			Phone:      "+8888888888",
			
			CreatedBy:  &user.UserID,
		},
		{
			BusinessID: business.BusinessID,
			FirstName:  "Jane",
			LastName:   "Johnson",
			Email:      "jane.johnson@example.com",
			Phone:      "+9999999999",
			
			CreatedBy:  &user.UserID,
		},
		{
			BusinessID: business.BusinessID,
			FirstName:  "Bob",
			LastName:   "Williams",
			Email:      "bob.williams@example.com",
			Phone:      "+0000000000",
			
			CreatedBy:  &user.UserID,
		},
	}

	for _, c := range clients {
		err := clientRepo.Create(ctx, c)
		require.NoError(t, err)
	}

	// Test search by first name
	results, err := clientRepo.Search(ctx, business.BusinessID, "john", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, results, 2, "Should find John Smith and Jane Johnson")

	// Test search by last name
	results, err = clientRepo.Search(ctx, business.BusinessID, "smith", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, results, 1, "Should find John Smith")

	// Test search by email
	results, err = clientRepo.Search(ctx, business.BusinessID, "jane.johnson", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, results, 1, "Should find Jane Johnson")
}