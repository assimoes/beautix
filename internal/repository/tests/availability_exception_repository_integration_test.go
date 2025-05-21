package tests

import (
	"context"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAvailabilityExceptionRepositoryIntegration_Create(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()

	// Create a new availability exception
	createdBy := testData.User.ID
	ctx := context.Background()
	startTime := time.Now().Add(24 * time.Hour) // Tomorrow
	endTime := startTime.Add(8 * time.Hour)     // 8 hours duration
	
	exception := &domain.AvailabilityException{
		BusinessID:     testData.Business.ID,
		StaffID:        testData.Staff.ID,
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
	err := repos.AvailabilityExceptionRepo.Create(ctx, exception)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, exception.ExceptionID, "Exception ID should be generated")
	assert.NotZero(t, exception.CreatedAt, "Created at timestamp should be set")

	// Verify the exception was created
	result, err := repos.AvailabilityExceptionRepo.GetByID(ctx, exception.ExceptionID)
	assert.NoError(t, err)
	assert.Equal(t, testData.Business.ID, result.BusinessID)
	assert.Equal(t, testData.Staff.ID, result.StaffID)
	assert.Equal(t, "time_off", result.ExceptionType)
	assert.Equal(t, startTime.Truncate(time.Second), result.StartTime.Truncate(time.Second))
	assert.Equal(t, endTime.Truncate(time.Second), result.EndTime.Truncate(time.Second))
	assert.True(t, result.IsFullDay)
	assert.False(t, result.IsRecurring)
	assert.Equal(t, "Annual leave", result.Notes)
	assert.Equal(t, createdBy, result.CreatedBy)
}

func TestAvailabilityExceptionRepositoryIntegration_GetByID(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	exception := createTestAvailabilityExceptionTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID)

	// Test GetByID
	ctx := context.Background()
	result, err := repos.AvailabilityExceptionRepo.GetByID(ctx, exception.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, exception.ID, result.ExceptionID)
	assert.Equal(t, testData.Business.ID, result.BusinessID)
	assert.Equal(t, testData.Staff.ID, result.StaffID)
	assert.Equal(t, string(exception.ExceptionType), result.ExceptionType)
	
	// Verify related entities are populated
	assert.NotNil(t, result.Staff)
	assert.Equal(t, testData.Staff.ID, result.Staff.StaffID)
}

func TestAvailabilityExceptionRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := repos.AvailabilityExceptionRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestAvailabilityExceptionRepositoryIntegration_GetByStaff(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	staff2 := createTestStaffTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	
	// Create 3 exceptions for staff1
	for i := 0; i < 3; i++ {
		createTestAvailabilityExceptionTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID)
	}
	
	// Create 2 exceptions for staff2
	for i := 0; i < 2; i++ {
		createTestAvailabilityExceptionTx(t, suite.Tx, testData.Business.ID, staff2.ID, testData.User.ID)
	}

	// Test GetByStaff
	ctx := context.Background()
	results1, err := repos.AvailabilityExceptionRepo.GetByStaff(ctx, testData.Staff.ID)
	assert.NoError(t, err)
	assert.Len(t, results1, 3)
	
	results2, err := repos.AvailabilityExceptionRepo.GetByStaff(ctx, staff2.ID)
	assert.NoError(t, err)
	assert.Len(t, results2, 2)
	
	// Verify all exceptions are for the correct staff
	for _, e := range results1 {
		assert.Equal(t, testData.Staff.ID, e.StaffID)
	}
	
	for _, e := range results2 {
		assert.Equal(t, staff2.ID, e.StaffID)
	}
}

func TestAvailabilityExceptionRepositoryIntegration_GetByStaffAndDateRange(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create exceptions with different date ranges
	now := time.Now()
	
	// Exception 1: Today
	createTestAvailabilityExceptionWithDatesTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID, 
		now, now.Add(4*time.Hour), false)
	
	// Exception 2: Tomorrow
	tomorrow := now.Add(24 * time.Hour)
	createTestAvailabilityExceptionWithDatesTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID, 
		tomorrow, tomorrow.Add(4*time.Hour), false)
	
	// Exception 3: Next week
	nextWeek := now.Add(7 * 24 * time.Hour)
	createTestAvailabilityExceptionWithDatesTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID, 
		nextWeek, nextWeek.Add(4*time.Hour), false)
	
	// Exception 4: Recurring
	createTestAvailabilityExceptionWithDatesTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID, 
		now.Add(14*24*time.Hour), now.Add(14*24*time.Hour + 4*time.Hour), true)

	// Test GetByStaffAndDateRange (today through tomorrow)
	ctx := context.Background()
	dayAfterTomorrow := now.Add(2 * 24 * time.Hour)
	results, err := repos.AvailabilityExceptionRepo.GetByStaffAndDateRange(ctx, testData.Staff.ID, now, dayAfterTomorrow)
	assert.NoError(t, err)
	
	// Should include exceptions 1, 2, and recurring (4) = 3 total
	assert.Len(t, results, 3)
	
	// Verify we have the correct mix of exceptions
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
	
	assert.Equal(t, 1, todayCount)
	assert.Equal(t, 1, tomorrowCount)
	assert.Equal(t, 1, recurringCount)
}

