package models

import "time"

// ProjectRequest represents a request to create a GCP project
//
//	@Example Basic {
//	  "project_id": "my-simple-project-2024",
//	  "display_name": "My Simple Project"
//	}
//
//	@Example Advanced {
//	  "project_id": "enterprise-app-prod-2024",
//	  "display_name": "Enterprise Application - Production",
//	  "parent_id": "123456789012",
//	  "parent_type": "organization",
//	  "labels": {
//	    "environment": "production",
//	    "team": "backend",
//	    "cost-center": "engineering",
//	    "compliance": "sox"
//	  }
//	}
type ProjectRequest struct {
	ProjectID   string            `json:"project_id" validate:"required,project_id" binding:"required" example:"my-dev-project-2024"`
	DisplayName string            `json:"display_name" validate:"required,min=1,max=100" binding:"required" example:"My Development Project"`
	ParentID    string            `json:"parent_id,omitempty" validate:"omitempty,numeric" example:"123456789012"`
	ParentType  string            `json:"parent_type,omitempty" validate:"omitempty,oneof=organization folder" example:"organization"` // "organization" or "folder"
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
