package sdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/stuartshay/gcp-automation-api/internal/models"
	"google.golang.org/api/option"
)

// Example demonstrates how to use the GCP Storage Client
func Example() {
	// Initialize the client
	ctx := context.Background()
	projectID := "your-gcp-project-id"

	// Option 1: Use Application Default Credentials
	client, err := NewGCPStorageClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to create storage client: %v\n", err)
		return
	}
	defer client.Close()

	// Option 2: Use service account key file
	// client, err := NewGCPStorageClient(ctx, projectID, option.WithCredentialsFile("path/to/service-account.json"))

	// Create a bucket
	bucketReq := &models.BucketRequest{
		Name:         "my-example-bucket-12345",
		Location:     "us-central1",
		StorageClass: "STANDARD",
		Labels: map[string]string{
			"environment": "development",
			"team":        "backend",
		},
		Versioning:               true,
		UniformBucketLevelAccess: true,
		PublicAccessPrevention:   "enforced",
	}

	bucket, err := client.CreateBucket(ctx, bucketReq)
	if err != nil {
		fmt.Printf("Failed to create bucket: %v\n", err)
		return
	}
	fmt.Printf("Created bucket: %s\n", bucket.Name)

	// Upload an object
	data := strings.NewReader("Hello, World! This is test data.")
	object, err := client.UploadObject(ctx, bucket.Name, "test-file.txt", data)
	if err != nil {
		fmt.Printf("Failed to upload object: %v\n", err)
		return
	}
	fmt.Printf("Uploaded object: %s\n", object.Name)

	// List objects in the bucket
	objects, err := client.ListObjects(ctx, bucket.Name, "")
	if err != nil {
		fmt.Printf("Failed to list objects: %v\n", err)
		return
	}
	fmt.Printf("Found %d objects in bucket\n", len(objects))

	// Check if object exists
	exists, err := client.ObjectExists(ctx, bucket.Name, "test-file.txt")
	if err != nil {
		fmt.Printf("Failed to check object existence: %v\n", err)
		return
	}
	fmt.Printf("Object exists: %t\n", exists)

	// Get object metadata
	metadata, err := client.GetObjectMetadata(ctx, bucket.Name, "test-file.txt")
	if err != nil {
		fmt.Printf("Failed to get object metadata: %v\n", err)
		return
	}
	fmt.Printf("Object size: %d bytes\n", metadata.Size)

	// Download the object
	reader, err := client.DownloadObject(ctx, bucket.Name, "test-file.txt")
	if err != nil {
		fmt.Printf("Failed to download object: %v\n", err)
		return
	}
	defer reader.Close()
	fmt.Println("Successfully downloaded object")

	// Update bucket settings
	updateReq := &models.BucketUpdateRequest{
		Labels: map[string]string{
			"environment": "production",
			"updated":     "true",
		},
	}

	updatedBucket, err := client.UpdateBucket(ctx, bucket.Name, updateReq)
	if err != nil {
		fmt.Printf("Failed to update bucket: %v\n", err)
		return
	}
	fmt.Printf("Updated bucket labels: %v\n", updatedBucket.Labels)

	// List all buckets
	buckets, err := client.ListBuckets(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to list buckets: %v\n", err)
		return
	}
	fmt.Printf("Found %d buckets in project\n", len(buckets))

	// Check if bucket exists
	bucketExists, err := client.BucketExists(ctx, bucket.Name)
	if err != nil {
		fmt.Printf("Failed to check bucket existence: %v\n", err)
		return
	}
	fmt.Printf("Bucket exists: %t\n", bucketExists)

	// Get bucket details
	bucketDetails, err := client.GetBucket(ctx, bucket.Name)
	if err != nil {
		fmt.Printf("Failed to get bucket details: %v\n", err)
		return
	}
	fmt.Printf("Bucket location: %s\n", bucketDetails.Location)

	// Clean up - delete object and bucket
	if err := client.DeleteObject(ctx, bucket.Name, "test-file.txt"); err != nil {
		fmt.Printf("Failed to delete object: %v\n", err)
		return
	}
	fmt.Println("Deleted object")

	if err := client.DeleteBucket(ctx, bucket.Name); err != nil {
		fmt.Printf("Failed to delete bucket: %v\n", err)
		return
	}
	fmt.Println("Deleted bucket")

	fmt.Println("Example completed successfully!")
}

// ExampleWithCredentials shows how to create a client with custom credentials
func ExampleWithCredentials() error {
	ctx := context.Background()
	projectID := "your-gcp-project-id"
	credentialsFile := "path/to/service-account.json"

	// Create client with credentials file
	client, err := NewGCPStorageClient(
		ctx,
		projectID,
		option.WithCredentialsFile(credentialsFile),
	)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	// Use the client...
	buckets, err := client.ListBuckets(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to list buckets: %w", err)
	}

	fmt.Printf("Found %d buckets\n", len(buckets))
	return nil
}

// ExampleErrorHandling demonstrates proper error handling patterns
func ExampleErrorHandling() {
	ctx := context.Background()
	projectID := "your-gcp-project-id"

	client, err := NewGCPStorageClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}
	defer client.Close()

	// Try to get a non-existent bucket
	_, err = client.GetBucket(ctx, "non-existent-bucket")
	if err != nil {
		fmt.Printf("Expected error getting non-existent bucket: %v\n", err)
	}

	// Try to create a bucket with invalid name
	invalidReq := &models.BucketRequest{
		Name:     "Invalid-Bucket-Name", // Invalid: contains uppercase
		Location: "us-central1",
	}

	_, err = client.CreateBucket(ctx, invalidReq)
	if err != nil {
		fmt.Printf("Expected validation error: %v\n", err)
	}

	// Try to upload to non-existent bucket
	data := strings.NewReader("test data")
	_, err = client.UploadObject(ctx, "non-existent-bucket", "test.txt", data)
	if err != nil {
		fmt.Printf("Expected error uploading to non-existent bucket: %v\n", err)
	}
}

// ExampleBulkOperations demonstrates how to perform bulk operations efficiently
func ExampleBulkOperations() error {
	ctx := context.Background()
	projectID := "your-gcp-project-id"

	client, err := NewGCPStorageClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	bucketName := "bulk-operations-example-bucket"

	// Create bucket
	bucketReq := &models.BucketRequest{
		Name:         bucketName,
		Location:     "us-central1",
		StorageClass: "STANDARD",
	}

	_, err = client.CreateBucket(ctx, bucketReq)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	// Upload multiple objects
	for i := 0; i < 10; i++ {
		objectName := fmt.Sprintf("file-%d.txt", i)
		data := strings.NewReader(fmt.Sprintf("Content of file %d", i))

		_, err := client.UploadObject(ctx, bucketName, objectName, data)
		if err != nil {
			fmt.Printf("Failed to upload %s: %v\n", objectName, err)
			continue
		}
		fmt.Printf("Uploaded %s\n", objectName)
	}

	// List all objects
	objects, err := client.ListObjects(ctx, bucketName, "file-")
	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}

	fmt.Printf("Found %d objects with 'file-' prefix\n", len(objects))

	// Delete all objects
	for _, obj := range objects {
		if err := client.DeleteObject(ctx, bucketName, obj.Name); err != nil {
			fmt.Printf("Failed to delete %s: %v\n", obj.Name, err)
			continue
		}
		fmt.Printf("Deleted %s\n", obj.Name)
	}

	// Delete bucket
	if err := client.DeleteBucket(ctx, bucketName); err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}

	fmt.Println("Bulk operations completed successfully!")
	return nil
}
