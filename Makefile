.PHONY: build run test test-coverage lint format clean migrate-up migrate-down migrate-create docker-up docker-down docker-logs docker-ps help generate-mocks tidy dev setup air air-install init build-release db-create db-reset db-dump db-restore install-tools all check

# Project variables
PROJECT_NAME := beautix
MAIN_PATH := ./cmd/api
BINARY_NAME := beautix-api
ENV_FILE := .env

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_FLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)"

# Database migration variables
MIGRATIONS_DIR := ./migrations
MIGRATION_NAME ?= migration
DB_NAME := beautix
TEST_DB_NAME := beautix_test
DB_USER := postgres
DB_PASS := postgres
DB_HOST := localhost
DB_PORT := 5432
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
TEST_DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(TEST_DB_NAME)?sslmode=disable

# Docker variables
DOCKER_COMPOSE_FILE := docker-compose.yml

# Air configuration
AIR_CONFIG := .air.toml

# Target: init - Initialize the project (install tools, setup docker, create database, run migrations)
init: install-tools docker-up db-create migrate-up

# Target: install-tools - Install all required development tools
install-tools: 
	@echo "Installing required development tools..."
	@go install github.com/golangci-lint/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/cosmtrek/air@latest
	@echo "All tools installed successfully!"

# Target: setup - Basic setup for the project (create .env file if not exists)
setup:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Creating .env file..."; \
		echo "# Application" > $(ENV_FILE); \
		echo "APP_ENV=development" >> $(ENV_FILE); \
		echo "APP_PORT=8090" >> $(ENV_FILE); \
		echo "APP_HOST=0.0.0.0" >> $(ENV_FILE); \
		echo "" >> $(ENV_FILE); \
		echo "# Database" >> $(ENV_FILE); \
		echo "DB_HOST=localhost" >> $(ENV_FILE); \
		echo "DB_PORT=5432" >> $(ENV_FILE); \
		echo "DB_USER=postgres" >> $(ENV_FILE); \
		echo "DB_PASSWORD=postgres" >> $(ENV_FILE); \
		echo "DB_NAME=beautix" >> $(ENV_FILE); \
		echo "DB_SSLMODE=disable" >> $(ENV_FILE); \
		echo "DATABASE_URL=postgres://postgres:postgres@localhost:5432/beautix?sslmode=disable" >> $(ENV_FILE); \
		echo "" >> $(ENV_FILE); \
		echo "# Test Database" >> $(ENV_FILE); \
		echo "TEST_DB_NAME=beautix_test" >> $(ENV_FILE); \
		echo "TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/beautix_test?sslmode=disable" >> $(ENV_FILE); \
		echo "" >> $(ENV_FILE); \
		echo "# JWT" >> $(ENV_FILE); \
		echo "JWT_SECRET=change_this_to_a_secure_secret_in_production" >> $(ENV_FILE); \
		echo "JWT_EXPIRATION=24h" >> $(ENV_FILE); \
	else \
		echo ".env file already exists."; \
	fi

# Target: all - Run all checks, tests, and build the project
all: format lint test build

# Target: check - Run quick checks before commit (format, lint, build)
check: format lint build

