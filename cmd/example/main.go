// main.go
package main

import (
	"fmt"

	"github.com/saadi925/go-val/internal/rules"
	"github.com/saadi925/go-val/internal/validator"
)

// User is a struct representing user data with validation tags.
type User struct {
	Username string `validate:"required,min_length=3,max_length=10"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,password"`
}

func main() {
	// Valid user data
	validUser := User{
		Username: "codestack",
		Email:    "user@example.com",
		Password: "Password123@#12",
	}

	// Create a new instance of the validator
	myValidator := validator.NewValidator()

	// Validate the valid user
	err := myValidator.Validate(validUser)
	if err != nil {
		fmt.Println("Error occurred, this user is invalid:", err)
	} else {
		fmt.Println("User is valid")
	}

	// Invalid user data
	invalidUser := User{
		Username: "short",
		Email:    "invalid-email",
		Password: "weak",
	}

	// Validate the invalid user
	err = myValidator.Validate(invalidUser)
	if err != nil {
		fmt.Println("Invalid user:", err) // Password is too short, Email is not valid
	}

	// Custom password rules
	customPasswordRules := &rules.PasswordRules{
		RequireUppercase:    false,
		RequireLowercase:    false,
		RequireDigits:       false,
		RequireSpecialChars: false,
	}

	// Set custom password rules
	myValidator.SetCustomPasswordRules(customPasswordRules)

	// Validate the invalid user with custom password rules
	err = myValidator.Validate(invalidUser)
	if err != nil {
		fmt.Println("Error:", err) // User is now considered valid with custom password rules
	}

	// Additional example: Validate a struct with nested structs
	type Address struct {
		City  string `validate:"required"`
		State string `validate:"required"`
	}

	type Profile struct {
		Name    string  `validate:"required"`
		Address Address `validate:"required"`
	}

	profileWithInvalidAddress := Profile{
		Name: "John Doe",
		Address: Address{
			City:  "",
			State: "New York",
		},
	}

	// Validate the profile with nested structs
	err = myValidator.Validate(profileWithInvalidAddress)
	if err != nil {
		fmt.Println("Invalid profile with nested structs:", err) // City is required
	}
}
