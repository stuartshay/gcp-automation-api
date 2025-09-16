package docs

import (
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// Example structs for Swagger UI dropdown - these will appear as separate schemas
// in the Swagger UI allowing users to choose between Basic and Advanced examples

// SwaggerBasicExamples contains all basic example models for API documentation
type SwaggerBasicExamples struct {
	BucketExample  models.BasicBucketRequest  `json:"bucket_example"`
	ProjectExample models.BasicProjectRequest `json:"project_example"`
	FolderExample  models.BasicFolderRequest  `json:"folder_example"`
}

// SwaggerAdvancedExamples contains all advanced example models for API documentation
type SwaggerAdvancedExamples struct {
	BucketExample  models.AdvancedBucketRequest  `json:"bucket_example"`
	ProjectExample models.AdvancedProjectRequest `json:"project_example"`
	FolderExample  models.AdvancedFolderRequest  `json:"folder_example"`
}
