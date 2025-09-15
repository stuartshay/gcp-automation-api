package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Port           string
	GCPProjectID   string
	GCPCredentials string
	Environment    string
	LogLevel       string
	LogFile        string
	EnableDebug    bool
	GCPRegion      string
	GCPZone        string
	// JWT Configuration
	JWTSecret          string
	JWTExpirationHours int
	GoogleClientID     string
	GoogleClientSecret string
	EnableGoogleAuth   bool
	// OAuth Configuration
	OAuthTokenURL     string
	OAuthRedirectURI  string
	OAuthCallbackPort string
	CredentialsDir    string
	CredentialsFile   string
}

// Load reads configuration from environment variables with defaults
func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		GCPProjectID:   getEnv("GCP_PROJECT_ID", ""),
		GCPCredentials: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
		Environment:    getEnv("ENVIRONMENT", "development"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogFile:        getEnv("LOG_FILE", "logs/app.log"),
		EnableDebug:    getEnvAsBool("ENABLE_DEBUG", false),
		GCPRegion:      getEnv("GCP_REGION", "us-central1"),
		GCPZone:        getEnv("GCP_ZONE", "us-central1-a"),
		// JWT Configuration
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		EnableGoogleAuth:   getEnvAsBool("ENABLE_GOOGLE_AUTH", true),
		// OAuth Configuration
		OAuthTokenURL:     getEnv("OAUTH_TOKEN_URL", "https://oauth2.googleapis.com/token"),
		OAuthRedirectURI:  getEnv("OAUTH_REDIRECT_URI", "http://localhost:8085/callback"),
		OAuthCallbackPort: getEnv("OAUTH_CALLBACK_PORT", "8085"),
		CredentialsDir:    getEnv("CREDENTIALS_DIR", ".gcp-automation"),
		CredentialsFile:   getEnv("CREDENTIALS_FILE", "credentials.json"),
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

// getEnvAsInt gets an environment variable as integer
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
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
