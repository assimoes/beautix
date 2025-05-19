package mock

import (
	"context"
	"errors"
	"time"

	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
)

// ServiceService is a mock implementation of domain.ServiceService
type ServiceService struct {
	services         []*domain.Service
	categories       []*domain.ServiceCategory
	providerService  *ProviderService
}

// NewServiceService creates a new mock service service with some test data
func NewServiceService(providerService *ProviderService) *ServiceService {
	// Create some mock service categories
	categories := []*domain.ServiceCategory{
		{
			ID:          uuid.MustParse("d1e2f3a4-b5c6-4d5e-6f7a-8b9c0d1e2f3a"),
			Name:        "Hair",
			Description: "Hair care services",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          uuid.MustParse("e2f3a4b5-c6d5-4e6f-7a8b-9c0d1e2f3a4b"),
			Name:        "Nails",
			Description: "Nail care services",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          uuid.MustParse("f3a4b5c6-d5e6-4f7a-8b9c-0d1e2f3a4b5c"),
			Name:        "Massage",
			Description: "Massage services",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
	}

	// Create some mock services
	hairCategoryID := categories[0].ID
	nailsCategoryID := categories[1].ID
	massageCategoryID := categories[2].ID

	services := []*domain.Service{
		{
			ID:          uuid.MustParse("a4b5c6d7-e8f9-4a1b-2c3d-4e5f6a7b8c9d"),
			ProviderID:  uuid.MustParse("a1b2c3d4-e5f6-4a5b-8c7d-9e8f7a6b5c4d"), // Beauty Salon A
			CategoryID:  &hairCategoryID,
			Name:        "Haircut",
			Description: "Professional haircut",
			Duration:    30,
			Price:       25.0,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
		},
		{
			ID:          uuid.MustParse("b5c6d7e8-f9a1-4b2c-3d4e-5f6a7b8c9d0e"),
			ProviderID:  uuid.MustParse("a1b2c3d4-e5f6-4a5b-8c7d-9e8f7a6b5c4d"), // Beauty Salon A
			CategoryID:  &hairCategoryID,
			Name:        "Hair Coloring",
			Description: "Professional hair coloring",
			Duration:    120,
			Price:       80.0,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
		},
		{
			ID:          uuid.MustParse("c6d7e8f9-a1b2-4c3d-4e5f-6a7b8c9d0e1f"),
			ProviderID:  uuid.MustParse("b2c3d4e5-f6a5-4b8c-7d9e-8f7a6b5c4d3e"), // Beauty Salon B
			CategoryID:  &nailsCategoryID,
			Name:        "Manicure",
			Description: "Basic manicure",
			Duration:    45,
			Price:       20.0,
			CreatedAt:   time.Now().Add(-6 * time.Hour),
		},
		{
			ID:          uuid.MustParse("d7e8f9a1-b2c3-4d4e-5f6a-7b8c9d0e1f2a"),
			ProviderID:  uuid.MustParse("b2c3d4e5-f6a5-4b8c-7d9e-8f7a6b5c4d3e"), // Beauty Salon B
			CategoryID:  &nailsCategoryID,
			Name:        "Pedicure",
			Description: "Basic pedicure",
			Duration:    60,
			Price:       30.0,
			CreatedAt:   time.Now().Add(-6 * time.Hour),
		},
		{
			ID:          uuid.MustParse("e8f9a1b2-c3d4-4e5f-6a7b-8c9d0e1f2a3b"),
			ProviderID:  uuid.MustParse("a1b2c3d4-e5f6-4a5b-8c7d-9e8f7a6b5c4d"), // Beauty Salon A
			CategoryID:  &massageCategoryID,
			Name:        "Relaxing Massage",
			Description: "Full body relaxing massage",
			Duration:    60,
			Price:       50.0,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
		},
	}

	return &ServiceService{
		services:        services,
		categories:      categories,
		providerService: providerService,
	}
}

// CreateService creates a new service
func (s *ServiceService) CreateService(ctx context.Context, input *domain.CreateServiceInput) (*domain.Service, error) {
	// Check if provider exists
	_, err := s.providerService.GetProvider(ctx, input.ProviderID)
	if err != nil {
		return nil, errors.New("provider not found")
	}

	// Check if category exists if provided
	if input.CategoryID != nil {
		categoryExists := false
		for _, category := range s.categories {
			if category.ID == *input.CategoryID {
				categoryExists = true
				break
			}
		}
		if !categoryExists {
			return nil, errors.New("category not found")
		}
	}

	// Create a new service
	service := &domain.Service{
		ID:          uuid.New(),
		ProviderID:  input.ProviderID,
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		Description: input.Description,
		Duration:    input.Duration,
		Price:       input.Price,
		CreatedAt:   time.Now(),
	}

	// Add the service to the mock database
	s.services = append(s.services, service)

	return service, nil
}

// GetService retrieves a service by ID
func (s *ServiceService) GetService(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	for _, service := range s.services {
		if service.ID == id {
			// Get provider information
			provider, err := s.providerService.GetProvider(ctx, service.ProviderID)
			if err == nil {
				service.Provider = provider
			}

			// Get category information if available
			if service.CategoryID != nil {
				for _, category := range s.categories {
					if category.ID == *service.CategoryID {
						service.Category = category
						break
					}
				}
			}

			return service, nil
		}
	}

	return nil, errors.New("service not found")
}

// UpdateService updates an existing service
func (s *ServiceService) UpdateService(ctx context.Context, id uuid.UUID, input *domain.UpdateServiceInput, updatedBy uuid.UUID) error {
	for i, service := range s.services {
		if service.ID == id {
			if input.CategoryID != nil {
				// Check if category exists
				categoryExists := false
				for _, category := range s.categories {
					if category.ID == *input.CategoryID {
						categoryExists = true
						break
					}
				}
				if !categoryExists {
					return errors.New("category not found")
				}
				service.CategoryID = input.CategoryID
			}

			if input.Name != nil {
				service.Name = *input.Name
			}
			if input.Description != nil {
				service.Description = *input.Description
			}
			if input.Duration != nil {
				service.Duration = *input.Duration
			}
			if input.Price != nil {
				service.Price = *input.Price
			}

			now := time.Now()
			service.UpdatedAt = &now
			service.UpdatedBy = &updatedBy

			s.services[i] = service
			return nil
		}
	}

	return errors.New("service not found")
}

// DeleteService marks a service as deleted
func (s *ServiceService) DeleteService(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	for i, service := range s.services {
		if service.ID == id {
			now := time.Now()
			service.DeletedAt = &now
			service.DeletedBy = &deletedBy

			s.services[i] = service
			return nil
		}
	}

	return errors.New("service not found")
}

// ListServicesByProvider retrieves a list of services for a provider with pagination
func (s *ServiceService) ListServicesByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) ([]*domain.Service, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Filter services by provider ID
	var filtered []*domain.Service
	for _, service := range s.services {
		if service.ProviderID == providerID && service.DeletedAt == nil {
			// Get provider information
			provider, err := s.providerService.GetProvider(ctx, service.ProviderID)
			if err == nil {
				service.Provider = provider
			}

			// Get category information if available
			if service.CategoryID != nil {
				for _, category := range s.categories {
					if category.ID == *service.CategoryID {
						service.Category = category
						break
					}
				}
			}

			filtered = append(filtered, service)
		}
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(filtered) {
		return []*domain.Service{}, nil
	}

	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], nil
}

