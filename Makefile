# Build variables
BINARY_NAME=gcp-automation-api
BINARY_PATH=./bin/$(BINARY_NAME)
AUTH_CLI_NAME=auth-cli
AUTH_CLI_PATH=./bin/$(AUTH_CLI_NAME)
GO_MODULE=github.com/stuartshay/gcp-automation-api

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

.PHONY: build build-auth-cli build-all-binaries clean test deps run run-auth-cli dev docker help

# Default target
all: clean deps test build-all-binaries

# Build the API server
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_PATH) ./cmd/server

# Build the auth CLI tool
build-auth-cli:
	@echo "Building $(AUTH_CLI_NAME)..."
	@mkdir -p bin
	$(GOBUILD) -ldflags="-w -s" -o $(AUTH_CLI_PATH) ./cmd/auth-cli

# Build both server and auth-cli
build-all-binaries: build build-auth-cli
	@echo "Built both $(BINARY_NAME) and $(AUTH_CLI_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/
	@rm -rf dist/

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./tests/...

# Run integration tests (mock mode)
test-integration:
	@echo "Running integration tests (mock mode)..."
	$(GOTEST) -v ./tests/integration/...

# Run integration tests with real GCP (requires credentials)
test-integration-real:
	@echo "Running integration tests with real GCP..."
	TEST_MODE=integration $(GOTEST) -v ./tests/integration/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./tests/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run integration tests with coverage
test-integration-coverage:
	@echo "Running integration tests with coverage..."
	$(GOTEST) -v -coverprofile=integration-coverage.out ./tests/integration/...
	$(GOCMD) tool cover -html=integration-coverage.out -o integration-coverage.html

# Run all tests (unit + integration)
test-all: test test-integration
	@echo "All tests completed"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application in development mode
dev:
	@echo "Running in development mode..."
	$(GOCMD) run ./cmd/server

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BINARY_PATH)

# Generate a development JWT token and print to terminal
generate-jwt: build-auth-cli
	@if [ -z "$$USER_ID" ] || [ -z "$$EMAIL" ] || [ -z "$$NAME" ]; then \
		echo "Usage: make generate-jwt USER_ID=<id> EMAIL=<email> NAME=<name>"; \
		exit 1; \
	fi
	@OUTPUT=`$(AUTH_CLI_PATH) test-token --user-id=$$USER_ID --email=$$EMAIL --name=$$NAME` ; \
	TOKEN=`echo "$$OUTPUT" | awk '/^Token:/ {print $2}'` ; \
	if [ -z "$$TOKEN" ]; then \
		echo "Failed to generate JWT token. CLI output:" ; \
		echo "$$OUTPUT" ; \
		exit 1; \
	else \
		echo "JWT Token: $$TOKEN" ; \
	fi
# Run the auth CLI tool
run-auth-cli: build-auth-cli
	@echo "Running $(AUTH_CLI_NAME)..."
	$(AUTH_CLI_PATH)

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	# Build API server for multiple platforms
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/server
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/server
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/server
	# Build auth-cli for multiple platforms
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o bin/$(AUTH_CLI_NAME)-linux-amd64 ./cmd/auth-cli
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o bin/$(AUTH_CLI_NAME)-darwin-amd64 ./cmd/auth-cli
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o bin/$(AUTH_CLI_NAME)-windows-amd64.exe ./cmd/auth-cli

	@echo "  generate-jwt             - Generate a development JWT token (usage: make generate-jwt USER_ID=<id> EMAIL=<email> NAME=<name>)"
# Lint the code
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format the code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install it with: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
	fi

# Docker build
docker:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file configs/.env $(BINARY_NAME):latest

# Help
help:
	@echo "Available targets:"
	@echo "  build                    - Build the API server"
	@echo "  build-auth-cli           - Build the auth CLI tool"
	@echo "  build-all-binaries       - Build both server and auth-cli"
	@echo "  clean                    - Clean build artifacts"
	@echo "  test                     - Run unit tests"
	@echo "  test-integration         - Run integration tests (mock mode)"
	@echo "  test-integration-real    - Run integration tests with real GCP"
	@echo "  test-coverage            - Run tests with coverage"
	@echo "  test-integration-coverage - Run integration tests with coverage"
	@echo "  test-all                 - Run all tests (unit + integration)"
	@echo "  deps                     - Install dependencies"
	@echo "  dev                      - Run in development mode"
	@echo "  run                      - Build and run the API server"
	@echo "  run-auth-cli             - Build and run the auth CLI tool"
	@echo "  build-all                - Build for multiple platforms (server + auth-cli)"
	@echo "  lint                     - Run linter"
	@echo "  fmt                      - Format code"
	@echo "  security                 - Check for security vulnerabilities"
	@echo "  docker                   - Build Docker image"
	@echo "  docker-run               - Run Docker container"
	@echo "  help                     - Show this help message"
