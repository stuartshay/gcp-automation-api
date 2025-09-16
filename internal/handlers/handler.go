package handlers

import (
	"github.com/stuartshay/gcp-automation-api/internal/services"
	"github.com/stuartshay/gcp-automation-api/internal/validators"
)

// Handler contains all HTTP handlers
type Handler struct {
	gcpService  services.GCPServiceInterface
	authService *services.AuthService
	validator   *validators.CustomValidator
}

// NewHandler creates a new handler instance
func NewHandler(gcpService services.GCPServiceInterface, authService *services.AuthService) *Handler {
	return &Handler{
		gcpService:  gcpService,
		authService: authService,
		validator:   validators.NewValidator(),
	}
}
