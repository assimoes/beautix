package models_test

import (
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvailabilityExceptionModel(t *testing.T) {
	// Connect to the test database using simple approach
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Clean up all tables comprehensively to avoid foreign key issues
	database.CleanupAllTables(t, testDB.DB)

	// Create a user
	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		ClerkID:   "clerk_availability_" + userID.String()[:8], // Unique ClerkID
		Email:     "availability_test@example.com",
		FirstName: "Availability",
		LastName:  "Test",
		Phone:     "+1234567890",
		Role:      models.UserRoleStaff,
		IsActive:  true,
	}

	// Save the user
	err = testDB.DB.Create(&user).Error
	assert.NoError(t, err, "Failed to create user")

	// Create a business
	businessID := uuid.New()
	business := models.Business{
		BaseModel: models.BaseModel{
			ID: businessID,
		},
		UserID:       userID,
		Name:         "availability-salon",
		BusinessType: "salon",
		DisplayName:  "Availability Salon",
		Address:      "123 Availability St",
		City:         "Lisbon",
		Country:      "Portugal",
		Phone:        "+351123456789",
		Email:        "contact@availability-salon.com",
		IsActive:     true,
	}

	// Save the business
	err = testDB.DB.Create(&business).Error
	assert.NoError(t, err, "Failed to create business")

	// Create a staff member
	staffID := uuid.New()
	joinDate := time.Now().Add(-90 * 24 * time.Hour) // 90 days ago
	staff := models.Staff{
		BaseModel: models.BaseModel{
			ID:        staffID,
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		BusinessID:     businessID,
		UserID:         userID,
		Position:       "Hair Stylist",
		Bio:            "Experienced hair stylist",
		SpecialtyAreas: models.SpecialtyAreas{"Hair Cutting", "Coloring"},
		IsActive:       true,
		EmploymentType: models.StaffEmploymentTypeFull,
		JoinDate:       joinDate,
		CommissionRate: 18.00,
	}

	// Save the staff
	err = testDB.DB.Create(&staff).Error
	assert.NoError(t, err, "Failed to create staff")

	// Create an availability exception for time off
	exceptionID := uuid.New()
	tomorrow := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	exception := models.AvailabilityException{
		BaseModel: models.BaseModel{
			ID:        exceptionID,
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		BusinessID:    businessID,
		StaffID:       staffID,
		ExceptionType: models.ExceptionTypeTimeOff,
		StartTime:     tomorrow.Add(9 * time.Hour),  // 9:00 AM
		EndTime:       tomorrow.Add(17 * time.Hour), // 5:00 PM
		IsFullDay:     true,
		IsRecurring:   false,
		Notes:         "Personal day off",
	}

	// Save the exception
	err = testDB.DB.Create(&exception).Error
	assert.NoError(t, err, "Failed to create availability exception")

	// Verify exception was created with ID
	var savedException models.AvailabilityException
	err = testDB.DB.First(&savedException, "id = ?", exceptionID).Error
	assert.NoError(t, err, "Failed to find availability exception")
	assert.Equal(t, exceptionID, savedException.ID)
	assert.Equal(t, businessID, savedException.BusinessID)
	assert.Equal(t, staffID, savedException.StaffID)
	assert.Equal(t, models.ExceptionTypeTimeOff, savedException.ExceptionType)
	assert.True(t, savedException.IsFullDay)
	assert.False(t, savedException.IsRecurring)
	assert.Equal(t, "Personal day off", savedException.Notes)

	// Test loaded relationships
	err = testDB.DB.Preload("Staff").Preload("Business").First(&savedException, "id = ?", exceptionID).Error
	if err != nil {
		t.Logf("Preload error: %v", err)
	}
	assert.NoError(t, err, "Failed to find availability exception with relationships")

	if savedException.Staff.ID == uuid.Nil {
		t.Logf("Staff not loaded, savedException.StaffID: %v", savedException.StaffID)
	}
	if savedException.Business.ID == uuid.Nil {
		t.Logf("Business not loaded, savedException.BusinessID: %v", savedException.BusinessID)
	}

	assert.Equal(t, staffID, savedException.Staff.ID)
	assert.Equal(t, "Hair Stylist", savedException.Staff.Position)
	assert.Equal(t, businessID, savedException.Business.ID)
	assert.Equal(t, "Availability Salon", savedException.Business.DisplayName)

	// Create a recurring exception (lunch break)
	nextWeek := time.Now().Add(7 * 24 * time.Hour).Truncate(24 * time.Hour)
	recurringException := models.AvailabilityException{
		BaseModel: models.BaseModel{
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		BusinessID:     businessID,
		StaffID:        staffID,
		ExceptionType:  models.ExceptionTypeCustomHours,
		StartTime:      nextWeek.Add(12 * time.Hour), // 12:00 PM
		EndTime:        nextWeek.Add(13 * time.Hour), // 1:00 PM
		IsFullDay:      false,
		IsRecurring:    true,
		RecurrenceRule: "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR", // Every weekday
		Notes:          "Lunch break",
	}

	// Save the recurring exception
	err = testDB.DB.Create(&recurringException).Error
	assert.NoError(t, err, "Failed to create recurring availability exception")

	// Verify that we now have two exceptions for the staff
	var exceptions []models.AvailabilityException
	err = testDB.DB.Where("staff_id = ?", staffID).Find(&exceptions).Error
	assert.NoError(t, err, "Failed to find availability exceptions")
	assert.Len(t, exceptions, 2, "Should have two availability exceptions")

	// Create a holiday exception (for multiple staff, but we'll test with just one)
	christmasDay := time.Date(time.Now().Year(), 12, 25, 0, 0, 0, 0, time.Local)
	holidayException := models.AvailabilityException{
		BaseModel: models.BaseModel{
			CreatedBy: &userID,
			UpdatedBy: &userID,
		},
		BusinessID:     businessID,
		StaffID:        staffID,
		ExceptionType:  models.ExceptionTypeHoliday,
		StartTime:      christmasDay,
		EndTime:        christmasDay.Add(24 * time.Hour),
		IsFullDay:      true,
		IsRecurring:    true,
		RecurrenceRule: "FREQ=YEARLY", // Every year
		Notes:          "Christmas Day",
	}

	// Save the holiday exception
	err = testDB.DB.Create(&holidayException).Error
	assert.NoError(t, err, "Failed to create holiday exception")

	// Test querying exceptions by date range
	startDate := tomorrow
	endDate := tomorrow.Add(24 * time.Hour)
	var dateRangeExceptions []models.AvailabilityException
	err = testDB.DB.Where("staff_id = ? AND start_time >= ? AND end_time <= ?", staffID, startDate, endDate).Find(&dateRangeExceptions).Error
	assert.NoError(t, err, "Failed to query exceptions by date range")
	assert.Len(t, dateRangeExceptions, 1, "Should have one exception in date range")
	assert.Equal(t, exceptionID, dateRangeExceptions[0].ID)

	// Test updating exception
	err = testDB.DB.Model(&exception).Updates(map[string]interface{}{
		"notes":       "Updated personal day off",
		"is_full_day": false,
		"start_time":  tomorrow.Add(12 * time.Hour), // 12:00 PM
		"end_time":    tomorrow.Add(17 * time.Hour), // 5:00 PM
	}).Error
	assert.NoError(t, err, "Failed to update availability exception")

	err = testDB.DB.First(&savedException, "id = ?", exceptionID).Error
	assert.NoError(t, err, "Failed to find updated availability exception")
	assert.Equal(t, "Updated personal day off", savedException.Notes)
	assert.False(t, savedException.IsFullDay)
	assert.Equal(t, tomorrow.Add(12*time.Hour).Format(time.RFC3339), savedException.StartTime.Format(time.RFC3339))

	// Test soft delete
	err = testDB.DB.Delete(&exception).Error
	assert.NoError(t, err, "Failed to soft delete availability exception")

	// Verify exception is soft deleted
	var deletedException models.AvailabilityException
	err = testDB.DB.Unscoped().First(&deletedException, "id = ?", exceptionID).Error
	assert.NoError(t, err, "Failed to find soft deleted availability exception")
	assert.False(t, deletedException.DeletedAt.Time.IsZero(), "DeletedAt should be set")

	// Verify we can't find the exception with normal queries
	err = testDB.DB.First(&models.AvailabilityException{}, "id = ?", exceptionID).Error
	assert.Error(t, err, "Should not find soft deleted availability exception")

	// Test that soft deleting a staff member doesn't cascade delete exceptions
	err = testDB.DB.Delete(&staff).Error
	assert.NoError(t, err, "Failed to soft delete staff")

	// Check that the recurring exception still exists (but can't be found with normal queries because of staff FK)
	var exceptionCount int64
	err = testDB.DB.Unscoped().Model(&models.AvailabilityException{}).Where("id = ?", recurringException.ID).Count(&exceptionCount).Error
	assert.NoError(t, err, "Failed to count exceptions")
	assert.Equal(t, int64(1), exceptionCount, "Exception should still exist")
}
