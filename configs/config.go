package configs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
type Config struct {
	App         AppConfig
	Database    DatabaseConfig
	Auth        AuthConfig
	Environment string
}

// AppConfig stores application configuration
type AppConfig struct {
	Host string
	Port string
}

// DatabaseConfig stores database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

// AuthConfig stores authentication configuration
type AuthConfig struct {
	Secret     string
	Expiration time.Duration
}

// LoadConfig reads configuration from environment variables or .env file
func LoadConfig() (*Config, error) {
	// Set default values
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_HOST", "0.0.0.0")
	viper.SetDefault("APP_PORT", "8090")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "beautix")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("JWT_SECRET", "change_this_to_a_secure_secret_in_production")
	viper.SetDefault("JWT_EXPIRATION", "24h")

	// Set environment variable prefix
	viper.SetEnvPrefix("")

	// Read from environment variables
	viper.AutomaticEnv()

	// Create config
	config := &Config{
		Environment: viper.GetString("APP_ENV"),
		App: AppConfig{
			Host: viper.GetString("APP_HOST"),
			Port: viper.GetString("APP_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Auth: AuthConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
	}

	// Set database URL
	config.Database.URL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	// Parse JWT expiration duration
	expiration, err := time.ParseDuration(viper.GetString("JWT_EXPIRATION"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION: %w", err)
	}
	config.Auth.Expiration = expiration

	return config, nil
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsTesting returns true if the application is running in testing mode
func (c *Config) IsTesting() bool {
	return c.Environment == "testing"
}
