package testdb

import (
	"fmt"

	"github.com/assimoes/beautix/internal/infrastructure/database"
)

// SetupTestSchema ensures the test database has necessary extensions
// Note: All DDL (triggers, functions) should be managed through migrations
func SetupTestSchema(db *database.DB) error {
	// Create necessary extensions only
	// These are typically idempotent and safe to run
	extensions := []string{
		"uuid-ossp",
		"pgcrypto",
	}

	for _, ext := range extensions {
		if err := db.Exec(fmt.Sprintf(`CREATE EXTENSION IF NOT EXISTS "%s"`, ext)).Error; err != nil {
			return fmt.Errorf("failed to create extension %s: %w", ext, err)
		}
	}

	// All other schema setup (triggers, functions, tables) should be handled by migrations
	// This ensures consistency between test and production environments

	return nil
}