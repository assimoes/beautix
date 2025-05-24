package testdb

import (
	"fmt"
	"strings"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/rs/zerolog/log"
)

// CleanDatabase removes all data from the test database
// This ensures each test starts with a clean slate
func CleanDatabase(db *database.DB) error {
	// Get list of all tables in public schema
	var tables []string
	query := `
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename NOT IN ('schema_migrations')
		ORDER BY tablename
	`
	
	if err := db.Raw(query).Pluck("tablename", &tables).Error; err != nil {
		return fmt.Errorf("failed to get table list: %w", err)
	}

	if len(tables) == 0 {
		return nil
	}

	// Disable foreign key constraints temporarily
	if err := db.Exec("SET session_replication_role = 'replica'").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	// Re-enable foreign key constraints after cleaning
	defer func() {
		if err := db.Exec("SET session_replication_role = 'origin'").Error; err != nil {
			log.Error().Err(err).Msg("failed to re-enable foreign key checks")
		}
	}()

	// Truncate all tables
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			log.Warn().Err(err).Str("table", table).Msg("Failed to truncate table")
			// Continue with other tables even if one fails
		}
	}

	// Reset sequences to ensure consistent IDs across tests
	if err := resetSequences(db); err != nil {
		log.Warn().Err(err).Msg("Failed to reset sequences")
	}

	return nil
}

// resetSequences resets all sequences to 1
func resetSequences(db *database.DB) error {
	var sequences []string
	query := `
		SELECT sequence_name 
		FROM information_schema.sequences 
		WHERE sequence_schema = 'public'
	`

	if err := db.Raw(query).Pluck("sequence_name", &sequences).Error; err != nil {
		return fmt.Errorf("failed to get sequences: %w", err)
	}

	for _, seq := range sequences {
		if err := db.Exec(fmt.Sprintf("ALTER SEQUENCE %s RESTART WITH 1", seq)).Error; err != nil {
			log.Warn().Err(err).Str("sequence", seq).Msg("Failed to reset sequence")
		}
	}

	return nil
}

// CleanSpecificTables cleans only specified tables
func CleanSpecificTables(db *database.DB, tables ...string) error {
	if len(tables) == 0 {
		return nil
	}

	// Disable foreign key constraints
	if err := db.Exec("SET session_replication_role = 'replica'").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	defer func() {
		if err := db.Exec("SET session_replication_role = 'origin'").Error; err != nil {
			log.Error().Err(err).Msg("failed to re-enable foreign key checks")
		}
	}()

	// Truncate specified tables
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}

// DeleteByID deletes a record by ID from a table
func DeleteByID(db *database.DB, table string, id interface{}) error {
	return db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", table), id).Error
}

// DeleteByCondition deletes records matching a condition
func DeleteByCondition(db *database.DB, table string, condition string, args ...interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, condition)
	return db.Exec(query, args...).Error
}

// CountRecords counts records in a table
func CountRecords(db *database.DB, table string) (int64, error) {
	var count int64
	err := db.Table(table).Count(&count).Error
	return count, err
}

// TableExists checks if a table exists
func TableExists(db *database.DB, tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM pg_tables 
			WHERE schemaname = 'public' 
			AND tablename = $1
		)
	`
	err := db.Raw(query, tableName).Scan(&exists).Error
	return exists, err
}

// GetTableList returns list of all tables in public schema
func GetTableList(db *database.DB) ([]string, error) {
	var tables []string
	query := `
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename NOT IN ('schema_migrations')
		ORDER BY tablename
	`
	err := db.Raw(query).Pluck("tablename", &tables).Error
	return tables, err
}

// VerifyCleanState verifies that all tables are empty
func VerifyCleanState(db *database.DB) error {
	tables, err := GetTableList(db)
	if err != nil {
		return fmt.Errorf("failed to get table list: %w", err)
	}

	var nonEmptyTables []string
	for _, table := range tables {
		count, err := CountRecords(db, table)
		if err != nil {
			return fmt.Errorf("failed to count records in %s: %w", table, err)
		}
		if count > 0 {
			nonEmptyTables = append(nonEmptyTables, fmt.Sprintf("%s(%d)", table, count))
		}
	}

	if len(nonEmptyTables) > 0 {
		return fmt.Errorf("tables not empty: %s", strings.Join(nonEmptyTables, ", "))
	}

	return nil
}