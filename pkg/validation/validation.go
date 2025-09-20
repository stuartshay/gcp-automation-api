// Package validation provides comprehensive validation utilities for GCP resources.
//
// This package is organized into sub-packages for different types of validation:
//   - gcp: GCP-specific validation (locations, buckets, objects, etc.)
//
// Example usage:
//
//	import "github.com/stuartshay/gcp-automation-api/pkg/validation/gcp"
//
//	// Validate a GCS bucket name
//	if err := gcp.ValidateBucketName("my-bucket"); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Validate a GCP location
//	if err := gcp.ValidateLocation("us-central1"); err != nil {
//	    log.Fatal(err)
//	}
package validation

// This file exists to provide package documentation.
// The actual validation implementations are in sub-packages.
