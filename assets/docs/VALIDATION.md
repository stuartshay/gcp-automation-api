# Request Validation System

## Overview

The GCP Automation API implements comprehensive request validation using the `go-playground/validator/v10` library. This ensures all API endpoints receive valid, properly formatted data before processing.

## Implementation

### Validation Components

1. **Custom Validator (`internal/validators/validators.go`)**
   - Implements Echo-compatible validator interface
   - Registers custom validation rules specific to GCP resources
   - Provides user-friendly error message formatting

2. **Request Models (`internal/models/models.go`)**
   - All request structs include validation tags
   - Uses both built-in and custom validation rules

3. **Handler Integration (`internal/handlers/handlers.go`)**
   - All endpoints use the custom validator
   - Returns structured error responses for validation failures

### Custom Validation Rules

#### `project_id`
- **Purpose**: Validates GCP project ID format
- **Rules**:
  - 6-30 characters
  - Lowercase letters, digits, hyphens only
  - Must start with a letter
  - Cannot end with a hyphen

#### `bucket_name`
- **Purpose**: Validates Google Cloud Storage bucket names
- **Rules**:
  - 3-63 characters
  - Lowercase letters, digits, dashes, underscores, dots
  - No consecutive dots
  - Cannot be formatted as an IP address

#### `label_key`
- **Purpose**: Validates GCP resource label keys
- **Rules**:
  - 1-63 characters
  - Start with lowercase letter
  - Lowercase letters, digits, underscores, dashes only

#### `label_value`
- **Purpose**: Validates GCP resource label values
- **Rules**:
  - 0-63 characters
  - Lowercase letters, digits, underscores, dashes only

#### `gcp_location`
- **Purpose**: Validates GCP regions and zones
- **Rules**: Must match known GCP location patterns

### Usage Examples

#### Project Request
```json
{
  "project_id": "my-gcp-project",
  "display_name": "My GCP Project",
  "parent_id": "123456789",
  "parent_type": "organization",
  "labels": {
    "environment": "production",
    "team": "backend"
  }
}
```

#### Bucket Request
```json
{
  "name": "my-storage-bucket",
  "location": "us-central1",
  "storage_class": "STANDARD",
  "labels": {
    "purpose": "data-lake",
    "retention": "long-term"
  },
  "versioning": true
}
```

### Error Response Format

When validation fails, the API returns a structured error response:

```json
{
  "error": "validation_failed",
  "message": "project_id must be a valid GCP project ID (6-30 chars, lowercase letters/digits/hyphens, start with letter, not end with hyphen)",
  "code": 400
}
```

### Testing

Comprehensive test suite in `internal/validators/validators_test.go` covers:
- Valid request scenarios
- Invalid field formats
- Missing required fields
- Edge cases for each validation rule

### Benefits

1. **Early Validation**: Catches invalid data before GCP API calls
2. **User-Friendly Errors**: Clear, actionable error messages
3. **Type Safety**: Compile-time validation of request structures
4. **Consistency**: Uniform validation across all endpoints
5. **Security**: Prevents injection and malformed data attacks
6. **Performance**: Reduces unnecessary GCP API calls

### Configuration

The validator is initialized once at startup and shared across all handlers for optimal performance. Custom validation functions are registered during initialization and cached for reuse.

## Best Practices

1. Always validate input at the API boundary
2. Provide clear, actionable error messages
3. Use GCP-specific validation rules for resource names
4. Test both positive and negative validation scenarios
5. Keep validation rules in sync with GCP requirements
