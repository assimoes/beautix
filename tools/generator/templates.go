package main

const resolverTemplate = `package graph

import (
	"errors"

	"github.com/graphql-go/graphql"

	"github.com/assimoes/beautix/internal/dto"
	{{if .GenerateService}}"github.com/assimoes/beautix/internal/service"{{end}}
)

// {{.EntityName}} Query Resolvers
func (r *Resolver) resolve{{.EntityName}}(p graphql.ResolveParams) (any, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, errors.New("id is required")
	}

	{{if .GenerateService}}{{.EntityNameLower}}, err := r.{{.EntityNameLower}}Service.GetByID(p.Context, id)
	if err != nil {
		return nil, err
	}

	return {{.EntityNameLower}}, nil{{else}}// TODO: Implement {{.EntityName}} retrieval logic
	return nil, errors.New("not implemented"){{end}}
}

func (r *Resolver) resolve{{.EntityNamePlural}}(p graphql.ResolveParams) (any, error) {
	// Default pagination values
	page := 1
	pageSize := 10

	// Extract pagination parameters if provided
	if limit, ok := p.Args["limit"].(int); ok && limit > 0 {
		pageSize = limit
	}
	if offset, ok := p.Args["offset"].(int); ok && offset >= 0 {
		page = (offset / pageSize) + 1
	}

	{{if .GenerateService}}{{.EntityNameLower}}s, _, err := r.{{.EntityNameLower}}Service.List(p.Context, page, pageSize)
	if err != nil {
		return nil, err
	}

	return {{.EntityNameLower}}s, nil{{else}}// TODO: Implement {{.EntityNamePlural}} listing logic
	return nil, errors.New("not implemented"){{end}}
}

// {{.EntityName}} Mutation Resolvers
func (r *Resolver) resolveCreate{{.EntityName}}(p graphql.ResolveParams) (any, error) {
	input, ok := p.Args["input"].(map[string]any)
	if !ok {
		return nil, errors.New("input is required")
	}

	{{if .GenerateService}}createDTO := dto.Create{{.EntityName}}DTO{}

	// TODO: Map input fields to createDTO
	// Example:
	// if name, ok := input["name"].(string); ok {
	//     createDTO.Name = name
	// }

	{{.EntityNameLower}}, err := r.{{.EntityNameLower}}Service.Create(p.Context, createDTO)
	if err != nil {
		return nil, err
	}

	return {{.EntityNameLower}}, nil{{else}}// TODO: Implement {{.EntityName}} creation logic
	return nil, errors.New("not implemented"){{end}}
}

func (r *Resolver) resolveUpdate{{.EntityName}}(p graphql.ResolveParams) (any, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, errors.New("id is required")
	}

	input, ok := p.Args["input"].(map[string]any)
	if !ok {
		return nil, errors.New("input is required")
	}

	{{if .GenerateService}}updateDTO := dto.Update{{.EntityName}}DTO{}

	// TODO: Map input fields to updateDTO
	// Example:
	// if name, ok := input["name"].(string); ok {
	//     updateDTO.Name = &name
	// }

	{{.EntityNameLower}}, err := r.{{.EntityNameLower}}Service.Update(p.Context, id, updateDTO)
	if err != nil {
		return nil, err
	}

	return {{.EntityNameLower}}, nil{{else}}// TODO: Implement {{.EntityName}} update logic
	return nil, errors.New("not implemented"){{end}}
}

func (r *Resolver) resolveDelete{{.EntityName}}(p graphql.ResolveParams) (any, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, errors.New("id is required")
	}

	{{if .GenerateService}}err := r.{{.EntityNameLower}}Service.Delete(p.Context, id)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"success": true,
		"message": "{{.EntityName}} deleted successfully",
	}, nil{{else}}// TODO: Implement {{.EntityName}} deletion logic
	return nil, errors.New("not implemented"){{end}}
}
`

