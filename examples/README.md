# GCP Storage Client SDK Examples

This directory contains complete, runnable examples demonstrating how to use the GCP Storage Client
SDK.

## Prerequisites

Before running any examples, ensure you have:

1. **Go 1.21+** installed
2. **GCP Project** with Cloud Storage API enabled
3. **Authentication** set up (one of the following):
   - Application Default Credentials: `gcloud auth application-default login`
   - Service Account Key: Set `GOOGLE_APPLICATION_CREDENTIALS` environment variable
   - Workload Identity (for GKE/Cloud Run environments)

## Examples

### 1. Basic Usage (`basic/`)

Demonstrates fundamental operations:

- Client initialization
- Creating and deleting buckets
- Basic object operations
- Error handling

```bash
cd examples/basic
go run main.go
```

### 2. Bucket Operations (`bucket-operations/`)

Comprehensive bucket management:

- Creating buckets with various configurations
- Updating bucket settings
- Listing and managing buckets
- Bucket existence checks

```bash
cd examples/bucket-operations
go run main.go
```

### 3. Object Operations (`object-operations/`)

Complete object lifecycle management:

- Uploading objects from various sources
- Downloading objects
- Object metadata operations
- Bulk operations
- Object listing with prefixes

```bash
cd examples/object-operations
go run main.go
```

### 4. Validation (`validation/`)

Input validation examples:

- Bucket name validation
- Location validation (static and dynamic)
- Storage class validation
- Error handling for invalid inputs

```bash
cd examples/validation
go run main.go
```

## Configuration

Each example uses placeholder values that you should replace:

- `your-gcp-project-id` → Your actual GCP project ID
- `example-bucket-name` → A unique bucket name
- `path/to/service-account.json` → Path to your service account key (if using)

## Environment Setup

For quick setup, you can export these environment variables:

```bash
export GCP_PROJECT_ID="your-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account.json"
```

Then modify the examples to use `os.Getenv("GCP_PROJECT_ID")` instead of hardcoded values.

## Running All Examples

To run all examples in sequence:

```bash
# From the repository root
make run-examples

# Or manually:
for dir in examples/*/; do
    echo "Running example in $dir"
    (cd "$dir" && go run main.go)
    echo "---"
done
```

## Error Handling

All examples include comprehensive error handling patterns that you should follow in your own code:

- Always check for errors from SDK operations
- Use proper context for timeouts and cancellation
- Handle validation errors gracefully
- Clean up resources in defer statements

## Best Practices Demonstrated

- **Resource cleanup**: All examples properly clean up created resources
- **Context usage**: Proper context handling for operations
- **Error handling**: Comprehensive error checking and reporting
- **Validation**: Input validation before API calls
- **Configuration**: Flexible client configuration options

## Troubleshooting

If you encounter issues:

1. **Authentication errors**: Ensure `gcloud auth application-default login` is run
2. **Project not found**: Verify your project ID is correct
3. **API not enabled**: Enable the Cloud Storage API in your GCP project
4. **Permissions**: Ensure your account has Storage Admin or appropriate IAM roles

For more details, see the main [SDK documentation](../pkg/sdk/README.md).
