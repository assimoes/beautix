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

func TestStaffPerformanceRepositoryIntegration_Create(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)

	// Create the repository
	repo := repository.NewStaffPerformanceRepository(testDB.DB)

	// Create a new staff performance record
	ctx := context.Background()
	startDate := time.Now().Add(-30 * 24 * time.Hour) // 30 days ago
	endDate := time.Now().Add(-1 * 24 * time.Hour)    // Yesterday
	
	performance := &domain.StaffPerformance{
		BusinessID:           business.ID,
		StaffID:              staff.ID,
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
	err = repo.Create(ctx, performance)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, performance.PerformanceID, "Performance ID should be generated")
	assert.NotZero(t, performance.CreatedAt, "Created at timestamp should be set")
	assert.NotZero(t, performance.UpdatedAt, "Updated at timestamp should be set")

	// Verify the performance record was created in the database
	var savedPerformance models.StaffPerformance
	err = testDB.First(&savedPerformance, "id = ?", performance.PerformanceID).Error
	assert.NoError(t, err)
	assert.Equal(t, business.ID, savedPerformance.BusinessID)
	assert.Equal(t, staff.ID, savedPerformance.StaffID)
	assert.Equal(t, models.PerformancePeriodMonthly, savedPerformance.Period)
	assert.Equal(t, startDate.Truncate(time.Second), savedPerformance.StartDate.Truncate(time.Second))
	assert.Equal(t, endDate.Truncate(time.Second), savedPerformance.EndDate.Truncate(time.Second))
	assert.Equal(t, 50, savedPerformance.TotalAppointments)
	assert.Equal(t, 45, savedPerformance.CompletedAppointments)
	assert.Equal(t, 3, savedPerformance.CanceledAppointments)
	assert.Equal(t, 2, savedPerformance.NoShowAppointments)
	assert.Equal(t, 1500.50, savedPerformance.TotalRevenue)
	assert.Equal(t, 4.8, savedPerformance.AverageRating)
	assert.Equal(t, 85.5, savedPerformance.ClientRetentionRate)
	assert.Equal(t, 10, savedPerformance.NewClients)
	assert.Equal(t, 35, savedPerformance.ReturnClients)
}

func TestStaffPerformanceRepositoryIntegration_GetByID(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	performance := createTestStaffPerformance(t, testDB.DB, business.ID, staff.ID)

	// Create the repository
	repo := repository.NewStaffPerformanceRepository(testDB.DB)

	// Test GetByID
	ctx := context.Background()
	result, err := repo.GetByID(ctx, performance.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, performance.ID, result.PerformanceID)
	assert.Equal(t, business.ID, result.BusinessID)
	assert.Equal(t, staff.ID, result.StaffID)
	assert.Equal(t, string(performance.Period), result.Period)
	assert.Equal(t, 50, result.TotalAppointments)
	assert.Equal(t, 1500.50, result.TotalRevenue)
	
	// Verify related entities are populated
	assert.NotNil(t, result.Staff)
	assert.Equal(t, staff.ID, result.Staff.StaffID)
}

func TestStaffPerformanceRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create the repository
	repo := repository.NewStaffPerformanceRepository(testDB.DB)

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffPerformanceRepositoryIntegration_GetByStaffAndPeriod(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	
	// Create performance records with different periods
	// Monthly - last month
	lastMonthStart := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	lastMonthEnd := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	monthly := createTestStaffPerformanceWithPeriod(t, testDB.DB, business.ID, staff.ID, 
		"monthly", lastMonthStart, lastMonthEnd)
	
	// Weekly - last week
	lastWeekStart := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	lastWeekEnd := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	createTestStaffPerformanceWithPeriod(t, testDB.DB, business.ID, staff.ID, 
		"weekly", lastWeekStart, lastWeekEnd)
	
	// Daily - yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	createTestStaffPerformanceWithPeriod(t, testDB.DB, business.ID, staff.ID, 
		"daily", yesterday, yesterday)

	// Create the repository
	repo := repository.NewStaffPerformanceRepository(testDB.DB)

	// Test GetByStaffAndPeriod
	ctx := context.Background()
	result, err := repo.GetByStaffAndPeriod(ctx, staff.ID, "monthly", lastMonthStart)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, monthly.ID, result.PerformanceID)
	assert.Equal(t, "monthly", result.Period)
	
	// Try with a non-existent period
	nonExistentDate := time.Now().AddDate(0, -2, 0) // 2 months ago
	result, err = repo.GetByStaffAndPeriod(ctx, staff.ID, "monthly", nonExistentDate)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffPerformanceRepositoryIntegration_GetByStaffAndDateRange(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business.ID, user.ID)
	
	// Create performance records across different date ranges for staff1
	now := time.Now()
	
	// Last month
	lastMonthStart := now.AddDate(0, -1, 0).Truncate(24 * time.Hour)
	lastMonthEnd := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
	monthly1 := createTestStaffPerformanceWithPeriod(t, testDB.DB, business.ID, staff1.ID, 
		"monthly", lastMonthStart, lastMonthEnd)
	
	// Last week
	lastWeekStart := now.AddDate(0, 0, -7).Truncate(24 * time.Hour)
	lastWeekEnd := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
	weekly1 := createTestStaffPerformanceWithPeriod(t, testDB.DB, business.ID, staff1.ID, 
		"weekly", lastWeekStart, lastWeekEnd)
	
	// Create a performance record for staff2
	createTestStaffPerformanceWithPeriod(t, testDB.DB, business.ID, staff2.ID, 
		"monthly", lastMonthStart, lastMonthEnd)

	// Create the repository
	repo := repository.NewStaffPerformanceRepository(testDB.DB)

	// Test GetByStaffAndDateRange
	ctx := context.Background()
	// Query range that covers both monthly and weekly records
	queryStart := lastMonthStart.AddDate(0, 0, -1) // Day before last month started
	queryEnd := now
	results, err := repo.GetByStaffAndDateRange(ctx, staff1.ID, queryStart, queryEnd)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	
	// Verify both records are included
	performanceIDs := []uuid.UUID{results[0].PerformanceID, results[1].PerformanceID}
	assert.Contains(t, performanceIDs, monthly1.ID)
	assert.Contains(t, performanceIDs, weekly1.ID)
	
	// Query range that only covers weekly record
	queryStart = lastWeekStart.AddDate(0, 0, -1) // Day before last week started
	queryEnd = lastWeekEnd.AddDate(0, 0, 1)      // Day after last week ended
	results, err = repo.GetByStaffAndDateRange(ctx, staff1.ID, queryStart, queryEnd)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, weekly1.ID, results[0].PerformanceID)
}

