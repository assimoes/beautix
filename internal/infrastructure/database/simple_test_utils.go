package database

import (
	"fmt"
	"testing"

	"github.com/assimoes/beautix/configs"
	"github.com/rs/zerolog/log"
)

// SimpleTestDB represents a simple test database connection without complex cleanup
type SimpleTestDB struct {
	*DB
	Config *configs.Config
}

// NewSimpleTestDB creates a simple test database connection that only does minimal setup
// NOTE: This assumes migrations have already been run on the test database using:
//   make migrate-up TEST_DATABASE_URL=postgres://user:pass@host:port/beautix_test?sslmode=disable
func NewSimpleTestDB(t *testing.T) (*SimpleTestDB, error) {
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

	// Just setup extensions - no table cleanup
	if err := setupBasicTestSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to setup test schema: %w", err)
	}

	testDB := &SimpleTestDB{
		DB:     db,
		Config: config,
	}

	// Only close connection on cleanup, no table operations
	t.Cleanup(func() {
		if err := testDB.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close test database connection")
		}
	})

	return testDB, nil
}

// setupBasicTestSchema sets up only extensions and critical functions
func setupBasicTestSchema(db *DB) error {
	// Create extensions
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		return fmt.Errorf("failed to create pgcrypto extension: %w", err)
	}

	// Apply trigger fixes without error checking (they may already exist)
	applyTriggerFixesQuiet(db)

	return nil
}

// applyTriggerFixesQuiet applies trigger fixes without failing on errors
func applyTriggerFixesQuiet(db *DB) {
	// Fix appointment overlap function
	appointmentTriggerFix := `
		CREATE OR REPLACE FUNCTION check_appointment_overlap()
		RETURNS TRIGGER AS $$
		DECLARE
		    conflict_count INTEGER;
		BEGIN
		    SELECT COUNT(*) INTO conflict_count
		    FROM public.appointments a
		    WHERE a.staff_id = NEW.staff_id
		    AND a.id != NEW.id
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
	db.Exec(appointmentTriggerFix)

	// Fix resource booking overlap function
	resourceTriggerFix := `
		CREATE OR REPLACE FUNCTION check_resource_booking_overlap()
		RETURNS TRIGGER AS $$
		DECLARE
		    conflict_count INTEGER;
		BEGIN
		    SELECT COUNT(*) INTO conflict_count
		    FROM public.resource_bookings rb
		    WHERE rb.resource_id = NEW.resource_id
		    AND rb.id != NEW.id
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
	db.Exec(resourceTriggerFix)
}

// CleanupBeforeTest provides a comprehensive cleanup that handles all foreign key dependencies
func CleanupBeforeTest(t *testing.T, db *DB, tables ...string) {
	// If specific tables are requested, clean just those in order
	if len(tables) > 0 {
		for _, table := range tables {
			if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
				t.Logf("Warning: failed to clean table %s: %v", table, err)
			}
		}
		return
	}

	// If no specific tables, do comprehensive cleanup in proper dependency order
	CleanupAllTables(t, db)
}

// CleanupAllTables performs a comprehensive cleanup of all tables in the correct order
func CleanupAllTables(t *testing.T, db *DB) {
	// Complete dependency order - children first, parents last
	allTables := []string{
		// Level 7: Most dependent children
		"reward_redemptions",
		"loyalty_transactions",
		"client_loyalty_memberships",
		"service_completions",
		"campaign_clients",
		"appointment_services",
		"waiting_list",
		"user_connected_accounts",
		"business_locations",
		"service_ratings",
		"staff_performance",
		"resource_bookings",
		
		// Level 6: Highly dependent
		"service_assignment",
		"availability_exception",
		
		// Level 5: Mid-level children
		"appointments",
		"loyalty_rewards",
		
		// Level 4: Basic children
		"services",
		"clients",
		"loyalty_programs",
		"campaigns",
		"resources",
		
		// Level 3: Categories and staff
		"service_categories",
		"staff",
		
		// Level 2: Business level
		"businesses",
		"providers",
		
		// Level 1: Root tables
		"users",
		
		// Level 0: System tables (usually not needed)
		"schema_migrations",
	}

	for _, table := range allTables {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			// Only log warning, don't fail test
			t.Logf("Warning: failed to clean table %s: %v", table, err)
		}
	}
}