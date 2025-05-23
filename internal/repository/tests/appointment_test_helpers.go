package tests

import (
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// createTestAppointmentDBTx creates a test appointment using the actual database schema within a transaction
func createTestAppointmentDBTx(t *testing.T, tx *gorm.DB, businessID, clientID, staffID, serviceID uuid.UUID, startTime time.Time, createdByID uuid.UUID) *models.Appointment {
	// Calculate end time (1 hour after start time)
	endTime := startTime.Add(time.Hour)

	appointment := &models.Appointment{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			CreatedBy: &createdByID,
		},
		BusinessID: businessID,
		ClientID:   clientID,
		StaffID:    staffID,
		ServiceID:  serviceID,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     models.AppointmentStatusScheduled,
		Notes:      "Test appointment",
	}

	err := tx.Create(appointment).Error
	require.NoError(t, err, "Failed to create test appointment")

	return appointment
}

// createTestServiceCompletionDBTx creates a test service completion using the actual database schema within a transaction
func createTestServiceCompletionDBTx(t *testing.T, tx *gorm.DB, appointmentID uuid.UUID, priceCharged float64, createdByID uuid.UUID) *models.ServiceCompletion {
	completion := &models.ServiceCompletion{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			CreatedBy: &createdByID,
		},
		AppointmentID:     appointmentID,
		PriceCharged:      priceCharged,
		PaymentMethod:     "card",
		ProviderConfirmed: true,
		ClientConfirmed:   true,
		CompletionDate:    func() *time.Time { t := time.Now(); return &t }(),
	}

	err := tx.Create(completion).Error
	require.NoError(t, err, "Failed to create test service completion")

	return completion
}
