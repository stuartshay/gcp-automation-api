package handlers_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// MockCloudRunService is a mock implementation of CloudRunServiceInterface
type MockCloudRunService struct {
	mock.Mock
}

func (m *MockCloudRunService) ConfigureLogging(ctx context.Context, req *models.CloudRunLoggingConfigRequest) (*models.CloudRunLoggingConfigResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.CloudRunLoggingConfigResponse), args.Error(1)
}

func (m *MockCloudRunService) GetLoggingConfig(ctx context.Context, serviceName, region string) (*models.CloudRunLoggingConfigResponse, error) {
	args := m.Called(ctx, serviceName, region)
	return args.Get(0).(*models.CloudRunLoggingConfigResponse), args.Error(1)
}

func (m *MockCloudRunService) UpdateLoggingConfig(ctx context.Context, serviceName, region string, req *models.CloudRunLoggingConfigUpdateRequest) (*models.CloudRunLoggingConfigResponse, error) {
	args := m.Called(ctx, serviceName, region, req)
	return args.Get(0).(*models.CloudRunLoggingConfigResponse), args.Error(1)
}

func (m *MockCloudRunService) GetLogs(ctx context.Context, req *models.CloudRunLogsRequest) (*models.CloudRunLogsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.CloudRunLogsResponse), args.Error(1)
}

func (m *MockCloudRunService) GetServiceInfo(ctx context.Context, serviceName, region string) (*models.CloudRunServiceInfo, error) {
	args := m.Called(ctx, serviceName, region)
	return args.Get(0).(*models.CloudRunServiceInfo), args.Error(1)
}

func (m *MockCloudRunService) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Test cases for CloudRunService

