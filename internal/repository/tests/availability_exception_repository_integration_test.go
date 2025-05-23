package tests

import (
	"context"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvailabilityExceptionRepositoryIntegration_Create(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	availabilityExceptionRepo := &repository.AvailabilityExceptionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "availability_create@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "create-availability-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "create_avail_staff@test.com")

	// Create a new availability exception
	startTime := time.Now().Add(24 * time.Hour) // Tomorrow
	endTime := startTime.Add(8 * time.Hour)     // 8 hours duration

	exception := &domain.AvailabilityException{
		BusinessID:     business.BusinessID,
		StaffID:        staff.StaffID,
		ExceptionType:  "time_off",
		StartTime:      startTime,
		EndTime:        endTime,
		IsFullDay:      true,
		IsRecurring:    false,
		RecurrenceRule: "",
		Notes:          "Annual leave",
		CreatedBy:      user.UserID,
	}

	// Test creation
	err = availabilityExceptionRepo.Create(ctx, exception)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, exception.ExceptionID)
	assert.NotZero(t, exception.CreatedAt)

	// Verify the exception was created
	result, err := availabilityExceptionRepo.GetByID(ctx, exception.ExceptionID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, staff.StaffID, result.StaffID)
	assert.Equal(t, "time_off", result.ExceptionType)
	assert.Equal(t, "Annual leave", result.Notes)
	assert.Equal(t, user.UserID, result.CreatedBy)
	assert.WithinDuration(t, startTime, result.StartTime, time.Second)
	assert.WithinDuration(t, endTime, result.EndTime, time.Second)
}

func TestAvailabilityExceptionRepositoryIntegration_GetByID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	availabilityExceptionRepo := &repository.AvailabilityExceptionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "availability_getbyid@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbyid-availability-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "getbyid_avail_staff@test.com")
	exception := createTestAvailabilityExceptionForIntegration(t, availabilityExceptionRepo, business.BusinessID, staff.StaffID, user.UserID)

	// Test GetByID
	result, err := availabilityExceptionRepo.GetByID(ctx, exception.ExceptionID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, exception.ExceptionID, result.ExceptionID)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, staff.StaffID, result.StaffID)
	assert.Equal(t, exception.ExceptionType, result.ExceptionType)
}

func TestAvailabilityExceptionRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	availabilityExceptionRepo := &repository.AvailabilityExceptionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := availabilityExceptionRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestAvailabilityExceptionRepositoryIntegration_Update(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	availabilityExceptionRepo := &repository.AvailabilityExceptionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "availability_update@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "update-availability-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "update_avail_staff@test.com")
	exception := createTestAvailabilityExceptionForIntegration(t, availabilityExceptionRepo, business.BusinessID, staff.StaffID, user.UserID)

	// Create update input
	updatedNotes := "Updated leave notes"
	newStartTime := time.Now().Add(48 * time.Hour)
	newEndTime := newStartTime.Add(4 * time.Hour)

	updateInput := &domain.UpdateAvailabilityExceptionInput{
		Notes:     &updatedNotes,
		StartTime: &newStartTime,
		EndTime:   &newEndTime,
	}

	// Test update
	err = availabilityExceptionRepo.Update(ctx, exception.ExceptionID, updateInput, user.UserID)
	require.NoError(t, err)

	// Verify the update
	updatedException, err := availabilityExceptionRepo.GetByID(ctx, exception.ExceptionID)
	require.NoError(t, err)
	assert.Equal(t, updatedNotes, updatedException.Notes)
	assert.WithinDuration(t, newStartTime, updatedException.StartTime, time.Second)
	assert.WithinDuration(t, newEndTime, updatedException.EndTime, time.Second)
}

