package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator package and is used by Echo for request validation.
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator initializes the validator, registers any custom validations,
// and returns a new instance of CustomValidator.
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Example of registering a custom validation:
	// v.RegisterValidation("currency", validCurrency)
	// You can add more custom validations here as needed.
	// Register custom validations
	v.RegisterValidation("regex", validCategoryName)

	return &CustomValidator{validator: v}
}

// Validate implements the echo.Validator interface.
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

// If needed, you can add helper functions for additional validations here.
// For instance, if you want a custom "currency" validation:
//
// var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
//     // Add your logic to check if the currency is supported, e.g.,
//     // return utils.IsSupportedCurrency(fl.Field().String())
//     return true // Simplified example
// }

// validCategoryName is a custom validation function for the "regex" tag
var validCategoryName validator.Func = func(fl validator.FieldLevel) bool {
	// Regex pattern: only letters and spaces
	re := regexp.MustCompile(`^[A-Za-z ]+$`)
	return re.MatchString(fl.Field().String())
}
