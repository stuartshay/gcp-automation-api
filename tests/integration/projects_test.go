package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stuartshay/gcp-automation-api/internal/middleware"
	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/tests/integration/mocks"
)

func TestProjectOperations(t *testing.T) {
	setup := SetupTestServer(t)
	defer CleanupTestResources(t, setup)

	// Setup authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(setup.AuthService.GetConfig())

	// Setup routes
	setup.Echo.POST("/api/v1/projects", setup.Handler.CreateProject, authMiddleware.RequireAuth())
	setup.Echo.GET("/api/v1/projects/:id", setup.Handler.GetProject, authMiddleware.RequireAuth())
	setup.Echo.DELETE("/api/v1/projects/:id", setup.Handler.DeleteProject, authMiddleware.RequireAuth())

	// Generate test JWT token
	token := GenerateTestJWT(t, setup.AuthService)

	t.Run("CreateProject", func(t *testing.T) {
		testCreateProject(t, setup, token)
	})

	t.Run("GetProject", func(t *testing.T) {
		testGetProject(t, setup, token)
	})

	t.Run("DeleteProject", func(t *testing.T) {
		testDeleteProject(t, setup, token)
	})

	t.Run("ProjectValidation", func(t *testing.T) {
		testProjectValidation(t, setup, token)
	})
}

