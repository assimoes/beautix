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

func TestAvailabilityExceptionRepositoryIntegration_Create(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Create a new availability exception
	createdBy := user.ID
	ctx := context.Background()
	startTime := time.Now().Add(24 * time.Hour) // Tomorrow
	endTime := startTime.Add(8 * time.Hour)     // 8 hours duration
	
	exception := &domain.AvailabilityException{
		BusinessID:     business.ID,
		StaffID:        staff.ID,
		ExceptionType:  "time_off",
		StartTime:      startTime,
		EndTime:        endTime,
		IsFullDay:      true,
		IsRecurring:    false,
		RecurrenceRule: "",
		Notes:          "Annual leave",
		CreatedBy:      createdBy,
	}

	// Test creation
	err = repo.Create(ctx, exception)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, exception.ExceptionID, "Exception ID should be generated")
	assert.NotZero(t, exception.CreatedAt, "Created at timestamp should be set")

	// Verify the exception was created in the database
	var savedException models.AvailabilityException
	err = testDB.First(&savedException, "id = ?", exception.ExceptionID).Error
	assert.NoError(t, err)
	assert.Equal(t, business.ID, savedException.BusinessID)
	assert.Equal(t, staff.ID, savedException.StaffID)
	assert.Equal(t, models.ExceptionTypeTimeOff, savedException.ExceptionType)
	assert.Equal(t, startTime.Truncate(time.Second), savedException.StartTime.Truncate(time.Second))
	assert.Equal(t, endTime.Truncate(time.Second), savedException.EndTime.Truncate(time.Second))
	assert.True(t, savedException.IsFullDay)
	assert.False(t, savedException.IsRecurring)
	assert.Equal(t, "Annual leave", savedException.Notes)
	assert.Equal(t, createdBy, *savedException.CreatedBy)
}

func TestAvailabilityExceptionRepositoryIntegration_GetByID(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	exception := createTestAvailabilityException(t, testDB.DB, business.ID, staff.ID, user.ID)

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Test GetByID
	ctx := context.Background()
	result, err := repo.GetByID(ctx, exception.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, exception.ID, result.ExceptionID)
	assert.Equal(t, business.ID, result.BusinessID)
	assert.Equal(t, staff.ID, result.StaffID)
	assert.Equal(t, string(exception.ExceptionType), result.ExceptionType)
	
	// Verify related entities are populated
	assert.NotNil(t, result.Staff)
	assert.Equal(t, staff.ID, result.Staff.StaffID)
}

func TestAvailabilityExceptionRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestAvailabilityExceptionRepositoryIntegration_GetByStaff(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business.ID, user.ID)
	
	// Create 3 exceptions for staff1
	for i := 0; i < 3; i++ {
		createTestAvailabilityException(t, testDB.DB, business.ID, staff1.ID, user.ID)
	}
	
	// Create 2 exceptions for staff2
	for i := 0; i < 2; i++ {
		createTestAvailabilityException(t, testDB.DB, business.ID, staff2.ID, user.ID)
	}

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Test GetByStaff
	ctx := context.Background()
	results1, err := repo.GetByStaff(ctx, staff1.ID)
	assert.NoError(t, err)
	assert.Len(t, results1, 3)
	
	results2, err := repo.GetByStaff(ctx, staff2.ID)
	assert.NoError(t, err)
	assert.Len(t, results2, 2)
	
	// Verify all exceptions are for the correct staff
	for _, e := range results1 {
		assert.Equal(t, staff1.ID, e.StaffID)
	}
	
	for _, e := range results2 {
		assert.Equal(t, staff2.ID, e.StaffID)
	}
}

