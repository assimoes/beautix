package tests

import (
	"context"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppointmentRepositoryIntegration_Create(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "appointment_create@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "create-appointment-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "create_appt_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "create_appt_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")
	
	// Create a new appointment
	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(time.Hour)

	appointment := &domain.Appointment{
		BusinessID: business.BusinessID,
		ClientID:   client.ID,
		StaffID:    staff.StaffID,
		ServiceID:  service.ID,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     "scheduled",
		Notes:      "Test appointment notes",
		CreatedBy:  &user.UserID,
	}

	// Test creation
	err = appointmentRepo.Create(ctx, appointment)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, appointment.ID)
	assert.NotZero(t, appointment.CreatedAt)

	// Verify the appointment was created
	result, err := appointmentRepo.GetByID(ctx, appointment.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, client.ID, result.ClientID)
	assert.Equal(t, staff.StaffID, result.StaffID)
	assert.Equal(t, service.ID, result.ServiceID)
	assert.Equal(t, "scheduled", result.Status)
	assert.Equal(t, "Test appointment notes", result.Notes)
	assert.Equal(t, user.UserID, *result.CreatedBy)
	assert.WithinDuration(t, startTime, result.StartTime, time.Second)
	assert.WithinDuration(t, endTime, result.EndTime, time.Second)
}

func TestAppointmentRepositoryIntegration_GetByID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "appointment_getbyid@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbyid-appointment-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "getbyid_appt_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "getbyid_appt_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")
	appointment := createTestAppointmentForService(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, service.ID, user.UserID)

	// Test GetByID
	result, err := appointmentRepo.GetByID(ctx, appointment.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, appointment.ID, result.ID)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, client.ID, result.ClientID)
	assert.Equal(t, staff.StaffID, result.StaffID)

	// Verify related entities are populated
	assert.NotNil(t, result.Client)
	assert.Equal(t, client.ID, result.Client.ID)
	
	assert.NotNil(t, result.Business)
	assert.Equal(t, business.BusinessID, result.Business.BusinessID)
}

func TestAppointmentRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := appointmentRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestAppointmentRepositoryIntegration_Update(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "appointment_update@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "update-appointment-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "update_appt_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "update_appt_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")
	appointment := createTestAppointmentForService(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, service.ID, user.UserID)

	// Create update input
	updatedStatus := "confirmed"
	updatedNotes := "Updated appointment notes"
	newStartTime := time.Now().Add(48 * time.Hour)
	newEndTime := newStartTime.Add(time.Hour)

	updateInput := &domain.UpdateAppointmentInput{
		Status:    &updatedStatus,
		Notes:     &updatedNotes,
		StartTime: &newStartTime,
		EndTime:   &newEndTime,
	}

	// Test update
	err = appointmentRepo.Update(ctx, appointment.ID, updateInput, user.UserID)
	require.NoError(t, err)

	// Verify the update
	updatedAppointment, err := appointmentRepo.GetByID(ctx, appointment.ID)
	require.NoError(t, err)
	assert.Equal(t, updatedStatus, updatedAppointment.Status)
	assert.Equal(t, updatedNotes, updatedAppointment.Notes)
	assert.WithinDuration(t, newStartTime, updatedAppointment.StartTime, time.Second)
	assert.WithinDuration(t, newEndTime, updatedAppointment.EndTime, time.Second)
}

func TestAppointmentRepositoryIntegration_Delete(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "appointment_delete@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "delete-appointment-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "delete_appt_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "delete_appt_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")
	appointment := createTestAppointmentForService(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, service.ID, user.UserID)

	// Test delete
	err = appointmentRepo.Delete(ctx, appointment.ID, user.UserID)
	require.NoError(t, err)

	// Verify the appointment is deleted (not found)
	deletedAppointment, err := appointmentRepo.GetByID(ctx, appointment.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedAppointment)
	assert.Contains(t, err.Error(), "not found")
}