const serviceTemplate = `package service

import (
	"context"
	"strings"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/dto"
	"github.com/assimoes/beautix/internal/infrastructure/validation"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// {{.EntityName}}Service defines the service interface for {{.EntityName}} operations
type {{.EntityName}}Service interface {
	domain.BaseService[dto.Create{{.EntityName}}DTO, dto.Update{{.EntityName}}DTO, dto.{{.EntityName}}ResponseDTO]
	// TODO: Add custom service methods here
	// Example:
	// GetByName(ctx context.Context, name string) (*dto.{{.EntityName}}ResponseDTO, error)
	// Search{{.EntityNamePlural}}(ctx context.Context, query string, limit int) ([]*dto.{{.EntityName}}ResponseDTO, error)
}

// {{.EntityNameLower}}ServiceImpl implements the {{.EntityName}}Service interface
type {{.EntityNameLower}}ServiceImpl struct {
	*BaseServiceImpl[domain.{{.EntityName}}, dto.Create{{.EntityName}}DTO, dto.Update{{.EntityName}}DTO, dto.{{.EntityName}}ResponseDTO]
	{{.EntityNameLower}}Repo domain.{{.EntityName}}Repository
	validator    *validator.Validate
}

// New{{.EntityName}}Service creates a new {{.EntityNameLower}} service
func New{{.EntityName}}Service({{.EntityNameLower}}Repo domain.{{.EntityName}}Repository, validator *validator.Validate) {{.EntityName}}Service {
	return &{{.EntityNameLower}}ServiceImpl{
		BaseServiceImpl: NewBaseService(
			{{.EntityNameLower}}Repo,
			func(createDTO dto.Create{{.EntityName}}DTO) (*domain.{{.EntityName}}, error) {
				{{.EntityNameLower}} := &domain.{{.EntityName}}{
					// TODO: Map DTO fields to domain entity
					// Example:
					// Name:     strings.TrimSpace(createDTO.Name),
					// IsActive: true,
				}
				return {{.EntityNameLower}}, {{.EntityNameLower}}.Validate()
			},
			func(entity *domain.{{.EntityName}}, updateDTO dto.Update{{.EntityName}}DTO) error {
				// TODO: Map update DTO fields to domain entity
				// Example:
				// if updateDTO.Name != nil {
				//     entity.Name = strings.TrimSpace(*updateDTO.Name)
				// }
				// if updateDTO.IsActive != nil {
				//     entity.IsActive = *updateDTO.IsActive
				// }
				return entity.Validate()
			},
			func(entity *domain.{{.EntityName}}) *dto.{{.EntityName}}ResponseDTO {
				return dto.To{{.EntityName}}ResponseDTO(entity)
			},
		),
		{{.EntityNameLower}}Repo: {{.EntityNameLower}}Repo,
		validator:    validator,
	}
}

// TODO: Implement custom service methods here
// Example:
// func (s *{{.EntityNameLower}}ServiceImpl) GetByName(ctx context.Context, name string) (*dto.{{.EntityName}}ResponseDTO, error) {
//     if name == "" {
//         return nil, validation.NewValidationError("name is required")
//     }
//
//     {{.EntityNameLower}}, err := s.{{.EntityNameLower}}Repo.FindByName(ctx, strings.ToLower(name))
//     if err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             return nil, NewNotFoundError("{{.EntityNameLower}}", "name", name)
//         }
//         return nil, NewServiceError("failed to retrieve {{.EntityNameLower}} by name", err)
//     }
//
//     return dto.To{{.EntityName}}ResponseDTO({{.EntityNameLower}}), nil
// }
`

