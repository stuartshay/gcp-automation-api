package validators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance with custom rules
func NewValidator() *CustomValidator {
	v := validator.New()

	// Register custom validation functions
	_ = v.RegisterValidation("project_id", validateProjectID)
	_ = v.RegisterValidation("bucket_name", validateBucketName)
	_ = v.RegisterValidation("label_key", validateLabelKey)
	_ = v.RegisterValidation("label_value", validateLabelValue)
	_ = v.RegisterValidation("gcp_location", validateGCPLocation)

	return &CustomValidator{validator: v}
}

// Validate validates a struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Convert validation errors to more user-friendly format
		return formatValidationError(err)
	}
	return nil
}

// validateProjectID validates GCP project ID format
func validateProjectID(fl validator.FieldLevel) bool {
	projectID := fl.Field().String()

	// GCP Project ID rules:
	// - 6-30 characters
	// - Lowercase letters, digits, and hyphens only
	// - Must start with a letter
	// - Cannot end with a hyphen
	if len(projectID) < 6 || len(projectID) > 30 {
		return false
	}

	projectIDRegex := regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)
	return projectIDRegex.MatchString(projectID)
}

// validateBucketName validates GCS bucket name format
func validateBucketName(fl validator.FieldLevel) bool {
	bucketName := fl.Field().String()

	// GCS Bucket name rules:
	// - 3-63 characters
	// - Lowercase letters, digits, dashes, underscores, and dots
	// - Must start and end with alphanumeric
	// - Cannot contain consecutive dots
	// - Cannot be formatted as IP address
	if len(bucketName) < 3 || len(bucketName) > 63 {
		return false
	}

	// Basic format check
	bucketNameRegex := regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*[a-z0-9]$`)
	if !bucketNameRegex.MatchString(bucketName) {
		return false
	}

	// Cannot contain consecutive dots
	if strings.Contains(bucketName, "..") {
		return false
	}

	// Cannot be formatted as IP address
	ipRegex := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)
	return !ipRegex.MatchString(bucketName)
}

// validateLabelKey validates GCP label keys
func validateLabelKey(fl validator.FieldLevel) bool {
	key := fl.Field().String()

	// GCP Label key rules:
	// - 1-63 characters
	// - Lowercase letters, digits, underscores, and dashes
	// - Must start with lowercase letter
	if len(key) == 0 || len(key) > 63 {
		return false
	}

	labelKeyRegex := regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
	return labelKeyRegex.MatchString(key)
}

// validateLabelValue validates GCP label values
func validateLabelValue(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// GCP Label value rules:
	// - 0-63 characters
	// - Lowercase letters, digits, underscores, and dashes
	if len(value) > 63 {
		return false
	}

	if len(value) == 0 {
		return true // Empty values are allowed
	}

	labelValueRegex := regexp.MustCompile(`^[a-z0-9_-]*$`)
	return labelValueRegex.MatchString(value)
}

// validateGCPLocation validates GCP location/region format
func validateGCPLocation(fl validator.FieldLevel) bool {
	location := fl.Field().String()

	// Common GCP locations/regions
	validLocations := map[string]bool{
		// Multi-regional
		"us":   true,
		"eu":   true,
		"asia": true,

		// Regional - US
		"us-central1": true,
		"us-east1":    true,
		"us-east4":    true,
		"us-west1":    true,
		"us-west2":    true,
		"us-west3":    true,
		"us-west4":    true,

		// Regional - Europe
		"europe-north1":   true,
		"europe-west1":    true,
		"europe-west2":    true,
		"europe-west3":    true,
		"europe-west4":    true,
		"europe-west6":    true,
		"europe-central2": true,

		// Regional - Asia Pacific
		"asia-east1":      true,
		"asia-east2":      true,
		"asia-northeast1": true,
		"asia-northeast2": true,
		"asia-northeast3": true,
		"asia-south1":     true,
		"asia-south2":     true,
		"asia-southeast1": true,
		"asia-southeast2": true,

		// Regional - Other
		"australia-southeast1":    true,
		"southamerica-east1":      true,
		"northamerica-northeast1": true,
	}

	return validLocations[strings.ToLower(location)]
}

// formatValidationError converts validator errors to user-friendly messages
func formatValidationError(err error) error {
	validationErrors := err.(validator.ValidationErrors)

	var messages []string
	for _, e := range validationErrors {
		message := getFieldErrorMessage(e)
		messages = append(messages, message)
	}

	return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
}

// getFieldErrorMessage returns user-friendly error message for a field
func getFieldErrorMessage(fe validator.FieldError) string {
	// Handle map key validation specially
	field := fe.Field()
	if strings.Contains(field, "[") && strings.Contains(field, "]") {
		// Extract the base field name (e.g., "Labels" from "Labels[environment]")
		baseField := strings.Split(field, "[")[0]
		field = convertFieldName(baseField)
	} else {
		field = convertFieldName(field)
	}

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "project_id":
		return fmt.Sprintf("%s must be a valid GCP project ID (6-30 chars, lowercase letters/digits/hyphens, start with letter, not end with hyphen)", field)
	case "bucket_name":
		return fmt.Sprintf("%s must be a valid GCS bucket name (3-63 chars, lowercase letters/digits/dashes/underscores/dots, no consecutive dots, not an IP address)", field)
	case "label_key":
		return "label_key must be a valid GCP label key (1-63 chars, start with lowercase letter, lowercase letters/digits/underscores/dashes only)"
	case "label_value":
		return "label_value must be a valid GCP label value (0-63 chars, lowercase letters/digits/underscores/dashes only)"
	case "gcp_location":
		return fmt.Sprintf("%s must be a valid GCP location/region", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// convertFieldName converts struct field names to snake_case for user-friendly messages
func convertFieldName(field string) string {
	switch field {
	case "ProjectID":
		return "project_id"
	case "DisplayName":
		return "display_name"
	case "ParentID":
		return "parent_id"
	case "ParentType":
		return "parent_type"
	case "StorageClass":
		return "storage_class"
	case "Name":
		return "name"
	case "Location":
		return "location"
	case "Labels":
		return "labels"
	case "Versioning":
		return "versioning"
	default:
		return strings.ToLower(field)
	}
}

// ValidationError represents a validation error response
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ParseValidationErrors parses validation errors into structured format
func ParseValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErr {
			validationErrors = append(validationErrors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: getFieldErrorMessage(e),
			})
		}
	}

	return validationErrors
}
