package main

const domainTemplate = `package domain

import (
	"context"
	"errors"
)

// {{.EntityName}} represents a {{.EntityNameLower}} in the system
type {{.EntityName}} struct {
	BaseModel
	// TODO: Add {{.EntityName}}-specific fields here
	// Example:
	// Name        string    ` + "`" + `gorm:"type:varchar(255);not null"` + "`" + `
	// Description *string   ` + "`" + `gorm:"type:text"` + "`" + `
	// IsActive    bool      ` + "`" + `gorm:"default:true"` + "`" + `
	// BusinessID  string    ` + "`" + `gorm:"type:uuid;not null;index"` + "`" + `
	// Business    *Business ` + "`" + `gorm:"foreignKey:BusinessID"` + "`" + `
}

// TableName overrides the table name for {{.EntityName}}
func ({{.EntityName}}) TableName() string {
	return "{{.EntityNameLower}}s"
}

// Validate performs validation on the {{.EntityName}} entity
func (e *{{.EntityName}}) Validate() error {
	// TODO: Add validation logic
	// Example:
	// if e.Name == "" {
	//     return errors.New("name is required")
	// }
	// if len(e.Name) > 255 {
	//     return errors.New("name cannot exceed 255 characters")
	// }
	return nil
}

// {{.EntityName}}Repository defines the repository interface for {{.EntityName}} operations
type {{.EntityName}}Repository interface {
	BaseRepository[{{.EntityName}}]
	// TODO: Add custom repository methods here
	// Example:
	// FindByName(ctx context.Context, name string) (*{{.EntityName}}, error)
	// FindByBusinessID(ctx context.Context, businessID string) ([]*{{.EntityName}}, error)
	// Search{{.EntityNamePlural}}(ctx context.Context, query string, limit int) ([]*{{.EntityName}}, error)
}

// {{.EntityName}}Service defines the service interface for {{.EntityName}} operations
// Note: The actual interface definition should be in the service package that imports both domain and dto
// This comment is here as a placeholder to document the expected service interface pattern:
//
// type {{.EntityName}}Service interface {
//     BaseService[dto.Create{{.EntityName}}DTO, dto.Update{{.EntityName}}DTO, dto.{{.EntityName}}ResponseDTO]
//     // Custom methods like:
//     // GetByName(ctx context.Context, name string) (*dto.{{.EntityName}}ResponseDTO, error)
//     // GetByBusiness(ctx context.Context, businessID string) ([]*dto.{{.EntityName}}ResponseDTO, error)
// }
`

const dtoTemplate = `package dto

import (
	"time"

	"github.com/assimoes/beautix/internal/domain"
)

// Create{{.EntityName}}DTO represents the data required to create a new {{.EntityNameLower}}
type Create{{.EntityName}}DTO struct {
	// TODO: Add fields required for creating a {{.EntityNameLower}}
	// Example:
	// Name        string  ` + "`" + `json:"name" validate:"required,min=2,max=255"` + "`" + `
	// Description *string ` + "`" + `json:"description,omitempty" validate:"omitempty,max=1000"` + "`" + `
	// BusinessID  string  ` + "`" + `json:"business_id" validate:"required,uuid"` + "`" + `
}

// Update{{.EntityName}}DTO represents the data that can be updated on a {{.EntityNameLower}}
type Update{{.EntityName}}DTO struct {
	// TODO: Add fields that can be updated (all should be pointers for partial updates)
	// Example:
	// Name        *string ` + "`" + `json:"name,omitempty" validate:"omitempty,min=2,max=255"` + "`" + `
	// Description *string ` + "`" + `json:"description,omitempty" validate:"omitempty,max=1000"` + "`" + `
	// IsActive    *bool   ` + "`" + `json:"is_active,omitempty"` + "`" + `
}

// {{.EntityName}}ResponseDTO represents the response data for a {{.EntityNameLower}}
type {{.EntityName}}ResponseDTO struct {
	BaseResponse
	// TODO: Add {{.EntityName}}-specific response fields
	// Example:
	// Name        string                ` + "`" + `json:"name"` + "`" + `
	// Description *string               ` + "`" + `json:"description,omitempty"` + "`" + `
	// IsActive    bool                  ` + "`" + `json:"is_active"` + "`" + `
	// BusinessID  string                ` + "`" + `json:"business_id"` + "`" + `
	// Business    *BusinessResponseDTO  ` + "`" + `json:"business,omitempty"` + "`" + `
}

// To{{.EntityName}}ResponseDTO converts a domain {{.EntityName}} to a {{.EntityName}}ResponseDTO
func To{{.EntityName}}ResponseDTO({{.EntityNameLower}} *domain.{{.EntityName}}) *{{.EntityName}}ResponseDTO {
	if {{.EntityNameLower}} == nil {
		return nil
	}

	return &{{.EntityName}}ResponseDTO{
		BaseResponse: BaseResponse{
			ID:        {{.EntityNameLower}}.ID,
			CreatedAt: {{.EntityNameLower}}.CreatedAt,
			UpdatedAt: {{.EntityNameLower}}.UpdatedAt,
		},
		// TODO: Map domain fields to DTO fields
		// Example:
		// Name:        {{.EntityNameLower}}.Name,
		// Description: {{.EntityNameLower}}.Description,
		// IsActive:    {{.EntityNameLower}}.IsActive,
		// BusinessID:  {{.EntityNameLower}}.BusinessID,
	}
}

// To{{.EntityNamePlural}}ResponseDTO converts a slice of domain {{.EntityNamePlural}} to {{.EntityName}}ResponseDTOs
func To{{.EntityNamePlural}}ResponseDTO({{.EntityNameLower}}s []*domain.{{.EntityName}}) []*{{.EntityName}}ResponseDTO {
	result := make([]*{{.EntityName}}ResponseDTO, len({{.EntityNameLower}}s))
	for i, {{.EntityNameLower}} := range {{.EntityNameLower}}s {
		result[i] = To{{.EntityName}}ResponseDTO({{.EntityNameLower}})
	}
	return result
}
`

