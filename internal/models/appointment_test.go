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
	appointment := Appointment{
		BusinessID:    uuid.New(),
		ClientID:      uuid.New(),
		StartTime:     time.Now(),
		EndTime:       time.Now().Add(1 * time.Hour),
		Status:        AppointmentStatusConfirmed,
		TotalDuration: 60,
		TotalPrice:    100.00,
		Discount:      10.00,
		Tax:           20.00,
		FinalPrice:    110.00,
		PaymentStatus: PaymentStatusPending,
		BookingSource: BookingSourceWebsite,
		Notifications: NotificationSettings{
			ClientReminded: true,
			ConfirmationSent: true,
		},
	}

	assert.Equal(t, AppointmentStatusConfirmed, appointment.Status)
	assert.Equal(t, PaymentStatusPending, appointment.PaymentStatus)
	assert.Equal(t, BookingSourceWebsite, appointment.BookingSource)
	assert.Equal(t, 60, appointment.TotalDuration)
	assert.Equal(t, 110.00, appointment.FinalPrice)
	assert.True(t, appointment.Notifications.ClientReminded)
}

func TestAppointmentService(t *testing.T) {
	// Test AppointmentService struct
	selectedOptions := AppointmentServiceOptions{
		{
			OptionID:        uuid.New(),
			OptionName:      "Color",
			ChoiceID:        uuid.New(),
			ChoiceName:      "Red",
			PriceAdjustment: 15.00,
			TimeAdjustment:  10,
		},
		{
			OptionID:        uuid.New(),
			OptionName:      "Treatment",
			ChoiceID:        uuid.New(),
			ChoiceName:      "Deep Conditioning",
			PriceAdjustment: 20.00,
			TimeAdjustment:  15,
		},
	}

	appointmentService := AppointmentService{
		AppointmentID:    uuid.New(),
		ServiceID:        uuid.New(),
		StaffID:          uuid.New(),
		Price:            75.00,
		Discount:         5.00,
		Duration:         45,
		StartTime:        time.Now(),
		EndTime:          time.Now().Add(45 * time.Minute),
		Status:           AppointmentStatusConfirmed,
		SelectedOptions:  selectedOptions,
	}

	assert.Equal(t, 75.00, appointmentService.Price)
	assert.Equal(t, 45, appointmentService.Duration)
	assert.Equal(t, 2, len(appointmentService.SelectedOptions))
	assert.Equal(t, "Red", appointmentService.SelectedOptions[0].ChoiceName)
	assert.Equal(t, 15.00, appointmentService.SelectedOptions[0].PriceAdjustment)
}

func TestAppointmentPayment(t *testing.T) {
	// Test AppointmentPayment struct
	payment := AppointmentPayment{
		AppointmentID: uuid.New(),
		Amount:        110.00,
		PaymentMethod: PaymentMethodCard,
		PaymentStatus: PaymentStatusPaid,
		TransactionID: "txn_123456789",
		PaymentDate:   time.Now(),
		ReceiptURL:    "https://example.com/receipts/123456789.pdf",
	}

	assert.Equal(t, 110.00, payment.Amount)
	assert.Equal(t, PaymentMethodCard, payment.PaymentMethod)
	assert.Equal(t, PaymentStatusPaid, payment.PaymentStatus)
	assert.Equal(t, "txn_123456789", payment.TransactionID)
}

func TestAppointmentForm(t *testing.T) {
	// Test AppointmentForm struct
	formData := FormData{
		"allergies":   []string{"Latex", "Peanuts"},
		"medications": []string{"Aspirin"},
		"consent":     true,
	}

	form := AppointmentForm{
		AppointmentID: uuid.New(),
		FormType:      "medical",
		FormData:      formData,
		IsCompleted:   true,
		CompletedAt:   func() *time.Time { now := time.Now(); return &now }(),
	}

	assert.Equal(t, "medical", form.FormType)
	assert.True(t, form.IsCompleted)
	assert.Equal(t, []string{"Latex", "Peanuts"}, form.FormData["allergies"])
}

