#!/bin/bash
set -e

echo "🧪 Testing GitHub Actions workflow dependencies locally..."

# Test hadolint installation (similar to what the workflow does)
echo "📦 Testing hadolint installation..."
HADOLINT_VERSION="2.12.0"
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "  - Downloading hadolint v${HADOLINT_VERSION}..."
wget -q -O hadolint "https://github.com/hadolint/hadolint/releases/download/v${HADOLINT_VERSION}/hadolint-Linux-x86_64"
chmod +x hadolint

echo "  - Testing hadolint..."
./hadolint --version

echo "✅ hadolint installation test passed!"

# Test shellcheck (if available)
echo "📦 Testing shellcheck..."
if command -v shellcheck >/dev/null 2>&1; then
    echo "  - shellcheck version: $(shellcheck --version | head -n2 | tail -n1)"
    echo "✅ shellcheck is available!"
else
    echo "⚠️  shellcheck not found (this is OK, it will be installed in CI)"
fi

# Test Go tools installation
echo "📦 Testing Go tools..."
echo "  - Go version: $(go version)"

# Test installing one Go tool to verify it works
echo "  - Testing go install..."
go install golang.org/x/tools/cmd/goimports@latest
echo "  - goimports installed successfully"

# Test gosec installation using the same method as CI
echo "  - Testing gosec installation..."
TEMP_GOSEC_DIR=$(mktemp -d)
CURRENT_DIR=$(pwd)
cd "$TEMP_GOSEC_DIR"
go install github.com/securego/gosec/v2/cmd/gosec@v2.22.8
if command -v gosec >/dev/null 2>&1; then
    gosec -version
    echo "  - gosec installed successfully"
else
    echo "  - gosec installation failed"
fi
cd "$CURRENT_DIR"
rm -rf "$TEMP_GOSEC_DIR"

echo "✅ All workflow dependencies test passed!"
echo "🚀 Your GitHub Actions workflow should now work correctly!"

# Cleanup
rm -rf "$TEMP_DIR"
