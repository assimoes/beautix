package database

import (
	"errors"
	"time"

	"github.com/assimoes/beautix/configs"
	"github.com/rs/zerolog/log"
)

var db *DB

// InitDB initializes the database connection
func InitDB(config *configs.Config) error {
	if db != nil {
		return errors.New("database already initialized")
	}

	log.Info().Msg("Initializing database connection")

	// Wait for database to be available
	if err := WaitForDB(config, 5, time.Second*2); err != nil {
		return err
	}

	// Initialize database connection
	var err error
	db, err = Initialize(config)
	if err != nil {
		return err
	}

	log.Info().Msg("Database connection initialized")
	return nil
}

// GetDB returns the database instance
func GetDB() (*DB, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}
	return db, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if db == nil {
		return nil
	}

	err := db.Close()
	if err != nil {
		return err
	}

	db = nil
	log.Info().Msg("Database connection closed")
	return nil
}
