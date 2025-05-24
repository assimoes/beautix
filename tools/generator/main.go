package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Config holds the generation configuration
type Config struct {
	EntityName       string
	EntityNameLower  string
	EntityNamePlural string
	GenerateDomain   bool
	GenerateDTO      bool
	GenerateService  bool
	GenerateRepo     bool
	GenerateTests    bool
}

func main() {
	fmt.Println("🚀 BeautiX Code Generator")
	fmt.Println("=========================")

	config := &Config{}
	
	// Get entity name
	config.EntityName = promptString("Enter entity name (e.g., Product, Service, Client): ")
	if config.EntityName == "" {
		fmt.Println("❌ Entity name is required")
		os.Exit(1)
	}
	
	config.EntityNameLower = strings.ToLower(config.EntityName)
	config.EntityNamePlural = config.EntityName + "s" // Simple pluralization

	// Ask about domain model generation
	config.GenerateDomain = promptBool("Generate domain model? (y/n): ")

	// Ask about DTO generation
	config.GenerateDTO = promptBool("Generate DTOs? (y/n): ")

	// Ask about service generation
	config.GenerateService = promptBool("Generate service layer? (y/n): ")
	
	// Ask about repository generation (only if service is being generated)
	if config.GenerateService {
		config.GenerateRepo = promptBool("Generate repository layer? (y/n): ")
	}

	// Ask about test generation
	config.GenerateTests = promptBool("Generate unit tests? (y/n): ")

	// Generate files
	fmt.Printf("\n🔨 Generating files for %s...\n\n", config.EntityName)

	if err := generateFiles(config); err != nil {
		fmt.Printf("❌ Error generating files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Code generation completed successfully!")
	fmt.Println("\n📝 Next steps:")
	
	step := 1
	
	if config.GenerateDomain {
		fmt.Printf("%d. Add domain fields to the generated model in internal/domain/%s.go\n", step, config.EntityNameLower)
		step++
		fmt.Printf("%d. Run migrations to create the %ss table\n", step, config.EntityNameLower)
		step++
	}
	
	if config.GenerateDTO {
		fmt.Printf("%d. Complete the DTO field mappings in internal/dto/%s.go\n", step, config.EntityNameLower)
		step++
	}
	
	fmt.Printf("%d. Update the GraphQL schema in pkg/graph/schema.go to register the new resolvers\n", step)
	step++
	
	if config.GenerateService {
		fmt.Printf("%d. Add the new service to the dependency injection in cmd/api/main.go\n", step)
		step++
		fmt.Printf("%d. Implement any custom service methods you added\n", step)
		step++
	}
	
	if config.GenerateRepo {
		fmt.Printf("%d. Add the new repository to the dependency injection in cmd/api/main.go\n", step)
		step++
		fmt.Printf("%d. Implement any custom repository methods you added\n", step)
		step++
	}
	
	fmt.Printf("%d. Complete the TODO comments in the generated files\n", step)
	step++
	
	if config.GenerateTests {
		fmt.Printf("%d. Run 'make generate-mocks' to generate mocks for the new interfaces\n", step)
		step++
		fmt.Printf("%d. Run the tests to ensure everything works correctly\n", step)
		step++
	}
	
	fmt.Printf("%d. Update the GraphQL resolver implementations with your business logic\n", step)
}

func promptString(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptBool(prompt string) bool {
	for {
		fmt.Print(prompt)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		
		if input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
		fmt.Println("Please enter 'y' or 'n'")
	}
}

func generateFiles(config *Config) error {
	// Generate domain model if requested
	if config.GenerateDomain {
		if err := generateDomain(config); err != nil {
			return fmt.Errorf("generating domain model: %w", err)
		}
	}

	// Generate DTOs if requested
	if config.GenerateDTO {
		if err := generateDTO(config); err != nil {
			return fmt.Errorf("generating DTOs: %w", err)
		}
	}

	// Generate GraphQL resolver
	if err := generateResolver(config); err != nil {
		return fmt.Errorf("generating resolver: %w", err)
	}

	// Generate service if requested
	if config.GenerateService {
		if err := generateService(config); err != nil {
			return fmt.Errorf("generating service: %w", err)
		}
	}

	// Generate repository if requested
	if config.GenerateRepo {
		if err := generateRepository(config); err != nil {
			return fmt.Errorf("generating repository: %w", err)
		}
	}

	// Generate tests if requested
	if config.GenerateTests {
		if err := generateTests(config); err != nil {
			return fmt.Errorf("generating tests: %w", err)
		}
	}

	return nil
}

func generateResolver(config *Config) error {
	fmt.Printf("📝 Generating GraphQL resolver...\n")
	
	tmpl, err := template.New("resolver").Parse(resolverTemplate)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("../../pkg/graph/%s_resolver.go", config.EntityNameLower)
	return writeTemplate(tmpl, config, filename)
}

func generateService(config *Config) error {
	fmt.Printf("📝 Generating service...\n")
	
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("../../internal/service/%s_service.go", config.EntityNameLower)
	return writeTemplate(tmpl, config, filename)
}

func generateRepository(config *Config) error {
	fmt.Printf("📝 Generating repository...\n")
	
	tmpl, err := template.New("repository").Parse(repositoryTemplate)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("../../internal/repository/%s_repository.go", config.EntityNameLower)
	return writeTemplate(tmpl, config, filename)
}

func generateDomain(config *Config) error {
	fmt.Printf("📝 Generating domain model...\n")
	
	tmpl, err := template.New("domain").Parse(domainTemplate)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("../../internal/domain/%s.go", config.EntityNameLower)
	return writeTemplate(tmpl, config, filename)
}

func generateDTO(config *Config) error {
	fmt.Printf("📝 Generating DTOs...\n")
	
	tmpl, err := template.New("dto").Parse(dtoTemplate)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("../../internal/dto/%s.go", config.EntityNameLower)
	return writeTemplate(tmpl, config, filename)
}

func generateTests(config *Config) error {
	fmt.Printf("📝 Generating tests...\n")
	
	// Generate domain test if domain is being generated
	if config.GenerateDomain {
		tmpl, err := template.New("domain_test").Parse(domainTestTemplate)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("../../internal/domain/%s_test.go", config.EntityNameLower)
		if err := writeTemplate(tmpl, config, filename); err != nil {
			return err
		}
	}

	// Generate DTO test if DTO is being generated
	if config.GenerateDTO {
		tmpl, err := template.New("dto_test").Parse(dtoTestTemplate)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("../../internal/dto/%s_test.go", config.EntityNameLower)
		if err := writeTemplate(tmpl, config, filename); err != nil {
			return err
		}
	}
	
	// Generate resolver test
	tmpl, err := template.New("resolver_test").Parse(resolverTestTemplate)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("../../pkg/graph/%s_resolver_test.go", config.EntityNameLower)
	if err := writeTemplate(tmpl, config, filename); err != nil {
		return err
	}

	// Generate service test if service is being generated
	if config.GenerateService {
		tmpl, err := template.New("service_test").Parse(serviceTestTemplate)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("../../internal/service/%s_service_test.go", config.EntityNameLower)
		if err := writeTemplate(tmpl, config, filename); err != nil {
			return err
		}
	}

	// Generate repository test if repository is being generated
	if config.GenerateRepo {
		tmpl, err := template.New("repository_test").Parse(repositoryTestTemplate)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("../../internal/repository/%s_repository_test.go", config.EntityNameLower)
		if err := writeTemplate(tmpl, config, filename); err != nil {
			return err
		}
	}

	return nil
}

func writeTemplate(tmpl *template.Template, config *Config, filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute template
	return tmpl.Execute(file, config)
}