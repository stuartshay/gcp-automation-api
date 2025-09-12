package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// ValidateProjectID validates a GCP project ID
func ValidateProjectID(projectID string) error {
	// GCP project ID requirements:
	// - Must be 6 to 30 lowercase letters, digits, or hyphens
	// - Must start with a letter
	// - Cannot end with a hyphen
	pattern := `^[a-z][a-z0-9-]{4,28}[a-z0-9]$`

	if matched, _ := regexp.MatchString(pattern, projectID); !matched {
		return fmt.Errorf("invalid project ID: must be 6-30 characters, start with letter, contain only lowercase letters, digits, and hyphens, and not end with hyphen")
	}

	return nil
}

// ValidateBucketName validates a GCS bucket name
func ValidateBucketName(bucketName string) error {
	// GCS bucket name requirements:
	// - Must be 3 to 63 characters
	// - Must start and end with a number or letter
	// - Can contain lowercase letters, numbers, and hyphens
	pattern := `^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$`

	if len(bucketName) < 3 || len(bucketName) > 63 {
		return fmt.Errorf("invalid bucket name: must be 3-63 characters")
	}

	if matched, _ := regexp.MatchString(pattern, bucketName); !matched {
		return fmt.Errorf("invalid bucket name: must contain only lowercase letters, numbers, and hyphens, and start/end with letter or number")
	}

	return nil
}

// PrettyPrintJSON returns a pretty-printed JSON string
func PrettyPrintJSON(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ContainsString checks if a string slice contains a string
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// RemoveString removes a string from a slice
func RemoveString(slice []string, str string) []string {
	var result []string
	for _, s := range slice {
		if s != str {
			result = append(result, s)
		}
	}
	return result
}
