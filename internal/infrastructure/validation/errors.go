package validation

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// Error implements the error interface
func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

// NewFieldValidationError creates a new validation error for a specific field
func NewFieldValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "validation failed"
	}
	
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	
	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Add adds a validation error
func (e *ValidationErrors) Add(err ValidationError) {
	*e = append(*e, err)
}

// AddField adds a field validation error
func (e *ValidationErrors) AddField(field, message string) {
	*e = append(*e, ValidationError{
		Field:   field,
		Message: message,
	})
}