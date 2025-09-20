// Package main demonstrates basic usage of the GCP Storage Client SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/pkg/sdk"
)

func main() {
	ctx := context.Background()
	projectID := getProjectID()

	fmt.Println("=== GCP Storage Client SDK - Basic Example ===")
	fmt.Printf("Project ID: %s\n", projectID)

	client, err := sdk.NewGCPStorageClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ Successfully created storage client")

	bucketName := fmt.Sprintf("sdk-basic-example-%d", generateUniqueSuffix())

	bucketReq := &models.BucketRequest{
		Name:         bucketName,
		Location:     "us-central1",
		StorageClass: "STANDARD",
		Labels: map[string]string{
			"example": "basic",
			"sdk":     "gcp-automation-api",
		},
	}

	fmt.Printf("\n🔄 Creating bucket: %s\n", bucketName)
	bucket, err := client.CreateBucket(ctx, bucketReq)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}
	fmt.Printf("✅ Created bucket: %s in %s\n", bucket.Name, bucket.Location)

	objectName := "hello-world.txt"
	content := "Hello, World! This is a test file from the GCP Storage SDK."
	data := strings.NewReader(content)

	fmt.Printf("\n🔄 Uploading object: %s\n", objectName)
	object, err := client.UploadObject(ctx, bucketName, objectName, data)
	if err != nil {
		log.Fatalf("Failed to upload object: %v", err)
	}
	fmt.Printf("✅ Uploaded object: %s (%d bytes)\n", object.Name, object.Size)

	fmt.Printf("\n🔄 Deleting object: %s\n", objectName)
	if err := client.DeleteObject(ctx, bucketName, objectName); err != nil {
		log.Fatalf("Failed to delete object: %v", err)
	}
	fmt.Printf("✅ Deleted object: %s\n", objectName)

	fmt.Printf("\n🔄 Deleting bucket: %s\n", bucketName)
	if err := client.DeleteBucket(ctx, bucketName); err != nil {
		log.Fatalf("Failed to delete bucket: %v", err)
	}
	fmt.Printf("✅ Deleted bucket: %s\n", bucketName)

	fmt.Println("\n🎉 Basic example completed successfully!")
}

func getProjectID() string {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = "your-gcp-project-id"
		fmt.Printf("💡 Tip: Set GCP_PROJECT_ID environment variable. Using default: %s\n", projectID)
	}
	return projectID
}

func generateUniqueSuffix() int64 {
	return int64(os.Getpid()) % 100000
}
