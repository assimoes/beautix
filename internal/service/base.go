package service

import (
	"context"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/assimoes/beautix/internal/dto"
)

// Context keys for user information
type ContextKey string

const (
	UserIDKey     ContextKey = "user_id"
	UserRoleKey   ContextKey = "user_role"
	BusinessIDKey ContextKey = "business_id"
	ClerkUserKey  ContextKey = "clerk_user"
)

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) *string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		return &userID
	}
	return nil
}

// GetUserRoleFromContext extracts user role from context
func GetUserRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(UserRoleKey).(string); ok {
		return role
	}
	return ""
}

// GetBusinessIDFromContext extracts business ID from context
func GetBusinessIDFromContext(ctx context.Context) *string {
	if businessID, ok := ctx.Value(BusinessIDKey).(string); ok && businessID != "" {
		return &businessID
	}
	return nil
}

// GetClerkUserFromContext extracts Clerk user from context
func GetClerkUserFromContext(ctx context.Context) any {
	return ctx.Value(ClerkUserKey)
}

// SetUserContext creates a new context with user information
func SetUserContext(ctx context.Context, userID, role string, businessID *string) context.Context {
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, UserRoleKey, role)
	if businessID != nil {
		ctx = context.WithValue(ctx, BusinessIDKey, *businessID)
	}
	return ctx
}

// ValidatePagination validates and sets default pagination values
func ValidatePagination(pagination *dto.PaginationRequest) *dto.PaginationRequest {
	if pagination == nil {
		return dto.DefaultPagination()
	}
	
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	
	if pagination.PageSize < 1 {
		pagination.PageSize = 20
	} else if pagination.PageSize > 100 {
		pagination.PageSize = 100
	}
	
	return pagination
}

// BaseServiceImpl provides a base implementation for services
type BaseServiceImpl[T, CreateDTO, UpdateDTO, ResponseDTO any] struct {
	repo             domain.BaseRepository[T]
	createConverter  func(CreateDTO) (*T, error)
	updateConverter  func(*T, UpdateDTO) error
	responseConverter func(*T) *ResponseDTO
}

// NewBaseService creates a new base service
func NewBaseService[T, CreateDTO, UpdateDTO, ResponseDTO any](
	repo domain.BaseRepository[T],
	createConverter func(CreateDTO) (*T, error),
	updateConverter func(*T, UpdateDTO) error,
	responseConverter func(*T) *ResponseDTO,
) *BaseServiceImpl[T, CreateDTO, UpdateDTO, ResponseDTO] {
	return &BaseServiceImpl[T, CreateDTO, UpdateDTO, ResponseDTO]{
		repo:              repo,
		createConverter:   createConverter,
		updateConverter:   updateConverter,
		responseConverter: responseConverter,
	}
}

// Create creates a new entity
func (s *BaseServiceImpl[T, CreateDTO, UpdateDTO, ResponseDTO]) Create(ctx context.Context, createDTO CreateDTO) (*ResponseDTO, error) {
	entity, err := s.createConverter(createDTO)
	if err != nil {
		return nil, err
	}
	
	err = s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	
	return s.responseConverter(entity), nil
}

// GetByID retrieves an entity by ID
func (s *BaseServiceImpl[T, CreateDTO, UpdateDTO, ResponseDTO]) GetByID(ctx context.Context, id string) (*ResponseDTO, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return s.responseConverter(entity), nil
}

// Update updates an existing entity
func (s *BaseServiceImpl[T, CreateDTO, UpdateDTO, ResponseDTO]) Update(ctx context.Context, id string, updateDTO UpdateDTO) (*ResponseDTO, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	err = s.updateConverter(entity, updateDTO)
	if err != nil {
		return nil, err
	}
	
	err = s.repo.Update(ctx, entity)
	if err != nil {
		return nil, err
	}
	
	return s.responseConverter(entity), nil
}

// Delete deletes an entity by ID
func (s *BaseServiceImpl[T, CreateDTO, UpdateDTO, ResponseDTO]) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves entities with pagination
func (s *BaseServiceImpl[T, CreateDTO, UpdateDTO, ResponseDTO]) List(ctx context.Context, page, pageSize int) ([]*ResponseDTO, int64, error) {
	entities, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	
	results := make([]*ResponseDTO, len(entities))
	for i, entity := range entities {
		results[i] = s.responseConverter(entity)
	}
	
	return results, total, nil
}