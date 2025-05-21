package tests

import (
	"context"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStaffPerformanceRepository_Create(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffPerformanceRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month
	
	performance := &domain.StaffPerformance{
		PerformanceID:        uuid.New(),
		BusinessID:           uuid.New(),
		StaffID:              uuid.New(),
		Period:               "monthly",
		StartDate:            startDate,
		EndDate:              endDate,
		TotalAppointments:    25,
		CompletedAppointments: 22,
		CanceledAppointments: 2,
		NoShowAppointments:   1,
		TotalRevenue:         1250.50,
		AverageRating:        4.8,
		ClientRetentionRate:  0.85,
		NewClients:           5,
		ReturnClients:        15,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	// Expectations
	mockRepo.On("Create", ctx, performance).Return(nil)

	// Act
	err := mockRepo.Create(ctx, performance)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffPerformanceRepository_GetByID(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffPerformanceRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	performanceID := uuid.New()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month
	
	expectedPerformance := &domain.StaffPerformance{
		PerformanceID:        performanceID,
		BusinessID:           uuid.New(),
		StaffID:              uuid.New(),
		Period:               "monthly",
		StartDate:            startDate,
		EndDate:              endDate,
		TotalAppointments:    25,
		CompletedAppointments: 22,
		CanceledAppointments: 2,
		NoShowAppointments:   1,
		TotalRevenue:         1250.50,
		AverageRating:        4.8,
		ClientRetentionRate:  0.85,
		NewClients:           5,
		ReturnClients:        15,
		CreatedAt:            now.Add(-24 * time.Hour),
		UpdatedAt:            now.Add(-24 * time.Hour),
	}

	mockRepo.On("GetByID", ctx, performanceID).Return(expectedPerformance, nil)

	// Act
	performance, err := mockRepo.GetByID(ctx, performanceID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPerformance, performance)
	mockRepo.AssertExpectations(t)
}

func TestStaffPerformanceRepository_GetByStaffAndPeriod(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffPerformanceRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	staffID := uuid.New()
	period := "monthly"
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month
	
	expectedPerformance := &domain.StaffPerformance{
		PerformanceID:        uuid.New(),
		BusinessID:           uuid.New(),
		StaffID:              staffID,
		Period:               period,
		StartDate:            startDate,
		EndDate:              endDate,
		TotalAppointments:    25,
		CompletedAppointments: 22,
		CanceledAppointments: 2,
		NoShowAppointments:   1,
		TotalRevenue:         1250.50,
		AverageRating:        4.8,
		ClientRetentionRate:  0.85,
		NewClients:           5,
		ReturnClients:        15,
		CreatedAt:            now.Add(-24 * time.Hour),
		UpdatedAt:            now.Add(-24 * time.Hour),
	}

	mockRepo.On("GetByStaffAndPeriod", ctx, staffID, period, startDate).Return(expectedPerformance, nil)

	// Act
	performance, err := mockRepo.GetByStaffAndPeriod(ctx, staffID, period, startDate)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPerformance, performance)
	mockRepo.AssertExpectations(t)
}

func TestStaffPerformanceRepository_GetByStaffAndDateRange(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffPerformanceRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	staffID := uuid.New()
	
	// First day of current month
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	
	// Last day of next month
	endDate := time.Date(now.Year(), now.Month()+2, 0, 23, 59, 59, 999999999, now.Location())
	
	expectedPerformances := []*domain.StaffPerformance{
		{
			PerformanceID:        uuid.New(),
			BusinessID:           uuid.New(),
			StaffID:              staffID,
			Period:               "monthly",
			StartDate:            startDate,
			EndDate:              time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location()),
			TotalAppointments:    25,
			CompletedAppointments: 22,
			CanceledAppointments: 2,
			NoShowAppointments:   1,
			TotalRevenue:         1250.50,
			AverageRating:        4.8,
			ClientRetentionRate:  0.85,
			NewClients:           5,
			ReturnClients:        15,
			CreatedAt:            now.Add(-24 * time.Hour),
			UpdatedAt:            now.Add(-24 * time.Hour),
		},
		{
			PerformanceID:        uuid.New(),
			BusinessID:           uuid.New(),
			StaffID:              staffID,
			Period:               "monthly",
			StartDate:            time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()),
			EndDate:              endDate,
			TotalAppointments:    30,
			CompletedAppointments: 27,
			CanceledAppointments: 3,
			NoShowAppointments:   0,
			TotalRevenue:         1500.75,
			AverageRating:        4.9,
			ClientRetentionRate:  0.90,
			NewClients:           7,
			ReturnClients:        20,
			CreatedAt:            now.Add(-12 * time.Hour),
			UpdatedAt:            now.Add(-12 * time.Hour),
		},
	}

	mockRepo.On("GetByStaffAndDateRange", ctx, staffID, startDate, endDate).Return(expectedPerformances, nil)

	// Act
	performances, err := mockRepo.GetByStaffAndDateRange(ctx, staffID, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPerformances, performances)
	assert.Len(t, performances, 2)
	mockRepo.AssertExpectations(t)
}

