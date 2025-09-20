# Swagger UI Authentication Guide

This document explains how to use the Swagger UI with the newly enabled Authorize button for testing
the GCP Automation API.

## Overview

The Swagger UI now includes an **Authorize** button that allows you to authenticate and test API
endpoints directly from the browser interface. This works in conjunction with the CLI-based
authentication system.

## How to Use

### Step 1: Authenticate with CLI Tool

First, you need to authenticate using the CLI tool to obtain a JWT token:

#### Option A: Google OAuth (Production)

```bash
# Authenticate with Google OAuth
./bin/auth-cli login

# Check authentication status
./bin/auth-cli status
```

#### Option B: Test Token (Development)

```bash
# Generate a test token for development
./bin/auth-cli test-token --user-id "test-123" --email "test@example.com" --name "Test User"

# Check authentication status
./bin/auth-cli status
```

### Step 2: Get Your JWT Token

Retrieve your current JWT token:

```bash
# Get the current token
./bin/auth-cli token
```

This will output something like:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidGVzdC0xMjMiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJuYW1lIjoiVGVzdCBVc2VyIiwiZ29vZ2xlX3N1YiI6InRlc3QtMTIzIiwiaXNzIjoiZ2NwLWF1dG9tYXRpb24tYXBpIiwic3ViIjoidGVzdC0xMjMiLCJhdWQiOlsiZ2NwLWF1dG9tYXRpb24tYXBpIl0sImV4cCI6MTc1Nzk1NTA1NiwibmJmIjoxNzU3ODY4NjU2LCJpYXQiOjE3NTc4Njg2NTZ9.hFFQiTonUUPRmNmTXmjdHUAtEh-IgdTonHYmY3dN0JU
```

### Step 3: Access Swagger UI

1. Start the API server:

   ```bash
   make run
   # or
   go run ./cmd/server
   ```

2. Open your browser and navigate to:

   ```
   http://localhost:8080/swagger/
   ```

### Step 4: Authorize in Swagger UI

1. **Click the "Authorize" button** in the top-right corner of the Swagger UI
2. **Enter your token** in the "Value" field using this format:

   ```
   Bearer YOUR_JWT_TOKEN_HERE
   ```

   For example:

   ```
   Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidGVzdC0xMjMiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJuYW1lIjoiVGVzdCBVc2VyIiwiZ29vZ2xlX3N1YiI6InRlc3QtMTIzIiwiaXNzIjoiZ2NwLWF1dG9tYXRpb24tYXBpIiwic3ViIjoidGVzdC0xMjMiLCJhdWQiOlsiZ2NwLWF1dG9tYXRpb24tYXBpIl0sImV4cCI6MTc1Nzk1NTA1NiwibmJmIjoxNzU3ODY4NjU2LCJpYXQiOjE3NTc4Njg2NTZ9.hFFQiTonUUPRmNmTXmjdHUAtEh-IgdTonHYmY3dN0JU
   ```

3. **Click "Authorize"** to save the token
4. **Click "Close"** to return to the API documentation

### Step 5: Test API Endpoints

Now you can test any API endpoint:

1. **Expand an endpoint** (e.g., "POST /projects")
2. **Click "Try it out"**
3. **Fill in the request body** with your data
4. **Click "Execute"**

The request will automatically include your authentication token in the `Authorization` header.

## Available Endpoints

All the following endpoints require authentication and can be tested through Swagger UI:

### Projects

- `POST /api/v1/projects` - Create a new GCP project
- `GET /api/v1/projects/{id}` - Get a project by ID
- `DELETE /api/v1/projects/{id}` - Delete a project

### Folders

- `POST /api/v1/folders` - Create a new GCP folder
- `GET /api/v1/folders/{id}` - Get a folder by ID
- `DELETE /api/v1/folders/{id}` - Delete a folder

### Buckets

- `POST /api/v1/buckets` - Create a new Cloud Storage bucket
- `GET /api/v1/buckets/{name}` - Get a bucket by name
- `DELETE /api/v1/buckets/{name}` - Delete a bucket

## Example: Creating a Project

1. **Authorize** using the steps above
2. **Navigate to** `POST /api/v1/projects`
3. **Click "Try it out"**
4. **Enter request body**:

   ```json
   {
     "project_id": "my-test-project-123",
     "display_name": "My Test Project",
     "parent_id": "123456789",
     "parent_type": "organization",
     "labels": {
       "environment": "test",
       "team": "platform"
     }
   }
   ```

5. **Click "Execute"**
6. **Review the response**

## Token Management

### Token Expiration

- JWT tokens expire after 24 hours by default
- When your token expires, you'll get a 401 Unauthorized response
- Use `auth-cli refresh` to refresh your token or `auth-cli login` to re-authenticate

### Checking Token Status

```bash
# Check if your token is still valid
./bin/auth-cli status

# Refresh an expired token
./bin/auth-cli refresh

# Get user profile information
./bin/auth-cli profile
```

### Logout

```bash
# Clear stored credentials
./bin/auth-cli logout
```

## Troubleshooting

### Common Issues

**401 Unauthorized Error**

- Check that you included "Bearer " before your token
- Verify your token hasn't expired with `auth-cli status`
- Try refreshing your token with `auth-cli refresh`

**Invalid Token Format**

- Ensure you're using the format: `Bearer YOUR_TOKEN_HERE`
- Make sure there's a space after "Bearer"
- Don't include any extra characters or line breaks

**Token Not Working**

- Generate a new token with `auth-cli login` or `auth-cli test-token`
- Copy the exact token output from `auth-cli token`
- Re-authorize in Swagger UI with the new token

### Getting Help

For CLI authentication issues, see:

- `../assets/docs/CLI_AUTHENTICATION.md` - Complete CLI authentication guide
- `./bin/auth-cli --help` - CLI help and commands

For API issues, see:

- `docs/API.md` - API documentation
- `../assets/docs/TROUBLESHOOTING.md` - General troubleshooting guide

## Security Notes

- **Never share your JWT tokens** - they provide full access to the API
- **Tokens are stored locally** in `~/.gcp-automation/credentials.json`
- **Use test tokens only in development** - never in production
- **Rotate tokens regularly** by re-authenticating

---

The Swagger UI with authentication provides a convenient way to test and explore the API while
maintaining the security of CLI-based authentication.
