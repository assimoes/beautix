package graph

import (
	"context"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/assimoes/beautix/internal/dto"
	"github.com/assimoes/beautix/internal/service"
)

// Mock implementations for testing
type mockUserService struct {
	users map[string]*dto.UserResponseDTO
}

func newMockUserService() *mockUserService {
	return &mockUserService{
		users: make(map[string]*dto.UserResponseDTO),
	}
}

func (m *mockUserService) Create(ctx context.Context, userDto dto.CreateUserDTO) (*dto.UserResponseDTO, error) {
	user := &dto.UserResponseDTO{
		BaseResponse: dto.BaseResponse{
			ID: "test-user-1",
		},
		Email:     userDto.Email,
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Phone:     userDto.Phone,
		IsActive:  true,
		FullName:  userDto.FirstName + " " + userDto.LastName,
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *mockUserService) GetByID(ctx context.Context, id string) (*dto.UserResponseDTO, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, service.NewNotFoundError("user", "id", id)
}

func (m *mockUserService) Update(ctx context.Context, id string, updateDTO dto.UpdateUserDTO) (*dto.UserResponseDTO, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, service.NewNotFoundError("user", "id", id)
	}

	if updateDTO.FirstName != nil {
		user.FirstName = *updateDTO.FirstName
	}
	if updateDTO.LastName != nil {
		user.LastName = *updateDTO.LastName
	}
	if updateDTO.Phone != nil {
		user.Phone = updateDTO.Phone
	}
	if updateDTO.IsActive != nil {
		user.IsActive = *updateDTO.IsActive
	}

	user.FullName = user.FirstName + " " + user.LastName
	return user, nil
}

