package mocks

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// MockGCPService is a mock implementation of GCPServiceInterface
type MockGCPService struct {
	mock.Mock
}

// CreateProject mocks the CreateProject method
func (m *MockGCPService) CreateProject(req *models.ProjectRequest) (*models.ProjectResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProjectResponse), args.Error(1)
}

// GetProject mocks the GetProject method
func (m *MockGCPService) GetProject(projectID string) (*models.ProjectResponse, error) {
	args := m.Called(projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProjectResponse), args.Error(1)
}

// DeleteProject mocks the DeleteProject method
func (m *MockGCPService) DeleteProject(projectID string) error {
	args := m.Called(projectID)
	return args.Error(0)
}

// CreateFolder mocks the CreateFolder method
func (m *MockGCPService) CreateFolder(req *models.FolderRequest) (*models.FolderResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FolderResponse), args.Error(1)
}

// GetFolder mocks the GetFolder method
func (m *MockGCPService) GetFolder(folderID string) (*models.FolderResponse, error) {
	args := m.Called(folderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FolderResponse), args.Error(1)
}

// DeleteFolder mocks the DeleteFolder method
func (m *MockGCPService) DeleteFolder(folderID string) error {
	args := m.Called(folderID)
	return args.Error(0)
}

// CreateBucket mocks the CreateBucket method
func (m *MockGCPService) CreateBucket(req *models.BucketRequest) (*models.BucketResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BucketResponse), args.Error(1)
}

// GetBucket mocks the GetBucket method
func (m *MockGCPService) GetBucket(bucketName string) (*models.BucketResponse, error) {
	args := m.Called(bucketName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BucketResponse), args.Error(1)
}

// DeleteBucket mocks the DeleteBucket method
func (m *MockGCPService) DeleteBucket(bucketName string) error {
	args := m.Called(bucketName)
	return args.Error(0)
}

// Close mocks the Close method
func (m *MockGCPService) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Helper methods for creating test data

// NewMockProjectResponse creates a mock project response
func NewMockProjectResponse(req *models.ProjectRequest) *models.ProjectResponse {
	return &models.ProjectResponse{
		ProjectID:     req.ProjectID,
		DisplayName:   req.DisplayName,
		ParentID:      req.ParentID,
		ParentType:    req.ParentType,
		State:         "ACTIVE",
		Labels:        req.Labels,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		ProjectNumber: 123456789,
	}
}

// NewMockFolderResponse creates a mock folder response
func NewMockFolderResponse(req *models.FolderRequest) *models.FolderResponse {
	return &models.FolderResponse{
		Name:        fmt.Sprintf("folders/%s", "mock-folder-id"),
		DisplayName: req.DisplayName,
		ParentID:    req.ParentID,
		ParentType:  req.ParentType,
		State:       "ACTIVE",
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
}

// NewMockBucketResponse creates a mock bucket response
func NewMockBucketResponse(req *models.BucketRequest) *models.BucketResponse {
	storageClass := req.StorageClass
	if storageClass == "" {
		storageClass = "STANDARD"
	}

	return &models.BucketResponse{
		Name:         req.Name,
		Location:     req.Location,
		StorageClass: storageClass,
		Labels:       req.Labels,
		Versioning:   req.Versioning,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
		SelfLink:     fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s", req.Name),
	}
}
