package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// GetBucketExamples provides example request bodies for bucket creation
// @Summary Get example request bodies for bucket creation
// @Description Provides both basic and advanced example request bodies for creating GCS buckets
// @Tags Examples
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=models.BasicBucketRequest} "Basic bucket creation example"
// @Success 200 {object} models.SuccessResponse{data=models.AdvancedBucketRequest} "Advanced bucket creation example"
// @Router /examples/buckets [get]
func (h *Handler) GetBucketExamples(c echo.Context) error {
	exampleType := c.QueryParam("type")

	switch exampleType {
	case "advanced":
		advanced := models.AdvancedBucketRequest{
			Name:         "my-enterprise-bucket-2024",
			Location:     "us-central1",
			StorageClass: "STANDARD",
			Labels: map[string]string{
				"environment": "production",
				"team":        "platform",
				"cost-center": "engineering",
			},
			Versioning: true,
			KMSKeyName: "projects/velvety-byway-327718/locations/us-central1/keyRings/bucket-encryption/cryptoKeys/bucket-key",
			RetentionPolicy: &models.RetentionPolicy{
				RetentionPeriodSeconds: 7776000, // 90 days in seconds
				IsLocked:               false,
			},
			UniformBucketLevelAccess: true,
			PublicAccessPrevention:   "enforced",
		}
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Message: "Advanced bucket creation example",
			Data:    advanced,
		})
	default:
		basic := models.BasicBucketRequest{
			Name:         "my-simple-storage-bucket",
			Location:     "us-central1",
			StorageClass: "STANDARD",
		}
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Message: "Basic bucket creation example",
			Data:    basic,
		})
	}
}

// GetProjectExamples provides example request bodies for project creation
// @Summary Get example request bodies for project creation
// @Description Provides both basic and advanced example request bodies for creating GCP projects
// @Tags Examples
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=models.BasicProjectRequest} "Basic project creation example"
// @Success 200 {object} models.SuccessResponse{data=models.AdvancedProjectRequest} "Advanced project creation example"
// @Router /examples/projects [get]
func (h *Handler) GetProjectExamples(c echo.Context) error {
	exampleType := c.QueryParam("type")

	switch exampleType {
	case "advanced":
		advanced := models.AdvancedProjectRequest{
			ProjectID:   "enterprise-app-prod-2024",
			DisplayName: "Enterprise Application - Production",
			ParentID:    "123456789012",
			ParentType:  "organization",
			Labels: map[string]string{
				"environment": "production",
				"team":        "backend",
				"cost-center": "engineering",
				"compliance":  "sox",
			},
		}
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Message: "Advanced project creation example",
			Data:    advanced,
		})
	default:
		basic := models.BasicProjectRequest{
			ProjectID:   "my-simple-project-2024",
			DisplayName: "My Simple Project",
		}
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Message: "Basic project creation example",
			Data:    basic,
		})
	}
}

// GetFolderExamples provides example request bodies for folder creation
// @Summary Get example request bodies for folder creation
// @Description Provides both basic and advanced example request bodies for creating GCP folders
// @Tags Examples
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=models.BasicFolderRequest} "Basic folder creation example"
// @Success 200 {object} models.SuccessResponse{data=models.AdvancedFolderRequest} "Advanced folder creation example"
// @Router /examples/folders [get]
func (h *Handler) GetFolderExamples(c echo.Context) error {
	exampleType := c.QueryParam("type")

	switch exampleType {
	case "advanced":
		advanced := models.AdvancedFolderRequest{
			DisplayName: "Production - North America Region",
			ParentID:    "987654321098",
			ParentType:  "folder",
		}
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Message: "Advanced folder creation example",
			Data:    advanced,
		})
	default:
		basic := models.BasicFolderRequest{
			DisplayName: "Development Environment",
			ParentID:    "123456789012",
			ParentType:  "organization",
		}
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Message: "Basic folder creation example",
			Data:    basic,
		})
	}
}
