# Build variables
BINARY_NAME=gcp-automation-api
BINARY_PATH=./bin/$(BINARY_NAME)
GO_MODULE=github.com/stuartshay/gcp-automation-api

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

.PHONY: build clean test deps run dev docker help

# Default target
all: clean deps test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_PATH) ./cmd/server

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

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./tests/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

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

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/server
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/server
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/server

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
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps         - Install dependencies"
	@echo "  dev          - Run in development mode"
	@echo "  run          - Build and run the application"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  security     - Check for security vulnerabilities"
	@echo "  docker       - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help message"