# Target: help - Display available commands
help:
	@echo "🌟 $(PROJECT_NAME) Makefile Help 🌟"
	@echo ""
	@echo "🚀 Main Commands:"
	@echo "  make init            - Initialize the project (install tools, setup docker, create database, run migrations)"
	@echo "  make setup           - Basic setup for the project (create .env file if not exists)"
	@echo "  make dev             - Run the application in development mode with live reload"
	@echo "  make build           - Build the application"
	@echo "  make run             - Run the application locally"
	@echo "  make all             - Run all checks, tests and build the project"
	@echo "  make check           - Run quick checks before commit (format, lint, build)"
	@echo ""
	@echo "🧪 Testing Commands:"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage"
	@echo "  make generate-mocks  - Generate mocks for testing"
	@echo ""
	@echo "🧹 Code Quality Commands:"
	@echo "  make lint            - Run linter"
	@echo "  make format          - Format code"
	@echo "  make tidy            - Clean up go.mod dependencies"
	@echo "  make clean           - Remove build artifacts"
	@echo ""
	@echo "🗄️ Database Commands:"
	@echo "  make migrate-create  - Create a new migration (use MIGRATION_NAME=name)"
	@echo "  make migrate-up      - Run all migrations"
	@echo "  make migrate-down    - Rollback the last migration"
	@echo "  make migrate-down-all - Rollback all migrations (interactive)"
	@echo "  make migrate-force-down-all - Rollback all migrations without confirmation"
	@echo "  make db-create       - Create database if not exists"
	@echo "  make db-reset        - Drop and recreate the database"
	@echo "  make db-dump         - Dump database to file (specify DB_DUMP=file.sql)"
	@echo "  make db-restore      - Restore database from file (specify DB_DUMP=file.sql)"
	@echo ""
	@echo "🔍 Database Inspection Commands:"
	@echo "  make db-list         - List all databases"
	@echo "  make db-tables       - List all tables in the database (use DATABASE=db_name for a different database)"
	@echo "  make db-schemas      - List all schemas in the database (use DATABASE=db_name for a different database)"
	@echo "  make db-extensions   - List all extensions in the database (use DATABASE=db_name for a different database)" 
	@echo "  make db-enums        - List all enum types in the database (use DATABASE=db_name for a different database)"
	@echo "  make db-describe     - Describe a specific table (use TABLE=table_name, DATABASE=db_name for a different database)"
	@echo "  make db-query        - Run a custom SQL query (use QUERY=\"SELECT * FROM table\", DATABASE=db_name for a different database)"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  make docker-up       - Start Docker containers"
	@echo "  make docker-down     - Stop and remove Docker containers"
	@echo "  make docker-logs     - Show Docker container logs (follow mode)"
	@echo "  make docker-ps       - List running Docker containers"
	@echo ""
	@echo "🛠️ Developer Tools:"
	@echo "  make air-install     - Install Air for live reloading"
	@echo "  make install-tools   - Install all required development tools"
	@echo "  make build-release   - Build the application with version information"
	@echo ""
	@echo "For more details, see the Makefile or README.md"

# Target: build - Build the application
build:
	@echo "Building $(PROJECT_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)

# Target: build-release - Build the application with version information
build-release:
	@echo "Building $(PROJECT_NAME) release version $(VERSION)..."
	@go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed: $(BINARY_NAME)"

# Target: run - Run the application
run:
	@echo "Running $(PROJECT_NAME)..."
	@go run $(MAIN_PATH)

# Target: dev - Run the application in development mode with live reload
dev: air-install
	@echo "Running $(PROJECT_NAME) in development mode..."
	@air -c $(AIR_CONFIG)

# Target: air-install - Install Air for live reloading
air-install:
	@if ! command -v air >/dev/null 2>&1; then \
		echo "Installing Air for live reloading..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@if [ ! -f $(AIR_CONFIG) ]; then \
		echo "Creating Air configuration file..."; \
		echo "root = \".\"\n" > $(AIR_CONFIG); \
		echo "tmp_dir = \"tmp\"\n" >> $(AIR_CONFIG); \
		echo "[build]\n" >> $(AIR_CONFIG); \
		echo "  cmd = \"go build -o ./tmp/$(BINARY_NAME) $(MAIN_PATH)\"\n" >> $(AIR_CONFIG); \
		echo "  bin = \"./tmp/$(BINARY_NAME)\"\n" >> $(AIR_CONFIG); \
		echo "  include_ext = [\"go\", \"sql\"]\n" >> $(AIR_CONFIG); \
		echo "  exclude_dir = [\"tmp\", \"vendor\"]\n" >> $(AIR_CONFIG); \
		echo "  delay = 1000\n" >> $(AIR_CONFIG); \
		echo "  kill_delay = 500\n" >> $(AIR_CONFIG); \
		echo "  stop_on_error = true\n" >> $(AIR_CONFIG); \
		echo "\n[log]\n" >> $(AIR_CONFIG); \
		echo "  time = true\n" >> $(AIR_CONFIG); \
		echo "\n[color]\n" >> $(AIR_CONFIG); \
		echo "  main = \"magenta\"\n" >> $(AIR_CONFIG); \
		echo "  watcher = \"cyan\"\n" >> $(AIR_CONFIG); \
		echo "  build = \"yellow\"\n" >> $(AIR_CONFIG); \
		echo "  runner = \"green\"\n" >> $(AIR_CONFIG); \
		echo "\n[misc]\n" >> $(AIR_CONFIG); \
		echo "  clean_on_exit = true\n" >> $(AIR_CONFIG); \
	fi


