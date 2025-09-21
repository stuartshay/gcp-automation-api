package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// CreateBucket handles bucket creation requests
// @Summary Create a new Cloud Storage bucket
// @Description Create a new Google Cloud Storage bucket with the specified parameters including advanced options for KMS encryption, retention policies, and access controls
// @Description
// @Description ## Example Usage:
// @Description ### Basic Example (models.BasicBucketRequest):
// @Description Simple bucket creation with minimal required fields
// @Description ### Advanced Example (models.AdvancedBucketRequest):
// @Description Enterprise bucket with security features, labels, and compliance settings
// @Tags Buckets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bucket body models.BucketRequest true "Bucket creation request"
// @Success 201 {object} models.SuccessResponse{data=models.BucketResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /buckets [post]
func (h *Handler) CreateBucket(c *gin.Context) {
	var req models.BucketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate the request
	if err := h.validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
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
func (h *Handler) GetBucket(c *gin.Context) {
	bucketName := c.Param("name")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing bucket name",
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
func (h *Handler) DeleteBucket(c *gin.Context) {
	bucketName := c.Param("name")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing bucket name",
			Message: "Bucket name is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.gcpService.DeleteBucket(bucketName); err != nil {
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
