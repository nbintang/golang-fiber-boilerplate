package validator

import "github.com/go-playground/validator/v10"

type serviceImpl struct {
	validator *validator.Validate
}

func NewService() Service {
	return &serviceImpl{validator: validator.New()}
}

func (v *serviceImpl) Struct(i any) error {
	return v.validator.Struct(i)
}
