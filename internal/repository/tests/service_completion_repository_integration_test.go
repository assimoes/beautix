package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceCompletionRepositoryIntegration_Create(t *testing.T) {
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
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceCompletionRepo := &repository.ServiceCompletionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_completion_create@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "create-service-completion-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "create_svc_completion_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "create_svc_completion_staff@test.com")
	
	// Create test appointment
	startTime := time.Now().Add(24 * time.Hour)
	appointment := createTestAppointmentForCompletion(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, user.UserID, startTime)

	// Create a new service completion
	completionDate := time.Now()
	completion := &domain.ServiceCompletion{
		AppointmentID:     appointment.ID,
		PriceCharged:      125.50,
		PaymentMethod:     "card",
		ProviderConfirmed: true,
		ClientConfirmed:   false,
		CompletionDate:    &completionDate,
		CreatedBy:         &user.UserID,
	}

	// Test creation
	err = serviceCompletionRepo.Create(ctx, completion)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, completion.ID)
	assert.NotZero(t, completion.CreatedAt)

	// Verify the completion was created
	result, err := serviceCompletionRepo.GetByID(ctx, completion.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, appointment.ID, result.AppointmentID)
	assert.Equal(t, 125.50, result.PriceCharged)
	assert.Equal(t, "card", result.PaymentMethod)
	assert.True(t, result.ProviderConfirmed)
	assert.False(t, result.ClientConfirmed)
	assert.Equal(t, user.UserID, *result.CreatedBy)
}

func TestServiceCompletionRepositoryIntegration_GetByID(t *testing.T) {
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
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceCompletionRepo := &repository.ServiceCompletionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_completion_getbyid@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbyid-service-completion-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "getbyid_svc_completion_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "getbyid_svc_completion_staff@test.com")
	
	// Create test appointment
	startTime := time.Now().Add(24 * time.Hour)
	appointment := createTestAppointmentForCompletion(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, user.UserID, startTime)
	
	// Create test service completion
	completion := createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment.ID, 100.0, user.UserID)

	// Test GetByID
	result, err := serviceCompletionRepo.GetByID(ctx, completion.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, completion.ID, result.ID)
	assert.Equal(t, appointment.ID, result.AppointmentID)
	assert.Equal(t, 100.0, result.PriceCharged)
	
	// Verify related appointment is populated
	assert.NotNil(t, result.Appointment)
	assert.Equal(t, appointment.ID, result.Appointment.ID)
}

func TestServiceCompletionRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	serviceCompletionRepo := &repository.ServiceCompletionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}

	// Test GetByID with non-existent ID
	ctx := context.Background()
	result, err := serviceCompletionRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestServiceCompletionRepositoryIntegration_GetByAppointmentID(t *testing.T) {
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
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceCompletionRepo := &repository.ServiceCompletionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_completion_getbyappointment@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "getbyappointment-service-completion-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "getbyappointment_svc_completion_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "getbyappointment_svc_completion_staff@test.com")
	
	// Create test appointment
	startTime := time.Now().Add(24 * time.Hour)
	appointment := createTestAppointmentForCompletion(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, user.UserID, startTime)
	
	// Create test service completion
	completion := createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment.ID, 150.0, user.UserID)

	// Test GetByAppointmentID
	result, err := serviceCompletionRepo.GetByAppointmentID(ctx, appointment.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, completion.ID, result.ID)
	assert.Equal(t, appointment.ID, result.AppointmentID)
	assert.Equal(t, 150.0, result.PriceCharged)
}

func TestServiceCompletionRepositoryIntegration_Update(t *testing.T) {
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
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceCompletionRepo := &repository.ServiceCompletionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_completion_update@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "update-service-completion-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "update_svc_completion_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "update_svc_completion_staff@test.com")
	
	// Create test appointment
	startTime := time.Now().Add(24 * time.Hour)
	appointment := createTestAppointmentForCompletion(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, user.UserID, startTime)
	
	// Create test service completion
	completion := createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment.ID, 100.0, user.UserID)

	// Create update input
	newPriceCharged := 175.50
	newPaymentMethod := "cash"
	providerConfirmed := true
	clientConfirmed := true
	updateInput := &domain.UpdateServiceCompletionInput{
		PriceCharged:      &newPriceCharged,
		PaymentMethod:     &newPaymentMethod,
		ProviderConfirmed: &providerConfirmed,
		ClientConfirmed:   &clientConfirmed,
	}

	// Test update
	err = serviceCompletionRepo.Update(ctx, completion.ID, updateInput, user.UserID)
	require.NoError(t, err)

	// Verify the update
	updatedCompletion, err := serviceCompletionRepo.GetByID(ctx, completion.ID)
	require.NoError(t, err)
	assert.Equal(t, newPriceCharged, updatedCompletion.PriceCharged)
	assert.Equal(t, newPaymentMethod, updatedCompletion.PaymentMethod)
	assert.True(t, updatedCompletion.ProviderConfirmed)
	assert.True(t, updatedCompletion.ClientConfirmed)
	assert.Equal(t, user.UserID, *updatedCompletion.UpdatedBy)
	assert.NotNil(t, updatedCompletion.UpdatedAt)
}

