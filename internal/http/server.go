package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/assimoes/beautix/pkg/graph"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Server represents the HTTP server
type Server struct {
	router     *gin.Engine
	graphQLHandler *graph.Handler
	port       string
	host       string
}

// NewServer creates a new HTTP server
func NewServer(graphQLHandler *graph.Handler, host, port string) *Server {
	// Set Gin mode based on environment
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	return &Server{
		router:     router,
		graphQLHandler: graphQLHandler,
		port:       port,
		host:       host,
	}
}

// Setup configures the HTTP routes
func (s *Server) Setup() {
	// Add CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// Add routes
	s.router.GET("/health", s.healthHandler)
	s.router.POST("/graphql", s.graphQLRouteHandler)
	s.router.GET("/graphql", s.graphQLRouteHandler)
	s.router.GET("/sandbox", s.sandboxHandler)

	// 404 handler
	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})
}

// Run starts the HTTP server
func (s *Server) Run() error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.host, s.port),
		Handler: s.router,
	}

	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Info().Msgf("Starting server on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking until a signal is received or server fails
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-shutdown:
		log.Info().Msg("Server is shutting down...")

		// Give outstanding requests a deadline to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			// Force close if graceful shutdown fails
			if err := server.Close(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %w", err)
			}
		}
	}

	return nil
}

// healthHandler handles health check requests
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// graphQLRouteHandler handles GraphQL requests
func (s *Server) graphQLRouteHandler(c *gin.Context) {
	s.graphQLHandler.ServeHTTP(c.Writer, c.Request)
}

// sandboxHandler serves the Apollo Sandbox
func (s *Server) sandboxHandler(c *gin.Context) {
	handler := graph.SandboxHandler()
	handler.ServeHTTP(c.Writer, c.Request)
}