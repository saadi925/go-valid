// internal/validator/rules/rules.go
package rules

import (
	"fmt"
	"reflect"
)

// MinLength validates if the field meets the minimum length requirement.
func MinLength(field reflect.Value, min int) error {
	switch field.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		if field.Len() < min {
			return fmt.Errorf("field must have a minimum length of %d characters", min)
		}
	default:
		return fmt.Errorf("unsupported type for MinLength validation: %v", field.Kind())
	}
	return nil
}

// MaxLength validates if the field meets the maximum length requirement.
func MaxLength(field reflect.Value, max int) error {
	switch field.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		if field.Len() > max {
			return fmt.Errorf("field must have a maximum length of %d characters", max)
		}
	default:
		return fmt.Errorf("unsupported type for MaxLength validation: %v", field.Kind())
	}
	return nil
}
