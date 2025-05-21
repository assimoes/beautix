package database

import (
	"context"
	"fmt"
	"time"

	"github.com/assimoes/beautix/configs"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is a wrapper around gorm.DB that provides additional functionality
type DB struct {
	*gorm.DB
}

// NewConnection creates a new database connection
func NewConnection(config *configs.Config) (*DB, error) {
	return NewConnectionWithConfig(config, nil)
}

// NewConnectionWithConfig creates a new database connection with a custom GORM config
func NewConnectionWithConfig(config *configs.Config, gormConfig *gorm.Config) (*DB, error) {
	if gormConfig == nil {
		logLevel := logger.Silent
		if config.IsDevelopment() {
			logLevel = logger.Info
		}

		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		}
	}

	db, err := gorm.Open(postgres.Open(config.Database.URL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &DB{DB: db}, nil
}

// Ping checks if the database connection is alive
func (db *DB) Ping() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// WithTransaction executes a function within a transaction
func (db *DB) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Interface("recover", r).Msg("Recovered from panic in transaction")
			panic(r) // Re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}