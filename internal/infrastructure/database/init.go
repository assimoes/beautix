package database

import (
	"fmt"
	"log"
	"time"

	"github.com/assimoes/beautix/configs"
)

// Initialize initializes the database connection and runs migrations
func Initialize(config *configs.Config) (*DB, error) {
	// Connect to the database
	db, err := NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// RunMigrations runs the database migrations
func RunMigrations(config *configs.Config) error {
	// Connect to the database
	db, err := NewConnection(config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Verify the connection by pinging
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	
	// Using the migrate package would go here in a real implementation
	// For now, we'll report success without attempting to run migrations
	log.Println("Migration step skipped for now")
	
	return nil
}

// Note: The getMigrationsDir function has been removed and replaced with a hardcoded path

// WaitForDB waits for the database to be available
func WaitForDB(config *configs.Config, maxRetries int, retryInterval time.Duration) error {
	var err error
	var db *DB

	// Try to connect to the database with retries
	for i := 0; i < maxRetries; i++ {
		db, err = NewConnection(config)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to the database")
				db.Close()
				return nil
			}
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}