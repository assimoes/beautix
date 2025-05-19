package mock

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
)

// ProviderService is a mock implementation of domain.ProviderService
type ProviderService struct {
	providers []*domain.Provider
	userService *UserService
}

// NewProviderService creates a new mock provider service with some test data
func NewProviderService(userService *UserService) *ProviderService {
	// Create some mock providers for testing
	providers := []*domain.Provider{
		{
			ID:               uuid.MustParse("a1b2c3d4-e5f6-4a5b-8c7d-9e8f7a6b5c4d"),
			UserID:           uuid.MustParse("5f6e8b3a-2d98-4d80-8f4e-6c79f236a3d2"), // provider@example.com
			BusinessName:     "Beauty Salon A",
			Description:      "A premium beauty salon offering a wide range of services",
			Address:          "Rua das Flores, 123",
			City:             "Lisboa",
			PostalCode:       "1000-123",
			Country:          "Portugal",
			Website:          "https://beautysalona.pt",
			LogoURL:          "https://example.com/logo.png",
			SubscriptionTier: "premium",
			CreatedAt:        time.Now().Add(-12 * time.Hour),
		},
		{
			ID:               uuid.MustParse("b2c3d4e5-f6a5-4b8c-7d9e-8f7a6b5c4d3e"),
			UserID:           uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479"), // admin@example.com
			BusinessName:     "Beauty Salon B",
			Description:      "A cozy beauty salon specializing in nail care",
			Address:          "Avenida da Liberdade, 456",
			City:             "Porto",
			PostalCode:       "4000-456",
			Country:          "Portugal",
			Website:          "https://beautysalonb.pt",
			LogoURL:          "https://example.com/logo2.png",
			SubscriptionTier: "basic",
			CreatedAt:        time.Now().Add(-6 * time.Hour),
		},
	}

	return &ProviderService{
		providers:   providers,
		userService: userService,
	}
}

// CreateProvider creates a new provider
func (s *ProviderService) CreateProvider(ctx context.Context, input *domain.CreateProviderInput) (*domain.Provider, error) {
	// Check if user exists
	user, err := s.userService.GetUser(ctx, input.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Create a new provider
	provider := &domain.Provider{
		ID:               uuid.New(),
		UserID:           input.UserID,
		BusinessName:     input.BusinessName,
		Description:      input.Description,
		Address:          input.Address,
		City:             input.City,
		PostalCode:       input.PostalCode,
		Country:          input.Country,
		Website:          input.Website,
		LogoURL:          input.LogoURL,
		SubscriptionTier: input.SubscriptionTier,
		CreatedAt:        time.Now(),
		User:             user,
	}

	// Add the provider to the mock database
	s.providers = append(s.providers, provider)

	return provider, nil
}

// GetProvider retrieves a provider by ID
func (s *ProviderService) GetProvider(ctx context.Context, id uuid.UUID) (*domain.Provider, error) {
	for _, provider := range s.providers {
		if provider.ID == id {
			// Get user information
			user, err := s.userService.GetUser(ctx, provider.UserID)
			if err == nil {
				provider.User = user
			}
			return provider, nil
		}
	}

	return nil, errors.New("provider not found")
}

// GetProviderByUserID retrieves a provider by user ID
func (s *ProviderService) GetProviderByUserID(ctx context.Context, userID uuid.UUID) (*domain.Provider, error) {
	for _, provider := range s.providers {
		if provider.UserID == userID {
			// Get user information
			user, err := s.userService.GetUser(ctx, provider.UserID)
			if err == nil {
				provider.User = user
			}
			return provider, nil
		}
	}

	return nil, errors.New("provider not found")
}

// UpdateProvider updates an existing provider
func (s *ProviderService) UpdateProvider(ctx context.Context, id uuid.UUID, input *domain.UpdateProviderInput, updatedBy uuid.UUID) error {
	for i, provider := range s.providers {
		if provider.ID == id {
			if input.BusinessName != nil {
				provider.BusinessName = *input.BusinessName
			}
			if input.Description != nil {
				provider.Description = *input.Description
			}
			if input.Address != nil {
				provider.Address = *input.Address
			}
			if input.City != nil {
				provider.City = *input.City
			}
			if input.PostalCode != nil {
				provider.PostalCode = *input.PostalCode
			}
			if input.Country != nil {
				provider.Country = *input.Country
			}
			if input.Website != nil {
				provider.Website = *input.Website
			}
			if input.LogoURL != nil {
				provider.LogoURL = *input.LogoURL
			}
			if input.SubscriptionTier != nil {
				provider.SubscriptionTier = *input.SubscriptionTier
			}

			now := time.Now()
			provider.UpdatedAt = &now
			provider.UpdatedBy = &updatedBy

			s.providers[i] = provider
			return nil
		}
	}

	return errors.New("provider not found")
}

// DeleteProvider marks a provider as deleted
func (s *ProviderService) DeleteProvider(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	for i, provider := range s.providers {
		if provider.ID == id {
			now := time.Now()
			provider.DeletedAt = &now
			provider.DeletedBy = &deletedBy

			s.providers[i] = provider
			return nil
		}
	}

	return errors.New("provider not found")
}

// ListProviders retrieves a list of providers with pagination
func (s *ProviderService) ListProviders(ctx context.Context, page, pageSize int) ([]*domain.Provider, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(s.providers) {
		return []*domain.Provider{}, nil
	}

	if end > len(s.providers) {
		end = len(s.providers)
	}

	// Get user information for each provider
	result := s.providers[start:end]
	for i, provider := range result {
		user, err := s.userService.GetUser(ctx, provider.UserID)
		if err == nil {
			result[i].User = user
		}
	}

	return result, nil
}

// SearchProviders searches providers by query
func (s *ProviderService) SearchProviders(ctx context.Context, query string, page, pageSize int) ([]*domain.Provider, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Filter providers by query
	var filtered []*domain.Provider
	for _, provider := range s.providers {
		if strings.Contains(strings.ToLower(provider.BusinessName), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(provider.Description), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(provider.City), strings.ToLower(query)) {
			filtered = append(filtered, provider)
		}
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(filtered) {
		return []*domain.Provider{}, nil
	}

	if end > len(filtered) {
		end = len(filtered)
	}

	// Get user information for each provider
	result := filtered[start:end]
	for i, provider := range result {
		user, err := s.userService.GetUser(ctx, provider.UserID)
		if err == nil {
			result[i].User = user
		}
	}

	return result, nil
}

// CountProviders returns the total number of providers
func (s *ProviderService) CountProviders(ctx context.Context) (int64, error) {
	return int64(len(s.providers)), nil
}