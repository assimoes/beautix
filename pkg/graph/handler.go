package graph

import (
	"context"
	"net/http"
	"strings"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/zerolog/log"
)

// Handler handles GraphQL requests
type Handler struct {
	handler       *handler.Handler
	userService   domain.UserService
	schema        *Schema
	resolver      *Resolver
}

// NewHandler creates a new GraphQL handler
func NewHandler(resolver *Resolver) (*Handler, error) {
	schema, err := NewSchema()
	if err != nil {
		return nil, err
	}

	resolver.SetupResolvers(schema)

	// Get the internal graphql.Schema object
	schemaObj := schema.Schema()
	
	h := handler.New(&handler.Config{
		Schema:     &schemaObj,
		Pretty:     true,
		GraphiQL:   false, // We'll use Apollo Sandbox instead
		Playground: false,
	})

	return &Handler{
		handler:  h,
		schema:   schema,
		resolver: resolver,
	}, nil
}

// ServeHTTP handles HTTP requests for GraphQL
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add authentication middleware
	r = r.WithContext(h.authMiddleware(r.Context(), r))
	
	// Handle OPTIONS for CORS
	if r.Method == http.MethodOptions {
		h.corsHandler(w, r)
		return
	}

	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Serve GraphQL
	h.handler.ServeHTTP(w, r)
}

// corsHandler handles CORS preflight requests
func (h *Handler) corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}

// authMiddleware handles authentication via JWT tokens
func (h *Handler) authMiddleware(ctx context.Context, r *http.Request) context.Context {
	// Extract the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ctx
	}

	// Bearer token format
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return ctx
	}

	token := parts[1]
	user, err := h.userService.ValidateToken(ctx, token)
	if err != nil {
		log.Error().Err(err).Str("token", token).Msg("Invalid token")
		return ctx
	}

	// Add user ID to context
	return context.WithValue(ctx, "currentUserID", user.ID)
}

// GraphQLPayload represents the JSON payload for a GraphQL request
type GraphQLPayload struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// ExecuteQuery executes a GraphQL query programmatically
func (h *Handler) ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}) (interface{}, error) {
	schemaObj := h.schema.Schema()
	params := graphql.Params{
		Schema:         schemaObj,
		RequestString:  query,
		VariableValues: variables,
		Context:        ctx,
	}
	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		return nil, result.Errors[0]
	}
	return result.Data, nil
}

// SetUserService sets the user service for token validation
func (h *Handler) SetUserService(userService domain.UserService) {
	h.userService = userService
}