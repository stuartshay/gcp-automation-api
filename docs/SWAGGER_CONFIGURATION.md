# Swagger Configuration Guide

This document explains how to configure the Swagger documentation for different deployment environments.

## Overview

The GCP Automation API uses dynamic Swagger configuration that allows you to set the correct host and scheme based on your deployment environment. This ensures that the Swagger UI displays the correct URLs for your API endpoints.

## Configuration

The Swagger configuration is controlled by two environment variables:

- `SWAGGER_HOST`: The hostname where your API is deployed
- `SWAGGER_SCHEME`: The protocol scheme (http or https)

## Environment-Specific Configuration

### Local Development

For local development, the default values are used:

```bash
SWAGGER_HOST=localhost:8080
SWAGGER_SCHEME=http
```

### Cloud Run Deployment

For your Cloud Run deployment, set these environment variables:

```bash
SWAGGER_HOST=gcp-automation-api-902997681858.us-central1.run.app
SWAGGER_SCHEME=https
```

## Setting Environment Variables in Cloud Run

### Option 1: Using gcloud CLI

```bash
gcloud run services update gcp-automation-api \
  --set-env-vars="SWAGGER_HOST=gcp-automation-api-902997681858.us-central1.run.app,SWAGGER_SCHEME=https" \
  --region=us-central1
```

### Option 2: Using Cloud Console

1. Go to Cloud Run in the Google Cloud Console
2. Select your service: `gcp-automation-api`
3. Click "Edit & Deploy New Revision"
4. Go to the "Variables & Secrets" tab
5. Add the following environment variables:
   - `SWAGGER_HOST`: `gcp-automation-api-902997681858.us-central1.run.app`
   - `SWAGGER_SCHEME`: `https`
6. Click "Deploy"

### Option 3: Using Cloud Run YAML

Add the environment variables to your Cloud Run service configuration:

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: gcp-automation-api
spec:
  template:
    spec:
      containers:
      - image: gcr.io/your-project/gcp-automation-api
        env:
        - name: SWAGGER_HOST
          value: "gcp-automation-api-902997681858.us-central1.run.app"
        - name: SWAGGER_SCHEME
          value: "https"
```

## How It Works

1. The application reads the `SWAGGER_HOST` and `SWAGGER_SCHEME` environment variables at startup
2. When a request is made to `/swagger/doc.json`, the custom handler intercepts it
3. The handler modifies the Swagger JSON response to use the configured host and scheme
4. The Swagger UI displays the correct URLs for testing API endpoints

## Verification

After deploying with the correct environment variables:

1. Visit your Swagger UI: `https://gcp-automation-api-902997681858.us-central1.run.app/swagger/index.html`
2. Check that the "Servers" dropdown shows only HTTPS with your correct domain
3. Verify that the "Try it out" functionality uses the correct URLs

## Troubleshooting

### Issue: Swagger still shows localhost

**Solution**: Ensure the environment variables are properly set and the service has been redeployed.

### Issue: Both HTTP and HTTPS schemes appear

**Solution**: Verify that `SWAGGER_SCHEME=https` is set (not `SWAGGER_SCHEME=http,https`).

### Issue: Wrong hostname in Swagger UI

**Solution**: Double-check the `SWAGGER_HOST` environment variable matches your actual Cloud Run service URL.

## Security Considerations

- Always use `https` for production deployments
- The Swagger UI will only show the configured scheme, improving security by preventing accidental HTTP requests in production
- Ensure your Cloud Run service is configured to redirect HTTP to HTTPS
