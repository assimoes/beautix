package graph

import (
	"github.com/graphql-go/graphql"
)

// Schema represents the GraphQL schema for the BeautyBiz API
type Schema struct {
	schema graphql.Schema
}

// NewSchema creates a new GraphQL schema
func NewSchema() (*Schema, error) {
	// Define base scalar types
	dateTimeType := graphql.NewScalar(graphql.ScalarConfig{
		Name:        "DateTime",
		Description: "ISO 8601 date-time format",
		Serialize: func(value interface{}) interface{} {
			return value
		},
	})

	// Define GraphQL types
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "A user of the BeautyBiz platform",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Unique identifier",
			},
			"email": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User's email address",
			},
			"firstName": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User's first name",
			},
			"lastName": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User's last name",
			},
			"phone": &graphql.Field{
				Type:        graphql.String,
				Description: "User's phone number",
			},
			"role": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User's role (user, admin, provider)",
			},
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(dateTimeType),
				Description: "When the user was created",
			},
		},
	})

	providerType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Provider",
		Description: "A service provider in the BeautyBiz platform",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Unique identifier",
			},
			"userId": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Associated user ID",
			},
			"businessName": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Provider's business name",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "Provider's business description",
			},
			"address": &graphql.Field{
				Type:        graphql.String,
				Description: "Provider's address",
			},
			"city": &graphql.Field{
				Type:        graphql.String,
				Description: "Provider's city",
			},
			"postalCode": &graphql.Field{
				Type:        graphql.String,
				Description: "Provider's postal code",
			},
			"country": &graphql.Field{
				Type:        graphql.String,
				Description: "Provider's country",
			},
			"website": &graphql.Field{
				Type:        graphql.String,
				Description: "Provider's website URL",
			},
			"logoUrl": &graphql.Field{
				Type:        graphql.String,
				Description: "Provider's logo URL",
			},
			"subscriptionTier": &graphql.Field{
				Type:        graphql.String,
				Description: "Provider's subscription tier",
			},
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(dateTimeType),
				Description: "When the provider was created",
			},
			"user": &graphql.Field{
				Type:        userType,
				Description: "The user associated with this provider",
			},
		},
	})

	serviceCategoryType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "ServiceCategory",
		Description: "A category for beauty services",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Unique identifier",
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Category name",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "Category description",
			},
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(dateTimeType),
				Description: "When the category was created",
			},
		},
	})

	serviceType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Service",
		Description: "A beauty service offered by a provider",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Unique identifier",
			},
			"providerId": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Provider ID",
			},
			"categoryId": &graphql.Field{
				Type:        graphql.ID,
				Description: "Category ID (optional)",
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Service name",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "Service description",
			},
			"duration": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Service duration in minutes",
			},
			"price": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Float),
				Description: "Service price",
			},
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(dateTimeType),
				Description: "When the service was created",
			},
			"provider": &graphql.Field{
				Type:        providerType,
				Description: "The provider offering this service",
			},
			"category": &graphql.Field{
				Type:        serviceCategoryType,
				Description: "The category of this service",
			},
		},
	})

	clientType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Client",
		Description: "A client of a service provider",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Unique identifier",
			},
			"userId": &graphql.Field{
				Type:        graphql.ID,
				Description: "Associated user ID (optional)",
			},
			"providerId": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Provider ID",
			},
			"firstName": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Client's first name",
			},
			"lastName": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Client's last name",
			},
			"email": &graphql.Field{
				Type:        graphql.String,
				Description: "Client's email address",
			},
			"phone": &graphql.Field{
				Type:        graphql.String,
				Description: "Client's phone number",
			},
			"notes": &graphql.Field{
				Type:        graphql.String,
				Description: "Notes about the client",
			},
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(dateTimeType),
				Description: "When the client was created",
			},
			"provider": &graphql.Field{
				Type:        providerType,
				Description: "The provider this client belongs to",
			},
			"user": &graphql.Field{
				Type:        userType,
				Description: "The user associated with this client",
			},
		},
	})

	appointmentType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Appointment",
		Description: "A scheduled appointment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Unique identifier",
			},
			"providerId": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Provider ID",
			},
			"clientId": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Client ID",
			},
			"serviceId": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Service ID",
			},
			"startTime": &graphql.Field{
				Type:        graphql.NewNonNull(dateTimeType),
				Description: "Appointment start time",
			},
			"endTime": &graphql.Field{
				Type:        graphql.NewNonNull(dateTimeType),
				Description: "Appointment end time",
			},
			"status": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Appointment status",
			},
			"notes": &graphql.Field{
				Type:        graphql.String,
				Description: "Notes about the appointment",
			},
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(dateTimeType),
				Description: "When the appointment was created",
			},
			"provider": &graphql.Field{
				Type:        providerType,
				Description: "The provider for this appointment",
			},
			"client": &graphql.Field{
				Type:        clientType,
				Description: "The client for this appointment",
			},
			"service": &graphql.Field{
				Type:        serviceType,
				Description: "The service for this appointment",
			},
		},
	})

	_ = graphql.NewObject(graphql.ObjectConfig{
		Name:        "ServiceCompletion",
		Description: "A record of a completed service",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Unique identifier",
			},
			"appointmentId": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Associated appointment ID",
			},
			"priceCharged": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Float),
				Description: "Price charged for the service",
			},
			"paymentMethod": &graphql.Field{
				Type:        graphql.String,
				Description: "Payment method used",
			},
			"providerConfirmed": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "Whether the provider confirmed completion",
			},
			"clientConfirmed": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "Whether the client confirmed completion",
			},
			"completionDate": &graphql.Field{
				Type:        dateTimeType,
				Description: "When the service was completed",
			},
			"createdAt": &graphql.Field{
				Type:        graphql.NewNonNull(dateTimeType),
				Description: "When the record was created",
			},
			"appointment": &graphql.Field{
				Type:        appointmentType,
				Description: "The associated appointment",
			},
		},
	})

	// Define Input Types
	createUserInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "CreateUserInput",
		Description: "Input for creating a user",
		Fields: graphql.InputObjectConfigFieldMap{
			"email": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User's email address",
			},
			"password": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User's password",
			},
			"firstName": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User's first name",
			},
			"lastName": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User's last name",
			},
			"phone": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "User's phone number",
			},
			"role": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User's role (user, admin, provider)",
			},
		},
	})

	updateUserInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "UpdateUserInput",
		Description: "Input for updating a user",
		Fields: graphql.InputObjectConfigFieldMap{
			"email": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "User's email address",
			},
			"password": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "User's password",
			},
			"firstName": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "User's first name",
			},
			"lastName": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "User's last name",
			},
			"phone": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "User's phone number",
			},
		},
	})

	createProviderInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "CreateProviderInput",
		Description: "Input for creating a provider",
		Fields: graphql.InputObjectConfigFieldMap{
			"userId": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "Associated user ID",
			},
			"businessName": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Provider's business name",
			},
			"description": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Provider's business description",
			},
			"address": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Provider's address",
			},
			"city": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Provider's city",
			},
			"postalCode": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Provider's postal code",
			},
			"country": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Provider's country",
			},
			"website": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Provider's website URL",
			},
			"logoUrl": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Provider's logo URL",
			},
			"subscriptionTier": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Provider's subscription tier",
			},
		},
	})

	// Define root query
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type:        userType,
				Description: "Get a user by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"users": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "List users with pagination",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"provider": &graphql.Field{
				Type:        providerType,
				Description: "Get a provider by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"providers": &graphql.Field{
				Type:        graphql.NewList(providerType),
				Description: "List providers with pagination",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"searchProviders": &graphql.Field{
				Type:        graphql.NewList(providerType),
				Description: "Search providers by query",
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"service": &graphql.Field{
				Type:        serviceType,
				Description: "Get a service by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"servicesByProvider": &graphql.Field{
				Type:        graphql.NewList(serviceType),
				Description: "List services by provider",
				Args: graphql.FieldConfigArgument{
					"providerId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"me": &graphql.Field{
				Type:        userType,
				Description: "Get the currently authenticated user",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
		},
	})

	// Define root mutation
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type:        userType,
				Description: "Create a new user",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(createUserInputType),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"updateUser": &graphql.Field{
				Type:        userType,
				Description: "Update an existing user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(updateUserInputType),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"deleteUser": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete a user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"login": &graphql.Field{
				Type:        graphql.String,
				Description: "Login and get JWT token",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			"createProvider": &graphql.Field{
				Type:        providerType,
				Description: "Create a new provider",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(createProviderInputType),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// This will be replaced by the resolver
					return nil, nil
				},
			},
			// More mutations would be defined here
		},
	})

	// Create schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	if err != nil {
		return nil, err
	}

	return &Schema{schema: schema}, nil
}

// Schema returns the GraphQL schema
func (s *Schema) Schema() graphql.Schema {
	return s.schema
}