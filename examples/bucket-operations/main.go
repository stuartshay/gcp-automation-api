// Package main demonstrates comprehensive bucket operations using the GCP Storage Client SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/pkg/sdk"
)

func main() {
	ctx := context.Background()
	projectID := getProjectID()

	fmt.Println("=== GCP Storage Client SDK - Bucket Operations Example ===")
	fmt.Printf("Project ID: %s\n", projectID)

	client, err := sdk.NewGCPStorageClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer client.Close()

	bucketName := fmt.Sprintf("sdk-bucket-ops-%d", generateUniqueSuffix())

	fmt.Printf("\n🔄 Creating bucket with comprehensive configuration: %s\n", bucketName)
	bucketReq := &models.BucketRequest{
		Name:         bucketName,
		Location:     "us-central1",
		StorageClass: "STANDARD",
		Labels: map[string]string{
			"environment": "development",
			"team":        "backend",
			"purpose":     "example",
		},
		Versioning:               true,
		UniformBucketLevelAccess: true,
		PublicAccessPrevention:   "enforced",
	}

	bucket, err := client.CreateBucket(ctx, bucketReq)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}
	fmt.Printf("✅ Created bucket: %s\n", bucket.Name)

	fmt.Printf("\n🔄 Checking if bucket exists: %s\n", bucketName)
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("Failed to check bucket existence: %v", err)
	}
	fmt.Printf("✅ Bucket exists: %t\n", exists)

	fmt.Printf("\n🔄 Getting bucket details: %s\n", bucketName)
	bucketDetails, err := client.GetBucket(ctx, bucketName)
	if err != nil {
		log.Fatalf("Failed to get bucket details: %v", err)
	}
	fmt.Printf("✅ Bucket details:\n")
	fmt.Printf("   Name: %s\n", bucketDetails.Name)
	fmt.Printf("   Location: %s\n", bucketDetails.Location)
	fmt.Printf("   Storage Class: %s\n", bucketDetails.StorageClass)

	fmt.Printf("\n🔄 Deleting bucket: %s\n", bucketName)
	if err := client.DeleteBucket(ctx, bucketName); err != nil {
		log.Fatalf("Failed to delete bucket: %v", err)
	}
	fmt.Printf("✅ Deleted bucket: %s\n", bucketName)

	fmt.Println("\n🎉 Bucket operations example completed successfully!")
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
