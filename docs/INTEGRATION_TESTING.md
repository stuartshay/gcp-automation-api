# Integration Testing Guide

This document provides comprehensive information about the integration testing framework for the GCP Automation API.

## Overview

The integration testing framework provides comprehensive testing for:
- **Projects** - GCP project creation, retrieval, and deletion
- **Folders** - GCP folder management operations
- **Buckets** - Cloud Storage bucket operations

## Test Architecture

### Framework Components

1. **Mock Services** - Testify-based mocks for GCP operations
2. **Test Utilities** - Helper functions for setup and assertions
3. **Test Fixtures** - JSON-based test data
4. **Environment Support** - Mock and real GCP testing modes

### Directory Structure

```
tests/
├── integration/
│   ├── setup_test.go           # Test setup and utilities
│   ├── projects_test.go        # Project integration tests
│   ├── folders_test.go         # Folder integration tests (planned)
│   ├── buckets_test.go         # Bucket integration tests (planned)
│   └── mocks/
│       └── gcp_service_mock.go # Mock GCP service
├── fixtures/
│   ├── project_requests.json   # Project test data
│   ├── folder_requests.json    # Folder test data
│   └── bucket_requests.json    # Bucket test data
└── e2e/                        # End-to-end tests (planned)
```

## Running Tests

### Mock Mode (Default)

Run tests with mocked GCP services:

```bash
# Run all integration tests
make test-integration

# Run with coverage
make test-integration-coverage

# Run specific test
go test -v ./tests/integration/ -run TestProjectOperations
```

### Real GCP Integration Mode

Run tests against real GCP services (requires credentials):

```bash
# Set environment variables
export TEST_MODE=integration
export TEST_PROJECT_ID=your-test-project-id
export TEST_BUCKET_PREFIX=your-test-prefix

# Run integration tests
make test-integration-real
```

### All Tests

```bash
# Run both unit and integration tests
make test-all
```

## Environment Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TEST_MODE` | Test mode: `mock` or `integration` | `mock` | No |
| `TEST_PROJECT_ID` | GCP project ID for real tests | - | Yes (for real GCP) |
| `TEST_BUCKET_PREFIX` | Prefix for test bucket names | `test-gcp-automation` | No |
| `GOOGLE_APPLICATION_CREDENTIALS` | Path to GCP service account key | - | Yes (for real GCP) |

### Mock Mode Configuration

```bash
# Default - no configuration needed
make test-integration
```

### Real GCP Mode Configuration

```bash
# Set up GCP credentials
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
export TEST_PROJECT_ID=my-test-project-123
export TEST_BUCKET_PREFIX=mycompany-test

# Run real integration tests
make test-integration-real
```

## Test Categories

### Project Tests

**File**: `tests/integration/projects_test.go`

**Test Cases**:
- ✅ Create project with organization parent
- ✅ Create project with folder parent
- ✅ Create minimal project
- ✅ Get existing project
- ✅ Get non-existent project
- ✅ Delete existing project
- ✅ Delete non-existent project
- ✅ Validation errors (missing fields, invalid types)
- ✅ Authentication required

**Example**:
```bash
go test -v ./tests/integration/ -run TestProjectOperations
```

### Folder Tests (Planned)

**File**: `tests/integration/folders_test.go`

**Test Cases**:
- Create folder with organization parent
- Create folder with folder parent
- Get existing folder
- Delete existing folder
- Validation errors
- Authentication required

### Bucket Tests (Planned)

**File**: `tests/integration/buckets_test.go`

**Test Cases**:
- Create bucket with different storage classes
- Create bucket with versioning
- Create bucket with labels
- Get existing bucket
- Delete existing bucket
- Validation errors
- Authentication required

## Test Utilities

### Setup Functions

```go
// SetupTestServer creates a test server with mock or real GCP service
setup := SetupTestServer(t)
defer CleanupTestResources(t, setup)

// Generate test JWT token
token := GenerateTestJWT(t, setup.AuthService)
```

### Assertion Functions

```go
// Assert successful response with data
response := AssertSuccessResponseWithData(t, body, "Project created successfully")

// Assert successful response (may not have data)
response := AssertSuccessResponse(t, body, "Project deleted successfully")

// Assert error response
AssertErrorResponse(t, body, http.StatusBadRequest, "Validation failed")
```

### Test Data Helpers

```go
// Generate unique test names
bucketName := GetTestBucketName("test-prefix")
projectID := GetTestProjectID("test-prefix")
folderName := GetTestFolderName("test-prefix")
```

## Test Fixtures

### Project Fixtures

**File**: `tests/fixtures/project_requests.json`