func TestAvailabilityExceptionRepositoryIntegration_GetByStaffAndDateRange(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	
	// Create exceptions with different date ranges
	now := time.Now()
	
	// Exception 1: Today
	exception1 := createTestAvailabilityExceptionWithDates(t, testDB.DB, business.ID, staff.ID, user.ID, 
		now, now.Add(4*time.Hour), false)
	
	// Exception 2: Tomorrow
	tomorrow := now.Add(24 * time.Hour)
	exception2 := createTestAvailabilityExceptionWithDates(t, testDB.DB, business.ID, staff.ID, user.ID, 
		tomorrow, tomorrow.Add(4*time.Hour), false)
	
	// Exception 3: Next week
	nextWeek := now.Add(7 * 24 * time.Hour)
	exception3 := createTestAvailabilityExceptionWithDates(t, testDB.DB, business.ID, staff.ID, user.ID, 
		nextWeek, nextWeek.Add(4*time.Hour), false)
	
	// Exception 4: Recurring
	recurring := createTestAvailabilityExceptionWithDates(t, testDB.DB, business.ID, staff.ID, user.ID, 
		now.Add(14*24*time.Hour), now.Add(14*24*time.Hour + 4*time.Hour), true)

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Test GetByStaffAndDateRange (today through tomorrow)
	ctx := context.Background()
	dayAfterTomorrow := now.Add(2 * 24 * time.Hour)
	results, err := repo.GetByStaffAndDateRange(ctx, staff.ID, now, dayAfterTomorrow)
	assert.NoError(t, err)
	
	// Should include exceptions 1, 2, and recurring (4)
	assert.Len(t, results, 3)
	
	exceptionIDs := []uuid.UUID{results[0].ExceptionID, results[1].ExceptionID, results[2].ExceptionID}
	assert.Contains(t, exceptionIDs, exception1.ID)
	assert.Contains(t, exceptionIDs, exception2.ID)
	assert.Contains(t, exceptionIDs, recurring.ID) // Recurring exception is included regardless of date
	assert.NotContains(t, exceptionIDs, exception3.ID) // Next week should not be included
}

func TestAvailabilityExceptionRepositoryIntegration_Update(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	exception := createTestAvailabilityException(t, testDB.DB, business.ID, staff.ID, user.ID)

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Create update input
	ctx := context.Background()
	updatedType := "holiday"
	updatedStartTime := time.Now().Add(48 * time.Hour) // Two days from now
	updatedEndTime := updatedStartTime.Add(8 * time.Hour)
	updatedIsFullDay := false
	updatedNotes := "Updated notes"
	
	updateInput := &domain.UpdateAvailabilityExceptionInput{
		ExceptionType: &updatedType,
		StartTime:     &updatedStartTime,
		EndTime:       &updatedEndTime,
		IsFullDay:     &updatedIsFullDay,
		Notes:         &updatedNotes,
	}

	// Test Update
	err = repo.Update(ctx, exception.ID, updateInput, user.ID)
	assert.NoError(t, err)

	// Verify the exception was updated in the database
	var updatedException models.AvailabilityException
	err = testDB.First(&updatedException, "id = ?", exception.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, models.ExceptionTypeHoliday, updatedException.ExceptionType)
	assert.Equal(t, updatedStartTime.Truncate(time.Second), updatedException.StartTime.Truncate(time.Second))
	assert.Equal(t, updatedEndTime.Truncate(time.Second), updatedException.EndTime.Truncate(time.Second))
	assert.Equal(t, updatedIsFullDay, updatedException.IsFullDay)
	assert.Equal(t, updatedNotes, updatedException.Notes)
	assert.Equal(t, user.ID, *updatedException.UpdatedBy)
	assert.NotNil(t, updatedException.UpdatedAt)
}

func TestAvailabilityExceptionRepositoryIntegration_Delete(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	exception := createTestAvailabilityException(t, testDB.DB, business.ID, staff.ID, user.ID)

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Test Delete
	ctx := context.Background()
	err = repo.Delete(ctx, exception.ID, user.ID)
	assert.NoError(t, err)

	// Verify the exception was soft deleted
	var deletedException models.AvailabilityException
	err = testDB.Unscoped().First(&deletedException, "id = ?", exception.ID).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedException.DeletedAt)
	assert.True(t, deletedException.DeletedAt.Valid)
	assert.Equal(t, user.ID, *deletedException.DeletedBy)
	
	// Verify that the exception is not returned in normal queries
	var count int64
	err = testDB.Model(&models.AvailabilityException{}).Where("id = ?", exception.ID).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestAvailabilityExceptionRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business1 := createTestBusiness(t, testDB.DB, user.ID)
	business2 := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business1.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business2.ID, user.ID)
	
	// Create 5 exceptions for business1
	for i := 0; i < 5; i++ {
		createTestAvailabilityException(t, testDB.DB, business1.ID, staff1.ID, user.ID)
	}
	
	// Create 3 exceptions for business2
	for i := 0; i < 3; i++ {
		createTestAvailabilityException(t, testDB.DB, business2.ID, staff2.ID, user.ID)
	}

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Test ListByBusiness with pagination
	ctx := context.Background()
	
	// Get first page of business1 exceptions (3 per page)
	page1, err := repo.ListByBusiness(ctx, business1.ID, 1, 3)
	assert.NoError(t, err)
	assert.Len(t, page1, 3)
	
	// Get second page of business1 exceptions
	page2, err := repo.ListByBusiness(ctx, business1.ID, 2, 3)
	assert.NoError(t, err)
	assert.Len(t, page2, 2)
	
	// Get all business2 exceptions
	business2Exceptions, err := repo.ListByBusiness(ctx, business2.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, business2Exceptions, 3)
	
	// Verify all exceptions are for the correct business
	for _, e := range page1 {
		assert.Equal(t, business1.ID, e.BusinessID)
	}
	
	for _, e := range page2 {
		assert.Equal(t, business1.ID, e.BusinessID)
	}
	
	for _, e := range business2Exceptions {
		assert.Equal(t, business2.ID, e.BusinessID)
	}
	
	// Make sure pages contain different records
	allIDs := make(map[uuid.UUID]bool)
	for _, e := range page1 {
		allIDs[e.ExceptionID] = true
	}
	for _, e := range page2 {
		allIDs[e.ExceptionID] = true
	}
	assert.Len(t, allIDs, 5)
}