func TestAppointmentFeedback(t *testing.T) {
	// Test AppointmentFeedback struct
	feedback := AppointmentFeedback{
		AppointmentID: uuid.New(),
		Rating:        5,
		Comments:      "Great service!",
		IsPublic:      true,
		FeedbackItems: FeedbackItems{
			ServiceQuality:      5,
			StaffProfessionalism: 5,
			Cleanliness:         4,
			ValueForMoney:       4,
			Atmosphere:          5,
			WouldRecommend:      true,
		},
	}

	assert.Equal(t, 5, feedback.Rating)
	assert.Equal(t, "Great service!", feedback.Comments)
	assert.True(t, feedback.IsPublic)
	assert.Equal(t, 5, feedback.FeedbackItems.ServiceQuality)
	assert.True(t, feedback.FeedbackItems.WouldRecommend)
}

func TestNotificationSettingsSerialization(t *testing.T) {
	// Test JSON serialization and deserialization
	now := time.Now()
	notifications := NotificationSettings{
		ClientReminded:    true,
		ClientReminderSent: now,
		ConfirmationSent:  true,
		ConfirmationTime:  now,
	}

	// Serialize to JSON
	jsonBytes, err := json.Marshal(notifications)
	assert.NoError(t, err)

	// Deserialize back
	var deserialized NotificationSettings
	err = json.Unmarshal(jsonBytes, &deserialized)
	assert.NoError(t, err)

	assert.Equal(t, notifications.ClientReminded, deserialized.ClientReminded)
	assert.Equal(t, notifications.ConfirmationSent, deserialized.ConfirmationSent)
}

func TestAppointmentServiceOptionsSerialization(t *testing.T) {
	// Test JSON serialization and deserialization for AppointmentServiceOptions
	options := AppointmentServiceOptions{
		{
			OptionID:        uuid.New(),
			OptionName:      "Length",
			ChoiceID:        uuid.New(),
			ChoiceName:      "Long",
			PriceAdjustment: 25.00,
			TimeAdjustment:  15,
		},
	}

	// Serialize to JSON
	jsonBytes, err := json.Marshal(options)
	assert.NoError(t, err)

	// Deserialize back
	var deserialized AppointmentServiceOptions
	err = json.Unmarshal(jsonBytes, &deserialized)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(deserialized))
	assert.Equal(t, options[0].OptionName, deserialized[0].OptionName)
	assert.Equal(t, options[0].ChoiceName, deserialized[0].ChoiceName)
	assert.Equal(t, options[0].PriceAdjustment, deserialized[0].PriceAdjustment)
}

func TestFormDataSerialization(t *testing.T) {
	// Test JSON serialization and deserialization for FormData
	formData := FormData{
		"name":    "John Doe",
		"age":     30,
		"history": []string{"Treatment A", "Treatment B"},
		"contact": map[string]interface{}{
			"email": "john@example.com",
			"phone": "123-456-7890",
		},
	}

	// Serialize to JSON
	jsonBytes, err := json.Marshal(formData)
	assert.NoError(t, err)

	// Deserialize back
	var deserialized FormData
	err = json.Unmarshal(jsonBytes, &deserialized)
	assert.NoError(t, err)

	assert.Equal(t, "John Doe", deserialized["name"])
	assert.Equal(t, float64(30), deserialized["age"])
	
	// Check the history array
	historyInterface, ok := deserialized["history"]
	assert.True(t, ok)
	history, ok := historyInterface.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(history))
	assert.Equal(t, "Treatment A", history[0])
	
	// Check the nested map
	contactInterface, ok := deserialized["contact"]
	assert.True(t, ok)
	contact, ok := contactInterface.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "john@example.com", contact["email"])
}

func TestFeedbackItemsSerialization(t *testing.T) {
	// Test JSON serialization and deserialization for FeedbackItems
	feedback := FeedbackItems{
		ServiceQuality:      5,
		StaffProfessionalism: 4,
		Cleanliness:         5,
		ValueForMoney:       3,
		Atmosphere:          4,
		WouldRecommend:      true,
	}

	// Serialize to JSON
	jsonBytes, err := json.Marshal(feedback)
	assert.NoError(t, err)

	// Deserialize back
	var deserialized FeedbackItems
	err = json.Unmarshal(jsonBytes, &deserialized)
	assert.NoError(t, err)

	assert.Equal(t, feedback.ServiceQuality, deserialized.ServiceQuality)
	assert.Equal(t, feedback.StaffProfessionalism, deserialized.StaffProfessionalism)
	assert.Equal(t, feedback.WouldRecommend, deserialized.WouldRecommend)
}