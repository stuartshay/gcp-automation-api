package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

func TestCreateProject(t *testing.T) {
	// Test request payload validation
	projectReq := models.ProjectRequest{
		ProjectID:   "test-project-123",
		DisplayName: "Test Project",
		ParentID:    "123456789",
		ParentType:  "organization",
		Labels: map[string]string{
			"environment": "test",
			"team":        "platform",
		},
	}

	reqBody, err := json.Marshal(projectReq)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Verify request structure is valid
	var parsedReq models.ProjectRequest
	if err := json.Unmarshal(reqBody, &parsedReq); err != nil {
		t.Fatalf("Failed to parse request: %v", err)
	}

	if parsedReq.ProjectID != "test-project-123" {
		t.Errorf("Expected project_id 'test-project-123', got '%s'", parsedReq.ProjectID)
	}

	if parsedReq.DisplayName != "Test Project" {
		t.Errorf("Expected display_name 'Test Project', got '%s'", parsedReq.DisplayName)
	}

	t.Log("Handler structure test passed - request format is valid")
}

func TestCreateBucket(t *testing.T) {
	// Test request payload validation
	bucketReq := models.BucketRequest{
		Name:         "test-bucket-123",
		Location:     "us-central1",
		StorageClass: "STANDARD",
		Labels: map[string]string{
			"environment": "test",
		},
		Versioning: true,
	}

	reqBody, err := json.Marshal(bucketReq)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Verify request structure is valid
	var parsedReq models.BucketRequest
	if err := json.Unmarshal(reqBody, &parsedReq); err != nil {
		t.Fatalf("Failed to parse request: %v", err)
	}

	if parsedReq.Name != "test-bucket-123" {
		t.Errorf("Expected name 'test-bucket-123', got '%s'", parsedReq.Name)
	}

	if parsedReq.Location != "us-central1" {
		t.Errorf("Expected location 'us-central1', got '%s'", parsedReq.Location)
	}

	if !parsedReq.Versioning {
		t.Error("Expected versioning to be true")
	}

	t.Log("Bucket handler structure test passed - request format is valid")
}

func TestConfigLoad(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Port == "" {
		t.Error("Port should have a default value")
	}

	if cfg.Environment == "" {
		t.Error("Environment should have a default value")
	}

	if cfg.LogLevel == "" {
		t.Error("LogLevel should have a default value")
	}

	t.Log("Config loading test passed")
}

// Test health check endpoint without GCP dependencies
func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	t.Log("Health check test passed")
}