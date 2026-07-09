package validations

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Valida que la contraseña tenga minimo 8 caracteres, un numero y un caracter especial
var ValidarContrasena validator.Func = func(fl validator.FieldLevel) bool {
	contrasena := fl.Field().String()

	min8 := len(contrasena) >= 8
	numero := regexp.MustCompile(`[0-9]`).MatchString(contrasena)
	especial := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(contrasena)

	return min8 && numero && especial
}