func testCreateProject(t *testing.T, setup *TestSetup, token string) {
	tests := []struct {
		name           string
		request        models.ProjectRequest
		mockSetup      func(*mocks.MockGCPService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid project creation",
			request: models.ProjectRequest{
				ProjectID:   "test-project-123",
				DisplayName: "Test Project",
				ParentID:    "123456789",
				ParentType:  "organization",
				Labels: map[string]string{
					"environment": "test",
					"team":        "platform",
				},
			},
			mockSetup: func(m *mocks.MockGCPService) {
				if setup.Config.UseRealGCP {
					return // Skip mock setup for real GCP tests
				}
				expectedResponse := mocks.NewMockProjectResponse(&models.ProjectRequest{
					ProjectID:   "test-project-123",
					DisplayName: "Test Project",
					ParentID:    "123456789",
					ParentType:  "organization",
					Labels: map[string]string{
						"environment": "test",
						"team":        "platform",
					},
				})
				m.On("CreateProject", mock.AnythingOfType("*models.ProjectRequest")).Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Minimal project creation",
			request: models.ProjectRequest{
				ProjectID:   "minimal-project-456",
				DisplayName: "Minimal Project",
			},
			mockSetup: func(m *mocks.MockGCPService) {
				if setup.Config.UseRealGCP {
					return
				}
				expectedResponse := mocks.NewMockProjectResponse(&models.ProjectRequest{
					ProjectID:   "minimal-project-456",
					DisplayName: "Minimal Project",
				})
				m.On("CreateProject", mock.AnythingOfType("*models.ProjectRequest")).Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Project with folder parent",
			request: models.ProjectRequest{
				ProjectID:   "folder-project-789",
				DisplayName: "Folder Parent Project",
				ParentID:    "987654321",
				ParentType:  "folder",
				Labels: map[string]string{
					"environment": "development",
				},
			},
			mockSetup: func(m *mocks.MockGCPService) {
				if setup.Config.UseRealGCP {
					return
				}
				expectedResponse := mocks.NewMockProjectResponse(&models.ProjectRequest{
					ProjectID:   "folder-project-789",
					DisplayName: "Folder Parent Project",
					ParentID:    "987654321",
					ParentType:  "folder",
					Labels: map[string]string{
						"environment": "development",
					},
				})
				m.On("CreateProject", mock.AnythingOfType("*models.ProjectRequest")).Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "GCP service error",
			request: models.ProjectRequest{
				ProjectID:   "error-project",
				DisplayName: "Error Project",
			},
			mockSetup: func(m *mocks.MockGCPService) {
				if setup.Config.UseRealGCP {
					return
				}
				m.On("CreateProject", mock.AnythingOfType("*models.ProjectRequest")).Return(nil, fmt.Errorf("GCP service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock if not using real GCP
			if !setup.Config.UseRealGCP && setup.MockService != nil {
				tt.mockSetup(setup.MockService)
			}

			// Create request
			reqBody, err := json.Marshal(tt.request)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			// Execute request
			setup.Echo.ServeHTTP(rec, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				AssertErrorResponse(t, rec.Body.Bytes(), tt.expectedStatus, tt.expectedError)
			} else {
				response := AssertSuccessResponseWithData(t, rec.Body.Bytes(), "Project created successfully")

				// Validate response data
				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, tt.request.ProjectID, data["project_id"])
				assert.Equal(t, tt.request.DisplayName, data["display_name"])

				if tt.request.ParentID != "" {
					assert.Equal(t, tt.request.ParentID, data["parent_id"])
					assert.Equal(t, tt.request.ParentType, data["parent_type"])
				}
			}

			// Reset mock expectations
			if !setup.Config.UseRealGCP && setup.MockService != nil {
				setup.MockService.ExpectedCalls = nil
				setup.MockService.Calls = nil
			}
		})
	}
}

func testGetProject(t *testing.T, setup *TestSetup, token string) {
	tests := []struct {
		name           string
		projectID      string
		mockSetup      func(*mocks.MockGCPService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "Get existing project",
			projectID: "test-project-123",
			mockSetup: func(m *mocks.MockGCPService) {
				if setup.Config.UseRealGCP {
					return
				}
				expectedResponse := &models.ProjectResponse{
					ProjectID:     "test-project-123",
					DisplayName:   "Test Project",
					State:         "ACTIVE",
					ProjectNumber: 123456789,
				}
				m.On("GetProject", "test-project-123").Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Get non-existent project",
			projectID: "non-existent-project",
			mockSetup: func(m *mocks.MockGCPService) {
				if setup.Config.UseRealGCP {
					return
				}
				m.On("GetProject", "non-existent-project").Return(nil, fmt.Errorf("project not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Project not found",
		},
		{
			name:      "Empty project ID",
			projectID: "",
			mockSetup: func(m *mocks.MockGCPService) {
				// No mock setup needed as validation happens before service call
			},
			expectedStatus: http.StatusNotFound, // Echo returns 404 for empty path params
			expectedError:  "Not Found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock if not using real GCP
			if !setup.Config.UseRealGCP && setup.MockService != nil {
				tt.mockSetup(setup.MockService)
			}

			// Create request
			url := fmt.Sprintf("/api/v1/projects/%s", tt.projectID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			// Execute request
			setup.Echo.ServeHTTP(rec, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				AssertErrorResponse(t, rec.Body.Bytes(), tt.expectedStatus, tt.expectedError)
			} else {
				response := AssertSuccessResponse(t, rec.Body.Bytes(), "Project retrieved successfully")

				// Validate response data
				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, tt.projectID, data["project_id"])
			}

			// Reset mock expectations
			if !setup.Config.UseRealGCP && setup.MockService != nil {
				setup.MockService.ExpectedCalls = nil
				setup.MockService.Calls = nil
			}
		})
	}
}

func testDeleteProject(t *testing.T, setup *TestSetup, token string) {
	tests := []struct {
		name           string
		projectID      string
		mockSetup      func(*mocks.MockGCPService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "Delete existing project",
			projectID: "test-project-123",
			mockSetup: func(m *mocks.MockGCPService) {
				if setup.Config.UseRealGCP {
					return
				}
				m.On("DeleteProject", "test-project-123").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Delete non-existent project",
			projectID: "non-existent-project",
			mockSetup: func(m *mocks.MockGCPService) {
				if setup.Config.UseRealGCP {
					return
				}
				m.On("DeleteProject", "non-existent-project").Return(fmt.Errorf("project not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to delete project",
		},
		{
			name:      "Empty project ID",
			projectID: "",
			mockSetup: func(m *mocks.MockGCPService) {
				// No mock setup needed as validation happens before service call
			},
			expectedStatus: http.StatusNotFound, // Echo returns 404 for empty path params
			expectedError:  "Not Found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock if not using real GCP
			if !setup.Config.UseRealGCP && setup.MockService != nil {
				tt.mockSetup(setup.MockService)
			}

			// Create request
			url := fmt.Sprintf("/api/v1/projects/%s", tt.projectID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			// Execute request
			setup.Echo.ServeHTTP(rec, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				AssertErrorResponse(t, rec.Body.Bytes(), tt.expectedStatus, tt.expectedError)
			} else {
				AssertSuccessResponse(t, rec.Body.Bytes(), "Project deleted successfully")
			}

			// Reset mock expectations
			if !setup.Config.UseRealGCP && setup.MockService != nil {
				setup.MockService.ExpectedCalls = nil
				setup.MockService.Calls = nil
			}
		})
	}
}

func testProjectValidation(t *testing.T, setup *TestSetup, token string) {
	tests := []struct {
		name           string
		request        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing project ID",
			request: map[string]interface{}{
				"display_name": "Missing Project ID",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		},
		{
			name: "Missing display name",
			request: map[string]interface{}{
				"project_id": "missing-display-name",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		},
		{
			name: "Invalid parent type",
			request: map[string]interface{}{
				"project_id":   "invalid-parent-project",
				"display_name": "Invalid Parent Type",
				"parent_id":    "123456789",
				"parent_type":  "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		},
		{
			name:           "Invalid JSON",
			request:        "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if str, ok := tt.request.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.request)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			// Execute request
			setup.Echo.ServeHTTP(rec, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, rec.Code)
			AssertErrorResponse(t, rec.Body.Bytes(), tt.expectedStatus, tt.expectedError)
		})
	}
}

func TestProjectAuthenticationRequired(t *testing.T) {
	setup := SetupTestServer(t)
	defer CleanupTestResources(t, setup)

	// Setup authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(setup.AuthService.GetConfig())

	// Setup routes
	setup.Echo.POST("/api/v1/projects", setup.Handler.CreateProject, authMiddleware.RequireAuth())
	setup.Echo.GET("/api/v1/projects/:id", setup.Handler.GetProject, authMiddleware.RequireAuth())
	setup.Echo.DELETE("/api/v1/projects/:id", setup.Handler.DeleteProject, authMiddleware.RequireAuth())

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
			setup.Echo.ServeHTTP(rec, req)

			// Should return 401 Unauthorized
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			AssertErrorResponse(t, rec.Body.Bytes(), http.StatusUnauthorized, "unauthorized")
		})
	}
}