```json
{
  "valid_project": {
    "project_id": "test-project-123",
    "display_name": "Test Project",
    "parent_id": "123456789",
    "parent_type": "organization",
    "labels": {
      "environment": "test",
      "team": "platform"
    }
  }
}
```

### Bucket Fixtures

**File**: `tests/fixtures/bucket_requests.json`

```json
{
  "valid_bucket": {
    "name": "test-bucket-123",
    "location": "us-central1",
    "storage_class": "STANDARD",
    "versioning": true
  }
}
```

## Mock Service Usage

### Setting Up Mocks

```go
func testCreateProject(t *testing.T, setup *TestSetup, token string) {
    // Setup mock expectations
    if !setup.Config.UseRealGCP && setup.MockService != nil {
        expectedResponse := mocks.NewMockProjectResponse(&req)
        setup.MockService.On("CreateProject", mock.AnythingOfType("*models.ProjectRequest")).Return(expectedResponse, nil)
    }

    // Execute test...

    // Reset mock expectations
    if !setup.Config.UseRealGCP && setup.MockService != nil {
        setup.MockService.ExpectedCalls = nil
        setup.MockService.Calls = nil
    }
}
```

### Mock Helper Functions

```go
// Create mock responses
projectResponse := mocks.NewMockProjectResponse(request)
folderResponse := mocks.NewMockFolderResponse(request)
bucketResponse := mocks.NewMockBucketResponse(request)
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Integration Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'

      # Mock tests (always run)
      - name: Run Integration Tests (Mock)
        run: make test-integration

      # Real GCP tests (only on main branch)
      - name: Run Integration Tests (Real GCP)
        if: github.ref == 'refs/heads/main'
        env:
          TEST_MODE: integration
          TEST_PROJECT_ID: ${{ secrets.TEST_PROJECT_ID }}
          GOOGLE_APPLICATION_CREDENTIALS: ${{ secrets.GCP_SA_KEY }}
        run: make test-integration-real
```

## Best Practices

### Test Organization

1. **Group Related Tests** - Use subtests for logical grouping
2. **Descriptive Names** - Use clear, descriptive test names
3. **Independent Tests** - Each test should be independent
4. **Cleanup Resources** - Always clean up test resources

### Mock Usage

1. **Reset Mocks** - Reset mock expectations between tests
2. **Specific Expectations** - Use specific mock expectations
3. **Error Testing** - Test both success and error scenarios
4. **Mock Helpers** - Use helper functions for common mock setups

### Real GCP Testing

1. **Separate Project** - Use a dedicated test project
2. **Resource Cleanup** - Implement proper cleanup
3. **Rate Limiting** - Be mindful of API rate limits
4. **Cost Management** - Monitor test resource costs

## Troubleshooting

### Common Issues

**Mock expectations not met**:
```bash
# Check mock setup and reset between tests
setup.MockService.ExpectedCalls = nil
setup.MockService.Calls = nil
```

**Authentication failures**:
```bash
# Verify JWT token generation
token := GenerateTestJWT(t, setup.AuthService)
req.Header.Set("Authorization", "Bearer "+token)
```

**Real GCP test failures**:
```bash
# Check credentials and project setup
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json
export TEST_PROJECT_ID=your-test-project
```

### Debug Mode

```bash
# Run with verbose output
go test -v ./tests/integration/ -run TestProjectOperations

# Run single test with debug
go test -v ./tests/integration/ -run TestProjectOperations/CreateProject/Valid_project_creation
```

## Extending Tests

### Adding New Test Cases

1. **Add Test Data** - Update fixtures with new test cases
2. **Implement Test** - Add test function with table-driven tests
3. **Mock Setup** - Add mock expectations for new scenarios
4. **Documentation** - Update this guide with new test information

### Adding New Resources

1. **Create Test File** - Add new test file (e.g., `folders_test.go`)
2. **Add Fixtures** - Create fixture file with test data
3. **Implement Mocks** - Add mock methods to `MockGCPService`
4. **Update Interface** - Add methods to `GCPServiceInterface`

## Performance Considerations

### Test Execution Time

- **Mock tests**: ~20ms per test suite
- **Real GCP tests**: ~2-5 seconds per operation
- **Parallel execution**: Tests run in parallel where possible

### Resource Usage

- **Mock mode**: Minimal resource usage
- **Real GCP mode**: Creates actual GCP resources (costs apply)
- **Cleanup**: Automatic cleanup prevents resource accumulation

## Security Considerations

### Credentials Management

- Never commit GCP credentials to version control
- Use environment variables for credential paths
- Rotate test service account keys regularly
- Limit test service account permissions

### Test Data

- Use non-sensitive test data only
- Avoid real customer data in tests
- Use randomized test resource names
- Clean up test resources promptly

---

For questions or issues with integration testing, please refer to the main project documentation or create an issue in the repository.
