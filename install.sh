#!/bin/bash

# GCP Automation API - Installation Script
# This script installs all necessary SDKs, tools, and prerequisites
# including Python virtual environment for pre-commit hooks

set -euo pipefail

# Configuration variables
GO_VERSION="1.24.7"  # Match the version in go.mod

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running on supported OS
check_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
        DISTRO=$(lsb_release -si 2>/dev/null || echo "Unknown")
        log_info "Detected OS: Linux ($DISTRO)"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
        log_info "Detected OS: macOS"
    else
        log_error "Unsupported OS: $OSTYPE"
        exit 1
    fi
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install system dependencies
install_system_deps() {
    log_info "Installing system dependencies..."

    if [[ "$OS" == "linux" ]]; then
        if command_exists apt-get; then
            sudo apt-get update
            sudo apt-get install -y \
                curl \
                wget \
                git \
                build-essential \
                ca-certificates \
                gnupg \
                lsb-release \
                python3 \
                python3-pip \
                python3-venv \
                python3-dev \
                software-properties-common \
                apt-transport-https
        elif command_exists yum; then
            sudo yum update -y
            sudo yum install -y \
                curl \
                wget \
                git \
                gcc \
                gcc-c++ \
                make \
                ca-certificates \
                python3 \
                python3-pip \
                python3-devel
        elif command_exists dnf; then
            sudo dnf update -y
            sudo dnf install -y \
                curl \
                wget \
                git \
                gcc \
                gcc-c++ \
                make \
                ca-certificates \
                python3 \
                python3-pip \
                python3-devel
        else
            log_error "Unsupported Linux distribution"
            exit 1
        fi
    elif [[ "$OS" == "macos" ]]; then
        if ! command_exists brew; then
            log_info "Installing Homebrew..."
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        fi
        brew update
        brew install curl wget git python3
    fi

    log_success "System dependencies installed"
}

# Install Go
install_go() {
    if command_exists go; then
        local current_version=$(go version | awk '{print $3}' | sed 's/go//')
        if [[ "$current_version" == "$GO_VERSION" ]]; then
            log_success "Go $GO_VERSION is already installed"
            return
        else
            log_warning "Go $current_version is installed, but $GO_VERSION is required"
        fi
    fi

    log_info "Installing Go $GO_VERSION..."

    local go_arch="amd64"
    if [[ $(uname -m) == "arm64" ]] || [[ $(uname -m) == "aarch64" ]]; then
        go_arch="arm64"
    fi

    local go_os="linux"
    if [[ "$OS" == "macos" ]]; then
        go_os="darwin"
    fi

    local go_package="go${GO_VERSION}.${go_os}-${go_arch}.tar.gz"
    local go_url="https://golang.org/dl/${go_package}"

    # Download and install Go
    cd /tmp
    wget "$go_url"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "$go_package"
    rm "$go_package"

    # Add Go to PATH if not already there
    if [[ ":$PATH:" != *":/usr/local/go/bin:"* ]]; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc 2>/dev/null || true
        export PATH=$PATH:/usr/local/go/bin
    fi

    # Add Go bin directory (for go install tools) to PATH
    local go_bin_dir="$(go env GOPATH)/bin"
    if [[ ":$PATH:" != *":$go_bin_dir:"* ]]; then
        echo "export PATH=\$PATH:$go_bin_dir" >> ~/.bashrc
        echo "export PATH=\$PATH:$go_bin_dir" >> ~/.zshrc 2>/dev/null || true
        export PATH=$PATH:$go_bin_dir
    fi

    log_success "Go $GO_VERSION installed successfully"
}

# Install Google Cloud SDK
install_gcloud() {
    if command_exists gcloud; then
        log_success "Google Cloud SDK is already installed"
        gcloud version
        return
    fi

    log_info "Installing Google Cloud SDK..."

    if [[ "$OS" == "linux" ]]; then
        # Add the Cloud SDK distribution URI as a package source
        echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list

        # Import the Google Cloud public key
        curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -

        # Update and install the Cloud SDK
        sudo apt-get update && sudo apt-get install -y google-cloud-cli

        # Install additional components
        sudo apt-get install -y google-cloud-cli-gke-gcloud-auth-plugin

    elif [[ "$OS" == "macos" ]]; then
        # Install via Homebrew
        brew install --cask google-cloud-sdk
    fi

    log_success "Google Cloud SDK installed successfully"
    log_info "Run 'gcloud init' to initialize and authenticate"
}

