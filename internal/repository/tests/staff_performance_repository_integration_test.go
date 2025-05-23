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

func TestStaffPerformanceRepositoryIntegration_Create(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffPerformanceRepo := &repository.StaffPerformanceRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "staff_performance_create@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "create-staff-performance-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "create_staff_perf@test.com")

	// Create a new staff performance record
	startDate := time.Now().Add(-30 * 24 * time.Hour) // 30 days ago
	endDate := time.Now().Add(-1 * 24 * time.Hour)    // Yesterday

	performance := &domain.StaffPerformance{
		BusinessID:            business.BusinessID,
		StaffID:               staff.StaffID,
		Period:                "monthly",
		StartDate:             startDate,
		EndDate:               endDate,
		TotalAppointments:     50,
		CompletedAppointments: 45,
		CanceledAppointments:  3,
		NoShowAppointments:    2,
		TotalRevenue:          1500.50,
		AverageRating:         4.8,
		ClientRetentionRate:   85.5,
		NewClients:            10,
		ReturnClients:         35,
	}

	// Test creation
	err = staffPerformanceRepo.Create(ctx, performance)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, performance.PerformanceID)
	assert.NotZero(t, performance.CreatedAt)
	assert.NotZero(t, performance.UpdatedAt)

	// Verify the performance record was created
	result, err := staffPerformanceRepo.GetByID(ctx, performance.PerformanceID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, staff.StaffID, result.StaffID)
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

func TestStaffPerformanceRepositoryIntegration_GetByID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffPerformanceRepo := &repository.StaffPerformanceRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "staff_performance_getbyid@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbyid-staff-performance-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "getbyid_staff_perf@test.com")
	
	// Create a staff performance record
	performance := createTestStaffPerformanceForIntegration(t, staffPerformanceRepo, business.BusinessID, staff.StaffID)

	// Test GetByID
	result, err := staffPerformanceRepo.GetByID(ctx, performance.PerformanceID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, performance.PerformanceID, result.PerformanceID)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, staff.StaffID, result.StaffID)
	assert.Equal(t, performance.Period, result.Period)
	assert.Equal(t, 50, result.TotalAppointments)
	assert.Equal(t, 1500.50, result.TotalRevenue)
	
	// Verify related staff is populated
	assert.NotNil(t, result.Staff)
	assert.Equal(t, staff.StaffID, result.Staff.StaffID)
}

func TestStaffPerformanceRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	staffPerformanceRepo := &repository.StaffPerformanceRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := staffPerformanceRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffPerformanceRepositoryIntegration_GetByStaffAndPeriod(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffPerformanceRepo := &repository.StaffPerformanceRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "staff_performance_byperiod@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "byperiod-staff-performance-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "byperiod_staff_perf@test.com")
	
	// Create performance records with different periods
	// Monthly - last month
	lastMonthStart := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	lastMonthEnd := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	monthly := createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business.BusinessID, staff.StaffID,
		"monthly", lastMonthStart, lastMonthEnd)
	
	// Weekly - last week
	lastWeekStart := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	lastWeekEnd := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business.BusinessID, staff.StaffID,
		"weekly", lastWeekStart, lastWeekEnd)
	
	// Daily - yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business.BusinessID, staff.StaffID,
		"daily", yesterday, yesterday)

	// Test GetByStaffAndPeriod
	result, err := staffPerformanceRepo.GetByStaffAndPeriod(ctx, staff.StaffID, "monthly", lastMonthStart)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, monthly.PerformanceID, result.PerformanceID)
	assert.Equal(t, "monthly", result.Period)

	// Try with a non-existent period
	nonExistentDate := time.Now().AddDate(0, -2, 0) // 2 months ago
	result, err = staffPerformanceRepo.GetByStaffAndPeriod(ctx, staff.StaffID, "monthly", nonExistentDate)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffPerformanceRepositoryIntegration_GetByStaffAndDateRange(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffPerformanceRepo := &repository.StaffPerformanceRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "staff_performance_daterange@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "daterange-staff-performance-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "daterange_staff_perf@test.com")
	
	// Define date ranges with clear separation to avoid flakiness
	// |----- January -----||----- February -----||----- March -----|
	januaryStart := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	januaryEnd := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)
	
	februaryStart := time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
	februaryEnd := time.Date(2023, 2, 28, 23, 59, 59, 0, time.UTC)
	
	marchStart := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)
	marchEnd := time.Date(2023, 3, 31, 23, 59, 59, 0, time.UTC)
	
	// Create performance records for different months
	januaryPerf := createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business.BusinessID, staff.StaffID,
		"monthly", januaryStart, januaryEnd)
	
	createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business.BusinessID, staff.StaffID,
		"monthly", februaryStart, februaryEnd)
	
	createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business.BusinessID, staff.StaffID,
		"monthly", marchStart, marchEnd)

	// Test GetByStaffAndDateRange
	// Case 1: Search for January only
	results, err := staffPerformanceRepo.GetByStaffAndDateRange(ctx, staff.StaffID, januaryStart, januaryEnd)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, januaryPerf.PerformanceID, results[0].PerformanceID)

	// Case 2: Search for January to February
	results, err = staffPerformanceRepo.GetByStaffAndDateRange(ctx, staff.StaffID, januaryStart, februaryEnd)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Case 3: Search for all three months
	results, err = staffPerformanceRepo.GetByStaffAndDateRange(ctx, staff.StaffID, januaryStart, marchEnd)
	require.NoError(t, err)
	assert.Len(t, results, 3)

	// Case 4: Search for a non-existent range
	nonExistentStart := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	nonExistentEnd := time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC)
	results, err = staffPerformanceRepo.GetByStaffAndDateRange(ctx, staff.StaffID, nonExistentStart, nonExistentEnd)
	require.NoError(t, err) // No error, just empty results
	assert.Len(t, results, 0)
}

