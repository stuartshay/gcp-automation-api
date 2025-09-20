# Copilot Environment Setup

## Overview

This document outlines the setup of the GitHub Actions environment "copilot" with secure access to
Google Cloud Platform resources.

## Environment Configuration

### Created Environment

- **Name**: `copilot`
- **Purpose**: Dedicated environment for Copilot-related CI/CD operations
- **Repository**: `stuartshay/gcp-automation-api`

### Environment Secrets

The following secrets have been configured for the `copilot` environment:

#### GCP_SA_KEY

- **Description**: Google Cloud service account key for CI/CD operations
- **Service Account**: `github-actions-ci@gcp-auto-api-250913.iam.gserviceaccount.com`
- **Permissions**:
  - Artifact Registry Writer
  - Storage Admin (for Artifact Registry)
  - Additional permissions as needed for deployment

## Security Best Practices

### Environment-Level Secrets

- Secrets are scoped to the `copilot` environment only
- Not accessible to other workflows or environments
- Provides better security isolation

### Service Account Permissions

- Follows principle of least privilege
- Only has permissions necessary for CI/CD operations
- Regular key rotation recommended

## Usage in Workflows

To use the copilot environment in GitHub Actions workflows:

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: copilot # This enables access to environment secrets
    steps:
      - name: Setup GCP Auth
        uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY }}
```

## Verification

### Check Environment Exists

```bash
gh api repos/stuartshay/gcp-automation-api/environments
```

### List Environment Secrets

```bash
gh secret list --env copilot
```

### Test Workflow Access

Create a test workflow that uses the `copilot` environment to verify secret access.

## Troubleshooting

### Common Issues

1. **Secret Not Found**
   - Verify the environment name is spelled correctly
   - Ensure the workflow specifies `environment: copilot`

2. **Permission Denied**
   - Check service account IAM permissions
   - Verify the service account key is valid

3. **Environment Not Found**
   - Confirm the environment was created successfully
   - Check repository permissions

4. **Firewall Rules Blocking metadata.google.internal**
   - **Problem**: GitHub Actions firewall blocks access to `metadata.google.internal`
   - **Symptoms**: Error message about DNS blocks for metadata.google.internal
   - **Solution**: Set environment variables to prevent GCP client initialization during build/test
     phases:

     ```yaml
     env:
       GOOGLE_APPLICATION_CREDENTIALS: ""
       GCP_PROJECT_ID: "mock-project"
     ```

   - **Explanation**: GCP clients attempt to use Application Default Credentials by accessing the
     metadata service. During build/test phases, we don't need real GCP access, so we disable it to
     avoid firewall blocks.

### GitHub Actions Firewall Configuration

If you encounter firewall blocking issues, ensure your workflow jobs are configured as follows:

#### Build/Test Jobs (No GCP Access Needed)

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    env:
      GOOGLE_APPLICATION_CREDENTIALS: ""
      GCP_PROJECT_ID: "mock-project"
    # ... rest of job
```

#### Deployment Jobs (GCP Access Required)

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: copilot
    steps:
      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY }}
      # ... rest of deployment steps
```

### Commands Used for Setup

```bash
# Create the copilot environment
gh api -X PUT repos/stuartshay/gcp-automation-api/environments/copilot --input - <<< '{}'

# Generate service account key
gcloud iam service-accounts keys create copilot-env-sa-key.json \
  --iam-account=github-actions-ci@gcp-auto-api-250913.iam.gserviceaccount.com

# Set environment secret
gh secret set GCP_SA_KEY --env copilot < copilot-env-sa-key.json

# Clean up key file
rm copilot-env-sa-key.json
```

## Next Steps

1. Update workflows to use the `copilot` environment
2. Test deployment using the new environment configuration
3. Monitor workflow runs for any authentication issues
4. Consider setting up additional environments (staging, production) with similar configuration

## References

- [GitHub Environments Documentation](https://docs.github.com/en/actions/deployment/targeting-different-environments/using-environments-for-deployment)
- [Google Cloud Authentication for GitHub Actions](https://github.com/google-github-actions/auth)
- [GitHub CLI Secrets Management](https://cli.github.com/manual/gh_secret)
