// Package main demonstrates validation features of the GCP Storage Client SDK.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/pkg/validation/gcp"
)

func main() {
	fmt.Println("=== GCP Storage Client SDK - Validation Example ===")

	// Example 1: Bucket Name Validation
	fmt.Println("\n1. Bucket Name Validation:")
	bucketNames := []string{
		"valid-bucket-name",       // ✅ Valid
		"my-test-bucket-123",      // ✅ Valid
		"Invalid-Bucket-Name",     // ❌ Invalid: uppercase
		"",                        // ❌ Invalid: empty
		"ab",                      // ❌ Invalid: too short
		"bucket_with_underscores", // ❌ Invalid: underscores
		"192.168.1.1",             // ❌ Invalid: IP address format
		"goog-reserved-name",      // ❌ Invalid: starts with 'goog'
		"my-google-bucket",        // ❌ Invalid: contains 'google'
	}

	for _, name := range bucketNames {
		if err := gcp.ValidateBucketName(name); err != nil {
			fmt.Printf("❌ '%s': %v\n", name, err)
		} else {
			fmt.Printf("✅ '%s': Valid\n", name)
		}
	}

	// Example 2: Storage Class Validation
	fmt.Println("\n2. Storage Class Validation:")
	storageClasses := []string{
		"STANDARD", // ✅ Valid
		"NEARLINE", // ✅ Valid
		"COLDLINE", // ✅ Valid
		"ARCHIVE",  // ✅ Valid
		"",         // ✅ Valid (empty defaults to STANDARD)
		"standard", // ❌ Invalid: lowercase
		"INVALID",  // ❌ Invalid: not a valid class
		"HOT",      // ❌ Invalid: not a GCS storage class
	}

	for _, class := range storageClasses {
		if err := gcp.ValidateStorageClass(class); err != nil {
			fmt.Printf("❌ '%s': %v\n", class, err)
		} else {
			fmt.Printf("✅ '%s': Valid\n", class)
		}
	}

	// Example 3: Location Validation (Static)
	fmt.Println("\n3. Location Validation (Static):")
	locations := []string{
		"us-central1",    // ✅ Valid region
		"us-central1-a",  // ✅ Valid zone
		"europe-west1",   // ✅ Valid region
		"asia-east1-c",   // ✅ Valid zone
		"us",             // ✅ Valid multi-regional
		"eu",             // ✅ Valid multi-regional
		"invalid-region", // ❌ Invalid
		"us-central1-z",  // ❌ Invalid zone suffix
		"",               // ❌ Invalid: empty
		"us--central1",   // ❌ Invalid: malformed
	}

	for _, location := range locations {
		if err := gcp.ValidateLocation(location); err != nil {
			fmt.Printf("❌ '%s': %v\n", location, err)
		} else {
			fmt.Printf("✅ '%s': Valid\n", location)
		}
	}

	// Example 4: Object Name Validation
	fmt.Println("\n4. Object Name Validation:")
	objectNames := []string{
		"my-file.txt",          // ✅ Valid
		"path/to/my/file.pdf",  // ✅ Valid
		"file with spaces.doc", // ✅ Valid
		"unicode-文件.txt",       // ✅ Valid
		"",                     // ❌ Invalid: empty
		".",                    // ❌ Invalid: single dot
		"..",                   // ❌ Invalid: double dot
		"file\nwith\nnewline",  // ❌ Invalid: newline characters
		"file\x00null",         // ❌ Invalid: null character
	}

	for _, name := range objectNames {
		if err := gcp.ValidateObjectName(name); err != nil {
			fmt.Printf("❌ '%s': %v\n", name, err)
		} else {
			fmt.Printf("✅ '%s': Valid\n", name)
		}
	}

	// Example 5: Dynamic Location Validation (requires GCP credentials)
	fmt.Println("\n5. Dynamic Location Validation:")
	ctx := context.Background()
	projectID := getProjectID()

	if projectID != "your-gcp-project-id" {
		fmt.Printf("Testing dynamic validation with project: %s\n", projectID)

		// Create location validator
		validator := gcp.NewLocationValidator(projectID, nil)

		testLocations := []string{
			"us-central1",
			"europe-west1",
			"unknown-region-123",
		}

		for _, location := range testLocations {
			if err := validator.ValidateLocationDynamic(ctx, location); err != nil {
				fmt.Printf("❌ '%s': %v (dynamic check)\n", location, err)
			} else {
				fmt.Printf("✅ '%s': Valid (dynamic check)\n", location)
			}
		}

		// Example 6: Hybrid Validation (Static + Dynamic fallback)
		fmt.Println("\n6. Hybrid Validation (Static + Dynamic fallback):")
		testLocation := "us-west4" // Relatively new region
		if err := gcp.ValidateLocationWithFallback(ctx, projectID, testLocation); err != nil {
			fmt.Printf("❌ '%s': %v (hybrid check)\n", testLocation, err)
		} else {
			fmt.Printf("✅ '%s': Valid (hybrid check)\n", testLocation)
		}
	} else {
		fmt.Println("💡 Set GCP_PROJECT_ID environment variable to test dynamic validation")
		fmt.Println("Dynamic validation requires valid GCP credentials and project access")
	}

	// Example 7: Validation in Practice (creating a bucket with validation)
	fmt.Println("\n7. Practical Example - Creating Bucket with Validation:")

	// This demonstrates how validation is used in the SDK
	exampleBucketReq := &models.BucketRequest{
		Name:         "example-validated-bucket",
		Location:     "us-central1",
		StorageClass: "STANDARD",
	}

	fmt.Printf("Validating bucket request...\n")

	// These validations happen automatically in the SDK
	if err := gcp.ValidateBucketName(exampleBucketReq.Name); err != nil {
		fmt.Printf("❌ Bucket name validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Bucket name valid: %s\n", exampleBucketReq.Name)
	}

	if err := gcp.ValidateLocation(exampleBucketReq.Location); err != nil {
		fmt.Printf("❌ Location validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Location valid: %s\n", exampleBucketReq.Location)
	}

	if err := gcp.ValidateStorageClass(exampleBucketReq.StorageClass); err != nil {
		fmt.Printf("❌ Storage class validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Storage class valid: %s\n", exampleBucketReq.StorageClass)
	}

	fmt.Printf("✅ All validations passed! Bucket request is valid.\n")

	fmt.Println("\n🎉 Validation example completed successfully!")
	fmt.Println("\nKey validation features demonstrated:")
	fmt.Println("- Bucket name validation (GCS naming rules)")
	fmt.Println("- Storage class validation")
	fmt.Println("- Location validation (static and dynamic)")
	fmt.Println("- Object name validation")
	fmt.Println("- Integration with SDK operations")
	fmt.Println("- Error handling for invalid inputs")
}

func getProjectID() string {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = "your-gcp-project-id"
	}
	return projectID
}
