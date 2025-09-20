package gcp

import (
	"strings"
	"testing"
	"time"
)

func TestValidateCloudRunServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		wantError   bool
	}{
		{
			name:        "valid service name",
			serviceName: "my-service",
			wantError:   false,
		},
		{
			name:        "valid service name with numbers",
			serviceName: "api-service-v1",
			wantError:   false,
		},
		{
			name:        "valid single letter",
			serviceName: "a",
			wantError:   false,
		},
		{
			name:        "empty service name",
			serviceName: "",
			wantError:   true,
		},
		{
			name:        "service name with uppercase",
			serviceName: "My-Service",
			wantError:   true,
		},
		{
			name:        "service name starting with number",
			serviceName: "1-service",
			wantError:   true,
		},
		{
			name:        "service name ending with hyphen",
			serviceName: "my-service-",
			wantError:   true,
		},
		{
			name:        "service name starting with hyphen",
			serviceName: "-my-service",
			wantError:   true,
		},
		{
			name:        "service name with underscores",
			serviceName: "my_service",
			wantError:   true,
		},
		{
			name:        "service name with dots",
			serviceName: "my.service",
			wantError:   true,
		},
		{
			name:        "service name too long",
			serviceName: strings.Repeat("a", 64),
			wantError:   true,
		},
		{
			name:        "service name starting with goog",
			serviceName: "goog-service",
			wantError:   true,
		},
		{
			name:        "service name containing google",
			serviceName: "my-google-service",
			wantError:   true,
		},
		{
			name:        "service name with consecutive hyphens",
			serviceName: "my--service",
			wantError:   false, // This is actually valid in Cloud Run
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCloudRunServiceName(tt.serviceName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateCloudRunServiceName() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateCloudRunRegion(t *testing.T) {
	tests := []struct {
		name      string
		region    string
		wantError bool
	}{
		{
			name:      "valid US region",
			region:    "us-central1",
			wantError: false,
		},
		{
			name:      "valid Europe region",
			region:    "europe-west1",
			wantError: false,
		},
		{
			name:      "valid Asia region",
			region:    "asia-east1",
			wantError: false,
		},
		{
			name:      "empty region",
			region:    "",
			wantError: true,
		},
		{
			name:      "invalid region",
			region:    "invalid-region",
			wantError: true,
		},
		{
			name:      "zone instead of region",
			region:    "us-central1-a",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCloudRunRegion(tt.region)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateCloudRunRegion() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateLogLevel(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  string
		wantError bool
	}{
		{
			name:      "valid log level DEBUG",
			logLevel:  "DEBUG",
			wantError: false,
		},
		{
			name:      "valid log level INFO",
			logLevel:  "INFO",
			wantError: false,
		},
		{
			name:      "valid log level ERROR",
			logLevel:  "ERROR",
			wantError: false,
		},
		{
			name:      "valid log level CRITICAL",
			logLevel:  "CRITICAL",
			wantError: false,
		},
		{
			name:      "empty log level",
			logLevel:  "",
			wantError: false, // Empty defaults to INFO
		},
		{
			name:      "lowercase log level",
			logLevel:  "debug",
			wantError: false, // Should be converted to uppercase
		},
		{
			name:      "invalid log level",
			logLevel:  "INVALID",
			wantError: true,
		},
		{
			name:      "numeric log level",
			logLevel:  "1",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLogLevel(tt.logLevel)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateLogLevel() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateRetentionDays(t *testing.T) {
	tests := []struct {
		name      string
		days      int
		wantError bool
	}{
		{
			name:      "valid retention 30 days",
			days:      30,
			wantError: false,
		},
		{
			name:      "valid retention 1 day",
			days:      1,
			wantError: false,
		},
		{
			name:      "valid retention 365 days",
			days:      365,
			wantError: false,
		},
		{
			name:      "valid retention 3653 days",
			days:      3653,
			wantError: false,
		},
		{
			name:      "zero days",
			days:      0,
			wantError: true,
		},
		{
			name:      "negative days",
			days:      -1,
			wantError: true,
		},
		{
			name:      "too many days",
			days:      3654,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRetentionDays(tt.days)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateRetentionDays() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateExportDestinationType(t *testing.T) {
	tests := []struct {
		name       string
		exportType string
		wantError  bool
	}{
		{
			name:       "valid bigquery",
			exportType: "bigquery",
			wantError:  false,
		},
		{
			name:       "valid cloud-storage",
			exportType: "cloud-storage",
			wantError:  false,
		},
		{
			name:       "valid pubsub",
			exportType: "pubsub",
			wantError:  false,
		},
		{
			name:       "valid uppercase",
			exportType: "BIGQUERY",
			wantError:  false,
		},
		{
			name:       "empty export type",
			exportType: "",
			wantError:  true,
		},
		{
			name:       "invalid export type",
			exportType: "elasticsearch",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExportDestinationType(tt.exportType)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateExportDestinationType() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateMetricName(t *testing.T) {
	tests := []struct {
		name       string
		metricName string
		wantError  bool
	}{
		{
			name:       "valid metric name",
			metricName: "error_count",
			wantError:  false,
		},
		{
			name:       "valid metric name with numbers",
			metricName: "requests_per_second",
			wantError:  false,
		},
		{
			name:       "valid single letter",
			metricName: "a",
			wantError:  false,
		},
		{
			name:       "empty metric name",
			metricName: "",
			wantError:  true,
		},
		{
			name:       "metric name starting with number",
			metricName: "1_error_count",
			wantError:  true,
		},
		{
			name:       "metric name with hyphens",
			metricName: "error-count",
			wantError:  true,
		},
		{
			name:       "metric name with dots",
			metricName: "error.count",
			wantError:  true,
		},
		{
			name:       "metric name too long",
			metricName: strings.Repeat("a", 101),
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricName(tt.metricName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateMetricName() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateLogFilter(t *testing.T) {
	tests := []struct {
		name      string
		filter    string
		wantError bool
	}{
		{
			name:      "valid severity filter",
			filter:    "severity >= ERROR",
			wantError: false,
		},
		{
			name:      "valid resource filter",
			filter:    "resource.type = \"cloud_run_revision\"",
			wantError: false,
		},
		{
			name:      "valid complex filter",
			filter:    "severity >= WARNING AND resource.type = \"cloud_run_revision\"",
			wantError: false,
		},
		{
			name:      "empty filter",
			filter:    "",
			wantError: true,
		},
		{
			name:      "filter too long",
			filter:    strings.Repeat("a", 20001),
			wantError: true,
		},
		{
			name:      "filter with invalid severity",
			filter:    "severity >= INVALID_LEVEL",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLogFilter(tt.filter)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateLogFilter() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAlertCondition(t *testing.T) {
	tests := []struct {
		name      string
		condition string
		wantError bool
	}{
		{
			name:      "valid condition with greater than",
			condition: "error_rate > 0.05",
			wantError: false,
		},
		{
			name:      "valid condition with less than",
			condition: "response_time < 1000",
			wantError: false,
		},
		{
			name:      "valid condition with equals",
			condition: "status_code = 500",
			wantError: false,
		},
		{
			name:      "empty condition",
			condition: "",
			wantError: true,
		},
		{
			name:      "condition too long",
			condition: strings.Repeat("a", 1001),
			wantError: true,
		},
		{
			name:      "condition without operator",
			condition: "error_rate 0.05",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAlertCondition(tt.condition)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAlertCondition() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateNotificationChannel(t *testing.T) {
	tests := []struct {
		name      string
		channelID string
		wantError bool
	}{
		{
			name:      "valid notification channel",
			channelID: "projects/my-project/notificationChannels/12345",
			wantError: false,
		},
		{
			name:      "valid notification channel with long ID",
			channelID: "projects/my-project-123/notificationChannels/abcd-efgh-ijkl",
			wantError: false,
		},
		{
			name:      "empty channel ID",
			channelID: "",
			wantError: true,
		},
		{
			name:      "channel ID without projects prefix",
			channelID: "notificationChannels/12345",
			wantError: true,
		},
		{
			name:      "channel ID without notificationChannels",
			channelID: "projects/my-project/12345",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotificationChannel(tt.channelID)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateNotificationChannel() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateTimeout(t *testing.T) {
	tests := []struct {
		name      string
		timeout   time.Duration
		wantError bool
	}{
		{
			name:      "valid timeout 30 seconds",
			timeout:   30 * time.Second,
			wantError: false,
		},
		{
			name:      "valid timeout 5 minutes",
			timeout:   5 * time.Minute,
			wantError: false,
		},
		{
			name:      "valid timeout 15 minutes",
			timeout:   15 * time.Minute,
			wantError: false,
		},
		{
			name:      "zero timeout",
			timeout:   0,
			wantError: false,
		},
		{
			name:      "negative timeout",
			timeout:   -1 * time.Second,
			wantError: true,
		},
		{
			name:      "timeout too long",
			timeout:   16 * time.Minute,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimeout(tt.timeout)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateTimeout() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// Benchmark tests for performance validation
func BenchmarkValidateCloudRunServiceName(b *testing.B) {
	serviceName := "my-test-service-with-a-reasonable-length-name"
	for i := 0; i < b.N; i++ {
		_ = ValidateCloudRunServiceName(serviceName)
	}
}

func BenchmarkValidateLogLevel(b *testing.B) {
	logLevel := "INFO"
	for i := 0; i < b.N; i++ {
		_ = ValidateLogLevel(logLevel)
	}
}

func BenchmarkValidateCloudRunRegion(b *testing.B) {
	region := "us-central1"
	for i := 0; i < b.N; i++ {
		_ = ValidateCloudRunRegion(region)
	}
}