func TestStaffPerformanceRepository_Update(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffPerformanceRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	performanceID := uuid.New()
	
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month
	
	// Updated performance data
	performance := &domain.StaffPerformance{
		PerformanceID:        performanceID,
		BusinessID:           uuid.New(),
		StaffID:              uuid.New(),
		Period:               "monthly",
		StartDate:            startDate,
		EndDate:              endDate,
		TotalAppointments:    28, // Updated
		CompletedAppointments: 25, // Updated
		CanceledAppointments: 3, // Updated
		NoShowAppointments:   0, // Updated
		TotalRevenue:         1350.50, // Updated
		AverageRating:        4.7, // Updated
		ClientRetentionRate:  0.88, // Updated
		NewClients:           6, // Updated
		ReturnClients:        19, // Updated
		CreatedAt:            now.Add(-24 * time.Hour),
		UpdatedAt:            now, // Updated
	}

	mockRepo.On("Update", ctx, performanceID, performance).Return(nil)

	// Act
	err := mockRepo.Update(ctx, performanceID, performance)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffPerformanceRepository_Delete(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffPerformanceRepository(t)
	ctx := context.Background()
	
	performanceID := uuid.New()

	mockRepo.On("Delete", ctx, performanceID).Return(nil)

	// Act
	err := mockRepo.Delete(ctx, performanceID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffPerformanceRepository_ListByBusiness(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffPerformanceRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	businessID := uuid.New()
	period := "monthly"
	page := 1
	pageSize := 10
	
	// First day of current month
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	
	// Last day of current month
	endDate := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location())
	
	expectedPerformances := []*domain.StaffPerformance{
		{
			PerformanceID:        uuid.New(),
			BusinessID:           businessID,
			StaffID:              uuid.New(),
			Period:               period,
			StartDate:            startDate,
			EndDate:              endDate,
			TotalAppointments:    25,
			CompletedAppointments: 22,
			CanceledAppointments: 2,
			NoShowAppointments:   1,
			TotalRevenue:         1250.50,
			AverageRating:        4.8,
			ClientRetentionRate:  0.85,
			NewClients:           5,
			ReturnClients:        15,
			CreatedAt:            now.Add(-24 * time.Hour),
			UpdatedAt:            now.Add(-24 * time.Hour),
		},
		{
			PerformanceID:        uuid.New(),
			BusinessID:           businessID,
			StaffID:              uuid.New(),
			Period:               period,
			StartDate:            startDate,
			EndDate:              endDate,
			TotalAppointments:    30,
			CompletedAppointments: 27,
			CanceledAppointments: 3,
			NoShowAppointments:   0,
			TotalRevenue:         1500.75,
			AverageRating:        4.9,
			ClientRetentionRate:  0.90,
			NewClients:           7,
			ReturnClients:        20,
			CreatedAt:            now.Add(-12 * time.Hour),
			UpdatedAt:            now.Add(-12 * time.Hour),
		},
	}

	mockRepo.On("ListByBusiness", ctx, businessID, period, page, pageSize).Return(expectedPerformances, nil)

	// Act
	performances, err := mockRepo.ListByBusiness(ctx, businessID, period, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPerformances, performances)
	assert.Len(t, performances, 2)
	mockRepo.AssertExpectations(t)
}