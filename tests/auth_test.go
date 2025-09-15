package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/handlers"
	authmiddleware "github.com/stuartshay/gcp-automation-api/internal/middleware"
	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/internal/services"
)

// setupTestServer creates a test Echo server with authentication
func setupTestServer(t *testing.T) (*echo.Echo, *handlers.Handler, *services.AuthService) {
	cfg := &config.Config{
		Port:               "8080",
		Environment:        "test",
		LogLevel:           "debug",
		JWTSecret:          "test-secret-key-for-testing-only",
		JWTExpirationHours: 24,
		EnableGoogleAuth:   false, // Disable for testing
		LogFile:            "logs/test.log",
	}

	// Initialize services
	gcpService, err := services.NewGCPService(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize GCP service: %v", err)
	}

	authService := services.NewAuthService(cfg)
	handler := handlers.NewHandler(gcpService, authService)

	// Create Echo instance
	e := echo.New()

	return e, handler, authService
}

// generateTestJWT creates a valid JWT token for testing
func generateTestJWT(t *testing.T, authService *services.AuthService) string {
	token, err := authService.GenerateTestJWT("test-user-123", "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to generate test JWT: %v", err)
	}
	return token
}

func TestJWTMiddleware(t *testing.T) {
	e, _, authService := setupTestServer(t)

	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-testing-only",
	}
	authMiddleware := authmiddleware.NewAuthMiddleware(cfg)

	// Create a test endpoint
	e.GET("/protected", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "access granted"})
	}, authMiddleware.RequireAuth())

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "No Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid or missing jwt token",
		},
		{
			name:           "Invalid token format",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid or missing jwt token",
		},
		{
			name:           "Valid token",
			authHeader:     "Bearer " + generateTestJWT(t, authService),
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, strings.ToLower(response.Message), strings.ToLower(tt.expectedError))
			}
		})
	}
}

func TestHealthEndpointNoAuth(t *testing.T) {
	e, _, _ := setupTestServer(t)

	// Health endpoint should not require authentication
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestProtectedEndpointRequiresAuth(t *testing.T) {
	e, _, _ := setupTestServer(t)

	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-testing-only",
	}
	authMiddleware := authmiddleware.NewAuthMiddleware(cfg)

	// Simulate a protected API endpoint
	e.GET("/api/v1/projects/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "project data"})
	}, authMiddleware.RequireAuth())

	// Test without authorization header
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/test", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, strings.ToLower(response.Message), "invalid or missing jwt")
}
