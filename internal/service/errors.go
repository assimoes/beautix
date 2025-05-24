package service

import (
	"fmt"
)

// ServiceError represents a service-level error
type ServiceError struct {
	Message string
	Cause   error
}

// Error implements the error interface
func (e ServiceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e ServiceError) Unwrap() error {
	return e.Cause
}

// NotFoundError represents a not found error
type NotFoundError struct {
	EntityType string
	Field      string
	Value      string
}

// Error implements the error interface
func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found with %s: %s", e.EntityType, e.Field, e.Value)
}

// NewServiceError creates a new service error
func NewServiceError(message string, cause error) error {
	return ServiceError{
		Message: message,
		Cause:   cause,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(entityType, field, value string) error {
	return NotFoundError{
		EntityType: entityType,
		Field:      field,
		Value:      value,
	}
}