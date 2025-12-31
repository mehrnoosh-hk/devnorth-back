package dto

import "fmt"

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
