package user

import "gopkg.in/go-playground/validator.v9"

func validationErrors(err error) []error {
	var errs []error
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, FieldError(err))
	}
	return errs
}
