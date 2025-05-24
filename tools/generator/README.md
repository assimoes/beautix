# BeautiX Code Generator

A code generation tool for quickly scaffolding new entities in the BeautiX project following the Clean Architecture pattern.

## Features

The generator can create:
- **Domain Models**: Entity definitions with GORM tags and validation
- **DTOs**: Data Transfer Objects for Create, Update, and Response
- **GraphQL Resolvers**: Query and mutation resolvers
- **Services**: Business logic layer with base service implementation
- **Repositories**: Data access layer with base repository implementation
- **Unit Tests**: Comprehensive test files for all generated components

## Installation

From the project root:

```bash
cd tools/generator
go build -o beautix-generator
```

Or use the Makefile from the project root:

```bash
make generate
```

## Usage

Run the generator from the `tools/generator` directory:

```bash
./beautix-generator
```

The generator will prompt you for:
1. **Entity name**: The name of your entity (e.g., Product, Service, Client)
2. **Generate domain model?**: Creates the domain entity in `internal/domain/`
3. **Generate DTOs?**: Creates DTOs in `internal/dto/`
4. **Generate service layer?**: Creates service implementation in `internal/service/`
5. **Generate repository layer?**: Creates repository implementation in `internal/repository/`
6. **Generate unit tests?**: Creates test files for all components

## Example

Creating a new "Product" entity:

```
🚀 BeautiX Code Generator
=========================
Enter entity name (e.g., Product, Service, Client): Product
Generate domain model? (y/n): y
Generate DTOs? (y/n): y
Generate service layer? (y/n): y
Generate repository layer? (y/n): y
Generate unit tests? (y/n): y

🔨 Generating files for Product...

📝 Generating domain model...
📝 Generating DTOs...
📝 Generating GraphQL resolver...
📝 Generating service...
📝 Generating repository...
📝 Generating tests...
✅ Code generation completed successfully!
```

## Generated Files

For an entity named "Product", the generator creates:

```
internal/
├── domain/
│   ├── product.go          # Domain model with validation
│   └── product_test.go     # Domain model tests
├── dto/
│   ├── product.go          # DTOs and conversion functions
│   └── product_test.go     # DTO conversion tests
├── service/
│   ├── product_service.go  # Service implementation
│   └── product_service_test.go
├── repository/
│   ├── product_repository.go  # Repository implementation
│   └── product_repository_test.go
pkg/
└── graph/
    ├── product_resolver.go     # GraphQL resolvers
    └── product_resolver_test.go
```

## Post-Generation Steps

After generation, you need to:

1. **Add domain fields** to the generated model
2. **Run migrations** to create the database table
3. **Complete DTO field mappings**
4. **Update GraphQL schema** in `pkg/graph/schema.go`
5. **Register services and repositories** in `cmd/api/main.go`
6. **Run `make generate-mocks`** to generate mocks for new interfaces
7. **Complete TODO comments** in generated files
8. **Run tests** to ensure everything works

## Templates

The generator uses Go templates located in:
- `templates.go`: Service, repository, and resolver templates
- `templates_domain.go`: Domain model and DTO templates

## Customization

To customize the generated code:
1. Edit the template files
2. Rebuild the generator
3. Run with your new templates

## Best Practices

- Use PascalCase for entity names (e.g., `ServiceCategory`, not `service_category`)
- The generator will automatically handle pluralization and case conversion
- Review and complete all TODO comments in generated files
- Always run tests after completing the implementation