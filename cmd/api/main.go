package main

import (
	"fmt"
	"os"

	"github.com/assimoes/beautix/configs"
	"github.com/assimoes/beautix/internal/http"
	"github.com/assimoes/beautix/internal/mock"
	"github.com/assimoes/beautix/pkg/graph"
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

	// Create mock services and resolver
	serviceFactory := mock.NewServiceFactory()
	resolver := graph.NewResolver(
		serviceFactory.UserService,
		serviceFactory.ProviderService,
		serviceFactory.ServiceCategoryService,
		serviceFactory.ServiceService,
		serviceFactory.ClientService,
		serviceFactory.AppointmentService,
		serviceFactory.ServiceCompletionService,
		serviceFactory.LoyaltyProgramService,
		serviceFactory.LoyaltyRewardService,
		serviceFactory.ClientLoyaltyService,
		serviceFactory.CampaignService,
	)

	// Create GraphQL handler
	graphQLHandler, err := graph.NewHandler(resolver)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create GraphQL handler")
	}
	
	// Set the user service for authentication
	graphQLHandler.SetUserService(serviceFactory.UserService)

	// Create and setup HTTP server
	server := http.NewServer(graphQLHandler, config.App.Host, config.App.Port)
	server.Setup()

	// Run the server
	if err := server.Run(); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}

