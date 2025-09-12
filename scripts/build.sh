#!/bin/bash

# Build script for GCP Automation API

set -e

echo "Starting build process..."

# Variables
BINARY_NAME="gcp-automation-api"
BUILD_DIR="./bin"
CMD_DIR="./cmd/server"

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# Download dependencies
echo "Downloading dependencies..."
go mod download
go mod tidy

# Run tests
echo "Running tests..."
go test -v ./tests/...

# Build the application
echo "Building application..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty)" \
  -o $BUILD_DIR/$BINARY_NAME \
  $CMD_DIR

# Make binary executable
chmod +x $BUILD_DIR/$BINARY_NAME

# Show build info
echo "Build completed successfully!"
echo "Binary: $BUILD_DIR/$BINARY_NAME"
echo "Size: $(du -h $BUILD_DIR/$BINARY_NAME | cut -f1)"

# Verify binary
if [ -x "$BUILD_DIR/$BINARY_NAME" ]; then
    echo "✓ Binary is executable"
else
    echo "✗ Binary is not executable"
    exit 1
fi

echo "Build process completed!"
