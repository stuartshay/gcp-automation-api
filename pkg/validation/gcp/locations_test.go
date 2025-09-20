package gcp

import (
	"context"
	"testing"
	"time"
)

func TestLocationValidator_ValidateLocationDynamic(t *testing.T) {
	// Skip this test if we don't have real GCP credentials
	// This is an integration test that requires actual GCP API access
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	projectID := "test-project" // This would need to be a real project ID
	validator := NewLocationValidator(projectID)

	tests := []struct {
		name      string
		location  string
		wantError bool
	}{
		{
			name:      "empty location",
			location:  "",
			wantError: true,
		},
		// Note: Real tests would require valid GCP credentials and project
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateLocationDynamic(ctx, tt.location)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateLocationDynamic() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestLocationValidator_Cache(t *testing.T) {
	validator := NewLocationValidator("test-project")

	// Test that cache TTL is set correctly
	if validator.cacheTTL != time.Hour {
		t.Errorf("Expected cache TTL to be 1 hour, got %v", validator.cacheTTL)
	}

	// Test that initial cache is empty
	if len(validator.regions) != 0 {
		t.Errorf("Expected empty regions cache, got %d items", len(validator.regions))
	}

	if len(validator.zones) != 0 {
		t.Errorf("Expected empty zones cache, got %d items", len(validator.zones))
	}
}

func TestValidateLocationWithFallback(t *testing.T) {
	ctx := context.Background()
	projectID := "test-project"

	tests := []struct {
		name      string
		location  string
		wantError bool
	}{
		{
			name:      "valid static location",
			location:  "us-central1",
			wantError: false, // Should pass static validation
		},
		{
			name:      "empty location",
			location:  "",
			wantError: true,
		},
		{
			name:      "invalid location",
			location:  "invalid-location",
			wantError: true, // Will fail both static and dynamic (since we don't have real creds)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLocationWithFallback(ctx, projectID, tt.location)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateLocationWithFallback() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