func TestStaffPerformanceRepositoryIntegration_Update(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	performance := createTestStaffPerformance(t, testDB.DB, business.ID, staff.ID)

	// Create the repository
	repo := repository.NewStaffPerformanceRepository(testDB.DB)

	// Create update input
	ctx := context.Background()
	updatedPerformance := &domain.StaffPerformance{
		PerformanceID:         performance.ID,
		BusinessID:            business.ID,
		StaffID:               staff.ID,
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
	err = repo.Update(ctx, performance.ID, updatedPerformance)
	assert.NoError(t, err)

	// Verify the performance record was updated in the database
	var updatedRecord models.StaffPerformance
	err = testDB.First(&updatedRecord, "id = ?", performance.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, 60, updatedRecord.TotalAppointments)
	assert.Equal(t, 55, updatedRecord.CompletedAppointments)
	assert.Equal(t, 1800.75, updatedRecord.TotalRevenue)
	assert.Equal(t, 4.9, updatedRecord.AverageRating)
	assert.Equal(t, 90.0, updatedRecord.ClientRetentionRate)
	assert.Equal(t, 15, updatedRecord.NewClients)
	assert.Equal(t, 40, updatedRecord.ReturnClients)
}

func TestStaffPerformanceRepositoryIntegration_Delete(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business := createTestBusiness(t, testDB.DB, user.ID)
	staff := createTestStaff(t, testDB.DB, business.ID, user.ID)
	performance := createTestStaffPerformance(t, testDB.DB, business.ID, staff.ID)

	// Create the repository
	repo := repository.NewStaffPerformanceRepository(testDB.DB)

	// Test Delete
	ctx := context.Background()
	err = repo.Delete(ctx, performance.ID)
	assert.NoError(t, err)

	// Verify the performance record was deleted (hard delete)
	var count int64
	err = testDB.Unscoped().Model(&models.StaffPerformance{}).Where("id = ?", performance.ID).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count, "Performance record should be hard deleted")
}

func TestStaffPerformanceRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Create test data
	user := createTestUser(t, testDB.DB)
	business1 := createTestBusiness(t, testDB.DB, user.ID)
	business2 := createTestBusiness(t, testDB.DB, user.ID)
	staff1 := createTestStaff(t, testDB.DB, business1.ID, user.ID)
	staff2 := createTestStaff(t, testDB.DB, business1.ID, user.ID)
	staff3 := createTestStaff(t, testDB.DB, business2.ID, user.ID)
	
	// Create performance records for business1
	// Monthly records for staff1
	createTestStaffPerformanceWithPeriod(t, testDB.DB, business1.ID, staff1.ID, 
		"monthly", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 0, -1))
	
	// Weekly records for staff1
	createTestStaffPerformanceWithPeriod(t, testDB.DB, business1.ID, staff1.ID, 
		"weekly", time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1))
	
	// Monthly records for staff2
	createTestStaffPerformanceWithPeriod(t, testDB.DB, business1.ID, staff2.ID, 
		"monthly", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 0, -1))
	
	// Create performance records for business2
	createTestStaffPerformanceWithPeriod(t, testDB.DB, business2.ID, staff3.ID, 
		"monthly", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 0, -1))

	// Create the repository
	repo := repository.NewStaffPerformanceRepository(testDB.DB)

	// Test ListByBusiness with all periods
	ctx := context.Background()
	allBusiness1Records, err := repo.ListByBusiness(ctx, business1.ID, "", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, allBusiness1Records, 3)
	
	// Test ListByBusiness filtered by period
	monthlyBusiness1Records, err := repo.ListByBusiness(ctx, business1.ID, "monthly", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, monthlyBusiness1Records, 2)
	
	weeklyBusiness1Records, err := repo.ListByBusiness(ctx, business1.ID, "weekly", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, weeklyBusiness1Records, 1)
	
	// Test ListByBusiness for business2
	business2Records, err := repo.ListByBusiness(ctx, business2.ID, "", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, business2Records, 1)
	
	// Verify records belong to the correct business
	for _, r := range allBusiness1Records {
		assert.Equal(t, business1.ID, r.BusinessID)
	}
	
	for _, r := range business2Records {
		assert.Equal(t, business2.ID, r.BusinessID)
	}
}

