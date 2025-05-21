package database_test

import (
	"testing"

	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTestTable(t *testing.T) {
	// Create a test database
	testDB, err := database.NewTestDB(t)
	require.NoError(t, err, "Failed to create test database")

	testTableName := "test_cleanup_table"
	testTableSchema := `CREATE TABLE test_cleanup_table (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	)`

	t.Run("VerifyTableCreation", func(t *testing.T) {
		// Step 1: Verify the table doesn't exist initially
		exists, err := database.TableExists(testDB.DB, testTableName)
		require.NoError(t, err, "Failed to check if table exists")
		assert.False(t, exists, "Table should not exist before test")

		// Step 2: Create the test table with cleanup registered for this subtest
		database.CreateTestTable(t, testDB.DB, testTableName, testTableSchema)

		// Step 3: Verify the table exists after creation
		exists, err = database.TableExists(testDB.DB, testTableName)
		require.NoError(t, err, "Failed to check if table exists")
		assert.True(t, exists, "Table should exist after creation")
	})

	// Step 4: Verify the table no longer exists after the subtest completes
	// (t.Cleanup runs when the test or subtest completes)
	existsAfter, errAfter := database.TableExists(testDB.DB, testTableName)
	require.NoError(t, errAfter, "Failed to check if table exists")
	assert.False(t, existsAfter, "Table should be dropped during cleanup")
}