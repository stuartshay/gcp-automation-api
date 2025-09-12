package services

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/models"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

// GCPService handles all GCP operations
type GCPService struct {
	config             *config.Config
	resourceManager    *cloudresourcemanager.Service
	storageClient      *storage.Client
	ctx                context.Context
}

// NewGCPService creates a new GCP service instance
func NewGCPService(cfg *config.Config) (*GCPService, error) {
	ctx := context.Background()
	
	var opts []option.ClientOption
	if cfg.GCPCredentials != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.GCPCredentials))
	}

	// Initialize Resource Manager client
	resourceManager, err := cloudresourcemanager.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource manager client: %w", err)
	}

	// Initialize Storage client
	storageClient, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &GCPService{
		config:             cfg,
		resourceManager:    resourceManager,
		storageClient:      storageClient,
		ctx:                ctx,
	}, nil
}

// CreateProject creates a new GCP project
func (s *GCPService) CreateProject(req *models.ProjectRequest) (*models.ProjectResponse, error) {
	project := &cloudresourcemanager.Project{
		ProjectId:   req.ProjectID,
		Name:        req.DisplayName,
		Labels:      req.Labels,
	}

	// Set parent if specified
	if req.ParentID != "" && req.ParentType != "" {
		switch req.ParentType {
		case "organization":
			project.Parent = &cloudresourcemanager.ResourceId{
				Type: "organization",
				Id:   req.ParentID,
			}
		case "folder":
			project.Parent = &cloudresourcemanager.ResourceId{
				Type: "folder",
				Id:   req.ParentID,
			}
		default:
			return nil, fmt.Errorf("invalid parent type: %s", req.ParentType)
		}
	}

	// Create the project
	op, err := s.resourceManager.Projects.Create(project).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Wait for operation to complete (simplified - in production, use polling)
	time.Sleep(2 * time.Second)

	// Get the created project
	createdProject, err := s.resourceManager.Projects.Get(req.ProjectID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get created project: %w", err)
	}

	response := &models.ProjectResponse{
		ProjectID:     createdProject.ProjectId,
		DisplayName:   createdProject.Name,
		State:         createdProject.LifecycleState,
		Labels:        createdProject.Labels,
		ProjectNumber: createdProject.ProjectNumber,
		CreateTime:    time.Now(), // Simplified - should parse from API
		UpdateTime:    time.Now(),
	}

	if createdProject.Parent != nil {
		response.ParentID = createdProject.Parent.Id
		response.ParentType = createdProject.Parent.Type
	}

	// Store operation details (simplified)
	_ = op

	return response, nil
}

// GetProject retrieves a GCP project
func (s *GCPService) GetProject(projectID string) (*models.ProjectResponse, error) {
	project, err := s.resourceManager.Projects.Get(projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	response := &models.ProjectResponse{
		ProjectID:     project.ProjectId,
		DisplayName:   project.Name,
		State:         project.LifecycleState,
		Labels:        project.Labels,
		ProjectNumber: project.ProjectNumber,
		CreateTime:    time.Now(), // Simplified
		UpdateTime:    time.Now(),
	}

	if project.Parent != nil {
		response.ParentID = project.Parent.Id
		response.ParentType = project.Parent.Type
	}

	return response, nil
}

// DeleteProject deletes a GCP project
func (s *GCPService) DeleteProject(projectID string) error {
	_, err := s.resourceManager.Projects.Delete(projectID).Do()
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}

// CreateFolder creates a new GCP folder (placeholder implementation)
func (s *GCPService) CreateFolder(req *models.FolderRequest) (*models.FolderResponse, error) {
	// This is a placeholder implementation
	// In a real implementation, you would use the Cloud Resource Manager API
	// to create folders, which requires additional permissions and setup
	
	response := &models.FolderResponse{
		Name:        fmt.Sprintf("folders/%s", "generated-id"),
		DisplayName: req.DisplayName,
		ParentID:    req.ParentID,
		ParentType:  req.ParentType,
		State:       "ACTIVE",
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	return response, nil
}

// GetFolder retrieves a GCP folder (placeholder implementation)
func (s *GCPService) GetFolder(folderID string) (*models.FolderResponse, error) {
	// Placeholder implementation
	response := &models.FolderResponse{
		Name:        fmt.Sprintf("folders/%s", folderID),
		DisplayName: "Sample Folder",
		State:       "ACTIVE",
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	return response, nil
}

// DeleteFolder deletes a GCP folder (placeholder implementation)
func (s *GCPService) DeleteFolder(folderID string) error {
	// Placeholder implementation
	return nil
}

// CreateBucket creates a new GCS bucket
func (s *GCPService) CreateBucket(req *models.BucketRequest) (*models.BucketResponse, error) {
	bucket := s.storageClient.Bucket(req.Name)
	
	// Set bucket attributes
	attrs := &storage.BucketAttrs{
		Location: req.Location,
		Labels:   req.Labels,
	}

	if req.StorageClass != "" {
		attrs.StorageClass = req.StorageClass
	}

	if req.Versioning {
		attrs.VersioningEnabled = true
	}

	// Create the bucket
	if err := bucket.Create(s.ctx, s.config.GCPProjectID, attrs); err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	// Get bucket attributes
	bucketAttrs, err := bucket.Attrs(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket attributes: %w", err)
	}

	response := &models.BucketResponse{
		Name:         bucketAttrs.Name,
		Location:     bucketAttrs.Location,
		StorageClass: bucketAttrs.StorageClass,
		Labels:       bucketAttrs.Labels,
		Versioning:   bucketAttrs.VersioningEnabled,
		CreateTime:   bucketAttrs.Created,
		UpdateTime:   bucketAttrs.Updated,
		SelfLink:     fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s", bucketAttrs.Name),
	}

	return response, nil
}

// GetBucket retrieves a GCS bucket
func (s *GCPService) GetBucket(bucketName string) (*models.BucketResponse, error) {
	bucket := s.storageClient.Bucket(bucketName)
	
	attrs, err := bucket.Attrs(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	response := &models.BucketResponse{
		Name:         attrs.Name,
		Location:     attrs.Location,
		StorageClass: attrs.StorageClass,
		Labels:       attrs.Labels,
		Versioning:   attrs.VersioningEnabled,
		CreateTime:   attrs.Created,
		UpdateTime:   attrs.Updated,
		SelfLink:     fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s", attrs.Name),
	}

	return response, nil
}

// DeleteBucket deletes a GCS bucket
func (s *GCPService) DeleteBucket(bucketName string) error {
	bucket := s.storageClient.Bucket(bucketName)
	
	if err := bucket.Delete(s.ctx); err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}

	return nil
}

// Close closes all GCP clients
func (s *GCPService) Close() error {
	if s.storageClient != nil {
		if err := s.storageClient.Close(); err != nil {
			return err
		}
	}
	return nil
}