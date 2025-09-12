# GCP Automation API

A REST API built with Go for automating GCP (Google Cloud Platform) resource management. This service provides endpoints for creating and managing GCP Projects, Folders, and Cloud Storage Buckets.

## Features

- **RESTful API** for GCP resource management
- **Enterprise-ready architecture** with proper separation of concerns
- **Comprehensive error handling** and logging
- **Docker containerization** for easy deployment
- **Unit testing infrastructure** 
- **Configuration management** via environment variables
- **Health check endpoints** for monitoring
- **Graceful shutdown** support

## Supported GCP Resources

- **Projects**: Create, retrieve, and delete GCP projects
- **Folders**: Organize resources with folder hierarchy
- **Storage Buckets**: Manage Cloud Storage buckets

## Project Structure

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models and DTOs
│   └── services/        # Business logic and GCP integration
├── pkg/                 # Public packages (reusable components)
├── api/v1/              # API specifications
├── configs/             # Configuration files
├── tests/               # Test files
├── scripts/             # Build and deployment scripts
├── docs/                # Documentation
├── Dockerfile           # Container definition
├── Makefile            # Build automation
└── go.mod              # Go module definition
```

## Quick Start

### Prerequisites

- Go 1.21 or later
- GCP Project with necessary APIs enabled
- Service Account with appropriate permissions
- Docker (optional, for containerization)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/stuartshay/gcp-automation-api.git
cd gcp-automation-api
```

2. Install dependencies:
```bash
make deps
```

3. Set up configuration:
```bash
cp configs/.env.example .env
# Edit .env with your GCP project details
```

4. Run the application:
```bash
make dev
```

### Configuration

Set the following environment variables:

```bash
# Server Configuration
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
ENABLE_DEBUG=true

# GCP Configuration
GCP_PROJECT_ID=your-gcp-project-id
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
GCP_REGION=us-central1
GCP_ZONE=us-central1-a
```

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Projects
- `POST /api/v1/projects` - Create a new project
- `GET /api/v1/projects/{id}` - Get project details
- `DELETE /api/v1/projects/{id}` - Delete a project

### Folders
- `POST /api/v1/folders` - Create a new folder
- `GET /api/v1/folders/{id}` - Get folder details
- `DELETE /api/v1/folders/{id}` - Delete a folder

### Buckets
- `POST /api/v1/buckets` - Create a new bucket
- `GET /api/v1/buckets/{name}` - Get bucket details
- `DELETE /api/v1/buckets/{name}` - Delete a bucket

## Usage Examples

### Create a Project

```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "my-new-project",
    "display_name": "My New Project",
    "parent_id": "123456789",
    "parent_type": "organization",
    "labels": {
      "environment": "production",
      "team": "platform"
    }
  }'
```

### Create a Storage Bucket

```bash
curl -X POST http://localhost:8080/api/v1/buckets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-storage-bucket",
    "location": "us-central1",
    "storage_class": "STANDARD",
    "labels": {
      "environment": "production"
    },
    "versioning": true
  }'
```

## Development

### Build Commands

```bash
# Build the application
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Run development server
make dev

# Format code
make fmt

# Lint code
make lint

# Build Docker image
make docker
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./tests/handlers_test.go
```

## Docker

### Build and Run

```bash
# Build Docker image
make docker

# Run with Docker
make docker-run
```

### Using Docker Compose

```yaml
version: '3.8'
services:
  gcp-automation-api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GCP_PROJECT_ID=your-project-id
      - PORT=8080
    volumes:
      - ./path/to/service-account.json:/app/service-account.json
    env_file:
      - .env
```

## Security Considerations

- Service account keys should be stored securely
- Use IAM roles with minimal required permissions
- Enable audit logging for GCP resource changes
- Implement rate limiting for production use
- Use HTTPS in production environments

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
- Create an issue in the GitHub repository
- Check the documentation in the `docs/` directory