package auth

import (
	"context"
	"fmt"

	"github.com/assimoes/beautix/internal/dto"
)

// ClerkClient is a placeholder for Clerk integration
// TODO: Implement proper Clerk integration in next task
type ClerkClient struct {
	// placeholder fields
}

// NewClerkClient creates a new Clerk client (placeholder)
func NewClerkClient() *ClerkClient {
	return &ClerkClient{}
}

// VerifyToken verifies a Clerk JWT token (placeholder)
func (c *ClerkClient) VerifyToken(ctx context.Context, token string) (*dto.ClerkUserDTO, error) {
	// TODO: Implement actual token verification
	return nil, fmt.Errorf("clerk integration not implemented yet")
}

// GetUser retrieves a user from Clerk (placeholder)
func (c *ClerkClient) GetUser(ctx context.Context, userID string) (*dto.ClerkUserDTO, error) {
	// TODO: Implement actual user retrieval
	return nil, fmt.Errorf("clerk integration not implemented yet")
}

// CreateUser creates a user in Clerk (placeholder)
func (c *ClerkClient) CreateUser(ctx context.Context, userData dto.ClerkUserDTO) (*dto.ClerkUserDTO, error) {
	// TODO: Implement actual user creation
	return nil, fmt.Errorf("clerk integration not implemented yet")
}

// UpdateUser updates a user in Clerk (placeholder)
func (c *ClerkClient) UpdateUser(ctx context.Context, userID string, userData dto.ClerkUserDTO) (*dto.ClerkUserDTO, error) {
	// TODO: Implement actual user update
	return nil, fmt.Errorf("clerk integration not implemented yet")
}

// DeleteUser deletes a user from Clerk (placeholder)
func (c *ClerkClient) DeleteUser(ctx context.Context, userID string) error {
	// TODO: Implement actual user deletion
	return fmt.Errorf("clerk integration not implemented yet")
}