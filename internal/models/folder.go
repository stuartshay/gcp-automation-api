package models

import "time"

// BasicFolderRequest represents a basic request to create a GCP folder (standard usage)
type BasicFolderRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=1,max=100" binding:"required" example:"Development Environment"`
	ParentID    string `json:"parent_id" validate:"required,numeric" binding:"required" example:"123456789012"`
	ParentType  string `json:"parent_type" validate:"required,oneof=organization folder" binding:"required" example:"organization"`
}

// AdvancedFolderRequest represents an advanced request to create a GCP folder (nested structure)
type AdvancedFolderRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=1,max=100" binding:"required" example:"Production - North America Region"`
	ParentID    string `json:"parent_id" validate:"required,numeric" binding:"required" example:"987654321098"`
	ParentType  string `json:"parent_type" validate:"required,oneof=organization folder" binding:"required" example:"folder"`
}

// FolderRequest represents a request to create a GCP folder
// @Example Basic {
//   "display_name": "Development Environment",
//   "parent_id": "123456789012",
//   "parent_type": "organization"
// }
// @Example Advanced {
//   "display_name": "Production - North America Region",
//   "parent_id": "987654321098",
//   "parent_type": "folder"
// }

// FolderRequest represents a request to create a GCP folder
//
//	@Example Basic {
//	  "display_name": "Development Environment",
//	  "parent_id": "123456789012",
//	  "parent_type": "organization"
//	}
//
//	@Example Advanced {
//	  "display_name": "Production - North America Region",
//	  "parent_id": "987654321098",
//	  "parent_type": "folder"
//	}
type FolderRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=1,max=100" binding:"required" example:"Development Environment"`
	ParentID    string `json:"parent_id" validate:"required,numeric" binding:"required" example:"123456789012"`
	ParentType  string `json:"parent_type" validate:"required,oneof=organization folder" binding:"required" example:"organization"` // "organization" or "folder"
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
