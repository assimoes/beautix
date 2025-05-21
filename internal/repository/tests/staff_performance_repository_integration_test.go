package tests

import (
	"context"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStaffPerformanceRepositoryIntegrationTx_Create(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()

	// Create a new staff performance record
	ctx := context.Background()
	startDate := time.Now().Add(-30 * 24 * time.Hour) // 30 days ago
	endDate := time.Now().Add(-1 * 24 * time.Hour)    // Yesterday
	
	performance := &domain.StaffPerformance{
		BusinessID:           testData.Business.ID,
		StaffID:              testData.Staff.ID,
		Period:               "monthly",
		StartDate:            startDate,
		EndDate:              endDate,
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

	// Test creation
	err := repos.StaffPerformanceRepo.Create(ctx, performance)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, performance.PerformanceID, "Performance ID should be generated")
	assert.NotZero(t, performance.CreatedAt, "Created at timestamp should be set")
	assert.NotZero(t, performance.UpdatedAt, "Updated at timestamp should be set")

	// Verify the performance record was created in the database using the repository
	result, err := repos.StaffPerformanceRepo.GetByID(ctx, performance.PerformanceID)
	assert.NoError(t, err)
	assert.Equal(t, testData.Business.ID, result.BusinessID)
	assert.Equal(t, testData.Staff.ID, result.StaffID)
	assert.Equal(t, "monthly", result.Period)
	assert.Equal(t, startDate.Truncate(time.Second), result.StartDate.Truncate(time.Second))
	assert.Equal(t, endDate.Truncate(time.Second), result.EndDate.Truncate(time.Second))
	assert.Equal(t, 50, result.TotalAppointments)
	assert.Equal(t, 45, result.CompletedAppointments)
	assert.Equal(t, 3, result.CanceledAppointments)
	assert.Equal(t, 2, result.NoShowAppointments)
	assert.Equal(t, 1500.50, result.TotalRevenue)
	assert.Equal(t, 4.8, result.AverageRating)
	assert.Equal(t, 85.5, result.ClientRetentionRate)
	assert.Equal(t, 10, result.NewClients)
	assert.Equal(t, 35, result.ReturnClients)
}

func TestStaffPerformanceRepositoryIntegrationTx_GetByID(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create a staff performance record
	performance := createTestStaffPerformanceTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID)

	// Test GetByID
	ctx := context.Background()
	result, err := repos.StaffPerformanceRepo.GetByID(ctx, performance.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, performance.ID, result.PerformanceID)
	assert.Equal(t, testData.Business.ID, result.BusinessID)
	assert.Equal(t, testData.Staff.ID, result.StaffID)
	assert.Equal(t, string(performance.Period), result.Period)
	assert.Equal(t, 50, result.TotalAppointments)
	assert.Equal(t, 1500.50, result.TotalRevenue)
	
	// Verify related entities are populated
	assert.NotNil(t, result.Staff)
	assert.Equal(t, testData.Staff.ID, result.Staff.StaffID)
}

func TestStaffPerformanceRepositoryIntegrationTx_GetByID_NotFound(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := repos.StaffPerformanceRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffPerformanceRepositoryIntegrationTx_GetByStaffAndPeriod(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create performance records with different periods
	// Monthly - last month
	lastMonthStart := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	lastMonthEnd := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	monthly := createTestStaffPerformanceWithPeriodTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, 
		"monthly", lastMonthStart, lastMonthEnd)
	
	// Weekly - last week
	lastWeekStart := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	lastWeekEnd := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	createTestStaffPerformanceWithPeriodTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, 
		"weekly", lastWeekStart, lastWeekEnd)
	
	// Daily - yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	createTestStaffPerformanceWithPeriodTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, 
		"daily", yesterday, yesterday)

	// Test GetByStaffAndPeriod
	ctx := context.Background()
	result, err := repos.StaffPerformanceRepo.GetByStaffAndPeriod(ctx, testData.Staff.ID, "monthly", lastMonthStart)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, monthly.ID, result.PerformanceID)
	assert.Equal(t, "monthly", result.Period)
	
	// Try with a non-existent period
	nonExistentDate := time.Now().AddDate(0, -2, 0) // 2 months ago
	result, err = repos.StaffPerformanceRepo.GetByStaffAndPeriod(ctx, testData.Staff.ID, "monthly", nonExistentDate)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

// Implementing a transaction-based version of the date range test
func TestStaffPerformanceRepositoryIntegrationTx_GetByStaffAndDateRange(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Define date ranges with clear separation to avoid flakiness
	// |----- January -----||----- February -----||----- March -----|
	januaryStart := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	januaryEnd := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)
	
	februaryStart := time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
	februaryEnd := time.Date(2023, 2, 28, 23, 59, 59, 0, time.UTC)
	
	marchStart := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)
	marchEnd := time.Date(2023, 3, 31, 23, 59, 59, 0, time.UTC)
	
	// Create performance records for different months
	// Create performance records for different months with distinct IDs
	januaryPerf := createTestStaffPerformanceWithPeriodTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, 
		"monthly", januaryStart, januaryEnd)
	
	createTestStaffPerformanceWithPeriodTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, 
		"monthly", februaryStart, februaryEnd)
	
	createTestStaffPerformanceWithPeriodTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, 
		"monthly", marchStart, marchEnd)

	// Test GetByStaffAndDateRange
	ctx := context.Background()
	
	// Case 1: Search for January only
	results, err := repos.StaffPerformanceRepo.GetByStaffAndDateRange(ctx, testData.Staff.ID, januaryStart, januaryEnd)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, januaryPerf.ID, results[0].PerformanceID)
	
	// Case 2: Search for January to February
	results, err = repos.StaffPerformanceRepo.GetByStaffAndDateRange(ctx, testData.Staff.ID, januaryStart, februaryEnd)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	
	// Case 3: Search for all three months
	results, err = repos.StaffPerformanceRepo.GetByStaffAndDateRange(ctx, testData.Staff.ID, januaryStart, marchEnd)
	assert.NoError(t, err)
	assert.Len(t, results, 3)
	
	// Case 4: Search for a non-existent range
	nonExistentStart := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	nonExistentEnd := time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC)
	results, err = repos.StaffPerformanceRepo.GetByStaffAndDateRange(ctx, testData.Staff.ID, nonExistentStart, nonExistentEnd)
	assert.NoError(t, err) // No error, just empty results
	assert.Len(t, results, 0)
}

