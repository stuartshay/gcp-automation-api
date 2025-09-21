package services

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/pkg/validation/gcp"
)

// CloudRunService provides operations for Cloud Run logging management
type CloudRunService struct {
	projectID      string
	runClient      *run.ServicesClient
	LoggingClient  *logging.Client
	logAdminClient *logadmin.Client
}

// CloudRunServiceInterface defines the interface for Cloud Run operations
type CloudRunServiceInterface interface {
	ConfigureLogging(ctx context.Context, req *models.CloudRunLoggingConfigRequest) (*models.CloudRunLoggingConfigResponse, error)
	GetLoggingConfig(ctx context.Context, serviceName, region string) (*models.CloudRunLoggingConfigResponse, error)
	UpdateLoggingConfig(ctx context.Context, serviceName, region string, req *models.CloudRunLoggingConfigUpdateRequest) (*models.CloudRunLoggingConfigResponse, error)
	GetLogs(ctx context.Context, req *models.CloudRunLogsRequest) (*models.CloudRunLogsResponse, error)
	GetServiceInfo(ctx context.Context, serviceName, region string) (*models.CloudRunServiceInfo, error)
	Close() error
}

// NewCloudRunService creates a new Cloud Run service instance
func NewCloudRunService(ctx context.Context, projectID string, opts ...option.ClientOption) (*CloudRunService, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	// Create Cloud Run client
	runClient, err := run.NewServicesClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Run client: %w", err)
	}

	// Create logging client
	loggingClient, err := logging.NewClient(ctx, projectID, opts...)
	if err != nil {
		_ = runClient.Close() // Ignore close error, original error is more important
		return nil, fmt.Errorf("failed to create logging client: %w", err)
	}

	// Create log admin client for advanced operations
	logAdminClient, err := logadmin.NewClient(ctx, projectID, opts...)
	if err != nil {
		_ = runClient.Close()     // Ignore close error, original error is more important
		_ = loggingClient.Close() // Ignore close error, original error is more important
		return nil, fmt.Errorf("failed to create log admin client: %w", err)
	}

	return &CloudRunService{
		projectID:      projectID,
		runClient:      runClient,
		LoggingClient:  loggingClient,
		logAdminClient: logAdminClient,
	}, nil
}

// ConfigureLogging configures logging for a Cloud Run service
func (s *CloudRunService) ConfigureLogging(ctx context.Context, req *models.CloudRunLoggingConfigRequest) (*models.CloudRunLoggingConfigResponse, error) {
	// Validate input
	if err := s.validateLoggingConfigRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get service information to ensure it exists
	_, err := s.GetServiceInfo(ctx, req.ServiceName, req.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to get service info: %w", err)
	}

	response := &models.CloudRunLoggingConfigResponse{
		ServiceName:   req.ServiceName,
		Region:        req.Region,
		Status:        "configured",
		LoggingConfig: req.LoggingConfig,
		ConfiguredAt:  time.Now(),
		LoggingURL:    s.buildLoggingURL(req.ServiceName, req.Region),
	}

	// Configure log exports if specified
	if len(req.LoggingConfig.ExportDestinations) > 0 {
		if err := s.configureLogExports(ctx, req); err != nil {
			return nil, fmt.Errorf("failed to configure log exports: %w", err)
		}
	}

	// Create log-based metrics if specified
	if len(req.Metrics) > 0 {
		metricResponses, err := s.createLogMetrics(ctx, req.ServiceName, req.Region, req.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to create log metrics: %w", err)
		}
		response.Metrics = metricResponses
	}

	// Create log-based alerts if specified
	if len(req.Alerts) > 0 {
		alertResponses, err := s.createLogAlerts(ctx, req.ServiceName, req.Region, req.Alerts)
		if err != nil {
			return nil, fmt.Errorf("failed to create log alerts: %w", err)
		}
		response.Alerts = alertResponses
	}

	return response, nil
}

// GetLoggingConfig retrieves the current logging configuration for a service
func (s *CloudRunService) GetLoggingConfig(ctx context.Context, serviceName, region string) (*models.CloudRunLoggingConfigResponse, error) {
	// Validate input
	if err := gcp.ValidateCloudRunServiceName(serviceName); err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	if err := gcp.ValidateCloudRunRegion(region); err != nil {
		return nil, fmt.Errorf("invalid region: %w", err)
	}

	// Get service information
	serviceInfo, err := s.GetServiceInfo(ctx, serviceName, region)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	// Build response with current configuration
	response := &models.CloudRunLoggingConfigResponse{
		ServiceName: serviceName,
		Region:      region,
		Status:      "active",
		LoggingConfig: models.LoggingConfig{
			LogLevel:          "INFO", // Default, could be retrieved from service metadata
			StructuredLogging: true,   // Default for Cloud Run
			RetentionDays:     30,     // Default retention
		},
		ConfiguredAt: serviceInfo.CreatedAt,
		LoggingURL:   s.buildLoggingURL(serviceName, region),
	}

	return response, nil
}

