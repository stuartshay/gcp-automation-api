package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/handlers"
	authmiddleware "github.com/stuartshay/gcp-automation-api/internal/middleware"
	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/internal/services"
	"github.com/stuartshay/gcp-automation-api/tests/integration/mocks"
)

// setupTestServer creates a test Gin server with authentication
func setupTestServer(t *testing.T) (*gin.Engine, *handlers.Handler, *services.AuthService) {
	cfg := &config.Config{
		Port:               "8080",
		Environment:        "test",
		LogLevel:           "debug",
		JWTSecret:          "test-secret-key-for-testing-only",
		JWTExpirationHours: 24,
		EnableGoogleAuth:   false, // Disable for testing
		LogFile:            "logs/test.log",
	}

	// Use mock GCPService for unit tests to avoid requiring real credentials
	mockGCPService := &mocks.MockGCPService{}
	authService := services.NewAuthService(cfg)
	handler := handlers.NewHandler(mockGCPService, authService)

	// Create Gin instance
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r, handler, authService
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
	r, _, authService := setupTestServer(t)

	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-testing-only",
	}
	authMiddleware := authmiddleware.NewAuthMiddleware(cfg)

	// Create a test endpoint
	r.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "access granted"})
	})

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
			expectedError:  "missing authorization header",
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

			r.ServeHTTP(rec, req)

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
	r, _, _ := setupTestServer(t)

	// Health endpoint should not require authentication
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestProtectedEndpointRequiresAuth(t *testing.T) {
	r, _, _ := setupTestServer(t)

	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-testing-only",
	}
	authMiddleware := authmiddleware.NewAuthMiddleware(cfg)

	// Simulate a protected API endpoint
	r.GET("/api/v1/projects/test", authMiddleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "project data"})
	})

	// Test without authorization header
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, strings.ToLower(response.Message), "missing authorization header")
}