func TestStaffPerformanceRepositoryIntegration_Update(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffPerformanceRepo := &repository.StaffPerformanceRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "staff_performance_update@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "update-staff-performance-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "update_staff_perf@test.com")
	
	// Create a staff performance record
	performance := createTestStaffPerformanceForIntegration(t, staffPerformanceRepo, business.BusinessID, staff.StaffID)

	// Create update input
	updatedPerformance := &domain.StaffPerformance{
		PerformanceID:         performance.PerformanceID,
		BusinessID:            business.BusinessID,
		StaffID:               staff.StaffID,
		Period:                performance.Period,
		StartDate:             performance.StartDate,
		EndDate:               performance.EndDate,
		TotalAppointments:     60, // Updated
		CompletedAppointments: 55, // Updated
		CanceledAppointments:  3,
		NoShowAppointments:    2,
		TotalRevenue:          1800.75, // Updated
		AverageRating:         4.9,     // Updated
		ClientRetentionRate:   90.0,    // Updated
		NewClients:            15,      // Updated
		ReturnClients:         40,      // Updated
	}

	// Test Update
	err = staffPerformanceRepo.Update(ctx, performance.PerformanceID, updatedPerformance)
	require.NoError(t, err)

	// Verify the performance record was updated
	result, err := staffPerformanceRepo.GetByID(ctx, performance.PerformanceID)
	require.NoError(t, err)
	assert.Equal(t, 60, result.TotalAppointments)
	assert.Equal(t, 55, result.CompletedAppointments)
	assert.Equal(t, 1800.75, result.TotalRevenue)
	assert.Equal(t, 4.9, result.AverageRating)
	assert.Equal(t, 90.0, result.ClientRetentionRate)
	assert.Equal(t, 15, result.NewClients)
	assert.Equal(t, 40, result.ReturnClients)
}

