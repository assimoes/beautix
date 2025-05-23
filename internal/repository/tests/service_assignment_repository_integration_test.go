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

func TestServiceAssignmentRepositoryIntegration_Create(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceAssignmentRepo := &repository.ServiceAssignmentRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_assignment_create@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "create-service-assignment-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "create_svc_assign_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")

	// Create a new service assignment
	assignment := &domain.ServiceAssignment{
		BusinessID: business.BusinessID,
		StaffID:    staff.StaffID,
		ServiceID:  service.ID,
		IsActive:   true,
		CreatedBy:  user.UserID,
	}

	// Test creation
	err = serviceAssignmentRepo.Create(ctx, assignment)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, assignment.AssignmentID)
	assert.NotZero(t, assignment.CreatedAt)

	// Verify the assignment was created
	result, err := serviceAssignmentRepo.GetByID(ctx, assignment.AssignmentID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, staff.StaffID, result.StaffID)
	assert.Equal(t, service.ID, result.ServiceID)
	assert.True(t, result.IsActive)
	assert.Equal(t, user.UserID, result.CreatedBy)
}

func TestServiceAssignmentRepositoryIntegration_GetByID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceAssignmentRepo := &repository.ServiceAssignmentRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_assignment_getbyid@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbyid-service-assignment-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "getbyid_svc_assign_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")
	assignment := createTestServiceAssignmentForIntegration(t, serviceAssignmentRepo, business.BusinessID, staff.StaffID, service.ID, user.UserID)

	// Test GetByID
	result, err := serviceAssignmentRepo.GetByID(ctx, assignment.AssignmentID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, assignment.AssignmentID, result.AssignmentID)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, staff.StaffID, result.StaffID)
	assert.Equal(t, service.ID, result.ServiceID)
}

func TestServiceAssignmentRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	serviceAssignmentRepo := &repository.ServiceAssignmentRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := serviceAssignmentRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestServiceAssignmentRepositoryIntegration_Update(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceAssignmentRepo := &repository.ServiceAssignmentRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_assignment_update@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "update-service-assignment-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "update_svc_assign_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")
	assignment := createTestServiceAssignmentForIntegration(t, serviceAssignmentRepo, business.BusinessID, staff.StaffID, service.ID, user.UserID)

	// Create update input
	isActive := false
	updateInput := &domain.UpdateServiceAssignmentInput{
		IsActive: &isActive,
	}

	// Test update
	err = serviceAssignmentRepo.Update(ctx, assignment.AssignmentID, updateInput, user.UserID)
	require.NoError(t, err)

	// Verify the update
	updatedAssignment, err := serviceAssignmentRepo.GetByID(ctx, assignment.AssignmentID)
	require.NoError(t, err)
	assert.False(t, updatedAssignment.IsActive)
}

func TestServiceAssignmentRepositoryIntegration_Delete(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceAssignmentRepo := &repository.ServiceAssignmentRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_assignment_delete@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "delete-service-assignment-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "delete_svc_assign_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")
	assignment := createTestServiceAssignmentForIntegration(t, serviceAssignmentRepo, business.BusinessID, staff.StaffID, service.ID, user.UserID)

	// Test delete
	err = serviceAssignmentRepo.Delete(ctx, assignment.AssignmentID, user.UserID)
	require.NoError(t, err)

	// Verify the assignment is deleted (not found)
	deletedAssignment, err := serviceAssignmentRepo.GetByID(ctx, assignment.AssignmentID)
	assert.Error(t, err)
	assert.Nil(t, deletedAssignment)
	assert.Contains(t, err.Error(), "not found")
}

func TestServiceAssignmentRepositoryIntegration_GetByStaff(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceAssignmentRepo := &repository.ServiceAssignmentRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_assignment_getbystaff@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbystaff-service-assignment-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "getbystaff_svc_assign_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service1 := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")
	service2 := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Color")

	// Create multiple assignments for the staff
	assignment1 := createTestServiceAssignmentForIntegration(t, serviceAssignmentRepo, business.BusinessID, staff.StaffID, service1.ID, user.UserID)
	assignment2 := createTestServiceAssignmentForIntegration(t, serviceAssignmentRepo, business.BusinessID, staff.StaffID, service2.ID, user.UserID)

	// Test GetByStaff
	results, err := serviceAssignmentRepo.GetByStaff(ctx, staff.StaffID)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify all assignments belong to the staff
	for _, assignment := range results {
		assert.Equal(t, staff.StaffID, assignment.StaffID)
	}

	// Verify we got our assignments
	assignmentIDs := make([]uuid.UUID, len(results))
	for i, assignment := range results {
		assignmentIDs[i] = assignment.AssignmentID
	}
	assert.Contains(t, assignmentIDs, assignment1.AssignmentID)
	assert.Contains(t, assignmentIDs, assignment2.AssignmentID)
}

func TestServiceAssignmentRepositoryIntegration_GetByService(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceAssignmentRepo := &repository.ServiceAssignmentRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_assignment_getbyservice@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbyservice-service-assignment-salon")
	staff1 := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "getbyservice_svc_assign_staff1@test.com")
	staff2 := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "getbyservice_svc_assign_staff2@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")

	// Create multiple assignments for the service
	assignment1 := createTestServiceAssignmentForIntegration(t, serviceAssignmentRepo, business.BusinessID, staff1.StaffID, service.ID, user.UserID)
	assignment2 := createTestServiceAssignmentForIntegration(t, serviceAssignmentRepo, business.BusinessID, staff2.StaffID, service.ID, user.UserID)

	// Test GetByService
	results, err := serviceAssignmentRepo.GetByService(ctx, service.ID)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify all assignments belong to the service
	for _, assignment := range results {
		assert.Equal(t, service.ID, assignment.ServiceID)
	}

	// Verify we got our assignments
	assignmentIDs := make([]uuid.UUID, len(results))
	for i, assignment := range results {
		assignmentIDs[i] = assignment.AssignmentID
	}
	assert.Contains(t, assignmentIDs, assignment1.AssignmentID)
	assert.Contains(t, assignmentIDs, assignment2.AssignmentID)
}

// Helper function to create a test service assignment
func createTestServiceAssignmentForIntegration(t *testing.T, serviceAssignmentRepo *repository.ServiceAssignmentRepository, businessID, staffID, serviceID, userID uuid.UUID) *domain.ServiceAssignment {
	assignment := &domain.ServiceAssignment{
		BusinessID: businessID,
		StaffID:    staffID,
		ServiceID:  serviceID,
		IsActive:   true,
		CreatedBy:  userID,
	}
	
	err := serviceAssignmentRepo.Create(context.Background(), assignment)
	require.NoError(t, err)
	
	return assignment
}