// internal/core/card/validation.go
package card

import (
	"fmt"

	"github.com/yourusername/cardgame/internal/core/card/base"
)

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
    ErrorTypeRequired   ErrorType = "required"
    ErrorTypeInvalid    ErrorType = "invalid"
    ErrorTypeRange      ErrorType = "range"
    ErrorTypeFormat     ErrorType = "format"
)

// ValidateBase performs basic card validation
func (b BaseCard) ValidateBase() *ValidationError {
    if b.Name == "" {
        return &ValidationError{
            Type:    ErrorTypeRequired,
            Message: "name cannot be empty",
            Field:   "name",
        }
    }
    if len(b.Name) > 40 {
        return &ValidationError{
            Type:    ErrorTypeRange,
            Message: "name exceeds maximum length of 40 characters",
            Field:   "name",
        }
    }
    if b.Cost < -1 { // -1 allowed for X costs
        return &ValidationError{
            Type:    ErrorTypeRange,
            Message: "cost cannot be negative (except -1 for X costs)",
            Field:   "cost",
        }
    }
    if b.Effect == "" {
        return &ValidationError{
            Type:    ErrorTypeRequired,
            Message: "effect cannot be empty",
            Field:   "effect",
        }
    }
    return nil
}