func TestServiceCompletionRepositoryIntegration_Delete(t *testing.T) {
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
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceCompletionRepo := &repository.ServiceCompletionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_completion_delete@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "delete-service-completion-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "delete_svc_completion_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "delete_svc_completion_staff@test.com")
	
	// Create test appointment
	startTime := time.Now().Add(24 * time.Hour)
	appointment := createTestAppointmentForCompletion(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, user.UserID, startTime)
	
	// Create test service completion
	completion := createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment.ID, 100.0, user.UserID)

	// Test delete
	err = serviceCompletionRepo.Delete(ctx, completion.ID, user.UserID)
	require.NoError(t, err)

	// Verify the completion is deleted (not found)
	deletedCompletion, err := serviceCompletionRepo.GetByID(ctx, completion.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedCompletion)
	assert.Contains(t, err.Error(), "not found")
}

func TestServiceCompletionRepositoryIntegration_ListByProvider(t *testing.T) {
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
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceCompletionRepo := &repository.ServiceCompletionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_completion_listbyprovider@test.com")
	business1 := createTestBusinessForService(t, businessRepo, user.UserID, "listbyprovider-service-completion-salon1")
	business2 := createTestBusinessForService(t, businessRepo, user.UserID, "listbyprovider-service-completion-salon2")
	client := createTestClientForAppointment(t, clientRepo, business1.BusinessID, user.UserID, "listbyprovider_svc_completion_client@test.com")
	staff1 := createTestStaffForAppointment(t, staffRepo, business1.BusinessID, user.UserID, "listbyprovider_svc_completion_staff1@test.com")
	staff2 := createTestStaffForAppointment(t, staffRepo, business2.BusinessID, user.UserID, "listbyprovider_svc_completion_staff2@test.com")
	
	// Create test appointments for both businesses
	startTime := time.Now().Add(24 * time.Hour)
	appointment1 := createTestAppointmentForCompletion(t, appointmentRepo, business1.BusinessID, client.ID, staff1.StaffID, user.UserID, startTime)
	appointment2 := createTestAppointmentForCompletion(t, appointmentRepo, business1.BusinessID, client.ID, staff1.StaffID, user.UserID, startTime.Add(time.Hour))
	appointment3 := createTestAppointmentForCompletion(t, appointmentRepo, business2.BusinessID, client.ID, staff2.StaffID, user.UserID, startTime.Add(2*time.Hour))
	
	// Create service completions within date range
	startDate := time.Now().AddDate(0, 0, -1)
	endDate := time.Now().AddDate(0, 0, 1)
	
	completion1 := createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment1.ID, 100.0, user.UserID)
	completion2 := createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment2.ID, 150.0, user.UserID)
	completion3 := createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment3.ID, 200.0, user.UserID)
	
	// Test ListByProvider for business1
	completionsBusiness1, err := serviceCompletionRepo.ListByProvider(ctx, business1.BusinessID, startDate, endDate, 1, 10)
	require.NoError(t, err)
	assert.Len(t, completionsBusiness1, 2)
	
	// Test ListByProvider for business2
	completionsBusiness2, err := serviceCompletionRepo.ListByProvider(ctx, business2.BusinessID, startDate, endDate, 1, 10)
	require.NoError(t, err)
	assert.Len(t, completionsBusiness2, 1)
	
	// Verify business1 completions
	expectedIDs := []uuid.UUID{completion1.ID, completion2.ID}
	for _, c := range completionsBusiness1 {
		assert.Contains(t, expectedIDs, c.ID)
	}
	
	// Verify business2 completion
	assert.Equal(t, completion3.ID, completionsBusiness2[0].ID)
}

func TestServiceCompletionRepositoryIntegration_Count(t *testing.T) {
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
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceCompletionRepo := &repository.ServiceCompletionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_completion_count@test.com")
	business1 := createTestBusinessForService(t, businessRepo, user.UserID, "count-service-completion-salon1")
	business2 := createTestBusinessForService(t, businessRepo, user.UserID, "count-service-completion-salon2")
	client := createTestClientForAppointment(t, clientRepo, business1.BusinessID, user.UserID, "count_svc_completion_client@test.com")
	staff1 := createTestStaffForAppointment(t, staffRepo, business1.BusinessID, user.UserID, "count_svc_completion_staff1@test.com")
	staff2 := createTestStaffForAppointment(t, staffRepo, business2.BusinessID, user.UserID, "count_svc_completion_staff2@test.com")
	
	// Create test appointments
	startTime := time.Now().Add(24 * time.Hour)
	appointment1 := createTestAppointmentForCompletion(t, appointmentRepo, business1.BusinessID, client.ID, staff1.StaffID, user.UserID, startTime)
	appointment2 := createTestAppointmentForCompletion(t, appointmentRepo, business1.BusinessID, client.ID, staff1.StaffID, user.UserID, startTime.Add(time.Hour))
	appointment3 := createTestAppointmentForCompletion(t, appointmentRepo, business2.BusinessID, client.ID, staff2.StaffID, user.UserID, startTime.Add(2*time.Hour))
	
	// Create service completions
	createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment1.ID, 100.0, user.UserID)
	createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment2.ID, 150.0, user.UserID)
	createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment3.ID, 200.0, user.UserID)
	
	// Test Count
	count, err := serviceCompletionRepo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count) // 3 total service completions
	
	// Test CountByProvider
	count1, err := serviceCompletionRepo.CountByProvider(ctx, business1.BusinessID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count1)
	
	count2, err := serviceCompletionRepo.CountByProvider(ctx, business2.BusinessID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count2)
}