func TestStaffPerformanceRepositoryIntegrationTx_Update(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create a staff performance record
	performance := createTestStaffPerformanceTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID)

	// Create update input
	ctx := context.Background()
	updatedPerformance := &domain.StaffPerformance{
		PerformanceID:         performance.ID,
		BusinessID:            testData.Business.ID,
		StaffID:               testData.Staff.ID,
		Period:                string(performance.Period),
		StartDate:             performance.StartDate,
		EndDate:               performance.EndDate,
		TotalAppointments:     60,          // Updated
		CompletedAppointments: 55,          // Updated
		CanceledAppointments:  3,
		NoShowAppointments:    2,
		TotalRevenue:          1800.75,     // Updated
		AverageRating:         4.9,         // Updated
		ClientRetentionRate:   90.0,        // Updated
		NewClients:            15,          // Updated
		ReturnClients:         40,          // Updated
	}

	// Test Update
	err := repos.StaffPerformanceRepo.Update(ctx, performance.ID, updatedPerformance)
	assert.NoError(t, err)

	// Verify the performance record was updated using the repository
	result, err := repos.StaffPerformanceRepo.GetByID(ctx, performance.ID)
	assert.NoError(t, err)
	assert.Equal(t, 60, result.TotalAppointments)
	assert.Equal(t, 55, result.CompletedAppointments)
	assert.Equal(t, 1800.75, result.TotalRevenue)
	assert.Equal(t, 4.9, result.AverageRating)
	assert.Equal(t, 90.0, result.ClientRetentionRate)
	assert.Equal(t, 15, result.NewClients)
	assert.Equal(t, 40, result.ReturnClients)
}

func TestStaffPerformanceRepositoryIntegrationTx_Delete(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create a staff performance record
	performance := createTestStaffPerformanceTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID)

	// Test Delete
	ctx := context.Background()
	err := repos.StaffPerformanceRepo.Delete(ctx, performance.ID)
	assert.NoError(t, err)

	// Verify the performance record was deleted using repository
	_, err = repos.StaffPerformanceRepo.GetByID(ctx, performance.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffPerformanceRepositoryIntegrationTx_ListByBusiness(t *testing.T) {
	// Initialize the transaction test suite
	suite := NewTransactionTestSuite(t)
	
	// Get repositories that use the transaction
	repos := suite.CreateTestRepositories()
	
	// Create test data
	testData := suite.CreateTestData()
	
	// Create a second staff member for the same business
	staff2 := createTestStaffWithPositionTx(t, suite.Tx, testData.Business.ID, testData.User.ID, "Senior Stylist")
	
	// Create a second business
	business2 := createTestBusinessTx(t, suite.Tx, testData.User.ID)
	staff3 := createTestStaffTx(t, suite.Tx, business2.ID, testData.User.ID)
	
	// Create performance records for business1
	// Monthly records for staff1
	createTestStaffPerformanceWithPeriodTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, 
		"monthly", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 0, -1))
	
	// Weekly records for staff1
	createTestStaffPerformanceWithPeriodTx(t, suite.Tx, testData.Business.ID, testData.Staff.ID, 
		"weekly", time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1))
	
	// Monthly records for staff2
	createTestStaffPerformanceWithPeriodTx(t, suite.Tx, testData.Business.ID, staff2.ID, 
		"monthly", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 0, -1))
	
	// Create performance records for business2
	createTestStaffPerformanceWithPeriodTx(t, suite.Tx, business2.ID, staff3.ID, 
		"monthly", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 0, -1))

	// Test ListByBusiness with all periods
	ctx := context.Background()
	allBusiness1Records, err := repos.StaffPerformanceRepo.ListByBusiness(ctx, testData.Business.ID, "", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, allBusiness1Records, 3)
	
	// Test ListByBusiness filtered by period
	monthlyBusiness1Records, err := repos.StaffPerformanceRepo.ListByBusiness(ctx, testData.Business.ID, "monthly", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, monthlyBusiness1Records, 2)
	
	weeklyBusiness1Records, err := repos.StaffPerformanceRepo.ListByBusiness(ctx, testData.Business.ID, "weekly", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, weeklyBusiness1Records, 1)
	
	// Test ListByBusiness for business2
	business2Records, err := repos.StaffPerformanceRepo.ListByBusiness(ctx, business2.ID, "", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, business2Records, 1)
	
	// Verify records belong to the correct business
	for _, r := range allBusiness1Records {
		assert.Equal(t, testData.Business.ID, r.BusinessID)
	}
	
	for _, r := range business2Records {
		assert.Equal(t, business2.ID, r.BusinessID)
	}
}