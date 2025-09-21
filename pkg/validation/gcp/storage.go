// Package gcp provides validation functions for Google Cloud Platform resources.
//
// This package includes validation for:
//   - Storage bucket names and object names
//   - GCP locations (regions and zones)
//   - Storage classes
//   - And other GCP-specific validation rules
//
// The package supports both static validation (fast, built-in rules) and
// dynamic validation (real-time validation against GCP APIs).
package gcp

import (
	"fmt"
	"net"
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

	// Cannot contain "google" in the name (since bucket names are already lowercase per regex)
	if strings.Contains(name, "google") {
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

// ValidateLocation validates a GCS location against known GCP regions and zones
func ValidateLocation(location string) error {
	if location == "" {
		return fmt.Errorf("location cannot be empty")
	}

	if len(location) < 2 {
		return fmt.Errorf("location must be at least 2 characters long")
	}

	// Validate against known GCP locations
	if isValidGCPLocation(location) {
		return nil
	}

	return fmt.Errorf("invalid GCP location: %s. Must be a valid GCP region (e.g., us-central1) or zone (e.g., us-central1-a)", location)
}

// isValidGCPLocation checks if the location is a valid GCP region or zone
func isValidGCPLocation(location string) bool {
	// Multi-regional locations
	multiRegional := map[string]bool{
		"us":   true,
		"eu":   true,
		"asia": true,
	}

	if multiRegional[location] {
		return true
	}

	// Common GCP regions and zones (as of 2024)
	validLocations := map[string]bool{
		// US regions
		"us-central1": true,
		"us-east1":    true,
		"us-east4":    true,
		"us-east5":    true,
		"us-south1":   true,
		"us-west1":    true,
		"us-west2":    true,
		"us-west3":    true,
		"us-west4":    true,

		// Europe regions
		"europe-central2":   true,
		"europe-north1":     true,
		"europe-southwest1": true,
		"europe-west1":      true,
		"europe-west2":      true,
		"europe-west3":      true,
		"europe-west4":      true,
		"europe-west6":      true,
		"europe-west8":      true,
		"europe-west9":      true,
		"europe-west10":     true,
		"europe-west12":     true,

		// Asia Pacific regions
		"asia-east1":              true,
		"asia-east2":              true,
		"asia-northeast1":         true,
		"asia-northeast2":         true,
		"asia-northeast3":         true,
		"asia-south1":             true,
		"asia-south2":             true,
		"asia-southeast1":         true,
		"asia-southeast2":         true,
		"australia-southeast1":    true,
		"australia-southeast2":    true,
		"northamerica-northeast1": true,
		"northamerica-northeast2": true,
		"southamerica-east1":      true,
		"southamerica-west1":      true,

		// Middle East and Africa
		"me-central1":   true,
		"me-central2":   true,
		"me-west1":      true,
		"africa-south1": true,
	}

	// Check if it's a known region
	if validLocations[location] {
		return true
	}

	// Check if it might be a zone (region + zone suffix like -a, -b, -c)
	return isValidZoneFormat(location, validLocations)
}

// isValidZoneFormat checks if the location follows the zone format (region-zone)
func isValidZoneFormat(location string, validRegions map[string]bool) bool {
	// Zone format: region-{a|b|c|d|e|f}
	if len(location) < 3 {
		return false
	}

	// Find the last dash and check if what follows is a valid zone suffix
	lastDash := strings.LastIndex(location, "-")
	if lastDash == -1 || lastDash == len(location)-1 {
		return false
	}

	region := location[:lastDash]
	zoneSuffix := location[lastDash+1:]

	// Valid zone suffixes
	validZoneSuffixes := map[string]bool{
		"a": true, "b": true, "c": true,
		"d": true, "e": true, "f": true,
	}

	return validRegions[region] && validZoneSuffixes[zoneSuffix]
}

// isIPAddress checks if a string is formatted as an IP address using Go's net package
func isIPAddress(s string) bool {
	// Use Go's built-in IP parsing which is more robust than regex
	return net.ParseIP(s) != nil
}

// WrapError wraps an error with additional context
func WrapError(operation, resource string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s %s: %w", operation, resource, err)
}