# Install Docker
install_docker() {
    if command_exists docker; then
        log_success "Docker is already installed"
        docker --version
        return
    fi

    log_info "Installing Docker..."

    if [[ "$OS" == "linux" ]]; then
        # Install Docker using official script
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        rm get-docker.sh

        # Add user to docker group
        sudo usermod -aG docker $USER

        # Install Docker Compose
        sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose

    elif [[ "$OS" == "macos" ]]; then
        log_info "Please install Docker Desktop from https://docs.docker.com/docker-for-mac/install/"
        log_warning "Manual installation required for macOS"
        return
    fi

    log_success "Docker installed successfully"
    log_warning "You may need to log out and back in for Docker group membership to take effect"
}

# Install Go tools and linters
install_go_tools() {
    log_info "Installing Go tools and linters..."

    # Ensure we have Go in PATH
    export PATH=$PATH:/usr/local/go/bin

    # Initialize Go modules if not exists
    if [[ ! -f "go.mod" ]]; then
        log_warning "go.mod not found, skipping go mod download"
    else
        # Ensure Go modules are available
        go mod download || log_warning "go mod download failed, continuing..."
    fi

    # Install commonly used Go tools
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/securego/gosec/v2/cmd/gosec@latest
    go install golang.org/x/tools/cmd/goimports@latest
    go install github.com/google/wire/cmd/wire@latest
    go install github.com/swaggo/swag/cmd/swag@latest

    log_success "Go tools installed successfully"
}

# Setup Python virtual environment for pre-commit
setup_python_venv() {
    log_info "Setting up Python virtual environment for pre-commit..."

    # Create .venv directory if it doesn't exist
    if [[ ! -d ".venv" ]]; then
        python3 -m venv .venv
        log_success "Python virtual environment created at .venv"
    else
        log_success "Python virtual environment already exists"
    fi

    # Activate virtual environment
    source .venv/bin/activate

    # Upgrade pip
    pip install --upgrade pip

    # Install pre-commit and other Python tools
    pip install pre-commit
    pip install black
    pip install flake8
    pip install isort
    pip install mypy
    pip install bandit

    log_success "Python tools installed in virtual environment"

    # Install Node.js and npm if not available (needed for markdownlint)
    if ! command_exists node; then
        log_info "Installing Node.js for markdown tools..."
        if [[ "$OS" == "linux" ]]; then
            curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
            sudo apt-get install -y nodejs
        elif [[ "$OS" == "macos" ]]; then
            brew install node
        fi
        log_success "Node.js installed"
    fi

    # Install markdownlint-cli and prettier globally
    npm install -g markdownlint-cli prettier

    log_success "Markdown tools installed globally"
}

