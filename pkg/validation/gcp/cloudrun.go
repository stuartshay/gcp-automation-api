package gcp

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	// cloudRunServiceNameRegex defines the valid Cloud Run service name pattern
	// Must be lowercase, start with letter, contain only letters, numbers, hyphens
	cloudRunServiceNameRegex = regexp.MustCompile(`^[a-z]([a-z0-9-]*[a-z0-9])?$`)

	// Valid log levels for Cloud Run services
	validLogLevels = map[string]bool{
		"DEFAULT":   true,
		"DEBUG":     true,
		"INFO":      true,
		"NOTICE":    true,
		"WARNING":   true,
		"ERROR":     true,
		"CRITICAL":  true,
		"ALERT":     true,
		"EMERGENCY": true,
	}

	// Valid Cloud Run regions (subset of common regions, can be extended)
	validCloudRunRegions = map[string]bool{
		"us-central1":             true,
		"us-east1":                true,
		"us-east4":                true,
		"us-west1":                true,
		"us-west2":                true,
		"us-west3":                true,
		"us-west4":                true,
		"europe-west1":            true,
		"europe-west2":            true,
		"europe-west3":            true,
		"europe-west4":            true,
		"europe-west6":            true,
		"europe-north1":           true,
		"asia-east1":              true,
		"asia-east2":              true,
		"asia-northeast1":         true,
		"asia-northeast2":         true,
		"asia-northeast3":         true,
		"asia-south1":             true,
		"asia-southeast1":         true,
		"asia-southeast2":         true,
		"australia-southeast1":    true,
		"northamerica-northeast1": true,
		"southamerica-east1":      true,
	}

	// Valid export destination types
	validExportTypes = map[string]bool{
		"bigquery":      true,
		"cloud-storage": true,
		"pubsub":        true,
	}
)

// ValidateCloudRunServiceName validates a Cloud Run service name according to GCP naming conventions
func ValidateCloudRunServiceName(serviceName string) error {
	if serviceName == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	if len(serviceName) < 1 || len(serviceName) > 63 {
		return fmt.Errorf("service name must be between 1 and 63 characters")
	}

	if !cloudRunServiceNameRegex.MatchString(serviceName) {
		return fmt.Errorf("service name must start with a letter, contain only lowercase letters, numbers, and hyphens, and end with a letter or number")
	}

	// Check for reserved prefixes
	if strings.HasPrefix(serviceName, "goog-") {
		return fmt.Errorf("service name cannot start with 'goog-'")
	}

	if strings.Contains(serviceName, "google") {
		return fmt.Errorf("service name cannot contain 'google'")
	}

	return nil
}

// ValidateCloudRunRegion validates a Cloud Run region
func ValidateCloudRunRegion(region string) error {
	if region == "" {
		return fmt.Errorf("region cannot be empty")
	}

	if !validCloudRunRegions[region] {
		return fmt.Errorf("invalid Cloud Run region: %s", region)
	}

	return nil
}

// ValidateLogLevel validates a log level for Cloud Run services
func ValidateLogLevel(logLevel string) error {
	if logLevel == "" {
		// Empty log level defaults to INFO, which is valid
		return nil
	}

	if !validLogLevels[strings.ToUpper(logLevel)] {
		return fmt.Errorf("invalid log level: %s. Valid levels are: DEFAULT, DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL, ALERT, EMERGENCY", logLevel)
	}

	return nil
}

// ValidateRetentionDays validates log retention days
func ValidateRetentionDays(days int) error {
	if days < 1 {
		return fmt.Errorf("retention days must be at least 1")
	}

	if days > 3653 { // ~10 years maximum
		return fmt.Errorf("retention days cannot exceed 3653 (10 years)")
	}

	return nil
}

// ValidateExportDestinationType validates the type of log export destination
func ValidateExportDestinationType(exportType string) error {
	if exportType == "" {
		return fmt.Errorf("export destination type cannot be empty")
	}

	if !validExportTypes[strings.ToLower(exportType)] {
		return fmt.Errorf("invalid export destination type: %s. Valid types are: bigquery, cloud-storage, pubsub", exportType)
	}

	return nil
}

// ValidateMetricName validates a log-based metric name
func ValidateMetricName(metricName string) error {
	if metricName == "" {
		return fmt.Errorf("metric name cannot be empty")
	}

	if len(metricName) > 100 {
		return fmt.Errorf("metric name cannot exceed 100 characters")
	}

	// Metric names should be valid identifiers
	metricNameRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !metricNameRegex.MatchString(metricName) {
		return fmt.Errorf("metric name must start with a letter and contain only letters, numbers, and underscores")
	}

	return nil
}

// ValidateLogFilter validates a Cloud Logging filter expression
func ValidateLogFilter(filter string) error {
	if filter == "" {
		return fmt.Errorf("log filter cannot be empty")
	}

	if len(filter) > 20000 {
		return fmt.Errorf("log filter cannot exceed 20,000 characters")
	}

	// Basic validation for common filter syntax
	// In a real implementation, you might want more sophisticated validation
	if strings.Contains(filter, "severity") {
		validSeverities := []string{"DEFAULT", "DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "CRITICAL", "ALERT", "EMERGENCY"}
		hasValidSeverity := false
		for _, severity := range validSeverities {
			if strings.Contains(filter, severity) {
				hasValidSeverity = true
				break
			}
		}
		if !hasValidSeverity {
			return fmt.Errorf("filter contains 'severity' but no valid severity level found")
		}
	}

	return nil
}

// ValidateAlertCondition validates an alert condition
func ValidateAlertCondition(condition string) error {
	if condition == "" {
		return fmt.Errorf("alert condition cannot be empty")
	}

	if len(condition) > 1000 {
		return fmt.Errorf("alert condition cannot exceed 1,000 characters")
	}

	// Basic validation for alert condition syntax
	if !strings.Contains(condition, ">") && !strings.Contains(condition, "<") && !strings.Contains(condition, "=") {
		return fmt.Errorf("alert condition must contain a comparison operator (>, <, =)")
	}

	return nil
}

// ValidateNotificationChannel validates a notification channel ID
func ValidateNotificationChannel(channelID string) error {
	if channelID == "" {
		return fmt.Errorf("notification channel ID cannot be empty")
	}

	// Notification channel IDs are typically in the format projects/{project}/notificationChannels/{channel_id}
	if !strings.HasPrefix(channelID, "projects/") {
		return fmt.Errorf("notification channel ID must start with 'projects/'")
	}

	if !strings.Contains(channelID, "/notificationChannels/") {
		return fmt.Errorf("notification channel ID must contain '/notificationChannels/'")
	}

	return nil
}

// ValidateTimeout validates a timeout duration
func ValidateTimeout(timeout time.Duration) error {
	if timeout < 0 {
		return fmt.Errorf("timeout cannot be negative")
	}

	if timeout > 15*time.Minute {
		return fmt.Errorf("timeout cannot exceed 15 minutes for Cloud Run services")
	}

	return nil
}
