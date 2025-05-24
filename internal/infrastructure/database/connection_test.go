//go:build !integration
// +build !integration

package database_test

import (
	"context"
	"errors"
	"testing"

	"github.com/assimoes/beautix/configs"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewConnection(t *testing.T) {
	config, err := configs.LoadConfig()
	require.NoError(t, err, "Failed to load configuration")

	// Test creating a new connection
	db, err := database.NewConnection(config)
	assert.NoError(t, err, "Should connect to database without error")
	assert.NotNil(t, db, "Database connection should not be nil")

	// Test pinging the database
	err = db.Ping()
	assert.NoError(t, err, "Should ping database without error")

	// Close the connection
	err = db.Close()
	assert.NoError(t, err, "Should close database connection without error")
}

func TestWithTransaction(t *testing.T) {
	// Create a test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to create test database")

	// Create a test table using the new helper
	database.CreateTestTable(t, testDB.DB, "test_transactions", `CREATE TABLE test_transactions (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	)`)

	// Test successful transaction
	err = testDB.WithTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := tx.Exec("INSERT INTO test_transactions (name) VALUES ('test1')").Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO test_transactions (name) VALUES ('test2')").Error; err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err, "Transaction should succeed")

	// Check if records were inserted
	var count int64
	err = testDB.Model(&struct{}{}).Table("test_transactions").Count(&count).Error
	require.NoError(t, err, "Failed to count records")
	assert.Equal(t, int64(2), count, "Should have inserted 2 records")

	// Test failed transaction
	err = testDB.WithTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := tx.Exec("INSERT INTO test_transactions (name) VALUES ('test3')").Error; err != nil {
			return err
		}
		return errors.New("test error")
	})
	assert.Error(t, err, "Transaction should fail")

	// Check if records were not inserted
	err = testDB.Model(&struct{}{}).Table("test_transactions").Count(&count).Error
	require.NoError(t, err, "Failed to count records")
	assert.Equal(t, int64(2), count, "Should still have only 2 records")
}
