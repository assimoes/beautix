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
)

func TestStaffRepositoryIntegration_Create(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create a test user
	user := createTestUser(t, testDB.DB)

	// Create a test business
	business := createTestBusiness(t, testDB.DB, user.ID)

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Create a new staff member
	createdBy := user.ID
	ctx := context.Background()
	staff := &domain.Staff{
		BusinessID:      business.ID,
		UserID:          user.ID,
		Position:        "Test Position",
		Bio:             "Test Bio",
		SpecialtyAreas:  []string{"Area 1", "Area 2"},
		ProfileImageURL: "http://example.com/profile.jpg",
		IsActive:        true,
		EmploymentType:  "full-time",
		JoinDate:        time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
		CommissionRate:  20.0,
		CreatedBy:       createdBy,
	}

	// Test creation
	err = staffRepo.Create(ctx, staff)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, staff.StaffID, "Staff ID should be generated")
	assert.NotZero(t, staff.CreatedAt, "Created at timestamp should be set")

	// Verify the staff was created in the database
	var savedStaff models.Staff
	err = testDB.First(&savedStaff, "id = ?", staff.StaffID).Error
	assert.NoError(t, err)
	assert.Equal(t, business.ID, savedStaff.BusinessID)
	assert.Equal(t, user.ID, savedStaff.UserID)
	assert.Equal(t, "Test Position", savedStaff.Position)
	assert.Equal(t, "Test Bio", savedStaff.Bio)
	assert.ElementsMatch(t, []string{"Area 1", "Area 2"}, savedStaff.SpecialtyAreas)
	assert.Equal(t, "http://example.com/profile.jpg", savedStaff.ProfileImageURL)
	assert.True(t, savedStaff.IsActive)
	assert.Equal(t, models.StaffEmploymentType("full-time"), savedStaff.EmploymentType)
	assert.Equal(t, float64(20.0), savedStaff.CommissionRate)
	assert.Equal(t, createdBy, *savedStaff.CreatedBy)
}

func TestStaffRepositoryIntegration_GetByID(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Test GetByID
	ctx := context.Background()
	result, err := staffRepo.GetByID(ctx, staff.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, staff.ID, result.StaffID)
	assert.Equal(t, business.ID, result.BusinessID)
	assert.Equal(t, user.ID, result.UserID)
	assert.Equal(t, staff.Position, result.Position)
	
	// Verify related entities are populated
	assert.NotNil(t, result.User)
	assert.Equal(t, user.ID, result.User.UserID)
	assert.Equal(t, user.Email, result.User.Email)
	
	assert.NotNil(t, result.Business)
	assert.Equal(t, business.ID, result.Business.BusinessID)
}

func TestStaffRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := staffRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffRepositoryIntegration_GetByUserID(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business1 := createTestBusiness(t, testDB.DB, user.ID)
	business2 := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business1.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business2.ID, user.ID)

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Test GetByUserID
	ctx := context.Background()
	results, err := staffRepo.GetByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	
	// Check that both staff members are returned
	staffIDs := []uuid.UUID{results[0].StaffID, results[1].StaffID}
	assert.Contains(t, staffIDs, staff1.ID)
	assert.Contains(t, staffIDs, staff2.ID)
}

func TestStaffRepositoryIntegration_GetByBusinessID(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user1 := createTestUser(t, testDB.DB)
	user2 := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user1.ID)
	staff1 := createTestStaff(t, testDB.DB, business.ID, user1.ID)
	staff2 := createTestStaff(t, testDB.DB, business.ID, user2.ID)

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Test GetByBusinessID
	ctx := context.Background()
	results, err := staffRepo.GetByBusinessID(ctx, business.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	
	// Check that both staff members are returned
	staffIDs := []uuid.UUID{results[0].StaffID, results[1].StaffID}
	assert.Contains(t, staffIDs, staff1.ID)
	assert.Contains(t, staffIDs, staff2.ID)
}

func TestStaffRepositoryIntegration_Update(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Create update input
	ctx := context.Background()
	updatedPosition := "Updated Position"
	updatedBio := "Updated Bio"
	updatedSpecialtyAreas := []string{"Updated Area 1", "Updated Area 2", "New Area"}
	updatedIsActive := false
	updatedCommissionRate := 25.0
	
	updateInput := &domain.UpdateStaffInput{
		Position:       &updatedPosition,
		Bio:            &updatedBio,
		SpecialtyAreas: &updatedSpecialtyAreas,
		IsActive:       &updatedIsActive,
		CommissionRate: &updatedCommissionRate,
	}

	// Test Update
	err = staffRepo.Update(ctx, staff.ID, updateInput, user.ID)
	assert.NoError(t, err)

	// Verify the staff was updated in the database
	var updatedStaff models.Staff
	err = testDB.First(&updatedStaff, "id = ?", staff.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, updatedPosition, updatedStaff.Position)
	assert.Equal(t, updatedBio, updatedStaff.Bio)
	assert.ElementsMatch(t, updatedSpecialtyAreas, updatedStaff.SpecialtyAreas)
	assert.Equal(t, updatedIsActive, updatedStaff.IsActive)
	assert.Equal(t, updatedCommissionRate, updatedStaff.CommissionRate)
	assert.Equal(t, user.ID, *updatedStaff.UpdatedBy)
	assert.NotNil(t, updatedStaff.UpdatedAt)
}

