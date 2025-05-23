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

func TestServiceCategoryRepositoryIntegration_Create(t *testing.T) {
	// Use simple test database
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)

	// Clean up all tables comprehensively to avoid foreign key issues
	database.CleanupAllTables(t, testDB.DB)

	// Create repository
	repo := &repository.ServiceCategoryRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}

	// Create minimal test business first (required for foreign key)
	businessID := uuid.New()
	userID := uuid.New()
	
	// Create user first
	err = testDB.Exec(`INSERT INTO users (id, clerk_id, email, first_name, last_name, phone, role, is_active) 
						VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, "test_clerk", "test@example.com", "Test", "User", "+1234567890", "owner", true).Error
	require.NoError(t, err)
	
	// Create business
	err = testDB.Exec(`INSERT INTO businesses (id, user_id, name, _business_type, display_name, address, city, country, phone, email, subscription_tier, is_active, is_verified, time_zone, currency, settings) 
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		businessID, userID, "test-business", "", "Test Business", "123 Test St", "Test City", "Test Country", "+9876543210", "test@business.com", "basic", true, false, "Europe/Lisbon", "EUR", "{}").Error
	require.NoError(t, err)

	category := &domain.ServiceCategory{
		BusinessID:  businessID,
		Name:        "Hair Services",
		Description: "All hair-related services",
	}

	err = repo.Create(context.Background(), category)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, category.ID)
	assert.NotZero(t, category.CreatedAt)
}

// TODO: Convert remaining tests to simple approach
// The other test functions were using transaction-based testing
// which caused deadlocks and hangs. They need to be rewritten
// to use the simple database.NewSimpleTestDB approach.