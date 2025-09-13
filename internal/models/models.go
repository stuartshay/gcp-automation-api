package models

import "time"

// ProjectRequest represents a request to create a GCP project
type ProjectRequest struct {
	ProjectID   string            `json:"project_id" validate:"required,project_id" binding:"required"`
	DisplayName string            `json:"display_name" validate:"required,min=1,max=100" binding:"required"`
	ParentID    string            `json:"parent_id,omitempty" validate:"omitempty,numeric"`
	ParentType  string            `json:"parent_type,omitempty" validate:"omitempty,oneof=organization folder"` // "organization" or "folder"
	Labels      map[string]string `json:"labels,omitempty" validate:"omitempty,dive,keys,label_key,endkeys,label_value"`
}

// ProjectResponse represents a GCP project response
type ProjectResponse struct {
	ProjectID     string            `json:"project_id"`
	DisplayName   string            `json:"display_name"`
	ParentID      string            `json:"parent_id,omitempty"`
	ParentType    string            `json:"parent_type,omitempty"`
	State         string            `json:"state"`
	Labels        map[string]string `json:"labels,omitempty"`
	CreateTime    time.Time         `json:"create_time"`
	UpdateTime    time.Time         `json:"update_time"`
	ProjectNumber int64             `json:"project_number"`
}

// FolderRequest represents a request to create a GCP folder
type FolderRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=1,max=100" binding:"required"`
	ParentID    string `json:"parent_id" validate:"required,numeric" binding:"required"`
	ParentType  string `json:"parent_type" validate:"required,oneof=organization folder" binding:"required"` // "organization" or "folder"
}

// FolderResponse represents a GCP folder response
type FolderResponse struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	ParentID    string    `json:"parent_id"`
	ParentType  string    `json:"parent_type"`
	State       string    `json:"state"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

// BucketRequest represents a request to create a GCS bucket
type BucketRequest struct {
	Name         string            `json:"name" validate:"required,bucket_name" binding:"required"`
	Location     string            `json:"location" validate:"required,gcp_location" binding:"required"`
	StorageClass string            `json:"storage_class,omitempty" validate:"omitempty,oneof=STANDARD NEARLINE COLDLINE ARCHIVE"`
	Labels       map[string]string `json:"labels,omitempty" validate:"omitempty,dive,keys,label_key,endkeys,label_value"`
	Versioning   bool              `json:"versioning,omitempty"`
}

// BucketResponse represents a GCS bucket response
type BucketResponse struct {
	Name         string            `json:"name"`
	Location     string            `json:"location"`
	StorageClass string            `json:"storage_class"`
	Labels       map[string]string `json:"labels,omitempty"`
	Versioning   bool              `json:"versioning"`
	CreateTime   time.Time         `json:"create_time"`
	UpdateTime   time.Time         `json:"update_time"`
	SelfLink     string            `json:"self_link"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
