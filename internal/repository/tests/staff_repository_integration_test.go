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

// createTestUserForStaff creates a test user for staff tests
func createTestUserForStaff(t *testing.T, userRepo *repository.UserRepository, email string) *domain.User {
	user := &domain.User{
		ClerkID:   uuid.New(),
		Email:     email,
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+351123456789",
		Role:      "staff",
		IsActive:  true,
	}

	err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)
	return user
}

// createTestBusinessForStaff creates a test business for staff tests
func createTestBusinessForStaff(t *testing.T, businessRepo *repository.BusinessRepository, ownerID uuid.UUID, name string) *domain.Business {
	business := &domain.Business{
		OwnerID:          ownerID,
		BusinessName:     name,
		BusinessType:     "beauty_salon",
		AddressLine1:     "123 Test Street",
		City:             "Lisbon",
		Country:          "Portugal",
		Phone:            "+351987654321",
		Email:            name + "@salon.com",
		TimeZone:         "Europe/Lisbon",
		SubscriptionPlan: "basic",
		IsActive:         true,
	}

	err := businessRepo.Create(context.Background(), business)
	require.NoError(t, err)
	return business
}

func TestStaffRepositoryIntegration_Create(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForStaff(t, userRepo, "staff_create@test.com")
	business := createTestBusinessForStaff(t, businessRepo, user.UserID, "create-staff-salon")

	// Create a new staff member
	staff := &domain.Staff{
		BusinessID:      business.BusinessID,
		UserID:          user.UserID,
		Position:        "Test Position",
		Bio:             "Test Bio",
		SpecialtyAreas:  []string{"Area 1", "Area 2"},
		ProfileImageURL: "http://example.com/profile.jpg",
		IsActive:        true,
		EmploymentType:  "full-time",
		JoinDate:        time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
		CommissionRate:  20.0,
		CreatedBy:       user.UserID,
	}

	// Test creation
	err = staffRepo.Create(ctx, staff)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, staff.StaffID, "Staff ID should be generated")
	assert.NotZero(t, staff.CreatedAt, "Created at timestamp should be set")

	// Verify the staff was created
	result, err := staffRepo.GetByID(ctx, staff.StaffID)
	assert.NoError(t, err)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, user.UserID, result.UserID)
	assert.Equal(t, "Test Position", result.Position)
	assert.Equal(t, "Test Bio", result.Bio)
	assert.ElementsMatch(t, []string{"Area 1", "Area 2"}, result.SpecialtyAreas)
	assert.Equal(t, "http://example.com/profile.jpg", result.ProfileImageURL)
	assert.True(t, result.IsActive)
	assert.Equal(t, "full-time", result.EmploymentType)
	assert.Equal(t, float64(20.0), result.CommissionRate)
	assert.Equal(t, user.UserID, result.CreatedBy)
}

func TestStaffRepositoryIntegration_GetByID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForStaff(t, userRepo, "staff_getbyid@test.com")
	business := createTestBusinessForStaff(t, businessRepo, user.UserID, "getbyid-staff-salon")

	// Create staff member using domain object
	staff := &domain.Staff{
		BusinessID:      business.BusinessID,
		UserID:          user.UserID,
		Position:        "Senior Stylist",
		Bio:             "Expert in hair coloring",
		SpecialtyAreas:  []string{"Hair", "Color"},
		IsActive:        true,
		EmploymentType:  "full-time",
		JoinDate:        time.Now().Add(-60 * 24 * time.Hour),
		CommissionRate:  25.0,
		CreatedBy:       user.UserID,
	}

	err = staffRepo.Create(ctx, staff)
	require.NoError(t, err)

	// Test GetByID
	result, err := staffRepo.GetByID(ctx, staff.StaffID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, staff.StaffID, result.StaffID)
	assert.Equal(t, business.BusinessID, result.BusinessID)
	assert.Equal(t, user.UserID, result.UserID)
	assert.Equal(t, "Senior Stylist", result.Position)

	// Verify related entities are populated
	assert.NotNil(t, result.User)
	assert.Equal(t, user.UserID, result.User.UserID)
	assert.NotEmpty(t, result.User.Email)

	assert.NotNil(t, result.Business)
	assert.Equal(t, business.BusinessID, result.Business.BusinessID)
}

