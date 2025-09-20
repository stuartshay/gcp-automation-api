package sdk

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// bucketNameRegex defines the valid bucket name pattern for GCS
	bucketNameRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-._]*[a-z0-9])?$`)

	// objectNameMaxLength is the maximum length for object names
	objectNameMaxLength = 1024
)

// ValidateBucketName validates a GCS bucket name according to GCS naming rules
func ValidateBucketName(name string) error {
	if name == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}

	if len(name) < 3 {
		return fmt.Errorf("bucket name must be at least 3 characters long")
	}

	if len(name) > 63 {
		return fmt.Errorf("bucket name must be 63 characters or less")
	}

	if !bucketNameRegex.MatchString(name) {
		return fmt.Errorf("bucket name contains invalid characters or format")
	}

	// Cannot start or end with periods or hyphens
	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") ||
		strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("bucket name cannot start or end with periods or hyphens")
	}

	// Cannot contain consecutive periods
	if strings.Contains(name, "..") {
		return fmt.Errorf("bucket name cannot contain consecutive periods")
	}

	// Cannot be formatted as an IP address
	if isIPAddress(name) {
		return fmt.Errorf("bucket name cannot be formatted as an IP address")
	}

	// Cannot start with "goog" prefix
	if strings.HasPrefix(name, "goog") {
		return fmt.Errorf("bucket name cannot start with 'goog' prefix")
	}

	// Cannot contain "google" in the name
	if strings.Contains(strings.ToLower(name), "google") {
		return fmt.Errorf("bucket name cannot contain 'google'")
	}

	return nil
}

// ValidateObjectName validates a GCS object name
func ValidateObjectName(name string) error {
	if name == "" {
		return fmt.Errorf("object name cannot be empty")
	}

	if len(name) > objectNameMaxLength {
		return fmt.Errorf("object name must be %d characters or less", objectNameMaxLength)
	}

	// Check for invalid characters
	invalidChars := []string{"\n", "\r", "\x00"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("object name contains invalid character")
		}
	}

	// Cannot be "." or ".."
	if name == "." || name == ".." {
		return fmt.Errorf("object name cannot be '.' or '..'")
	}

	return nil
}

// ValidateStorageClass validates a GCS storage class
func ValidateStorageClass(class string) error {
	validClasses := []string{"STANDARD", "NEARLINE", "COLDLINE", "ARCHIVE"}

	if class == "" {
		return nil // Empty is valid, will use default
	}

	for _, valid := range validClasses {
		if class == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid storage class: %s. Valid classes are: %s",
		class, strings.Join(validClasses, ", "))
}

// ValidateLocation validates a GCS location
func ValidateLocation(location string) error {
	if location == "" {
		return fmt.Errorf("location cannot be empty")
	}

	// Basic validation - in a real implementation, you might want to
	// validate against a list of valid GCP regions/zones
	if len(location) < 2 {
		return fmt.Errorf("location must be at least 2 characters long")
	}

	return nil
}

// isIPAddress checks if a string is formatted as an IP address
func isIPAddress(s string) bool {
	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if !ipRegex.MatchString(s) {
		return false
	}

	// Additional validation for valid IP ranges (0-255)
	parts := strings.Split(s, ".")
	for _, part := range parts {
		if len(part) > 3 {
			return false
		}
		if part[0] == '0' && len(part) > 1 {
			return false
		}
		// Validate octet range (0-255)
		var octet int
		if _, err := fmt.Sscanf(part, "%d", &octet); err != nil {
			return false
		}
		if octet > 255 {
			return false
		}
	}

	return true
}

// WrapError wraps an error with additional context
func WrapError(operation, resource string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s %s: %w", operation, resource, err)
}
