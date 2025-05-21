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

func TestAvailabilityExceptionRepository_Create(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAvailabilityExceptionRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	exception := &domain.AvailabilityException{
		ExceptionID:    uuid.New(),
		BusinessID:     uuid.New(),
		StaffID:        uuid.New(),
		ExceptionType:  "time_off",
		StartTime:      now.Add(24 * time.Hour),
		EndTime:        now.Add(48 * time.Hour),
		IsFullDay:      true,
		IsRecurring:    false,
		RecurrenceRule: "",
		Notes:          "Annual leave",
		CreatedAt:      now,
		CreatedBy:      uuid.New(),
	}

	// Expectations
	mockRepo.On("Create", ctx, exception).Return(nil)

	// Act
	err := mockRepo.Create(ctx, exception)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAvailabilityExceptionRepository_GetByID(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAvailabilityExceptionRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	exceptionID := uuid.New()
	
	expectedException := &domain.AvailabilityException{
		ExceptionID:    exceptionID,
		BusinessID:     uuid.New(),
		StaffID:        uuid.New(),
		ExceptionType:  "time_off",
		StartTime:      now.Add(24 * time.Hour),
		EndTime:        now.Add(48 * time.Hour),
		IsFullDay:      true,
		IsRecurring:    false,
		RecurrenceRule: "",
		Notes:          "Annual leave",
		CreatedAt:      now,
		CreatedBy:      uuid.New(),
	}

	mockRepo.On("GetByID", ctx, exceptionID).Return(expectedException, nil)

	// Act
	exception, err := mockRepo.GetByID(ctx, exceptionID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedException, exception)
	mockRepo.AssertExpectations(t)
}

func TestAvailabilityExceptionRepository_GetByStaff(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAvailabilityExceptionRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	staffID := uuid.New()
	
	expectedExceptions := []*domain.AvailabilityException{
		{
			ExceptionID:    uuid.New(),
			BusinessID:     uuid.New(),
			StaffID:        staffID,
			ExceptionType:  "time_off",
			StartTime:      now.Add(24 * time.Hour),
			EndTime:        now.Add(48 * time.Hour),
			IsFullDay:      true,
			IsRecurring:    false,
			CreatedAt:      now.Add(-24 * time.Hour),
			CreatedBy:      uuid.New(),
		},
		{
			ExceptionID:    uuid.New(),
			BusinessID:     uuid.New(),
			StaffID:        staffID,
			ExceptionType:  "custom_hours",
			StartTime:      now.Add(72 * time.Hour),
			EndTime:        now.Add(76 * time.Hour),
			IsFullDay:      false,
			IsRecurring:    false,
			CreatedAt:      now.Add(-12 * time.Hour),
			CreatedBy:      uuid.New(),
		},
	}

	mockRepo.On("GetByStaff", ctx, staffID).Return(expectedExceptions, nil)

	// Act
	exceptions, err := mockRepo.GetByStaff(ctx, staffID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedExceptions, exceptions)
	assert.Len(t, exceptions, 2)
	mockRepo.AssertExpectations(t)
}

func TestAvailabilityExceptionRepository_GetByStaffAndDateRange(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAvailabilityExceptionRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	staffID := uuid.New()
	startDate := now.Add(24 * time.Hour)
	endDate := now.Add(96 * time.Hour)
	
	expectedExceptions := []*domain.AvailabilityException{
		{
			ExceptionID:    uuid.New(),
			BusinessID:     uuid.New(),
			StaffID:        staffID,
			ExceptionType:  "time_off",
			StartTime:      now.Add(24 * time.Hour),
			EndTime:        now.Add(48 * time.Hour),
			IsFullDay:      true,
			IsRecurring:    false,
			CreatedAt:      now.Add(-24 * time.Hour),
			CreatedBy:      uuid.New(),
		},
		{
			ExceptionID:    uuid.New(),
			BusinessID:     uuid.New(),
			StaffID:        staffID,
			ExceptionType:  "custom_hours",
			StartTime:      now.Add(72 * time.Hour),
			EndTime:        now.Add(76 * time.Hour),
			IsFullDay:      false,
			IsRecurring:    false,
			CreatedAt:      now.Add(-12 * time.Hour),
			CreatedBy:      uuid.New(),
		},
	}

	mockRepo.On("GetByStaffAndDateRange", ctx, staffID, startDate, endDate).Return(expectedExceptions, nil)

	// Act
	exceptions, err := mockRepo.GetByStaffAndDateRange(ctx, staffID, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedExceptions, exceptions)
	assert.Len(t, exceptions, 2)
	mockRepo.AssertExpectations(t)
}

func TestAvailabilityExceptionRepository_Update(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAvailabilityExceptionRepository(t)
	ctx := context.Background()
	
	exceptionID := uuid.New()
	updatedBy := uuid.New()
	
	notes := "Updated annual leave"
	isFullDay := false
	
	updateInput := &domain.UpdateAvailabilityExceptionInput{
		Notes:     &notes,
		IsFullDay: &isFullDay,
	}

	mockRepo.On("Update", ctx, exceptionID, updateInput, updatedBy).Return(nil)

	// Act
	err := mockRepo.Update(ctx, exceptionID, updateInput, updatedBy)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAvailabilityExceptionRepository_Delete(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAvailabilityExceptionRepository(t)
	ctx := context.Background()
	
	exceptionID := uuid.New()
	deletedBy := uuid.New()

	mockRepo.On("Delete", ctx, exceptionID, deletedBy).Return(nil)

	// Act
	err := mockRepo.Delete(ctx, exceptionID, deletedBy)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAvailabilityExceptionRepository_ListByBusiness(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAvailabilityExceptionRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	businessID := uuid.New()
	page := 1
	pageSize := 10
	
	expectedExceptions := []*domain.AvailabilityException{
		{
			ExceptionID:    uuid.New(),
			BusinessID:     businessID,
			StaffID:        uuid.New(),
			ExceptionType:  "time_off",
			StartTime:      now.Add(24 * time.Hour),
			EndTime:        now.Add(48 * time.Hour),
			IsFullDay:      true,
			IsRecurring:    false,
			CreatedAt:      now.Add(-24 * time.Hour),
			CreatedBy:      uuid.New(),
		},
		{
			ExceptionID:    uuid.New(),
			BusinessID:     businessID,
			StaffID:        uuid.New(),
			ExceptionType:  "custom_hours",
			StartTime:      now.Add(72 * time.Hour),
			EndTime:        now.Add(76 * time.Hour),
			IsFullDay:      false,
			IsRecurring:    false,
			CreatedAt:      now.Add(-12 * time.Hour),
			CreatedBy:      uuid.New(),
		},
	}

	mockRepo.On("ListByBusiness", ctx, businessID, page, pageSize).Return(expectedExceptions, nil)

	// Act
	exceptions, err := mockRepo.ListByBusiness(ctx, businessID, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedExceptions, exceptions)
	assert.Len(t, exceptions, 2)
	mockRepo.AssertExpectations(t)
}

func TestAvailabilityExceptionRepository_ListByBusinessAndDateRange(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAvailabilityExceptionRepository(t)
	ctx := context.Background()
	
	now := time.Now()
	businessID := uuid.New()
	startDate := now.Add(24 * time.Hour)
	endDate := now.Add(96 * time.Hour)
	page := 1
	pageSize := 10
	
	expectedExceptions := []*domain.AvailabilityException{
		{
			ExceptionID:    uuid.New(),
			BusinessID:     businessID,
			StaffID:        uuid.New(),
			ExceptionType:  "time_off",
			StartTime:      now.Add(24 * time.Hour),
			EndTime:        now.Add(48 * time.Hour),
			IsFullDay:      true,
			IsRecurring:    false,
			CreatedAt:      now.Add(-24 * time.Hour),
			CreatedBy:      uuid.New(),
		},
		{
			ExceptionID:    uuid.New(),
			BusinessID:     businessID,
			StaffID:        uuid.New(),
			ExceptionType:  "custom_hours",
			StartTime:      now.Add(72 * time.Hour),
			EndTime:        now.Add(76 * time.Hour),
			IsFullDay:      false,
			IsRecurring:    false,
			CreatedAt:      now.Add(-12 * time.Hour),
			CreatedBy:      uuid.New(),
		},
	}

	mockRepo.On("ListByBusinessAndDateRange", ctx, businessID, startDate, endDate, page, pageSize).Return(expectedExceptions, nil)

	// Act
	exceptions, err := mockRepo.ListByBusinessAndDateRange(ctx, businessID, startDate, endDate, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedExceptions, exceptions)
	assert.Len(t, exceptions, 2)
	mockRepo.AssertExpectations(t)
}

func TestAvailabilityExceptionRepository_CountByBusiness(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAvailabilityExceptionRepository(t)
	ctx := context.Background()
	
	businessID := uuid.New()
	expectedCount := int64(8)

	mockRepo.On("CountByBusiness", ctx, businessID).Return(expectedCount, nil)

	// Act
	count, err := mockRepo.CountByBusiness(ctx, businessID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}