// ListServicesByCategory retrieves a list of services for a category with pagination
func (s *ServiceService) ListServicesByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int) ([]*domain.Service, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Filter services by category ID
	var filtered []*domain.Service
	for _, service := range s.services {
		if service.CategoryID != nil && *service.CategoryID == categoryID && service.DeletedAt == nil {
			// Get provider information
			provider, err := s.providerService.GetProvider(ctx, service.ProviderID)
			if err == nil {
				service.Provider = provider
			}

			// Get category information
			for _, category := range s.categories {
				if category.ID == categoryID {
					service.Category = category
					break
				}
			}

			filtered = append(filtered, service)
		}
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(filtered) {
		return []*domain.Service{}, nil
	}

	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], nil
}

// CountServices returns the total number of services
func (s *ServiceService) CountServices(ctx context.Context) (int64, error) {
	count := 0
	for _, service := range s.services {
		if service.DeletedAt == nil {
			count++
		}
	}
	return int64(count), nil
}

// CountServicesByProvider returns the total number of services for a provider
func (s *ServiceService) CountServicesByProvider(ctx context.Context, providerID uuid.UUID) (int64, error) {
	count := 0
	for _, service := range s.services {
		if service.ProviderID == providerID && service.DeletedAt == nil {
			count++
		}
	}
	return int64(count), nil
}