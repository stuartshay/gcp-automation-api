# Cloud Run Deployment Guide

## Overview

This document describes the Cloud Run deployment setup for the GCP Automation API. The API is
automatically deployed to Google Cloud Run using a CI/CD pipeline that builds Docker images and
deploys them to a managed Cloud Run service.

## Cloud Run Service Details

### Service Configuration

- **Service Name**: `gcp-automation-api`
- **Region**: `us-central1`
- **Platform**: Managed
- **URL**: <https://gcp-automation-api-902997681858.us-central1.run.app>
- **Authentication**: Allow unauthenticated requests
- **Port**: 8080 (automatically configured by Cloud Run)

### Resource Allocation

- **Memory**: 1 GiB
- **CPU**: 1 vCPU
- **Min Instances**: 0 (scale to zero when not in use)
- **Max Instances**: 10
- **Concurrency**: Default (1000 concurrent requests per instance)

### Environment Variables

The following environment variables are automatically configured for the Cloud Run service:

- `ENVIRONMENT=production`
- `LOG_LEVEL=info`
- `GCP_PROJECT_ID=gcp-auto-api-250913`
- `PORT` (automatically set by Cloud Run)

## Deployment Pipeline

### Workflow Overview

The deployment is triggered automatically when code is pushed to the `master` branch. The workflow
consists of:

1. **Lint and Test**: Code quality checks and unit tests
2. **Build**: Compile the Go binary
3. **Docker Build**: Create and push Docker image to Artifact Registry
4. **Deploy**: Deploy the new image to Cloud Run
5. **Verify**: Test the deployed service

### Deployment Job

```yaml
deploy:
  name: Deploy to Cloud Run
  runs-on: ubuntu-latest
  needs: [docker]
  if: github.event_name == 'push' && github.ref == 'refs/heads/master'
  environment: copilot
```

### Authentication

The deployment uses the `copilot` environment with the `GCP_SA_KEY` secret containing service
account credentials with the following permissions:

- `roles/run.admin` - Deploy and manage Cloud Run services
- `roles/iam.serviceAccountUser` - Act as service accounts
- `roles/artifactregistry.writer` - Pull Docker images from Artifact Registry

## API Endpoints

Once deployed, the following endpoints are available:

### Health Check

```bash
curl https://gcp-automation-api-902997681858.us-central1.run.app/health
```

**Response:**

```json
{
  "status": "healthy"
}
```

### API Documentation

- **Swagger UI**: <https://gcp-automation-api-902997681858.us-central1.run.app/swagger/index.html>
- **OpenAPI Spec**: Available through the Swagger UI

### Available Endpoints

- **Base URL**: `https://gcp-automation-api-902997681858.us-central1.run.app/api/v1`

#### Projects

- `POST /api/v1/projects` - Create GCP project
- `GET /api/v1/projects/{id}` - Get project details
- `DELETE /api/v1/projects/{id}` - Delete project

#### Folders

- `POST /api/v1/folders` - Create GCP folder
- `GET /api/v1/folders/{id}` - Get folder details
- `DELETE /api/v1/folders/{id}` - Delete folder

#### Storage Buckets

- `POST /api/v1/buckets` - Create Cloud Storage bucket
- `GET /api/v1/buckets/{name}` - Get bucket details
- `DELETE /api/v1/buckets/{name}` - Delete bucket

## Manual Deployment

### Prerequisites

1. **Authenticate with Google Cloud**:

   ```bash
   gcloud auth login
   gcloud config set project gcp-auto-api-250913
   ```

2. **Configure Docker for Artifact Registry**:

   ```bash
   gcloud auth configure-docker us-central1-docker.pkg.dev
   ```

### Deploy from CLI

```bash
# Deploy the latest image
gcloud run deploy gcp-automation-api \
  --image=us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api:latest \
  --platform=managed \
  --region=us-central1 \
  --allow-unauthenticated \
  --port=8080 \
  --set-env-vars="ENVIRONMENT=production,LOG_LEVEL=info,GCP_PROJECT_ID=gcp-auto-api-250913" \
  --memory=1Gi \
  --cpu=1 \
  --max-instances=10 \
  --min-instances=0 \
  --project=gcp-auto-api-250913
```

### Deploy Specific Version

```bash
# Deploy a specific image tag
gcloud run deploy gcp-automation-api \
  --image=us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api:master-<commit-sha> \
  --region=us-central1 \
  --project=gcp-auto-api-250913
```

## Monitoring and Logging

### Cloud Run Logs

View service logs using the Google Cloud Console or CLI:

```bash
# View recent logs
gcloud logs read "resource.type=cloud_run_revision AND resource.labels.service_name=gcp-automation-api" \
  --project=gcp-auto-api-250913 \
  --limit=50

# Follow logs in real-time
gcloud logs tail "resource.type=cloud_run_revision AND resource.labels.service_name=gcp-automation-api" \
  --project=gcp-auto-api-250913
```

