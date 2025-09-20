# VS Code Configuration Guide

## Overview

This project includes comprehensive VS Code settings configured for Go development with zsh as the
default terminal shell.

## Configuration Files

### `.vscode/settings.json`

- **Terminal**: Configured to use zsh as default shell
- **Go Extension**: Optimized settings for Go development
- **Editor**: Consistent formatting and linting
- **File Associations**: Proper syntax highlighting for project files
- **Python**: Configuration for pre-commit tools

### `.vscode/launch.json`

- **Launch Server**: Debug configuration for running the main server
- **Debug Tests**: Debug configuration for running tests
- **Debug Validators Tests**: Specific configuration for validator tests

### `.vscode/tasks.json`

- **Build Server**: Compile the Go application
- **Run Server**: Build and run the server with development settings
- **Run Tests**: Execute all tests with verbose output
- **Lint Code**: Run golangci-lint
- **Security Scan**: Run gosec security scanner
- **Generate Swagger Docs**: Run swag init
- **Pre-commit Check**: Run all pre-commit hooks
- **Docker Build**: Build Docker image

### `.vscode/extensions.json`

Recommended extensions for optimal development experience:

#### Essential Extensions

- **golang.go**: Go language support
- **ms-python.python**: Python support for pre-commit tools
- **ms-azuretools.vscode-docker**: Docker support
- **redhat.vscode-yaml**: YAML language support
- **davidanson.vscode-markdownlint**: Markdown linting

#### Development Tools

- **ms-python.flake8**: Python linting
- **ms-python.black-formatter**: Python code formatting
- **ms-python.isort**: Python import sorting
- **esbenp.prettier-vscode**: Code formatting
- **eamodio.gitlens**: Enhanced Git support

### `.vscode/go.code-snippets`

Custom code snippets for this project:

- `gcphandler`: Create GCP API handler with validation
- `validtest`: Create validation test function
- `customval`: Create custom validator function
- `errresp`: Return error response
- `succresp`: Return success response

### `gcp-automation-api.code-workspace`

Workspace configuration with:

- Terminal settings for zsh
- Go-specific paths and settings
- File watcher exclusions
- Extension recommendations

## Terminal Configuration

The configuration sets zsh as the default terminal with these features:

### Shell Profiles

```json
{
  "zsh": {
    "path": "/usr/bin/zsh",
    "args": ["-l"]
  },
  "bash": {
    "path": "/bin/bash",
    "args": ["-l"]
  }
}
```

### Environment Variables

- `WORKSPACE_ROOT`: Points to the workspace folder
- Development environment variables loaded from `.env`

## Go Development Features

### Code Intelligence

- **Auto-completion**: IntelliSense for Go code
- **Go to Definition**: Navigate to function/type definitions
- **Find References**: Find all references to symbols
- **Rename Symbol**: Refactor symbol names across codebase

### Code Quality

- **Format on Save**: Automatic code formatting with goimports
- **Lint on Save**: Real-time linting with golangci-lint
- **Organize Imports**: Automatic import organization
- **Error Detection**: Real-time error highlighting

### Testing Support

- **Test Explorer**: Run individual tests or test suites
- **Test Coverage**: View code coverage in editor
- **Debug Tests**: Set breakpoints and debug test code

### Build and Run

- **Build Tasks**: One-click build and run
- **Debug Configuration**: Step-through debugging
- **Hot Reload**: Automatic rebuild on file changes

## Usage Instructions

### Opening the Project

1. Open VS Code
2. File → Open Workspace from File
3. Select `gcp-automation-api.code-workspace`

### First Time Setup

1. Install recommended extensions when prompted
2. Run `Ctrl+Shift+P` → "Go: Install/Update Tools"
3. Activate development environment: `source activate-dev.sh`

### Development Workflow

#### Running the Server

- **Command Palette**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Run Server"
- **Terminal**: `make dev`
- **Debug**: Press `F5` to start debugging

#### Running Tests

- **Command Palette**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Run Tests"
- **Terminal**: `make test`
- **Debug**: Use "Debug Tests" launch configuration

#### Code Quality Checks

- **Lint**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Lint Code"
- **Security**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Security Scan"
- **Pre-commit**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Pre-commit Check"

### Keyboard Shortcuts

#### Go-Specific

- `F12`: Go to Definition
- `Shift+F12`: Find All References
- `F2`: Rename Symbol
- `Ctrl+Shift+O`: Go to Symbol in File
- `Ctrl+T`: Go to Symbol in Workspace

#### Build and Debug

- `Ctrl+Shift+B`: Run Build Task
- `F5`: Start Debugging
- `Ctrl+F5`: Run Without Debugging
- `Shift+F5`: Stop Debugging

#### Terminal

- `Ctrl+Shift+\``: New Terminal
- `Ctrl+\``: Toggle Terminal Panel

### Code Snippets Usage

Type the snippet prefix and press `Tab`:

- `gcphandler` + Tab: Create GCP API handler
- `validtest` + Tab: Create validation test
- `customval` + Tab: Create custom validator
- `errresp` + Tab: Create error response
- `succresp` + Tab: Create success response

## Customization

### Terminal Shell

To change the default shell, edit `.vscode/settings.json`:

```json
{
  "terminal.integrated.defaultProfile.linux": "bash" // or "zsh"
}
```

### Go Tools

To change Go tools, edit `.vscode/settings.json`:

```json
{
  "go.formatTool": "gofmt", // or "goimports"
  "go.lintTool": "staticcheck" // or "golangci-lint"
}
```

### Theme and Appearance

Install and configure themes via:

1. `Ctrl+Shift+X` → Search for themes
2. `Ctrl+Shift+P` → "Preferences: Color Theme"

## Troubleshooting

### Zsh Not Working

1. Verify zsh is installed: `which zsh`
2. Check VS Code settings for correct path
3. Restart VS Code after configuration changes

### Go Extension Issues

1. `Ctrl+Shift+P` → "Go: Install/Update Tools"
2. Check Go version: `go version`
3. Verify GOPATH and GOROOT settings

### Path Issues

1. Ensure `$HOME/go/bin` is in PATH
2. Run `source activate-dev.sh`
3. Check environment variables in terminal

### Performance Issues

1. Exclude large directories in file watcher settings
2. Disable unused extensions
3. Increase VS Code memory limit if needed

## Integration with Project Tools

### Pre-commit Hooks

VS Code tasks integrate with pre-commit hooks:

- Automatic formatting on save
- Real-time linting feedback
- One-click pre-commit execution

### Docker Support

- Dockerfile syntax highlighting and linting
- Build and run containers from VS Code
- Integrated Docker Explorer

### Git Integration

- Built-in Git support with GitLens
- Diff view and merge conflict resolution
- Branch management and history view

This configuration provides a complete, production-ready development environment for the GCP
Automation API project.
