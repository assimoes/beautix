package tests

import (
	"context"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This is an example of how to use the TransactionTestSuite to create
// truly idempotent tests that are isolated from each other
func TestStaffRepositoryTx_Create(t *testing.T) {
	// Create a test suite that manages transactions
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Define a new staff member to create
	ctx := context.Background()
	staff := &domain.Staff{
		BusinessID:      testData.Business.ID,
		UserID:          testData.User.ID,
		Position:        "Transaction Test Position",
		Bio:             "Test Bio",
		SpecialtyAreas:  []string{"Area 1", "Area 2"},
		ProfileImageURL: "http://example.com/profile.jpg",
		IsActive:        true,
		EmploymentType:  "full-time",
		JoinDate:        time.Now().Add(-30 * 24 * time.Hour),
		CommissionRate:  20.0,
		CreatedBy:       testData.User.ID,
	}

	// Test creation
	err := repos.StaffRepo.Create(ctx, staff)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, staff.StaffID, "Staff ID should be generated")
	assert.NotZero(t, staff.CreatedAt, "Created at timestamp should be set")
	
	// Verify the staff was created in the database
	// Note: We can query using the transaction to verify the data was created correctly
	createdStaff, err := repos.StaffRepo.GetByID(ctx, staff.StaffID)
	assert.NoError(t, err)
	assert.Equal(t, "Transaction Test Position", createdStaff.Position)
	
	// When this test completes, the transaction will be automatically rolled back
	// so no data will persist in the database
}

// This test demonstrates multiple repository operations within a transaction
func TestMultipleRepositoryOperationsTx(t *testing.T) {
	// Create a test suite that manages transactions
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	ctx := context.Background()
	
	// 1. Create a staff member
	staff := &domain.Staff{
		BusinessID:      testData.Business.ID,
		UserID:          testData.User.ID,
		Position:        "Transaction Test Position",
		Bio:             "Test Bio",
		SpecialtyAreas:  []string{"Area 1", "Area 2"},
		ProfileImageURL: "http://example.com/profile.jpg",
		IsActive:        true,
		EmploymentType:  "full-time",
		JoinDate:        time.Now().Add(-30 * 24 * time.Hour),
		CommissionRate:  20.0,
		CreatedBy:       testData.User.ID,
	}
	err := repos.StaffRepo.Create(ctx, staff)
	require.NoError(t, err)
	
	// 2. Create an availability exception for this staff member
	startTime := time.Now().Add(24 * time.Hour) // Tomorrow
	endTime := startTime.Add(8 * time.Hour)     // 8 hours duration
	
	exception := &domain.AvailabilityException{
		BusinessID:     testData.Business.ID,
		StaffID:        staff.StaffID,
		ExceptionType:  "time_off",
		StartTime:      startTime,
		EndTime:        endTime,
		IsFullDay:      true,
		IsRecurring:    false,
		RecurrenceRule: "",
		Notes:          "Annual leave",
		CreatedBy:      testData.User.ID,
	}
	
	err = repos.AvailabilityExceptionRepo.Create(ctx, exception)
	require.NoError(t, err)
	
	// 3. Create a performance record
	performance := &domain.StaffPerformance{
		BusinessID:           testData.Business.ID,
		StaffID:              staff.StaffID,
		Period:               "monthly",
		StartDate:            time.Now().Add(-30 * 24 * time.Hour),
		EndDate:              time.Now().Add(-1 * 24 * time.Hour),
		TotalAppointments:    50,
		CompletedAppointments: 45,
		CanceledAppointments: 3,
		NoShowAppointments:   2,
		TotalRevenue:         1500.50,
		AverageRating:        4.8,
		ClientRetentionRate:  85.5,
		NewClients:           10,
		ReturnClients:        35,
	}
	
	err = repos.StaffPerformanceRepo.Create(ctx, performance)
	require.NoError(t, err)
	
	// All these operations happened in a single transaction
	// and will be rolled back when the test completes
}