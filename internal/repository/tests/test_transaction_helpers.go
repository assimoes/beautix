package tests

import (
	"testing"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/models"
	"github.com/assimoes/beautix/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TransactionTestSuite provides a wrapped database connection that automatically
// runs tests within transactions that are rolled back after test completion.
// This ensures tests are truly idempotent regardless of what data was created
// during test execution.
type TransactionTestSuite struct {
	DB       *database.DB
	Tx       *gorm.DB
	t        *testing.T
	teardown func()
}

// NewTransactionTestSuite creates a new transaction test suite that will ensure
// each test runs in its own transaction that will be rolled back after test completion
func NewTransactionTestSuite(t *testing.T) *TransactionTestSuite {
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to connect to test database")

	// Start a transaction
	tx := testDB.Begin()
	
	// Override the default test cleanup to rollback the transaction
	// and then run the original cleanup
	originalCleanup := t.Cleanup
	
	teardown := func() {
		// Always rollback the transaction
		tx.Rollback()
		
		// No need to manually call cleanup since t.Cleanup handles this
	}
	
	t.Cleanup(teardown)
	
	return &TransactionTestSuite{
		DB:       testDB,
		Tx:       tx,
		t:        t,
		teardown: teardown,
	}
}

// CreateTestRepositories returns a set of repository implementations that all use
// the transaction connection, ensuring all operations are contained in a transaction
// that will be rolled back when the test completes
func (ts *TransactionTestSuite) CreateTestRepositories() *TestRepositories {
	return &TestRepositories{
		StaffRepo:                  repository.NewStaffRepository(ts.Tx),
		ServiceAssignmentRepo:      repository.NewServiceAssignmentRepository(ts.Tx),
		AvailabilityExceptionRepo:  repository.NewAvailabilityExceptionRepository(ts.Tx),
		StaffPerformanceRepo:       repository.NewStaffPerformanceRepository(ts.Tx),
	}
}

// TestRepositories contains all repository implementations for testing
type TestRepositories struct {
	StaffRepo                  domain.StaffRepository
	ServiceAssignmentRepo      domain.ServiceAssignmentRepository
	AvailabilityExceptionRepo  domain.AvailabilityExceptionRepository
	StaffPerformanceRepo       domain.StaffPerformanceRepository
}

// CreateTestData creates a complete set of test data (user, business, staff)
// to use for testing, all within the transaction
func (ts *TransactionTestSuite) CreateTestData() *TestData {
	user := createTestUserTx(ts.t, ts.Tx)
	business := createTestBusinessTx(ts.t, ts.Tx, user.ID)
	staff := createTestStaffTx(ts.t, ts.Tx, business.ID, user.ID)
	
	return &TestData{
		User:     user,
		Business: business,
		Staff:    staff,
	}
}

// TestData contains common test data used by tests
type TestData struct {
	User     *models.User
	Business *models.Business
	Staff    *models.Staff
}