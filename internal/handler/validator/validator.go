package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.socket9.com/sritrang/super-driver/backend/internal/handler/common"
)

// CustomValidator adapts go-playground/validator to Echo's Validator interface.
type CustomValidator struct {
	v *validator.Validate
}

// New creates a CustomValidator and registers any custom validation rules.
func New() *CustomValidator {
	v := validator.New()
	return &CustomValidator{v: v}
}

// Validate implements echo.Validator. Returns common.Error on failure so the error
// handler maps it to HTTP 400 with field-level sub-errors.
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.v.Struct(i); err != nil {
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			return common.NewInternalServerError()
		}
		var errs []common.SubError
		for _, e := range ve {
			errs = append(errs, common.SubError{
				Field:   e.Field(),
				Message: e.Tag(),
			})
		}
		return common.NewBadRequest(errs)
	}
	return nil
}

// Register adds this validator to the Echo instance.
func Register(e *echo.Echo) {
	e.Validator = New()
}
