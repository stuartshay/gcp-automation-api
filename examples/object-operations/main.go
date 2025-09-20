// Package main demonstrates object operations using the GCP Storage Client SDK.
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

	fmt.Println("=== Object Operations Example ===")

	client, err := sdk.NewGCPStorageClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer client.Close()

	bucketName := fmt.Sprintf("sdk-obj-ops-%d", generateUniqueSuffix())
	bucketReq := &models.BucketRequest{
		Name:         bucketName,
		Location:     "us-central1",
		StorageClass: "STANDARD",
	}

	_, err = client.CreateBucket(ctx, bucketReq)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}
	fmt.Printf("✅ Created bucket: %s\n", bucketName)

	content := "Hello, World! This is a test file."
	reader := strings.NewReader(content)
	object, err := client.UploadObject(ctx, bucketName, "test.txt", reader)
	if err != nil {
		log.Fatalf("Failed to upload object: %v", err)
	}
	fmt.Printf("✅ Uploaded: %s (%d bytes)\n", object.Name, object.Size)

	if err := client.DeleteObject(ctx, bucketName, "test.txt"); err != nil {
		log.Printf("Failed to delete object: %v", err)
	}

	if err := client.DeleteBucket(ctx, bucketName); err != nil {
		log.Fatalf("Failed to delete bucket: %v", err)
	}
	fmt.Printf("✅ Deleted bucket: %s\n", bucketName)

	fmt.Println("🎉 Object operations example completed!")
}

func getProjectID() string {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = "your-gcp-project-id"
	}
	return projectID
}

func generateUniqueSuffix() int64 {
	return int64(os.Getpid()) % 100000
}
