package inits

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func ValidateInit() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func MyValidate(s interface{}) error {
	if err := Validate.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
		}

		errMsgs := make([]string, 0)

		for _, elm := range err.(validator.ValidationErrors) {
			param := elm.Param()
			if param != "" {
				param = fmt.Sprintf(" %s", param)
			}
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s] is '%s%s'",
				elm.Field(),
				elm.Tag(),
				param,
			))
		}

		return errors.New(strings.Join(errMsgs, ", "))
	}

	return nil
}
