package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// RegisterCustomValidations registers custom validation rules
func RegisterCustomValidations() {
	// Example of a custom validation: "is-cool"
	validate.RegisterValidation("is-cool", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "cool"
	})
}

// FormatValidationErrors formats validation errors into a map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		tag := e.Tag()
		param := e.Param()

		// Handle built-in and custom tags
		switch tag {
		case "required":
			errors[field] = fmt.Sprintf("%s is required", field)
		case "email":
			errors[field] = fmt.Sprintf("%s must be a valid email address", field)
		case "min":
			errors[field] = fmt.Sprintf("%s must be at least %s characters long", field, param)
		case "max":
			errors[field] = fmt.Sprintf("%s must be at most %s characters long", field, param)
		case "gte":
			errors[field] = fmt.Sprintf("%s must be greater than or equal to %s", field, param)
		case "lte":
			errors[field] = fmt.Sprintf("%s must be less than or equal to %s", field, param)
		case "is-cool":
			errors[field] = fmt.Sprintf("%s must be 'cool'", field)
		default:
			errors[field] = fmt.Sprintf("%s is invalid", field)
		}
	}

	return errors
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(input interface{}) (bool, map[string]string) {
	err := validate.Struct(input)
	if err != nil {
		return false, FormatValidationErrors(err)
	}
	return true, nil
}
