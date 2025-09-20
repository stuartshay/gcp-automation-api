package models

import "time"

// BucketRequest represents a request to create a GCS bucket
//
//	@Example Basic {
//	  "name": "my-simple-storage-bucket",
//	  "location": "us-central1",
//	  "storage_class": "STANDARD"
//	}
//
//	@Example Advanced {
//	  "name": "my-enterprise-bucket-2024",
//	  "location": "us-central1",
//	  "storage_class": "STANDARD",
//	  "labels": {
//	    "environment": "production",
//	    "team": "platform",
//	    "cost-center": "engineering"
//	  },
//	  "versioning": true,
//	  "kms_key_name": "projects/velvety-byway-327718/locations/us-central1/keyRings/bucket-encryption/cryptoKeys/bucket-key",
//	  "retention_policy": {
//	    "retention_period_seconds": 7776000,
//	    "is_locked": false
//	  },
//	  "uniform_bucket_level_access": true,
//	  "public_access_prevention": "enforced"
//	}
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

// BucketUpdateRequest represents a request to update a GCS bucket
type BucketUpdateRequest struct {
	Labels                   map[string]string `json:"labels,omitempty"`
	Versioning               *bool             `json:"versioning,omitempty"`
	KMSKeyName               string            `json:"kms_key_name,omitempty"`
	RetentionPolicy          *RetentionPolicy  `json:"retention_policy,omitempty"`
	UniformBucketLevelAccess *bool             `json:"uniform_bucket_level_access,omitempty"`
	PublicAccessPrevention   string            `json:"public_access_prevention,omitempty"`
}

// ObjectResponse represents a GCS object response
type ObjectResponse struct {
	Name         string            `json:"name"`
	Bucket       string            `json:"bucket"`
	Size         int64             `json:"size"`
	ContentType  string            `json:"content_type"`
	MD5Hash      string            `json:"md5_hash"`
	CRC32C       string            `json:"crc32c"`
	CreateTime   time.Time         `json:"create_time"`
	UpdateTime   time.Time         `json:"update_time"`
	Generation   int64             `json:"generation"`
	StorageClass string            `json:"storage_class"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	SelfLink     string            `json:"self_link"`
}

// LifecyclePolicy represents a bucket lifecycle policy
type LifecyclePolicy struct {
	Rules []LifecycleRule `json:"rules"`
}

// LifecycleRule represents a single lifecycle rule
type LifecycleRule struct {
	Action    LifecycleAction    `json:"action"`
	Condition LifecycleCondition `json:"condition"`
}

// LifecycleAction represents a lifecycle action
type LifecycleAction struct {
	Type         string `json:"type"` // Delete, SetStorageClass
	StorageClass string `json:"storage_class,omitempty"`
}

// LifecycleCondition represents a lifecycle condition
type LifecycleCondition struct {
	Age                   int      `json:"age,omitempty"`
	CreatedBefore         string   `json:"created_before,omitempty"`
	IsLive                *bool    `json:"is_live,omitempty"`
	MatchesStorageClass   []string `json:"matches_storage_class,omitempty"`
	NumberOfNewerVersions int      `json:"number_of_newer_versions,omitempty"`
	MatchesPrefix         []string `json:"matches_prefix,omitempty"`
	MatchesSuffix         []string `json:"matches_suffix,omitempty"`
}

// IAMPolicy represents an IAM policy
type IAMPolicy struct {
	Bindings []IAMBinding `json:"bindings"`
	Etag     string       `json:"etag"`
	Version  int          `json:"version"`
}

// IAMBinding represents an IAM binding
type IAMBinding struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}
