package validator

import "github.com/go-playground/validator"

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	return &Validator{
		validate: v,
	}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	err := v.validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}
