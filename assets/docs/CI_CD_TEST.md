# CI/CD Pipeline Test

This file was created to test the GitHub Actions CI/CD pipeline with Google Cloud Artifact Registry
integration.

**Test Date**: September 13, 2025 **Purpose**: Verify that the `GCP_SA_KEY` secret is properly
configured and the pipeline can authenticate with Google Cloud.

## Expected Results

- ✅ Lint and Test jobs should pass
- ✅ Docker build job should authenticate with Google Cloud
- ✅ Docker image should be pushed to Artifact Registry:
  `us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api`

## Artifact Registry Details

- **Registry**: us-central1-docker.pkg.dev
- **Project**: gcp-auto-api-250913
- **Repository**: gcp-automation-api
- **Image Name**: gcp-automation-api

If this test succeeds, the CI/CD pipeline is fully operational with Google Cloud Artifact Registry.
