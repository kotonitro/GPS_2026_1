package validations

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ValidationsConfig() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			nombre := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if nombre == "-" {
				return ""
			}
			return nombre
		})

		v.RegisterValidation("rut_valido", ValidarRUT)
		v.RegisterValidation("contrasena_segura", ValidarContrasena)
	}
}
