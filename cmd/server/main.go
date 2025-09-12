package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/handlers"
	"github.com/stuartshay/gcp-automation-api/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize services
	gcpService, err := services.NewGCPService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize GCP service: %v", err)
	}

	// Initialize handlers
	handler := handlers.NewHandler(gcpService)

	// Setup router
	router := setupRouter(handler)

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
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

func setupRouter(handler *handlers.Handler) *gin.Engine {
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
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

	return router
}