func TestAvailabilityExceptionRepositoryIntegration_ListByBusinessAndDateRange(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business.ID, user.ID)
	
	// Create exceptions with different date ranges
	now := time.Now()
	
	// Today (staff1)
	createTestAvailabilityExceptionWithDates(t, testDB.DB, business.ID, staff1.ID, user.ID, 
		now, now.Add(4*time.Hour), false)
	
	// Tomorrow (staff1)
	tomorrow := now.Add(24 * time.Hour)
	createTestAvailabilityExceptionWithDates(t, testDB.DB, business.ID, staff1.ID, user.ID, 
		tomorrow, tomorrow.Add(4*time.Hour), false)
	
	// Next week (staff1)
	nextWeek := now.Add(7 * 24 * time.Hour)
	createTestAvailabilityExceptionWithDates(t, testDB.DB, business.ID, staff1.ID, user.ID, 
		nextWeek, nextWeek.Add(4*time.Hour), false)
	
	// Today (staff2)
	createTestAvailabilityExceptionWithDates(t, testDB.DB, business.ID, staff2.ID, user.ID, 
		now, now.Add(4*time.Hour), false)
	
	// Recurring (staff2)
	createTestAvailabilityExceptionWithDates(t, testDB.DB, business.ID, staff2.ID, user.ID, 
		now.Add(14*24*time.Hour), now.Add(14*24*time.Hour + 4*time.Hour), true)

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Test ListByBusinessAndDateRange (today through tomorrow)
	ctx := context.Background()
	dayAfterTomorrow := now.Add(2 * 24 * time.Hour)
	results, err := repo.ListByBusinessAndDateRange(ctx, business.ID, now, dayAfterTomorrow, 1, 10)
	assert.NoError(t, err)
	
	// Should include 3 exceptions within the date range + 1 recurring
	assert.Len(t, results, 4)
	
	// Count by type
	todayCount := 0
	tomorrowCount := 0
	recurringCount := 0
	
	for _, e := range results {
		if e.IsRecurring {
			recurringCount++
		} else if e.StartTime.Day() == now.Day() {
			todayCount++
		} else if e.StartTime.Day() == tomorrow.Day() {
			tomorrowCount++
		}
	}
	
	assert.Equal(t, 2, todayCount) // One from each staff
	assert.Equal(t, 1, tomorrowCount)
	assert.Equal(t, 1, recurringCount)
}

func TestAvailabilityExceptionRepositoryIntegration_CountByBusiness(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business1 := createTestBusiness(t, testDB.DB, user.ID)
	business2 := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business1.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business2.ID, user.ID)
	
	// Create 4 exceptions for business1
	for i := 0; i < 4; i++ {
		createTestAvailabilityException(t, testDB.DB, business1.ID, staff1.ID, user.ID)
	}
	
	// Create 2 exceptions for business2
	for i := 0; i < 2; i++ {
		createTestAvailabilityException(t, testDB.DB, business2.ID, staff2.ID, user.ID)
	}

	// Create the repository
	repo := repository.NewAvailabilityExceptionRepository(testDB.DB)

	// Test CountByBusiness
	ctx := context.Background()
	count1, err := repo.CountByBusiness(ctx, business1.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(4), count1)
	
	count2, err := repo.CountByBusiness(ctx, business2.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count2)
}

