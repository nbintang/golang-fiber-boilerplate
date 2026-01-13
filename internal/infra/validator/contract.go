package validator

type Service interface {
	Struct(i any) error
}
