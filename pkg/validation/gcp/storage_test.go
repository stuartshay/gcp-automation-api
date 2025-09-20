package gcp

import (
	"fmt"
	"testing"
)

func TestValidateBucketName(t *testing.T) {
	tests := []struct {
		name       string
		bucketName string
		wantError  bool
	}{
		{
			name:       "valid bucket name",
			bucketName: "valid-bucket-name",
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
			bucketName: "this-bucket-name-is-way-too-long-to-be-valid-according-to-gcs-rules",
			wantError:  true,
		},
		{
			name:       "bucket name with uppercase",
			bucketName: "Invalid-Bucket-Name",
			wantError:  true,
		},
		{
			name:       "bucket name starting with period",
			bucketName: ".invalid-bucket-name",
			wantError:  true,
		},
		{
			name:       "bucket name ending with period",
			bucketName: "invalid-bucket-name.",
			wantError:  true,
		},
		{
			name:       "bucket name starting with hyphen",
			bucketName: "-invalid-bucket-name",
			wantError:  true,
		},
		{
			name:       "bucket name ending with hyphen",
			bucketName: "invalid-bucket-name-",
			wantError:  true,
		},
		{
			name:       "bucket name with consecutive periods",
			bucketName: "invalid..bucket-name",
			wantError:  true,
		},
		{
			name:       "bucket name formatted as IP",
			bucketName: "192.168.1.1",
			wantError:  true,
		},
		{
			name:       "bucket name starting with goog",
			bucketName: "goog-bucket-name",
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

func TestValidateObjectName(t *testing.T) {
	tests := []struct {
		name       string
		objectName string
		wantError  bool
	}{
		{
			name:       "valid object name",
			objectName: "valid-object-name.txt",
			wantError:  false,
		},
		{
			name:       "empty object name",
			objectName: "",
			wantError:  true,
		},
		{
			name:       "too long object name",
			objectName: string(make([]byte, 1025)),
			wantError:  true,
		},
		{
			name:       "object name with newline",
			objectName: "invalid\nobject-name",
			wantError:  true,
		},
		{
			name:       "object name with carriage return",
			objectName: "invalid\robject-name",
			wantError:  true,
		},
		{
			name:       "object name with null character",
			objectName: "invalid\x00object-name",
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
			wantError:    false,
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
		// Valid regions
		{
			name:      "valid US region",
			location:  "us-central1",
			wantError: false,
		},
		{
			name:      "valid Europe region",
			location:  "europe-west1",
			wantError: false,
		},
		{
			name:      "valid Asia region",
			location:  "asia-east1",
			wantError: false,
		},
		{
			name:      "valid multi-regional",
			location:  "us",
			wantError: false,
		},
		// Valid zones
		{
			name:      "valid US zone",
			location:  "us-central1-a",
			wantError: false,
		},
		{
			name:      "valid Europe zone",
			location:  "europe-west1-b",
			wantError: false,
		},
		{
			name:      "valid Asia zone",
			location:  "asia-east1-c",
			wantError: false,
		},
		// Invalid cases
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
		{
			name:      "invalid region",
			location:  "invalid-region",
			wantError: true,
		},
		{
			name:      "invalid zone format",
			location:  "us-central1-z",
			wantError: true,
		},
		{
			name:      "invalid zone base",
			location:  "invalid-region-a",
			wantError: true,
		},
		{
			name:      "malformed location",
			location:  "us--central1",
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
	tests := []struct {
		name      string
		operation string
		resource  string
		err       error
		want      string
		wantNil   bool
	}{
		{
			name:      "nil error",
			operation: "creating",
			resource:  "bucket",
			err:       nil,
			wantNil:   true,
		},
		{
			name:      "non-nil error",
			operation: "creating",
			resource:  "bucket",
			err:       fmt.Errorf("original error"),
			want:      "creating bucket: original error",
			wantNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.operation, tt.resource, tt.err)
			if tt.wantNil {
				if result != nil {
					t.Errorf("WrapError() = %v, want nil", result)
				}
			} else {
				if result == nil {
					t.Errorf("WrapError() = nil, want error")
				} else if result.Error() != tt.want {
					t.Errorf("WrapError() = %v, want %v", result.Error(), tt.want)
				}
			}
		})
	}
}

// Benchmark the static validation performance
func BenchmarkValidateLocation(b *testing.B) {
	location := "us-central1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateLocation(location)
	}
}

// Benchmark the zone validation performance
func BenchmarkValidateLocationZone(b *testing.B) {
	location := "us-central1-a"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateLocation(location)
	}
}

// Benchmark bucket name validation
func BenchmarkValidateBucketName(b *testing.B) {
	bucketName := "valid-bucket-name"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateBucketName(bucketName)
	}
}
