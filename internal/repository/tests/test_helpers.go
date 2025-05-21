package tests

import (
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Helper functions to create test data for repository integration tests
// Each function is available in two versions:
// 1. The original version using a direct DB connection (for backward compatibility)
// 2. A TX version that accepts a transaction for true test isolation

// createTestUser creates a test user in the database
func createTestUser(t *testing.T, db *database.DB) *models.User {
	user := &models.User{
		ClerkID:   uuid.New().String(), // Use UUID as a unique clerk_id
		Email:     uuid.New().String() + "@example.com", // Ensure unique email
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+1234567890",
		Role:      models.UserRoleStaff,
		IsActive:  true,
	}
	
	err := db.Create(user).Error
	require.NoError(t, err, "Failed to create test user")
	
	return user
}

// createTestUserTx creates a test user within a transaction
func createTestUserTx(t *testing.T, tx *gorm.DB) *models.User {
	user := &models.User{
		ClerkID:   uuid.New().String(), // Use UUID as a unique clerk_id
		Email:     uuid.New().String() + "@example.com", // Ensure unique email
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+1234567890",
		Role:      models.UserRoleStaff,
		IsActive:  true,
	}
	
	err := tx.Create(user).Error
	require.NoError(t, err, "Failed to create test user")
	
	return user
}

// createTestBusiness creates a test business in the database
func createTestBusiness(t *testing.T, db *database.DB, userID uuid.UUID) *models.Business {
	businessName := "test-business-" + uuid.New().String()
	business := &models.Business{
		UserID:           userID,
		Name:             businessName,
		DisplayName:      "Test Business",
		Description:      "Test business description",
		Address:          "123 Test St",
		City:             "Test City",
		Country:          "Test Country",
		PostalCode:       "12345",
		Phone:            "+9876543210",
		Email:            businessName + "@example.com",
		SubscriptionTier: models.SubscriptionTierBasic,
		IsActive:         true,
	}
	
	err := db.Create(business).Error
	require.NoError(t, err, "Failed to create test business")
	
	return business
}

// createTestBusinessTx creates a test business within a transaction
func createTestBusinessTx(t *testing.T, tx *gorm.DB, userID uuid.UUID) *models.Business {
	businessName := "test-business-" + uuid.New().String()
	business := &models.Business{
		UserID:           userID,
		Name:             businessName,
		DisplayName:      "Test Business",
		Description:      "Test business description",
		Address:          "123 Test St",
		City:             "Test City",
		Country:          "Test Country",
		PostalCode:       "12345",
		Phone:            "+9876543210",
		Email:            businessName + "@example.com",
		SubscriptionTier: models.SubscriptionTierBasic,
		IsActive:         true,
	}
	
	err := tx.Create(business).Error
	require.NoError(t, err, "Failed to create test business")
	
	return business
}

// createTestStaff creates a test staff member in the database
func createTestStaff(t *testing.T, db *database.DB, businessID, userID uuid.UUID) *models.Staff {
	staff := &models.Staff{
		BusinessID:      businessID,
		UserID:          userID,
		Position:        "Test Position",
		Bio:             "Test Bio",
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
	
	err := db.Create(staff).Error
	require.NoError(t, err, "Failed to create test staff")
	
	return staff
}

// createTestStaffTx creates a test staff member within a transaction
func createTestStaffTx(t *testing.T, tx *gorm.DB, businessID, userID uuid.UUID) *models.Staff {
	staff := &models.Staff{
		BusinessID:      businessID,
		UserID:          userID,
		Position:        "Test Position",
		Bio:             "Test Bio",
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
	require.NoError(t, err, "Failed to create test staff")
	
	return staff
}

// createTestStaffWithPosition creates a test staff member with a specific position
func createTestStaffWithPosition(t *testing.T, db *database.DB, businessID, userID uuid.UUID, position string) *models.Staff {
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
	
	err := db.Create(staff).Error
	require.NoError(t, err, "Failed to create test staff with position: "+position)
	
	return staff
}

// createTestService creates a test service in the database
func createTestService(t *testing.T, db *database.DB, businessID, createdByID uuid.UUID) *models.Service {
	// Create a service category first
	category := &models.ServiceCategory{
		Name:        "Test Category",
		Description: "Test category description",
	}
	
	err := db.Create(category).Error
	require.NoError(t, err, "Failed to create test service category")

	// Now create the service with a reference to the category
	categoryID := category.ID
	service := &models.Service{
		BusinessID:  businessID,
		CategoryID:  &categoryID,
		Name:        "Test Service " + uuid.New().String(), // Ensure unique name
		Description: "Test service description",
		Duration:    60, // 60 minutes
		Price:       100.0,
		IsActive:    true,
	}

	// Set created_by
	service.CreatedBy = &createdByID

	err = db.Create(service).Error
	require.NoError(t, err, "Failed to create test service")

	return service
}

// createTestServiceAssignment creates a test service assignment in the database
func createTestServiceAssignment(t *testing.T, db *database.DB, businessID, staffID, serviceID, createdByID uuid.UUID) *models.ServiceAssignment {
	assignment := &models.ServiceAssignment{
		BusinessID: businessID,
		StaffID:    staffID,
		ServiceID:  serviceID,
		IsActive:   true,
	}

	// Set created_by
	assignment.CreatedBy = &createdByID

	err := db.Create(assignment).Error
	require.NoError(t, err, "Failed to create test service assignment")

	return assignment
}

// createTestAvailabilityException creates a test availability exception in the database
func createTestAvailabilityException(t *testing.T, db *database.DB, businessID, staffID, userID uuid.UUID) *models.AvailabilityException {
	startTime := time.Now().Add(24 * time.Hour) // Tomorrow
	endTime := startTime.Add(8 * time.Hour)     // 8 hours duration
	
	return createTestAvailabilityExceptionWithDates(t, db, businessID, staffID, userID, startTime, endTime, false)
}

// createTestAvailabilityExceptionWithDates creates a test availability exception with specific dates
func createTestAvailabilityExceptionWithDates(
	t *testing.T, 
	db *database.DB, 
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
	
	err := db.Create(exception).Error
	require.NoError(t, err, "Failed to create test availability exception")
	
	return exception
}

// createTestStaffPerformance creates a test staff performance record in the database
func createTestStaffPerformance(t *testing.T, db *database.DB, businessID, staffID uuid.UUID) *models.StaffPerformance {
	startDate := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour) // Last month
	endDate := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)   // Yesterday
	
	return createTestStaffPerformanceWithPeriod(t, db, businessID, staffID, "monthly", startDate, endDate)
}

// createTestStaffPerformanceWithPeriod creates a test staff performance record with specific period and dates
func createTestStaffPerformanceWithPeriod(
	t *testing.T, 
	db *database.DB, 
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
	
	err := db.Create(performance).Error
	require.NoError(t, err, "Failed to create test staff performance record")
	
	return performance
}