# GCP Automation API

[![Pre-commit Checks](https://github.com/stuartshay/gcp-automation-api/actions/workflows/pre-commit.yml/badge.svg)](https://github.com/stuartshay/gcp-automation-api/actions/workflows/pre-commit.yml)

A Go-based REST API for automating Google Cloud Platform operations.

Demo: [https://gcp-automation-api-902997681858.us-central1.run.app](https://gcp-automation-api-902997681858.us-central1.run.app/swagger/index.html#/)

## Quick Start

```bash
# Install dependencies and set up development environment
./install.sh

# Activate development environment
source activate-dev.sh

# Copy and configure environment variables
cp .env.example .env
# Edit .env with your GCP settings

# Run the development server
make dev
```

## Documentation

📚 **All documentation is located in the [`docs/`](./docs/) folder.**

- **[Project Documentation](./docs/PROJECT_README.md)** - Complete project overview, architecture, and detailed setup
- **[API Documentation](./docs/API.md)** - REST API endpoints and usage
- **[Setup Documentation](./docs/)** - Installation and configuration guides

## Project Structure

```text
├── cmd/server/          # Application entry point
├── internal/            # Private application code
├── pkg/                 # Public library code
├── api/v1/             # API specifications
├── docs/               # 📋 All documentation files
├── tests/              # Test files
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

For detailed information, please refer to the [complete documentation](./docs/) in the `docs/` folder.
