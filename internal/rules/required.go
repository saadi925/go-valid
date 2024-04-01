package rules

import (
	"fmt"
	"reflect"
)

// Required validates if the field is required (not empty or zero value).
func Required(field reflect.Value) error {
	switch field.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		if field.Len() == 0 {
			return fmt.Errorf("field is required and must not be empty")
		}
	case reflect.Ptr, reflect.Interface:
		if field.IsNil() {
			return fmt.Errorf("field is required and must not be nil")
		}
	default:
		zeroValue := reflect.Zero(field.Type())
		if reflect.DeepEqual(field.Interface(), zeroValue.Interface()) {
			return fmt.Errorf("field is required and must not be zero value")
		}
	}
	return nil
}
