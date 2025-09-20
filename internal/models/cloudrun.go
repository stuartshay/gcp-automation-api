package models

import (
	"time"
)

// CloudRunLoggingConfigRequest represents a request to configure logging for a Cloud Run service
type CloudRunLoggingConfigRequest struct {
	ServiceName   string        `json:"service_name" binding:"required" example:"my-api-service"`
	Region        string        `json:"region" binding:"required" example:"us-central1"`
	LoggingConfig LoggingConfig `json:"logging_config" binding:"required"`
	Metrics       []LogMetric   `json:"metrics,omitempty"`
	Alerts        []LogAlert    `json:"alerts,omitempty"`
}

// LoggingConfig represents the logging configuration for a Cloud Run service
type LoggingConfig struct {
	LogLevel           string              `json:"log_level,omitempty" example:"INFO"`
	StructuredLogging  bool                `json:"structured_logging" example:"true"`
	RetentionDays      int                 `json:"retention_days" example:"30"`
	ExportDestinations []ExportDestination `json:"export_destinations,omitempty"`
	CustomFields       map[string]string   `json:"custom_fields,omitempty" example:"environment:production,team:backend"`
	SamplingRate       float64             `json:"sampling_rate,omitempty" example:"0.1"`
}

// ExportDestination represents a destination for log exports
type ExportDestination struct {
	Type    string            `json:"type" binding:"required" example:"bigquery"`
	Dataset string            `json:"dataset,omitempty" example:"logs_dataset"`
	Table   string            `json:"table,omitempty" example:"cloudrun_logs"`
	Bucket  string            `json:"bucket,omitempty" example:"logs-bucket"`
	Topic   string            `json:"topic,omitempty" example:"projects/my-project/topics/logs"`
	Filter  string            `json:"filter,omitempty" example:"severity >= WARNING"`
	Labels  map[string]string `json:"labels,omitempty" example:"environment:production"`
}

// LogMetric represents a log-based metric configuration
type LogMetric struct {
	Name        string            `json:"name" binding:"required" example:"error_rate"`
	Description string            `json:"description,omitempty" example:"Rate of error logs"`
	Filter      string            `json:"filter" binding:"required" example:"severity >= ERROR"`
	Type        string            `json:"type" example:"counter"`
	Labels      map[string]string `json:"labels,omitempty" example:"service:api"`
}

// LogAlert represents a log-based alert configuration
type LogAlert struct {
	Name                 string   `json:"name" binding:"required" example:"high_error_rate"`
	Description          string   `json:"description,omitempty" example:"Alert when error rate exceeds threshold"`
	Condition            string   `json:"condition" binding:"required" example:"error_rate > 0.05"`
	NotificationChannels []string `json:"notification_channels" example:"projects/my-project/notificationChannels/12345"`
	Enabled              bool     `json:"enabled" example:"true"`
}

// CloudRunLoggingConfigResponse represents the response after configuring logging
type CloudRunLoggingConfigResponse struct {
	ServiceName   string              `json:"service_name" example:"my-api-service"`
	Region        string              `json:"region" example:"us-central1"`
	Status        string              `json:"status" example:"configured"`
	LoggingConfig LoggingConfig       `json:"logging_config"`
	Metrics       []LogMetricResponse `json:"metrics,omitempty"`
	Alerts        []LogAlertResponse  `json:"alerts,omitempty"`
	ConfiguredAt  time.Time           `json:"configured_at" example:"2025-09-20T10:00:00Z"`
	LoggingURL    string              `json:"logging_url,omitempty" example:"https://console.cloud.google.com/logs/query"`
}

// LogMetricResponse represents the response for a created log metric
type LogMetricResponse struct {
	Name        string            `json:"name" example:"error_rate"`
	Description string            `json:"description" example:"Rate of error logs"`
	Filter      string            `json:"filter" example:"severity >= ERROR"`
	Type        string            `json:"type" example:"counter"`
	Labels      map[string]string `json:"labels,omitempty"`
	MetricURL   string            `json:"metric_url,omitempty" example:"https://console.cloud.google.com/monitoring/metrics-explorer"`
	CreatedAt   time.Time         `json:"created_at" example:"2025-09-20T10:00:00Z"`
}

// LogAlertResponse represents the response for a created log alert
type LogAlertResponse struct {
	Name                 string    `json:"name" example:"high_error_rate"`
	Description          string    `json:"description" example:"Alert when error rate exceeds threshold"`
	Condition            string    `json:"condition" example:"error_rate > 0.05"`
	NotificationChannels []string  `json:"notification_channels"`
	Enabled              bool      `json:"enabled" example:"true"`
	AlertURL             string    `json:"alert_url,omitempty" example:"https://console.cloud.google.com/monitoring/alerting"`
	CreatedAt            time.Time `json:"created_at" example:"2025-09-20T10:00:00Z"`
}