### Service Metrics

Monitor service performance in the Google Cloud Console:

1. Go to Cloud Run in the Google Cloud Console
2. Select the `gcp-automation-api` service
3. View metrics for:
   - Request count and latency
   - Instance count
   - Memory and CPU usage
   - Error rates

### Health Monitoring

Set up monitoring alerts for:

- **Health Check Failures**: Monitor `/health` endpoint
- **High Error Rates**: 4xx/5xx response codes
- **High Latency**: Response time > threshold
- **Instance Scaling**: Unusual scaling patterns

## Security Configuration

### Network Security

- **HTTPS Only**: All traffic is automatically encrypted with TLS
- **No VPC**: Service runs in Google's default network (public internet)
- **Firewall**: No additional firewall rules needed (Cloud Run handles this)

### Authentication Options

#### Current: Allow Unauthenticated

The service currently allows unauthenticated requests for development and testing.

#### Enable Authentication (Optional)

To require authentication:

```bash
# Remove unauthenticated access
gcloud run services remove-iam-policy-binding gcp-automation-api \
  --member="allUsers" \
  --role="roles/run.invoker" \
  --region=us-central1 \
  --project=gcp-auto-api-250913

# Add specific users/service accounts
gcloud run services add-iam-policy-binding gcp-automation-api \
  --member="user:example@domain.com" \
  --role="roles/run.invoker" \
  --region=us-central1 \
  --project=gcp-auto-api-250913
```

## Troubleshooting

### Common Issues

#### 1. Deployment Fails

**Symptoms**: Deployment job fails in GitHub Actions

**Solutions**:

- Check service account permissions
- Verify Docker image exists in Artifact Registry
- Check Cloud Run API is enabled
- Review deployment logs

#### 2. Service Doesn't Start

**Symptoms**: Cloud Run service shows errors, container fails to start

**Solutions**:

- Check container logs in Cloud Console
- Verify PORT environment variable configuration
- Test Docker image locally
- Check resource allocation (memory/CPU)

#### 3. Health Check Failures

**Symptoms**: `/health` endpoint returns errors

**Solutions**:

- Check application logs
- Verify database connections (if any)
- Check GCP credentials configuration
- Review environment variables

#### 4. Permission Denied Errors

**Symptoms**: API returns 403 errors for GCP operations

**Solutions**:

- Check service account running the Cloud Run service
- Verify IAM permissions for GCP resources
- Ensure proper credential configuration

### Debugging Commands

```bash
# Check service status
gcloud run services describe gcp-automation-api \
  --region=us-central1 \
  --project=gcp-auto-api-250913

# View service configuration
gcloud run services describe gcp-automation-api \
  --region=us-central1 \
  --project=gcp-auto-api-250913 \
  --format="export"

# Test local container
docker run -p 8080:8080 \
  -e ENVIRONMENT=development \
  -e LOG_LEVEL=debug \
  us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api:latest
```

## Cost Optimization

### Scaling Configuration

The current configuration is optimized for cost:

- **Scale to Zero**: Service scales down to 0 instances when not in use
- **Minimal Resources**: 1 GiB memory, 1 vCPU is sufficient for the API
- **Request-based Billing**: Only pay for actual request processing time

### Cost Monitoring

Monitor costs in Google Cloud Console:

1. Go to Billing → Reports
2. Filter by:
   - Service: Cloud Run
   - Project: gcp-auto-api-250913
   - Time range: Last 30 days

### Cost Reduction Tips

1. **Optimize Cold Starts**: Use warm-up requests for production
2. **Resource Right-sizing**: Monitor actual usage and adjust memory/CPU
3. **Request Batching**: Combine multiple operations when possible
4. **Caching**: Implement caching to reduce processing time

## Next Steps

### Production Readiness

1. **Enable Authentication**: Configure proper IAM for production use
2. **Custom Domain**: Set up custom domain with SSL certificate
3. **Monitoring**: Set up comprehensive monitoring and alerting
4. **Backup Strategy**: Implement backup for configuration and data
5. **Disaster Recovery**: Plan for service recovery procedures

### Performance Optimization

1. **Connection Pooling**: Optimize database connections
2. **Caching**: Implement Redis or Memorystore for caching
3. **CDN**: Use Cloud CDN for static assets
4. **Load Testing**: Perform load testing to validate performance

### Security Enhancements

1. **VPC Integration**: Move to private VPC if needed
2. **Firewall Rules**: Implement additional network security
3. **Secret Management**: Use Secret Manager for sensitive data
4. **Audit Logging**: Enable comprehensive audit logging

## References

- [Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Cloud Run Pricing](https://cloud.google.com/run/pricing)
- [Cloud Run Best Practices](https://cloud.google.com/run/docs/best-practices)
- [Cloud Run Security](https://cloud.google.com/run/docs/securing)
- [GitHub Actions for Google Cloud](https://github.com/google-github-actions)
