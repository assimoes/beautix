package tests

import (
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// createTestStaffWithPositionTx creates a test staff member with a specific position within a transaction
func createTestStaffWithPositionTx(t *testing.T, tx *gorm.DB, businessID, userID uuid.UUID, position string) *models.Staff {
	staff := &models.Staff{
		BusinessID:      businessID,
		UserID:          userID,
		Position:        position,
		Bio:             "Bio for " + position,
		SpecialtyAreas:  models.SpecialtyAreas{"Test Area 1", "Test Area 2"},
		ProfileImageURL: "http://example.com/test.jpg",
		IsActive:        true,
		EmploymentType:  models.StaffEmploymentTypeFull,
		JoinDate:        time.Now().Add(-30 * 24 * time.Hour),
		CommissionRate:  15.0,
	}

	// Set created_by
	createdBy := userID
	staff.CreatedBy = &createdBy

	err := tx.Create(staff).Error
	require.NoError(t, err, "Failed to create test staff with position: "+position)

	return staff
}

// createTestServiceWithCategoryTx creates a test service with a new category in the database within a transaction
func createTestServiceWithCategoryTx(t *testing.T, tx *gorm.DB, businessID, createdByID uuid.UUID) *models.Service {
	// Create a service category first
	category := &models.ServiceCategory{
		BusinessID:  businessID,
		Name:        "Test Category",
		Description: "Test category description",
	}

	err := tx.Create(category).Error
	require.NoError(t, err, "Failed to create test service category")

	// Now create the service with a reference to the category
	service := &models.Service{
		BusinessID:  businessID,
		Category:    category.Name,
		Name:        "Test Service " + uuid.New().String(), // Ensure unique name
		Description: "Test service description",
		Duration:    60, // 60 minutes
		Price:       100.0,
		IsActive:    true,
	}

	// Set created_by
	service.CreatedBy = &createdByID

	err = tx.Create(service).Error
	require.NoError(t, err, "Failed to create test service")

	return service
}

// createTestServiceAssignmentTx creates a test service assignment in the database within a transaction
func createTestServiceAssignmentTx(t *testing.T, tx *gorm.DB, businessID, staffID, serviceID, createdByID uuid.UUID) *models.ServiceAssignment {
	assignment := &models.ServiceAssignment{
		BusinessID: businessID,
		StaffID:    staffID,
		ServiceID:  serviceID,
		IsActive:   true,
	}

	// Set created_by
	assignment.CreatedBy = &createdByID

	err := tx.Create(assignment).Error
	require.NoError(t, err, "Failed to create test service assignment")

	return assignment
}

// createTestAvailabilityExceptionTx creates a test availability exception in the database within a transaction
func createTestAvailabilityExceptionTx(t *testing.T, tx *gorm.DB, businessID, staffID, userID uuid.UUID) *models.AvailabilityException {
	startTime := time.Now().Add(24 * time.Hour) // Tomorrow
	endTime := startTime.Add(8 * time.Hour)     // 8 hours duration

	return createTestAvailabilityExceptionWithDatesTx(t, tx, businessID, staffID, userID, startTime, endTime, false)
}

// createTestAvailabilityExceptionWithDatesTx creates a test availability exception with specific dates within a transaction
func createTestAvailabilityExceptionWithDatesTx(
	t *testing.T,
	tx *gorm.DB,
	businessID,
	staffID,
	userID uuid.UUID,
	startTime,
	endTime time.Time,
	isRecurring bool,
) *models.AvailabilityException {
	exception := &models.AvailabilityException{
		BusinessID:     businessID,
		StaffID:        staffID,
		ExceptionType:  models.ExceptionTypeTimeOff,
		StartTime:      startTime,
		EndTime:        endTime,
		IsFullDay:      true,
		IsRecurring:    isRecurring,
		RecurrenceRule: "",
		Notes:          "Test exception",
	}

	// Set created_by
	createdBy := userID
	exception.CreatedBy = &createdBy

	err := tx.Create(exception).Error
	require.NoError(t, err, "Failed to create test availability exception")

	return exception
}

// createTestStaffPerformanceTx creates a test staff performance record in the database within a transaction
func createTestStaffPerformanceTx(t *testing.T, tx *gorm.DB, businessID, staffID uuid.UUID) *models.StaffPerformance {
	startDate := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour) // Last month
	endDate := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)   // Yesterday

	return createTestStaffPerformanceWithPeriodTx(t, tx, businessID, staffID, "monthly", startDate, endDate)
}

// createTestStaffPerformanceWithPeriodTx creates a test staff performance record with specific period and dates within a transaction
func createTestStaffPerformanceWithPeriodTx(
	t *testing.T,
	tx *gorm.DB,
	businessID,
	staffID uuid.UUID,
	period string,
	startDate,
	endDate time.Time,
) *models.StaffPerformance {
	performance := &models.StaffPerformance{
		BusinessID:            businessID,
		StaffID:               staffID,
		Period:                models.PerformancePeriod(period),
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

	err := tx.Create(performance).Error
	require.NoError(t, err, "Failed to create test staff performance record")

	return performance
}
