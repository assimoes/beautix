package graph

import (
	"github.com/assimoes/beautix/internal/domain"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

// Resolver handles GraphQL resolvers
type Resolver struct {
	userService              domain.UserService
	providerService          domain.ProviderService
	serviceCategoryService   domain.ServiceCategoryService
	serviceService           domain.ServiceService
	clientService            domain.ClientService
	appointmentService       domain.AppointmentService
	serviceCompletionService domain.ServiceCompletionService
	loyaltyProgramService    domain.LoyaltyProgramService
	loyaltyRewardService     domain.LoyaltyRewardService
	clientLoyaltyService     domain.ClientLoyaltyService
	campaignService          domain.CampaignService
}

// NewResolver creates a new resolver with all required services
func NewResolver(
	userService domain.UserService,
	providerService domain.ProviderService,
	serviceCategoryService domain.ServiceCategoryService,
	serviceService domain.ServiceService,
	clientService domain.ClientService,
	appointmentService domain.AppointmentService,
	serviceCompletionService domain.ServiceCompletionService,
	loyaltyProgramService domain.LoyaltyProgramService,
	loyaltyRewardService domain.LoyaltyRewardService,
	clientLoyaltyService domain.ClientLoyaltyService,
	campaignService domain.CampaignService,
) *Resolver {
	return &Resolver{
		userService:              userService,
		providerService:          providerService,
		serviceCategoryService:   serviceCategoryService,
		serviceService:           serviceService,
		clientService:            clientService,
		appointmentService:       appointmentService,
		serviceCompletionService: serviceCompletionService,
		loyaltyProgramService:    loyaltyProgramService,
		loyaltyRewardService:     loyaltyRewardService,
		clientLoyaltyService:     clientLoyaltyService,
		campaignService:          campaignService,
	}
}

// SetupResolvers sets up all resolver functions for the GraphQL schema
func (r *Resolver) SetupResolvers(schema *Schema) {
	// Get the internal schema object
	schemaObj := schema.Schema()
	
	// Get query and mutation fields
	queryType := schemaObj.QueryType()
	if queryType != nil {
		queryFields := queryType.Fields()
		r.setupQueryResolvers(queryFields)
	}
	
	mutationType := schemaObj.MutationType()
	if mutationType != nil {
		mutationFields := mutationType.Fields()
		r.setupMutationResolvers(mutationFields)
	}
}

// setupQueryResolvers configures resolvers for all query fields
func (r *Resolver) setupQueryResolvers(queryFields graphql.FieldDefinitionMap) {
	// User queries
	if field, ok := queryFields["user"]; ok {
		field.Resolve = r.resolveUser
	}
	if field, ok := queryFields["users"]; ok {
		field.Resolve = r.resolveUsers
	}

	// Provider queries
	if field, ok := queryFields["provider"]; ok {
		field.Resolve = r.resolveProvider
	}
	if field, ok := queryFields["providers"]; ok {
		field.Resolve = r.resolveProviders
	}
	if field, ok := queryFields["searchProviders"]; ok {
		field.Resolve = r.resolveSearchProviders
	}

	// Service queries
	if field, ok := queryFields["service"]; ok {
		field.Resolve = r.resolveService
	}
	if field, ok := queryFields["servicesByProvider"]; ok {
		field.Resolve = r.resolveServicesByProvider
	}

	// Current user
	if field, ok := queryFields["me"]; ok {
		field.Resolve = r.resolveMe
	}
}

// setupMutationResolvers configures resolvers for all mutation fields
func (r *Resolver) setupMutationResolvers(mutationFields graphql.FieldDefinitionMap) {
	// User mutations
	if field, ok := mutationFields["createUser"]; ok {
		field.Resolve = r.resolveCreateUser
	}
	if field, ok := mutationFields["updateUser"]; ok {
		field.Resolve = r.resolveUpdateUser
	}
	if field, ok := mutationFields["deleteUser"]; ok {
		field.Resolve = r.resolveDeleteUser
	}
	if field, ok := mutationFields["login"]; ok {
		field.Resolve = r.resolveLogin
	}

	// Provider mutations
	if field, ok := mutationFields["createProvider"]; ok {
		field.Resolve = r.resolveCreateProvider
	}
}

// User resolvers
func (r *Resolver) resolveUser(p graphql.ResolveParams) (interface{}, error) {
	id, err := uuid.Parse(p.Args["id"].(string))
	if err != nil {
		return nil, err
	}

	return r.userService.GetUser(p.Context, id)
}

func (r *Resolver) resolveUsers(p graphql.ResolveParams) (interface{}, error) {
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)
	page := (offset / limit) + 1

	return r.userService.ListUsers(p.Context, page, limit)
}

func (r *Resolver) resolveMe(p graphql.ResolveParams) (interface{}, error) {
	// Extract user from context (would be set by auth middleware)
	userID, ok := p.Context.Value("currentUserID").(uuid.UUID)
	if !ok {
		return nil, nil // Not authenticated
	}

	return r.userService.GetUser(p.Context, userID)
}

func (r *Resolver) resolveCreateUser(p graphql.ResolveParams) (interface{}, error) {
	input, _ := p.Args["input"].(map[string]interface{})
	createInput := &domain.CreateUserInput{
		Email:     input["email"].(string),
		Password:  input["password"].(string),
		FirstName: input["firstName"].(string),
		LastName:  input["lastName"].(string),
		Role:      input["role"].(string),
	}

	if phone, ok := input["phone"]; ok && phone != nil {
		createInput.Phone = phone.(string)
	}

	return r.userService.CreateUser(p.Context, createInput)
}

