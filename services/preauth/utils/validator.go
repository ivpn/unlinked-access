package utils

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

func NewValidator() Validator {
	v := Validator{validator.New()}

	return v
}

func ValidateUUID(uuid string) bool {
	err := validator.New().Var(uuid, "required,uuid")
	return err == nil
}
