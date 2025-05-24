package graph

import (
	"errors"

	"github.com/graphql-go/graphql"

	"github.com/assimoes/beautix/internal/dto"
	"github.com/assimoes/beautix/internal/service"
)

// Resolver contains the GraphQL resolvers
type Resolver struct {
	userService service.UserService
	authService service.AuthService
}

// NewResolver creates a new GraphQL resolver
func NewResolver(userService service.UserService, authService service.AuthService) *Resolver {
	return &Resolver{
		userService: userService,
		authService: authService,
	}
}

// User Query Resolvers
func (r *Resolver) resolveUser(p graphql.ResolveParams) (any, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, errors.New("id is required")
	}

	user, err := r.userService.GetByID(p.Context, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Resolver) resolveUsers(p graphql.ResolveParams) (any, error) {
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

	users, _, err := r.userService.List(p.Context, page, pageSize)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Resolver) resolveUserByEmail(p graphql.ResolveParams) (any, error) {
	email, ok := p.Args["email"].(string)
	if !ok {
		return nil, errors.New("email is required")
	}

	user, err := r.userService.GetByEmail(p.Context, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Resolver) resolveSearchUsers(p graphql.ResolveParams) (any, error) {
	query, ok := p.Args["query"].(string)
	if !ok {
		return nil, errors.New("query is required")
	}

	limit := 50
	if l, ok := p.Args["limit"].(int); ok && l > 0 {
		limit = l
	}

	users, err := r.userService.SearchUsers(p.Context, query, limit)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Resolver) resolveCurrentUser(p graphql.ResolveParams) (any, error) {
	// Extract user from context (set by authentication middleware)
	userID := service.GetUserIDFromContext(p.Context)
	if userID == nil {
		return nil, errors.New("user not authenticated")
	}

	user, err := r.userService.GetByID(p.Context, *userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// User Mutation Resolvers
func (r *Resolver) resolveCreateUser(p graphql.ResolveParams) (any, error) {
	input, ok := p.Args["input"].(map[string]any)
	if !ok {
		return nil, errors.New("input is required")
	}

	createDTO := dto.CreateUserDTO{}

	if email, ok := input["email"].(string); ok {
		createDTO.Email = email
	}
	if clerkId, ok := input["clerkId"].(string); ok {
		createDTO.ClerkID = &clerkId
	}
	if firstName, ok := input["firstName"].(string); ok {
		createDTO.FirstName = firstName
	}
	if lastName, ok := input["lastName"].(string); ok {
		createDTO.LastName = lastName
	}
	if phone, ok := input["phone"].(string); ok {
		createDTO.Phone = &phone
	}

	user, err := r.userService.Create(p.Context, createDTO)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Resolver) resolveUpdateUser(p graphql.ResolveParams) (any, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, errors.New("id is required")
	}

	input, ok := p.Args["input"].(map[string]any)
	if !ok {
		return nil, errors.New("input is required")
	}

	updateDTO := dto.UpdateUserDTO{}

	if firstName, ok := input["firstName"].(string); ok {
		updateDTO.FirstName = &firstName
	}
	if lastName, ok := input["lastName"].(string); ok {
		updateDTO.LastName = &lastName
	}
	if phone, ok := input["phone"].(string); ok {
		updateDTO.Phone = &phone
	}
	if isActive, ok := input["isActive"].(bool); ok {
		updateDTO.IsActive = &isActive
	}

	user, err := r.userService.Update(p.Context, id, updateDTO)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Resolver) resolveDeleteUser(p graphql.ResolveParams) (any, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, errors.New("id is required")
	}

	err := r.userService.Delete(p.Context, id)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"success": true,
		"message": "User deleted successfully",
	}, nil
}