const repositoryTemplate = `package repository

import (
	"context"
	"strings"

	"github.com/assimoes/beautix/internal/domain"
	"gorm.io/gorm"
)

// {{.EntityNameLower}}RepositoryImpl implements the {{.EntityName}}Repository interface
type {{.EntityNameLower}}RepositoryImpl struct {
	*BaseRepositoryImpl[domain.{{.EntityName}}]
}

// New{{.EntityName}}Repository creates a new {{.EntityNameLower}} repository
func New{{.EntityName}}Repository(db *gorm.DB) domain.{{.EntityName}}Repository {
	return &{{.EntityNameLower}}RepositoryImpl{
		BaseRepositoryImpl: &BaseRepositoryImpl[domain.{{.EntityName}}]{db: db},
	}
}

// TODO: Implement custom repository methods here
// Example:
// func (r *{{.EntityNameLower}}RepositoryImpl) FindByName(ctx context.Context, name string) (*domain.{{.EntityName}}, error) {
//     var {{.EntityNameLower}} domain.{{.EntityName}}
//     err := r.db.WithContext(ctx).
//         Where("name = ?", strings.ToLower(name)).
//         First(&{{.EntityNameLower}}).Error
//     if err != nil {
//         return nil, err
//     }
//     return &{{.EntityNameLower}}, nil
// }
//
// func (r *{{.EntityNameLower}}RepositoryImpl) FindByBusinessID(ctx context.Context, businessID string) ([]*domain.{{.EntityName}}, error) {
//     var {{.EntityNameLower}}s []*domain.{{.EntityName}}
//     err := r.db.WithContext(ctx).
//         Where("business_id = ? AND deleted_at IS NULL", businessID).
//         Find(&{{.EntityNameLower}}s).Error
//     return {{.EntityNameLower}}s, err
// }
//
// func (r *{{.EntityNameLower}}RepositoryImpl) Search{{.EntityNamePlural}}(ctx context.Context, query string, limit int) ([]*domain.{{.EntityName}}, error) {
//     var {{.EntityNameLower}}s []*domain.{{.EntityName}}
//     err := r.db.WithContext(ctx).
//         Where("name ILIKE ? AND deleted_at IS NULL", "%"+query+"%").
//         Limit(limit).
//         Find(&{{.EntityNameLower}}s).Error
//     return {{.EntityNameLower}}s, err
// }
`

