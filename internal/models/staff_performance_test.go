package models_test

import (
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestStaffPerformanceModel(t *testing.T) {
	// Run test within transaction for true isolation
	database.NewTestDBWithTransaction(t, func(testDB *gorm.DB) {
		// Auto-migrate the models
		err := testDB.AutoMigrate(
			&models.User{},
			&models.Business{},
			&models.Staff{},
			&models.StaffPerformance{},
		)
		require.NoError(t, err, "Failed to migrate models")

		// Create a user
		userID := uuid.New()
		user := models.User{
			BaseModel: models.BaseModel{
				ID: userID,
			},
			ClerkID:   "clerk_performance_test",
			Email:     "performance_test@example.com",
			FirstName: "Performance",
			LastName:  "Test",
			Phone:     "+1234567890",
			Role:      models.UserRoleStaff,
			IsActive:  true,
		}

		// Save the user
		err = testDB.Create(&user).Error
		assert.NoError(t, err, "Failed to create user")

		// Create a business
		businessID := uuid.New()
		business := models.Business{
			BaseModel: models.BaseModel{
				ID: businessID,
			},
			UserID:           userID,
			Name:             "performance-salon",
			DisplayName:      "Performance Salon",
			Description:      "A salon for testing staff performance metrics",
			Address:          "123 Performance St",
			City:             "Lisbon",
			Country:          "Portugal",
			IsActive:         true,
		}

		// Save the business
		err = testDB.Create(&business).Error
		assert.NoError(t, err, "Failed to create business")

		// Create a staff member
		staffID := uuid.New()
		joinDate := time.Now().Add(-6 * 30 * 24 * time.Hour) // 6 months ago
		staff := models.Staff{
			BaseModel: models.BaseModel{
				ID: staffID,
			},
			BusinessID:     businessID,
			UserID:         userID,
			Position:       "Massage Therapist",
			Bio:            "Experienced massage therapist",
			SpecialtyAreas: models.SpecialtyAreas{"Swedish Massage", "Deep Tissue", "Hot Stone"},
			IsActive:       true,
			EmploymentType: models.StaffEmploymentTypeFull,
			JoinDate:       joinDate,
			CommissionRate: 25.00,
		}

		// Save the staff
		err = testDB.Create(&staff).Error
		assert.NoError(t, err, "Failed to create staff")

		// Create a staff performance record for April 2025
		performanceID := uuid.New()
		performance := models.StaffPerformance{
			BaseModel: models.BaseModel{
				ID: performanceID,
			},
			BusinessID:            businessID,
			StaffID:               staffID,
			Period:                models.PerformancePeriodMonthly,
			StartDate:             time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			EndDate:               time.Date(2025, 4, 30, 23, 59, 59, 999000000, time.UTC),
			TotalAppointments:     50,
			CompletedAppointments: 45,
			CanceledAppointments:  3,
			NoShowAppointments:    2,
			TotalRevenue:          3750.50,
			AverageRating:         4.7,
			ClientRetentionRate:   88.5,
			NewClients:            12,
			ReturnClients:         38,
		}

		// Save the performance record
		err = testDB.Create(&performance).Error
		assert.NoError(t, err, "Failed to create staff performance")

		// Verify performance was created with ID
		var savedPerformance models.StaffPerformance
		err = testDB.First(&savedPerformance, "id = ?", performanceID).Error
		assert.NoError(t, err, "Failed to find staff performance")
		assert.Equal(t, performanceID, savedPerformance.ID)
		assert.Equal(t, businessID, savedPerformance.BusinessID)
		assert.Equal(t, staffID, savedPerformance.StaffID)
		assert.Equal(t, models.PerformancePeriodMonthly, savedPerformance.Period)
		assert.Equal(t, 50, savedPerformance.TotalAppointments)
		assert.Equal(t, 45, savedPerformance.CompletedAppointments)
		assert.Equal(t, 3, savedPerformance.CanceledAppointments)
		assert.Equal(t, 2, savedPerformance.NoShowAppointments)
		assert.Equal(t, 3750.50, savedPerformance.TotalRevenue)
		assert.Equal(t, 4.7, savedPerformance.AverageRating)
		assert.Equal(t, 88.5, savedPerformance.ClientRetentionRate)
		assert.Equal(t, 12, savedPerformance.NewClients)
		assert.Equal(t, 38, savedPerformance.ReturnClients)

		// Test loaded relationships
		err = testDB.Preload("Staff").Preload("Business").First(&savedPerformance, "id = ?", performanceID).Error
		assert.NoError(t, err, "Failed to find staff performance with relationships")
		assert.Equal(t, staffID, savedPerformance.Staff.ID)
		assert.Equal(t, "Massage Therapist", savedPerformance.Staff.Position)
		assert.Equal(t, businessID, savedPerformance.Business.ID)
		assert.Equal(t, "Performance Salon", savedPerformance.Business.DisplayName)

		// Create a weekly performance record
		// Calculate first day of week
		currentTime := time.Now()
		weekday := int(currentTime.Weekday())
		if weekday == 0 { // If today is Sunday
			weekday = 7
		}
		firstDayOfWeek := currentTime.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
		lastDayOfWeek := firstDayOfWeek.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

		weeklyPerformance := models.StaffPerformance{
			BusinessID:            businessID,
			StaffID:               staffID,
			Period:                models.PerformancePeriodWeekly,
			StartDate:             firstDayOfWeek,
			EndDate:               lastDayOfWeek,
			TotalAppointments:     12,
			CompletedAppointments: 10,
			CanceledAppointments:  1,
			NoShowAppointments:    1,
			TotalRevenue:          825.00,
			AverageRating:         4.5,
			ClientRetentionRate:   75.0,
			NewClients:            3,
			ReturnClients:         9,
		}

		// Save the weekly performance record
		err = testDB.Create(&weeklyPerformance).Error
		assert.NoError(t, err, "Failed to create weekly staff performance")

		// Verify that we now have two performance records for the staff
		var performances []models.StaffPerformance
		err = testDB.Where("staff_id = ?", staffID).Find(&performances).Error
		assert.NoError(t, err, "Failed to find staff performances")
		assert.Len(t, performances, 2, "Should have two staff performance records")

		// Test querying performance by period
		var monthlyPerformances []models.StaffPerformance
		err = testDB.Where("staff_id = ? AND period = ?", staffID, models.PerformancePeriodMonthly).Find(&monthlyPerformances).Error
		assert.NoError(t, err, "Failed to query performances by period")
		assert.Len(t, monthlyPerformances, 1, "Should have one monthly performance record")
		assert.Equal(t, performanceID, monthlyPerformances[0].ID)

		// Test updating performance metrics
		err = testDB.Model(&performance).Updates(map[string]interface{}{
			"total_appointments":      52,
			"completed_appointments":  48,
			"canceled_appointments":   2,
			"no_show_appointments":    2,
			"total_revenue":           4000.00,
			"average_rating":          4.8,
			"client_retention_rate":   90.0,
			"new_clients":             14,
			"return_clients":          38,
		}).Error
		assert.NoError(t, err, "Failed to update staff performance")

		err = testDB.First(&savedPerformance, "id = ?", performanceID).Error
		assert.NoError(t, err, "Failed to find updated staff performance")
		assert.Equal(t, 52, savedPerformance.TotalAppointments)
		assert.Equal(t, 48, savedPerformance.CompletedAppointments)
		assert.Equal(t, 4000.00, savedPerformance.TotalRevenue)
		assert.Equal(t, 4.8, savedPerformance.AverageRating)

		// Test soft delete
		err = testDB.Delete(&performance).Error
		assert.NoError(t, err, "Failed to soft delete staff performance")

		// Verify performance is soft deleted
		var deletedPerformance models.StaffPerformance
		err = testDB.Unscoped().First(&deletedPerformance, "id = ?", performanceID).Error
		assert.NoError(t, err, "Failed to find soft deleted staff performance")
		assert.False(t, deletedPerformance.DeletedAt.Time.IsZero(), "DeletedAt should be set")

		// Verify we can't find the performance with normal queries
		err = testDB.First(&models.StaffPerformance{}, "id = ?", performanceID).Error
		assert.Error(t, err, "Should not find soft deleted staff performance")

		// Test that soft deleting a staff member doesn't cascade delete performance records
		err = testDB.Delete(&staff).Error
		assert.NoError(t, err, "Failed to soft delete staff")

		// Check that the weekly performance record still exists (but can't be found with normal queries because of staff FK)
		var performanceCount int64
		err = testDB.Unscoped().Model(&models.StaffPerformance{}).Where("id = ?", weeklyPerformance.ID).Count(&performanceCount).Error
		assert.NoError(t, err, "Failed to count performance records")
		assert.Equal(t, int64(1), performanceCount, "Performance record should still exist")
	})
}