func TestServiceCompletionRepositoryIntegration_GetProviderRevenue(t *testing.T) {
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
	appointmentRepo := &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	serviceCompletionRepo := &repository.ServiceCompletionRepository{
		BaseRepository: repository.NewBaseRepository(testDB.DB),
	}
	
	ctx := context.Background()
	
	// Create test data
	user := createTestUserForService(t, userRepo, "service_completion_revenue@test.com")
	business := createTestBusinessForService(t, businessRepo, user.UserID, "revenue-service-completion-salon")
	client := createTestClientForAppointment(t, clientRepo, business.BusinessID, user.UserID, "revenue_svc_completion_client@test.com")
	staff := createTestStaffForAppointment(t, staffRepo, business.BusinessID, user.UserID, "revenue_svc_completion_staff@test.com")
	
	// Create test appointments
	startTime := time.Now().Add(24 * time.Hour)
	appointment1 := createTestAppointmentForCompletion(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, user.UserID, startTime)
	appointment2 := createTestAppointmentForCompletion(t, appointmentRepo, business.BusinessID, client.ID, staff.StaffID, user.UserID, startTime.Add(time.Hour))
	
	// Create service completions within date range
	startDate := time.Now().AddDate(0, 0, -1)
	endDate := time.Now().AddDate(0, 0, 1)
	
	createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment1.ID, 100.0, user.UserID)
	createTestServiceCompletionForIntegration(t, serviceCompletionRepo, appointment2.ID, 250.50, user.UserID)
	
	// Test GetProviderRevenue
	revenue, err := serviceCompletionRepo.GetProviderRevenue(ctx, business.BusinessID, startDate, endDate)
	require.NoError(t, err)
	assert.Equal(t, 350.50, revenue) // 100.0 + 250.50 = 350.50
}

// Helper function to create a test appointment for service completion tests
var appointmentCounterForCompletion int = 0
func createTestAppointmentForCompletion(t *testing.T, appointmentRepo *repository.AppointmentRepository, businessID, clientID, staffID, userID uuid.UUID, startTime time.Time) *domain.Appointment {
	appointmentCounterForCompletion++
	
	// First create repositories needed for service creation
	testDB := appointmentRepo.BaseRepository.DB()
	categoryRepo := &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(testDB)}
	serviceRepo := &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(testDB)}
	
	// Create a service category with unique name
	categoryName := fmt.Sprintf("Completion Test Services %d", appointmentCounterForCompletion)
	category := createTestServiceCategory(t, categoryRepo, businessID, categoryName)
	
	// Create a service
	service := createTestServiceForAppointment(t, serviceRepo, businessID, &category.ID, userID, "Test Service for Completion")
	// Add unique offset to avoid conflicts
	adjustedStartTime := startTime.Add(time.Duration(appointmentCounterForCompletion*5) * time.Minute)
	endTime := adjustedStartTime.Add(time.Hour)
	
	appointment := &domain.Appointment{
		BusinessID:     businessID,
		ClientID:       clientID,
		StaffID:        staffID,
		ServiceID:      service.ID,
		StartTime:      adjustedStartTime,
		EndTime:        endTime,
		Status:         "confirmed",
		Notes:          "Test appointment for service completion",
		CreatedBy:      &userID,
	}
	
	err := appointmentRepo.Create(context.Background(), appointment)
	require.NoError(t, err)
	
	return appointment
}

// Helper function to create a test service completion
func createTestServiceCompletionForIntegration(t *testing.T, serviceCompletionRepo *repository.ServiceCompletionRepository, appointmentID uuid.UUID, amount float64, userID uuid.UUID) *domain.ServiceCompletion {
	completionDate := time.Now()
	completion := &domain.ServiceCompletion{
		AppointmentID:     appointmentID,
		PriceCharged:      amount,
		PaymentMethod:     "card",
		ProviderConfirmed: true,
		ClientConfirmed:   true,
		CompletionDate:    &completionDate,
		CreatedBy:         &userID,
	}
	
	err := serviceCompletionRepo.Create(context.Background(), completion)
	require.NoError(t, err)
	
	return completion
}