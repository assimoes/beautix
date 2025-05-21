package tests

import (
	"context"
	"testing"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServiceAssignmentRepositoryIntegration_Create(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	service := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)

	// Create a new service assignment
	createdBy := testData.User.ID
	ctx := context.Background()
	assignment := &domain.ServiceAssignment{
		BusinessID: testData.Business.ID,
		StaffID:    testData.Staff.ID,
		ServiceID:  service.ID,
		IsActive:   true,
		CreatedBy:  createdBy,
	}

	// Test creation
	err := repos.ServiceAssignmentRepo.Create(ctx, assignment)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, assignment.AssignmentID, "Assignment ID should be generated")
	assert.NotZero(t, assignment.CreatedAt, "Created at timestamp should be set")

	// Verify the assignment was created
	result, err := repos.ServiceAssignmentRepo.GetByID(ctx, assignment.AssignmentID)
	assert.NoError(t, err)
	assert.Equal(t, testData.Business.ID, result.BusinessID)
	assert.Equal(t, testData.Staff.ID, result.StaffID)
	assert.Equal(t, service.ID, result.ServiceID)
	assert.True(t, result.IsActive)
	assert.Equal(t, createdBy, result.CreatedBy)
}

func TestServiceAssignmentRepositoryIntegration_GetByID(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	service := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	assignment := createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, service.ID, testData.User.ID)

	// Test GetByID
	ctx := context.Background()
	result, err := repos.ServiceAssignmentRepo.GetByID(ctx, assignment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, assignment.ID, result.AssignmentID)
	assert.Equal(t, testData.Business.ID, result.BusinessID)
	assert.Equal(t, testData.Staff.ID, result.StaffID)
	assert.Equal(t, service.ID, result.ServiceID)
	assert.True(t, result.IsActive)

	// Verify related entities are populated
	assert.NotNil(t, result.Staff)
	assert.Equal(t, testData.Staff.ID, result.Staff.StaffID)
}

func TestServiceAssignmentRepositoryIntegration_GetByStaffAndService(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	service := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	assignment := createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, service.ID, testData.User.ID)

	// Test GetByStaffAndService
	ctx := context.Background()
	result, err := repos.ServiceAssignmentRepo.GetByStaffAndService(ctx, testData.Staff.ID, service.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, assignment.ID, result.AssignmentID)
	assert.Equal(t, testData.Staff.ID, result.StaffID)
	assert.Equal(t, service.ID, result.ServiceID)
}

func TestServiceAssignmentRepositoryIntegration_GetByStaff(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	service1 := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	service2 := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	assignment1 := createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, service1.ID, testData.User.ID)
	assignment2 := createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, service2.ID, testData.User.ID)

	// Test GetByStaff
	ctx := context.Background()
	results, err := repos.ServiceAssignmentRepo.GetByStaff(ctx, testData.Staff.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Check that both assignments are returned
	assignmentIDs := []uuid.UUID{results[0].AssignmentID, results[1].AssignmentID}
	assert.Contains(t, assignmentIDs, assignment1.ID)
	assert.Contains(t, assignmentIDs, assignment2.ID)

	// Check that all assignments are for the correct staff
	for _, a := range results {
		assert.Equal(t, testData.Staff.ID, a.StaffID)
	}
}

func TestServiceAssignmentRepositoryIntegration_GetByService(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	staff2 := createTestStaffTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	service := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	assignment1 := createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, service.ID, testData.User.ID)
	assignment2 := createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, staff2.ID, service.ID, testData.User.ID)

	// Test GetByService
	ctx := context.Background()
	results, err := repos.ServiceAssignmentRepo.GetByService(ctx, service.ID)
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
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	service := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	assignment := createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, service.ID, testData.User.ID)

	// Create update input
	ctx := context.Background()
	isActive := false
	updateInput := &domain.UpdateServiceAssignmentInput{
		IsActive: &isActive,
	}

	// Test Update
	err := repos.ServiceAssignmentRepo.Update(ctx, assignment.ID, updateInput, testData.User.ID)
	assert.NoError(t, err)

	// Verify the assignment was updated
	result, err := repos.ServiceAssignmentRepo.GetByID(ctx, assignment.ID)
	assert.NoError(t, err)
	assert.Equal(t, isActive, result.IsActive)
	assert.Equal(t, testData.User.ID, *result.UpdatedBy)
	assert.NotNil(t, result.UpdatedAt)
}

func TestServiceAssignmentRepositoryIntegration_Delete(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	service := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	assignment := createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, service.ID, testData.User.ID)

	// Test Delete
	ctx := context.Background()
	err := repos.ServiceAssignmentRepo.Delete(ctx, assignment.ID, testData.User.ID)
	assert.NoError(t, err)

	// Verify the assignment was soft deleted
	_, err = repos.ServiceAssignmentRepo.GetByID(ctx, assignment.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestServiceAssignmentRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	staff2 := createTestStaffTx(t, suite.Tx, business2.ID, testData.User.ID)

	// Create 3 assignments for business1
	for i := 0; i < 3; i++ {
		s := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
		createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, s.ID, testData.User.ID)
	}

	// Create 2 assignments for business2
	for i := 0; i < 2; i++ {
		s := createTestServiceTx(t, suite.Tx, business2.ID, testData.User.ID)
		createTestServiceAssignmentTx(t, suite.Tx, business2.ID, staff2.ID, s.ID, testData.User.ID)
	}

	// Test ListByBusiness
	ctx := context.Background()
	assignmentsBusiness1, err := repos.ServiceAssignmentRepo.ListByBusiness(ctx, testData.Business.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, assignmentsBusiness1, 3)

	assignmentsBusiness2, err := repos.ServiceAssignmentRepo.ListByBusiness(ctx, business2.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, assignmentsBusiness2, 2)

	// Verify all assignments are for the correct business
	for _, a := range assignmentsBusiness1 {
		assert.Equal(t, testData.Business.ID, a.BusinessID)
	}

	for _, a := range assignmentsBusiness2 {
		assert.Equal(t, business2.ID, a.BusinessID)
	}
}

func TestServiceAssignmentRepositoryIntegration_CountByBusiness(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	staff2 := createTestStaffTx(t, suite.Tx, business2.ID, testData.User.ID)

	// Create 3 assignments for business1
	for i := 0; i < 3; i++ {
		s := createTestServiceTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
		createTestServiceAssignmentTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, s.ID, testData.User.ID)
	}

	// Create 2 assignments for business2
	for i := 0; i < 2; i++ {
		s := createTestServiceTx(t, suite.Tx, business2.ID, testData.User.ID)
		createTestServiceAssignmentTx(t, suite.Tx, business2.ID, staff2.ID, s.ID, testData.User.ID)
	}

	// Test CountByBusiness
	ctx := context.Background()
	count1, err := repos.ServiceAssignmentRepo.CountByBusiness(ctx, testData.Business.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count1)

	count2, err := repos.ServiceAssignmentRepo.CountByBusiness(ctx, business2.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count2)
}