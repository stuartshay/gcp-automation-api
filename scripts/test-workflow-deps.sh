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

echo "✅ All workflow dependencies test passed!"
echo "🚀 Your GitHub Actions workflow should now work correctly!"

# Cleanup
cd - >/dev/null
rm -rf "$TEMP_DIR"