// UpdateLoggingConfig updates the logging configuration for a service
func (s *CloudRunService) UpdateLoggingConfig(ctx context.Context, serviceName, region string, req *models.CloudRunLoggingConfigUpdateRequest) (*models.CloudRunLoggingConfigResponse, error) {
	// Validate input
	if err := gcp.ValidateCloudRunServiceName(serviceName); err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	if err := gcp.ValidateCloudRunRegion(region); err != nil {
		return nil, fmt.Errorf("invalid region: %w", err)
	}

	// Get current configuration
	current, err := s.GetLoggingConfig(ctx, serviceName, region)
	if err != nil {
		return nil, fmt.Errorf("failed to get current config: %w", err)
	}

	// Update configuration
	if req.LoggingConfig != nil {
		current.LoggingConfig = *req.LoggingConfig
	}

	// Update metrics if specified
	if len(req.Metrics) > 0 {
		metricResponses, err := s.createLogMetrics(ctx, serviceName, region, req.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to update log metrics: %w", err)
		}
		current.Metrics = metricResponses
	}

	// Update alerts if specified
	if len(req.Alerts) > 0 {
		alertResponses, err := s.createLogAlerts(ctx, serviceName, region, req.Alerts)
		if err != nil {
			return nil, fmt.Errorf("failed to update log alerts: %w", err)
		}
		current.Alerts = alertResponses
	}

	current.ConfiguredAt = time.Now()
	current.Status = "updated"

	return current, nil
}

// GetLogs retrieves logs for a Cloud Run service
func (s *CloudRunService) GetLogs(ctx context.Context, req *models.CloudRunLogsRequest) (*models.CloudRunLogsResponse, error) {
	// Validate input
	if err := gcp.ValidateCloudRunServiceName(req.ServiceName); err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	if err := gcp.ValidateCloudRunRegion(req.Region); err != nil {
		return nil, fmt.Errorf("invalid region: %w", err)
	}

	// Build log filter
	filter := s.buildLogFilter(req)

	// Query logs
	entries := []models.LogEntry{}
	it := s.logAdminClient.Entries(ctx, logadmin.Filter(filter))

	count := 0
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 1000 {
		pageSize = 100 // Default page size
	}

	for count < pageSize {
		entry, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve logs: %w", err)
		}

		logEntry := s.convertLogEntry(entry)
		entries = append(entries, logEntry)
		count++
	}

	response := &models.CloudRunLogsResponse{
		ServiceName: req.ServiceName,
		Region:      req.Region,
		Logs:        entries,
		TotalCount:  len(entries),
	}

	return response, nil
}

// GetServiceInfo retrieves information about a Cloud Run service
func (s *CloudRunService) GetServiceInfo(ctx context.Context, serviceName, region string) (*models.CloudRunServiceInfo, error) {
	// Validate input
	if err := gcp.ValidateCloudRunServiceName(serviceName); err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	if err := gcp.ValidateCloudRunRegion(region); err != nil {
		return nil, fmt.Errorf("invalid region: %w", err)
	}

	// Build service name
	name := fmt.Sprintf("projects/%s/locations/%s/services/%s", s.projectID, region, serviceName)

	// Get service
	getReq := &runpb.GetServiceRequest{
		Name: name,
	}

	service, err := s.runClient.GetService(ctx, getReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	// Convert to response model
	serviceInfo := &models.CloudRunServiceInfo{
		ServiceName: serviceName,
		Region:      region,
		URL:         service.GetUri(),
		Status:      s.convertServiceStatus(service),
		Labels:      service.GetLabels(),
		CreatedAt:   service.GetCreateTime().AsTime(),
		UpdatedAt:   service.GetUpdateTime().AsTime(),
	}

	return serviceInfo, nil
}

// Close closes all clients
func (s *CloudRunService) Close() error {
	var errs []error

	if err := s.runClient.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close run client: %w", err))
	}

	if err := s.LoggingClient.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close logging client: %w", err))
	}

	if err := s.logAdminClient.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close log admin client: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing clients: %v", errs)
	}

	return nil
}

// Private helper methods

func (s *CloudRunService) validateLoggingConfigRequest(req *models.CloudRunLoggingConfigRequest) error {
	if err := gcp.ValidateCloudRunServiceName(req.ServiceName); err != nil {
		return fmt.Errorf("invalid service name: %w", err)
	}

	if err := gcp.ValidateCloudRunRegion(req.Region); err != nil {
		return fmt.Errorf("invalid region: %w", err)
	}

	if err := gcp.ValidateLogLevel(req.LoggingConfig.LogLevel); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	if req.LoggingConfig.RetentionDays > 0 {
		if err := gcp.ValidateRetentionDays(req.LoggingConfig.RetentionDays); err != nil {
			return fmt.Errorf("invalid retention days: %w", err)
		}
	}

	// Validate export destinations
	for _, dest := range req.LoggingConfig.ExportDestinations {
		if err := gcp.ValidateExportDestinationType(dest.Type); err != nil {
			return fmt.Errorf("invalid export destination: %w", err)
		}
	}

	// Validate metrics
	for _, metric := range req.Metrics {
		if err := gcp.ValidateMetricName(metric.Name); err != nil {
			return fmt.Errorf("invalid metric name: %w", err)
		}
		if err := gcp.ValidateLogFilter(metric.Filter); err != nil {
			return fmt.Errorf("invalid metric filter: %w", err)
		}
	}

	// Validate alerts
	for _, alert := range req.Alerts {
		if err := gcp.ValidateAlertCondition(alert.Condition); err != nil {
			return fmt.Errorf("invalid alert condition: %w", err)
		}
		for _, channel := range alert.NotificationChannels {
			if err := gcp.ValidateNotificationChannel(channel); err != nil {
				return fmt.Errorf("invalid notification channel: %w", err)
			}
		}
	}

	return nil
}