func TestAvailabilityExceptionRepositoryIntegration_Update(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	exception := createTestAvailabilityExceptionTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID)

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
	err := repos.AvailabilityExceptionRepo.Update(ctx, exception.ID, updateInput, testData.User.ID)
	assert.NoError(t, err)

	// Verify the exception was updated
	updatedResult, err := repos.AvailabilityExceptionRepo.GetByID(ctx, exception.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedType, updatedResult.ExceptionType)
	assert.Equal(t, updatedStartTime.Truncate(time.Second), updatedResult.StartTime.Truncate(time.Second))
	assert.Equal(t, updatedEndTime.Truncate(time.Second), updatedResult.EndTime.Truncate(time.Second))
	assert.Equal(t, updatedIsFullDay, updatedResult.IsFullDay)
	assert.Equal(t, updatedNotes, updatedResult.Notes)
	assert.Equal(t, testData.User.ID, *updatedResult.UpdatedBy)
	assert.NotNil(t, updatedResult.UpdatedAt)
}

func TestAvailabilityExceptionRepositoryIntegration_Delete(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	exception := createTestAvailabilityExceptionTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID)

	// Test Delete
	ctx := context.Background()
	err := repos.AvailabilityExceptionRepo.Delete(ctx, exception.ID, testData.User.ID)
	assert.NoError(t, err)

	// Verify the exception was soft deleted
	_, err = repos.AvailabilityExceptionRepo.GetByID(ctx, exception.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestAvailabilityExceptionRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	staff2 := createTestStaffTx(t, suite.Tx, business2.ID, testData.User.ID)
	
	// Create 5 exceptions for business1
	for i := 0; i < 5; i++ {
		createTestAvailabilityExceptionTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID)
	}
	
	// Create 3 exceptions for business2
	for i := 0; i < 3; i++ {
		createTestAvailabilityExceptionTx(t, suite.Tx, business2.ID, staff2.ID, testData.User.ID)
	}

	// Test ListByBusiness with pagination
	ctx := context.Background()
	
	// Get first page of business1 exceptions (3 per page)
	page1, err := repos.AvailabilityExceptionRepo.ListByBusiness(ctx, testData.Business.ID, 1, 3)
	assert.NoError(t, err)
	assert.Len(t, page1, 3)
	
	// Get second page of business1 exceptions
	page2, err := repos.AvailabilityExceptionRepo.ListByBusiness(ctx, testData.Business.ID, 2, 3)
	assert.NoError(t, err)
	assert.Len(t, page2, 2)
	
	// Get all business2 exceptions
	business2Exceptions, err := repos.AvailabilityExceptionRepo.ListByBusiness(ctx, business2.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, business2Exceptions, 3)
	
	// Verify all exceptions are for the correct business
	for _, e := range page1 {
		assert.Equal(t, testData.Business.ID, e.BusinessID)
	}
	
	for _, e := range page2 {
		assert.Equal(t, testData.Business.ID, e.BusinessID)
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
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	staff2 := createTestStaffTx(t, suite.Tx, testData.Business.ID, testData.User.ID)
	
	// Create exceptions with different date ranges
	now := time.Now()
	
	// Today (staff1)
	createTestAvailabilityExceptionWithDatesTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID, 
		now, now.Add(4*time.Hour), false)
	
	// Tomorrow (staff1)
	tomorrow := now.Add(24 * time.Hour)
	createTestAvailabilityExceptionWithDatesTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID, 
		tomorrow, tomorrow.Add(4*time.Hour), false)
	
	// Next week (staff1)
	nextWeek := now.Add(7 * 24 * time.Hour)
	createTestAvailabilityExceptionWithDatesTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID, 
		nextWeek, nextWeek.Add(4*time.Hour), false)
	
	// Today (staff2)
	createTestAvailabilityExceptionWithDatesTx(t, suite.Tx, testData.Business.ID, staff2.ID, testData.User.ID, 
		now, now.Add(4*time.Hour), false)
	
	// Recurring (staff2)
	createTestAvailabilityExceptionWithDatesTx(t, suite.Tx, testData.Business.ID, staff2.ID, testData.User.ID, 
		now.Add(14*24*time.Hour), now.Add(14*24*time.Hour + 4*time.Hour), true)

	// Test ListByBusinessAndDateRange (today through tomorrow)
	ctx := context.Background()
	dayAfterTomorrow := now.Add(2 * 24 * time.Hour)
	results, err := repos.AvailabilityExceptionRepo.ListByBusinessAndDateRange(ctx, testData.Business.ID, now, dayAfterTomorrow, 1, 10)
	assert.NoError(t, err)
	
	// Should include 3 exceptions within the date range + 1 recurring = 4 total
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
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	staff2 := createTestStaffTx(t, suite.Tx, business2.ID, testData.User.ID)
	
	// Create 4 exceptions for business1
	for i := 0; i < 4; i++ {
		createTestAvailabilityExceptionTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, testData.User.ID)
	}
	
	// Create 2 exceptions for business2
	for i := 0; i < 2; i++ {
		createTestAvailabilityExceptionTx(t, suite.Tx, business2.ID, staff2.ID, testData.User.ID)
	}

	// Test CountByBusiness
	ctx := context.Background()
	count1, err := repos.AvailabilityExceptionRepo.CountByBusiness(ctx, testData.Business.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(4), count1)
	
	count2, err := repos.AvailabilityExceptionRepo.CountByBusiness(ctx, business2.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count2)
}