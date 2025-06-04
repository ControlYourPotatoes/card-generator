package validation

import "fmt"

// ValidationError represents a card validation error
type ValidationError struct {
	Type    ErrorType
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.Field)
}

// ErrorType represents the type of validation error
type ErrorType string

const (
	ErrorTypeRequired ErrorType = "required"
	ErrorTypeInvalid  ErrorType = "invalid"
	ErrorTypeRange    ErrorType = "range"
	ErrorTypeFormat   ErrorType = "format"
)

// NewValidationError creates a new ValidationError
func NewValidationError(errType ErrorType, message, field string) *ValidationError {
	return &ValidationError{
		Type:    errType,
		Message: message,
		Field:   field,
	}
}
