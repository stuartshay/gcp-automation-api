package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/internal/services"
	"github.com/stuartshay/gcp-automation-api/internal/validators"
)

// Handler contains all HTTP handlers
type Handler struct {
	gcpService *services.GCPService
	validator  *validators.CustomValidator
}

// NewHandler creates a new handler instance
func NewHandler(gcpService *services.GCPService) *Handler {
	return &Handler{
		gcpService: gcpService,
		validator:  validators.NewValidator(),
	}
}

// CreateProject handles project creation requests
// @Summary Create a new GCP project
// @Description Create a new Google Cloud Platform project with the specified parameters
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body models.ProjectRequest true "Project creation request"
// @Success 201 {object} models.SuccessResponse{data=models.ProjectResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects [post]
func (h *Handler) CreateProject(c echo.Context) error {
	var req models.ProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	// Validate the request
	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	project, err := h.gcpService.CreateProject(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create project",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Project created successfully",
		Data:    project,
	})
}

// GetProject handles project retrieval requests
// @Summary Get a GCP project by ID
// @Description Retrieve details of a Google Cloud Platform project by its project ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} models.SuccessResponse{data=models.ProjectResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /projects/{id} [get]
func (h *Handler) GetProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Project ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	project, err := h.gcpService.GetProject(projectID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Project not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Project retrieved successfully",
		Data:    project,
	})
}

// DeleteProject handles project deletion requests
// @Summary Delete a GCP project
// @Description Delete a Google Cloud Platform project by its project ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects/{id} [delete]
func (h *Handler) DeleteProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Project ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	err := h.gcpService.DeleteProject(projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete project",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Project deleted successfully",
	})
}

// CreateFolder handles folder creation requests
// @Summary Create a new GCP folder
// @Description Create a new Google Cloud Platform folder with the specified parameters
// @Tags Folders
// @Accept json
// @Produce json
// @Param folder body models.FolderRequest true "Folder creation request"
// @Success 201 {object} models.SuccessResponse{data=models.FolderResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /folders [post]
func (h *Handler) CreateFolder(c echo.Context) error {
	var req models.FolderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	// Validate the request
	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	folder, err := h.gcpService.CreateFolder(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create folder",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Folder created successfully",
		Data:    folder,
	})
}

// GetFolder handles folder retrieval requests
// @Summary Get a GCP folder by ID
// @Description Retrieve details of a Google Cloud Platform folder by its folder ID
// @Tags Folders
// @Accept json
// @Produce json
// @Param id path string true "Folder ID"
// @Success 200 {object} models.SuccessResponse{data=models.FolderResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /folders/{id} [get]
func (h *Handler) GetFolder(c echo.Context) error {
	folderID := c.Param("id")
	if folderID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Folder ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	folder, err := h.gcpService.GetFolder(folderID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Folder not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Folder retrieved successfully",
		Data:    folder,
	})
}

// DeleteFolder handles folder deletion requests
// @Summary Delete a GCP folder
// @Description Delete a Google Cloud Platform folder by its folder ID
// @Tags Folders
// @Accept json
// @Produce json
// @Param id path string true "Folder ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /folders/{id} [delete]
func (h *Handler) DeleteFolder(c echo.Context) error {
	folderID := c.Param("id")
	if folderID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Folder ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	err := h.gcpService.DeleteFolder(folderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete folder",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Folder deleted successfully",
	})
}

// CreateBucket handles bucket creation requests
// @Summary Create a new Cloud Storage bucket
// @Description Create a new Google Cloud Storage bucket with the specified parameters
// @Tags Buckets
// @Accept json
// @Produce json
// @Param bucket body models.BucketRequest true "Bucket creation request"
// @Success 201 {object} models.SuccessResponse{data=models.BucketResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /buckets [post]
func (h *Handler) CreateBucket(c echo.Context) error {
	var req models.BucketRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	// Validate the request
	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	bucket, err := h.gcpService.CreateBucket(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create bucket",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Bucket created successfully",
		Data:    bucket,
	})
}

// GetBucket handles bucket retrieval requests
// @Summary Get a Cloud Storage bucket by name
// @Description Retrieve details of a Google Cloud Storage bucket by its bucket name
// @Tags Buckets
// @Accept json
// @Produce json
// @Param name path string true "Bucket name"
// @Success 200 {object} models.SuccessResponse{data=models.BucketResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /buckets/{name} [get]
func (h *Handler) GetBucket(c echo.Context) error {
	bucketName := c.Param("name")
	if bucketName == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Bucket name is required",
			Code:    http.StatusBadRequest,
		})
	}

	bucket, err := h.gcpService.GetBucket(bucketName)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Bucket not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Bucket retrieved successfully",
		Data:    bucket,
	})
}

// DeleteBucket handles bucket deletion requests
// @Summary Delete a Cloud Storage bucket
// @Description Delete a Google Cloud Storage bucket by its bucket name
// @Tags Buckets
// @Accept json
// @Produce json
// @Param name path string true "Bucket name"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /buckets/{name} [delete]
func (h *Handler) DeleteBucket(c echo.Context) error {
	bucketName := c.Param("name")
	if bucketName == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Bucket name is required",
			Code:    http.StatusBadRequest,
		})
	}

	err := h.gcpService.DeleteBucket(bucketName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete bucket",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Bucket deleted successfully",
	})
}