func TestStaffPerformanceRepositoryIntegration_Delete(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffPerformanceRepo := &repository.StaffPerformanceRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "staff_performance_delete@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "delete-staff-performance-salon")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "delete_staff_perf@test.com")
	
	// Create a staff performance record
	performance := createTestStaffPerformanceForIntegration(t, staffPerformanceRepo, business.BusinessID, staff.StaffID)

	// Test Delete
	err = staffPerformanceRepo.Delete(ctx, performance.PerformanceID)
	require.NoError(t, err)

	// Verify the performance record was deleted
	_, err = staffPerformanceRepo.GetByID(ctx, performance.PerformanceID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffPerformanceRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffPerformanceRepo := &repository.StaffPerformanceRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "staff_performance_listbybusiness@test.com")
	business1 := createTestBusinessForService(t, businessRepo, user.UserID, "listbybusiness-staff-performance-salon1")
	business2 := createTestBusinessForService(t, businessRepo, user.UserID, "listbybusiness-staff-performance-salon2")
	staff1 := createTestStaffForAppointment(t, staffRepo, business1.BusinessID, user.UserID, "listbybusiness_staff_perf1@test.com")
	staff2 := createTestStaffForAppointment(t, staffRepo, business1.BusinessID, user.UserID, "listbybusiness_staff_perf2@test.com")
	staff3 := createTestStaffForAppointment(t, staffRepo, business2.BusinessID, user.UserID, "listbybusiness_staff_perf3@test.com")
	
	// Create performance records for business1
	// Monthly records for staff1
	createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business1.BusinessID, staff1.StaffID,
		"monthly", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 0, -1))
	
	// Weekly records for staff1
	createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business1.BusinessID, staff1.StaffID,
		"weekly", time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1))
	
	// Monthly records for staff2
	createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business1.BusinessID, staff2.StaffID,
		"monthly", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 0, -1))
	
	// Create performance records for business2
	createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, business2.BusinessID, staff3.StaffID,
		"monthly", time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 0, -1))

	// Test ListByBusiness with all periods
	allBusiness1Records, err := staffPerformanceRepo.ListByBusiness(ctx, business1.BusinessID, "", 1, 10)
	require.NoError(t, err)
	assert.Len(t, allBusiness1Records, 3)

	// Test ListByBusiness filtered by period
	monthlyBusiness1Records, err := staffPerformanceRepo.ListByBusiness(ctx, business1.BusinessID, "monthly", 1, 10)
	require.NoError(t, err)
	assert.Len(t, monthlyBusiness1Records, 2)

	weeklyBusiness1Records, err := staffPerformanceRepo.ListByBusiness(ctx, business1.BusinessID, "weekly", 1, 10)
	require.NoError(t, err)
	assert.Len(t, weeklyBusiness1Records, 1)

	// Test ListByBusiness for business2
	business2Records, err := staffPerformanceRepo.ListByBusiness(ctx, business2.BusinessID, "", 1, 10)
	require.NoError(t, err)
	assert.Len(t, business2Records, 1)

	// Verify records belong to the correct business
	for _, r := range allBusiness1Records {
		assert.Equal(t, business1.BusinessID, r.BusinessID)
	}

	for _, r := range business2Records {
		assert.Equal(t, business2.BusinessID, r.BusinessID)
	}
}

// Helper function to create a test staff performance record
func createTestStaffPerformanceForIntegration(t *testing.T, staffPerformanceRepo *repository.StaffPerformanceRepository, businessID, staffID uuid.UUID) *domain.StaffPerformance {
	startDate := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour) // Last month
	endDate := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)   // Yesterday
	
	return createTestStaffPerformanceWithPeriodForIntegration(t, staffPerformanceRepo, businessID, staffID, "monthly", startDate, endDate)
}

// Helper function to create a test staff performance record with specific period and dates
func createTestStaffPerformanceWithPeriodForIntegration(t *testing.T, staffPerformanceRepo *repository.StaffPerformanceRepository, businessID, staffID uuid.UUID, period string, startDate, endDate time.Time) *domain.StaffPerformance {
	performance := &domain.StaffPerformance{
		BusinessID:            businessID,
		StaffID:               staffID,
		Period:                period,
		StartDate:             startDate,
		EndDate:               endDate,
		TotalAppointments:     50,
		CompletedAppointments: 45,
		CanceledAppointments:  3,
		NoShowAppointments:    2,
		TotalRevenue:          1500.50,
		AverageRating:         4.8,
		ClientRetentionRate:   85.5,
		NewClients:            10,
		ReturnClients:         35,
	}
	
	err := staffPerformanceRepo.Create(context.Background(), performance)
	require.NoError(t, err)
	
	return performance
}