package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// CreateFolder handles folder creation requests
// @Summary Create a new GCP folder
// @Description Create a new Google Cloud Platform folder with the specified parameters
// @Description
// @Description ## Example Usage:
// @Description ### Basic Example (models.BasicFolderRequest):
// @Description Simple folder under organization root
// @Description ### Advanced Example (models.AdvancedFolderRequest):
// @Description Nested folder structure for complex hierarchies
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
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
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Success 200 {object} models.SuccessResponse{data=models.FolderResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /folders/{id} [get]
func (h *Handler) GetFolder(c echo.Context) error {
	folderID := c.Param("id")
	if folderID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing folder ID",
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
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /folders/{id} [delete]
func (h *Handler) DeleteFolder(c echo.Context) error {
	folderID := c.Param("id")
	if folderID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing folder ID",
			Message: "Folder ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	if err := h.gcpService.DeleteFolder(folderID); err != nil {
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
