# GCP Storage Client SDK

A comprehensive Go SDK for managing Google Cloud Storage resources including buckets and objects.

## Overview

This SDK provides a clean, well-validated interface for Google Cloud Storage operations. It supports:

- **Bucket Operations**: Create, Read, Update, Delete, List buckets
- **Object Operations**: Upload, Download, Delete, List objects
- **Advanced Features**: KMS encryption, retention policies, IAM controls
- **Validation**: Comprehensive input validation following GCS naming rules
- **Error Handling**: Structured error handling with contextual information

## Installation

```bash
go get github.com/stuartshay/gcp-automation-api/pkg/sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "strings"

    "github.com/stuartshay/gcp-automation-api/pkg/sdk"
    "github.com/stuartshay/gcp-automation-api/internal/models"
)

func main() {
    ctx := context.Background()

    // Create client
    client, err := sdk.NewGCPStorageClient(ctx, "your-project-id")
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Create a bucket
    bucket, err := client.CreateBucket(ctx, &models.BucketRequest{
        Name:         "my-test-bucket",
        Location:     "us-central1",
        StorageClass: "STANDARD",
    })
    if err != nil {
        panic(err)
    }

    // Upload an object
    data := strings.NewReader("Hello, World!")
    obj, err := client.UploadObject(ctx, bucket.Name, "hello.txt", data)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Uploaded %s to %s\n", obj.Name, bucket.Name)
}
```

## Authentication

### Application Default Credentials (Recommended)

```go
client, err := sdk.NewGCPStorageClient(ctx, projectID)
```

### Service Account Key File

```go
client, err := sdk.NewGCPStorageClient(
    ctx,
    projectID,
    option.WithCredentialsFile("path/to/service-account.json"),
)
```

### Environment Variables

Set `GOOGLE_APPLICATION_CREDENTIALS` to point to your service account key file:

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account.json"
```

## API Reference

### StorageClient Interface

```go
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

    // Advanced features (simplified implementations)
    SetBucketLifecycle(ctx context.Context, bucketName string, lifecycle *models.LifecyclePolicy) error
    GetBucketLifecycle(ctx context.Context, bucketName string) (*models.LifecyclePolicy, error)
    DeleteBucketLifecycle(ctx context.Context, bucketName string) error
    SetBucketIAM(ctx context.Context, bucketName string, policy *models.IAMPolicy) error
    GetBucketIAM(ctx context.Context, bucketName string) (*models.IAMPolicy, error)
    TestBucketIAM(ctx context.Context, bucketName string, permissions []string) ([]string, error)

    // Cleanup
    Close() error
}
```

## Bucket Operations

### Create Bucket

```go
bucket, err := client.CreateBucket(ctx, &models.BucketRequest{
    Name:         "my-bucket-name",
    Location:     "us-central1",
    StorageClass: "STANDARD",
    Labels: map[string]string{
        "environment": "production",
        "team":        "backend",
    },
    Versioning:               true,
    UniformBucketLevelAccess: true,
    PublicAccessPrevention:   "enforced",
    KMSKeyName:               "projects/PROJECT_ID/locations/LOCATION/keyRings/RING_ID/cryptoKeys/KEY_ID",
    RetentionPolicy: &models.RetentionPolicy{
        RetentionPeriodSeconds: 3600, // 1 hour
        IsLocked:               false,
    },
})
```

### Get Bucket

```go
bucket, err := client.GetBucket(ctx, "bucket-name")
if err != nil {
    // Handle error
}
fmt.Printf("Bucket: %s, Location: %s\n", bucket.Name, bucket.Location)
```

### List Buckets

```go
buckets, err := client.ListBuckets(ctx, "project-id")
if err != nil {
    // Handle error
}
for _, bucket := range buckets {
    fmt.Printf("Bucket: %s\n", bucket.Name)
}
```

### Update Bucket

```go
versioning := false
updated, err := client.UpdateBucket(ctx, "bucket-name", &models.BucketUpdateRequest{
    Versioning: &versioning,
    Labels: map[string]string{
        "updated": "true",
    },
})
```

### Delete Bucket

```go
err := client.DeleteBucket(ctx, "bucket-name")
```

## Object Operations

### Upload Object

```go
data := strings.NewReader("file content")
obj, err := client.UploadObject(ctx, "bucket-name", "path/to/file.txt", data)
if err != nil {
    // Handle error
}
fmt.Printf("Uploaded: %s (%d bytes)\n", obj.Name, obj.Size)
```

### Download Object

```go
reader, err := client.DownloadObject(ctx, "bucket-name", "path/to/file.txt")
if err != nil {
    // Handle error
}
defer reader.Close()

// Read the content
content, err := io.ReadAll(reader)
```

### List Objects

```go
// List all objects
objects, err := client.ListObjects(ctx, "bucket-name", "")

// List objects with prefix
objects, err := client.ListObjects(ctx, "bucket-name", "path/to/")
```

### Delete Object

```go
err := client.DeleteObject(ctx, "bucket-name", "path/to/file.txt")
```

## Validation

The SDK includes comprehensive validation for:

### Bucket Names
- Must be 3-63 characters long
- Can only contain lowercase letters, numbers, hyphens, and periods
- Cannot start/end with hyphens or periods
- Cannot contain consecutive periods
- Cannot be formatted as IP addresses
- Cannot start with "goog" or contain "google"

### Object Names
- Must be 1-1024 characters long
- Cannot contain newline, carriage return, or null characters
- Cannot be "." or ".."

### Storage Classes
- STANDARD, NEARLINE, COLDLINE, ARCHIVE

## Error Handling

All errors are wrapped with contextual information:

```go
bucket, err := client.CreateBucket(ctx, req)
if err != nil {
    // Error format: "creating bucket bucket-name: detailed error message"
    fmt.Printf("Error: %v\n", err)

    // Check for specific error types
    if strings.Contains(err.Error(), "invalid") {
        // Handle validation error
    }
}
```

## Testing

Run the test suite:

```bash
go test ./pkg/sdk/...
```

Run with coverage:

```bash
go test -cover ./pkg/sdk/...
```

Run benchmarks:

```bash
go test -bench=. ./pkg/sdk/...
```

## Examples

See the `examples.go` file for comprehensive usage examples including:

- Basic operations
- Error handling patterns
- Bulk operations
- Custom credentials

## Best Practices

1. **Always use context**: Pass context for timeouts and cancellation
2. **Close readers**: Always close io.ReadCloser from download operations
3. **Validate inputs**: The SDK validates inputs, but check your data first
4. **Handle errors**: Check all errors and handle them appropriately
5. **Use defer**: Always defer client.Close() after creation
6. **Batch operations**: For multiple operations, reuse the same client

## Limitations

- Lifecycle policy management is simplified (returns not implemented errors)
- IAM policy management is simplified (basic implementation only)
- Some advanced GCS features are not yet implemented

## Contributing

1. Follow Go conventions and best practices
2. Add tests for new functionality
3. Update documentation
4. Ensure validation is comprehensive
5. Handle errors appropriately

## License

This project is licensed under the MIT License.
