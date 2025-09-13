package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

func TestProjectRequestValidation(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		request   models.ProjectRequest
		expectErr bool
		errMsg    string
	}{
		{
			name: "Valid project request",
			request: models.ProjectRequest{
				ProjectID:   "my-test-project-123",
				DisplayName: "My Test Project",
				ParentID:    "123456789",
				ParentType:  "organization",
				Labels: map[string]string{
					"environment": "test",
					"team":        "platform",
				},
			},
			expectErr: false,
		},
		{
			name: "Invalid project ID - too short",
			request: models.ProjectRequest{
				ProjectID:   "test",
				DisplayName: "My Test Project",
			},
			expectErr: true,
			errMsg:    "project_id must be a valid GCP project ID",
		},
		{
			name: "Invalid project ID - uppercase",
			request: models.ProjectRequest{
				ProjectID:   "My-Test-Project",
				DisplayName: "My Test Project",
			},
			expectErr: true,
			errMsg:    "project_id must be a valid GCP project ID",
		},
		{
			name: "Invalid project ID - ends with hyphen",
			request: models.ProjectRequest{
				ProjectID:   "my-test-project-",
				DisplayName: "My Test Project",
			},
			expectErr: true,
			errMsg:    "project_id must be a valid GCP project ID",
		},
		{
			name: "Missing required fields",
			request: models.ProjectRequest{
				Labels: map[string]string{
					"environment": "test",
				},
			},
			expectErr: true,
			errMsg:    "project_id is required",
		},
		{
			name: "Invalid parent type",
			request: models.ProjectRequest{
				ProjectID:   "my-test-project-123",
				DisplayName: "My Test Project",
				ParentType:  "invalid",
			},
			expectErr: true,
			errMsg:    "parent_type must be one of: organization folder",
		},
		{
			name: "Invalid label key",
			request: models.ProjectRequest{
				ProjectID:   "my-test-project-123",
				DisplayName: "My Test Project",
				Labels: map[string]string{
					"Environment": "test", // Capital letter not allowed
				},
			},
			expectErr: true,
			errMsg:    "label_key must be a valid GCP label key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(&tt.request)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBucketRequestValidation(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		request   models.BucketRequest
		expectErr bool
		errMsg    string
	}{
		{
			name: "Valid bucket request",
			request: models.BucketRequest{
				Name:         "my-test-bucket-123",
				Location:     "us-central1",
				StorageClass: "STANDARD",
				Labels: map[string]string{
					"environment": "test",
				},
				Versioning: true,
			},
			expectErr: false,
		},
		{
			name: "Invalid bucket name - too short",
			request: models.BucketRequest{
				Name:     "ab",
				Location: "us-central1",
			},
			expectErr: true,
			errMsg:    "name must be a valid GCS bucket name",
		},
		{
			name: "Invalid bucket name - uppercase",
			request: models.BucketRequest{
				Name:     "My-Test-Bucket",
				Location: "us-central1",
			},
			expectErr: true,
			errMsg:    "name must be a valid GCS bucket name",
		},
		{
			name: "Invalid bucket name - consecutive dots",
			request: models.BucketRequest{
				Name:     "my..test.bucket",
				Location: "us-central1",
			},
			expectErr: true,
			errMsg:    "name must be a valid GCS bucket name",
		},
		{
			name: "Invalid bucket name - IP address format",
			request: models.BucketRequest{
				Name:     "192.168.1.1",
				Location: "us-central1",
			},
			expectErr: true,
			errMsg:    "name must be a valid GCS bucket name",
		},
		{
			name: "Invalid location",
			request: models.BucketRequest{
				Name:     "my-test-bucket",
				Location: "invalid-location",
			},
			expectErr: true,
			errMsg:    "location must be a valid GCP location/region",
		},
		{
			name: "Invalid storage class",
			request: models.BucketRequest{
				Name:         "my-test-bucket",
				Location:     "us-central1",
				StorageClass: "INVALID",
			},
			expectErr: true,
			errMsg:    "storage_class must be one of: STANDARD NEARLINE COLDLINE ARCHIVE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(&tt.request)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFolderRequestValidation(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		request   models.FolderRequest
		expectErr bool
		errMsg    string
	}{
		{
			name: "Valid folder request",
			request: models.FolderRequest{
				DisplayName: "My Test Folder",
				ParentID:    "123456789",
				ParentType:  "organization",
			},
			expectErr: false,
		},
		{
			name: "Missing display name",
			request: models.FolderRequest{
				ParentID:   "123456789",
				ParentType: "organization",
			},
			expectErr: true,
			errMsg:    "display_name is required",
		},
		{
			name: "Invalid parent ID - non-numeric",
			request: models.FolderRequest{
				DisplayName: "My Test Folder",
				ParentID:    "abc123",
				ParentType:  "organization",
			},
			expectErr: true,
			errMsg:    "parent_id must be numeric",
		},
		{
			name: "Invalid parent type",
			request: models.FolderRequest{
				DisplayName: "My Test Folder",
				ParentID:    "123456789",
				ParentType:  "project",
			},
			expectErr: true,
			errMsg:    "parent_type must be one of: organization folder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(&tt.request)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
