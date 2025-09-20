package sdk

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stuartshay/gcp-automation-api/internal/models"
)

func TestValidateBucketName(t *testing.T) {
	tests := []struct {
		name       string
		bucketName string
		wantError  bool
	}{
		{
			name:       "valid bucket name",
			bucketName: "my-test-bucket-123",
			wantError:  false,
		},
		{
			name:       "empty bucket name",
			bucketName: "",
			wantError:  true,
		},
		{
			name:       "too short bucket name",
			bucketName: "ab",
			wantError:  true,
		},
		{
			name:       "too long bucket name",
			bucketName: strings.Repeat("a", 64),
			wantError:  true,
		},
		{
			name:       "bucket name with uppercase",
			bucketName: "My-Test-Bucket",
			wantError:  true,
		},
		{
			name:       "bucket name starting with period",
			bucketName: ".my-bucket",
			wantError:  true,
		},
		{
			name:       "bucket name ending with period",
			bucketName: "my-bucket.",
			wantError:  true,
		},
		{
			name:       "bucket name starting with hyphen",
			bucketName: "-my-bucket",
			wantError:  true,
		},
		{
			name:       "bucket name ending with hyphen",
			bucketName: "my-bucket-",
			wantError:  true,
		},
		{
			name:       "bucket name with consecutive periods",
			bucketName: "my..bucket",
			wantError:  true,
		},
		{
			name:       "bucket name formatted as IP",
			bucketName: "192.168.1.1",
			wantError:  true,
		},
		{
			name:       "bucket name starting with goog",
			bucketName: "goog-bucket",
			wantError:  true,
		},
		{
			name:       "bucket name containing google",
			bucketName: "my-google-bucket",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBucketName(tt.bucketName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateBucketName() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestIsIPAddress(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantIP bool
	}{
		{
			name:   "valid IP address",
			input:  "192.168.1.1",
			wantIP: true,
		},
		{
			name:   "invalid IP with out of range octets",
			input:  "999.999.999.999",
			wantIP: false,
		},
		{
			name:   "valid IP edge case",
			input:  "255.255.255.255",
			wantIP: true,
		},
		{
			name:   "invalid IP with leading zeros",
			input:  "192.168.01.1",
			wantIP: false,
		},
		{
			name:   "not an IP address",
			input:  "my-bucket-name",
			wantIP: false,
		},
		{
			name:   "IP with invalid format",
			input:  "192.168.1",
			wantIP: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIPAddress(tt.input)
			if result != tt.wantIP {
				t.Errorf("isIPAddress(%s) = %v, want %v", tt.input, result, tt.wantIP)
			}
		})
	}
}

func TestValidateObjectName(t *testing.T) {
	tests := []struct {
		name       string
		objectName string
		wantError  bool
	}{
		{
			name:       "valid object name",
			objectName: "path/to/my-file.txt",
			wantError:  false,
		},
		{
			name:       "empty object name",
			objectName: "",
			wantError:  true,
		},
		{
			name:       "too long object name",
			objectName: strings.Repeat("a", 1025),
			wantError:  true,
		},
		{
			name:       "object name with newline",
			objectName: "file\nwith\nnewline",
			wantError:  true,
		},
		{
			name:       "object name with carriage return",
			objectName: "file\rwith\rcarriage",
			wantError:  true,
		},
		{
			name:       "object name with null character",
			objectName: "file\x00with\x00null",
			wantError:  true,
		},
		{
			name:       "object name as single period",
			objectName: ".",
			wantError:  true,
		},
		{
			name:       "object name as double period",
			objectName: "..",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateObjectName(tt.objectName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateObjectName() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateStorageClass(t *testing.T) {
	tests := []struct {
		name         string
		storageClass string
		wantError    bool
	}{
		{
			name:         "valid storage class STANDARD",
			storageClass: "STANDARD",
			wantError:    false,
		},
		{
			name:         "valid storage class NEARLINE",
			storageClass: "NEARLINE",
			wantError:    false,
		},
		{
			name:         "valid storage class COLDLINE",
			storageClass: "COLDLINE",
			wantError:    false,
		},
		{
			name:         "valid storage class ARCHIVE",
			storageClass: "ARCHIVE",
			wantError:    false,
		},
		{
			name:         "empty storage class",
			storageClass: "",
			wantError:    false, // Empty is valid
		},
		{
			name:         "invalid storage class",
			storageClass: "INVALID",
			wantError:    true,
		},
		{
			name:         "lowercase storage class",
			storageClass: "standard",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStorageClass(tt.storageClass)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStorageClass() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateLocation(t *testing.T) {
	tests := []struct {
		name      string
		location  string
		wantError bool
	}{
		{
			name:      "valid location",
			location:  "us-central1",
			wantError: false,
		},
		{
			name:      "empty location",
			location:  "",
			wantError: true,
		},
		{
			name:      "too short location",
			location:  "a",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLocation(tt.location)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateLocation() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	wrappedErr := WrapError("creating", "bucket-name", originalErr)

	if wrappedErr == nil {
		t.Error("WrapError() should not return nil for non-nil error")
	}

	expectedPrefix := "creating bucket-name:"
	if !strings.Contains(wrappedErr.Error(), expectedPrefix) {
		t.Errorf("WrapError() error = %v, should contain %s", wrappedErr, expectedPrefix)
	}

	// Test with nil error
	nilWrapped := WrapError("creating", "bucket-name", nil)
	if nilWrapped != nil {
		t.Error("WrapError() should return nil for nil error")
	}
}

// Mock tests for GCPStorageClient methods
// Note: These tests would typically use a mock GCS client for full testing

func TestGCPStorageClient_CreateBucket_Validation(t *testing.T) {
	// This test focuses on validation logic without requiring actual GCP credentials
	// We'll test the validation logic by checking what type of error we get

	tests := []struct {
		name              string
		request           *models.BucketRequest
		wantValidationErr bool
	}{
		{
			name: "valid request",
			request: &models.BucketRequest{
				Name:         "valid-bucket-name",
				Location:     "us-central1",
				StorageClass: "STANDARD",
			},
			wantValidationErr: false,
		},
		{
			name: "invalid bucket name",
			request: &models.BucketRequest{
				Name:         "Invalid-Bucket-Name",
				Location:     "us-central1",
				StorageClass: "STANDARD",
			},
			wantValidationErr: true,
		},
		{
			name: "invalid location",
			request: &models.BucketRequest{
				Name:         "valid-bucket-name",
				Location:     "",
				StorageClass: "STANDARD",
			},
			wantValidationErr: true,
		},
		{
			name: "invalid storage class",
			request: &models.BucketRequest{
				Name:         "valid-bucket-name",
				Location:     "us-central1",
				StorageClass: "INVALID",
			},
			wantValidationErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation functions directly instead of the full client
			bucketErr := ValidateBucketName(tt.request.Name)
			locationErr := ValidateLocation(tt.request.Location)
			storageErr := ValidateStorageClass(tt.request.StorageClass)

			hasValidationErr := bucketErr != nil || locationErr != nil || storageErr != nil

			if hasValidationErr != tt.wantValidationErr {
				t.Errorf("Validation error = %v, wantValidationErr %v", hasValidationErr, tt.wantValidationErr)
			}
		})
	}
}

// Note: The following tests would require actual GCP credentials and connections
// They are commented out for unit testing purposes, but can be enabled for integration testing

/*
func TestGCPStorageClient_GetBucket_Validation(t *testing.T) {
	// Integration test - requires actual GCP connection
}

func TestGCPStorageClient_DeleteBucket_Validation(t *testing.T) {
	// Integration test - requires actual GCP connection
}

func TestGCPStorageClient_ObjectOperations_Validation(t *testing.T) {
	// Integration test - requires actual GCP connection
}
*/

// Benchmark tests for validation functions
func BenchmarkValidateBucketName(b *testing.B) {
	bucketName := "my-test-bucket-with-a-reasonable-length-name"

	for i := 0; i < b.N; i++ {
		_ = ValidateBucketName(bucketName)
	}
}

func BenchmarkValidateObjectName(b *testing.B) {
	objectName := "path/to/my/object/file.txt"

	for i := 0; i < b.N; i++ {
		_ = ValidateObjectName(objectName)
	}
}

func BenchmarkValidateStorageClass(b *testing.B) {
	storageClass := "STANDARD"

	for i := 0; i < b.N; i++ {
		_ = ValidateStorageClass(storageClass)
	}
}

// Note: Helper functions for creating test requests can be added here when needed
// for integration tests that require actual GCP connections
