package services

import (
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// GCPServiceInterface defines the interface for GCP operations
type GCPServiceInterface interface {
	// Project operations
	CreateProject(req *models.ProjectRequest) (*models.ProjectResponse, error)
	GetProject(projectID string) (*models.ProjectResponse, error)
	DeleteProject(projectID string) error

	// Folder operations
	CreateFolder(req *models.FolderRequest) (*models.FolderResponse, error)
	GetFolder(folderID string) (*models.FolderResponse, error)
	DeleteFolder(folderID string) error

	// Bucket operations
	CreateBucket(req *models.BucketRequest) (*models.BucketResponse, error)
	GetBucket(bucketName string) (*models.BucketResponse, error)
	DeleteBucket(bucketName string) error

	// Cleanup
	Close() error
}

// Ensure GCPService implements the interface
var _ GCPServiceInterface = (*GCPService)(nil)
