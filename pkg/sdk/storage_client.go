package sdk

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/pkg/validation/gcp"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCPStorageClient implements the StorageClient interface for Google Cloud Storage
type GCPStorageClient struct {
	client    *storage.Client
	projectID string
	ctx       context.Context
}

// NewGCPStorageClient creates a new GCP Storage Client
func NewGCPStorageClient(ctx context.Context, projectID string, opts ...option.ClientOption) (*GCPStorageClient, error) {
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &GCPStorageClient{
		client:    client,
		projectID: projectID,
		ctx:       ctx,
	}, nil
}

// CreateBucket creates a new GCS bucket
func (c *GCPStorageClient) CreateBucket(ctx context.Context, req *models.BucketRequest) (*models.BucketResponse, error) {
	// Validate request
	if err := gcp.ValidateBucketName(req.Name); err != nil {
		return nil, gcp.WrapError("creating bucket", req.Name, err)
	}

	if err := gcp.ValidateLocation(req.Location); err != nil {
		return nil, gcp.WrapError("creating bucket", req.Name, err)
	}

	if err := gcp.ValidateStorageClass(req.StorageClass); err != nil {
		return nil, gcp.WrapError("creating bucket", req.Name, err)
	}

	bucket := c.client.Bucket(req.Name)

	// Set basic bucket attributes
	attrs := &storage.BucketAttrs{
		Location: req.Location,
		Labels:   req.Labels,
	}

	if req.StorageClass != "" {
		attrs.StorageClass = req.StorageClass
	}

	if req.Versioning {
		attrs.VersioningEnabled = true
	}

	// Phase 1 Advanced Options - Security & Compliance

	// KMS Encryption
	if req.KMSKeyName != "" {
		attrs.Encryption = &storage.BucketEncryption{
			DefaultKMSKeyName: req.KMSKeyName,
		}
	}

	// Retention Policy
	if req.RetentionPolicy != nil {
		attrs.RetentionPolicy = &storage.RetentionPolicy{
			RetentionPeriod: time.Duration(req.RetentionPolicy.RetentionPeriodSeconds) * time.Second,
			IsLocked:        req.RetentionPolicy.IsLocked,
		}
	}

	// Uniform Bucket-Level Access
	if req.UniformBucketLevelAccess {
		attrs.UniformBucketLevelAccess = storage.UniformBucketLevelAccess{
			Enabled: true,
		}
	}

	// Public Access Prevention
	if req.PublicAccessPrevention != "" {
		switch req.PublicAccessPrevention {
		case "enforced":
			attrs.PublicAccessPrevention = storage.PublicAccessPreventionEnforced
		case "inherited":
			attrs.PublicAccessPrevention = storage.PublicAccessPreventionInherited
		case "unspecified":
			attrs.PublicAccessPrevention = storage.PublicAccessPreventionUnspecified
		}
	}

	// Create the bucket
	if err := bucket.Create(ctx, c.projectID, attrs); err != nil {
		return nil, gcp.WrapError("creating bucket", req.Name, err)
	}

	// Get bucket attributes to return complete information
	bucketAttrs, err := bucket.Attrs(ctx)
	if err != nil {
		return nil, gcp.WrapError("getting bucket attributes after creation", req.Name, err)
	}

	return c.mapBucketAttrsToResponse(bucketAttrs), nil
}

// GetBucket retrieves a GCS bucket
func (c *GCPStorageClient) GetBucket(ctx context.Context, bucketName string) (*models.BucketResponse, error) {
	if err := gcp.ValidateBucketName(bucketName); err != nil {
		return nil, gcp.WrapError("getting bucket", bucketName, err)
	}

	bucket := c.client.Bucket(bucketName)

	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		return nil, gcp.WrapError("getting bucket", bucketName, err)
	}

	return c.mapBucketAttrsToResponse(attrs), nil
}

// DeleteBucket deletes a GCS bucket
func (c *GCPStorageClient) DeleteBucket(ctx context.Context, bucketName string) error {
	if err := gcp.ValidateBucketName(bucketName); err != nil {
		return gcp.WrapError("deleting bucket", bucketName, err)
	}

	bucket := c.client.Bucket(bucketName)

	if err := bucket.Delete(ctx); err != nil {
		return gcp.WrapError("deleting bucket", bucketName, err)
	}

	return nil
}

