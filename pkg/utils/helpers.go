package utils

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the given int
func IntPtr(i int) *int {
	return &i
}

// BoolPtr returns a pointer to the given bool
func BoolPtr(b bool) *bool {
	return &b
}

// TimePtr returns a pointer to the given time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// SafeString returns the value of a string pointer or empty string if nil
func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// SafeInt returns the value of an int pointer or 0 if nil
func SafeInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// SafeBool returns the value of a bool pointer or false if nil
func SafeBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// SafeTime returns the value of a time pointer or zero time if nil
func SafeTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}

// IsValidEmail checks if a string is a valid email address
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// NormalizeEmail normalizes an email address (lowercase, trim)
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizeName normalizes a name (trim, title case)
func NormalizeName(name string) string {
	return strings.TrimSpace(strings.Title(strings.ToLower(name)))
}

// IsValidPhoneNumber checks if a string is a valid phone number (basic validation)
func IsValidPhoneNumber(phone string) bool {
	// Remove common formatting characters
	cleaned := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")
	
	// Basic validation: should have at least 10 digits, optionally starting with +
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{9,14}$`)
	return phoneRegex.MatchString(cleaned)
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}

// IsValidTimeZone checks if a string is a valid timezone
func IsValidTimeZone(tz string) bool {
	_, err := time.LoadLocation(tz)
	return err == nil
}

// TruncateString truncates a string to the specified length
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}

// ContainsString checks if a slice contains a string
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveString removes a string from a slice
func RemoveString(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// UniqueStrings returns a slice with unique strings
func UniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// GetDefaultIfEmpty returns the default value if the string is empty
func GetDefaultIfEmpty(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}