const domainTestTemplate = `package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew{{.EntityName}}(t *testing.T) {
	{{.EntityNameLower}} := &{{.EntityName}}{
		// TODO: Set required fields
		// Example:
		// Name:       "Test {{.EntityName}}",
		// BusinessID: "test-business-id",
		// IsActive:   true,
	}

	assert.NotNil(t, {{.EntityNameLower}})
	// TODO: Add assertions for fields
	// Example:
	// assert.Equal(t, "Test {{.EntityName}}", {{.EntityNameLower}}.Name)
	// assert.Equal(t, "test-business-id", {{.EntityNameLower}}.BusinessID)
	// assert.True(t, {{.EntityNameLower}}.IsActive)
}

func Test{{.EntityName}}_TableName(t *testing.T) {
	{{.EntityNameLower}} := {{.EntityName}}{}
	assert.Equal(t, "{{.EntityNameLower}}s", {{.EntityNameLower}}.TableName())
}

func Test{{.EntityName}}_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *{{.EntityName}}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid {{.EntityNameLower}}",
			setup: func() *{{.EntityName}} {
				return &{{.EntityName}}{
					// TODO: Set valid fields
					// Example:
					// Name:       "Valid {{.EntityName}}",
					// BusinessID: "test-business-id",
				}
			},
			wantErr: false,
		},
		// TODO: Add validation test cases
		// Example:
		// {
		//     name: "empty name",
		//     setup: func() *{{.EntityName}} {
		//         return &{{.EntityName}}{
		//             Name:       "",
		//             BusinessID: "test-business-id",
		//         }
		//     },
		//     wantErr: true,
		//     errMsg:  "name is required",
		// },
		// {
		//     name: "name too long",
		//     setup: func() *{{.EntityName}} {
		//         return &{{.EntityName}}{
		//             Name:       strings.Repeat("a", 256),
		//             BusinessID: "test-business-id",
		//         }
		//     },
		//     wantErr: true,
		//     errMsg:  "name cannot exceed 255 characters",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			{{.EntityNameLower}} := tt.setup()
			err := {{.EntityNameLower}}.Validate()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
`

const dtoTestTemplate = `package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/assimoes/beautix/internal/domain"
)

func TestTo{{.EntityName}}ResponseDTO(t *testing.T) {
	t.Run("nil {{.EntityNameLower}}", func(t *testing.T) {
		result := To{{.EntityName}}ResponseDTO(nil)
		assert.Nil(t, result)
	})

	t.Run("valid {{.EntityNameLower}}", func(t *testing.T) {
		now := time.Now()
		{{.EntityNameLower}} := &domain.{{.EntityName}}{
			BaseModel: domain.BaseModel{
				ID:        "test-id",
				CreatedAt: now,
				UpdatedAt: now,
			},
			// TODO: Set domain fields
			// Example:
			// Name:        "Test {{.EntityName}}",
			// Description: StringPtr("Test description"),
			// IsActive:    true,
			// BusinessID:  "test-business-id",
		}

		result := To{{.EntityName}}ResponseDTO({{.EntityNameLower}})

		require.NotNil(t, result)
		assert.Equal(t, "test-id", result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		// TODO: Add assertions for mapped fields
		// Example:
		// assert.Equal(t, "Test {{.EntityName}}", result.Name)
		// assert.Equal(t, "Test description", *result.Description)
		// assert.True(t, result.IsActive)
		// assert.Equal(t, "test-business-id", result.BusinessID)
	})
}

func TestTo{{.EntityNamePlural}}ResponseDTO(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		result := To{{.EntityNamePlural}}ResponseDTO([]*domain.{{.EntityName}}{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("multiple {{.EntityNameLower}}s", func(t *testing.T) {
		now := time.Now()
		{{.EntityNameLower}}s := []*domain.{{.EntityName}}{
			{
				BaseModel: domain.BaseModel{
					ID:        "test-id-1",
					CreatedAt: now,
					UpdatedAt: now,
				},
				// TODO: Set fields for first {{.EntityNameLower}}
			},
			{
				BaseModel: domain.BaseModel{
					ID:        "test-id-2",
					CreatedAt: now,
					UpdatedAt: now,
				},
				// TODO: Set fields for second {{.EntityNameLower}}
			},
		}

		result := To{{.EntityNamePlural}}ResponseDTO({{.EntityNameLower}}s)

		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "test-id-1", result[0].ID)
		assert.Equal(t, "test-id-2", result[1].ID)
	})
}

// Helper function for string pointers
func StringPtr(s string) *string {
	return &s
}

// Helper function for bool pointers
func BoolPtr(b bool) *bool {
	return &b
}

// Helper function for int pointers
func IntPtr(i int) *int {
	return &i
}
`