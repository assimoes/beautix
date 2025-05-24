package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/assimoes/beautix/pkg/errors"
)

// Validator is a global validator instance
var Validator *validator.Validate

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// init initializes the validator
func init() {
	Validator = validator.New()
	
	// Register custom tag name function to use JSON tags
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	// Register custom validations
	registerCustomValidations()
}

// registerCustomValidations registers custom validation rules
func registerCustomValidations() {
	// Register UUID validation
	Validator.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		return IsValidUUID(fl.Field().String())
	})
	
	// Register role validation
	Validator.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		role := fl.Field().String()
		validRoles := []string{"admin", "owner", "staff", "user"}
		for _, validRole := range validRoles {
			if role == validRole {
				return true
			}
		}
		return false
	})
	
	// Register currency validation
	Validator.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		currency := fl.Field().String()
		validCurrencies := []string{"EUR", "USD", "GBP", "BRL"}
		for _, validCurrency := range validCurrencies {
			if currency == validCurrency {
				return true
			}
		}
		return false
	})
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(s interface{}) error {
	err := Validator.Struct(s)
	if err == nil {
		return nil
	}
	
	var validationErrors []ValidationError
	
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   err.Field(),
			Tag:     err.Tag(),
			Value:   err.Param(),
			Message: getValidationMessage(err),
		})
	}
	
	return errors.NewValidationError("Validation failed", validationErrors)
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email address"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters long"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters long"
	case "uuid":
		return err.Field() + " must be a valid UUID"
	case "role":
		return err.Field() + " must be one of: admin, owner, staff, user"
	case "currency":
		return err.Field() + " must be a valid currency code (EUR, USD, GBP, BRL)"
	case "len":
		return err.Field() + " must be exactly " + err.Param() + " characters long"
	default:
		return err.Field() + " is invalid"
	}
}