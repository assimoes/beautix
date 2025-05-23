package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAppointmentModel(t *testing.T) {
	// Test Appointment Struct
	businessID := uuid.New()
	clientID := uuid.New()
	staffID := uuid.New()
	serviceID := uuid.New()
	estimatedPrice := 100.00
	actualPrice := 110.00

	appointment := Appointment{
		BusinessID:         businessID,
		ClientID:           clientID,
		StaffID:            staffID,
		ServiceID:          serviceID,
		StartTime:          time.Now(),
		EndTime:            time.Now().Add(1 * time.Hour),
		Status:             AppointmentStatusConfirmed,
		Notes:              "Test appointment notes",
		EstimatedPrice:     &estimatedPrice,
		ActualPrice:        &actualPrice,
		PaymentStatus:      PaymentStatusPending,
		ClientConfirmed:    true,
		StaffConfirmed:     false,
		CancellationReason: "",
	}

	assert.Equal(t, AppointmentStatusConfirmed, appointment.Status)
	assert.Equal(t, PaymentStatusPending, appointment.PaymentStatus)
	assert.Equal(t, businessID, appointment.BusinessID)
	assert.Equal(t, clientID, appointment.ClientID)
	assert.Equal(t, staffID, appointment.StaffID)
	assert.Equal(t, serviceID, appointment.ServiceID)
	assert.Equal(t, 100.00, *appointment.EstimatedPrice)
	assert.Equal(t, 110.00, *appointment.ActualPrice)
	assert.True(t, appointment.ClientConfirmed)
	assert.False(t, appointment.StaffConfirmed)
	assert.Equal(t, 60, appointment.Duration()) // 1 hour = 60 minutes
}

func TestAppointmentStatus(t *testing.T) {
	tests := []struct {
		status   AppointmentStatus
		expected string
	}{
		{AppointmentStatusScheduled, "scheduled"},
		{AppointmentStatusConfirmed, "confirmed"},
		{AppointmentStatusInProgress, "in_progress"},
		{AppointmentStatusCompleted, "completed"},
		{AppointmentStatusCancelled, "cancelled"},
		{AppointmentStatusNoShow, "no_show"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, string(test.status))
	}
}

func TestPaymentStatus(t *testing.T) {
	tests := []struct {
		status   PaymentStatus
		expected string
	}{
		{PaymentStatusPending, "pending"},
		{PaymentStatusPaid, "paid"},
		{PaymentStatusPartial, "partial"},
		{PaymentStatusRefunded, "refunded"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, string(test.status))
	}
}

func TestPaymentMethod(t *testing.T) {
	tests := []struct {
		method   PaymentMethod
		expected string
	}{
		{PaymentMethodCash, "cash"},
		{PaymentMethodCard, "card"},
		{PaymentMethodTransfer, "transfer"},
		{PaymentMethodOther, "other"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, string(test.method))
	}
}

func TestAppointmentMethods(t *testing.T) {
	appointment := Appointment{
		Status:          AppointmentStatusCompleted,
		ClientConfirmed: true,
		StaffConfirmed:  true,
	}

	// Test IsCompleted
	assert.True(t, appointment.IsCompleted())

	// Test IsConfirmed
	assert.True(t, appointment.IsConfirmed())

	// Test IsCancelled
	assert.False(t, appointment.IsCancelled())

	// Change status to cancelled
	appointment.Status = AppointmentStatusCancelled
	assert.True(t, appointment.IsCancelled())
	assert.False(t, appointment.IsCompleted())

	// Test with incomplete confirmation
	appointment.StaffConfirmed = false
	assert.False(t, appointment.IsConfirmed())
}

func TestAppointmentDuration(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	endTime := startTime.Add(90 * time.Minute)

	appointment := Appointment{
		StartTime: startTime,
		EndTime:   endTime,
	}

	assert.Equal(t, 90, appointment.Duration())
}

func TestAppointmentJSON(t *testing.T) {
	estimatedPrice := 100.00
	paymentMethod := PaymentMethodCard

	appointment := Appointment{
		BusinessID:     uuid.New(),
		ClientID:       uuid.New(),
		StaffID:        uuid.New(),
		ServiceID:      uuid.New(),
		StartTime:      time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		EndTime:        time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
		Status:         AppointmentStatusConfirmed,
		Notes:          "Test appointment",
		EstimatedPrice: &estimatedPrice,
		PaymentMethod:  &paymentMethod,
		PaymentStatus:  PaymentStatusPaid,
	}

	// Test JSON marshaling
	jsonBytes, err := json.Marshal(appointment)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Test JSON unmarshaling
	var unmarshaled Appointment
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, appointment.BusinessID, unmarshaled.BusinessID)
	assert.Equal(t, appointment.Status, unmarshaled.Status)
	assert.Equal(t, appointment.PaymentStatus, unmarshaled.PaymentStatus)
	assert.Equal(t, appointment.Notes, unmarshaled.Notes)

	if appointment.EstimatedPrice != nil && unmarshaled.EstimatedPrice != nil {
		assert.Equal(t, *appointment.EstimatedPrice, *unmarshaled.EstimatedPrice)
	}

	if appointment.PaymentMethod != nil && unmarshaled.PaymentMethod != nil {
		assert.Equal(t, *appointment.PaymentMethod, *unmarshaled.PaymentMethod)
	}
}

func TestAppointmentTableName(t *testing.T) {
	appointment := Appointment{}
	assert.Equal(t, "appointments", appointment.TableName())
}

func TestAppointmentWithNilValues(t *testing.T) {
	appointment := Appointment{
		BusinessID:         uuid.New(),
		ClientID:           uuid.New(),
		StaffID:            uuid.New(),
		ServiceID:          uuid.New(),
		StartTime:          time.Now(),
		EndTime:            time.Now().Add(1 * time.Hour),
		Status:             AppointmentStatusScheduled,
		PaymentStatus:      PaymentStatusPending,
		EstimatedPrice:     nil, // Test nil pointer
		ActualPrice:        nil, // Test nil pointer
		PaymentMethod:      nil, // Test nil pointer
		ClientConfirmed:    false,
		StaffConfirmed:     false,
		CancellationReason: "",
	}

	// Should not panic with nil values
	assert.Equal(t, AppointmentStatusScheduled, appointment.Status)
	assert.Equal(t, PaymentStatusPending, appointment.PaymentStatus)
	assert.Nil(t, appointment.EstimatedPrice)
	assert.Nil(t, appointment.ActualPrice)
	assert.Nil(t, appointment.PaymentMethod)
	assert.False(t, appointment.IsConfirmed())
	assert.False(t, appointment.IsCompleted())
	assert.False(t, appointment.IsCancelled())
}
