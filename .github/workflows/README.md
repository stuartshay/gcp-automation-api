# GitHub Actions Workflows

This directory contains GitHub Actions workflows for the GCP Automation API project.

## Workflows

### 1. Pre-commit Checks (`pre-commit.yml`)

**Triggers:**

- Push to `master` or `develop` branches
- Pull requests to `master` or `develop` branches

**Purpose:** Runs pre-commit hooks on all files to ensure code quality and consistency.

**What it does:**

- Sets up Go, Python, and Node.js environments
- Installs all required tools (golangci-lint, gosec, swag, markdownlint, etc.)
- Runs pre-commit hooks defined in `.pre-commit-config.yaml`
- Uploads artifacts on failure for debugging

### 2. Continuous Integration (`ci.yml`)

**Triggers:**

- Push to `master` or `develop` branches
- Pull requests to `master` or `develop` branches

**Purpose:** Comprehensive CI pipeline with linting, testing, building, and security scanning.

**Jobs:**

- **Lint:** Code quality checks with golangci-lint and gosec
- **Test:** Unit tests with coverage reporting
- **Build:** Binary compilation and artifact upload
- **Docker:** Container image building and publishing (on push only)
- **Security:** Trivy vulnerability scanning

### 3. Dependency Updates (`dependency-update.yml`)

**Triggers:**

- Daily schedule (6:00 AM UTC)
- Manual workflow dispatch

**Purpose:** Automatically updates Go dependencies and pre-commit hooks.

**What it does:**

- Updates Go modules to latest versions
- Updates pre-commit hooks
- Creates pull requests for changes
- Runs tests to ensure updates don't break anything

### 4. Release (`release.yml`)

**Triggers:**

- Git tags matching `v*`
- Manual workflow dispatch with version input

**Purpose:** Creates releases with multi-platform binaries and Docker images.

**What it does:**

- Builds binaries for multiple platforms (Linux, macOS, Windows)
- Creates Docker images for multiple architectures
- Generates changelog
- Creates GitHub release with artifacts

## Required Secrets

The workflows use the following secrets:

- `GITHUB_TOKEN` - Automatically provided by GitHub
- `CODECOV_TOKEN` - (Optional) For coverage reporting to Codecov

## Branch Protection

It's recommended to set up branch protection rules for `master` and `develop` branches:

- Require status checks to pass before merging
- Require branches to be up to date before merging
- Include administrators in restrictions
- Required status checks:
  - `pre-commit`
  - `lint`
  - `test`
  - `build`

## Usage

### Running Pre-commit Locally

```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Run on all files
pre-commit run --all-files
```

### Creating a Release

1. **Automatic (recommended):**

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **Manual:**
   - Go to Actions → Release workflow
   - Click "Run workflow"
   - Enter version (e.g., `v1.0.0`)

### Viewing Build Status

- Check the Actions tab in the GitHub repository
- Each workflow run shows detailed logs and artifacts
- Failed runs include debug information and artifacts

## Customization

To modify the workflows:

1. Edit the respective `.yml` files in this directory
2. Test changes in a feature branch first
3. Workflows will automatically run on the next push/PR

## Troubleshooting

Common issues and solutions:

1. **Go tools not found:** The workflow installs tools in each run. If issues persist, check Go
   version compatibility.

2. **Pre-commit failures:** Run `pre-commit run --all-files` locally to debug issues before pushing.

3. **Docker build failures:** Ensure the Dockerfile is valid and all dependencies are available.

4. **Permission errors:** Check that the `GITHUB_TOKEN` has sufficient permissions for the
   repository.
