package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

	return &CustomValidator{validator: v}
}

// Validate implements the echo.Validator interface.
func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Return a 400 Bad Request error if validation fails.
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
