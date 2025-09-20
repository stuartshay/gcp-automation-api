# GCP Automation API

[![Pre-commit Checks](https://github.com/stuartshay/gcp-automation-api/actions/workflows/pre-commit.yml/badge.svg)](https://github.com/stuartshay/gcp-automation-api/actions/workflows/pre-commit.yml)

A Go-based REST API for automating Google Cloud Platform operations with **JWT authentication**.

Demo:
[https://gcp-automation-api-902997681858.us-central1.run.app](https://gcp-automation-api-902997681858.us-central1.run.app/swagger/index.html#/)

## 🔐 Authentication

The API uses **CLI-based authentication** for enhanced security. Authentication is handled by the
`auth-cli` tool, not through HTTP endpoints.

### Quick Start with CLI Authentication

```bash
# Build the CLI authentication tool
make build-auth-cli

# Generate a test JWT token (development only)
./bin/auth-cli test-token --user-id "test-user" --email "test@example.com" --name "Test User"

# Get the token for API requests
TOKEN=$(./bin/auth-cli token)

# Use the token in API requests
curl -X GET http://localhost:8080/api/v1/projects/my-project \
  -H "Authorization: Bearer $TOKEN"
```

### Production Authentication

```bash
# Authenticate with Google OAuth (requires GOOGLE_CLIENT_ID setup)
./bin/auth-cli login

# Check authentication status
./bin/auth-cli status

# Use token with API
TOKEN=$(./bin/auth-cli token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/projects
```

### Swagger UI

1. Get your token: `./bin/auth-cli token`
2. Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
3. Click **"Authorize"** button
4. Enter: `Bearer YOUR_JWT_TOKEN`
5. Test endpoints interactively

See [CLI Authentication Documentation](./assets/docs/CLI_AUTHENTICATION.md) for complete
authentication details.

## Quick Start

```bash
# Fast, non-interactive bootstrap (installs Go/Python tools, writes .env)
./scripts/setup-dev-env.sh --project YOUR_GCP_PROJECT_ID

# Activate development environment
source activate-dev.sh

# (Optional) Inspect or tweak the generated .env file
${EDITOR:-nano} .env

# Run the development server
make dev
```

> **Need a full workstation setup?** The legacy `./install.sh` script is still available if you need
> system packages such as Docker or the Google Cloud SDK installed globally.

## Documentation

📚 **All documentation is located in the [`docs/`](./docs/) folder.**

- **[Project Documentation](./assets/docs/PROJECT_README.md)** - Complete project overview,
  architecture, and detailed setup
- **[API Documentation](./docs/API.md)** - REST API endpoints and usage
- **[Setup Documentation](./docs/)** - Installation and configuration guides

## Project Structure

```text
├── cmd/
│   ├── server/          # API server entry point
│   └── auth-cli/        # CLI authentication tool
├── internal/            # Private application code
├── pkg/                 # Public library code
├── api/v1/             # API specifications
├── docs/               # 📋 All documentation files
├── tests/              # Test files
├── bin/                # Built binaries
└── configs/            # Configuration files
```

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Build application
make build

# Run pre-commit hooks
pre-commit run --all-files
```

## Documentation Organization

> **Rule**: All markup files and documentation must be placed in the `docs/` folder.

See [`docs/README.md`](./docs/README.md) for complete documentation organization guidelines.

---

For detailed information, please refer to the [complete documentation](./docs/) in the `docs/`
folder.
