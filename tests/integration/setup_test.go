package integration

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/handlers"
	"github.com/stuartshay/gcp-automation-api/internal/services"
	"github.com/stuartshay/gcp-automation-api/tests/integration/mocks"
)

// TestConfig holds configuration for integration tests
type TestConfig struct {
	UseRealGCP    bool
	TestProjectID string
	BucketPrefix  string
}

// TestSetup holds the test setup components
type TestSetup struct {
	Router      *gin.Engine
	Handler     *handlers.Handler
	MockService *mocks.MockGCPService
	Config      *TestConfig
	AuthService *services.AuthService
}

// SetupTestServer creates a test server with either mock or real GCP service
func SetupTestServer(t *testing.T) *TestSetup {
	// Load test configuration
	testConfig := loadTestConfig()

	// Create Gin instance
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Load application config
	cfg := &config.Config{
		Port:               "8080",
		Environment:        "test",
		LogLevel:           "debug",
		JWTSecret:          "test-secret-key-for-testing-only",
		JWTExpirationHours: 24,
		EnableGoogleAuth:   false,
		LogFile:            "logs/test.log",
		GCPProjectID:       testConfig.TestProjectID,
	}

	// Initialize auth service
	authService := services.NewAuthService(cfg)

	var gcpService services.GCPServiceInterface
	var mockService *mocks.MockGCPService

	if testConfig.UseRealGCP {
		// Use real GCP service for integration tests
		realService, err := services.NewGCPService(cfg)
		if err != nil {
			t.Fatalf("Failed to initialize real GCP service: %v", err)
		}
		gcpService = realService
	} else {
		// Use mock service for unit tests
		mockService = &mocks.MockGCPService{}
		gcpService = mockService
	}

	// Initialize handlers
	handler := handlers.NewHandler(gcpService, authService)

	return &TestSetup{
		Router:      r,
		Handler:     handler,
		MockService: mockService,
		Config:      testConfig,
		AuthService: authService,
	}
}

// loadTestConfig loads test configuration from environment variables
func loadTestConfig() *TestConfig {
	useRealGCP := os.Getenv("TEST_MODE") == "integration"
	testProjectID := os.Getenv("TEST_PROJECT_ID")
	bucketPrefix := os.Getenv("TEST_BUCKET_PREFIX")

	if bucketPrefix == "" {
		bucketPrefix = "test-gcp-automation"
	}

	if useRealGCP && testProjectID == "" {
		// If using real GCP but no project ID is set, fall back to mock
		useRealGCP = false
	}

	return &TestConfig{
		UseRealGCP:    useRealGCP,
		TestProjectID: testProjectID,
		BucketPrefix:  bucketPrefix,
	}
}

// GenerateTestJWT creates a valid JWT token for testing
func GenerateTestJWT(t *testing.T, authService *services.AuthService) string {
	token, err := authService.GenerateTestJWT("test-user-123", "test@example.com", "Test User")
	assert.NoError(t, err)
	return token
}

// AssertSuccessResponseWithData validates a successful API response that should contain data
func AssertSuccessResponseWithData(t *testing.T, body []byte, expectedMessage string) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedMessage, response["message"])
	assert.Contains(t, response, "data")

	return response
}

// AssertSuccessResponse validates a successful API response (may or may not have data)
func AssertSuccessResponse(t *testing.T, body []byte, expectedMessage string) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedMessage, response["message"])

	return response
}

// AssertErrorResponse validates an error API response
func AssertErrorResponse(t *testing.T, body []byte, expectedCode int, expectedError string) {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Handle different error response formats
	if code, exists := response["code"]; exists {
		// Our API error format
		assert.Equal(t, float64(expectedCode), code)
		assert.Contains(t, response["error"], expectedError)
	} else if message, exists := response["message"]; exists {
		// Echo's default error format
		assert.Contains(t, message, expectedError)
	} else {
		// Fallback - just check the body contains the error
		assert.Contains(t, string(body), expectedError)
	}
}

// CleanupTestResources cleans up any test resources (for integration tests)
func CleanupTestResources(t *testing.T, setup *TestSetup) {
	if setup.Config.UseRealGCP {
		// Add cleanup logic for real GCP resources if needed
		t.Log("Cleaning up test resources...")
	}
}

// resetMockExpectations resets the mock expectations and calls for the test setup
func resetMockExpectations(setup *TestSetup) {
	if !setup.Config.UseRealGCP && setup.MockService != nil {
		// Use the proper testify method to reset the mock
		setup.MockService.Mock.ExpectedCalls = nil
		setup.MockService.Mock.Calls = nil
	}
}

// GetTestBucketName generates a unique test bucket name
func GetTestBucketName(prefix string) string {
	return prefix + "-" + generateRandomString(8)
}

// GetTestProjectID generates a unique test project ID
func GetTestProjectID(prefix string) string {
	return prefix + "-" + generateRandomString(8)
}

// GetTestFolderName generates a unique test folder name
func GetTestFolderName(prefix string) string {
	return prefix + "-" + generateRandomString(8)
}

// generateRandomString generates a cryptographically secure random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		// Use crypto/rand for secure random generation
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to timestamp-based generation if crypto/rand fails
			seed := time.Now().UnixNano() + int64(i)
			b[i] = charset[seed%int64(len(charset))]
		} else {
			b[i] = charset[n.Int64()]
		}
	}
	return string(b)
}

// TestEnvironment represents the test environment type
type TestEnvironment string

const (
	// MockEnvironment uses mocked services
	MockEnvironment TestEnvironment = "mock"
	// IntegrationEnvironment uses real GCP services
	IntegrationEnvironment TestEnvironment = "integration"
)

// GetTestEnvironment returns the current test environment
func GetTestEnvironment() TestEnvironment {
	if os.Getenv("TEST_MODE") == "integration" {
		return IntegrationEnvironment
	}
	return MockEnvironment
}

// SkipIfNotIntegration skips the test if not running in integration mode
func SkipIfNotIntegration(t *testing.T) {
	if GetTestEnvironment() != IntegrationEnvironment {
		t.Skip("Skipping integration test - set TEST_MODE=integration to run")
	}
}

// SkipIfNotMock skips the test if not running in mock mode
func SkipIfNotMock(t *testing.T) {
	if GetTestEnvironment() != MockEnvironment {
		t.Skip("Skipping mock test - unset TEST_MODE or set TEST_MODE=mock to run")
	}
}
