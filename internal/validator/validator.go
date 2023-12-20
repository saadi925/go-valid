package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/saadi925/go-valid/internal/rules"
)

// ValidationErrors is a custom type for validation errors.
type ValidationErrors []error

// Error implements the error interface for ValidationErrors.
func (ve ValidationErrors) Error() string {
	var errorStrings []string
	for _, err := range ve {
		errorStrings = append(errorStrings, err.Error())
	}
	return strings.Join(errorStrings, "\n")
}

// Validator is the main validation struct.
type Validator struct {
	mu                  sync.Mutex
	customPasswordRules *rules.PasswordRules
}

// NewValidator creates a new instance of the Validator with optional custom password rules.
func NewValidator() *Validator {
	return &Validator{}
}

// SetCustomPasswordRules sets custom password rules for the validator.
func (v *Validator) SetCustomPasswordRules(customPasswordRules *rules.PasswordRules) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.customPasswordRules = customPasswordRules
}

// Validate performs validation on the provided struct.
func (v *Validator) Validate(data interface{}) error {
	var validationErrors ValidationErrors
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Struct {
		err := v.structValidate(val, &validationErrors)
		if err != nil {
			return err
		}
	} else {
		validationErrors = append(validationErrors, fmt.Errorf("unsupported type for validation: %v", val.Kind()))
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func (v *Validator) structValidate(val reflect.Value, validationErrors *ValidationErrors) error {
	var wg sync.WaitGroup

	for i := 0; i < val.NumField(); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			field := val.Field(i)
			fieldType := val.Type().Field(i)
			tag := fieldType.Tag.Get("validate")
			jsonTag := fieldType.Tag.Get("json")
			rulesList := strings.Split(tag, ",")

			for _, rule := range rulesList {
				parts := strings.SplitN(rule, "=", 2)
				name := parts[0]
				var options string
				if len(parts) == 2 {
					options = parts[1]
				}

				switch name {
				case "required":
					err := rules.Required(field)
					if err != nil {
						v.mu.Lock()
						*validationErrors = append(*validationErrors, fmt.Errorf("%s is required", fieldType.Name))
						v.mu.Unlock()
					}
				case "email":
					err := rules.Email(field)
					if err != nil {
						v.mu.Lock()
						*validationErrors = append(*validationErrors, fmt.Errorf("%s must be a valid email address", fieldType.Name))
						v.mu.Unlock()
					}
				case "min_length":
					min, err := strconv.Atoi(options)
					if err != nil {
						v.mu.Lock()
						*validationErrors = append(*validationErrors, fmt.Errorf("invalid min_length options: %s", options))
						v.mu.Unlock()
					}
					err = rules.MinLength(field, min)
					if err != nil {
						v.mu.Lock()
						*validationErrors = append(*validationErrors, fmt.Errorf("%s", err))
						v.mu.Unlock()
					}
				case "max_length":
					max, err := strconv.Atoi(options)
					if err != nil {
						v.mu.Lock()
						*validationErrors = append(*validationErrors, fmt.Errorf("invalid max_length options: %s", options))
						v.mu.Unlock()
					}
					err = rules.MaxLength(field, max)
					if err != nil {
						v.mu.Lock()
						*validationErrors = append(*validationErrors, fmt.Errorf("%s", err))
						v.mu.Unlock()
					}
				case "password":
					passwordRules, err := v.parsePasswordRules(options)
					if err != nil {
						v.mu.Lock()
						*validationErrors = append(*validationErrors, fmt.Errorf("invalid password options: %s", options))
						v.mu.Unlock()
					}
					err = rules.Password(field.Interface(), *passwordRules)
					if err != nil {
						v.mu.Lock()
						*validationErrors = append(*validationErrors, fmt.Errorf("%s", err))
						v.mu.Unlock()
					}
				}
			}

			if jsonTag != "" {
				fmt.Printf("JSON tag for %s: %s\n", fieldType.Name, jsonTag)
			}

			if field.Kind() == reflect.Struct {
				err := v.structValidate(field, validationErrors)
				if err != nil {
					v.mu.Lock()
					switch v := err.(type) {
					case *ValidationErrors:
						*validationErrors = append(*validationErrors, *v...)
					default:
						*validationErrors = append(*validationErrors, fmt.Errorf("%s: %s", fieldType.Name, err))
					}
					v.mu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	return nil
}

func (v *Validator) parsePasswordRules(options string) (*rules.PasswordRules, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if options == "" && v.customPasswordRules == nil {
		return rules.DefaultPasswordRules(), nil
	}

	customRules := *rules.DefaultPasswordRules()

	if v.customPasswordRules != nil {
		customRules = *v.customPasswordRules
	}

	if options != "" {
		parts := strings.Split(options, ",")
		for _, part := range parts {
			subparts := strings.SplitN(part, "=", 2)
			if len(subparts) != 2 {
				return rules.DefaultPasswordRules(), errors.New("invalid password rule format")
			}
			switch subparts[0] {
			case "min_length":
				customRules.MinLength, _ = strconv.Atoi(subparts[1])
			case "require_digits":
				customRules.RequireDigits, _ = strconv.ParseBool(subparts[1])
			case "require_uppercase":
				customRules.RequireUppercase, _ = strconv.ParseBool(subparts[1])
			case "require_lowercase":
				customRules.RequireLowercase, _ = strconv.ParseBool(subparts[1])
			case "require_special_chars":
				customRules.RequireSpecialChars, _ = strconv.ParseBool(subparts[1])
			}
		}
	}

	return &customRules, nil
}