# Setup pre-commit configuration
setup_precommit() {
    log_info "Setting up pre-commit configuration..."

    # Create markdownlint configuration
    cat > .markdownlint.json << 'EOF'
{
  "default": true,
  "MD013": false,
  "MD025": false,
  "MD022": {
    "lines_above": 1,
    "lines_below": 1
  },
  "MD024": {
    "siblings_only": true
  },
  "MD026": {
    "punctuation": ".,;:!?"
  },
  "MD029": {
    "style": "ordered"
  },
  "MD031": true,
  "MD032": true,
  "MD034": false,
  "MD036": false,
  "MD040": false,
  "MD041": false,
  "MD046": {
    "style": "fenced"
  }
}
EOF

    # Create markdownlint ignore file
    cat > .markdownlintignore << 'EOF'
# Ignore auto-generated files
docs/swagger.json
docs/swagger.yaml
docs/docs.go

# Ignore files with intentional long lines
assets/docs/CLI_AUTHENTICATION.md
assets/docs/COPILOT_ENVIRONMENT_SETUP.md

# Ignore GitHub Copilot chatmode files (special formatting)
.github/chatmodes/*.md
EOF

    log_success "Created markdownlint configuration files"

    # Create .pre-commit-config.yaml
    cat > .pre-commit-config.yaml << 'EOF'
# Pre-commit configuration for GCP Automation API
repos:
  # Go formatting and linting
  - repo: local
    hooks:
      - id: go-fmt
        name: go fmt
        entry: gofmt
        language: system
        args: [-w]
        files: \.go$

      - id: go-imports
        name: go imports
        entry: goimports
        language: system
        args: [-w]
        files: \.go$

      - id: go-vet
        name: go vet
        entry: go vet
        language: system
        args: [./...]
        files: \.go$
        pass_filenames: false

      - id: go-test
        name: go test
        entry: go test
        language: system
        args: [-v, ./tests/...]
        files: \.go$
        pass_filenames: false

      - id: golangci-lint
        name: golangci-lint
        entry: golangci-lint run
        language: system
        args: [--fix]
        files: \.go$
        pass_filenames: false

  # Security scanning
  - repo: local
    hooks:
      - id: gosec
        name: gosec security scan
        entry: gosec
        language: system
        args: [-quiet, ./...]
        files: \.go$
        pass_filenames: false

  # General file checks
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-merge-conflict
      - id: check-added-large-files
      - id: detect-private-key
      - id: check-case-conflict

  # Dockerfile linting
  - repo: https://github.com/hadolint/hadolint
    rev: v2.12.0
    hooks:
      - id: hadolint-docker
        args: [--ignore, DL3008, --ignore, DL3009]

  # Markdown linting and formatting
  - repo: https://github.com/igorshubovych/markdownlint-cli
    rev: v0.45.0
    hooks:
      - id: markdownlint
        name: markdownlint (check)
        args: [--config, .markdownlint.json]
        exclude: \.markdownlintignore$
      - id: markdownlint
        name: markdownlint-fix
        alias: markdownlint-fix
        args: [--fix, --config, .markdownlint.json]
        stages: [manual]
        exclude: \.markdownlintignore$

  # Enhanced Markdown formatting with Prettier
  - repo: https://github.com/pre-commit/mirrors-prettier
    rev: v4.0.0-alpha.8
    hooks:
      - id: prettier
        name: prettier (markdown)
        types: [markdown]
        args: [--prose-wrap, always, --print-width, "100"]
        exclude: ^docs/swagger\.yaml$

  # YAML formatting
  - repo: https://github.com/pre-commit/mirrors-prettier
    rev: v3.0.0
    hooks:
      - id: prettier
        types: [yaml]
EOF

    # Activate virtual environment and install pre-commit hooks
    source .venv/bin/activate
    pre-commit install

    log_success "Pre-commit hooks installed"
    log_info "Run 'pre-commit run --all-files' to check all files"
}

# Create development scripts
create_dev_scripts() {
    log_info "Creating development scripts and directories..."

    # Create logs directory
    mkdir -p logs
    log_success "Created logs directory"

    # Create activate script for easy venv activation
    cat > activate-dev.sh << 'EOF'
#!/bin/bash
# Activate development environment
source .venv/bin/activate
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

# Load environment variables from .env file
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    export $(grep -v '^#' .env | grep -v '^$' | xargs)
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
EOF
    chmod +x activate-dev.sh

    # Create environment template
    if [[ ! -f ".env" ]]; then
        cat > .env.example << 'EOF'
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

# Optional: Service Account Key (base64 encoded)
# GCP_SERVICE_ACCOUNT_KEY=

# Optional: Specific GCP APIs to enable
# GCP_APIS=compute.googleapis.com,storage-component.googleapis.com,cloudresourcemanager.googleapis.com
EOF
        log_success "Created .env.example template"
        log_info "Copy .env.example to .env and fill in your GCP credentials"
    fi
}

# Main installation function
main() {
    log_info "Starting GCP Automation API installation..."
    log_info "This script will install:"
    log_info "  - System dependencies"
    log_info "  - Go $GO_VERSION"
    log_info "  - Google Cloud SDK"
    log_info "  - Docker"
    log_info "  - Go development tools"
    log_info "  - Node.js and markdown tools (markdownlint, prettier)"
    log_info "  - Python virtual environment"
    log_info "  - Pre-commit hooks"
    echo

    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Installation cancelled"
        exit 0
    fi

    check_os
    install_system_deps
    install_go
    install_gcloud
    install_docker
    install_go_tools
    setup_python_venv
    setup_precommit
    create_dev_scripts

    echo
    log_success "Installation completed successfully!"
    echo
    log_info "Next steps:"
    log_info "1. Run 'gcloud init' to authenticate with Google Cloud"
    log_info "2. Copy .env.example to .env and configure your GCP settings"
    log_info "3. Run 'source activate-dev.sh' to activate the development environment"
    log_info "4. Run 'make deps' to install Go dependencies"
    log_info "5. Run 'make dev' to start the development server"
    echo
    log_info "For Docker users: You may need to log out and back in for Docker group membership to take effect"
}

# Run main function
main "$@"
