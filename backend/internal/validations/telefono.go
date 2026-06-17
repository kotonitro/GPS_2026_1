package validations

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var ValidarTelefono validator.Func = func(fl validator.FieldLevel) bool {
	telefono := fl.Field().String()

	//912345678, 9 1234 5678,      9-1234-5678, +56912345678, +569 1234 5678
	patron := regexp.MustCompile(`^(\+56)?[\s.-]?9[\s.-]?\d{4}[\s.-]?\d{4}$`)
	
	return patron.MatchString(telefono)
}
