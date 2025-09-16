package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// CreateProject handles project creation requests
// @Summary Create a new GCP project
// @Description Create a new Google Cloud Platform project with the specified parameters
// @Description
// @Description ## Example Usage:
// @Description ### Basic Example (models.BasicProjectRequest):
// @Description Simple project creation with just ID and display name
// @Description ### Advanced Example (models.AdvancedProjectRequest):
// @Description Enterprise project with organization hierarchy, labels, and governance
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
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
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} models.SuccessResponse{data=models.ProjectResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects/{id} [get]
func (h *Handler) GetProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing project ID",
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
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects/{id} [delete]
func (h *Handler) DeleteProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing project ID",
			Message: "Project ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	if err := h.gcpService.DeleteProject(projectID); err != nil {
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