func (m *mockUserService) Delete(ctx context.Context, id string) error {
	if _, exists := m.users[id]; !exists {
		return service.NewNotFoundError("user", "id", id)
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserService) List(ctx context.Context, page, pageSize int) ([]*dto.UserResponseDTO, int64, error) {
	users := make([]*dto.UserResponseDTO, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, int64(len(users)), nil
}

func (m *mockUserService) GetByEmail(ctx context.Context, email string) (*dto.UserResponseDTO, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, service.NewNotFoundError("user", "email", email)
}

func (m *mockUserService) GetByClerkID(ctx context.Context, clerkID string) (*dto.UserResponseDTO, error) {
	for _, user := range m.users {
		if user.ClerkID != nil && *user.ClerkID == clerkID {
			return user, nil
		}
	}
	return nil, service.NewNotFoundError("user", "clerk_id", clerkID)
}

func (m *mockUserService) UpdateClerkID(ctx context.Context, userID, clerkID string) error {
	if _, exists := m.users[userID]; !exists {
		return service.NewNotFoundError("user", "id", userID)
	}
	m.users[userID].ClerkID = &clerkID
	return nil
}

func (m *mockUserService) GetWithBusinesses(ctx context.Context, userID string) (*dto.UserWithBusinessesDTO, error) {
	user, err := m.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.UserWithBusinessesDTO{
		UserResponseDTO: *user,
		Businesses:      []*dto.BusinessResponseDTO{},
	}, nil
}

func (m *mockUserService) SearchUsers(ctx context.Context, query string, limit int) ([]*dto.UserResponseDTO, error) {
	var results []*dto.UserResponseDTO
	for _, user := range m.users {
		if user.FirstName == query || user.LastName == query || user.Email == query {
			results = append(results, user)
		}
	}
	return results, nil
}

func (m *mockUserService) ActivateUser(ctx context.Context, userID string) error {
	if user, exists := m.users[userID]; exists {
		user.IsActive = true
		return nil
	}
	return service.NewNotFoundError("user", "id", userID)
}

func (m *mockUserService) DeactivateUser(ctx context.Context, userID string) error {
	if user, exists := m.users[userID]; exists {
		user.IsActive = false
		return nil
	}
	return service.NewNotFoundError("user", "id", userID)
}

type mockAuthService struct{}

func (m *mockAuthService) RegisterWithClerk(ctx context.Context, clerkUserData dto.ClerkUserDTO) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func (m *mockAuthService) SyncClerkUser(ctx context.Context, clerkID string, userData dto.ClerkUserDTO) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func (m *mockAuthService) GetCurrentUser(ctx context.Context, clerkID string) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func (m *mockAuthService) VerifyToken(ctx context.Context, token string) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func (m *mockAuthService) GetUserFromContext(ctx context.Context) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func setupTestSchema() (graphql.Schema, *mockUserService) {
	mockUserSvc := newMockUserService()
	mockAuthSvc := &mockAuthService{}

	resolver := NewResolver(mockUserSvc, mockAuthSvc)
	schema, err := CreateSchema(resolver)
	if err != nil {
		panic(err)
	}

	return schema, mockUserSvc
}

func TestGraphQLUserQueries(t *testing.T) {
	schema, mockUserSvc := setupTestSchema()

	// Create a test user first
	testUser := &dto.UserResponseDTO{
		BaseResponse: dto.BaseResponse{
			ID: "test-user-1",
		},
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		IsActive:  true,
		FullName:  "John Doe",
	}
	mockUserSvc.users["test-user-1"] = testUser

	t.Run("Query user by ID", func(t *testing.T) {
		query := `
			query {
				user(id: "test-user-1") {
					id
					email
					firstName
					lastName
					fullName
					isActive
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       context.Background(),
		})

		require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		user, ok := data["user"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "test-user-1", user["id"])
		assert.Equal(t, "test@example.com", user["email"])
		assert.Equal(t, "John", user["firstName"])
		assert.Equal(t, "Doe", user["lastName"])
		assert.Equal(t, "John Doe", user["fullName"])
		assert.Equal(t, true, user["isActive"])
	})

	t.Run("Query user by email", func(t *testing.T) {
		query := `
			query {
				userByEmail(email: "test@example.com") {
					id
					email
					firstName
					lastName
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       context.Background(),
		})

		require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		user, ok := data["userByEmail"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "test-user-1", user["id"])
		assert.Equal(t, "test@example.com", user["email"])
	})

	t.Run("Query users list", func(t *testing.T) {
		query := `
			query {
				users {
					id
					email
					firstName
					lastName
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       context.Background(),
		})

		require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		users, ok := data["users"].([]interface{})
		require.True(t, ok)
		assert.Len(t, users, 1)

		user, ok := users[0].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "test-user-1", user["id"])
	})

	t.Run("Search users", func(t *testing.T) {
		query := `
			query {
				searchUsers(query: "John") {
					id
					firstName
					lastName
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       context.Background(),
		})

		require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		users, ok := data["searchUsers"].([]interface{})
		require.True(t, ok)
		assert.Len(t, users, 1)
	})
}

func TestGraphQLUserMutations(t *testing.T) {
	schema, _ := setupTestSchema()

	t.Run("Create user", func(t *testing.T) {
		mutation := `
			mutation {
				createUser(input: {
					email: "newuser@example.com"
					firstName: "Jane"
					lastName: "Smith"
				}) {
					id
					email
					firstName
					lastName
					fullName
					isActive
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: mutation,
			Context:       context.Background(),
		})

		require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		user, ok := data["createUser"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "test-user-1", user["id"])
		assert.Equal(t, "newuser@example.com", user["email"])
		assert.Equal(t, "Jane", user["firstName"])
		assert.Equal(t, "Smith", user["lastName"])
		assert.Equal(t, "Jane Smith", user["fullName"])
		assert.Equal(t, true, user["isActive"])
	})

	t.Run("Update user", func(t *testing.T) {
		// First create a user
		createMutation := `
			mutation {
				createUser(input: {
					email: "update@example.com"
					firstName: "Update"
					lastName: "Test"
				}) {
					id
				}
			}
		`

		createResult := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: createMutation,
			Context:       context.Background(),
		})
		require.Empty(t, createResult.Errors)

		// Then update the user
		updateMutation := `
			mutation {
				updateUser(id: "test-user-1", input: {
					firstName: "Updated"
					lastName: "User"
				}) {
					id
					firstName
					lastName
					fullName
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: updateMutation,
			Context:       context.Background(),
		})

		require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		user, ok := data["updateUser"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "Updated", user["firstName"])
		assert.Equal(t, "User", user["lastName"])
		assert.Equal(t, "Updated User", user["fullName"])
	})

	t.Run("Delete user", func(t *testing.T) {
		// First create a user
		createMutation := `
			mutation {
				createUser(input: {
					email: "delete@example.com"
					firstName: "Delete"
					lastName: "Test"
				}) {
					id
				}
			}
		`

		createResult := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: createMutation,
			Context:       context.Background(),
		})
		require.Empty(t, createResult.Errors)

		// Then delete the user
		deleteMutation := `
			mutation {
				deleteUser(id: "test-user-1") {
					success
					message
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: deleteMutation,
			Context:       context.Background(),
		})

		require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		deleteResult, ok := data["deleteUser"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, true, deleteResult["success"])
		assert.Equal(t, "User deleted successfully", deleteResult["message"])
	})
}

func TestGraphQLErrorHandling(t *testing.T) {
	schema, _ := setupTestSchema()

	t.Run("Query non-existent user", func(t *testing.T) {
		query := `
			query {
				user(id: "non-existent") {
					id
					email
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       context.Background(),
		})

		require.NotEmpty(t, result.Errors)
		// GraphQL may return data with null values when there are errors
		if result.Data != nil {
			data, ok := result.Data.(map[string]interface{})
			if ok {
				assert.Nil(t, data["user"])
			}
		}
	})

	t.Run("Create user with missing required field", func(t *testing.T) {
		mutation := `
			mutation {
				createUser(input: {
					firstName: "Missing"
					lastName: "Email"
				}) {
					id
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: mutation,
			Context:       context.Background(),
		})

		// GraphQL schema validation should catch missing required fields
		require.NotEmpty(t, result.Errors)
		assert.Contains(t, result.Errors[0].Message, "email")
	})

	t.Run("Query with missing required argument", func(t *testing.T) {
		query := `
			query {
				user {
					id
					email
				}
			}
		`

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       context.Background(),
		})

		require.NotEmpty(t, result.Errors)
		assert.Contains(t, result.Errors[0].Message, "id")
	})
}
