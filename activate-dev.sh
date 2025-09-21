#!/bin/bash
# Activate development environment
source .venv/bin/activate
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

# Load environment variables from .env file
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    export "$(grep -v '^#' .env | grep -v '^$' | xargs)"
fi

# Verify Go installation
if command -v go >/dev/null 2>&1; then
    echo "Go version: $(go version)"
else
    echo "❌ Go not found in PATH"
fi

# Verify Go tools
if command -v golangci-lint >/dev/null 2>&1; then
    echo "golangci-lint version: $(golangci-lint version | head -1)"
else
    echo "❌ golangci-lint not found in PATH"
fi

# Verify markdown tools
if command -v markdownlint >/dev/null 2>&1; then
    echo "markdownlint version: $(markdownlint --version)"
else
    echo "❌ markdownlint not found in PATH"
fi

if command -v prettier >/dev/null 2>&1; then
    echo "prettier version: $(prettier --version)"
else
    echo "❌ prettier not found in PATH"
fi

echo "Development environment activated!"
echo "Available commands:"
echo "  make dev          - Run development server"
echo "  make test         - Run tests"
echo "  make lint         - Run linter"
echo "  pre-commit run    - Run pre-commit hooks"
echo "  golangci-lint run - Run Go linter"
echo "  gosec ./...       - Run security scanner"
echo "  deactivate        - Exit virtual environment"
echo ""
echo "Environment:"
echo "  PORT=${PORT:-not set}"
echo "  GCP_PROJECT_ID=${GCP_PROJECT_ID:-not set}"
echo "  ENVIRONMENT=${ENVIRONMENT:-not set}"