// ListBuckets lists all buckets in the project
func (c *GCPStorageClient) ListBuckets(ctx context.Context, projectID string) ([]*models.BucketResponse, error) {
	if projectID == "" {
		projectID = c.projectID
	}

	var buckets []*models.BucketResponse
	it := c.client.Buckets(ctx, projectID)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list buckets: %w", err)
		}

		buckets = append(buckets, c.mapBucketAttrsToResponse(attrs))
	}

	return buckets, nil
}

// BucketExists checks if a bucket exists
func (c *GCPStorageClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	bucket := c.client.Bucket(bucketName)
	_, err := bucket.Attrs(ctx)
	if err != nil {
		if err == storage.ErrBucketNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	return true, nil
}

// UpdateBucket updates a GCS bucket (simplified version)
func (c *GCPStorageClient) UpdateBucket(ctx context.Context, bucketName string, req *models.BucketUpdateRequest) (*models.BucketResponse, error) {
	bucket := c.client.Bucket(bucketName)

	// Get current attributes first
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current bucket attributes: %w", err)
	}

	// For now, only support versioning updates
	if req.Versioning != nil {
		attrsToUpdate := storage.BucketAttrsToUpdate{
			VersioningEnabled: req.Versioning,
		}

		attrs, err = bucket.Update(ctx, attrsToUpdate)
		if err != nil {
			return nil, fmt.Errorf("failed to update bucket: %w", err)
		}
	}

	return c.mapBucketAttrsToResponse(attrs), nil
}

// UploadObject uploads an object to a bucket
func (c *GCPStorageClient) UploadObject(ctx context.Context, bucketName, objectName string, data io.Reader) (*models.ObjectResponse, error) {
	if err := gcp.ValidateBucketName(bucketName); err != nil {
		return nil, gcp.WrapError("uploading object", bucketName+"/"+objectName, err)
	}

	if err := gcp.ValidateObjectName(objectName); err != nil {
		return nil, gcp.WrapError("uploading object", bucketName+"/"+objectName, err)
	}

	bucket := c.client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	writer := obj.NewWriter(ctx)

	if _, err := io.Copy(writer, data); err != nil {
		return nil, gcp.WrapError("uploading object", bucketName+"/"+objectName, err)
	}

	if err := writer.Close(); err != nil {
		return nil, gcp.WrapError("closing object writer after upload", bucketName+"/"+objectName, err)
	}

	// Get object attributes
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, gcp.WrapError("getting object attributes after upload", bucketName+"/"+objectName, err)
	}

	return c.mapObjectAttrsToResponse(attrs), nil
}

// DownloadObject downloads an object from a bucket
func (c *GCPStorageClient) DownloadObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	if err := gcp.ValidateBucketName(bucketName); err != nil {
		return nil, gcp.WrapError("downloading object", bucketName+"/"+objectName, err)
	}

	if err := gcp.ValidateObjectName(objectName); err != nil {
		return nil, gcp.WrapError("downloading object", bucketName+"/"+objectName, err)
	}

	bucket := c.client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, gcp.WrapError("downloading object", bucketName+"/"+objectName, err)
	}

	return reader, nil
}

// DeleteObject deletes an object from a bucket
func (c *GCPStorageClient) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	if err := gcp.ValidateBucketName(bucketName); err != nil {
		return gcp.WrapError("deleting object", bucketName+"/"+objectName, err)
	}

	if err := gcp.ValidateObjectName(objectName); err != nil {
		return gcp.WrapError("deleting object", bucketName+"/"+objectName, err)
	}

	bucket := c.client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	if err := obj.Delete(ctx); err != nil {
		return gcp.WrapError("deleting object", bucketName+"/"+objectName, err)
	}

	return nil
}

// ListObjects lists objects in a bucket
func (c *GCPStorageClient) ListObjects(ctx context.Context, bucketName string, prefix string) ([]*models.ObjectResponse, error) {
	bucket := c.client.Bucket(bucketName)

	query := &storage.Query{Prefix: prefix}
	it := bucket.Objects(ctx, query)

	var objects []*models.ObjectResponse
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		objects = append(objects, c.mapObjectAttrsToResponse(attrs))
	}

	return objects, nil
}

// ObjectExists checks if an object exists
func (c *GCPStorageClient) ObjectExists(ctx context.Context, bucketName, objectName string) (bool, error) {
	bucket := c.client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	_, err := obj.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}

// GetObjectMetadata retrieves object metadata
func (c *GCPStorageClient) GetObjectMetadata(ctx context.Context, bucketName, objectName string) (*models.ObjectResponse, error) {
	bucket := c.client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}

	return c.mapObjectAttrsToResponse(attrs), nil
}

