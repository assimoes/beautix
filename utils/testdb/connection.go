package testdb

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/assimoes/beautix/configs"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global mutex to ensure sequential test execution
var testMutex sync.Mutex

// TestDB represents a test database instance
type TestDB struct {
	DB     *database.DB
	Config *configs.Config
	t      *testing.T
}

// NewTestDB creates a new test database connection
// This ensures sequential execution and proper cleanup
func NewTestDB(t *testing.T) *TestDB {
	// Acquire lock to ensure sequential execution
	testMutex.Lock()
	
	// Load configuration
	config, err := configs.LoadConfig()
	require.NoError(t, err, "Failed to load config")

	// Update config to use test database
	config.Database.DBName = "beautix_test"
	config.Database.URL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	// Connect to test database with silent logging
	db, err := database.NewConnectionWithConfig(config, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "Failed to connect to test database")

	testDB := &TestDB{
		DB:     db,
		Config: config,
		t:      t,
	}

	// Setup cleanup that will run after test completes
	t.Cleanup(func() {
		// Clean all data before closing connection
		if err := CleanDatabase(testDB.DB); err != nil {
			log.Error().Err(err).Msg("Failed to clean database after test")
		}

		// Close database connection
		if err := testDB.DB.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close test database connection")
		}

		// Release lock after cleanup is complete
		testMutex.Unlock()
	})

	// Initial cleanup to ensure pristine state
	err = CleanDatabase(testDB.DB)
	require.NoError(t, err, "Failed to clean database before test")

	// Ensure schema is ready
	err = SetupTestSchema(testDB.DB)
	require.NoError(t, err, "Failed to setup test schema")

	return testDB
}

// GetDB returns the underlying GORM database instance
func (tdb *TestDB) GetDB() *database.DB {
	return tdb.DB
}

// WaitForCondition waits for a condition to be true with timeout
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("Timeout waiting for condition: %s", message)
}

// AssertEventuallyConsistent checks that a condition becomes true within a timeout
func AssertEventuallyConsistent(t *testing.T, condition func() bool, timeout time.Duration) {
	WaitForCondition(t, condition, timeout, "Condition did not become true within timeout")
}