//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/assimoes/beautix/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestClientLoyaltyRepositoryIntegration_Create(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientLoyaltyRepo := &repository.ClientLoyaltyRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForLoyalty(t, userRepo, "client_loyalty@test.com")
	business := createTestBusinessForLoyalty(t, businessRepo, "Loyalty Client Test Salon")
	client := createTestClientForLoyalty(t, clientRepo, user.UserID, business.BusinessID)
	program := createTestLoyaltyProgram(t, testDB.DB.DB, business.BusinessID, "Test Program")
	
	clientLoyalty := &domain.ClientLoyalty{
		ID:        uuid.New(),
		ClientID:  client.ID,
		ProgramID: program.ID,
		Points:    50,
		CreatedAt: time.Now(),
	}

	err = clientLoyaltyRepo.Create(ctx, clientLoyalty)
	require.NoError(t, err)

	// Verify in database
	var dbClientLoyalty models.ClientLoyalty
	err = testDB.DB.DB.Where("id = ?", clientLoyalty.ID).First(&dbClientLoyalty).Error
	require.NoError(t, err)
	assert.Equal(t, client.ID, dbClientLoyalty.ClientID)
	assert.Equal(t, program.ID, dbClientLoyalty.ProgramID)
	assert.Equal(t, 50, dbClientLoyalty.Points)
	assert.NotEmpty(t, dbClientLoyalty.CardNumber)
}