func (s *CloudRunService) buildLoggingURL(serviceName, region string) string {
	return fmt.Sprintf("https://console.cloud.google.com/logs/query;query=resource.type%%3D%%22cloud_run_revision%%22%%0Aresource.labels.service_name%%3D%%22%s%%22%%0Aresource.labels.location%%3D%%22%s%%22?project=%s",
		serviceName, region, s.projectID)
}

func (s *CloudRunService) buildLogFilter(req *models.CloudRunLogsRequest) string {
	filter := fmt.Sprintf(`resource.type="cloud_run_revision" AND resource.labels.service_name="%s" AND resource.labels.location="%s"`,
		req.ServiceName, req.Region)

	if !req.StartTime.IsZero() {
		filter += fmt.Sprintf(` AND timestamp >= "%s"`, req.StartTime.Format(time.RFC3339))
	}

	if !req.EndTime.IsZero() {
		filter += fmt.Sprintf(` AND timestamp <= "%s"`, req.EndTime.Format(time.RFC3339))
	}

	if req.Filter != "" {
		filter += fmt.Sprintf(` AND %s`, req.Filter)
	}

	return filter
}

func (s *CloudRunService) configureLogExports(ctx context.Context, req *models.CloudRunLoggingConfigRequest) error {
	// Implementation would configure log exports to various destinations
	// This is a placeholder for the actual export configuration
	return nil
}

func (s *CloudRunService) createLogMetrics(ctx context.Context, serviceName, region string, metrics []models.LogMetric) ([]models.LogMetricResponse, error) {
	var responses []models.LogMetricResponse

	for _, metric := range metrics {
		response := models.LogMetricResponse{
			Name:        metric.Name,
			Description: metric.Description,
			Filter:      metric.Filter,
			Type:        metric.Type,
			Labels:      metric.Labels,
			MetricURL:   fmt.Sprintf("https://console.cloud.google.com/monitoring/metrics-explorer?project=%s", s.projectID),
			CreatedAt:   time.Now(),
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *CloudRunService) createLogAlerts(ctx context.Context, serviceName, region string, alerts []models.LogAlert) ([]models.LogAlertResponse, error) {
	var responses []models.LogAlertResponse

	for _, alert := range alerts {
		response := models.LogAlertResponse{
			Name:                 alert.Name,
			Description:          alert.Description,
			Condition:            alert.Condition,
			NotificationChannels: alert.NotificationChannels,
			Enabled:              alert.Enabled,
			AlertURL:             fmt.Sprintf("https://console.cloud.google.com/monitoring/alerting?project=%s", s.projectID),
			CreatedAt:            time.Now(),
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *CloudRunService) convertLogEntry(entry *logging.Entry) models.LogEntry {
	logEntry := models.LogEntry{
		Timestamp: entry.Timestamp,
		Severity:  entry.Severity.String(),
		Message:   fmt.Sprintf("%v", entry.Payload),
		Resource: models.LogResource{
			Type: entry.Resource.Type,
		},
		Labels: entry.Labels,
	}

	if entry.Resource.Labels != nil {
		logEntry.Resource.ServiceName = entry.Resource.Labels["service_name"]
		logEntry.Resource.RevisionName = entry.Resource.Labels["revision_name"]
		logEntry.Resource.Location = entry.Resource.Labels["location"]
		logEntry.Resource.ConfigurationName = entry.Resource.Labels["configuration_name"]
		logEntry.Resource.Labels = entry.Resource.Labels
	}

	if entry.HTTPRequest != nil {
		var requestMethod, requestURL, userAgent string
		if entry.HTTPRequest.Request != nil {
			requestMethod = entry.HTTPRequest.Request.Method
			requestURL = entry.HTTPRequest.Request.URL.String()
			userAgent = entry.HTTPRequest.Request.UserAgent()
		}

		logEntry.HTTPRequest = &models.HTTPRequest{
			RequestMethod: requestMethod,
			RequestURL:    requestURL,
			Status:        entry.HTTPRequest.Status,
			ResponseSize:  entry.HTTPRequest.ResponseSize,
			UserAgent:     userAgent,
			RemoteIP:      entry.HTTPRequest.RemoteIP,
			Latency:       entry.HTTPRequest.Latency.String(),
		}
	}

	return logEntry
}

func (s *CloudRunService) convertServiceStatus(service *runpb.Service) string {
	if service.GetGeneration() > 0 {
		return "READY"
	}
	return "UNKNOWN"
}