func TestStaffRepositoryIntegration_GetByID_NotFound(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repository
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Test GetByID with non-existent ID
	result, err := staffRepo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffRepositoryIntegration_GetByUserID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user
	user := createTestUserForStaff(t, userRepo, "staff_getbyuserid@test.com")

	// Create two businesses
	business1 := createTestBusinessForStaff(t, businessRepo, user.UserID, "business1-staff")
	business2 := createTestBusinessForStaff(t, businessRepo, user.UserID, "business2-staff")

	// Create staff members for both businesses
	staff1 := &domain.Staff{
		BusinessID:     business1.BusinessID,
		UserID:         user.UserID,
		Position:       "Stylist",
		IsActive:       true,
		EmploymentType: "full-time",
		JoinDate:       time.Now().Add(-30 * 24 * time.Hour),
		CreatedBy:      user.UserID,
	}
	err = staffRepo.Create(ctx, staff1)
	require.NoError(t, err)

	staff2 := &domain.Staff{
		BusinessID:     business2.BusinessID,
		UserID:         user.UserID,
		Position:       "Manager",
		IsActive:       true,
		EmploymentType: "full-time",
		JoinDate:       time.Now().Add(-15 * 24 * time.Hour),
		CreatedBy:      user.UserID,
	}
	err = staffRepo.Create(ctx, staff2)
	require.NoError(t, err)

	// Test GetByUserID
	results, err := staffRepo.GetByUserID(ctx, user.UserID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Check that both staff members are returned
	staffIDs := []uuid.UUID{results[0].StaffID, results[1].StaffID}
	assert.Contains(t, staffIDs, staff1.StaffID)
	assert.Contains(t, staffIDs, staff2.StaffID)
}

func TestStaffRepositoryIntegration_GetByBusinessID(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create two users
	user1 := createTestUserForStaff(t, userRepo, "staff1_getbybusiness@test.com")
	user2 := createTestUserForStaff(t, userRepo, "staff2_getbybusiness@test.com")

	// Create business
	business := createTestBusinessForStaff(t, businessRepo, user1.UserID, "getbybusiness-salon")

	// Create staff members for the business
	staff1 := &domain.Staff{
		BusinessID:     business.BusinessID,
		UserID:         user1.UserID,
		Position:       "Senior Stylist",
		IsActive:       true,
		EmploymentType: "full-time",
		JoinDate:       time.Now().Add(-60 * 24 * time.Hour),
		CreatedBy:      user1.UserID,
	}
	err = staffRepo.Create(ctx, staff1)
	require.NoError(t, err)

	staff2 := &domain.Staff{
		BusinessID:     business.BusinessID,
		UserID:         user2.UserID,
		Position:       "Junior Stylist",
		IsActive:       true,
		EmploymentType: "part-time",
		JoinDate:       time.Now().Add(-10 * 24 * time.Hour),
		CreatedBy:      user1.UserID,
	}
	err = staffRepo.Create(ctx, staff2)
	require.NoError(t, err)

	// Create inactive staff member - note: we need to use a different user to avoid unique constraint
	user3 := createTestUserForStaff(t, userRepo, "staff3_getbybusiness@test.com")
	staff3 := &domain.Staff{
		BusinessID:     business.BusinessID,
		UserID:         user3.UserID,
		Position:       "Intern",
		IsActive:       false,
		EmploymentType: "intern",
		JoinDate:       time.Now().Add(-5 * 24 * time.Hour),
		CreatedBy:      user1.UserID,
	}
	err = staffRepo.Create(ctx, staff3)
	require.NoError(t, err)

	// Test GetByBusinessID (returns all staff members for a business)
	results, err := staffRepo.GetByBusinessID(ctx, business.BusinessID)
	assert.NoError(t, err)
	assert.Len(t, results, 3, "Should return all staff members")

	// Verify all staff members are returned with correct business ID
	for _, s := range results {
		assert.Equal(t, business.BusinessID, s.BusinessID)
	}

	// Check that all staff members are returned
	staffIDs := []uuid.UUID{results[0].StaffID, results[1].StaffID, results[2].StaffID}
	assert.Contains(t, staffIDs, staff1.StaffID)
	assert.Contains(t, staffIDs, staff2.StaffID)
	assert.Contains(t, staffIDs, staff3.StaffID)
}

func TestStaffRepositoryIntegration_Update(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForStaff(t, userRepo, "staff_update@test.com")
	business := createTestBusinessForStaff(t, businessRepo, user.UserID, "update-staff-salon")

	// Create staff member
	staff := &domain.Staff{
		BusinessID:      business.BusinessID,
		UserID:          user.UserID,
		Position:        "Junior Stylist",
		Bio:             "Original bio",
		SpecialtyAreas:  []string{"Hair"},
		IsActive:        true,
		EmploymentType:  "part-time",
		JoinDate:        time.Now().Add(-30 * 24 * time.Hour),
		CommissionRate:  15.0,
		CreatedBy:       user.UserID,
	}
	err = staffRepo.Create(ctx, staff)
	require.NoError(t, err)

	// Update the staff member
	position := "Senior Stylist"
	bio := "Updated bio with more experience"
	specialtyAreas := []string{"Hair", "Color", "Extensions"}
	employmentType := "full-time"
	commissionRate := 25.0
	isActive := true

	updateInput := &domain.UpdateStaffInput{
		Position:       &position,
		Bio:            &bio,
		SpecialtyAreas: &specialtyAreas,
		EmploymentType: &employmentType,
		CommissionRate: &commissionRate,
		IsActive:       &isActive,
	}

	err = staffRepo.Update(ctx, staff.StaffID, updateInput, user.UserID)
	assert.NoError(t, err)

	// Verify the update
	updated, err := staffRepo.GetByID(ctx, staff.StaffID)
	assert.NoError(t, err)
	assert.Equal(t, "Senior Stylist", updated.Position)
	assert.Equal(t, "Updated bio with more experience", updated.Bio)
	assert.ElementsMatch(t, []string{"Hair", "Color", "Extensions"}, updated.SpecialtyAreas)
	assert.Equal(t, "full-time", updated.EmploymentType)
	assert.Equal(t, float64(25.0), updated.CommissionRate)
	assert.True(t, updated.IsActive)
}

func TestStaffRepositoryIntegration_Delete(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test user and business
	user := createTestUserForStaff(t, userRepo, "staff_delete@test.com")
	business := createTestBusinessForStaff(t, businessRepo, user.UserID, "delete-staff-salon")

	// Create staff member
	staff := &domain.Staff{
		BusinessID:     business.BusinessID,
		UserID:         user.UserID,
		Position:       "Stylist",
		IsActive:       true,
		EmploymentType: "full-time",
		JoinDate:       time.Now().Add(-30 * 24 * time.Hour),
		CreatedBy:      user.UserID,
	}
	err = staffRepo.Create(ctx, staff)
	require.NoError(t, err)

	// Delete the staff member
	err = staffRepo.Delete(ctx, staff.StaffID, user.UserID)
	assert.NoError(t, err)

	// Verify deletion
	result, err := staffRepo.GetByID(ctx, staff.StaffID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaffRepositoryIntegration_ListWithPagination(t *testing.T) {
	// Create test database connection
	testDB, err := database.NewSimpleTestDB(t)
	require.NoError(t, err)
	
	// Clean all tables before test
	database.CleanupAllTables(t, testDB.DB)
	
	// Create repositories
	userRepo := &repository.UserRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	businessRepo := &repository.BusinessRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	staffRepo := &repository.StaffRepository{BaseRepository: repository.NewBaseRepository(testDB.DB)}
	ctx := context.Background()

	// Create test users
	users := make([]*domain.User, 5)
	for i := 0; i < 5; i++ {
		users[i] = createTestUserForStaff(t, userRepo, uuid.New().String()+"@test.com")
	}

	// Create business
	business := createTestBusinessForStaff(t, businessRepo, users[0].UserID, "pagination-salon")

	// Create multiple staff members
	for i := 0; i < 5; i++ {
		staff := &domain.Staff{
			BusinessID:     business.BusinessID,
			UserID:         users[i].UserID,
			Position:       "Stylist",
			IsActive:       true,
			EmploymentType: "full-time",
			JoinDate:       time.Now().Add(-time.Duration(i*10) * 24 * time.Hour),
			CreatedBy:      users[0].UserID,
		}
		err = staffRepo.Create(ctx, staff)
		require.NoError(t, err)
	}

	// Test pagination
	page1, err := staffRepo.List(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, err := staffRepo.List(ctx, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, page2, 2)

	page3, err := staffRepo.List(ctx, 3, 2)
	assert.NoError(t, err)
	assert.Len(t, page3, 1)

	// Test count
	count, err := staffRepo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}