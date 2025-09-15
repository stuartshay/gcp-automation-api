# CLI Authentication Tool

The GCP Automation API now uses a CLI-based authentication system instead of exposing authentication endpoints through the HTTP API. This approach provides enhanced security by keeping authentication logic local and not exposing sensitive operations over the network.

## Overview

The `auth-cli` tool handles all authentication operations including:
- Google OAuth login
- JWT token generation and management
- Token refresh
- User profile management
- Development test tokens

## Installation

Build the CLI tool from source:

```bash
# Using Makefile (recommended)
make build-auth-cli

# Or build manually
go build -o bin/auth-cli ./cmd/auth-cli

# Build both server and CLI
make build-all-binaries
```

## Configuration

The CLI tool uses the same environment variables as the main API server:

### Required Environment Variables

```bash
# Google OAuth Configuration
export GOOGLE_CLIENT_ID="your-google-client-id"
export GOOGLE_CLIENT_SECRET="your-google-client-secret"

# JWT Configuration
export JWT_SECRET="your-jwt-secret-key"
export JWT_EXPIRATION_HOURS="24"

# Optional
export ENVIRONMENT="development"  # or "production"
```

## Commands

### `auth-cli login`

Performs Google OAuth authentication using your browser.

```bash
./bin/auth-cli login
```

This command:
1. Opens your default browser to Google's OAuth page
2. Starts a local callback server on port 8085
3. Exchanges the OAuth code for a Google ID token
4. Generates a JWT token using the API's authentication service
5. Stores credentials locally in `~/.gcp-automation/credentials.json`

### `auth-cli status`

Shows current authentication status.

```bash
./bin/auth-cli status
```

Example output:
```
Status: Authenticated
User: John Doe (john@example.com)
Token Type: Bearer
Expires: 2025-09-15T16:50:56Z
Time remaining: 23h46m0s
```

### `auth-cli token`

Displays the current JWT access token.

```bash
./bin/auth-cli token
```

Use this token with API requests:
```bash
TOKEN=$(./bin/auth-cli token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/projects
```

### `auth-cli profile`

Shows detailed user profile information.

```bash
./bin/auth-cli profile
```

Example output:
```
User Profile:
  Name: John Doe
  Email: john@example.com
  ID: 123456789
  Picture: https://lh3.googleusercontent.com/...
  Token Expires: 2025-09-15T16:50:56Z
```

### `auth-cli refresh`

Refreshes the current JWT token to extend its expiration.

```bash
./bin/auth-cli refresh
```

### `auth-cli test-token`

Generates a test JWT token for development purposes (only works in non-production environments).

```bash
./bin/auth-cli test-token --user-id "test-123" --email "test@example.com" --name "Test User"
```

### `auth-cli logout`

Clears all stored authentication credentials.

```bash
./bin/auth-cli logout
```

## Usage with API

Once authenticated, use the token with API requests:

### Using curl
```bash
# Get token
TOKEN=$(./bin/auth-cli token)

# Make API request
curl -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     http://localhost:8080/api/v1/projects/my-project
```

### Using scripts
```bash
#!/bin/bash

# Check authentication status
if ! ./bin/auth-cli status | grep -q "Authenticated"; then
    echo "Please authenticate first: ./bin/auth-cli login"
    exit 1
fi

# Get token and make API call
TOKEN=$(./bin/auth-cli token)
curl -H "Authorization: Bearer $TOKEN" \
     -d '{"project_id":"my-new-project","display_name":"My Project"}' \
     -H "Content-Type: application/json" \
     http://localhost:8080/api/v1/projects
```

## Security Features

### Local Storage
- Credentials are stored in `~/.gcp-automation/credentials.json`
- File permissions are set to 0600 (owner read/write only)
- Credentials include token expiration tracking

### OAuth Security
- Uses PKCE (Proof Key for Code Exchange) flow
- State parameter validation prevents CSRF attacks
- Local callback server with automatic shutdown
- 5-minute timeout for authentication flow

### Token Management
- JWT tokens have configurable expiration (default 24 hours)
- Automatic expiration checking
- Secure token refresh without re-authentication
- Production environment restrictions for test tokens

## Troubleshooting

### Authentication Issues

**Browser doesn't open automatically:**
```bash
# The CLI will display the URL to visit manually
Opening browser for Google authentication...
If the browser doesn't open automatically, visit: https://accounts.google.com/o/oauth2/v2/auth?...
```

**Port 8085 already in use:**
The CLI uses port 8085 for the OAuth callback. If this port is busy, stop the conflicting service or wait for the authentication to timeout and try again.

**Google OAuth errors:**
- Verify `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` are correctly set
- Ensure the redirect URI `http://localhost:8085/callback` is configured in your Google OAuth application
- Check that the Google OAuth application is enabled

### Token Issues

**Token expired:**
```bash
./bin/auth-cli refresh
# or
./bin/auth-cli login
```

**Invalid token:**
```bash
./bin/auth-cli logout
./bin/auth-cli login
```

**Test token in production:**
Test tokens are disabled in production environments. Use the regular login flow instead.

### Configuration Issues

**Missing environment variables:**
```bash
# Check required variables are set
echo $GOOGLE_CLIENT_ID
echo $JWT_SECRET

# Set missing variables
export GOOGLE_CLIENT_ID="your-client-id"
export JWT_SECRET="your-secret-key"
```

## Migration from API Authentication

If you were previously using the HTTP API authentication endpoints:

1. **Remove API calls** to `/auth/login`, `/auth/refresh`, `/auth/profile`, and `/auth/test-token`
2. **Use CLI instead**:
   - Replace API login calls with `auth-cli login`
   - Replace token refresh calls with `auth-cli refresh`
   - Replace profile calls with `auth-cli profile`
   - Replace test token generation with `auth-cli test-token`
3. **Update scripts** to use `auth-cli token` to get the JWT for API requests
4. **Update CI/CD** pipelines to use the CLI for authentication

## Examples

### Complete Workflow
```bash
# 1. Build and run auth-cli (shows help)
make run-auth-cli

# 2. Authenticate
./bin/auth-cli login

# 3. Check status
./bin/auth-cli status

# 4. Get token for API use
TOKEN=$(./bin/auth-cli token)

# 5. Use token with API
curl -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"project_id":"test-project","display_name":"Test Project"}' \
     http://localhost:8080/api/v1/projects

# 6. Refresh token when needed
./bin/auth-cli refresh

# 7. Logout when done
./bin/auth-cli logout
```

### Development Workflow
```bash
# For development/testing, use test tokens
./bin/auth-cli test-token --user-id "dev-user" --email "dev@company.com" --name "Developer"

# Use the test token
TOKEN=$(./bin/auth-cli token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/projects