func TestStaffRepositoryIntegration_Delete(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Test Delete
	ctx := context.Background()
	err = staffRepo.Delete(ctx, staff.ID, user.ID)
	assert.NoError(t, err)

	// Verify the staff was soft deleted
	var deletedStaff models.Staff
	err = testDB.Unscoped().First(&deletedStaff, "id = ?", staff.ID).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedStaff.DeletedAt)
	assert.True(t, deletedStaff.DeletedAt.Valid)
	assert.Equal(t, user.ID, *deletedStaff.DeletedBy)
	
	// Verify that the staff is not returned in normal queries
	var count int64
	err = testDB.Model(&models.Staff{}).Where("id = ?", staff.ID).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestStaffRepositoryIntegration_List(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data - 5 staff members
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	for i := 0; i < 5; i++ {
		createTestStaff(t, testDB.DB, business.ID, user.ID)
	}

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Test List with pagination
	ctx := context.Background()
	page1, err := staffRepo.List(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, page1, 2)
	
	page2, err := staffRepo.List(ctx, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, page2, 2)
	
	page3, err := staffRepo.List(ctx, 3, 2)
	assert.NoError(t, err)
	assert.Len(t, page3, 1)
	
	// Make sure each page has different records
	allIDs := make(map[uuid.UUID]bool)
	for _, s := range page1 {
		allIDs[s.StaffID] = true
	}
	for _, s := range page2 {
		allIDs[s.StaffID] = true
	}
	for _, s := range page3 {
		allIDs[s.StaffID] = true
	}
	assert.Len(t, allIDs, 5)
}

func TestStaffRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business1 := createTestBusiness(t, testDB.DB, user.ID)
	business2 := createTestBusiness(t, testDB.DB, user.ID)
	
	// Create 3 staff for business1
	for i := 0; i < 3; i++ {
		createTestStaff(t, testDB.DB, business1.ID, user.ID)
	}
	
	// Create 2 staff for business2
	for i := 0; i < 2; i++ {
		createTestStaff(t, testDB.DB, business2.ID, user.ID)
	}

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Test ListByBusiness
	ctx := context.Background()
	staffBusiness1, err := staffRepo.ListByBusiness(ctx, business1.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, staffBusiness1, 3)
	
	staffBusiness2, err := staffRepo.ListByBusiness(ctx, business2.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, staffBusiness2, 2)
	
	// Verify all staff members are for the correct business
	for _, s := range staffBusiness1 {
		assert.Equal(t, business1.ID, s.BusinessID)
	}
	
	for _, s := range staffBusiness2 {
		assert.Equal(t, business2.ID, s.BusinessID)
	}
}

// Skipping the search test as it requires database-specific search functionality
func TestStaffRepositoryIntegration_Search(t *testing.T) {
	t.Skip("Skipping search test as it requires database-specific functionality")
}

func TestStaffRepositoryIntegration_Count(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business1 := createTestBusiness(t, testDB.DB, user.ID)
	business2 := createTestBusiness(t, testDB.DB, user.ID)
	
	// Create 3 staff for business1
	for i := 0; i < 3; i++ {
		createTestStaff(t, testDB.DB, business1.ID, user.ID)
	}
	
	// Create 2 staff for business2
	for i := 0; i < 2; i++ {
		createTestStaff(t, testDB.DB, business2.ID, user.ID)
	}

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Test Count
	ctx := context.Background()
	count, err := staffRepo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestStaffRepositoryIntegration_CountByBusiness(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business1 := createTestBusiness(t, testDB.DB, user.ID)
	business2 := createTestBusiness(t, testDB.DB, user.ID)
	
	// Create 3 staff for business1
	for i := 0; i < 3; i++ {
		createTestStaff(t, testDB.DB, business1.ID, user.ID)
	}
	
	// Create 2 staff for business2
	for i := 0; i < 2; i++ {
		createTestStaff(t, testDB.DB, business2.ID, user.ID)
	}

	// Create the repository
	staffRepo := repository.NewStaffRepository(testDB.DB)

	// Test CountByBusiness
	ctx := context.Background()
	count1, err := staffRepo.CountByBusiness(ctx, business1.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count1)
	
	count2, err := staffRepo.CountByBusiness(ctx, business2.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count2)
}

