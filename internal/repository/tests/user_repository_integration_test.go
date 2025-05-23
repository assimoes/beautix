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

// TestUserRepositoryIntegration_CreateUser tests creating a user
func TestUserRepositoryIntegration_CreateUser(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	user := &domain.User{
		ClerkID:            uuid.New(),
		Email:              "test_" + uuid.New().String()[:8] + "@example.com",
		FirstName:          "John",
		LastName:           "Doe",
		Phone:              "+1234567890",
		Role:               "user",
		IsActive:           true,
		EmailVerified:      false,
		LanguagePreference: "en",
	}

	// Test creation
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.UserID)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

// TestUserRepositoryIntegration_CreateUserDuplicate tests creating a user with duplicate email
func TestUserRepositoryIntegration_CreateUserDuplicate(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create first user
	user1 := &domain.User{
		ClerkID:   uuid.New(),
		Email:     "duplicate@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "user",
		IsActive:  true,
	}

	err = userRepo.Create(ctx, user1)
	require.NoError(t, err)

	// Try to create second user with same email
	user2 := &domain.User{
		ClerkID:   uuid.New(),
		Email:     "duplicate@example.com",
		FirstName: "Jane",
		LastName:  "Smith",
		Role:      "user",
		IsActive:  true,
	}

	err = userRepo.Create(ctx, user2)
	assert.Error(t, err, "Should fail due to unique constraint on email")
}

// TestUserRepositoryIntegration_GetUserByID tests retrieving a user by ID
func TestUserRepositoryIntegration_GetUserByID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create a user first
	originalUser := &domain.User{
		ClerkID:   uuid.New(),
		Email:     "getbyid_" + uuid.New().String()[:8] + "@example.com",
		FirstName: "Get",
		LastName:  "ByID",
		Phone:     "+1234567890",
		Role:      "user",
		IsActive:  true,
	}

	err = userRepo.Create(ctx, originalUser)
	require.NoError(t, err)

	// Test retrieval
	retrievedUser, err := userRepo.GetByID(ctx, originalUser.UserID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	
	// Verify all fields match
	assert.Equal(t, originalUser.UserID, retrievedUser.UserID)
	assert.Equal(t, originalUser.ClerkID, retrievedUser.ClerkID)
	assert.Equal(t, originalUser.Email, retrievedUser.Email)
	assert.Equal(t, originalUser.FirstName, retrievedUser.FirstName)
	assert.Equal(t, originalUser.LastName, retrievedUser.LastName)
	assert.Equal(t, originalUser.Phone, retrievedUser.Phone)
	assert.Equal(t, originalUser.Role, retrievedUser.Role)
	assert.Equal(t, originalUser.IsActive, retrievedUser.IsActive)

	// Test non-existent ID
	nonExistentID := uuid.New()
	retrievedUser, err = userRepo.GetByID(ctx, nonExistentID)
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
}

// TestUserRepositoryIntegration_GetUserByEmail tests retrieving a user by email
func TestUserRepositoryIntegration_GetUserByEmail(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create a user first
	originalUser := &domain.User{
		ClerkID:   uuid.New(),
		Email:     "getbyemail_" + uuid.New().String()[:8] + "@example.com",
		FirstName: "Get",
		LastName:  "ByEmail",
		Phone:     "+1234567890",
		Role:      "user",
		IsActive:  true,
	}

	err = userRepo.Create(ctx, originalUser)
	require.NoError(t, err)

	// Test retrieval by email
	retrievedUser, err := userRepo.GetByEmail(ctx, originalUser.Email)
	require.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, originalUser.UserID, retrievedUser.UserID)
	assert.Equal(t, originalUser.Email, retrievedUser.Email)

	// Test non-existent email
	retrievedUser, err = userRepo.GetByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
}

// TestUserRepositoryIntegration_UpdateUser tests updating a user
func TestUserRepositoryIntegration_UpdateUser(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		ClerkID:   uuid.New(),
		Email:     "update_" + uuid.New().String()[:8] + "@example.com",
		FirstName: "Original",
		LastName:  "Name",
		Phone:     "+1234567890",
		Role:      "user",
		IsActive:  true,
	}

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Update the user
	firstName := "Updated"
	lastName := "User"
	phone := "+0987654321"
	isActive := false
	
	updateInput := &domain.UpdateUserInput{
		FirstName: &firstName,
		LastName:  &lastName,
		Phone:     &phone,
		IsActive:  &isActive,
	}

	err = userRepo.Update(ctx, user.UserID, updateInput)
	require.NoError(t, err)

	// Retrieve and verify
	updatedUser, err := userRepo.GetByID(ctx, user.UserID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", updatedUser.FirstName)
	assert.Equal(t, "User", updatedUser.LastName)
	assert.Equal(t, "+0987654321", updatedUser.Phone)
	assert.False(t, updatedUser.IsActive)
	assert.True(t, updatedUser.UpdatedAt.After(user.CreatedAt))
}

// TestUserRepositoryIntegration_DeleteUser tests deleting a user
func TestUserRepositoryIntegration_DeleteUser(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		ClerkID:   uuid.New(),
		Email:     "delete_" + uuid.New().String()[:8] + "@example.com",
		FirstName: "Delete",
		LastName:  "Me",
		Phone:     "+1234567890",
		Role:      "user",
		IsActive:  true,
	}

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Delete the user
	err = userRepo.Delete(ctx, user.UserID)
	require.NoError(t, err)

	// Try to retrieve - should fail
	retrievedUser, err := userRepo.GetByID(ctx, user.UserID)
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
}

// TestUserRepositoryIntegration_ListUsers tests listing users with pagination
func TestUserRepositoryIntegration_ListUsers(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create multiple users
	users := []*domain.User{
		{
			ClerkID:   uuid.New(),
			Email:     "list1_" + uuid.New().String()[:8] + "@example.com",
			FirstName: "User",
			LastName:  "One",
			Role:      "user",
			IsActive:  true,
		},
		{
			ClerkID:   uuid.New(),
			Email:     "list2_" + uuid.New().String()[:8] + "@example.com",
			FirstName: "User",
			LastName:  "Two",
			Role:      "admin",
			IsActive:  true,
		},
		{
			ClerkID:   uuid.New(),
			Email:     "list3_" + uuid.New().String()[:8] + "@example.com",
			FirstName: "User",
			LastName:  "Three",
			Role:      "user",
			IsActive:  false,
		},
	}

	for _, user := range users {
		err := userRepo.Create(ctx, user)
		require.NoError(t, err)
	}

	// Test listing with pagination
	page1Users, err := userRepo.List(ctx, 1, 2)
	require.NoError(t, err)
	assert.Len(t, page1Users, 2)

	page2Users, err := userRepo.List(ctx, 2, 2)
	require.NoError(t, err)
	assert.Len(t, page2Users, 1)

	// Test count
	count, err := userRepo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}