// Package main GCP Automation API
//
// This is a GCP Automation API built with Go that provides RESTful endpoints
// for automating Google Cloud Platform resource management.
//
// The service supports creating, retrieving, and managing:
// - GCP Projects
// - GCP Folders
// - Cloud Storage Buckets
//
// @title GCP Automation API
// @version 1.0
// @description RESTful API for automating Google Cloud Platform resource management
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api/v1
//
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Authorization: Bearer {token}"
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/stuartshay/gcp-automation-api/docs" // Import generated docs
	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/handlers"
	authmiddleware "github.com/stuartshay/gcp-automation-api/internal/middleware"
	"github.com/stuartshay/gcp-automation-api/internal/services"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
)

// setupLogging configures logging to write to both file and console
func setupLogging(cfg *config.Config) error {
	// Create logs directory if it doesn't exist
	logDir := filepath.Dir(cfg.LogFile)
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Set up multi-writer to write to both file and console
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// Set log format with timestamp
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Logging configured - writing to: %s", cfg.LogFile)
	return nil
}

// createDynamicSwaggerHandler creates a custom Swagger handler that serves swagger.json file with examples
func createDynamicSwaggerHandler(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if this is a request for the swagger.json
		if c.Request().URL.Path == "/swagger/doc.json" {
			// Read the swagger.json file directly (contains x-examples)
			swaggerFile := "docs/swagger.json"
			swaggerData, err := os.ReadFile(swaggerFile)
			if err != nil {
				log.Printf("Failed to read swagger.json file: %v", err)
				// Fall back to embedded docs
				doc := swag.GetSwagger("swagger")
				if doc != nil {
					var swaggerSpec map[string]interface{}
					if err := json.Unmarshal([]byte(doc.ReadDoc()), &swaggerSpec); err != nil {
						log.Printf("Failed to unmarshal swagger JSON: %v", err)
					} else {
						// Update host and schemes
						swaggerSpec["host"] = cfg.SwaggerHost
						swaggerSpec["schemes"] = []string{cfg.SwaggerScheme}
						return c.JSON(http.StatusOK, swaggerSpec)
					}
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load swagger documentation")
			}

			// Parse the file content and update host/schemes
			var swaggerSpec map[string]interface{}
			if err := json.Unmarshal(swaggerData, &swaggerSpec); err != nil {
				log.Printf("Failed to unmarshal swagger.json file: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Invalid swagger.json format")
			}

			// Update host and schemes to match configuration
			swaggerSpec["host"] = cfg.SwaggerHost
			swaggerSpec["schemes"] = []string{cfg.SwaggerScheme}

			// Return the modified JSON with x-examples intact
			return c.JSON(http.StatusOK, swaggerSpec)
		}

		// For all other swagger requests, use the default handler
		return echoSwagger.WrapHandler(c)
	}
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logging
	if err := setupLogging(cfg); err != nil {
		log.Fatalf("Failed to setup logging: %v", err)
	}

	// Initialize services
	gcpService, err := services.NewGCPService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize GCP service: %v", err)
	}

	// Initialize authentication service
	authService := services.NewAuthService(cfg)

	// Initialize handlers
	handler := handlers.NewHandler(gcpService, authService)

	// Setup router
	router := setupRouter(handler, authService, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRouter(handler *handlers.Handler, authService *services.AuthService, cfg *config.Config) *echo.Echo {
	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Custom logging middleware to write to our log file
	if cfg.LogFile != "" {
		logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err == nil {
			e.Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
		}
	}

	// Swagger endpoint with dynamic configuration
	e.GET("/swagger/*", createDynamicSwaggerHandler(cfg))

	// Health check endpoint (no authentication required)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	// Create authentication middleware
	authMiddleware := authmiddleware.NewAuthMiddleware(cfg)

	// API v1 routes (all require authentication)
	v1 := e.Group("/api/v1")
	v1.Use(authMiddleware.RequireAuth())
	{
		// Project endpoints
		projects := v1.Group("/projects")
		{
			projects.POST("", handler.CreateProject)
			projects.GET("/:id", handler.GetProject)
			projects.DELETE("/:id", handler.DeleteProject)
		}

		// Folder endpoints
		folders := v1.Group("/folders")
		{
			folders.POST("", handler.CreateFolder)
			folders.GET("/:id", handler.GetFolder)
			folders.DELETE("/:id", handler.DeleteFolder)
		}

		// Bucket endpoints
		buckets := v1.Group("/buckets")
		{
			buckets.POST("", handler.CreateBucket)
			buckets.GET("/:name", handler.GetBucket)
			buckets.DELETE("/:name", handler.DeleteBucket)
		}
	}

	return e
}
