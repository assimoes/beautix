package testdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/stretchr/testify/require"
)

// GraphQLTestHelper provides utilities for GraphQL integration testing
type GraphQLTestHelper struct {
	Schema graphql.Schema
	t      *testing.T
}

// NewGraphQLTestHelper creates a new GraphQL test helper
func NewGraphQLTestHelper(t *testing.T, schema graphql.Schema) *GraphQLTestHelper {
	return &GraphQLTestHelper{
		Schema: schema,
		t:      t,
	}
}

// ExecuteQuery executes a GraphQL query and returns the result
func (h *GraphQLTestHelper) ExecuteQuery(query string, variables map[string]interface{}) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         h.Schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        context.Background(),
	})
	return result
}

// ExecuteQueryWithContext executes a GraphQL query with a custom context
func (h *GraphQLTestHelper) ExecuteQueryWithContext(ctx context.Context, query string, variables map[string]interface{}) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         h.Schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        ctx,
	})
	return result
}

// MustExecuteQuery executes a query and fails the test if there are errors
func (h *GraphQLTestHelper) MustExecuteQuery(query string, variables map[string]interface{}) map[string]interface{} {
	result := h.ExecuteQuery(query, variables)
	require.Empty(h.t, result.Errors, "GraphQL query returned errors: %v", result.Errors)
	
	data, ok := result.Data.(map[string]interface{})
	require.True(h.t, ok, "Result data is not a map")
	
	return data
}

// ExpectError executes a query expecting an error
func (h *GraphQLTestHelper) ExpectError(query string, variables map[string]interface{}, expectedMessage ...string) []gqlerrors.FormattedError {
	result := h.ExecuteQuery(query, variables)
	require.NotEmpty(h.t, result.Errors, "Expected GraphQL errors but got none")
	
	// If expected message is provided, check for it
	if len(expectedMessage) > 0 {
		found := false
		for _, err := range result.Errors {
			if contains(err.Message, expectedMessage[0]) {
				found = true
				break
			}
		}
		require.True(h.t, found, "No error contains expected message: %s", expectedMessage[0])
	}
	
	return result.Errors
}

// TestContext provides a context with test-specific values
type TestContext struct {
	UserID   string
	TenantID string
	Roles    []string
}

// WithTestContext creates a context with test values
func WithTestContext(ctx context.Context, testCtx TestContext) context.Context {
	if testCtx.UserID != "" {
		ctx = context.WithValue(ctx, "userID", testCtx.UserID)
	}
	if testCtx.TenantID != "" {
		ctx = context.WithValue(ctx, "tenantID", testCtx.TenantID)
	}
	if len(testCtx.Roles) > 0 {
		ctx = context.WithValue(ctx, "roles", testCtx.Roles)
	}
	return ctx
}

// AssertGraphQLResponse provides fluent assertions for GraphQL responses
type AssertGraphQLResponse struct {
	t      *testing.T
	result *graphql.Result
}

// NewAssertGraphQLResponse creates a new response asserter
func NewAssertGraphQLResponse(t *testing.T, result *graphql.Result) *AssertGraphQLResponse {
	return &AssertGraphQLResponse{t: t, result: result}
}

// NoErrors asserts there are no GraphQL errors
func (a *AssertGraphQLResponse) NoErrors() *AssertGraphQLResponse {
	require.Empty(a.t, a.result.Errors, "GraphQL query returned errors: %v", a.result.Errors)
	return a
}

// HasErrors asserts there are GraphQL errors
func (a *AssertGraphQLResponse) HasErrors() *AssertGraphQLResponse {
	require.NotEmpty(a.t, a.result.Errors, "Expected GraphQL errors but got none")
	return a
}

// ErrorContains asserts that an error message contains a substring
func (a *AssertGraphQLResponse) ErrorContains(substring string) *AssertGraphQLResponse {
	require.NotEmpty(a.t, a.result.Errors, "No errors to check")
	found := false
	for _, err := range a.result.Errors {
		if contains(err.Message, substring) {
			found = true
			break
		}
	}
	require.True(a.t, found, "No error contains substring: %s", substring)
	return a
}

// DataEquals asserts the data equals expected value
func (a *AssertGraphQLResponse) DataEquals(expected interface{}) *AssertGraphQLResponse {
	require.Equal(a.t, expected, a.result.Data)
	return a
}

// DataHasKey asserts the data contains a key
func (a *AssertGraphQLResponse) DataHasKey(key string) *AssertGraphQLResponse {
	data, ok := a.result.Data.(map[string]interface{})
	require.True(a.t, ok, "Result data is not a map")
	_, exists := data[key]
	require.True(a.t, exists, "Data does not contain key: %s", key)
	return a
}

// GetData returns the data as a map
func (a *AssertGraphQLResponse) GetData() map[string]interface{} {
	data, ok := a.result.Data.(map[string]interface{})
	require.True(a.t, ok, "Result data is not a map")
	return data
}

// GetField returns a specific field from the data
func (a *AssertGraphQLResponse) GetField(field string) interface{} {
	data := a.GetData()
	value, exists := data[field]
	require.True(a.t, exists, "Field %s does not exist in data", field)
	return value
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && 
		(s == substr || s[:len(substr)] == substr || 
		(len(s) > len(substr) && contains(s[1:], substr)))
}

// WaitForDatabase waits for database operations to complete
func WaitForDatabase(t *testing.T, timeout time.Duration) {
	// Small delay to ensure database operations complete
	time.Sleep(10 * time.Millisecond)
}

// RunSequential ensures tests run sequentially by acquiring the global lock
func RunSequential(t *testing.T, testFunc func()) {
	testMutex.Lock()
	defer testMutex.Unlock()
	testFunc()
}

// QueryBuilder helps build GraphQL queries
type QueryBuilder struct {
	operation string
	name      string
	args      map[string]interface{}
	fields    []string
}

// NewQuery creates a new query builder
func NewQuery(name string) *QueryBuilder {
	return &QueryBuilder{
		operation: "query",
		name:      name,
		args:      make(map[string]interface{}),
		fields:    make([]string, 0),
	}
}

// NewMutation creates a new mutation builder
func NewMutation(name string) *QueryBuilder {
	return &QueryBuilder{
		operation: "mutation",
		name:      name,
		args:      make(map[string]interface{}),
		fields:    make([]string, 0),
	}
}

// WithArg adds an argument
func (qb *QueryBuilder) WithArg(name string, value interface{}) *QueryBuilder {
	qb.args[name] = value
	return qb
}

// WithFields adds fields to select
func (qb *QueryBuilder) WithFields(fields ...string) *QueryBuilder {
	qb.fields = append(qb.fields, fields...)
	return qb
}

// Build constructs the GraphQL query string
func (qb *QueryBuilder) Build() string {
	// This is a simplified builder - in production you'd want proper escaping
	query := qb.operation + " {\n  " + qb.name
	
	// Add arguments
	if len(qb.args) > 0 {
		query += "("
		first := true
		for k, v := range qb.args {
			if !first {
				query += ", "
			}
			first = false
			query += k + ": "
			// Simple type handling - expand as needed
			switch v := v.(type) {
			case string:
				query += `"` + v + `"`
			case int, int64, float64:
				query += fmt.Sprintf("%v", v)
			case bool:
				query += fmt.Sprintf("%v", v)
			default:
				query += fmt.Sprintf("%v", v)
			}
		}
		query += ")"
	}
	
	// Add fields
	if len(qb.fields) > 0 {
		query += " {\n"
		for _, field := range qb.fields {
			query += "    " + field + "\n"
		}
		query += "  }"
	}
	
	query += "\n}"
	return query
}