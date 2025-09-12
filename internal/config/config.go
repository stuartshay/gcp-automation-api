package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Port             string
	GCPProjectID     string
	GCPCredentials   string
	Environment      string
	LogLevel         string
	EnableDebug      bool
	GCPRegion        string
	GCPZone          string
}

// Load reads configuration from environment variables with defaults
func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		GCPProjectID:   getEnv("GCP_PROJECT_ID", ""),
		GCPCredentials: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
		Environment:    getEnv("ENVIRONMENT", "development"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		EnableDebug:    getEnvAsBool("ENABLE_DEBUG", false),
		GCPRegion:      getEnv("GCP_REGION", "us-central1"),
		GCPZone:        getEnv("GCP_ZONE", "us-central1-a"),
	}

	return cfg, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if running in development environment  
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}