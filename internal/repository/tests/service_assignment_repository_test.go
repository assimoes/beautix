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
)

func TestServiceAssignmentRepository_Create(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	assignment := &domain.ServiceAssignment{
		AssignmentID: uuid.New(),
		BusinessID:   uuid.New(),
		StaffID:      uuid.New(),
		ServiceID:    uuid.New(),
		IsActive:     true,
		CreatedAt:    time.Now(),
		CreatedBy:    uuid.New(),
	}

	// Expectations
	mockRepo.On("Create", ctx, assignment).Return(nil)

	// Act
	err := mockRepo.Create(ctx, assignment)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestServiceAssignmentRepository_Create_Error(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	assignment := &domain.ServiceAssignment{
		AssignmentID: uuid.New(),
		BusinessID:   uuid.New(),
		StaffID:      uuid.New(),
		ServiceID:    uuid.New(),
		IsActive:     true,
		CreatedAt:    time.Now(),
		CreatedBy:    uuid.New(),
	}

	expectedErr := errors.New("database error")
	mockRepo.On("Create", ctx, assignment).Return(expectedErr)

	// Act
	err := mockRepo.Create(ctx, assignment)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestServiceAssignmentRepository_GetByID(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	assignmentID := uuid.New()
	expectedAssignment := &domain.ServiceAssignment{
		AssignmentID: assignmentID,
		BusinessID:   uuid.New(),
		StaffID:      uuid.New(),
		ServiceID:    uuid.New(),
		IsActive:     true,
		CreatedAt:    time.Now().Add(-24 * time.Hour),
		CreatedBy:    uuid.New(),
	}

	mockRepo.On("GetByID", ctx, assignmentID).Return(expectedAssignment, nil)

	// Act
	assignment, err := mockRepo.GetByID(ctx, assignmentID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAssignment, assignment)
	mockRepo.AssertExpectations(t)
}

func TestServiceAssignmentRepository_GetByStaffAndService(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	staffID := uuid.New()
	serviceID := uuid.New()
	
	expectedAssignment := &domain.ServiceAssignment{
		AssignmentID: uuid.New(),
		BusinessID:   uuid.New(),
		StaffID:      staffID,
		ServiceID:    serviceID,
		IsActive:     true,
		CreatedAt:    time.Now().Add(-24 * time.Hour),
		CreatedBy:    uuid.New(),
	}

	mockRepo.On("GetByStaffAndService", ctx, staffID, serviceID).Return(expectedAssignment, nil)

	// Act
	assignment, err := mockRepo.GetByStaffAndService(ctx, staffID, serviceID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAssignment, assignment)
	mockRepo.AssertExpectations(t)
}

func TestServiceAssignmentRepository_GetByStaff(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	staffID := uuid.New()
	
	expectedAssignments := []*domain.ServiceAssignment{
		{
			AssignmentID: uuid.New(),
			BusinessID:   uuid.New(),
			StaffID:      staffID,
			ServiceID:    uuid.New(),
			IsActive:     true,
			CreatedAt:    time.Now().Add(-48 * time.Hour),
			CreatedBy:    uuid.New(),
		},
		{
			AssignmentID: uuid.New(),
			BusinessID:   uuid.New(),
			StaffID:      staffID,
			ServiceID:    uuid.New(),
			IsActive:     true,
			CreatedAt:    time.Now().Add(-24 * time.Hour),
			CreatedBy:    uuid.New(),
		},
	}

	mockRepo.On("GetByStaff", ctx, staffID).Return(expectedAssignments, nil)

	// Act
	assignments, err := mockRepo.GetByStaff(ctx, staffID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAssignments, assignments)
	assert.Len(t, assignments, 2)
	mockRepo.AssertExpectations(t)
}

func TestServiceAssignmentRepository_GetByService(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	serviceID := uuid.New()
	
	expectedAssignments := []*domain.ServiceAssignment{
		{
			AssignmentID: uuid.New(),
			BusinessID:   uuid.New(),
			StaffID:      uuid.New(),
			ServiceID:    serviceID,
			IsActive:     true,
			CreatedAt:    time.Now().Add(-48 * time.Hour),
			CreatedBy:    uuid.New(),
		},
		{
			AssignmentID: uuid.New(),
			BusinessID:   uuid.New(),
			StaffID:      uuid.New(),
			ServiceID:    serviceID,
			IsActive:     true,
			CreatedAt:    time.Now().Add(-24 * time.Hour),
			CreatedBy:    uuid.New(),
		},
	}

	mockRepo.On("GetByService", ctx, serviceID).Return(expectedAssignments, nil)

	// Act
	assignments, err := mockRepo.GetByService(ctx, serviceID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAssignments, assignments)
	assert.Len(t, assignments, 2)
	mockRepo.AssertExpectations(t)
}

func TestServiceAssignmentRepository_Update(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	assignmentID := uuid.New()
	updatedBy := uuid.New()
	
	isActive := false
	updateInput := &domain.UpdateServiceAssignmentInput{
		IsActive: &isActive,
	}

	mockRepo.On("Update", ctx, assignmentID, updateInput, updatedBy).Return(nil)

	// Act
	err := mockRepo.Update(ctx, assignmentID, updateInput, updatedBy)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestServiceAssignmentRepository_Delete(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	assignmentID := uuid.New()
	deletedBy := uuid.New()

	mockRepo.On("Delete", ctx, assignmentID, deletedBy).Return(nil)

	// Act
	err := mockRepo.Delete(ctx, assignmentID, deletedBy)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestServiceAssignmentRepository_ListByBusiness(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	businessID := uuid.New()
	page := 1
	pageSize := 10
	
	expectedAssignments := []*domain.ServiceAssignment{
		{
			AssignmentID: uuid.New(),
			BusinessID:   businessID,
			StaffID:      uuid.New(),
			ServiceID:    uuid.New(),
			IsActive:     true,
			CreatedAt:    time.Now().Add(-48 * time.Hour),
			CreatedBy:    uuid.New(),
		},
		{
			AssignmentID: uuid.New(),
			BusinessID:   businessID,
			StaffID:      uuid.New(),
			ServiceID:    uuid.New(),
			IsActive:     true,
			CreatedAt:    time.Now().Add(-24 * time.Hour),
			CreatedBy:    uuid.New(),
		},
	}

	mockRepo.On("ListByBusiness", ctx, businessID, page, pageSize).Return(expectedAssignments, nil)

	// Act
	assignments, err := mockRepo.ListByBusiness(ctx, businessID, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAssignments, assignments)
	assert.Len(t, assignments, 2)
	mockRepo.AssertExpectations(t)
}

func TestServiceAssignmentRepository_CountByBusiness(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewServiceAssignmentRepository(t)
	ctx := context.Background()
	
	businessID := uuid.New()
	expectedCount := int64(12)

	mockRepo.On("CountByBusiness", ctx, businessID).Return(expectedCount, nil)

	// Act
	count, err := mockRepo.CountByBusiness(ctx, businessID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}