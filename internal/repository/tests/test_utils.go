package tests

import (
	"testing"
	"time"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/assimoes/beautix/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// SimpleTestSuite provides a simple test setup that cleans tables before each test
type SimpleTestSuite struct {
	t      *testing.T
	testDB *database.SimpleTestDB
}

// NewSimpleTestSuite creates a new simple test suite
func NewSimpleTestSuite(t *testing.T) *SimpleTestSuite {
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")
	
	return &SimpleTestSuite{
		t:      t,
		testDB: testDB,
	}
}

// CleanTables cleans the specified tables in the correct order
func (s *SimpleTestSuite) CleanTables(tables ...string) {
	// If no specific tables provided, clean all tables
	if len(tables) == 0 {
		database.CleanupAllTables(s.t, s.testDB.DB)
		return
	}
	
	// Clean specific tables
	database.CleanupBeforeTest(s.t, s.testDB.DB, tables...)
}

// CreateTestRepositories creates all repository instances for testing
func (s *SimpleTestSuite) CreateTestRepositories() *SimpleTestRepositories {
	return &SimpleTestRepositories{
		UserRepo:                  &repository.UserRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		BusinessRepo:              &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		StaffRepo:                 &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		ServiceAssignmentRepo:     &repository.ServiceAssignmentRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		AvailabilityExceptionRepo: &repository.AvailabilityExceptionRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		StaffPerformanceRepo:      &repository.StaffPerformanceRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		ClientRepo:                &repository.ClientRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		ServiceRepo:               &repository.ServiceRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		ServiceCategoryRepo:       &repository.ServiceCategoryRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		AppointmentRepo:           &repository.AppointmentRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
		ServiceCompletionRepo:     &repository.ServiceCompletionRepository{BaseRepository: repository.NewBaseRepository(s.testDB.DB)},
	}
}

// GetDB returns the test database
func (s *SimpleTestSuite) GetDB() *database.SimpleTestDB {
	return s.testDB
}

// SimpleTestRepositories holds all repository instances for testing
type SimpleTestRepositories struct {
	UserRepo                  *repository.UserRepository
	BusinessRepo              *repository.BusinessRepository
	StaffRepo                 *repository.StaffRepository
	ServiceAssignmentRepo     *repository.ServiceAssignmentRepository
	AvailabilityExceptionRepo *repository.AvailabilityExceptionRepository
	StaffPerformanceRepo      *repository.StaffPerformanceRepository
	ClientRepo                *repository.ClientRepository
	ServiceRepo               *repository.ServiceRepository
	ServiceCategoryRepo       *repository.ServiceCategoryRepository
	AppointmentRepo           *repository.AppointmentRepository
	ServiceCompletionRepo     *repository.ServiceCompletionRepository
}

// Helper functions to create test data for repository integration tests
// Each function is available in two versions:
// 1. The original version using a direct DB connection (for backward compatibility)
// 2. A TX version that accepts a transaction for true test isolation

// createTestUser creates a test user in the database
func createTestUser(t *testing.T, db *database.DB) *models.User {
	user := &models.User{
		ClerkID:   uuid.New().String(),                  // Use UUID as a unique clerk_id
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
		ClerkID:   uuid.New().String(),                  // Use UUID as a unique clerk_id
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
		UserID:      userID,
		Name:        businessName,
		DisplayName: "Test Business",

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
		UserID:      userID,
		Name:        businessName,
		DisplayName: "Test Business",

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

// createTestClient creates a test client in the database
func createTestClient(t *testing.T, db *database.DB, businessID uuid.UUID, userID *uuid.UUID, createdByID *uuid.UUID) *models.Client {
	clientEmail := uuid.New().String() + "@example.com" // Ensure unique email
	client := &models.Client{
		BusinessID:       businessID,
		UserID:           userID,
		FirstName:        "Test",
		LastName:         "Client",
		Email:            clientEmail,
		Phone:            "+1987654321",
		Notes:            "Test client notes",
		IsActive:         true,
		AcceptsMarketing: true,
	}

	// Set created_by if provided
	if createdByID != nil {
		client.CreatedBy = createdByID
	}

	err := db.Create(client).Error
	require.NoError(t, err, "Failed to create test client")

	return client
}

// createTestClientTx creates a test client in the database within a transaction
func createTestClientTx(t *testing.T, tx *gorm.DB, businessID uuid.UUID, userID *uuid.UUID, createdByID *uuid.UUID) *models.Client {
	clientEmail := uuid.New().String() + "@example.com" // Ensure unique email
	client := &models.Client{
		BusinessID:       businessID,
		UserID:           userID,
		FirstName:        "Test",
		LastName:         "Client",
		Email:            clientEmail,
		Phone:            "+1987654321",
		Notes:            "Test client notes",
		IsActive:         true,
		AcceptsMarketing: true,
	}

	// Set created_by if provided
	if createdByID != nil {
		client.CreatedBy = createdByID
	}

	err := tx.Create(client).Error
	require.NoError(t, err, "Failed to create test client")

	return client
}

// createTestServiceCategoryTx creates a test service category in the database within a transaction
func createTestServiceCategoryTx(t *testing.T, tx *gorm.DB, businessID uuid.UUID, createdByID *uuid.UUID) *models.ServiceCategory {
	category := &models.ServiceCategory{
		BusinessID:  businessID,
		Name:        "Test Category " + uuid.New().String()[0:8],
		Description: "Test category description",
	}

	// Note: ServiceCategory doesn't have audit fields

	err := tx.Create(category).Error
	require.NoError(t, err, "Failed to create test service category")

	return category
}

// createTestServiceTx creates a test service in the database within a transaction
func createTestServiceTx(t *testing.T, tx *gorm.DB, businessID uuid.UUID, categoryID *uuid.UUID, createdByID *uuid.UUID) *models.Service {
	service := &models.Service{
		BusinessID:  businessID,
		Name:        "Test Service " + uuid.New().String()[0:8],
		Description: "Test service description",
		Duration:    60, // 60 minutes
		Price:       100.0,
		IsActive:    true,
	}

	// Set created_by if provided
	if createdByID != nil {
		service.CreatedBy = createdByID
	}

	err := tx.Create(service).Error
	require.NoError(t, err, "Failed to create test service")

	return service
}

// createTestAppointmentTx creates a test appointment in the database within a transaction
func createTestAppointmentTx(t *testing.T, tx *gorm.DB, businessID, clientID uuid.UUID, startTime time.Time, createdByID *uuid.UUID) *models.Appointment {
	// Calculate end time (1 hour after start time)
	endTime := startTime.Add(time.Hour)

	appointment := &models.Appointment{
		BusinessID:      businessID,
		ClientID:        clientID,
		StaffID:         uuid.New(), // Generate a staff ID for now
		ServiceID:       uuid.New(), // Generate a service ID for now
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          models.AppointmentStatusConfirmed,
		Notes:           "Test appointment",
		EstimatedPrice:  func() *float64 { p := 100.0; return &p }(),
		PaymentStatus:   models.PaymentStatusPending,
		ClientConfirmed: true,
		StaffConfirmed:  true,
	}

	// Set created_by if provided
	if createdByID != nil {
		appointment.CreatedBy = createdByID
	}

	err := tx.Create(appointment).Error
	require.NoError(t, err, "Failed to create test appointment")

	return appointment
}