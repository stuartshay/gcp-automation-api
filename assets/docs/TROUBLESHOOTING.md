# Troubleshooting Guide

## Common Issues and Solutions

### golangci-lint not found

**Problem**: `golangci-lint` executable not found in PATH

**Symptoms**:

```bash
Executable `golangci-lint` not found
```

**Solution**:

1. Ensure Go tools are installed:

   ```bash
   ./install.sh
   ```

2. Add Go bin directory to PATH:

   ```bash
   export PATH=$PATH:$HOME/go/bin
   ```

3. For permanent fix, the following should be in your shell config:

   ```bash
   # In ~/.bashrc or ~/.zshrc
   export PATH=$PATH:$HOME/go/bin
   ```

4. Activate development environment:

   ```bash
   source activate-dev.sh
   ```

**Verification**:

```bash
which golangci-lint
golangci-lint version
```

### Pre-commit Hook Failures

**Problem**: Pre-commit hooks fail on first run

**Symptoms**:

- `swag init` shows "files were modified"
- `prettier` shows "files were modified"
- `end-of-file-fixer` shows "files were modified"

**Solution**: These are normal formatting fixes. Run pre-commit multiple times until all
auto-formatting stabilizes:

```bash
pre-commit run --all-files
pre-commit run --all-files  # Run again
```

### Go Tools Not in PATH

**Problem**: Go tools like `gosec`, `goimports` not found

**Solution**:

1. Install tools:

   ```bash
   go install github.com/securego/gosec/v2/cmd/gosec@latest
   go install golang.org/x/tools/cmd/goimports@latest
   ```

2. Ensure `$HOME/go/bin` is in PATH as described above

### Environment Variables Not Loaded

**Problem**: `.env` file variables not available

**Solution**:

1. Ensure `.env` file exists (copy from `.env.example`)
2. Use activation script:

   ```bash
   source activate-dev.sh
   ```

### Docker Group Permissions

**Problem**: Docker commands require sudo

**Solution**:

1. Add user to docker group (done by install.sh):

   ```bash
   sudo usermod -aG docker $USER
   ```

2. Log out and back in, or use:

   ```bash
   newgrp docker
   ```

## Development Workflow

### Recommended Development Setup

1. **Initial Setup**:

   ```bash
   ./install.sh
   cp .env.example .env
   # Edit .env with your GCP credentials
   ```

2. **Daily Development**:

   ```bash
   source activate-dev.sh
   make dev
   ```

3. **Before Committing**:

   ```bash
   pre-commit run --all-files
   go test ./...
   ```

### Useful Commands

- **Check all tools are available**:

  ```bash
  source activate-dev.sh
  ```

- **Run specific linter**:

  ```bash
  golangci-lint run
  gosec ./...
  ```

- **Format code**:

  ```bash
  gofmt -w .
  goimports -w .
  ```

- **Generate Swagger docs**:

  ```bash
  swag init
  ```

## PATH Configuration

The development environment requires these directories in PATH:

1. `/usr/local/go/bin` - Go compiler and standard tools
2. `$HOME/go/bin` - Go tools installed with `go install`
3. `.venv/bin` - Python virtual environment (when activated)

The `activate-dev.sh` script ensures all required paths are set correctly.
