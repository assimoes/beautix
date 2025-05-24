package graph

import (
	"github.com/graphql-go/graphql"
	
	"github.com/assimoes/beautix/internal/dto"
)

// UserType represents the GraphQL User type
var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "A user in the system",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The unique identifier of the user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.ID, nil
				}
				return nil, nil
			},
		},
		"email": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The email address of the user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.Email, nil
				}
				return nil, nil
			},
		},
		"clerkId": &graphql.Field{
			Type:        graphql.String,
			Description: "The Clerk ID of the user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.ClerkID, nil
				}
				return nil, nil
			},
		},
		"firstName": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The first name of the user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.FirstName, nil
				}
				return nil, nil
			},
		},
		"lastName": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The last name of the user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.LastName, nil
				}
				return nil, nil
			},
		},
		"fullName": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The full name of the user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.FullName, nil
				}
				return nil, nil
			},
		},
		"phone": &graphql.Field{
			Type:        graphql.String,
			Description: "The phone number of the user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.Phone, nil
				}
				return nil, nil
			},
		},
		"isActive": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Whether the user is active",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.IsActive, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.DateTime),
			Description: "When the user was created",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.DateTime),
			Description: "When the user was last updated",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(*dto.UserResponseDTO); ok {
					return user.UpdatedAt, nil
				}
				return nil, nil
			},
		},
	},
})

// BusinessRoleEnum represents the GraphQL BusinessRole enum for staff roles within businesses
var BusinessRoleEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "BusinessRole",
	Description: "The role of a user within a specific business",
	Values: graphql.EnumValueConfigMap{
		"owner": &graphql.EnumValueConfig{
			Value:       "owner",
			Description: "Business owner with full access",
		},
		"manager": &graphql.EnumValueConfig{
			Value:       "manager",
			Description: "Business manager with administrative access",
		},
		"employee": &graphql.EnumValueConfig{
			Value:       "employee",
			Description: "Business employee with limited access",
		},
		"assistant": &graphql.EnumValueConfig{
			Value:       "assistant",
			Description: "Business assistant with basic access",
		},
	},
})

// BusinessType represents the GraphQL Business type
var BusinessType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Business",
	Description: "A business entity",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The unique identifier of the business",
		},
		"userId": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The ID of the business owner",
		},
		"name": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The name of the business",
		},
		"displayName": &graphql.Field{
			Type:        graphql.String,
			Description: "The display name of the business",
		},
		"businessType": &graphql.Field{
			Type:        graphql.String,
			Description: "The type of business",
		},
		"email": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The email address of the business",
		},
		"website": &graphql.Field{
			Type:        graphql.String,
			Description: "The website URL of the business",
		},
		"logoUrl": &graphql.Field{
			Type:        graphql.String,
			Description: "The logo URL of the business",
		},
		"coverPhotoUrl": &graphql.Field{
			Type:        graphql.String,
			Description: "The cover photo URL of the business",
		},
		"isVerified": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Whether the business is verified",
		},
		"currency": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The currency used by the business",
		},
		"timeZone": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The timezone of the business",
		},
		"isActive": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Whether the business is active",
		},
		"subscriptionTier": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The subscription tier of the business",
		},
		"displayNameValue": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The effective display name (displayName or name)",
		},
		"createdAt": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.DateTime),
			Description: "When the business was created",
		},
		"updatedAt": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.DateTime),
			Description: "When the business was last updated",
		},
	},
})

// CreateUserInput represents the GraphQL input for creating a user
var CreateUserInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "CreateUserInput",
	Description: "Input for creating a new user",
	Fields: graphql.InputObjectConfigFieldMap{
		"email": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The email address of the user",
		},
		"clerkId": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The Clerk ID of the user",
		},
		"firstName": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The first name of the user",
		},
		"lastName": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The last name of the user",
		},
		"phone": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The phone number of the user",
		},
	},
})

// UpdateUserInput represents the GraphQL input for updating a user
var UpdateUserInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "UpdateUserInput",
	Description: "Input for updating an existing user",
	Fields: graphql.InputObjectConfigFieldMap{
		"firstName": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The first name of the user",
		},
		"lastName": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The last name of the user",
		},
		"phone": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The phone number of the user",
		},
		"isActive": &graphql.InputObjectFieldConfig{
			Type:        graphql.Boolean,
			Description: "Whether the user is active",
		},
	},
})

// CreateBusinessInput represents the GraphQL input for creating a business
var CreateBusinessInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "CreateBusinessInput",
	Description: "Input for creating a new business",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The name of the business",
		},
		"displayName": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The display name of the business",
		},
		"businessType": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The type of business",
		},
		"email": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The email address of the business",
		},
		"website": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The website URL of the business",
		},
		"currency": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The currency used by the business",
		},
		"timeZone": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The timezone of the business",
		},
	},
})

// UpdateBusinessInput represents the GraphQL input for updating a business
var UpdateBusinessInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "UpdateBusinessInput",
	Description: "Input for updating an existing business",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The name of the business",
		},
		"displayName": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The display name of the business",
		},
		"businessType": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The type of business",
		},
		"email": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The email address of the business",
		},
		"website": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The website URL of the business",
		},
		"logoUrl": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The logo URL of the business",
		},
		"coverPhotoUrl": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The cover photo URL of the business",
		},
		"currency": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The currency used by the business",
		},
		"timeZone": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The timezone of the business",
		},
		"isActive": &graphql.InputObjectFieldConfig{
			Type:        graphql.Boolean,
			Description: "Whether the business is active",
		},
	},
})