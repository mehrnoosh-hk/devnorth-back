package dto

import (
	"fmt"
	"time"
)

// JSONSerializable is a marker interface for types that can be safely JSON-serialized
// This prevents accidentally passing channels, functions, or other non-serializable types
type JSONSerializable interface {
	isJSONSerializable()
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ErrFieldRequired creates a validation error for required fields
func ErrFieldRequired(field string) error {
	return ValidationError{
		Field:   field,
		Message: "field is required",
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// Implement JSONSerializable for all response types
func (ValidationError) isJSONSerializable()  {}
func (ErrorResponse) isJSONSerializable()    {}
func (HealthResponse) isJSONSerializable()   {}
