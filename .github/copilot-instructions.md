# GitHub Copilot Instructions for GCP Automation API

## Project Overview

This is a **GCP Automation API** built with **Go 1.24.7** that provides RESTful endpoints for
automating Google Cloud Platform resource management. The service supports creating, retrieving, and
managing GCP Projects, Folders, and Cloud Storage Buckets.

## Architecture & Technology Stack

### **Backend Framework**

- **Language**: Go 1.24.7
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **GCP SDK**: Google Cloud Go SDK (cloud.google.com/go/\*)
- **Architecture**: Clean architecture with separation of concerns

### **Project Structure**

```text
├── cmd/server/          # Application entry point and main.go
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers (controllers)
│   ├── models/          # Data models and DTOs
│   └── services/        # Business logic and GCP service integration
├── pkg/                 # Public packages (reusable components)
├── api/v1/              # OpenAPI specifications
├── tests/               # Unit and integration tests
└── configs/             # Configuration files
```

### **Development Environment**

- **Environment Variables**: Loaded from `.env` file
- **Authentication**: Google Cloud Application Default Credentials (ADC)
- **Port**: 8080 (configurable via PORT env var)
- **Project ID**: velvety-byway-327718
- **Region**: us-central1

## API Endpoints & Patterns

### **RESTful API Structure**

Base URL: `http://localhost:8080/api/v1`

#### **Projects**

- `POST /api/v1/projects` - Create GCP project
- `GET /api/v1/projects/{id}` - Get project details
- `DELETE /api/v1/projects/{id}` - Delete project

#### **Folders**

- `POST /api/v1/folders` - Create GCP folder
- `GET /api/v1/folders/{id}` - Get folder details
- `DELETE /api/v1/folders/{id}` - Delete folder

#### **Storage Buckets**

- `POST /api/v1/buckets` - Create Cloud Storage bucket
- `GET /api/v1/buckets/{name}` - Get bucket details
- `DELETE /api/v1/buckets/{name}` - Delete bucket

#### **Health Check**

- `GET /health` - Service health status

## Code Style & Conventions

### **Go Code Standards**

- Follow **Go best practices** and **idiomatic Go** patterns
- Use **structured logging** with appropriate log levels
- Implement **proper error handling** with descriptive error messages
- Apply **clean architecture** principles with clear separation of concerns
- Use **dependency injection** for service components

### **Naming Conventions**

- **Handlers**: Use descriptive names like `CreateProject`, `GetBucket`
- **Services**: Business logic services like `GCPService`, `ProjectService`
- **Models**: Use clear model names like `ProjectRequest`, `BucketResponse`
- **Functions**: Use verb-noun patterns like `validateRequest`, `handleError`

### **Error Handling Patterns**

```go
// Prefer structured error responses
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}

// Use proper HTTP status codes
c.JSON(http.StatusBadRequest, ErrorResponse{
    Error:   "validation_failed",
    Message: "Invalid request parameters",
    Code:    400,
})
```

### **Request/Response Models**

```go
// Use proper validation tags
type CreateProjectRequest struct {
    ProjectID   string            `json:"project_id" binding:"required"`
    DisplayName string            `json:"display_name" binding:"required"`
    ParentID    string            `json:"parent_id"`
    ParentType  string            `json:"parent_type"` // "organization" or "folder"
    Labels      map[string]string `json:"labels"`
}
```

## GCP Integration Guidelines

### **Service Initialization**

- Initialize GCP clients in the `services` package
- Use **Application Default Credentials** for authentication
- Handle **context** properly for timeouts and cancellation
- Implement **proper resource cleanup**

### **Common GCP Patterns**

