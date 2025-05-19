# BeautyBiz GraphQL API

BeautyBiz is a comprehensive provider-centric platform designed specifically for aestheticians in Portugal. This GraphQL API serves as the backend for the BeautyBiz platform, providing functionality for client management, scheduling, business operations, and client acquisition.

## Project Structure

The project follows clean architecture principles:

```
beautix/
  ├── cmd/                   # Application entry points
  │   └── api/               # API server entry point
  ├── configs/               # Configuration
  ├── internal/              # Private application code
  │   ├── domain/            # Business entities and interfaces
  │   ├── service/           # Business logic implementation
  │   ├── repository/        # Data access implementation
  │   └── http/              # HTTP server
  ├── migrations/            # Database migrations
  ├── pkg/                   # Public packages
  │   └── graph/             # GraphQL implementation
  ├── scripts/               # Scripts for development and deployment
  └── test/                  # Test utilities
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL client (optional)

### Setup and Running

1. Clone the repository

```bash
git clone https://github.com/assimoes/beautix.git
cd beautix
```

2. Start the database with Docker Compose

```bash
make docker-up
```

3. Run database migrations

```bash
make migrate-up
```

4. Run the application

```bash
make run
```

The API will be available at http://localhost:8090/graphql, and you can explore it using the Apollo Sandbox at http://localhost:8090/sandbox.

## Development

### Create a Migration

```bash
make migrate-create MIGRATION_NAME=your_migration_name
```

### Run Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
```

### Lint Code

```bash
make lint
```

### Format Code

```bash
make format
```

### Build the Application

```bash
make build
```

## API Documentation

The GraphQL API provides the following main operations:

### Authentication

- `login(email: String!, password: String!): String!` - Authenticates a user and returns a JWT token

### User Management

- `createUser(input: CreateUserInput!): User`
- `updateUser(id: ID!, input: UpdateUserInput!): User`
- `deleteUser(id: ID!): Boolean`
- `user(id: ID!): User`
- `users(limit: Int, offset: Int): [User]`

### Provider Management

- `createProvider(input: CreateProviderInput!): Provider`
- `provider(id: ID!): Provider`
- `providers(limit: Int, offset: Int): [Provider]`
- `searchProviders(query: String!, limit: Int, offset: Int): [Provider]`

### Services

- `service(id: ID!): Service`
- `servicesByProvider(providerId: ID!, limit: Int, offset: Int): [Service]`

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Application
APP_ENV=development
APP_PORT=8090
APP_HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=beautix
DB_SSLMODE=disable
DATABASE_URL=postgres://postgres:postgres@localhost:5432/beautix?sslmode=disable

# Test Database
TEST_DB_NAME=beautix_test
TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/beautix_test?sslmode=disable

# JWT
JWT_SECRET=change_this_to_a_secure_secret_in_production
JWT_EXPIRATION=24h
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.