func TestAvailabilityExceptionRepositoryIntegration_Delete(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	availabilityExceptionRepo := &repository.AvailabilityExceptionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "availability_delete@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "delete-availability-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "delete_avail_staff@test.com")
	exception := createTestAvailabilityExceptionForIntegration(t, availabilityExceptionRepo, business.BusinessID, staff.StaffID, user.UserID)

	// Test delete
	err = availabilityExceptionRepo.Delete(ctx, exception.ExceptionID, user.UserID)
	require.NoError(t, err)

	// Verify the exception is deleted (not found)
	deletedException, err := availabilityExceptionRepo.GetByID(ctx, exception.ExceptionID)
	assert.Error(t, err)
	assert.Nil(t, deletedException)
	assert.Contains(t, err.Error(), "not found")
}

func TestAvailabilityExceptionRepositoryIntegration_ListByStaff(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	availabilityExceptionRepo := &repository.AvailabilityExceptionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "availability_listbystaff@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "listbystaff-availability-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "listbystaff_avail_staff@test.com")

	// Create multiple exceptions for the staff
	exception1 := createTestAvailabilityExceptionForIntegration(t, availabilityExceptionRepo, business.BusinessID, staff.StaffID, user.UserID)
	exception2 := createTestAvailabilityExceptionForIntegration(t, availabilityExceptionRepo, business.BusinessID, staff.StaffID, user.UserID)

	// Test GetByStaff
	results, err := availabilityExceptionRepo.GetByStaff(ctx, staff.StaffID)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify all exceptions belong to the staff
	for _, exception := range results {
		assert.Equal(t, staff.StaffID, exception.StaffID)
	}

	// Verify we got our exceptions
	exceptionIDs := make([]uuid.UUID, len(results))
	for i, exception := range results {
		exceptionIDs[i] = exception.ExceptionID
	}
	assert.Contains(t, exceptionIDs, exception1.ExceptionID)
	assert.Contains(t, exceptionIDs, exception2.ExceptionID)
}

func TestAvailabilityExceptionRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	availabilityExceptionRepo := &repository.AvailabilityExceptionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "availability_listbybusiness@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "listbybusiness-availability-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "listbybusiness_avail_staff@test.com")

	// Create multiple exceptions for the business
	exception1 := createTestAvailabilityExceptionForIntegration(t, availabilityExceptionRepo, business.BusinessID, staff.StaffID, user.UserID)
	exception2 := createTestAvailabilityExceptionForIntegration(t, availabilityExceptionRepo, business.BusinessID, staff.StaffID, user.UserID)

	// Test ListByBusiness
	results, err := availabilityExceptionRepo.ListByBusiness(ctx, business.BusinessID, 1, 10)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify all exceptions belong to the business
	for _, exception := range results {
		assert.Equal(t, business.BusinessID, exception.BusinessID)
	}

	// Verify we got our exceptions
	exceptionIDs := make([]uuid.UUID, len(results))
	for i, exception := range results {
		exceptionIDs[i] = exception.ExceptionID
	}
	assert.Contains(t, exceptionIDs, exception1.ExceptionID)
	assert.Contains(t, exceptionIDs, exception2.ExceptionID)
}

// Helper function to create a test availability exception with incrementing times
var availabilityExceptionCounter int = 0
func createTestAvailabilityExceptionForIntegration(t *testing.T, availabilityExceptionRepo *repository.AvailabilityExceptionRepository, businessID, staffID, userID uuid.UUID) *domain.AvailabilityException {
	availabilityExceptionCounter++
	startTime := time.Now().Add(time.Duration(24+availabilityExceptionCounter*12) * time.Hour)
	endTime := startTime.Add(8 * time.Hour)
	
	exception := &domain.AvailabilityException{
		BusinessID:     businessID,
		StaffID:        staffID,
		ExceptionType:  "time_off",
		StartTime:      startTime,
		EndTime:        endTime,
		IsFullDay:      true,
		IsRecurring:    false,
		RecurrenceRule: "",
		Notes:          "Test availability exception",
		CreatedBy:      userID,
	}
	
	err := availabilityExceptionRepo.Create(context.Background(), exception)
	require.NoError(t, err)
	
	return exception
}