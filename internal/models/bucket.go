package models

import "time"

// BucketRequest represents a request to create a GCS bucket
type BucketRequest struct {
	Name         string            `json:"name" validate:"required,bucket_name" binding:"required" example:"my-project-storage-bucket"`
	Location     string            `json:"location" validate:"required,gcp_location" binding:"required" example:"us-central1"`
	StorageClass string            `json:"storage_class,omitempty" validate:"omitempty,oneof=STANDARD NEARLINE COLDLINE ARCHIVE" example:"STANDARD"`
	Labels       map[string]string `json:"labels,omitempty" validate:"omitempty,dive,keys,label_key,endkeys,label_value"`
	Versioning   bool              `json:"versioning,omitempty" example:"true"`

	// Phase 1 Advanced Options - Security & Compliance
	KMSKeyName               string           `json:"kms_key_name,omitempty" validate:"omitempty" example:"projects/my-project/locations/us-central1/keyRings/my-keyring/cryptoKeys/my-key"`
	RetentionPolicy          *RetentionPolicy `json:"retention_policy,omitempty" validate:"omitempty"`
	UniformBucketLevelAccess bool             `json:"uniform_bucket_level_access,omitempty" example:"true"`
	PublicAccessPrevention   string           `json:"public_access_prevention,omitempty" validate:"omitempty,oneof=inherited enforced unspecified" example:"enforced"`
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

	// Phase 1 Advanced Options - Security & Compliance
	KMSKeyName               string           `json:"kms_key_name,omitempty"`
	RetentionPolicy          *RetentionPolicy `json:"retention_policy,omitempty"`
	UniformBucketLevelAccess bool             `json:"uniform_bucket_level_access,omitempty"`
	PublicAccessPrevention   string           `json:"public_access_prevention,omitempty"`
}
