package validator

import (
	validatorGo "github.com/go-playground/validator/v10"
	"regexp"
)

func registerCustomValidator(v *validatorGo.Validate) {
	if err := v.RegisterValidation("name", ValidateName); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("email_address", ValidateEmail); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("address", ValidateAddress); err != nil {
		panic(err)
	}
}

func ValidateName(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile(`^[a-zA-Z\s,. ]*$`)
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func ValidateEmail(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+")
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func ValidateAddress(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile(`^[a-zA-Z0-9-',.()\s\/]+$`)
		return regex.MatchString(fl.Field().String())
	}
	return true
}
