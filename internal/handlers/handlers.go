package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/internal/services"
)

// Handler contains all HTTP handlers
type Handler struct {
	gcpService *services.GCPService
}

// NewHandler creates a new handler instance
func NewHandler(gcpService *services.GCPService) *Handler {
	return &Handler{
		gcpService: gcpService,
	}
}

// CreateProject handles project creation requests
func (h *Handler) CreateProject(c *gin.Context) {
	var req models.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	project, err := h.gcpService.CreateProject(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create project",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Project created successfully",
		Data:    project,
	})
}

// GetProject handles project retrieval requests
func (h *Handler) GetProject(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Project ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	project, err := h.gcpService.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Project not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Project retrieved successfully",
		Data:    project,
	})
}

// DeleteProject handles project deletion requests
func (h *Handler) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Project ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.gcpService.DeleteProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete project",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Project deleted successfully",
	})
}

// CreateFolder handles folder creation requests
func (h *Handler) CreateFolder(c *gin.Context) {
	var req models.FolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	folder, err := h.gcpService.CreateFolder(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create folder",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Folder created successfully",
		Data:    folder,
	})
}

// GetFolder handles folder retrieval requests
func (h *Handler) GetFolder(c *gin.Context) {
	folderID := c.Param("id")
	if folderID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Folder ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	folder, err := h.gcpService.GetFolder(folderID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Folder not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Folder retrieved successfully",
		Data:    folder,
	})
}

// DeleteFolder handles folder deletion requests
func (h *Handler) DeleteFolder(c *gin.Context) {
	folderID := c.Param("id")
	if folderID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Folder ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.gcpService.DeleteFolder(folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete folder",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Folder deleted successfully",
	})
}

// CreateBucket handles bucket creation requests
func (h *Handler) CreateBucket(c *gin.Context) {
	var req models.BucketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	bucket, err := h.gcpService.CreateBucket(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create bucket",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Bucket created successfully",
		Data:    bucket,
	})
}

// GetBucket handles bucket retrieval requests
func (h *Handler) GetBucket(c *gin.Context) {
	bucketName := c.Param("name")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Bucket name is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	bucket, err := h.gcpService.GetBucket(bucketName)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Bucket not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Bucket retrieved successfully",
		Data:    bucket,
	})
}

// DeleteBucket handles bucket deletion requests
func (h *Handler) DeleteBucket(c *gin.Context) {
	bucketName := c.Param("name")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Bucket name is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.gcpService.DeleteBucket(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete bucket",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Bucket deleted successfully",
	})
}