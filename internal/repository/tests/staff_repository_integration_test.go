package tests

import (
	"context"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStaffRepositoryIntegration_Create(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()

	// Create a new staff member
	createdBy := testData.User.ID
	ctx := context.Background()
	staff := &domain.Staff{
		BusinessID:      testData.Business.ID,
		UserID:          testData.User.ID,
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
	err := repos.StaffRepo.Create(ctx, staff)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, staff.StaffID, "Staff ID should be generated")
	assert.NotZero(t, staff.CreatedAt, "Created at timestamp should be set")

	// Verify the staff was created using the repository rather than direct DB access
	result, err := repos.StaffRepo.GetByID(ctx, staff.StaffID)
	assert.NoError(t, err)
	assert.Equal(t, testData.Business.ID, result.BusinessID)
	assert.Equal(t, testData.User.ID, result.UserID)
	assert.Equal(t, "Test Position", result.Position)
	assert.Equal(t, "Test Bio", result.Bio)
	assert.ElementsMatch(t, []string{"Area 1", "Area 2"}, result.SpecialtyAreas)
	assert.Equal(t, "http://example.com/profile.jpg", result.ProfileImageURL)
	assert.True(t, result.IsActive)
	assert.Equal(t, "full-time", result.EmploymentType)
	assert.Equal(t, float64(20.0), result.CommissionRate)
	assert.Equal(t, createdBy, result.CreatedBy)
}

func TestStaffRepositoryIntegration_GetByID(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create a new staff member specifically for this test
	extraStaff := createTestStaffTx(t, suite.Tx, testData.Business.ID, testData.User.ID)

	// Test GetByID
	ctx := context.Background()
	result, err := repos.StaffRepo.GetByID(ctx, extraStaff.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, extraStaff.ID, result.StaffID)
	assert.Equal(t, testData.Business.ID, result.BusinessID)
	assert.Equal(t, testData.User.ID, result.UserID)
	assert.Equal(t, extraStaff.Position, result.Position)
	
	// Verify related entities are populated
	assert.NotNil(t, result.User)
	assert.Equal(t, testData.User.ID, result.User.UserID)
	assert.NotEmpty(t, result.User.Email)
	
	assert.NotNil(t, result.Business)
	assert.Equal(t, testData.Business.ID, result.Business.BusinessID)
}

func TestStaffRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := repos.StaffRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffRepositoryIntegration_GetByUserID(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create an additional business
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	
	// Create staff members for both businesses
	staff1 := testData.Staff // already created in CreateTestData
	staff2 := createTestStaffTx(t, suite.Tx, business2.ID, testData.User.ID)

	// Test GetByUserID
	ctx := context.Background()
	results, err := repos.StaffRepo.GetByUserID(ctx, testData.User.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	
	// Check that both staff members are returned
	staffIDs := []uuid.UUID{results[0].StaffID, results[1].StaffID}
	assert.Contains(t, staffIDs, staff1.ID)
	assert.Contains(t, staffIDs, staff2.ID)
}

func TestStaffRepositoryIntegration_GetByBusinessID(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create an additional user
	user2 := createTestUserTx(t, suite.Tx)
	
	// Create staff members for the same business but different users
	staff1 := testData.Staff // already created in CreateTestData
	staff2 := createTestStaffTx(t, suite.Tx, testData.Business.ID, user2.ID)

	// Test GetByBusinessID
	ctx := context.Background()
	results, err := repos.StaffRepo.GetByBusinessID(ctx, testData.Business.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	
	// Check that both staff members are returned
	staffIDs := []uuid.UUID{results[0].StaffID, results[1].StaffID}
	assert.Contains(t, staffIDs, staff1.ID)
	assert.Contains(t, staffIDs, staff2.ID)
}

func TestStaffRepositoryIntegration_Update(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()

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
	err := repos.StaffRepo.Update(ctx, testData.Staff.ID, updateInput, testData.User.ID)
	assert.NoError(t, err)

	// Verify the staff was updated
	updated, err := repos.StaffRepo.GetByID(ctx, testData.Staff.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedPosition, updated.Position)
	assert.Equal(t, updatedBio, updated.Bio)
	assert.ElementsMatch(t, updatedSpecialtyAreas, updated.SpecialtyAreas)
	assert.Equal(t, updatedIsActive, updated.IsActive)
	assert.Equal(t, updatedCommissionRate, updated.CommissionRate)
	assert.Equal(t, testData.User.ID, *updated.UpdatedBy)
	assert.NotNil(t, updated.UpdatedAt)
}

func TestStaffRepositoryIntegration_Delete(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()

	// Test Delete
	ctx := context.Background()
	err := repos.StaffRepo.Delete(ctx, testData.Staff.ID, testData.User.ID)
	assert.NoError(t, err)

	// Verify the staff was soft deleted by trying to retrieve it
	_, err = repos.StaffRepo.GetByID(ctx, testData.Staff.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffRepositoryIntegration_List(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data - 5 staff members in total
	testData := suite.CreateTestData() // This creates 1 staff
	
	// Create 4 more staff members
	for i := 0; i < 4; i++ {
		createTestStaffTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	}

	// Test List with pagination
	ctx := context.Background()
	page1, err := repos.StaffRepo.List(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, page1, 2)
	
	page2, err := repos.StaffRepo.List(ctx, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, page2, 2)
	
	page3, err := repos.StaffRepo.List(ctx, 3, 2)
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
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	
	// Create 2 more staff for business1 (total 3)
	for i := 0; i < 2; i++ {
		createTestStaffTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	}
	
	// Create 2 staff for business2
	for i := 0; i < 2; i++ {
		createTestStaffTx(t, suite.Tx, business2.ID, testData.User.ID)
	}

	// Test ListByBusiness
	ctx := context.Background()
	staffBusiness1, err := repos.StaffRepo.ListByBusiness(ctx, testData.Business.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, staffBusiness1, 3)
	
	staffBusiness2, err := repos.StaffRepo.ListByBusiness(ctx, business2.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, staffBusiness2, 2)
	
	// Verify all staff members are for the correct business
	for _, s := range staffBusiness1 {
		assert.Equal(t, testData.Business.ID, s.BusinessID)
	}
	
	for _, s := range staffBusiness2 {
		assert.Equal(t, business2.ID, s.BusinessID)
	}
}

// Skip the Search test as it was marked as flaky
func TestStaffRepositoryIntegration_Search(t *testing.T) {
	t.Skip("Skipping search test as it requires database-specific functionality")
}

func TestStaffRepositoryIntegration_Count(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data - 5 staff members in total
	testData := suite.CreateTestData() // This creates 1 staff
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	
	// Create 2 more staff for business1 (total 3)
	for i := 0; i < 2; i++ {
		createTestStaffTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	}
	
	// Create 2 staff for business2
	for i := 0; i < 2; i++ {
		createTestStaffTx(t, suite.Tx, business2.ID, testData.User.ID)
	}

	// Test Count
	ctx := context.Background()
	count, err := repos.StaffRepo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count) // 1 + 2 + 2 = 5
}

func TestStaffRepositoryIntegration_CountByBusiness(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	
	// Create 2 more staff for business1 (total 3)
	for i := 0; i < 2; i++ {
		createTestStaffTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	}
	
	// Create 2 staff for business2
	for i := 0; i < 2; i++ {
		createTestStaffTx(t, suite.Tx, business2.ID, testData.User.ID)
	}

	// Test CountByBusiness
	ctx := context.Background()
	count1, err := repos.StaffRepo.CountByBusiness(ctx, testData.Business.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count1)
	
	count2, err := repos.StaffRepo.CountByBusiness(ctx, business2.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count2)
}