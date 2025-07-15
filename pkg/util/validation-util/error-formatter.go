package validation_util

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// FormatValidationErrors formats validation errors into a slice of maps.
// Each map contains the field name and the corresponding error message.
func FormatValidationErrors(err error) []map[string]string {
	var errors []map[string]string

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			// Customize the message based on tag
			var message string
			switch fe.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", fe.Field())
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", fe.Field())
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
			default:
				message = fmt.Sprintf("%s is not valid", fe.Field())
			}

			errors = append(errors, map[string]string{
				"field":   fe.Field(),
				"message": message,
			})
		}
	}

	return errors
}
