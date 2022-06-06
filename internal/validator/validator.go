package validating

import (
	"gopkg.in/go-playground/validator.v9"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

func (v *Validator) RegisterRules(maxWorkers int) error {
	err := v.validator.RegisterValidation("worker", func(fl validator.FieldLevel) bool {
		return int(fl.Field().Int()) <= maxWorkers
	})
	return err
}

func(v *Validator) ValidateStruct(data interface{}) error {
	if err := v.validator.Struct(data); err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}
