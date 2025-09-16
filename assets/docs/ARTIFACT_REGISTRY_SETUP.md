# Google Cloud Artifact Registry Setup

This document describes the setup for using Google Cloud Artifact Registry to store Docker images for the GCP Automation API project.

## 🏗️ Infrastructure Setup

### Artifact Registry Repository
- **Registry**: `us-central1-docker.pkg.dev`
- **Project**: `gcp-auto-api-250913`
- **Repository**: `gcp-automation-api`
- **Format**: Docker
- **Location**: `us-central1`

### Full Image Path
```
us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api
```

## 🔑 Authentication & Permissions

### Local Development
1. **Authenticate with gcloud**:
   ```bash
   gcloud auth login
   gcloud config set project gcp-auto-api-250913
   ```

2. **Configure Docker authentication**:
   ```bash
   gcloud auth configure-docker us-central1-docker.pkg.dev
   ```

3. **Test push access**:
   ```bash
   docker tag your-image:latest us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api:test
   docker push us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api:test
   ```

### GitHub Actions CI/CD
1. **Service Account**: `github-actions-ci@gcp-auto-api-250913.iam.gserviceaccount.com`
2. **IAM Role**: `roles/artifactregistry.writer`

## 🔐 GitHub Secrets Configuration

To enable the CI/CD pipeline, add the following secret to your GitHub repository:

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add the following secret:

**Secret Name**: `GCP_SA_KEY`
**Secret Value**: _(Copy the entire JSON content from the service account key created during setup)_

⚠️ **SECURITY NOTE**: The service account key JSON contains sensitive information. Never commit this to your repository. It should only be stored in GitHub Secrets.

## 🚀 CI/CD Pipeline

The GitHub Actions workflow (`.github/workflows/ci.yml`) has been configured to:

1. **Authenticate** with Google Cloud using the service account key
2. **Configure** Docker to use gcloud as credential helper
3. **Build** the Docker image
4. **Push** to Artifact Registry on pushes to `master` or `develop` branches

### Image Tags
The pipeline creates the following tags:
- `latest` (for master branch)
- `{branch-name}` (for branch pushes)
- `{branch-name}-{sha}` (for commit-specific versions)

## 📋 CLI Commands Used

### Enable APIs
```bash
gcloud services enable artifactregistry.googleapis.com
```

### Create Repository
```bash
gcloud artifacts repositories create gcp-automation-api \
  --repository-format=docker \
  --location=us-central1 \
  --description="Docker repository for GCP Automation API"
```

### Create Service Account
```bash
gcloud iam service-accounts create github-actions-ci \
  --display-name="GitHub Actions CI Service Account" \
  --description="Service account for GitHub Actions CI/CD pipeline"
```

### Grant Permissions
```bash
gcloud projects add-iam-policy-binding gcp-auto-api-250913 \
  --member="serviceAccount:github-actions-ci@gcp-auto-api-250913.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"
```

### Create Service Account Key
```bash
gcloud iam service-accounts keys create github-actions-key.json \
  --iam-account=github-actions-ci@gcp-auto-api-250913.iam.gserviceaccount.com
```

## 🔍 Verification

### List Images
```bash
gcloud artifacts docker images list us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api
```

### List Tags for Specific Image
```bash
gcloud artifacts docker tags list us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api
```

## 🛡️ Security Best Practices

1. **Service Account Key**: Store only in GitHub Secrets, never commit to repository
2. **Least Privilege**: Service account has only `artifactregistry.writer` role
3. **Repository Access**: Artifact Registry repository is private by default
4. **Key Rotation**: Regularly rotate service account keys

## 🐳 Docker Image Management

### Pull Image
```bash
docker pull us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api:latest
```

### Run Container
```bash
docker run -p 8080:8080 us-central1-docker.pkg.dev/gcp-auto-api-250913/gcp-automation-api/gcp-automation-api:latest
```

## ❓ Troubleshooting

### Authentication Issues
- Ensure `gcloud auth configure-docker us-central1-docker.pkg.dev` has been run
- Verify you have the correct IAM permissions
- Check that the project ID is correct

### CI/CD Issues
- Verify `GCP_SA_KEY` secret is properly configured in GitHub
- Check service account has `artifactregistry.writer` role
- Ensure Artifact Registry API is enabled

### Permission Denied
- Check IAM policy bindings: `gcloud projects get-iam-policy gcp-auto-api-250913`
- Verify service account exists: `gcloud iam service-accounts list`

For additional support, refer to the [Google Cloud Artifact Registry documentation](https://cloud.google.com/artifact-registry/docs).