// CloudRunLogsRequest represents a request to retrieve logs for a Cloud Run service
type CloudRunLogsRequest struct {
	ServiceName string    `json:"service_name" form:"service_name" binding:"required" example:"my-api-service"`
	Region      string    `json:"region" form:"region" binding:"required" example:"us-central1"`
	StartTime   time.Time `json:"start_time" form:"start_time" example:"2025-09-20T09:00:00Z"`
	EndTime     time.Time `json:"end_time" form:"end_time" example:"2025-09-20T10:00:00Z"`
	Filter      string    `json:"filter" form:"filter" example:"severity >= WARNING"`
	PageSize    int       `json:"page_size" form:"page_size" example:"100"`
	PageToken   string    `json:"page_token" form:"page_token" example:""`
}

// CloudRunLogsResponse represents the response containing Cloud Run service logs
type CloudRunLogsResponse struct {
	ServiceName   string     `json:"service_name" example:"my-api-service"`
	Region        string     `json:"region" example:"us-central1"`
	Logs          []LogEntry `json:"logs"`
	NextPageToken string     `json:"next_page_token,omitempty" example:"abc123"`
	TotalCount    int        `json:"total_count" example:"150"`
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp   time.Time         `json:"timestamp" example:"2025-09-20T10:00:00Z"`
	Severity    string            `json:"severity" example:"INFO"`
	Message     string            `json:"message" example:"Request processed successfully"`
	Resource    LogResource       `json:"resource"`
	Labels      map[string]string `json:"labels,omitempty"`
	HTTPRequest *HTTPRequest      `json:"http_request,omitempty"`
	SourceFile  string            `json:"source_file,omitempty" example:"main.go"`
	SourceLine  int               `json:"source_line,omitempty" example:"42"`
	TraceID     string            `json:"trace_id,omitempty" example:"projects/my-project/traces/12345"`
}

// LogResource represents the resource information for a log entry
type LogResource struct {
	Type              string            `json:"type" example:"cloud_run_revision"`
	ServiceName       string            `json:"service_name" example:"my-api-service"`
	RevisionName      string            `json:"revision_name" example:"my-api-service-00001"`
	Location          string            `json:"location" example:"us-central1"`
	ConfigurationName string            `json:"configuration_name" example:"my-api-service"`
	Labels            map[string]string `json:"labels,omitempty"`
}

// HTTPRequest represents HTTP request information in logs
type HTTPRequest struct {
	RequestMethod string `json:"request_method" example:"GET"`
	RequestURL    string `json:"request_url" example:"/api/v1/health"`
	Status        int    `json:"status" example:"200"`
	ResponseSize  int64  `json:"response_size" example:"1024"`
	UserAgent     string `json:"user_agent" example:"curl/7.68.0"`
	RemoteIP      string `json:"remote_ip" example:"203.0.113.1"`
	Latency       string `json:"latency" example:"0.123s"`
}

// CloudRunLoggingConfigUpdateRequest represents a request to update logging configuration
type CloudRunLoggingConfigUpdateRequest struct {
	LoggingConfig *LoggingConfig `json:"logging_config,omitempty"`
	Metrics       []LogMetric    `json:"metrics,omitempty"`
	Alerts        []LogAlert     `json:"alerts,omitempty"`
}

// CloudRunServiceInfo represents basic information about a Cloud Run service
type CloudRunServiceInfo struct {
	ServiceName string            `json:"service_name" example:"my-api-service"`
	Region      string            `json:"region" example:"us-central1"`
	URL         string            `json:"url" example:"https://my-api-service-hash-uc.a.run.app"`
	Status      string            `json:"status" example:"READY"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreatedAt   time.Time         `json:"created_at" example:"2025-09-20T09:00:00Z"`
	UpdatedAt   time.Time         `json:"updated_at" example:"2025-09-20T10:00:00Z"`
}

// ErrorResponse represents an error response for Cloud Run operations
type CloudRunErrorResponse struct {
	Error     string `json:"error" example:"validation_failed"`
	Message   string `json:"message" example:"Invalid service name format"`
	Code      int    `json:"code" example:"400"`
	Details   string `json:"details,omitempty" example:"Service name must contain only lowercase letters, numbers, and hyphens"`
	Timestamp string `json:"timestamp" example:"2025-09-20T10:00:00Z"`
}
