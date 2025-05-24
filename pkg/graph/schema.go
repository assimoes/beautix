package graph

import (
	"github.com/graphql-go/graphql"
)

// CreateSchema creates the GraphQL schema
func CreateSchema(resolver *Resolver) (graphql.Schema, error) {
	// Define the root query type
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// User queries
			"user": &graphql.Field{
				Type:        UserType,
				Description: "Get a user by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The ID of the user to retrieve",
					},
				},
				Resolve: resolver.resolveUser,
			},
			"users": &graphql.Field{
				Type:        graphql.NewList(UserType),
				Description: "Get a list of users with pagination",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						Description:  "Maximum number of users to return",
						DefaultValue: 10,
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						Description:  "Number of users to skip",
						DefaultValue: 0,
					},
				},
				Resolve: resolver.resolveUsers,
			},
			"userByEmail": &graphql.Field{
				Type:        UserType,
				Description: "Get a user by email address",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The email address of the user",
					},
				},
				Resolve: resolver.resolveUserByEmail,
			},
			"searchUsers": &graphql.Field{
				Type:        graphql.NewList(UserType),
				Description: "Search users by name or email",
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The search query",
					},
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						Description:  "Maximum number of users to return",
						DefaultValue: 50,
					},
				},
				Resolve: resolver.resolveSearchUsers,
			},
			"currentUser": &graphql.Field{
				Type:        UserType,
				Description: "Get the currently authenticated user",
				Resolve:     resolver.resolveCurrentUser,
			},
		},
	})

	// Define the root mutation type
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			// User mutations
			"createUser": &graphql.Field{
				Type:        UserType,
				Description: "Create a new user",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(CreateUserInput),
						Description: "The user data",
					},
				},
				Resolve: resolver.resolveCreateUser,
			},
			"updateUser": &graphql.Field{
				Type:        UserType,
				Description: "Update an existing user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The ID of the user to update",
					},
					"input": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(UpdateUserInput),
						Description: "The updated user data",
					},
				},
				Resolve: resolver.resolveUpdateUser,
			},
			"deleteUser": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "DeleteResult",
					Fields: graphql.Fields{
						"success": &graphql.Field{
							Type: graphql.NewNonNull(graphql.Boolean),
						},
						"message": &graphql.Field{
							Type: graphql.String,
						},
					},
				}),
				Description: "Delete a user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The ID of the user to delete",
					},
				},
				Resolve: resolver.resolveDeleteUser,
			},
		},
	})

	// Create the schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	if err != nil {
		return graphql.Schema{}, err
	}

	return schema, nil
}