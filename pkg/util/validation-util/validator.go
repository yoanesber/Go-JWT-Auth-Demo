package validation_util

import (
	"reflect"
	"strings"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

var (
	once     sync.Once
	validate *validator.Validate
)

// Init initializes the validator and registers custom validations.
func Init() bool {
	isSuccess := true
	once.Do(func() {
		validate = validator.New()

		// Register tag name function to use JSON field names if available
		// instead of struct field names
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			tag := fld.Tag.Get("json")
			if tag == "-" || tag == "" {
				return fld.Name // fallback ke nama field struct
			}
			return strings.Split(tag, ",")[0]
		})
	})

	return isSuccess
}

// GetValidator returns the initialized validator instance.
func GetValidator() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}

// ClearValidator clears the validator instance.
// This function can be used to reset the validator for re-initialization.
func ClearValidator() {
	once = sync.Once{} // Reset the once to allow re-initialization
	validate = nil     // Clear the validator instance
}