# Target: test - Run tests
test:
	@echo "Running tests..."
	@go test ./... -count=1 -p 1

# Target: test-coverage - Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./... -count=1 -p 1
	@go tool cover -html=coverage.out

# Target: lint - Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

# Target: format - Format code
format:
	@echo "Formatting code..."
	@go fmt ./...

# Target: clean - Remove build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out
	@rm -rf ./tmp

# Target: docker-up - Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) up -d

# Target: docker-down - Stop and remove Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) down

# Target: docker-logs - Show Docker container logs
docker-logs:
	@echo "Showing Docker logs..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Target: docker-ps - List running Docker containers
docker-ps:
	@echo "Listing Docker containers..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) ps

# Target: db-create - Create database if not exists
db-create: docker-up
	@echo "Creating database if not exists..."
	@docker exec beautix_postgres psql -U $(DB_USER) -tc "SELECT 1 FROM pg_database WHERE datname = '$(DB_NAME)'" | grep -q 1 || \
		docker exec beautix_postgres psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"
	@docker exec beautix_postgres psql -U $(DB_USER) -tc "SELECT 1 FROM pg_database WHERE datname = '$(TEST_DB_NAME)'" | grep -q 1 || \
		docker exec beautix_postgres psql -U $(DB_USER) -c "CREATE DATABASE $(TEST_DB_NAME);"
	@echo "Databases created or already exist."

# Target: db-reset - Drop and recreate the database
db-reset: docker-up
	@echo "Resetting database..."
	@docker exec beautix_postgres psql -U $(DB_USER) -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@docker exec beautix_postgres psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"
	@echo "Database reset completed."

# Target: db-dump - Dump database to file
db-dump: docker-up
	@if [ -z "$(DB_DUMP)" ]; then \
		echo "Error: Please specify DB_DUMP=file.sql"; \
		exit 1; \
	fi
	@echo "Dumping database to $(DB_DUMP)..."
	@docker exec beautix_postgres pg_dump -U $(DB_USER) -d $(DB_NAME) -F p > $(DB_DUMP)
	@echo "Database dump completed."

# Target: db-restore - Restore database from file
db-restore: docker-up
	@if [ -z "$(DB_DUMP)" ]; then \
		echo "Error: Please specify DB_DUMP=file.sql"; \
		exit 1; \
	fi
	@echo "Restoring database from $(DB_DUMP)..."
	@docker exec -i beautix_postgres psql -U $(DB_USER) -d $(DB_NAME) < $(DB_DUMP)
	@echo "Database restore completed."

# Target: db-list - List all databases in the PostgreSQL instance
db-list: docker-up
	@echo "Listing all databases..."
	@docker exec beautix_postgres psql -U $(DB_USER) -c '\l'

# Target: db-tables - List all tables in the specified database (specify DATABASE=db_name to use a different database)
db-tables: docker-up
	@TARGET_DB="$(DB_NAME)"; \
	if [ ! -z "$(DATABASE)" ]; then \
		TARGET_DB="$(DATABASE)"; \
	fi; \
	echo "Listing tables in database $$TARGET_DB..."; \
	docker exec beautix_postgres psql -U $(DB_USER) -d $$TARGET_DB -c '\dt'

# Target: db-schemas - List all schemas in the specified database (specify DATABASE=db_name to use a different database)
db-schemas: docker-up
	@TARGET_DB="$(DB_NAME)"; \
	if [ ! -z "$(DATABASE)" ]; then \
		TARGET_DB="$(DATABASE)"; \
	fi; \
	echo "Listing schemas in database $$TARGET_DB..."; \
	docker exec beautix_postgres psql -U $(DB_USER) -d $$TARGET_DB -c '\dn'

# Target: db-extensions - List all extensions in the specified database (specify DATABASE=db_name to use a different database)
db-extensions: docker-up
	@TARGET_DB="$(DB_NAME)"; \
	if [ ! -z "$(DATABASE)" ]; then \
		TARGET_DB="$(DATABASE)"; \
	fi; \
	echo "Listing extensions in database $$TARGET_DB..."; \
	docker exec beautix_postgres psql -U $(DB_USER) -d $$TARGET_DB -c '\dx'

