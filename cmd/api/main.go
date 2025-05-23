package main

import (
	"os"

	"github.com/assimoes/beautix/configs"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Version information - will be set during build
var (
	Version    = "dev"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)

func main() {
	// Configure logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file, using environment variables")
	}

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set log level based on environment
	logLevel := zerolog.InfoLevel
	if config.IsDevelopment() {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	log.Info().
		Str("version", Version).
		Str("build_time", BuildTime).
		Str("commit", CommitHash).
		Str("environment", config.Environment).
		Msg("Starting BeautyBiz API")

	// Initialize database
	if err := database.InitDB(config); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer database.CloseDB()

	// For now, we're focusing on database setup and will implement HTTP server later
	log.Info().Msg("Database setup complete. HTTP server implementation pending.")

	// Keep the application running until interrupted
	select {}
}