func TestAppointmentRepositoryIntegration_ListByBusiness(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "appointment_listbybusiness@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "listbybusiness-appointment-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "listbybusiness_appt_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "listbybusiness_appt_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")

	// Create multiple appointments
	appointment1 := createTestAppointmentForService(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, service.ID, user.UserID)
	appointment2 := createTestAppointmentForService(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, service.ID, user.UserID)

	// Test ListByBusiness
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now().Add(48 * time.Hour)
	results, err := appointmentRepo.ListByBusiness(ctx, business.BusinessID, startDate, endDate, 1, 10)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify all appointments belong to the business
	for _, appt := range results {
		assert.Equal(t, business.BusinessID, appt.BusinessID)
	}

	// Verify we got our appointments
	appointmentIDs := make([]uuid.UUID, len(results))
	for i, appt := range results {
		appointmentIDs[i] = appt.ID
	}
	assert.Contains(t, appointmentIDs, appointment1.ID)
	assert.Contains(t, appointmentIDs, appointment2.ID)
}

func TestAppointmentRepositoryIntegration_ListByClient(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	clientRepo := &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "appointment_listbyclient@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "listbyclient-appointment-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "listbyclient_appt_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "listbyclient_appt_staff@test.com")
	category := createTestServiceCategory(t, categoryRepo, business.BusinessID, "Hair Services")
	service := createTestServiceForAppointment(t, serviceRepo, business.BusinessID, &category.ID, user.UserID, "Hair Cut")

	// Create multiple appointments for the client
	appointment1 := createTestAppointmentForService(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, service.ID, user.UserID)
	appointment2 := createTestAppointmentForService(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, service.ID, user.UserID)

	// Test ListByClient
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now().Add(48 * time.Hour)
	results, err := appointmentRepo.ListByClient(ctx, client.ID, startDate, endDate, 1, 10)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify all appointments belong to the client
	for _, appt := range results {
		assert.Equal(t, client.ID, appt.ClientID)
	}

	// Verify we got our appointments
	appointmentIDs := make([]uuid.UUID, len(results))
	for i, appt := range results {
		appointmentIDs[i] = appt.ID
	}
	assert.Contains(t, appointmentIDs, appointment1.ID)
	assert.Contains(t, appointmentIDs, appointment2.ID)
}

// Helper functions for creating test data

func createTestClientForAppointment(t *testing.T, clientRepo *repository.ClientRepository, businessID, userID uuid.UUID, email string) *domain.Client {
	client := &domain.Client{
		BusinessID: businessID,
		FirstName:  "Test",
		LastName:   "Client",
		Email:      email,
		Phone:      "+351123456789",
		CreatedBy:  &userID,
	}
	
	err := clientRepo.Create(context.Background(), client)
	require.NoError(t, err)
	
	return client
}

func createTestStaffForAppointment(t *testing.T, staffRepo *repository.StaffRepository, businessID, userID uuid.UUID, email string) *domain.Staff {
	staff := &domain.Staff{
		BusinessID: businessID,
		UserID:     userID,
		Position:   "Hairstylist",
		IsActive:   true,
		CreatedBy:  userID,
	}
	
	err := staffRepo.Create(context.Background(), staff)
	require.NoError(t, err)
	
	return staff
}

func createTestServiceForAppointment(t *testing.T, serviceRepo *repository.ServiceRepository, businessID uuid.UUID, categoryID *uuid.UUID, userID uuid.UUID, name string) *domain.Service {
	service := &domain.Service{
		BusinessID:  businessID,
		CategoryID:  categoryID,
		Name:        name,
		Description: "Professional hair cut",
		Duration:    45,
		Price:       50.00,
		CreatedBy:   &userID,
	}

	err := serviceRepo.Create(context.Background(), service)
	require.NoError(t, err)
	return service
}

// Helper function to create a test appointment with incrementing times
var appointmentCounter int = 0
func createTestAppointmentForService(t *testing.T, appointmentRepo *repository.AppointmentRepository, businessID, clientID, staffID, serviceID, userID uuid.UUID) *domain.Appointment {
	appointmentCounter++
	startTime := time.Now().Add(time.Duration(24+appointmentCounter*2) * time.Hour)
	endTime := startTime.Add(time.Hour)
	
	appointment := &domain.Appointment{
		BusinessID: businessID,
		ClientID:   clientID,
		StaffID:    staffID,
		ServiceID:  serviceID,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     "scheduled",
		Notes:      "Test appointment",
		CreatedBy:  &userID,
	}
	
	err := appointmentRepo.Create(context.Background(), appointment)
	require.NoError(t, err)
	
	return appointment
}