const resolverTestTemplate = `package graph

import (
	"context"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/assimoes/beautix/internal/dto"
	{{if .GenerateService}}"github.com/assimoes/beautix/internal/service"{{end}}
)

{{if .GenerateService}}// Mock implementations for testing
type mock{{.EntityName}}Service struct {
	{{.EntityNameLower}}s map[string]*dto.{{.EntityName}}ResponseDTO
}

func newMock{{.EntityName}}Service() *mock{{.EntityName}}Service {
	return &mock{{.EntityName}}Service{
		{{.EntityNameLower}}s: make(map[string]*dto.{{.EntityName}}ResponseDTO),
	}
}

func (m *mock{{.EntityName}}Service) Create(ctx context.Context, createDTO dto.Create{{.EntityName}}DTO) (*dto.{{.EntityName}}ResponseDTO, error) {
	{{.EntityNameLower}} := &dto.{{.EntityName}}ResponseDTO{
		BaseResponse: dto.BaseResponse{
			ID: "test-{{.EntityNameLower}}-1",
		},
		// TODO: Map createDTO fields to response
		// Example:
		// Name:     createDTO.Name,
		// IsActive: true,
	}
	m.{{.EntityNameLower}}s[{{.EntityNameLower}}.ID] = {{.EntityNameLower}}
	return {{.EntityNameLower}}, nil
}

func (m *mock{{.EntityName}}Service) GetByID(ctx context.Context, id string) (*dto.{{.EntityName}}ResponseDTO, error) {
	if {{.EntityNameLower}}, exists := m.{{.EntityNameLower}}s[id]; exists {
		return {{.EntityNameLower}}, nil
	}
	return nil, service.NewNotFoundError("{{.EntityNameLower}}", "id", id)
}

func (m *mock{{.EntityName}}Service) Update(ctx context.Context, id string, updateDTO dto.Update{{.EntityName}}DTO) (*dto.{{.EntityName}}ResponseDTO, error) {
	{{.EntityNameLower}}, exists := m.{{.EntityNameLower}}s[id]
	if !exists {
		return nil, service.NewNotFoundError("{{.EntityNameLower}}", "id", id)
	}

	// TODO: Apply updates from updateDTO
	// Example:
	// if updateDTO.Name != nil {
	//     {{.EntityNameLower}}.Name = *updateDTO.Name
	// }
	// if updateDTO.IsActive != nil {
	//     {{.EntityNameLower}}.IsActive = *updateDTO.IsActive
	// }

	return {{.EntityNameLower}}, nil
}

func (m *mock{{.EntityName}}Service) Delete(ctx context.Context, id string) error {
	if _, exists := m.{{.EntityNameLower}}s[id]; !exists {
		return service.NewNotFoundError("{{.EntityNameLower}}", "id", id)
	}
	delete(m.{{.EntityNameLower}}s, id)
	return nil
}

func (m *mock{{.EntityName}}Service) List(ctx context.Context, page, pageSize int) ([]*dto.{{.EntityName}}ResponseDTO, int64, error) {
	{{.EntityNameLower}}s := make([]*dto.{{.EntityName}}ResponseDTO, 0, len(m.{{.EntityNameLower}}s))
	for _, {{.EntityNameLower}} := range m.{{.EntityNameLower}}s {
		{{.EntityNameLower}}s = append({{.EntityNameLower}}s, {{.EntityNameLower}})
	}
	return {{.EntityNameLower}}s, int64(len({{.EntityNameLower}}s)), nil
}{{end}}

type mockAuth{{.EntityName}}Service struct{}

func (m *mockAuth{{.EntityName}}Service) RegisterWithClerk(ctx context.Context, clerkUserData dto.ClerkUserDTO) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func (m *mockAuth{{.EntityName}}Service) SyncClerkUser(ctx context.Context, clerkID string, userData dto.ClerkUserDTO) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func (m *mockAuth{{.EntityName}}Service) GetCurrentUser(ctx context.Context, clerkID string) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func (m *mockAuth{{.EntityName}}Service) VerifyToken(ctx context.Context, token string) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func (m *mockAuth{{.EntityName}}Service) GetUserFromContext(ctx context.Context) (*dto.UserResponseDTO, error) {
	return nil, nil
}

func setupTest{{.EntityName}}Schema() (graphql.Schema, {{if .GenerateService}}*mock{{.EntityName}}Service{{else}}*mockAuth{{.EntityName}}Service{{end}}) {
	{{if .GenerateService}}mock{{.EntityName}}Svc := newMock{{.EntityName}}Service(){{end}}
	mockAuthSvc := &mockAuth{{.EntityName}}Service{}

	resolver := NewResolver({{if not .GenerateService}}nil{{else}}nil{{end}}, mockAuthSvc)
	schema, err := CreateSchema(resolver)
	if err != nil {
		panic(err)
	}

	return schema, {{if .GenerateService}}mock{{.EntityName}}Svc{{else}}mockAuthSvc{{end}}
}

func TestGraphQL{{.EntityName}}Queries(t *testing.T) {
	schema, {{if .GenerateService}}mock{{.EntityName}}Svc{{else}}_{{end}} := setupTest{{.EntityName}}Schema()

	{{if .GenerateService}}// Create a test {{.EntityNameLower}} first
	test{{.EntityName}} := &dto.{{.EntityName}}ResponseDTO{
		BaseResponse: dto.BaseResponse{
			ID: "test-{{.EntityNameLower}}-1",
		},
		// TODO: Set test data fields
		// Example:
		// Name:     "Test {{.EntityName}}",
		// IsActive: true,
	}
	mock{{.EntityName}}Svc.{{.EntityNameLower}}s["test-{{.EntityNameLower}}-1"] = test{{.EntityName}}{{end}}

	t.Run("Query {{.EntityNameLower}} by ID", func(t *testing.T) {
		query := ` + "`" + `
			query {
				{{.EntityNameLower}}(id: "test-{{.EntityNameLower}}-1") {
					id
					# TODO: Add other fields to query
					# Example:
					# name
					# isActive
				}
			}
		` + "`" + `

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       context.Background(),
		})

		{{if .GenerateService}}require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		{{.EntityNameLower}}, ok := data["{{.EntityNameLower}}"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "test-{{.EntityNameLower}}-1", {{.EntityNameLower}}["id"])
		// TODO: Add assertions for other fields
		// Example:
		// assert.Equal(t, "Test {{.EntityName}}", {{.EntityNameLower}}["name"])
		// assert.Equal(t, true, {{.EntityNameLower}}["isActive"]){{else}}require.NotEmpty(t, result.Errors){{end}}
	})

	t.Run("Query {{.EntityNameLower}}s list", func(t *testing.T) {
		query := ` + "`" + `
			query {
				{{.EntityNameLower}}s {
					id
					# TODO: Add other fields to query
				}
			}
		` + "`" + `

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       context.Background(),
		})

		{{if .GenerateService}}require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		{{.EntityNameLower}}s, ok := data["{{.EntityNameLower}}s"].([]interface{})
		require.True(t, ok)
		assert.Len(t, {{.EntityNameLower}}s, 1)

		{{.EntityNameLower}}, ok := {{.EntityNameLower}}s[0].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "test-{{.EntityNameLower}}-1", {{.EntityNameLower}}["id"]){{else}}require.NotEmpty(t, result.Errors){{end}}
	})
}

func TestGraphQL{{.EntityName}}Mutations(t *testing.T) {
	schema, _ := setupTest{{.EntityName}}Schema()

	t.Run("Create {{.EntityNameLower}}", func(t *testing.T) {
		mutation := ` + "`" + `
			mutation {
				create{{.EntityName}}(input: {
					# TODO: Add required input fields
					# Example:
					# name: "New {{.EntityName}}"
				}) {
					id
					# TODO: Add other fields to return
				}
			}
		` + "`" + `

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: mutation,
			Context:       context.Background(),
		})

		{{if .GenerateService}}require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		{{.EntityNameLower}}, ok := data["create{{.EntityName}}"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "test-{{.EntityNameLower}}-1", {{.EntityNameLower}}["id"])
		// TODO: Add assertions for other fields{{else}}require.NotEmpty(t, result.Errors){{end}}
	})

	t.Run("Update {{.EntityNameLower}}", func(t *testing.T) {
		mutation := ` + "`" + `
			mutation {
				update{{.EntityName}}(id: "test-{{.EntityNameLower}}-1", input: {
					# TODO: Add fields to update
					# Example:
					# name: "Updated {{.EntityName}}"
				}) {
					id
					# TODO: Add other fields to return
				}
			}
		` + "`" + `

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: mutation,
			Context:       context.Background(),
		})

		{{if .GenerateService}}require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		{{.EntityNameLower}}, ok := data["update{{.EntityName}}"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "test-{{.EntityNameLower}}-1", {{.EntityNameLower}}["id"])
		// TODO: Add assertions for updated fields{{else}}require.NotEmpty(t, result.Errors){{end}}
	})

	t.Run("Delete {{.EntityNameLower}}", func(t *testing.T) {
		mutation := ` + "`" + `
			mutation {
				delete{{.EntityName}}(id: "test-{{.EntityNameLower}}-1") {
					success
					message
				}
			}
		` + "`" + `

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: mutation,
			Context:       context.Background(),
		})

		{{if .GenerateService}}require.Empty(t, result.Errors)

		data, ok := result.Data.(map[string]interface{})
		require.True(t, ok)

		deleteResult, ok := data["delete{{.EntityName}}"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, true, deleteResult["success"])
		assert.Equal(t, "{{.EntityName}} deleted successfully", deleteResult["message"]){{else}}require.NotEmpty(t, result.Errors){{end}}
	})
}
`

