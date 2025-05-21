package tests

import (
	"context"
	"testing"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/assimoes/beautix/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceAssignmentRepositoryIntegration_Create(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	service := createTestService(t, testDB.DB, business.ID, user.ID)

	// Create the repository
	assignmentRepo := repository.NewServiceAssignmentRepository(testDB.DB)

	// Create a new service assignment
	createdBy := user.ID
	ctx := context.Background()
	assignment := &domain.ServiceAssignment{
		BusinessID: business.ID,
		StaffID:    staff.ID,
		ServiceID:  service.ID,
		IsActive:   true,
		CreatedBy:  createdBy,
	}

	// Test creation
	err = assignmentRepo.Create(ctx, assignment)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, assignment.AssignmentID, "Assignment ID should be generated")
	assert.NotZero(t, assignment.CreatedAt, "Created at timestamp should be set")

	// Verify the assignment was created in the database
	var savedAssignment models.ServiceAssignment
	err = testDB.First(&savedAssignment, "id = ?", assignment.AssignmentID).Error
	assert.NoError(t, err)
	assert.Equal(t, business.ID, savedAssignment.BusinessID)
	assert.Equal(t, staff.ID, savedAssignment.StaffID)
	assert.Equal(t, service.ID, savedAssignment.ServiceID)
	assert.True(t, savedAssignment.IsActive)
	assert.Equal(t, createdBy, *savedAssignment.CreatedBy)
}

func TestServiceAssignmentRepositoryIntegration_GetByID(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	service := createTestService(t, testDB.DB, business.ID, user.ID)
	assignment := createTestServiceAssignment(t, testDB.DB, business.ID, staff.ID, service.ID, user.ID)

	// Create the repository
	assignmentRepo := repository.NewServiceAssignmentRepository(testDB.DB)

	// Test GetByID
	ctx := context.Background()
	result, err := assignmentRepo.GetByID(ctx, assignment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, assignment.ID, result.AssignmentID)
	assert.Equal(t, business.ID, result.BusinessID)
	assert.Equal(t, staff.ID, result.StaffID)
	assert.Equal(t, service.ID, result.ServiceID)
	assert.True(t, result.IsActive)

	// Verify related entities are populated
	assert.NotNil(t, result.Staff)
	assert.Equal(t, staff.ID, result.Staff.StaffID)
}

func TestServiceAssignmentRepositoryIntegration_GetByStaffAndService(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	service := createTestService(t, testDB.DB, business.ID, user.ID)
	assignment := createTestServiceAssignment(t, testDB.DB, business.ID, staff.ID, service.ID, user.ID)

	// Create the repository
	assignmentRepo := repository.NewServiceAssignmentRepository(testDB.DB)

	// Test GetByStaffAndService
	ctx := context.Background()
	result, err := assignmentRepo.GetByStaffAndService(ctx, staff.ID, service.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, assignment.ID, result.AssignmentID)
	assert.Equal(t, staff.ID, result.StaffID)
	assert.Equal(t, service.ID, result.ServiceID)
}