func TestClientLoyaltyRepositoryIntegration_GetByClientAndProgram(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientLoyaltyRepo := &repository.ClientLoyaltyRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForLoyalty(t, userRepo, "get_loyalty@test.com")
	business := createTestBusinessForLoyalty(t, businessRepo, "Get Client Loyalty Test")
	client := createTestClientForLoyalty(t, clientRepo, user.UserID, business.BusinessID)
	program := createTestLoyaltyProgram(t, testDB.DB.DB, business.BusinessID, "Test Program")
	loyalty := createTestClientLoyalty(t, testDB.DB.DB, client.ID, program.ID, business.BusinessID)

	// Test GetByClientAndProgram
	result, err := clientLoyaltyRepo.GetByClientAndProgram(ctx, client.ID, program.ID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, loyalty.ID, result.ID)
	assert.Equal(t, loyalty.Points, result.Points)

	// Test not found
	result, err = clientLoyaltyRepo.GetByClientAndProgram(ctx, uuid.New(), program.ID)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestClientLoyaltyRepositoryIntegration_UpdatePoints(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientLoyaltyRepo := &repository.ClientLoyaltyRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForLoyalty(t, userRepo, "update_points@test.com")
	business := createTestBusinessForLoyalty(t, businessRepo, "Update Points Test")
	client := createTestClientForLoyalty(t, clientRepo, user.UserID, business.BusinessID)
	program := createTestLoyaltyProgram(t, testDB.DB.DB, business.BusinessID, "Test Program")
	createTestClientLoyalty(t, testDB.DB.DB, client.ID, program.ID, business.BusinessID)
	updaterID := uuid.New()

	// Update points
	newPoints := 150
	err = clientLoyaltyRepo.UpdatePoints(ctx, client.ID, program.ID, newPoints, updaterID)
	require.NoError(t, err)

	// Verify update
	updated, err := clientLoyaltyRepo.GetByClientAndProgram(ctx, client.ID, program.ID)
	require.NoError(t, err)
	assert.Equal(t, newPoints, updated.Points)
	assert.Equal(t, updaterID, *updated.UpdatedBy)

	// Verify last activity date was updated
	var dbLoyalty models.ClientLoyalty
	err = testDB.DB.DB.Where("client_id = ? AND program_id = ?", client.ID, program.ID).First(&dbLoyalty).Error
	require.NoError(t, err)
	assert.NotNil(t, dbLoyalty.LastActivityDate)
}

func TestClientLoyaltyRepositoryIntegration_ListByClient(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientLoyaltyRepo := &repository.ClientLoyaltyRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForLoyalty(t, userRepo, "list_by_client@test.com")
	business := createTestBusinessForLoyalty(t, businessRepo, "List By Client Test")
	client := createTestClientForLoyalty(t, clientRepo, user.UserID, business.BusinessID)
	
	// Create multiple programs and client loyalties
	program1 := createTestLoyaltyProgram(t, testDB.DB.DB, business.BusinessID, "Program 1")
	program2 := createTestLoyaltyProgram(t, testDB.DB.DB, business.BusinessID, "Program 2")
	program3 := createTestLoyaltyProgram(t, testDB.DB.DB, business.BusinessID, "Program 3")
	
	loyalty1 := createTestClientLoyalty(t, testDB.DB.DB, client.ID, program1.ID, business.BusinessID)
	time.Sleep(10 * time.Millisecond)
	loyalty2 := createTestClientLoyalty(t, testDB.DB.DB, client.ID, program2.ID, business.BusinessID)
	time.Sleep(10 * time.Millisecond)
	loyalty3 := createTestClientLoyalty(t, testDB.DB.DB, client.ID, program3.ID, business.BusinessID)

	// Create loyalty for another client
	otherUser := createTestUserForLoyalty(t, userRepo, "other_client@test.com")
	otherClient := createTestClientForLoyalty(t, clientRepo, otherUser.UserID, business.BusinessID)
	createTestClientLoyalty(t, testDB.DB.DB, otherClient.ID, program1.ID, business.BusinessID)

	// List loyalties for client
	loyalties, err := clientLoyaltyRepo.ListByClient(ctx, client.ID)
	require.NoError(t, err)
	assert.Len(t, loyalties, 3)

	// Verify ordering (newest first)
	assert.Equal(t, loyalty3.ID, loyalties[0].ID)
	assert.Equal(t, loyalty2.ID, loyalties[1].ID)
	assert.Equal(t, loyalty1.ID, loyalties[2].ID)
}

func TestClientLoyaltyRepositoryIntegration_ListByProgram(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientLoyaltyRepo := &repository.ClientLoyaltyRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	business := createTestBusinessForLoyalty(t, businessRepo, "List By Program Test")
	program := createTestLoyaltyProgram(t, testDB.DB.DB, business.BusinessID, "Test Program")
	
	// Create multiple clients with different points
	user1 := createTestUserForLoyalty(t, userRepo, "client1@test.com")
	client1 := createTestClientForLoyalty(t, clientRepo, user1.UserID, business.BusinessID)
	loyalty1 := createTestClientLoyaltyWithPoints(t, testDB.DB.DB, client1.ID, program.ID, business.BusinessID, 100)
	
	user2 := createTestUserForLoyalty(t, userRepo, "client2@test.com")
	client2 := createTestClientForLoyalty(t, clientRepo, user2.UserID, business.BusinessID)
	loyalty2 := createTestClientLoyaltyWithPoints(t, testDB.DB.DB, client2.ID, program.ID, business.BusinessID, 200)
	
	user3 := createTestUserForLoyalty(t, userRepo, "client3@test.com")
	client3 := createTestClientForLoyalty(t, clientRepo, user3.UserID, business.BusinessID)
	loyalty3 := createTestClientLoyaltyWithPoints(t, testDB.DB.DB, client3.ID, program.ID, business.BusinessID, 150)

	// List clients in program
	loyalties, err := clientLoyaltyRepo.ListByProgram(ctx, program.ID, 1, 10)
	require.NoError(t, err)
	assert.Len(t, loyalties, 3)

	// Verify ordering (by points descending)
	assert.Equal(t, loyalty2.ID, loyalties[0].ID) // 200 points
	assert.Equal(t, loyalty3.ID, loyalties[1].ID) // 150 points
	assert.Equal(t, loyalty1.ID, loyalties[2].ID) // 100 points

	// Test pagination
	loyalties, err = clientLoyaltyRepo.ListByProgram(ctx, program.ID, 2, 2)
	require.NoError(t, err)
	assert.Len(t, loyalties, 1)
	assert.Equal(t, loyalty1.ID, loyalties[0].ID)
}

func TestClientLoyaltyRepositoryIntegration_Delete(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientLoyaltyRepo := &repository.ClientLoyaltyRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForLoyalty(t, userRepo, "delete_loyalty@test.com")
	business := createTestBusinessForLoyalty(t, businessRepo, "Delete Client Loyalty Test")
	client := createTestClientForLoyalty(t, clientRepo, user.UserID, business.BusinessID)
	program := createTestLoyaltyProgram(t, testDB.DB.DB, business.BusinessID, "Test Program")
	loyalty := createTestClientLoyalty(t, testDB.DB.DB, client.ID, program.ID, business.BusinessID)
	deleterID := uuid.New()

	// Delete client loyalty
	err = clientLoyaltyRepo.Delete(ctx, client.ID, program.ID, deleterID)
	require.NoError(t, err)

	// Verify soft delete
	deleted, err := clientLoyaltyRepo.GetByClientAndProgram(ctx, client.ID, program.ID)
	require.NoError(t, err)
	assert.Nil(t, deleted)

	// Verify in database
	var dbLoyalty models.ClientLoyalty
	err = testDB.DB.DB.Unscoped().Where("id = ?", loyalty.ID).First(&dbLoyalty).Error
	require.NoError(t, err)
	assert.NotNil(t, dbLoyalty.DeletedAt)
	assert.Equal(t, deleterID, *dbLoyalty.DeletedBy)
	assert.False(t, dbLoyalty.IsActive)
	assert.Equal(t, "inactive", dbLoyalty.MembershipStatus)
}

// Helper functions

func createTestUserForLoyalty(t *testing.T, userRepo *repository.UserRepository, email string) *domain.User {
	user := &domain.User{
		ClerkID:   uuid.New(),
		Email:     email,
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+351123456789",
		Role:      "client",
		IsActive:  true,
	}

	err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)
	return user
}

func createTestClientForLoyalty(t *testing.T, clientRepo *repository.ClientRepository, userID, businessID uuid.UUID) *domain.Client {
	client := &domain.Client{
		UserID:     &userID,
		BusinessID: businessID,
		FirstName:  "Test",
		LastName:   "Client",
		Email:      "test_client_" + uuid.New().String()[:8] + "@example.com",
		Phone:      "+351123456789",
	}

	err := clientRepo.Create(context.Background(), client)
	require.NoError(t, err)
	return client
}

func createTestClientLoyalty(t *testing.T, db *gorm.DB, clientID, programID, businessID uuid.UUID) *models.ClientLoyalty {
	return createTestClientLoyaltyWithPoints(t, db, clientID, programID, businessID, 100)
}

func createTestClientLoyaltyWithPoints(t *testing.T, db *gorm.DB, clientID, programID, businessID uuid.UUID, points int) *models.ClientLoyalty {
	loyalty := &models.ClientLoyalty{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
		},
		BusinessID:       businessID,
		ClientID:         clientID,
		ProgramID:        programID,
		Points:           points,
		EnrollmentDate:   time.Now(),
		MembershipStatus: "active",
		IsActive:         true,
		CardNumber:       "CARD-" + uuid.New().String()[:8],
	}

	err := db.Create(loyalty).Error
	require.NoError(t, err)
	return loyalty
}