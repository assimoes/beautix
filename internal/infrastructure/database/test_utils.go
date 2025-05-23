package database

import (
	"fmt"
	"testing"

	"github.com/assimoes/beautix/configs"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestDB represents a test database connection
type TestDB struct {
	*DB
	Config *configs.Config
}

// NewTestDB creates a connection to the test database
func NewTestDB(t *testing.T) (*TestDB, error) {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Update config to use the test database
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

	// Connect to the test database
	db, err := NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Reset the test database - truncate all tables to ensure a clean slate
	if err := cleanTestDatabase(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to clean test database: %w", err)
	}

	// Run basic schema setup
	if err := setupTestSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to setup test schema: %w", err)
	}

	// Register cleanup function
	testDB := &TestDB{
		DB:     db,
		Config: config,
	}

	// Clean up after the test
	t.Cleanup(func() {
		if err := cleanTestDatabase(db); err != nil {
			log.Error().Err(err).Msg("Failed to clean test database after test")
		}
		if err := testDB.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close test database connection")
		}
	})

	return testDB, nil
}

// cleanTestDatabase cleans all tables in the test database
func cleanTestDatabase(db *DB) error {
	// Get a list of all tables in the public schema
	var tables []string
	if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Pluck("tablename", &tables).Error; err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Skip if no tables exist
	if len(tables) == 0 {
		return nil
	}

	// Disable foreign key checks
	if err := db.Exec("SET session_replication_role = 'replica'").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	defer func() {
		// Re-enable foreign key checks in a deferred function to ensure it runs
		if err := db.Exec("SET session_replication_role = 'origin'").Error; err != nil {
			log.Error().Err(err).Msg("failed to re-enable foreign key checks")
		}
	}()

	// First try to truncate all tables
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE \"%s\" CASCADE", table)).Error; err != nil {
			log.Warn().Err(err).Str("table", table).Msg("Failed to truncate table, will try dropping it")
		}
	}

	// Get the list again to see if any test tables remain (like test_transactions)
	// that might have been created during tests
	if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Pluck("tablename", &tables).Error; err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Drop specific test tables that might cause issues
	for _, table := range tables {
		// Only drop tables that are obviously test tables
		if table == "test_transactions" || table == "test_table" {
			if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS \"%s\" CASCADE", table)).Error; err != nil {
				log.Warn().Err(err).Str("table", table).Msg("Failed to drop test table")
			}
		}
	}

	return nil
}

// setupTestSchema sets up the basic schema needed for tests
func setupTestSchema(db *DB) error {
	// Create extensions
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		return fmt.Errorf("failed to create pgcrypto extension: %w", err)
	}

	// Apply critical trigger fixes that are needed for tests to work
	// This ensures test database has the correct trigger functions even if migrations aren't run
	if err := applyTriggerFixes(db); err != nil {
		return fmt.Errorf("failed to apply trigger fixes: %w", err)
	}

	return nil
}

// applyTriggerFixes applies the trigger function fixes needed for tests
func applyTriggerFixes(db *DB) error {
	// Fix the appointment overlap check function to use 'id' instead of 'appointment_id'
	appointmentTriggerFix := `
		CREATE OR REPLACE FUNCTION check_appointment_overlap()
		RETURNS TRIGGER AS $$
		DECLARE
		    conflict_count INTEGER;
		BEGIN
		    -- Check for overlapping appointments for the same staff
		    SELECT COUNT(*) INTO conflict_count
		    FROM public.appointments a
		    WHERE a.staff_id = NEW.staff_id
		    AND a.id != NEW.id  -- Fixed: use 'id' instead of 'appointment_id'
		    AND a.deleted_at IS NULL
		    AND a.status NOT IN ('cancelled', 'no-show')
		    AND (
		        (NEW.start_time, NEW.end_time) OVERLAPS (a.start_time, a.end_time)
		    );
		    
		    IF conflict_count > 0 THEN
		        RAISE EXCEPTION 'Appointment time conflicts with existing appointment for this staff member';
		    END IF;
		    
		    RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`

	if err := db.Exec(appointmentTriggerFix).Error; err != nil {
		return fmt.Errorf("failed to fix appointment trigger function: %w", err)
	}

	// Fix the resource booking overlap check function to use 'id' instead of 'booking_id'
	resourceTriggerFix := `
		CREATE OR REPLACE FUNCTION check_resource_booking_overlap()
		RETURNS TRIGGER AS $$
		DECLARE
		    conflict_count INTEGER;
		BEGIN
		    -- Check for overlapping resource bookings
		    SELECT COUNT(*) INTO conflict_count
		    FROM public.resource_bookings rb
		    WHERE rb.resource_id = NEW.resource_id
		    AND rb.id != NEW.id  -- Fixed: use 'id' instead of 'booking_id'
		    AND rb.deleted_at IS NULL
		    AND rb.status NOT IN ('cancelled')
		    AND (
		        (NEW.start_time, NEW.end_time) OVERLAPS (rb.start_time, rb.end_time)
		    );
		    
		    IF conflict_count > 0 THEN
		        RAISE EXCEPTION 'Resource booking time conflicts with existing booking for this resource';
		    END IF;
		    
		    RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`

	if err := db.Exec(resourceTriggerFix).Error; err != nil {
		return fmt.Errorf("failed to fix resource booking trigger function: %w", err)
	}

	return nil
}

// WithTestTransaction runs a test within a transaction that is rolled back
func WithTestTransaction(t *testing.T, db *DB, fn func(tx *gorm.DB)) {
	tx := db.Begin()
	defer tx.Rollback()

	fn(tx)

	if t.Failed() {
		t.Log("Test failed, transaction rolled back")
	}
}

// CreateTestTable creates a test table with the given name and schema
// and ensures it's dropped after the test completes
func CreateTestTable(t *testing.T, db *DB, tableName, schema string) {
	// First drop the table if it exists
	err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)).Error
	if err != nil {
		t.Fatalf("Failed to drop test table %s: %v", tableName, err)
	}

	// Create the table
	err = db.Exec(schema).Error
	if err != nil {
		t.Fatalf("Failed to create test table %s: %v", tableName, err)
	}

	// Make sure the table is dropped after the test
	t.Cleanup(func() {
		err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)).Error
		if err != nil {
			log.Error().Err(err).Str("table", tableName).Msg("Failed to drop test table during cleanup")
		}
	})
}

// TableExists checks if a table exists in the database
func TableExists(db *DB, tableName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename = $1
	)`

	err := db.Raw(query, tableName).Scan(&exists).Error
	return exists, err
}

// NewTestDBWithTransaction creates a test database connection and runs the test
// within a transaction that is rolled back after completion. This ensures
// true test isolation and prevents interference between tests.
func NewTestDBWithTransaction(t *testing.T, testFn func(db *gorm.DB)) {
	// Load configuration
	config, err := configs.LoadConfig()
	require.NoError(t, err, "Failed to load config")

	// Update config to use the test database
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

	// Connect to the test database WITHOUT global cleanup
	db, err := NewConnection(config)
	require.NoError(t, err, "Failed to connect to test database")

	// Ensure database connection is closed after test
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close test database connection")
		}
	})

	// Run basic schema setup (only extensions, no cleanup)
	err = setupTestSchema(db)
	require.NoError(t, err, "Failed to setup test schema")

	// Start a transaction
	tx := db.Begin()

	// Ensure transaction is always rolled back
	t.Cleanup(func() {
		tx.Rollback()
	})

	// Run the test function with the transaction
	testFn(tx)
}