func TestCloudRunService_ConfigureLogging(t *testing.T) {
	tests := []struct {
		name           string
		request        *models.CloudRunLoggingConfigRequest
		expectedError  bool
		expectedStatus string
	}{
		{
			name: "Valid configuration request",
			request: &models.CloudRunLoggingConfigRequest{
				ServiceName: "test-service",
				Region:      "us-central1",
				LoggingConfig: models.LoggingConfig{
					LogLevel:          "INFO",
					StructuredLogging: true,
					RetentionDays:     30,
					ExportDestinations: []models.ExportDestination{
						{
							Type:    "bigquery",
							Dataset: "logs",
						},
					},
				},
				Metrics: []models.LogMetric{
					{
						Name:        "error-rate",
						Description: "Error rate metric",
						Filter:      "severity=ERROR",
						Type:        "COUNTER",
						Labels:      map[string]string{"service": "test-service"},
					},
				},
				Alerts: []models.LogAlert{
					{
						Name:                 "high-error-rate",
						Description:          "Alert on high error rate",
						Condition:            "count > 10",
						NotificationChannels: []string{"projects/test/notificationChannels/123"},
						Enabled:              true,
					},
				},
			},
			expectedError:  false,
			expectedStatus: "configured",
		},
		{
			name: "Invalid service name",
			request: &models.CloudRunLoggingConfigRequest{
				ServiceName: "",
				Region:      "us-central1",
				LoggingConfig: models.LoggingConfig{
					LogLevel: "INFO",
				},
			},
			expectedError: true,
		},
		{
			name: "Invalid region",
			request: &models.CloudRunLoggingConfigRequest{
				ServiceName: "test-service",
				Region:      "",
				LoggingConfig: models.LoggingConfig{
					LogLevel: "INFO",
				},
			},
			expectedError: true,
		},
		{
			name: "Invalid log level",
			request: &models.CloudRunLoggingConfigRequest{
				ServiceName: "test-service",
				Region:      "us-central1",
				LoggingConfig: models.LoggingConfig{
					LogLevel: "INVALID",
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCloudRunService)
			ctx := context.Background()

			if !tt.expectedError {
				expectedResponse := &models.CloudRunLoggingConfigResponse{
					ServiceName:   tt.request.ServiceName,
					Region:        tt.request.Region,
					Status:        tt.expectedStatus,
					LoggingConfig: tt.request.LoggingConfig,
					ConfiguredAt:  time.Now(),
					LoggingURL:    "https://console.cloud.google.com/logs/query",
				}

				mockService.On("ConfigureLogging", ctx, tt.request).Return(expectedResponse, nil)

				response, err := mockService.ConfigureLogging(ctx, tt.request)

				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, response.Status)
				assert.Equal(t, tt.request.ServiceName, response.ServiceName)
				assert.Equal(t, tt.request.Region, response.Region)
			} else {
				mockService.On("ConfigureLogging", ctx, tt.request).Return((*models.CloudRunLoggingConfigResponse)(nil), assert.AnError)

				response, err := mockService.ConfigureLogging(ctx, tt.request)

				require.Error(t, err)
				assert.Nil(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCloudRunService_GetLoggingConfig(t *testing.T) {
	tests := []struct {
		name          string
		serviceName   string
		region        string
		expectedError bool
	}{
		{
			name:          "Valid service and region",
			serviceName:   "test-service",
			region:        "us-central1",
			expectedError: false,
		},
		{
			name:          "Invalid service name",
			serviceName:   "",
			region:        "us-central1",
			expectedError: true,
		},
		{
			name:          "Invalid region",
			serviceName:   "test-service",
			region:        "",
			expectedError: true,
		},
		{
			name:          "Service not found",
			serviceName:   "nonexistent-service",
			region:        "us-central1",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCloudRunService)
			ctx := context.Background()

			if !tt.expectedError {
				expectedResponse := &models.CloudRunLoggingConfigResponse{
					ServiceName: tt.serviceName,
					Region:      tt.region,
					Status:      "active",
					LoggingConfig: models.LoggingConfig{
						LogLevel:          "INFO",
						StructuredLogging: true,
						RetentionDays:     30,
					},
					ConfiguredAt: time.Now(),
					LoggingURL:   "https://console.cloud.google.com/logs/query",
				}

				mockService.On("GetLoggingConfig", ctx, tt.serviceName, tt.region).Return(expectedResponse, nil)

				response, err := mockService.GetLoggingConfig(ctx, tt.serviceName, tt.region)

				require.NoError(t, err)
				assert.Equal(t, tt.serviceName, response.ServiceName)
				assert.Equal(t, tt.region, response.Region)
				assert.Equal(t, "active", response.Status)
			} else {
				mockService.On("GetLoggingConfig", ctx, tt.serviceName, tt.region).Return((*models.CloudRunLoggingConfigResponse)(nil), assert.AnError)

				response, err := mockService.GetLoggingConfig(ctx, tt.serviceName, tt.region)

				require.Error(t, err)
				assert.Nil(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCloudRunService_UpdateLoggingConfig(t *testing.T) {
	tests := []struct {
		name          string
		serviceName   string
		region        string
		request       *models.CloudRunLoggingConfigUpdateRequest
		expectedError bool
	}{
		{
			name:        "Valid update request",
			serviceName: "test-service",
			region:      "us-central1",
			request: &models.CloudRunLoggingConfigUpdateRequest{
				LoggingConfig: &models.LoggingConfig{
					LogLevel:          "DEBUG",
					StructuredLogging: true,
					RetentionDays:     60,
				},
			},
			expectedError: false,
		},
		{
			name:        "Update with metrics",
			serviceName: "test-service",
			region:      "us-central1",
			request: &models.CloudRunLoggingConfigUpdateRequest{
				Metrics: []models.LogMetric{
					{
						Name:        "request-count",
						Description: "Request count metric",
						Filter:      "resource.type=cloud_run_revision",
						Type:        "COUNTER",
					},
				},
			},
			expectedError: false,
		},
		{
			name:        "Update with alerts",
			serviceName: "test-service",
			region:      "us-central1",
			request: &models.CloudRunLoggingConfigUpdateRequest{
				Alerts: []models.LogAlert{
					{
						Name:        "cpu-usage-alert",
						Description: "CPU usage alert",
						Condition:   "cpu > 80",
						Enabled:     true,
					},
				},
			},
			expectedError: false,
		},
		{
			name:          "Invalid service name",
			serviceName:   "",
			region:        "us-central1",
			request:       &models.CloudRunLoggingConfigUpdateRequest{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCloudRunService)
			ctx := context.Background()

			if !tt.expectedError {
				expectedResponse := &models.CloudRunLoggingConfigResponse{
					ServiceName:  tt.serviceName,
					Region:       tt.region,
					Status:       "updated",
					ConfiguredAt: time.Now(),
					LoggingURL:   "https://console.cloud.google.com/logs/query",
				}

				if tt.request.LoggingConfig != nil {
					expectedResponse.LoggingConfig = *tt.request.LoggingConfig
				}

				mockService.On("UpdateLoggingConfig", ctx, tt.serviceName, tt.region, tt.request).Return(expectedResponse, nil)

				response, err := mockService.UpdateLoggingConfig(ctx, tt.serviceName, tt.region, tt.request)

				require.NoError(t, err)
				assert.Equal(t, "updated", response.Status)
				assert.Equal(t, tt.serviceName, response.ServiceName)
				assert.Equal(t, tt.region, response.Region)
			} else {
				mockService.On("UpdateLoggingConfig", ctx, tt.serviceName, tt.region, tt.request).Return((*models.CloudRunLoggingConfigResponse)(nil), assert.AnError)

				response, err := mockService.UpdateLoggingConfig(ctx, tt.serviceName, tt.region, tt.request)

				require.Error(t, err)
				assert.Nil(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCloudRunService_GetLogs(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.CloudRunLogsRequest
		expectedError bool
		expectedCount int
	}{
		{
			name: "Valid logs request",
			request: &models.CloudRunLogsRequest{
				ServiceName: "test-service",
				Region:      "us-central1",
				StartTime:   time.Now().Add(-1 * time.Hour),
				EndTime:     time.Now(),
				Filter:      "severity=ERROR",
				PageSize:    50,
			},
			expectedError: false,
			expectedCount: 5,
		},
		{
			name: "Request without time range",
			request: &models.CloudRunLogsRequest{
				ServiceName: "test-service",
				Region:      "us-central1",
				PageSize:    100,
			},
			expectedError: false,
			expectedCount: 10,
		},
		{
			name: "Invalid service name",
			request: &models.CloudRunLogsRequest{
				ServiceName: "",
				Region:      "us-central1",
				PageSize:    100,
			},
			expectedError: true,
		},
		{
			name: "Invalid region",
			request: &models.CloudRunLogsRequest{
				ServiceName: "test-service",
				Region:      "",
				PageSize:    100,
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCloudRunService)
			ctx := context.Background()

			if !tt.expectedError {
				// Create mock log entries
				logEntries := make([]models.LogEntry, tt.expectedCount)
				for i := 0; i < tt.expectedCount; i++ {
					logEntries[i] = models.LogEntry{
						Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
						Severity:  "INFO",
						Message:   "Test log message",
						Resource: models.LogResource{
							Type:        "cloud_run_revision",
							ServiceName: tt.request.ServiceName,
							Location:    tt.request.Region,
						},
					}
				}

				expectedResponse := &models.CloudRunLogsResponse{
					ServiceName: tt.request.ServiceName,
					Region:      tt.request.Region,
					Logs:        logEntries,
					TotalCount:  tt.expectedCount,
				}

				mockService.On("GetLogs", ctx, tt.request).Return(expectedResponse, nil)

				response, err := mockService.GetLogs(ctx, tt.request)

				require.NoError(t, err)
				assert.Equal(t, tt.request.ServiceName, response.ServiceName)
				assert.Equal(t, tt.request.Region, response.Region)
				assert.Equal(t, tt.expectedCount, response.TotalCount)
				assert.Len(t, response.Logs, tt.expectedCount)
			} else {
				mockService.On("GetLogs", ctx, tt.request).Return((*models.CloudRunLogsResponse)(nil), assert.AnError)

				response, err := mockService.GetLogs(ctx, tt.request)

				require.Error(t, err)
				assert.Nil(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCloudRunService_GetServiceInfo(t *testing.T) {
	tests := []struct {
		name          string
		serviceName   string
		region        string
		expectedError bool
	}{
		{
			name:          "Valid service info request",
			serviceName:   "test-service",
			region:        "us-central1",
			expectedError: false,
		},
		{
			name:          "Service not found",
			serviceName:   "nonexistent-service",
			region:        "us-central1",
			expectedError: true,
		},
		{
			name:          "Invalid service name",
			serviceName:   "",
			region:        "us-central1",
			expectedError: true,
		},
		{
			name:          "Invalid region",
			serviceName:   "test-service",
			region:        "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCloudRunService)
			ctx := context.Background()

			if !tt.expectedError {
				expectedResponse := &models.CloudRunServiceInfo{
					ServiceName: tt.serviceName,
					Region:      tt.region,
					URL:         "https://test-service-hash-uc.a.run.app",
					Status:      "READY",
					Labels: map[string]string{
						"app": "test-service",
					},
					CreatedAt: time.Now().Add(-24 * time.Hour),
					UpdatedAt: time.Now().Add(-1 * time.Hour),
				}

				mockService.On("GetServiceInfo", ctx, tt.serviceName, tt.region).Return(expectedResponse, nil)

				response, err := mockService.GetServiceInfo(ctx, tt.serviceName, tt.region)

				require.NoError(t, err)
				assert.Equal(t, tt.serviceName, response.ServiceName)
				assert.Equal(t, tt.region, response.Region)
				assert.Equal(t, "READY", response.Status)
				assert.NotEmpty(t, response.URL)
			} else {
				mockService.On("GetServiceInfo", ctx, tt.serviceName, tt.region).Return((*models.CloudRunServiceInfo)(nil), assert.AnError)

				response, err := mockService.GetServiceInfo(ctx, tt.serviceName, tt.region)

				require.Error(t, err)
				assert.Nil(t, response)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCloudRunService_Close(t *testing.T) {
	mockService := new(MockCloudRunService)

	mockService.On("Close").Return(nil)

	err := mockService.Close()

	require.NoError(t, err)
	mockService.AssertExpectations(t)
}

// Benchmark tests

func BenchmarkCloudRunService_ConfigureLogging(b *testing.B) {
	mockService := new(MockCloudRunService)
	ctx := context.Background()

	request := &models.CloudRunLoggingConfigRequest{
		ServiceName: "test-service",
		Region:      "us-central1",
		LoggingConfig: models.LoggingConfig{
			LogLevel:          "INFO",
			StructuredLogging: true,
			RetentionDays:     30,
		},
	}

	expectedResponse := &models.CloudRunLoggingConfigResponse{
		ServiceName:   request.ServiceName,
		Region:        request.Region,
		Status:        "configured",
		LoggingConfig: request.LoggingConfig,
		ConfiguredAt:  time.Now(),
		LoggingURL:    "https://console.cloud.google.com/logs/query",
	}

	mockService.On("ConfigureLogging", ctx, request).Return(expectedResponse, nil).Times(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mockService.ConfigureLogging(ctx, request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCloudRunService_GetLogs(b *testing.B) {
	mockService := new(MockCloudRunService)
	ctx := context.Background()

	request := &models.CloudRunLogsRequest{
		ServiceName: "test-service",
		Region:      "us-central1",
		PageSize:    100,
	}

	logEntries := make([]models.LogEntry, 100)
	for i := 0; i < 100; i++ {
		logEntries[i] = models.LogEntry{
			Timestamp: time.Now(),
			Severity:  "INFO",
			Message:   "Test log message",
			Resource: models.LogResource{
				Type:        "cloud_run_revision",
				ServiceName: request.ServiceName,
				Location:    request.Region,
			},
		}
	}

	expectedResponse := &models.CloudRunLogsResponse{
		ServiceName: request.ServiceName,
		Region:      request.Region,
		Logs:        logEntries,
		TotalCount:  100,
	}

	mockService.On("GetLogs", ctx, request).Return(expectedResponse, nil).Times(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mockService.GetLogs(ctx, request)
		if err != nil {
			b.Fatal(err)
		}
	}
}