func (r *Resolver) resolveUpdateUser(p graphql.ResolveParams) (interface{}, error) {
	id, err := uuid.Parse(p.Args["id"].(string))
	if err != nil {
		return nil, err
	}

	input, _ := p.Args["input"].(map[string]interface{})
	updateInput := &domain.UpdateUserInput{}

	if email, ok := input["email"]; ok && email != nil {
		emailStr := email.(string)
		updateInput.Email = &emailStr
	}
	if password, ok := input["password"]; ok && password != nil {
		passwordStr := password.(string)
		updateInput.Password = &passwordStr
	}
	if firstName, ok := input["firstName"]; ok && firstName != nil {
		firstNameStr := firstName.(string)
		updateInput.FirstName = &firstNameStr
	}
	if lastName, ok := input["lastName"]; ok && lastName != nil {
		lastNameStr := lastName.(string)
		updateInput.LastName = &lastNameStr
	}
	if phone, ok := input["phone"]; ok && phone != nil {
		phoneStr := phone.(string)
		updateInput.Phone = &phoneStr
	}

	// Get current user ID from context for the updatedBy field
	updatedBy, _ := p.Context.Value("currentUserID").(uuid.UUID)

	if err := r.userService.UpdateUser(p.Context, id, updateInput, updatedBy); err != nil {
		return nil, err
	}

	return r.userService.GetUser(p.Context, id)
}

func (r *Resolver) resolveDeleteUser(p graphql.ResolveParams) (interface{}, error) {
	id, err := uuid.Parse(p.Args["id"].(string))
	if err != nil {
		return nil, err
	}

	// Get current user ID from context for the deletedBy field
	deletedBy, _ := p.Context.Value("currentUserID").(uuid.UUID)

	return true, r.userService.DeleteUser(p.Context, id, deletedBy)
}

func (r *Resolver) resolveLogin(p graphql.ResolveParams) (interface{}, error) {
	email := p.Args["email"].(string)
	password := p.Args["password"].(string)

	user, err := r.userService.Authenticate(p.Context, email, password)
	if err != nil {
		return nil, err
	}

	return r.userService.GenerateToken(p.Context, user)
}

// Provider resolvers
func (r *Resolver) resolveProvider(p graphql.ResolveParams) (interface{}, error) {
	id, err := uuid.Parse(p.Args["id"].(string))
	if err != nil {
		return nil, err
	}

	return r.providerService.GetProvider(p.Context, id)
}

func (r *Resolver) resolveProviders(p graphql.ResolveParams) (interface{}, error) {
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)
	page := (offset / limit) + 1

	return r.providerService.ListProviders(p.Context, page, limit)
}

func (r *Resolver) resolveSearchProviders(p graphql.ResolveParams) (interface{}, error) {
	query := p.Args["query"].(string)
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)
	page := (offset / limit) + 1

	return r.providerService.SearchProviders(p.Context, query, page, limit)
}

func (r *Resolver) resolveCreateProvider(p graphql.ResolveParams) (interface{}, error) {
	input, _ := p.Args["input"].(map[string]interface{})
	
	userID, err := uuid.Parse(input["userId"].(string))
	if err != nil {
		return nil, err
	}

	createInput := &domain.CreateProviderInput{
		UserID:       userID,
		BusinessName: input["businessName"].(string),
	}

	if description, ok := input["description"]; ok && description != nil {
		createInput.Description = description.(string)
	}
	if address, ok := input["address"]; ok && address != nil {
		createInput.Address = address.(string)
	}
	if city, ok := input["city"]; ok && city != nil {
		createInput.City = city.(string)
	}
	if postalCode, ok := input["postalCode"]; ok && postalCode != nil {
		createInput.PostalCode = postalCode.(string)
	}
	if country, ok := input["country"]; ok && country != nil {
		createInput.Country = country.(string)
	} else {
		createInput.Country = "Portugal" // Default country
	}
	if website, ok := input["website"]; ok && website != nil {
		createInput.Website = website.(string)
	}
	if logoUrl, ok := input["logoUrl"]; ok && logoUrl != nil {
		createInput.LogoURL = logoUrl.(string)
	}
	if subscriptionTier, ok := input["subscriptionTier"]; ok && subscriptionTier != nil {
		createInput.SubscriptionTier = subscriptionTier.(string)
	} else {
		createInput.SubscriptionTier = "free" // Default tier
	}

	return r.providerService.CreateProvider(p.Context, createInput)
}

// Service resolvers
func (r *Resolver) resolveService(p graphql.ResolveParams) (interface{}, error) {
	id, err := uuid.Parse(p.Args["id"].(string))
	if err != nil {
		return nil, err
	}

	return r.serviceService.GetService(p.Context, id)
}

func (r *Resolver) resolveServicesByProvider(p graphql.ResolveParams) (interface{}, error) {
	providerID, err := uuid.Parse(p.Args["providerId"].(string))
	if err != nil {
		return nil, err
	}

	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)
	page := (offset / limit) + 1

	return r.serviceService.ListServicesByProvider(p.Context, providerID, page, limit)
}