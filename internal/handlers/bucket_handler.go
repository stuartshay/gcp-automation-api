package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// CreateBucket handles bucket creation requests
// @Summary Create a new Cloud Storage bucket
// @Description Create a new Google Cloud Storage bucket with the specified parameters including advanced options for KMS encryption, retention policies, and access controls
// @Tags Buckets
// @Accept json
// @Produce json
// @Security BearerAuth
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
// @Security BearerAuth
// @Param name path string true "Bucket name"
// @Success 200 {object} models.SuccessResponse{data=models.BucketResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /buckets/{name} [get]
func (h *Handler) GetBucket(c echo.Context) error {
	bucketName := c.Param("name")
	if bucketName == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing bucket name",
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
// @Security BearerAuth
// @Param name path string true "Bucket name"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /buckets/{name} [delete]
func (h *Handler) DeleteBucket(c echo.Context) error {
	bucketName := c.Param("name")
	if bucketName == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing bucket name",
			Message: "Bucket name is required",
			Code:    http.StatusBadRequest,
		})
	}

	if err := h.gcpService.DeleteBucket(bucketName); err != nil {
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
