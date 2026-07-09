package validations

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Verifica que el RUT sea un RUT valido en Chile
var ValidarRUT validator.Func = func(fl validator.FieldLevel) bool {
	rutOriginal := fl.Field().String()

	rut := strings.ReplaceAll(rutOriginal, ".", "")
	rut = strings.ReplaceAll(rut, "-", "")
	rut = strings.ToUpper(rut)

	if len(rut) < 8 {
		return false
	}

	cuerpo := rut[:len(rut)-1]
	dvIngresado := string(rut[len(rut)-1])

	suma := 0
	multiplicador := 2

	for i := len(cuerpo) - 1; i >= 0; i-- {

		digito, err := strconv.Atoi(string(cuerpo[i]))
		if err != nil {
			return false
		}

		suma += digito * multiplicador

		multiplicador++
		if multiplicador > 7 {
			multiplicador = 2
		}
	}

	resto := suma % 11
	resultado := 11 - resto

	dvCalculado := ""
	if resultado == 11 {
		dvCalculado = "0"
	} else if resultado == 10 {
		dvCalculado = "K"
	} else {
		dvCalculado = strconv.Itoa(resultado)
	}

	return dvCalculado == dvIngresado
}
