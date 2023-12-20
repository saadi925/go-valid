package rules

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// PasswordRules defines the rules for password validation.
type PasswordRules struct {
	MinLength           int
	RequireDigits       bool
	RequireUppercase    bool
	RequireLowercase    bool
	RequireSpecialChars bool
}

// DefaultPasswordRules returns the default password rules.
func DefaultPasswordRules() *PasswordRules {
	return &PasswordRules{
		MinLength:           8,
		RequireDigits:       true,
		RequireSpecialChars: true,
		RequireUppercase:    true,
		RequireLowercase:    true,
	}
}

// Email validates if the field is a valid email address.
func Email(field reflect.Value) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	email := field.String()

	match, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return fmt.Errorf("error validating email: %v", err)
	}

	if !match {
		return fmt.Errorf("field must be a valid email address")
	}

	return nil
}

// Password validates if the field meets the password rules.
func Password(field interface{}, rules PasswordRules) error {
	password, ok := field.(string)
	if !ok {
		return fmt.Errorf("invalid field type for password validation")
	}

	if len(password) < rules.MinLength {
		return fmt.Errorf("password must have a minimum length of %d characters", rules.MinLength)
	}

	if rules.RequireDigits && !containsDigits(password) {
		return fmt.Errorf("password must contain at least one digit")
	}

	if rules.RequireUppercase && !containsUppercase(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if rules.RequireLowercase && !containsLowercase(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if rules.RequireSpecialChars && !containsSpecialChars(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// Helper functions for password validation
// Add these functions to the rules package as well
func containsDigits(s string) bool {
	for _, char := range s {
		if char >= '0' && char <= '9' {
			return true
		}
	}
	return false
}

func containsUppercase(s string) bool {
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, char := range s {
		if char >= 'a' && char <= 'z' {
			return true
		}
	}
	return false
}

func containsSpecialChars(s string) bool {
	specialChars := "~!@#$%^&*()-_+=<>?/[]{}|"
	for _, char := range s {
		if strings.ContainsRune(specialChars, char) {
			return true
		}
	}
	return false
}
