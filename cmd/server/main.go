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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/handlers"
	authmiddleware "github.com/stuartshay/gcp-automation-api/internal/middleware"
	"github.com/stuartshay/gcp-automation-api/internal/services"
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

	// Debug: print all registered routes
	for _, ri := range router.Routes() {
		log.Printf("Registered route: %s %s -> %s", ri.Method, ri.Path, ri.Handler)
	}

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

func setupRouter(handler *handlers.Handler, authService *services.AuthService, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Custom logging middleware to write to our log file
	if cfg.LogFile != "" {
		logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err == nil {
			gin.DefaultWriter = io.MultiWriter(os.Stdout, logFile)
		}
	}

	// GCP logging client for middleware (inject as Logger interface)
	var logger handlers.Logger
	{
		loggingClient, err := logging.NewClient(context.Background(), cfg.GCPProjectID)
		if err != nil {
			log.Fatalf("Failed to initialize GCP logging client: %v", err)
		}
		cloudLogger := loggingClient.Logger("cloudrun-api")
		logger = &handlers.LoggerAdapter{Logger: cloudLogger}
	}

	// Gin middleware to inject GCP logger as Logger interface
	r.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	})

	// Serve static files from /static directory
	r.Static("/static", "./static")

	// Custom Swagger UI endpoint
	r.GET("/swagger/", func(c *gin.Context) {
		c.File("static/swagger-ui.html")
	})

	// Serve OpenAPI spec for Swagger UI
	r.GET("/swagger/doc.json", func(c *gin.Context) {
		log.Println("DEBUG: /swagger/doc.json handler invoked")
		absPath, err := filepath.Abs("static/doc.json")
		if err != nil {
			log.Printf("ERROR: Failed to resolve doc.json path: %v", err)
			c.String(http.StatusInternalServerError, "Failed to resolve doc.json path")
			return
		}
		log.Printf("DEBUG: Serving file at %s", absPath)
		c.File(absPath)
	})

	// Redirect common swagger URLs to the correct path
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/")
	})
	r.GET("/swagger/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/")
	})

	// Health check endpoint (no authentication required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// NoRoute handler: return 404 for unregistered routes only
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "route not found",
			"code":    http.StatusNotFound,
		})
	})

	// Create authentication middleware
	authMiddleware := authmiddleware.NewAuthMiddleware(cfg)

	// API v1 routes (all require authentication)
	v1 := r.Group("/api/v1")
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

	return r
}
