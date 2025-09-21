package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stuartshay/gcp-automation-api/internal/middleware"
	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/tests/integration/mocks"
)

func TestProjectOperations(t *testing.T) {
	setup := SetupTestServer(t)
	defer CleanupTestResources(t, setup)

	// Setup authentication middleware
	ginAuthMiddleware := middleware.NewAuthMiddleware(setup.AuthService.GetConfig())
	r := setup.Router
	v1 := r.Group("/api/v1")
	v1.Use(ginAuthMiddleware.RequireAuth())
	{
		projects := v1.Group("/projects")
		{
			projects.POST("", setup.Handler.CreateProject)
			projects.GET(":id", setup.Handler.GetProject)
			projects.DELETE(":id", setup.Handler.DeleteProject)
		}
	}

	// Generate test JWT token
	token := GenerateTestJWT(t, setup.AuthService)

	t.Run("CreateProject", func(t *testing.T) {
		tests := []struct {
			name           string
			request        models.ProjectRequest
			mockSetup      func(*mocks.MockGCPService)
			expectedStatus int
			expectedError  string
		}{
			// ...copy test cases from original testCreateProject...
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resetMockExpectations(setup)
				if tt.mockSetup != nil {
					tt.mockSetup(setup.MockService)
				}
				reqBody, err := json.Marshal(tt.request)
				assert.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				if tt.expectedError != "" {
					AssertErrorResponse(t, rec.Body.Bytes(), tt.expectedStatus, tt.expectedError)
				} else {
					response := AssertSuccessResponseWithData(t, rec.Body.Bytes(), "Project created successfully")
					data, ok := response["data"].(map[string]interface{})
					assert.True(t, ok)
					assert.Equal(t, tt.request.ProjectID, data["project_id"])
					assert.Equal(t, tt.request.DisplayName, data["display_name"])
					if tt.request.ParentID != "" {
						assert.Equal(t, tt.request.ParentID, data["parent_id"])
						assert.Equal(t, tt.request.ParentType, data["parent_type"])
					}
				}
				resetMockExpectations(setup)
			})
		}
	})

	t.Run("GetProject", func(t *testing.T) {
		// ...implement Gin-based get project tests...
	})

	t.Run("DeleteProject", func(t *testing.T) {
		// ...implement Gin-based delete project tests...
	})

	t.Run("ProjectValidation", func(t *testing.T) {
		// ...implement Gin-based project validation tests...
	})
}

func TestProjectAuthenticationRequired(t *testing.T) {
	setup := SetupTestServer(t)
	defer CleanupTestResources(t, setup)

	// Setup authentication middleware
	ginAuthMiddleware := middleware.NewAuthMiddleware(setup.AuthService.GetConfig())
	r := setup.Router
	v1 := r.Group("/api/v1")
	v1.Use(ginAuthMiddleware.RequireAuth())
	{
		projects := v1.Group("/projects")
		{
			projects.POST("", setup.Handler.CreateProject)
			projects.GET(":id", setup.Handler.GetProject)
			projects.DELETE(":id", setup.Handler.DeleteProject)
		}
	}

	tests := []struct {
		name   string
		method string
		url    string
		body   string
	}{
		{
			name:   "Create project without auth",
			method: http.MethodPost,
			url:    "/api/v1/projects",
			body:   `{"project_id":"test","display_name":"Test"}`,
		},
		{
			name:   "Get project without auth",
			method: http.MethodGet,
			url:    "/api/v1/projects/test",
		},
		{
			name:   "Delete project without auth",
			method: http.MethodDelete,
			url:    "/api/v1/projects/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.url, bytes.NewBufferString(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.url, nil)
			}

			rec := httptest.NewRecorder()

			// Execute request without authorization header
			setup.Router.ServeHTTP(rec, req)

			// Should return 401 Unauthorized
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			AssertErrorResponse(t, rec.Body.Bytes(), http.StatusUnauthorized, "unauthorized")
		})
	}
}
