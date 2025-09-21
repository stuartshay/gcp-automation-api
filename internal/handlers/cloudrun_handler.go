package handlers

import (
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"

	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/internal/services"
	"github.com/stuartshay/gcp-automation-api/pkg/validation/gcp"
)

// CloudRunHandler handles Cloud Run logging related HTTP requests
type CloudRunHandler struct {
	cloudRunService services.CloudRunServiceInterface
}

// NewCloudRunHandler creates a new Cloud Run handler
func NewCloudRunHandler(cloudRunService services.CloudRunServiceInterface) *CloudRunHandler {
	return &CloudRunHandler{
		cloudRunService: cloudRunService,
	}
}

// ConfigureLogging configures logging for a Cloud Run service and logs request/response
// @Summary Configure Cloud Run logging
// @Description Configure logging settings for a Cloud Run service including log level, retention, exports, metrics, and alerts
// @Tags cloudrun
// @Accept json
// @Produce json
// @Param request body models.CloudRunLoggingConfigRequest true "Logging configuration request"
// @Success 200 {object} models.CloudRunLoggingConfigResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/cloudrun/logging/configure [post]
func (h *CloudRunHandler) ConfigureLogging(c *gin.Context) {
	var req models.CloudRunLoggingConfigRequest
	logger := c.MustGet("logger").(*logging.Logger)
	start := time.Now()

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log(logging.Entry{
			Severity: logging.Error,
			Payload: map[string]interface{}{
				"error":       "invalid_request",
				"message":     err.Error(),
				"request":     c.Request.URL.Path,
				"method":      c.Request.Method,
				"duration_ms": time.Since(start).Milliseconds(),
			},
		})
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate required fields
	if req.ServiceName == "" {
		logger.Log(logging.Entry{
			Severity: logging.Error,
			Payload: map[string]interface{}{
				"error":       "validation_failed",
				"message":     "Service name is required",
				"request":     c.Request.URL.Path,
				"method":      c.Request.Method,
				"duration_ms": time.Since(start).Milliseconds(),
			},
		})
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Service name is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if req.Region == "" {
		logger.Log(logging.Entry{
			Severity: logging.Error,
			Payload: map[string]interface{}{
				"error":       "validation_failed",
				"message":     "Region is required",
				"request":     c.Request.URL.Path,
				"method":      c.Request.Method,
				"duration_ms": time.Since(start).Milliseconds(),
			},
		})
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Region is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.cloudRunService.ConfigureLogging(c.Request.Context(), &req)
	if err != nil {
		logger.Log(logging.Entry{
			Severity: logging.Error,
			Payload: map[string]interface{}{
				"error":       "configuration_failed",
				"message":     err.Error(),
				"request":     c.Request.URL.Path,
				"method":      c.Request.Method,
				"duration_ms": time.Since(start).Milliseconds(),
			},
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "configuration_failed",
			Message: "Failed to configure logging: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	logger.Log(logging.Entry{
		Severity: logging.Info,
		Payload: map[string]interface{}{
			"message":      "ConfigureLogging success",
			"service_name": req.ServiceName,
			"region":       req.Region,
			"request":      c.Request.URL.Path,
			"method":       c.Request.Method,
			"duration_ms":  time.Since(start).Milliseconds(),
		},
	})
	c.JSON(http.StatusOK, response)
}

// GetLoggingConfig retrieves the current logging configuration for a Cloud Run service
// @Summary Get Cloud Run logging configuration
// @Description Retrieve the current logging configuration for a Cloud Run service
// @Tags cloudrun
// @Produce json
// @Param serviceName path string true "Cloud Run service name"
// @Param region path string true "Cloud Run service region"
// @Success 200 {object} models.CloudRunLoggingConfigResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/cloudrun/logging/{serviceName}/{region} [get]
func (h *CloudRunHandler) GetLoggingConfig(c *gin.Context) {
	serviceName := c.Param("serviceName")
	region := c.Param("region")

	if serviceName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Service name is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if region == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Region is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate service name and region
	if err := gcp.ValidateCloudRunServiceName(serviceName); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Invalid service name: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := gcp.ValidateCloudRunRegion(region); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Invalid region: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.cloudRunService.GetLoggingConfig(c.Request.Context(), serviceName, region)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "service_not_found",
			Message: "Failed to get logging configuration: " + err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateLoggingConfig updates the logging configuration for a Cloud Run service
// @Summary Update Cloud Run logging configuration
// @Description Update the logging configuration for a Cloud Run service
// @Tags cloudrun
// @Accept json
// @Produce json
// @Param serviceName path string true "Cloud Run service name"
// @Param region path string true "Cloud Run service region"
// @Param request body models.CloudRunLoggingConfigUpdateRequest true "Logging configuration update request"
// @Success 200 {object} models.CloudRunLoggingConfigResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/cloudrun/logging/{serviceName}/{region} [patch]
func (h *CloudRunHandler) UpdateLoggingConfig(c *gin.Context) {
	serviceName := c.Param("serviceName")
	region := c.Param("region")

	var req models.CloudRunLoggingConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if serviceName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Service name is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if region == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Region is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate service name and region
	if err := gcp.ValidateCloudRunServiceName(serviceName); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Invalid service name: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := gcp.ValidateCloudRunRegion(region); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Invalid region: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.cloudRunService.UpdateLoggingConfig(c.Request.Context(), serviceName, region, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "update_failed",
			Message: "Failed to update logging configuration: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetLogs retrieves logs for a Cloud Run service
// @Summary Get Cloud Run logs
// @Description Retrieve logs for a Cloud Run service with optional filtering and pagination
// @Tags cloudrun
// @Produce json
// @Param serviceName path string true "Cloud Run service name"
// @Param region path string true "Cloud Run service region"
// @Param startTime query string false "Start time for logs (RFC3339 format)"
// @Param endTime query string false "End time for logs (RFC3339 format)"
// @Param filter query string false "Additional log filter"
// @Param pageSize query int false "Number of logs to return (default: 100, max: 1000)"
// @Success 200 {object} models.CloudRunLogsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/cloudrun/logs/{serviceName}/{region} [get]
func (h *CloudRunHandler) GetLogs(c *gin.Context) {
	serviceName := c.Param("serviceName")
	region := c.Param("region")

	if serviceName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Service name is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if region == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Region is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate service name and region
	if err := gcp.ValidateCloudRunServiceName(serviceName); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Invalid service name: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := gcp.ValidateCloudRunRegion(region); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Invalid region: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Build request from query parameters
	req := &models.CloudRunLogsRequest{
		ServiceName: serviceName,
		Region:      region,
		Filter:      c.Query("filter"),
		PageSize:    100, // Default
	}

	// Parse start time
	if startTimeStr := c.Query("startTime"); startTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "validation_failed",
				Message: "Invalid startTime format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)",
				Code:    http.StatusBadRequest,
			})
			return
		}
		req.StartTime = startTime
	}

	// Parse end time
	if endTimeStr := c.Query("endTime"); endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "validation_failed",
				Message: "Invalid endTime format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)",
				Code:    http.StatusBadRequest,
			})
			return
		}
		req.EndTime = endTime
	}

	// Parse page size
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "validation_failed",
				Message: "Invalid pageSize format. Must be a number",
				Code:    http.StatusBadRequest,
			})
			return
		}
		if pageSize <= 0 || pageSize > 1000 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "validation_failed",
				Message: "pageSize must be between 1 and 1000",
				Code:    http.StatusBadRequest,
			})
			return
		}
		req.PageSize = pageSize
	}

	response, err := h.cloudRunService.GetLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "logs_retrieval_failed",
			Message: "Failed to retrieve logs: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetServiceInfo retrieves information about a Cloud Run service
