package mock

import "github.com/assimoes/beautix/internal/domain"

// ServiceFactory creates mock services for the application
type ServiceFactory struct {
	UserService              domain.UserService
	ProviderService          domain.ProviderService
	ServiceCategoryService   domain.ServiceCategoryService
	ServiceService           domain.ServiceService
	ClientService            domain.ClientService
	AppointmentService       domain.AppointmentService
	ServiceCompletionService domain.ServiceCompletionService
	LoyaltyProgramService    domain.LoyaltyProgramService
	LoyaltyRewardService     domain.LoyaltyRewardService
	ClientLoyaltyService     domain.ClientLoyaltyService
	CampaignService          domain.CampaignService
}

// NewServiceFactory creates a new service factory with mock implementations
func NewServiceFactory() *ServiceFactory {
	// Create the user service
	userService := NewUserService()

	// Create the provider service with the user service
	providerService := NewProviderService(userService)

	// Create the service service with the provider service
	serviceService := NewServiceService(providerService)

	// For now, we don't need to implement all services
	// We just return nil for the ones we haven't implemented yet

	return &ServiceFactory{
		UserService:              userService,
		ProviderService:          providerService,
		ServiceService:           serviceService,
		ServiceCategoryService:   nil,
		ClientService:            nil,
		AppointmentService:       nil,
		ServiceCompletionService: nil,
		LoyaltyProgramService:    nil,
		LoyaltyRewardService:     nil,
		ClientLoyaltyService:     nil,
		CampaignService:          nil,
	}
}