# API Documentation

## Overview

The GCP Automation API provides RESTful endpoints for managing Google Cloud Platform resources including Projects, Folders, and Cloud Storage Buckets.

## Interactive API Documentation

The API includes **Swagger UI** for interactive documentation and testing:

- **Swagger UI**: [http://localhost:8090/swagger/index.html](http://localhost:8090/swagger/index.html)
- **Swagger JSON**: [http://localhost:8090/swagger/doc.json](http://localhost:8090/swagger/doc.json)

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

The API uses Google Cloud Service Account authentication. Ensure you have:

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

The API implements basic rate limiting to prevent abuse. In production, consider implementing more sophisticated rate limiting based on your needs.

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

1. **Project IDs**: Must be globally unique, 6-30 characters, lowercase letters, digits, and hyphens only
2. **Bucket Names**: Must be globally unique, 3-63 characters, follow DNS naming conventions
3. **Labels**: Use consistent labeling strategy for resource organization and cost tracking
4. **Error Handling**: Always check HTTP status codes and error messages
5. **Idempotency**: Some operations may not be idempotent; check existing resources before creation

## Monitoring

The API includes a health check endpoint at `/health` that can be used for monitoring and load balancer health checks.

## Logs

Application logs include structured JSON logging with request IDs for tracing and debugging.
