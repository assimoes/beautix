package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStaffModel(t *testing.T) {
	// Connect to the test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate the models
	err = testDB.AutoMigrate(&models.User{}, &models.Business{}, &models.Staff{}, &models.ServiceAssignment{}, &models.AvailabilityException{}, &models.StaffPerformance{})
	require.NoError(t, err, "Failed to migrate models")

	// Create a user first (since staff requires a user)
	userID := uuid.New()
	user := models.User{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		ClerkID:   "clerk_staff_test",
		Email:     "staff_test@example.com",
		FirstName: "Staff",
		LastName:  "Member",
		Phone:     "+1234567890",
		Role:      models.UserRoleStaff,
		IsActive:  true,
	}

	// Save the user
	err = testDB.Create(&user).Error
	assert.NoError(t, err, "Failed to create user")

	// Create a business (since staff requires a business)
	businessID := uuid.New()
	business := models.Business{
		BaseModel: models.BaseModel{
			ID: businessID,
		},
		UserID:           userID,
		Name:             "staff-salon-1",
		DisplayName:      "Staff Salon Example",
		Description:      "A salon for testing staff models",
		Address:          "123 Main St",
		City:             "Lisbon",
		Country:          "Portugal",
		PostalCode:       "1000-100",
		Phone:            "+351123456789",
		Email:            "contact@staffsalon.com",
		SubscriptionTier: models.SubscriptionTierPro,
		IsActive:         true,
	}

	// Save the business
	err = testDB.Create(&business).Error
	assert.NoError(t, err, "Failed to create business")

	// Create staff
	staffID := uuid.New()
	joinDate := time.Now().Add(-30 * 24 * time.Hour) // 30 days ago
	staff := models.Staff{
		BaseModel: models.BaseModel{
			ID: staffID,
		},
		BusinessID:     businessID,
		UserID:         userID,
		Position:       "Senior Stylist",
		Bio:            "An experienced hair stylist with 5 years of experience",
		SpecialtyAreas: models.SpecialtyAreas{"Hair Coloring", "Hair Cutting", "Styling"},
		ProfileImageURL: "https://example.com/staff/profile.jpg",
		IsActive:       true,
		EmploymentType: models.StaffEmploymentTypeFull,
		JoinDate:       joinDate,
		CommissionRate: 25.00, // 25%
	}

	// Save the staff
	err = testDB.Create(&staff).Error
	assert.NoError(t, err, "Failed to create staff")

	// Verify staff was created with ID
	var savedStaff models.Staff
	err = testDB.First(&savedStaff, "id = ?", staffID).Error
	assert.NoError(t, err, "Failed to find staff")
	assert.Equal(t, staffID, savedStaff.ID)
	assert.Equal(t, "Senior Stylist", savedStaff.Position)
	assert.Equal(t, models.StaffEmploymentTypeFull, savedStaff.EmploymentType)
	assert.Equal(t, 25.00, savedStaff.CommissionRate)

	// Test JSONB field
	assert.Len(t, savedStaff.SpecialtyAreas, 3)
	assert.Contains(t, savedStaff.SpecialtyAreas, "Hair Coloring")

	// Test loaded relationships
	err = testDB.Preload("User").Preload("Business").First(&savedStaff, "id = ?", staffID).Error
	assert.NoError(t, err, "Failed to find staff with relationships")
	assert.Equal(t, userID, savedStaff.User.ID)
	assert.Equal(t, "Staff", savedStaff.User.FirstName)
	assert.Equal(t, businessID, savedStaff.Business.ID)
	assert.Equal(t, "Staff Salon Example", savedStaff.Business.DisplayName)

	// Create availability exception
	exceptionID := uuid.New()
	tomorrow := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	exception := models.AvailabilityException{
		BaseModel: models.BaseModel{
			ID: exceptionID,
		},
		BusinessID:    businessID,
		StaffID:       staffID,
		ExceptionType: models.ExceptionTypeTimeOff,
		StartTime:     tomorrow.Add(9 * time.Hour),  // 9:00 AM
		EndTime:       tomorrow.Add(17 * time.Hour), // 5:00 PM
		IsFullDay:     true,
		Notes:         "Personal day off",
	}

	// Save the exception
	err = testDB.Create(&exception).Error
	assert.NoError(t, err, "Failed to create availability exception")

	// Verify exception was created
	var savedException models.AvailabilityException
	err = testDB.Preload("Staff").First(&savedException, "id = ?", exceptionID).Error
	assert.NoError(t, err, "Failed to find availability exception")
	assert.Equal(t, exceptionID, savedException.ID)
	assert.Equal(t, models.ExceptionTypeTimeOff, savedException.ExceptionType)
	assert.Equal(t, staffID, savedException.Staff.ID)
	assert.True(t, savedException.IsFullDay)

	// Create staff performance record
	performanceID := uuid.New()
	firstDayOfMonth := time.Now().Truncate(24 * time.Hour)
	firstDayOfMonth = time.Date(firstDayOfMonth.Year(), firstDayOfMonth.Month(), 1, 0, 0, 0, 0, firstDayOfMonth.Location())
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	
	performance := models.StaffPerformance{
		BaseModel: models.BaseModel{
			ID: performanceID,
		},
		BusinessID:            businessID,
		StaffID:               staffID,
		Period:                models.PerformancePeriodMonthly,
		StartDate:             firstDayOfMonth,
		EndDate:               lastDayOfMonth,
		TotalAppointments:     45,
		CompletedAppointments: 40,
		CanceledAppointments:  3,
		NoShowAppointments:    2,
		TotalRevenue:          2500.50,
		AverageRating:         4.8,
		ClientRetentionRate:   85.5,
		NewClients:            15,
		ReturnClients:         30,
	}

	// Save the performance record
	err = testDB.Create(&performance).Error
	assert.NoError(t, err, "Failed to create staff performance")

	// Verify performance was created
	var savedPerformance models.StaffPerformance
	err = testDB.Preload("Staff").First(&savedPerformance, "id = ?", performanceID).Error
	assert.NoError(t, err, "Failed to find staff performance")
	assert.Equal(t, performanceID, savedPerformance.ID)
	assert.Equal(t, models.PerformancePeriodMonthly, savedPerformance.Period)
	assert.Equal(t, 45, savedPerformance.TotalAppointments)
	assert.Equal(t, 4.8, savedPerformance.AverageRating)
	assert.Equal(t, staffID, savedPerformance.Staff.ID)

	// Test serialization and deserialization of JSON fields
	jsonBytes, err := json.Marshal(staff)
	assert.NoError(t, err, "Failed to marshal staff to JSON")

	var unmarshaledStaff models.Staff
	err = json.Unmarshal(jsonBytes, &unmarshaledStaff)
	assert.NoError(t, err, "Failed to unmarshal staff from JSON")
	assert.Equal(t, staff.Position, unmarshaledStaff.Position)
	assert.Equal(t, staff.SpecialtyAreas, unmarshaledStaff.SpecialtyAreas)

	// Test soft delete
	err = testDB.Delete(&staff).Error
	assert.NoError(t, err, "Failed to soft delete staff")

	// Verify staff is soft deleted
	var deletedStaff models.Staff
	err = testDB.Unscoped().First(&deletedStaff, "id = ?", staffID).Error
	assert.NoError(t, err, "Failed to find soft deleted staff")
	assert.False(t, deletedStaff.DeletedAt.Time.IsZero(), "DeletedAt should be set")

	// Verify we can't find the staff with normal queries
	err = testDB.First(&models.Staff{}, "id = ?", staffID).Error
	assert.Error(t, err, "Should not find soft deleted staff")

	// Verify cascade effect - exceptions and performance should still exist
	// since they are related to business as well and may be needed for historical purposes
	var exceptionCount int64
	err = testDB.Model(&models.AvailabilityException{}).Where("staff_id = ?", staffID).Count(&exceptionCount).Error
	assert.NoError(t, err, "Failed to count exceptions")
	assert.Equal(t, int64(1), exceptionCount, "Exception should still exist")

	var performanceCount int64
	err = testDB.Model(&models.StaffPerformance{}).Where("staff_id = ?", staffID).Count(&performanceCount).Error
	assert.NoError(t, err, "Failed to count performance records")
	assert.Equal(t, int64(1), performanceCount, "Performance record should still exist")
}