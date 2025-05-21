package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStaffRepository_Create(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	staff := &domain.Staff{
		StaffID:         uuid.New(),
		BusinessID:      uuid.New(),
		UserID:          uuid.New(),
		Position:        "Hair Stylist",
		Bio:             "Experienced stylist with 5 years in the industry",
		SpecialtyAreas:  []string{"Hair Coloring", "Cutting"},
		ProfileImageURL: "https://example.com/profile.jpg",
		IsActive:        true,
		EmploymentType:  "full-time",
		JoinDate:        time.Now().Add(-365 * 24 * time.Hour), // 1 year ago
		CommissionRate:  0.2, // 20%
		CreatedAt:       time.Now(),
		CreatedBy:       uuid.New(),
	}

	// Expectations
	mockRepo.On("Create", ctx, staff).Return(nil)

	// Act
	err := mockRepo.Create(ctx, staff)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_Create_Error(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	staff := &domain.Staff{
		StaffID:         uuid.New(),
		BusinessID:      uuid.New(),
		UserID:          uuid.New(),
		Position:        "Hair Stylist",
		EmploymentType:  "full-time",
		JoinDate:        time.Now(),
		CreatedAt:       time.Now(),
		CreatedBy:       uuid.New(),
	}

	expectedErr := errors.New("database error")
	mockRepo.On("Create", ctx, staff).Return(expectedErr)

	// Act
	err := mockRepo.Create(ctx, staff)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_GetByID(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	staffID := uuid.New()
	expectedStaff := &domain.Staff{
		StaffID:         staffID,
		BusinessID:      uuid.New(),
		UserID:          uuid.New(),
		Position:        "Hair Stylist",
		EmploymentType:  "full-time",
		JoinDate:        time.Now().Add(-365 * 24 * time.Hour),
		IsActive:        true,
		CreatedAt:       time.Now().Add(-365 * 24 * time.Hour),
		CreatedBy:       uuid.New(),
	}

	mockRepo.On("GetByID", ctx, staffID).Return(expectedStaff, nil)

	// Act
	staff, err := mockRepo.GetByID(ctx, staffID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStaff, staff)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	staffID := uuid.New()
	expectedErr := errors.New("staff not found")

	mockRepo.On("GetByID", ctx, staffID).Return(nil, expectedErr)

	// Act
	staff, err := mockRepo.GetByID(ctx, staffID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, staff)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_GetByUserID(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	userID := uuid.New()
	expectedStaff := []*domain.Staff{
		{
			StaffID:         uuid.New(),
			BusinessID:      uuid.New(),
			UserID:          userID,
			Position:        "Hair Stylist",
			EmploymentType:  "full-time",
			JoinDate:        time.Now().Add(-365 * 24 * time.Hour),
			IsActive:        true,
			CreatedAt:       time.Now().Add(-365 * 24 * time.Hour),
			CreatedBy:       uuid.New(),
		},
		{
			StaffID:         uuid.New(),
			BusinessID:      uuid.New(),
			UserID:          userID,
			Position:        "Nail Technician",
			EmploymentType:  "part-time",
			JoinDate:        time.Now().Add(-180 * 24 * time.Hour),
			IsActive:        true,
			CreatedAt:       time.Now().Add(-180 * 24 * time.Hour),
			CreatedBy:       uuid.New(),
		},
	}

	mockRepo.On("GetByUserID", ctx, userID).Return(expectedStaff, nil)

	// Act
	staff, err := mockRepo.GetByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStaff, staff)
	assert.Len(t, staff, 2)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_GetByUserID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	userID := uuid.New()
	mockRepo.On("GetByUserID", ctx, userID).Return([]*domain.Staff{}, nil)

	// Act
	staff, err := mockRepo.GetByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, staff)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_GetByBusinessID(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	businessID := uuid.New()
	expectedStaff := []*domain.Staff{
		{
			StaffID:         uuid.New(),
			BusinessID:      businessID,
			UserID:          uuid.New(),
			Position:        "Hair Stylist",
			EmploymentType:  "full-time",
			JoinDate:        time.Now().Add(-365 * 24 * time.Hour),
			IsActive:        true,
			CreatedAt:       time.Now().Add(-365 * 24 * time.Hour),
			CreatedBy:       uuid.New(),
		},
		{
			StaffID:         uuid.New(),
			BusinessID:      businessID,
			UserID:          uuid.New(),
			Position:        "Nail Technician",
			EmploymentType:  "part-time",
			JoinDate:        time.Now().Add(-180 * 24 * time.Hour),
			IsActive:        true,
			CreatedAt:       time.Now().Add(-180 * 24 * time.Hour),
			CreatedBy:       uuid.New(),
		},
	}

	mockRepo.On("GetByBusinessID", ctx, businessID).Return(expectedStaff, nil)

	// Act
	staff, err := mockRepo.GetByBusinessID(ctx, businessID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStaff, staff)
	assert.Len(t, staff, 2)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_Update(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	staffID := uuid.New()
	updatedBy := uuid.New()
	
	position := "Senior Hair Stylist"
	bio := "Award-winning stylist with 8 years experience"
	specialtyAreas := []string{"Hair Coloring", "Cutting", "Wedding Styling"}
	isActive := true
	employmentType := "full-time"
	
	updateInput := &domain.UpdateStaffInput{
		Position:       &position,
		Bio:            &bio,
		SpecialtyAreas: &specialtyAreas,
		IsActive:       &isActive,
		EmploymentType: &employmentType,
	}

	mockRepo.On("Update", ctx, staffID, updateInput, updatedBy).Return(nil)

	// Act
	err := mockRepo.Update(ctx, staffID, updateInput, updatedBy)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_Update_NotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	staffID := uuid.New()
	updatedBy := uuid.New()
	
	position := "Senior Hair Stylist"
	updateInput := &domain.UpdateStaffInput{
		Position: &position,
	}

	expectedErr := errors.New("staff not found")
	mockRepo.On("Update", ctx, staffID, updateInput, updatedBy).Return(expectedErr)

	// Act
	err := mockRepo.Update(ctx, staffID, updateInput, updatedBy)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_Delete(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	staffID := uuid.New()
	deletedBy := uuid.New()

	mockRepo.On("Delete", ctx, staffID, deletedBy).Return(nil)

	// Act
	err := mockRepo.Delete(ctx, staffID, deletedBy)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	staffID := uuid.New()
	deletedBy := uuid.New()

	expectedErr := errors.New("staff not found")
	mockRepo.On("Delete", ctx, staffID, deletedBy).Return(expectedErr)

	// Act
	err := mockRepo.Delete(ctx, staffID, deletedBy)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_List(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	page := 1
	pageSize := 10
	
	expectedStaff := []*domain.Staff{
		{
			StaffID:        uuid.New(),
			BusinessID:     uuid.New(),
			UserID:         uuid.New(),
			Position:       "Hair Stylist",
			EmploymentType: "full-time",
			JoinDate:       time.Now().Add(-365 * 24 * time.Hour),
			IsActive:       true,
		},
		{
			StaffID:        uuid.New(),
			BusinessID:     uuid.New(),
			UserID:         uuid.New(),
			Position:       "Nail Technician",
			EmploymentType: "part-time",
			JoinDate:       time.Now().Add(-180 * 24 * time.Hour),
			IsActive:       true,
		},
	}

	mockRepo.On("List", ctx, page, pageSize).Return(expectedStaff, nil)

	// Act
	staff, err := mockRepo.List(ctx, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStaff, staff)
	assert.Len(t, staff, 2)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_ListByBusiness(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	businessID := uuid.New()
	page := 1
	pageSize := 10
	
	expectedStaff := []*domain.Staff{
		{
			StaffID:        uuid.New(),
			BusinessID:     businessID,
			UserID:         uuid.New(),
			Position:       "Hair Stylist",
			EmploymentType: "full-time",
			JoinDate:       time.Now().Add(-365 * 24 * time.Hour),
			IsActive:       true,
		},
		{
			StaffID:        uuid.New(),
			BusinessID:     businessID,
			UserID:         uuid.New(),
			Position:       "Nail Technician",
			EmploymentType: "part-time",
			JoinDate:       time.Now().Add(-180 * 24 * time.Hour),
			IsActive:       true,
		},
	}

	mockRepo.On("ListByBusiness", ctx, businessID, page, pageSize).Return(expectedStaff, nil)

	// Act
	staff, err := mockRepo.ListByBusiness(ctx, businessID, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStaff, staff)
	assert.Len(t, staff, 2)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_Search(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	businessID := uuid.New()
	query := "stylist"
	page := 1
	pageSize := 10
	
	expectedStaff := []*domain.Staff{
		{
			StaffID:        uuid.New(),
			BusinessID:     businessID,
			UserID:         uuid.New(),
			Position:       "Hair Stylist",
			EmploymentType: "full-time",
			JoinDate:       time.Now().Add(-365 * 24 * time.Hour),
			IsActive:       true,
		},
		{
			StaffID:        uuid.New(),
			BusinessID:     businessID,
			UserID:         uuid.New(),
			Position:       "Senior Stylist",
			EmploymentType: "full-time",
			JoinDate:       time.Now().Add(-500 * 24 * time.Hour),
			IsActive:       true,
		},
	}

	mockRepo.On("Search", ctx, businessID, query, page, pageSize).Return(expectedStaff, nil)

	// Act
	staff, err := mockRepo.Search(ctx, businessID, query, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStaff, staff)
	assert.Len(t, staff, 2)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_Count(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	expectedCount := int64(25)
	mockRepo.On("Count", ctx).Return(expectedCount, nil)

	// Act
	count, err := mockRepo.Count(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

func TestStaffRepository_CountByBusiness(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewStaffRepository(t)
	ctx := context.Background()
	
	businessID := uuid.New()
	expectedCount := int64(5)
	mockRepo.On("CountByBusiness", ctx, businessID).Return(expectedCount, nil)

	// Act
	count, err := mockRepo.CountByBusiness(ctx, businessID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}