const serviceTestTemplate = `package service

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/dto"
)

// Mock{{.EntityName}}Repository is a mock implementation of {{.EntityName}}Repository
type Mock{{.EntityName}}Repository struct {
	mock.Mock
}

func (m *Mock{{.EntityName}}Repository) Create(ctx context.Context, entity *domain.{{.EntityName}}) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *Mock{{.EntityName}}Repository) GetByID(ctx context.Context, id string) (*domain.{{.EntityName}}, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.{{.EntityName}}), args.Error(1)
}

func (m *Mock{{.EntityName}}Repository) Update(ctx context.Context, id string, entity *domain.{{.EntityName}}) error {
	args := m.Called(ctx, id, entity)
	return args.Error(0)
}

func (m *Mock{{.EntityName}}Repository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *Mock{{.EntityName}}Repository) List(ctx context.Context, page, pageSize int) ([]*domain.{{.EntityName}}, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*domain.{{.EntityName}}), args.Get(1).(int64), args.Error(2)
}

// TODO: Add mock methods for custom repository methods
// Example:
// func (m *Mock{{.EntityName}}Repository) FindByName(ctx context.Context, name string) (*domain.{{.EntityName}}, error) {
//     args := m.Called(ctx, name)
//     return args.Get(0).(*domain.{{.EntityName}}), args.Error(1)
// }

func TestNew{{.EntityName}}Service(t *testing.T) {
	mockRepo := &Mock{{.EntityName}}Repository{}
	validator := validator.New()

	service := New{{.EntityName}}Service(mockRepo, validator)

	assert.NotNil(t, service)
}

func Test{{.EntityName}}Service_Create(t *testing.T) {
	mockRepo := &Mock{{.EntityName}}Repository{}
	validator := validator.New()
	service := New{{.EntityName}}Service(mockRepo, validator)

	createDTO := dto.Create{{.EntityName}}DTO{
		// TODO: Set required fields for creation
		// Example:
		// Name: "Test {{.EntityName}}",
	}

	expected{{.EntityName}} := &domain.{{.EntityName}}{
		BaseModel: domain.BaseModel{ID: "test-id"},
		// TODO: Set expected fields
		// Example:
		// Name:     "Test {{.EntityName}}",
		// IsActive: true,
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.{{.EntityName}}")).Return(nil)
	mockRepo.On("GetByID", mock.Anything, "test-id").Return(expected{{.EntityName}}, nil)

	result, err := service.Create(context.Background(), createDTO)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	// TODO: Add assertions for other fields

	mockRepo.AssertExpectations(t)
}

func Test{{.EntityName}}Service_GetByID(t *testing.T) {
	mockRepo := &Mock{{.EntityName}}Repository{}
	validator := validator.New()
	service := New{{.EntityName}}Service(mockRepo, validator)

	expected{{.EntityName}} := &domain.{{.EntityName}}{
		BaseModel: domain.BaseModel{ID: "test-id"},
		// TODO: Set expected fields
	}

	mockRepo.On("GetByID", mock.Anything, "test-id").Return(expected{{.EntityName}}, nil)

	result, err := service.GetByID(context.Background(), "test-id")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)

	mockRepo.AssertExpectations(t)
}

func Test{{.EntityName}}Service_Update(t *testing.T) {
	mockRepo := &Mock{{.EntityName}}Repository{}
	validator := validator.New()
	service := New{{.EntityName}}Service(mockRepo, validator)

	existing{{.EntityName}} := &domain.{{.EntityName}}{
		BaseModel: domain.BaseModel{ID: "test-id"},
		// TODO: Set existing fields
	}

	updateDTO := dto.Update{{.EntityName}}DTO{
		// TODO: Set fields to update
		// Example:
		// Name: StringPtr("Updated {{.EntityName}}"),
	}

	mockRepo.On("GetByID", mock.Anything, "test-id").Return(existing{{.EntityName}}, nil)
	mockRepo.On("Update", mock.Anything, "test-id", mock.AnythingOfType("*domain.{{.EntityName}}")).Return(nil)

	result, err := service.Update(context.Background(), "test-id", updateDTO)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)

	mockRepo.AssertExpectations(t)
}

func Test{{.EntityName}}Service_Delete(t *testing.T) {
	mockRepo := &Mock{{.EntityName}}Repository{}
	validator := validator.New()
	service := New{{.EntityName}}Service(mockRepo, validator)

	existing{{.EntityName}} := &domain.{{.EntityName}}{
		BaseModel: domain.BaseModel{ID: "test-id"},
	}

	mockRepo.On("GetByID", mock.Anything, "test-id").Return(existing{{.EntityName}}, nil)
	mockRepo.On("Delete", mock.Anything, "test-id").Return(nil)

	err := service.Delete(context.Background(), "test-id")

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func Test{{.EntityName}}Service_List(t *testing.T) {
	mockRepo := &Mock{{.EntityName}}Repository{}
	validator := validator.New()
	service := New{{.EntityName}}Service(mockRepo, validator)

	expected{{.EntityNamePlural}} := []*domain.{{.EntityName}}{
		{
			BaseModel: domain.BaseModel{ID: "test-id-1"},
			// TODO: Set fields
		},
		{
			BaseModel: domain.BaseModel{ID: "test-id-2"},
			// TODO: Set fields
		},
	}

	mockRepo.On("List", mock.Anything, 1, 10).Return(expected{{.EntityNamePlural}}, int64(2), nil)

	result, total, err := service.List(context.Background(), 1, 10)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)

	mockRepo.AssertExpectations(t)
}

// Helper function for string pointers in tests
func StringPtr(s string) *string {
	return &s
}

// Helper function for bool pointers in tests
func BoolPtr(b bool) *bool {
	return &b
}
`

