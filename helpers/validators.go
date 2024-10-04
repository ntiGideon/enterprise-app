package helpers

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

func RequestValidators(entity interface{}) error {
	validate := validator.New()
	err := validate.Struct(entity)
	if err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		return errs
	}
	return nil
}
