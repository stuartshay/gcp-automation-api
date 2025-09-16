package models

// RetentionPolicy represents bucket retention policy configuration
type RetentionPolicy struct {
	RetentionPeriodSeconds int64 `json:"retention_period_seconds" validate:"min=1,max=3155760000" example:"86400"` // 1 second to 100 years
	IsLocked               bool  `json:"is_locked" example:"false"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
