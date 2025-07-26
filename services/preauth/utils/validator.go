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
