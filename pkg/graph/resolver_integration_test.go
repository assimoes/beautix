//go:build integration
// +build integration

package graph_test

import (
	"testing"

	"github.com/assimoes/beautix/internal/repository"
	"github.com/assimoes/beautix/internal/service"
	"github.com/assimoes/beautix/pkg/graph"
	"github.com/assimoes/beautix/utils/testdb"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain ensures tests run sequentially
func TestMain(m *testing.M) {
	// Run tests sequentially to avoid database conflicts
	m.Run()
}

func TestUserQueriesIntegration(t *testing.T) {
	// Setup test database
	db := testdb.NewTestDB(t)
	gormDB := db.GetDB().DB // Get the embedded gorm.DB

	// Create validator
	validate := validator.New()

	// Create repositories
	userRepo := repository.NewUserRepository(gormDB)
	businessRepo := repository.NewBusinessRepository(gormDB)
	staffRepo := repository.NewStaffRepository(gormDB)

	// Create services
	userService := service.NewUserService(userRepo, businessRepo, staffRepo, validate)
	authService := service.NewAuthService(userRepo, nil, gormDB) // Using nil for clerk client in tests

	// Create resolver and schema
	resolver := graph.NewResolver(userService, authService)
	schema, err := graph.CreateSchema(resolver)
	require.NoError(t, err)

	// Create test helper
	helper := testdb.NewGraphQLTestHelper(t, schema)

	// Create fixtures
	fixtures := testdb.NewFixtureBuilder(db.GetDB())

	t.Run("Query user by ID", func(t *testing.T) {
		// Create test user
		testUser, err := fixtures.NewUser().
			WithEmail("test@example.com").
			WithName("John", "Doe").
			Build()
		require.NoError(t, err)

		// Query for the user
		query := `
			query($id: String!) {
				user(id: $id) {
					id
					email
					firstName
					lastName
					fullName
					isActive
				}
			}
		`

		data := helper.MustExecuteQuery(query, map[string]interface{}{
			"id": testUser.ID,
		})

		user := data["user"].(map[string]interface{})
		assert.Equal(t, testUser.ID, user["id"])
		assert.Equal(t, "test@example.com", user["email"])
		assert.Equal(t, "John", user["firstName"])
		assert.Equal(t, "Doe", user["lastName"])
		assert.Equal(t, "John Doe", user["fullName"])
		assert.Equal(t, true, user["isActive"])
	})

	t.Run("Query user by email", func(t *testing.T) {
		// Create test user
		testUser, err := fixtures.NewUser().
			WithEmail("byemail@example.com").
			WithName("Jane", "Smith").
			Build()
		require.NoError(t, err)

		// Query for the user by email
		query := `
			query($email: String!) {
				userByEmail(email: $email) {
					id
					email
					firstName
					lastName
				}
			}
		`

		data := helper.MustExecuteQuery(query, map[string]interface{}{
			"email": "byemail@example.com",
		})

		user := data["userByEmail"].(map[string]interface{})
		assert.Equal(t, testUser.ID, user["id"])
		assert.Equal(t, "byemail@example.com", user["email"])
		assert.Equal(t, "Jane", user["firstName"])
		assert.Equal(t, "Smith", user["lastName"])
	})

	t.Run("List users with pagination", func(t *testing.T) {
		// Create multiple test users
		for i := 0; i < 5; i++ {
			_, err := fixtures.CreateUser()
			require.NoError(t, err)
		}

		// Query for users with pagination
		query := `
			query {
				users(limit: 3, offset: 0) {
					id
					email
				}
			}
		`

		data := helper.MustExecuteQuery(query, nil)

		users := data["users"].([]interface{})
		assert.GreaterOrEqual(t, len(users), 3)
	})

	t.Run("Search users", func(t *testing.T) {
		// Create test user with specific name
		_, err := fixtures.NewUser().
			WithEmail("searchtest@example.com").
			WithName("SearchFirst", "SearchLast").
			Build()
		require.NoError(t, err)

		// Search for users
		query := `
			query($query: String!) {
				searchUsers(query: $query) {
					id
					firstName
					lastName
					email
				}
			}
		`

		data := helper.MustExecuteQuery(query, map[string]interface{}{
			"query": "SearchFirst",
		})

		users := data["searchUsers"].([]interface{})
		assert.GreaterOrEqual(t, len(users), 1)

		// Verify the searched user is in results
		found := false
		for _, u := range users {
			user := u.(map[string]interface{})
			if user["firstName"] == "SearchFirst" {
				found = true
				assert.Equal(t, "SearchLast", user["lastName"])
				assert.Equal(t, "searchtest@example.com", user["email"])
				break
			}
		}
		assert.True(t, found, "Expected user not found in search results")
	})

	t.Run("Query non-existent user", func(t *testing.T) {
		query := `
			query {
				user(id: "00000000-0000-0000-0000-000000000000") {
					id
					email
				}
			}
		`

		errors := helper.ExpectError(query, nil, "not found")
		assert.NotEmpty(t, errors)
	})
}

func TestUserMutationsIntegration(t *testing.T) {
	// Setup test database
	db := testdb.NewTestDB(t)
	gormDB := db.GetDB().DB // Get the embedded gorm.DB

	// Create validator
	validate := validator.New()

	// Create repositories
	userRepo := repository.NewUserRepository(gormDB)
	businessRepo := repository.NewBusinessRepository(gormDB)
	staffRepo := repository.NewStaffRepository(gormDB)

	// Create services
	userService := service.NewUserService(userRepo, businessRepo, staffRepo, validate)
	authService := service.NewAuthService(userRepo, nil, gormDB)

	// Create resolver and schema
	resolver := graph.NewResolver(userService, authService)
	schema, err := graph.CreateSchema(resolver)
	require.NoError(t, err)

	// Create test helper
	helper := testdb.NewGraphQLTestHelper(t, schema)

	// Create fixtures
	fixtures := testdb.NewFixtureBuilder(db.GetDB())

	t.Run("Create user", func(t *testing.T) {
		mutation := `
			mutation($input: CreateUserInput!) {
				createUser(input: $input) {
					id
					email
					firstName
					lastName
					fullName
					isActive
				}
			}
		`

		data := helper.MustExecuteQuery(mutation, map[string]interface{}{
			"input": map[string]interface{}{
				"email":     "newuser@example.com",
				"firstName": "New",
				"lastName":  "User",
				"phone":     "+1234567890",
			},
		})

		user := data["createUser"].(map[string]interface{})
		assert.NotEmpty(t, user["id"])
		assert.Equal(t, "newuser@example.com", user["email"])
		assert.Equal(t, "New", user["firstName"])
		assert.Equal(t, "User", user["lastName"])
		assert.Equal(t, "New User", user["fullName"])
		assert.Equal(t, true, user["isActive"])
	})

	t.Run("Update user", func(t *testing.T) {
		// Create initial user
		testUser, err := fixtures.NewUser().
			WithEmail("updatetest@example.com").
			WithName("Original", "Name").
			Build()
		require.NoError(t, err)

		// Update the user
		mutation := `
			mutation($id: String!, $input: UpdateUserInput!) {
				updateUser(id: $id, input: $input) {
					id
					firstName
					lastName
					fullName
				}
			}
		`

		data := helper.MustExecuteQuery(mutation, map[string]interface{}{
			"id": testUser.ID,
			"input": map[string]interface{}{
				"firstName": "Updated",
				"lastName":  "UserName",
			},
		})

		user := data["updateUser"].(map[string]interface{})
		assert.Equal(t, testUser.ID, user["id"])
		assert.Equal(t, "Updated", user["firstName"])
		assert.Equal(t, "UserName", user["lastName"])
		assert.Equal(t, "Updated UserName", user["fullName"])
	})

	t.Run("Delete user", func(t *testing.T) {
		// Create user to delete
		testUser, err := fixtures.NewUser().
			WithEmail("deletetest@example.com").
			Build()
		require.NoError(t, err)

		// Delete the user
		mutation := `
			mutation($id: String!) {
				deleteUser(id: $id) {
					success
					message
				}
			}
		`

		data := helper.MustExecuteQuery(mutation, map[string]interface{}{
			"id": testUser.ID,
		})

		result := data["deleteUser"].(map[string]interface{})
		assert.Equal(t, true, result["success"])
		assert.Equal(t, "User deleted successfully", result["message"])

		// Verify user is deleted by trying to query it
		query := `
			query($id: String!) {
				user(id: $id) {
					id
				}
			}
		`

		errors := helper.ExpectError(query, map[string]interface{}{
			"id": testUser.ID,
		}, "not found")
		assert.NotEmpty(t, errors)
	})

	// Note: activateUser and deactivateUser mutations are not implemented in the current schema
	// This test can be re-enabled when those mutations are added
	/*
		t.Run("Activate and deactivate user", func(t *testing.T) {
			// Test disabled - activateUser/deactivateUser mutations not in schema
		})
	*/
}

// TestBusinessRelatedQueriesIntegration is commented out for now
// The userWithBusinesses query doesn't exist in the current schema
// This test can be re-enabled when the query is implemented
/*
func TestBusinessRelatedQueriesIntegration(t *testing.T) {
	// This test is disabled because userWithBusinesses query is not implemented
}
*/

func TestErrorHandlingIntegration(t *testing.T) {
	// Setup test database
	db := testdb.NewTestDB(t)
	gormDB := db.GetDB().DB // Get the embedded gorm.DB

	// Create validator
	validate := validator.New()

	// Create repositories
	userRepo := repository.NewUserRepository(gormDB)
	businessRepo := repository.NewBusinessRepository(gormDB)
	staffRepo := repository.NewStaffRepository(gormDB)

	// Create services
	userService := service.NewUserService(userRepo, businessRepo, staffRepo, validate)
	authService := service.NewAuthService(userRepo, nil, gormDB)

	// Create resolver and schema
	resolver := graph.NewResolver(userService, authService)
	schema, err := graph.CreateSchema(resolver)
	require.NoError(t, err)

	// Create test helper
	helper := testdb.NewGraphQLTestHelper(t, schema)

	// Create fixtures
	fixtures := testdb.NewFixtureBuilder(db.GetDB())

	t.Run("Create user with duplicate email", func(t *testing.T) {
		// Create initial user
		_, err := fixtures.NewUser().
			WithEmail("duplicate@example.com").
			Build()
		require.NoError(t, err)

		// Try to create another user with same email
		mutation := `
			mutation {
				createUser(input: {
					email: "duplicate@example.com"
					firstName: "Duplicate"
					lastName: "User"
				}) {
					id
				}
			}
		`

		errors := helper.ExpectError(mutation, nil)
		assert.NotEmpty(t, errors)
	})

	t.Run("Update non-existent user", func(t *testing.T) {
		mutation := `
			mutation {
				updateUser(id: "00000000-0000-0000-0000-000000000000", input: {
					firstName: "Test"
				}) {
					id
				}
			}
		`

		errors := helper.ExpectError(mutation, nil, "not found")
		assert.NotEmpty(t, errors)
	})

	t.Run("Invalid query syntax", func(t *testing.T) {
		query := `
			query {
				user {
					id
				}
			}
		`

		// This should fail because 'id' argument is required
		errors := helper.ExpectError(query, nil)
		assert.NotEmpty(t, errors)
	})
}