// @Summary Get Cloud Run service information
// @Description Retrieve detailed information about a Cloud Run service
// @Tags cloudrun
// @Produce json
// @Param serviceName path string true "Cloud Run service name"
// @Param region path string true "Cloud Run service region"
// @Success 200 {object} models.CloudRunServiceInfo
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/cloudrun/service/{serviceName}/{region} [get]
func (h *CloudRunHandler) GetServiceInfo(c *gin.Context) {
	serviceName := c.Param("serviceName")
	region := c.Param("region")

	if serviceName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Service name is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if region == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Region is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate service name and region
	if err := gcp.ValidateCloudRunServiceName(serviceName); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Invalid service name: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := gcp.ValidateCloudRunRegion(region); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_failed",
			Message: "Invalid region: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.cloudRunService.GetServiceInfo(c.Request.Context(), serviceName, region)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "service_not_found",
			Message: "Failed to get service information: " + err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers all Cloud Run routes with the given router group
func (h *CloudRunHandler) RegisterRoutes(rg *gin.RouterGroup) {
	cloudRun := rg.Group("/cloudrun")
	{
		// Logging endpoints
		logging := cloudRun.Group("/logging")
		{
			logging.POST("/configure", h.ConfigureLogging)
			logging.GET("/:serviceName/:region", h.GetLoggingConfig)
			logging.PATCH("/:serviceName/:region", h.UpdateLoggingConfig)
		}

		// Log retrieval endpoints
		cloudRun.GET("/logs/:serviceName/:region", h.GetLogs)

		// Service information endpoints
		cloudRun.GET("/service/:serviceName/:region", h.GetServiceInfo)
	}
}