```go
// Example service pattern
type GCPService struct {
    projectID          string
    resourceManagerClient *resourcemanager.ProjectsClient
    storageClient      *storage.Client
}

// Context usage
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

### **GCP Resource Management**

- **Projects**: Use Cloud Resource Manager API
- **Folders**: Use Cloud Resource Manager API
- **Storage**: Use Cloud Storage API
- Always check **resource existence** before operations
- Implement **proper IAM handling**

## Testing Guidelines

### **Test Structure**

- **Unit tests**: Focus on business logic in `services`
- **Integration tests**: Test API endpoints with mock GCP services
- **Test files**: Use `*_test.go` naming convention
- Place tests in `tests/` directory

### **Test Patterns**

```go
// Use table-driven tests
func TestCreateProject(t *testing.T) {
    tests := []struct {
        name     string
        request  CreateProjectRequest
        expected int
        hasError bool
    }{
        // test cases
    }
}
```

## Configuration Management

### **Environment Variables**

```bash
# Server Configuration
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
ENABLE_DEBUG=true

# GCP Configuration
GCP_PROJECT_ID=velvety-byway-327718
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
GCP_REGION=us-central1
GCP_ZONE=us-central1-a
```

### **Configuration Loading**

- Load config in `internal/config/config.go`
- Use environment variables with sensible defaults
- Validate required configuration at startup

## Development Workflow

### **Getting Started**

```bash
# Setup environment
source activate-dev.sh

# Start development server
make dev

# Run tests
make test

# Build application
make build
```

### **Code Quality Tools**

- **Linting**: `golangci-lint run`
- **Security**: `gosec ./...`
- **Formatting**: `gofmt` and `goimports`
- **Pre-commit hooks**: Configured for automatic checks

## Security Considerations

### **Authentication & Authorization**

- Use **GCP IAM** for access control
- Validate **service account permissions**
- Implement **proper credential management**
- Never log sensitive information

### **Input Validation**

- Validate all request parameters
- Sanitize user input
- Use **Gin's binding** for automatic validation
- Implement **rate limiting** for production

### **GCP Best Practices**

- Follow **principle of least privilege**
- Use **resource-level IAM** when possible
- Implement **audit logging**
- Handle **quota limitations** gracefully

## Docker & Deployment

### **Containerization**

- **Multi-stage builds** for smaller images
- **Non-root user** for security
- **Health checks** in Dockerfile
- **Environment-based configuration**

### **Production Considerations**

- Set `GIN_MODE=release` for production
- Configure **proper logging** levels
- Implement **graceful shutdown**
- Use **structured configuration**

## Common Patterns & Examples

### **Handler Pattern**

```go
func (h *Handler) CreateProject(c *gin.Context) {
    var req CreateProjectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "invalid_request",
            Message: err.Error(),
        })
        return
    }

    // Business logic in service layer
    result, err := h.gcpService.CreateProject(c.Request.Context(), req)
    if err != nil {
        // Handle error appropriately
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error: "project_creation_failed",
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, result)
}
```

### **Service Pattern**

```go
func (s *GCPService) CreateProject(ctx context.Context, req CreateProjectRequest) (*ProjectResponse, error) {
    // Validate business rules
    if err := s.validateProjectRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Call GCP API
    project := &resourcemanagerpb.Project{
        ProjectId:   req.ProjectID,
        DisplayName: req.DisplayName,
        Labels:      req.Labels,
    }

    op, err := s.resourceManagerClient.CreateProject(ctx, &resourcemanagerpb.CreateProjectRequest{
        Project: project,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create project: %w", err)
    }

    // Handle long-running operation
    result, err := op.Wait(ctx)
    if err != nil {
        return nil, fmt.Errorf("project creation failed: %w", err)
    }

    return &ProjectResponse{
        ProjectID:   result.GetProjectId(),
        DisplayName: result.GetDisplayName(),
        State:       result.GetState().String(),
    }, nil
}
```

## Tips for GitHub Copilot

### **Context Awareness**

- When suggesting code, consider the **clean architecture** pattern
- Prefer **dependency injection** over global variables
- Always include **proper error handling**
- Consider **context cancellation** for GCP operations

### **GCP-Specific Suggestions**

- Use **official GCP Go SDK** patterns
- Include **proper resource cleanup** (defer statements)
- Consider **long-running operations** for GCP resources
- Implement **retry logic** for transient failures

### **Code Generation Preferences**

- Generate **complete functions** with error handling
- Include **proper logging** statements
- Use **structured JSON responses**
- Follow **Go conventions** for package organization

This documentation should help GitHub Copilot provide more accurate and contextually relevant
suggestions for your GCP Automation API project!
