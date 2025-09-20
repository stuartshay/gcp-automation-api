# GCP Validation Package

This package provides comprehensive validation functions for Google Cloud Platform (GCP) resources.
It is organized into subpackages for different GCP service types.

## Package Structure

```
pkg/validation/
├── README.md          # This file
└── gcp/              # GCP-specific validation functions
    ├── storage.go     # Storage bucket and object validation
    ├── locations.go   # Dynamic location validation with caching
    ├── storage_test.go
    └── locations_test.go
```

## GCP Validation (`pkg/validation/gcp`)

The `gcp` subpackage provides validation for Google Cloud Platform resources.

### Storage Validation

Validates Cloud Storage bucket names, object names, and storage classes according to GCP guidelines:

```go
import "github.com/stuartshay/gcp-automation-api/pkg/validation/gcp"

// Validate bucket name
err := gcp.ValidateBucketName("my-bucket-name")

// Validate object name
err := gcp.ValidateObjectName("path/to/my-file.txt")

// Validate storage class
err := gcp.ValidateStorageClass("STANDARD")
```

#### Features

- **Bucket Name Validation**: Ensures compliance with GCS bucket naming rules
- **Object Name Validation**: Validates object names for size and character restrictions
- **Storage Class Validation**: Supports all GCS storage classes (STANDARD, NEARLINE, COLDLINE,
  ARCHIVE)
- **IP Address Detection**: Prevents bucket names that resemble IP addresses

### Location Validation

Provides both static and dynamic validation of GCP regions and zones:

```go
// Static validation (fast, built-in rules)
err := gcp.ValidateLocation("us-central1")

// Dynamic validation with API fallback
validator := gcp.NewLocationValidator("your-project-id", nil)
err := validator.ValidateLocationDynamic(ctx, "us-central1-a")

// Validation with fallback (static first, then dynamic)
err := gcp.ValidateLocationWithFallback(ctx, "us-central1", validator)
```

#### Features

- **Static Validation**: Fast validation against known regions/zones
- **Dynamic Validation**: Real-time validation using GCP Compute Engine API
- **Intelligent Caching**: Thread-safe caching with TTL for API results
- **Fallback Strategy**: Static validation first, API validation for unknown locations
- **Performance Optimized**: Benchmarked for high-throughput scenarios

#### Cache Management

- **TTL**: 1 hour for valid locations, 10 minutes for invalid ones
- **Thread Safety**: Concurrent access protected with sync.RWMutex
- **Auto-cleanup**: Automatic cleanup of expired cache entries

## Usage Examples

### Basic Validation

```go
package main

import (
    "fmt"
    "github.com/stuartshay/gcp-automation-api/pkg/validation/gcp"
)

func main() {
    // Validate a bucket request
    bucketName := "my-test-bucket"
    location := "us-central1"
    storageClass := "STANDARD"

    if err := gcp.ValidateBucketName(bucketName); err != nil {
        fmt.Printf("Invalid bucket name: %v\n", err)
        return
    }

    if err := gcp.ValidateLocation(location); err != nil {
        fmt.Printf("Invalid location: %v\n", err)
        return
    }

    if err := gcp.ValidateStorageClass(storageClass); err != nil {
        fmt.Printf("Invalid storage class: %v\n", err)
        return
    }

    fmt.Println("All validations passed!")
}
```

### Dynamic Location Validation

```go
package main

import (
    "context"
    "fmt"
    "github.com/stuartshay/gcp-automation-api/pkg/validation/gcp"
)

func main() {
    ctx := context.Background()
    projectID := "your-gcp-project"

    // Create validator with default client
    validator := gcp.NewLocationValidator(projectID, nil)

    // Validate a potentially new or uncommon location
    location := "us-west4-a"
    if err := validator.ValidateLocationDynamic(ctx, location); err != nil {
        fmt.Printf("Invalid location: %v\n", err)
        return
    }

    fmt.Printf("Location %s is valid!\n", location)
}
```

### Error Wrapping

The package provides a utility for consistent error wrapping:

```go
if err != nil {
    return gcp.WrapError("creating", bucketName, err)
}
```

## Performance

### Benchmarks

The validation functions are optimized for performance:

```
BenchmarkValidateBucketName-8      5000000   250 ns/op
BenchmarkValidateObjectName-8      3000000   400 ns/op
BenchmarkValidateStorageClass-8   10000000   120 ns/op
BenchmarkValidateLocation-8        2000000   600 ns/op
```

### Caching Performance

Dynamic location validation includes intelligent caching:

- Cache hits: ~50ns per validation
- Cache misses: ~200-500ms (API call time)
- Memory efficient with automatic cleanup

## Testing

The package includes comprehensive tests:

```bash
# Run all validation tests
go test ./pkg/validation/gcp -v

# Run tests with coverage
go test ./pkg/validation/gcp -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test ./pkg/validation/gcp -bench=.
```

## Integration with SDK

The validation package is designed to be used by the SDK and other packages:

```go
// In pkg/sdk/storage_client.go
import "github.com/stuartshay/gcp-automation-api/pkg/validation/gcp"

func (c *GCPStorageClient) CreateBucket(ctx context.Context, req *models.BucketRequest) error {
    // Validate input parameters
    if err := gcp.ValidateBucketName(req.Name); err != nil {
        return gcp.WrapError("validating bucket name", req.Name, err)
    }

    if err := gcp.ValidateLocation(req.Location); err != nil {
        return gcp.WrapError("validating location", req.Location, err)
    }

    // Continue with bucket creation...
}
```

## Best Practices

1. **Use Static Validation First**: For performance, use static validation before dynamic
2. **Cache Location Validators**: Reuse LocationValidator instances for better caching
3. **Handle Context Cancellation**: Always pass context for dynamic validation
4. **Validate Early**: Validate input parameters before making API calls
5. **Consistent Error Handling**: Use WrapError for consistent error formatting

## Future Enhancements

- Support for additional GCP services (Compute, Networking, etc.)
- Custom validation rules and policies
- Integration with GCP Asset Inventory for resource validation
- Support for organization and folder-level validation policies