// SetBucketLifecycle sets the lifecycle policy for a bucket (simplified implementation)
func (c *GCPStorageClient) SetBucketLifecycle(ctx context.Context, bucketName string, lifecycle *models.LifecyclePolicy) error {
	// Simplified implementation - just return not implemented for now
	return fmt.Errorf("lifecycle policy management not implemented yet")
}

// GetBucketLifecycle gets the lifecycle policy for a bucket (simplified implementation)
func (c *GCPStorageClient) GetBucketLifecycle(ctx context.Context, bucketName string) (*models.LifecyclePolicy, error) {
	// Return empty lifecycle policy for now
	return &models.LifecyclePolicy{Rules: []models.LifecycleRule{}}, nil
}

// DeleteBucketLifecycle deletes the lifecycle policy for a bucket (simplified implementation)
func (c *GCPStorageClient) DeleteBucketLifecycle(ctx context.Context, bucketName string) error {
	// Simplified implementation - just return not implemented for now
	return fmt.Errorf("lifecycle policy management not implemented yet")
}

// SetBucketIAM sets the IAM policy for a bucket (simplified implementation)
func (c *GCPStorageClient) SetBucketIAM(ctx context.Context, bucketName string, policy *models.IAMPolicy) error {
	// Simplified implementation - just return not implemented for now
	return fmt.Errorf("IAM policy management not implemented yet")
}

// GetBucketIAM gets the IAM policy for a bucket (simplified implementation)
func (c *GCPStorageClient) GetBucketIAM(ctx context.Context, bucketName string) (*models.IAMPolicy, error) {
	// Return empty IAM policy for now
	return &models.IAMPolicy{
		Bindings: []models.IAMBinding{},
		Etag:     "",
		Version:  1,
	}, nil
}

// TestBucketIAM tests IAM permissions for a bucket (simplified implementation)
func (c *GCPStorageClient) TestBucketIAM(ctx context.Context, bucketName string, permissions []string) ([]string, error) {
	bucket := c.client.Bucket(bucketName)
	handle := bucket.IAM()

	perms, err := handle.TestPermissions(ctx, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to test bucket IAM permissions: %w", err)
	}

	return perms, nil
}

// Close closes the storage client
func (c *GCPStorageClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// Helper methods for mapping GCP types to our models

func (c *GCPStorageClient) mapBucketAttrsToResponse(attrs *storage.BucketAttrs) *models.BucketResponse {
	response := &models.BucketResponse{
		Name:         attrs.Name,
		Location:     attrs.Location,
		StorageClass: attrs.StorageClass,
		Labels:       attrs.Labels,
		Versioning:   attrs.VersioningEnabled,
		CreateTime:   attrs.Created,
		UpdateTime:   attrs.Updated,
		SelfLink:     fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s", attrs.Name),
	}

	// Map advanced options
	if attrs.Encryption != nil {
		response.KMSKeyName = attrs.Encryption.DefaultKMSKeyName
	}

	if attrs.RetentionPolicy != nil {
		response.RetentionPolicy = &models.RetentionPolicy{
			RetentionPeriodSeconds: int64(attrs.RetentionPolicy.RetentionPeriod.Seconds()),
			IsLocked:               attrs.RetentionPolicy.IsLocked,
		}
	}

	response.UniformBucketLevelAccess = attrs.UniformBucketLevelAccess.Enabled

	switch attrs.PublicAccessPrevention {
	case storage.PublicAccessPreventionEnforced:
		response.PublicAccessPrevention = "enforced"
	case storage.PublicAccessPreventionInherited:
		response.PublicAccessPrevention = "inherited"
	case storage.PublicAccessPreventionUnspecified:
		response.PublicAccessPrevention = "unspecified"
	}

	return response
}

func (c *GCPStorageClient) mapObjectAttrsToResponse(attrs *storage.ObjectAttrs) *models.ObjectResponse {
	return &models.ObjectResponse{
		Name:         attrs.Name,
		Bucket:       attrs.Bucket,
		Size:         attrs.Size,
		ContentType:  attrs.ContentType,
		MD5Hash:      fmt.Sprintf("%x", attrs.MD5),
		CRC32C:       fmt.Sprintf("%x", attrs.CRC32C),
		CreateTime:   attrs.Created,
		UpdateTime:   attrs.Updated,
		Generation:   attrs.Generation,
		StorageClass: attrs.StorageClass,
		Metadata:     attrs.Metadata,
		SelfLink:     fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s/o/%s", attrs.Bucket, attrs.Name),
	}
}

// Note: parseTimeString function removed as it was unused
// Can be re-added when lifecycle policy management is fully implemented