# Target: db-enums - List all enum types in the database (specify DATABASE=db_name to use a different database)
db-enums: docker-up
	@TARGET_DB="$(DB_NAME)"; \
	if [ ! -z "$(DATABASE)" ]; then \
		TARGET_DB="$(DATABASE)"; \
	fi; \
	echo "Listing enum types in database $$TARGET_DB..."; \
	docker exec beautix_postgres psql -U $(DB_USER) -d $$TARGET_DB -c "SELECT n.nspname AS schema, t.typname AS type, e.enumlabel AS value FROM pg_type t JOIN pg_enum e ON t.oid = e.enumtypid JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace ORDER BY schema, type, e.enumsortorder;"

# Target: db-describe - Describe a specific table (specify TABLE=table_name, DATABASE=db_name to use a different database)
db-describe: docker-up
	@if [ -z "$(TABLE)" ]; then \
		echo "Error: Please specify TABLE=table_name"; \
		exit 1; \
	fi
	@TARGET_DB="$(DB_NAME)"; \
	if [ ! -z "$(DATABASE)" ]; then \
		TARGET_DB="$(DATABASE)"; \
	fi; \
	echo "Describing table $(TABLE) in database $$TARGET_DB..."; \
	docker exec beautix_postgres psql -U $(DB_USER) -d $$TARGET_DB -c "\d $(TABLE)"

# Target: db-query - Run a custom SQL query (specify QUERY="SELECT * FROM users", DATABASE=db_name to use a different database)
db-query: docker-up
	@if [ -z "$(QUERY)" ]; then \
		echo "Error: Please specify QUERY=\"SELECT * FROM users\""; \
		exit 1; \
	fi
	@TARGET_DB="$(DB_NAME)"; \
	if [ ! -z "$(DATABASE)" ]; then \
		TARGET_DB="$(DATABASE)"; \
	fi; \
	echo "Running query in database $$TARGET_DB..."; \
	docker exec beautix_postgres psql -U $(DB_USER) -d $$TARGET_DB -c "$(QUERY)"

# Set the migrate binary path
MIGRATE_BIN := $(shell go env GOPATH)/bin/migrate

# Target: migrate-create - Create a new migration
migrate-create:
	@echo "Creating migration $(MIGRATION_NAME)..."
	@if [ ! -f "$(MIGRATE_BIN)" ]; then \
		echo "Installing migrate command..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(MIGRATION_NAME)

# Target: migrate-up - Run all migrations
migrate-up: docker-up
	@echo "Running migrations up..."
	@if [ ! -f "$(MIGRATE_BIN)" ]; then \
		echo "Installing migrate command..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

# Target: migrate-down - Rollback the last migration
migrate-down: docker-up
	@echo "Running migrations down..."
	@if [ ! -f "$(MIGRATE_BIN)" ]; then \
		echo "Installing migrate command..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

# Target: migrate-down-all - Rollback all migrations (interactive)
migrate-down-all: docker-up
	@echo "Rolling back all migrations..."
	@if [ ! -f "$(MIGRATE_BIN)" ]; then \
		echo "Installing migrate command..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@echo "This will drop all tables in the database. Are you sure? [y/N]"
	@read -r confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down; \
	else \
		echo "Migration rollback cancelled."; \
	fi

# Target: migrate-force-down-all - Rollback all migrations without confirmation (for CI/CD)
migrate-force-down-all: docker-up
	@echo "Forcing rollback of all migrations..."
	@if [ ! -f "$(MIGRATE_BIN)" ]; then \
		echo "Installing migrate command..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force 0

# Target: generate-mocks - Generate mocks for testing
generate-mocks:
	@echo "Generating mocks..."
	@if ! command -v mockery >/dev/null 2>&1; then \
		echo "Installing mockery..."; \
		go install github.com/vektra/mockery/v2@latest; \
	fi
	@mockery --dir=./internal/domain --all --output=./internal/mocks

# Target: tidy - Clean up go.mod
tidy:
	@echo "Tidying go modules..."
	@go mod tidy