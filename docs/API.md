# API Documentation

## Overview

The GCP Automation API provides RESTful endpoints for managing Google Cloud Platform resources
including Projects, Folders, and Cloud Storage Buckets.

## Interactive API Documentation

The API includes **Swagger UI** for interactive documentation and testing:

- **Swagger UI**:
  [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Swagger JSON**: [http://localhost:8080/swagger/doc.json](http://localhost:8080/swagger/doc.json)

The Swagger UI provides:

- Complete API endpoint documentation
- Interactive request/response testing
- Model schemas and examples
- Authentication information

## Framework & Architecture

The API is built with:

- **Echo v4** - High performance web framework
- **Swaggo** - Automated Swagger documentation generation
- **Clean Architecture** - Separation of concerns with handlers, services, and models
- **Structured Logging** - JSON logging with request tracing

## Authentication

The API uses **JWT (JSON Web Token) authentication** with optional Google OAuth integration.

### JWT Authentication

All API endpoints (except `/health`) require JWT authentication:

1. **Obtain a JWT Token**: Use one of the authentication methods below
2. **Include in Requests**: Add the `Authorization: Bearer <token>` header to all requests
3. **Token Expiration**: Tokens are valid for 24 hours by default

### Authentication Methods

#### 1. Google OAuth Login (Production)

```bash
POST /auth/login
Content-Type: application/json

{
  "google_id_token": "your-google-id-token-here"
}
```

**Response:**

```json
{
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

#### 2. Test Token Generation (Development Only)

```bash
POST /auth/test-token
Content-Type: application/json

{
  "user_id": "test-user-123",
  "email": "test@example.com",
  "name": "Test User"
}
```

**Response:**

```json
{
  "message": "Test token generated successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

### Using JWT Tokens

Include the JWT token in the Authorization header for all API requests:

```bash
curl -X GET "http://localhost:8080/api/v1/projects/my-project" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Swagger UI Authentication

The Swagger UI includes a convenient **"Authorize" button** for JWT authentication:

1. Open Swagger UI:
   [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
2. Click the **"Authorize"** button (lock icon)
3. Enter your JWT token in the format: `Bearer <your-token>`
4. Click **"Authorize"**
5. All subsequent API calls will automatically include the token

### Authentication Endpoints

| Endpoint           | Method | Purpose                        | Auth Required |
| ------------------ | ------ | ------------------------------ | ------------- |
| `/auth/login`      | POST   | Google OAuth login             | No            |
| `/auth/test-token` | POST   | Generate test token (dev only) | No            |
| `/auth/refresh`    | POST   | Refresh JWT token              | Yes           |
| `/auth/profile`    | GET    | Get user profile               | Yes           |

### Environment Variables

Configure JWT authentication with these environment variables:

```bash
# JWT Configuration
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION_HOURS=24

# Google OAuth (optional)
ENABLE_GOOGLE_AUTH=true
GOOGLE_CLIENT_ID=your-google-client-id
```

### Error Responses

Authentication errors return standardized responses:

```json
{
  "error": "unauthorized",
  "message": "Invalid or missing JWT token",
  "code": 401
}
```

## GCP Service Account Authentication

For backend GCP operations, the API uses Google Cloud Service Account authentication:

1. A GCP Service Account with appropriate permissions
2. Service Account JSON key file
3. Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable

## Required GCP Permissions

### For Projects

- `resourcemanager.projects.create`
- `resourcemanager.projects.delete`
- `resourcemanager.projects.get`

### For Folders

- `resourcemanager.folders.create`
- `resourcemanager.folders.delete`
- `resourcemanager.folders.get`

### For Storage Buckets

- `storage.buckets.create`
- `storage.buckets.delete`
- `storage.buckets.get`

## Rate Limiting

The API implements basic rate limiting to prevent abuse. In production, consider implementing more
sophisticated rate limiting based on your needs.

## Error Handling

All errors follow the standard format:

```json
{
  "error": "Error Type",
  "message": "Detailed error message",
  "code": 400
}
```

## Request/Response Examples

### Create Project

**Request:**

```bash
POST /api/v1/projects
Content-Type: application/json

{
  "project_id": "my-new-project-123",
  "display_name": "My New Project",
  "parent_id": "123456789",
  "parent_type": "organization",
  "labels": {
    "environment": "production",
    "team": "platform",
    "cost-center": "engineering"
  }
}
```

**Response:**

```json
{
  "message": "Project created successfully",
  "data": {
    "project_id": "my-new-project-123",
    "display_name": "My New Project",
    "parent_id": "123456789",
    "parent_type": "organization",
    "state": "ACTIVE",
    "labels": {
      "environment": "production",
      "team": "platform",
      "cost-center": "engineering"
    },
    "create_time": "2023-12-01T10:00:00Z",
    "update_time": "2023-12-01T10:00:00Z",
    "project_number": 123456789012
  }
}
```

### Create Storage Bucket

**Request:**

```bash
POST /api/v1/buckets
Content-Type: application/json

{
  "name": "my-data-bucket-123",
  "location": "us-central1",
  "storage_class": "STANDARD",
  "labels": {
    "environment": "production",
    "purpose": "data-storage"
  },
  "versioning": true
}
```

**Response:**

```json
{
  "message": "Bucket created successfully",
  "data": {
    "name": "my-data-bucket-123",
    "location": "us-central1",
    "storage_class": "STANDARD",
    "labels": {
      "environment": "production",
      "purpose": "data-storage"
    },
    "versioning": true,
    "create_time": "2023-12-01T10:00:00Z",
    "update_time": "2023-12-01T10:00:00Z",
    "self_link": "https://www.googleapis.com/storage/v1/b/my-data-bucket-123"
  }
}
```

## Best Practices

1. **Project IDs**: Must be globally unique, 6-30 characters, lowercase letters, digits, and hyphens
   only
2. **Bucket Names**: Must be globally unique, 3-63 characters, follow DNS naming conventions
3. **Labels**: Use consistent labeling strategy for resource organization and cost tracking
4. **Error Handling**: Always check HTTP status codes and error messages
5. **Idempotency**: Some operations may not be idempotent; check existing resources before creation

## Monitoring

The API includes a health check endpoint at `/health` that can be used for monitoring and load
balancer health checks.

## Logs

Application logs include structured JSON logging with request IDs for tracing and debugging.
