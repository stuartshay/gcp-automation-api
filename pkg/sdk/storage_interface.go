package sdk

import (
	"context"
	"io"

	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// StorageClient defines the interface for Cloud Storage operations
type StorageClient interface {
	// Bucket operations
	CreateBucket(ctx context.Context, req *models.BucketRequest) (*models.BucketResponse, error)
	GetBucket(ctx context.Context, bucketName string) (*models.BucketResponse, error)
	DeleteBucket(ctx context.Context, bucketName string) error
	ListBuckets(ctx context.Context, projectID string) ([]*models.BucketResponse, error)
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	UpdateBucket(ctx context.Context, bucketName string, req *models.BucketUpdateRequest) (*models.BucketResponse, error)

	// Object operations
	UploadObject(ctx context.Context, bucketName, objectName string, data io.Reader) (*models.ObjectResponse, error)
	DownloadObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, bucketName, objectName string) error
	ListObjects(ctx context.Context, bucketName string, prefix string) ([]*models.ObjectResponse, error)
	ObjectExists(ctx context.Context, bucketName, objectName string) (bool, error)
	GetObjectMetadata(ctx context.Context, bucketName, objectName string) (*models.ObjectResponse, error)

	// Lifecycle management
	SetBucketLifecycle(ctx context.Context, bucketName string, lifecycle *models.LifecyclePolicy) error
	GetBucketLifecycle(ctx context.Context, bucketName string) (*models.LifecyclePolicy, error)
	DeleteBucketLifecycle(ctx context.Context, bucketName string) error

	// Access control
	SetBucketIAM(ctx context.Context, bucketName string, policy *models.IAMPolicy) error
	GetBucketIAM(ctx context.Context, bucketName string) (*models.IAMPolicy, error)
	TestBucketIAM(ctx context.Context, bucketName string, permissions []string) ([]string, error)

	// Cleanup
	Close() error
}