const repositoryTestTemplate = `package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/infrastructure/database"
)

func TestNew{{.EntityName}}Repository(t *testing.T) {
	// Initialize test database
	database.InitTestDB()
	defer database.CloseDB()

	db, err := database.GetDB()
	require.NoError(t, err)

	repo := New{{.EntityName}}Repository(db.DB)
	assert.NotNil(t, repo)
}

func Test{{.EntityName}}Repository_Create(t *testing.T) {
	// Initialize test database
	database.InitTestDB()
	defer database.CloseDB()

	db, err := database.GetDB()
	require.NoError(t, err)

	repo := New{{.EntityName}}Repository(db.DB)

	{{.EntityNameLower}} := &domain.{{.EntityName}}{
		// TODO: Set required fields
		// Example:
		// Name:     "Test {{.EntityName}}",
		// IsActive: true,
	}

	err = repo.Create(context.Background(), {{.EntityNameLower}})
	require.NoError(t, err)
	assert.NotEmpty(t, {{.EntityNameLower}}.ID)
}

func Test{{.EntityName}}Repository_GetByID(t *testing.T) {
	// Initialize test database
	database.InitTestDB()
	defer database.CloseDB()

	db, err := database.GetDB()
	require.NoError(t, err)

	repo := New{{.EntityName}}Repository(db.DB)

	// Create test {{.EntityNameLower}}
	{{.EntityNameLower}} := &domain.{{.EntityName}}{
		// TODO: Set required fields
	}

	err = repo.Create(context.Background(), {{.EntityNameLower}})
	require.NoError(t, err)

	// Get {{.EntityNameLower}} by ID
	retrieved{{.EntityName}}, err := repo.GetByID(context.Background(), {{.EntityNameLower}}.ID)
	require.NoError(t, err)
	assert.Equal(t, {{.EntityNameLower}}.ID, retrieved{{.EntityName}}.ID)
	// TODO: Add assertions for other fields
}

func Test{{.EntityName}}Repository_Update(t *testing.T) {
	// Initialize test database
	database.InitTestDB()
	defer database.CloseDB()

	db, err := database.GetDB()
	require.NoError(t, err)

	repo := New{{.EntityName}}Repository(db.DB)

	// Create test {{.EntityNameLower}}
	{{.EntityNameLower}} := &domain.{{.EntityName}}{
		// TODO: Set required fields
	}

	err = repo.Create(context.Background(), {{.EntityNameLower}})
	require.NoError(t, err)

	// Update {{.EntityNameLower}}
	// TODO: Update fields
	// Example:
	// {{.EntityNameLower}}.Name = "Updated {{.EntityName}}"

	err = repo.Update(context.Background(), {{.EntityNameLower}}.ID, {{.EntityNameLower}})
	require.NoError(t, err)

	// Verify update
	updated{{.EntityName}}, err := repo.GetByID(context.Background(), {{.EntityNameLower}}.ID)
	require.NoError(t, err)
	// TODO: Add assertions for updated fields
	// Example:
	// assert.Equal(t, "Updated {{.EntityName}}", updated{{.EntityName}}.Name)
}

func Test{{.EntityName}}Repository_Delete(t *testing.T) {
	// Initialize test database
	database.InitTestDB()
	defer database.CloseDB()

	db, err := database.GetDB()
	require.NoError(t, err)

	repo := New{{.EntityName}}Repository(db.DB)

	// Create test {{.EntityNameLower}}
	{{.EntityNameLower}} := &domain.{{.EntityName}}{
		// TODO: Set required fields
	}

	err = repo.Create(context.Background(), {{.EntityNameLower}})
	require.NoError(t, err)

	// Delete {{.EntityNameLower}}
	err = repo.Delete(context.Background(), {{.EntityNameLower}}.ID)
	require.NoError(t, err)

	// Verify deletion (should return error)
	_, err = repo.GetByID(context.Background(), {{.EntityNameLower}}.ID)
	assert.Error(t, err)
}

func Test{{.EntityName}}Repository_List(t *testing.T) {
	// Initialize test database
	database.InitTestDB()
	defer database.CloseDB()

	db, err := database.GetDB()
	require.NoError(t, err)

	repo := New{{.EntityName}}Repository(db.DB)

	// Create test {{.EntityNameLower}}s
	{{.EntityNameLower}}1 := &domain.{{.EntityName}}{
		// TODO: Set required fields
	}
	{{.EntityNameLower}}2 := &domain.{{.EntityName}}{
		// TODO: Set required fields
	}

	err = repo.Create(context.Background(), {{.EntityNameLower}}1)
	require.NoError(t, err)
	err = repo.Create(context.Background(), {{.EntityNameLower}}2)
	require.NoError(t, err)

	// List {{.EntityNameLower}}s
	{{.EntityNameLower}}s, total, err := repo.List(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len({{.EntityNameLower}}s), 2)
	assert.GreaterOrEqual(t, total, int64(2))
}

// TODO: Add tests for custom repository methods
// Example:
// func Test{{.EntityName}}Repository_FindByName(t *testing.T) {
//     // Initialize test database
//     database.InitTestDB()
//     defer database.CloseDB()
//
//     db, err := database.GetDB()
//     require.NoError(t, err)
//
//     repo := New{{.EntityName}}Repository(db.DB)
//
//     // Create test {{.EntityNameLower}}
//     {{.EntityNameLower}} := &domain.{{.EntityName}}{
//         Name: "Test {{.EntityName}}",
//     }
//
//     err = repo.Create(context.Background(), {{.EntityNameLower}})
//     require.NoError(t, err)
//
//     // Find by name
//     found{{.EntityName}}, err := repo.FindByName(context.Background(), "Test {{.EntityName}}")
//     require.NoError(t, err)
//     assert.Equal(t, {{.EntityNameLower}}.ID, found{{.EntityName}}.ID)
//     assert.Equal(t, "Test {{.EntityName}}", found{{.EntityName}}.Name)
// }
`