package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/assimoes/beautix/configs"
	"github.com/assimoes/beautix/internal/infrastructure/auth"
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"github.com/assimoes/beautix/internal/repository"
	"github.com/assimoes/beautix/internal/service"
	"github.com/assimoes/beautix/pkg/graph"
	"github.com/go-playground/validator/v10"
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

	// Get database instance
	db, err := database.GetDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database instance")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	businessRepo := repository.NewBusinessRepository(db.DB)
	staffRepo := repository.NewStaffRepository(db.DB)

	// Initialize services
	validator := validator.New()
	clerkClient := auth.NewClerkClient()
	authService := service.NewAuthService(userRepo, clerkClient, db.DB)
	userService := service.NewUserService(userRepo, businessRepo, staffRepo, validator)

	// Initialize GraphQL resolver and schema
	resolver := graph.NewResolver(userService, authService)
	schema, err := graph.CreateSchema(resolver)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create GraphQL schema")
	}

	// Setup HTTP server with routes
	mux := http.NewServeMux()

	// GraphQL endpoint
	mux.Handle("/graphql", graph.Handler(schema))

	// GraphQL Sandbox (Apollo Studio)
	mux.Handle("/sandbox", graph.SandboxHandler("http://localhost:8090/graphql"))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","version":"` + Version + `"}`))
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.App.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info().
			Str("port", config.App.Port).
			Str("graphql_endpoint", "/graphql").
			Str("sandbox_endpoint", "/sandbox").
			Str("health_endpoint", "/health").
			Msg("Starting HTTP server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Gracefully shutdown the server with a 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