func TestServiceAssignmentRepositoryIntegration_GetByStaff(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	service1 := createTestService(t, testDB.DB, business.ID, user.ID)
	service2 := createTestService(t, testDB.DB, business.ID, user.ID)
	assignment1 := createTestServiceAssignment(t, testDB.DB, business.ID, staff.ID, service1.ID, user.ID)
	assignment2 := createTestServiceAssignment(t, testDB.DB, business.ID, staff.ID, service2.ID, user.ID)

	// Create the repository
	assignmentRepo := repository.NewServiceAssignmentRepository(testDB.DB)

	// Test GetByStaff
	ctx := context.Background()
	results, err := assignmentRepo.GetByStaff(ctx, staff.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Check that both assignments are returned
	assignmentIDs := []uuid.UUID{results[0].AssignmentID, results[1].AssignmentID}
	assert.Contains(t, assignmentIDs, assignment1.ID)
	assert.Contains(t, assignmentIDs, assignment2.ID)

	// Check that all assignments are for the correct staff
	for _, a := range results {
		assert.Equal(t, staff.ID, a.StaffID)
	}
}

func TestServiceAssignmentRepositoryIntegration_GetByService(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business.ID, user.ID)
	service := createTestService(t, testDB.DB, business.ID, user.ID)
	assignment1 := createTestServiceAssignment(t, testDB.DB, business.ID, staff1.ID, service.ID, user.ID)
	assignment2 := createTestServiceAssignment(t, testDB.DB, business.ID, staff2.ID, service.ID, user.ID)

	// Create the repository
	assignmentRepo := repository.NewServiceAssignmentRepository(testDB.DB)

	// Test GetByService
	ctx := context.Background()
	results, err := assignmentRepo.GetByService(ctx, service.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Check that both assignments are returned
	assignmentIDs := []uuid.UUID{results[0].AssignmentID, results[1].AssignmentID}
	assert.Contains(t, assignmentIDs, assignment1.ID)
	assert.Contains(t, assignmentIDs, assignment2.ID)

	// Check that all assignments are for the correct service
	for _, a := range results {
		assert.Equal(t, service.ID, a.ServiceID)
	}
}

func TestServiceAssignmentRepositoryIntegration_Update(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	service := createTestService(t, testDB.DB, business.ID, user.ID)
	assignment := createTestServiceAssignment(t, testDB.DB, business.ID, staff.ID, service.ID, user.ID)

	// Create the repository
	assignmentRepo := repository.NewServiceAssignmentRepository(testDB.DB)

	// Create update input
	ctx := context.Background()
	isActive := false
	updateInput := &domain.UpdateServiceAssignmentInput{
		IsActive: &isActive,
	}

	// Test Update
	err = assignmentRepo.Update(ctx, assignment.ID, updateInput, user.ID)
	assert.NoError(t, err)

	// Verify the assignment was updated in the database
	var updatedAssignment models.ServiceAssignment
	err = testDB.First(&updatedAssignment, "id = ?", assignment.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, isActive, updatedAssignment.IsActive)
	assert.Equal(t, user.ID, *updatedAssignment.UpdatedBy)
	assert.NotNil(t, updatedAssignment.UpdatedAt)
}

func TestServiceAssignmentRepositoryIntegration_Delete(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	service := createTestService(t, testDB.DB, business.ID, user.ID)
	assignment := createTestServiceAssignment(t, testDB.DB, business.ID, staff.ID, service.ID, user.ID)

	// Create the repository
	assignmentRepo := repository.NewServiceAssignmentRepository(testDB.DB)

	// Test Delete
	ctx := context.Background()
	err = assignmentRepo.Delete(ctx, assignment.ID, user.ID)
	assert.NoError(t, err)

	// Verify the assignment was soft deleted
	var deletedAssignment models.ServiceAssignment
	err = testDB.Unscoped().First(&deletedAssignment, "id = ?", assignment.ID).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedAssignment.DeletedAt)
	assert.True(t, deletedAssignment.DeletedAt.Valid)
	assert.Equal(t, user.ID, *deletedAssignment.DeletedBy)

	// Verify that the assignment is not returned in normal queries
	var count int64
	err = testDB.Model(&models.ServiceAssignment{}).Where("id = ?", assignment.ID).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestServiceAssignmentRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business1 := createTestBusiness(t, testDB.DB, user.ID)
	business2 := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business1.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business2.ID, user.ID)
	// Create services for each business
	createTestService(t, testDB.DB, business1.ID, user.ID)
	createTestService(t, testDB.DB, business2.ID, user.ID)

	// Create 3 assignments for business1
	for i := 0; i < 3; i++ {
		s := createTestService(t, testDB.DB, business1.ID, user.ID)
		createTestServiceAssignment(t, testDB.DB, business1.ID, staff1.ID, s.ID, user.ID)
	}

	// Create 2 assignments for business2
	for i := 0; i < 2; i++ {
		s := createTestService(t, testDB.DB, business2.ID, user.ID)
		createTestServiceAssignment(t, testDB.DB, business2.ID, staff2.ID, s.ID, user.ID)
	}

	// Create the repository
	assignmentRepo := repository.NewServiceAssignmentRepository(testDB.DB)

	// Test ListByBusiness
	ctx := context.Background()
	assignmentsBusiness1, err := assignmentRepo.ListByBusiness(ctx, business1.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, assignmentsBusiness1, 3)

	assignmentsBusiness2, err := assignmentRepo.ListByBusiness(ctx, business2.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, assignmentsBusiness2, 2)

	// Verify all assignments are for the correct business
	for _, a := range assignmentsBusiness1 {
		assert.Equal(t, business1.ID, a.BusinessID)
	}

	for _, a := range assignmentsBusiness2 {
		assert.Equal(t, business2.ID, a.BusinessID)
	}
}

func TestServiceAssignmentRepositoryIntegration_CountByBusiness(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business1 := createTestBusiness(t, testDB.DB, user.ID)
	business2 := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business1.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business2.ID, user.ID)

	// Create 3 assignments for business1
	for i := 0; i < 3; i++ {
		s := createTestService(t, testDB.DB, business1.ID, user.ID)
		createTestServiceAssignment(t, testDB.DB, business1.ID, staff1.ID, s.ID, user.ID)
	}

	// Create 2 assignments for business2
	for i := 0; i < 2; i++ {
		s := createTestService(t, testDB.DB, business2.ID, user.ID)
		createTestServiceAssignment(t, testDB.DB, business2.ID, staff2.ID, s.ID, user.ID)
	}

	// Create the repository
	assignmentRepo := repository.NewServiceAssignmentRepository(testDB.DB)

	// Test CountByBusiness
	ctx := context.Background()
	count1, err := assignmentRepo.CountByBusiness(ctx, business1.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count1)

	count2, err := assignmentRepo.CountByBusiness(ctx, business2.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count2)
}

