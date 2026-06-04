package empleados

import (
	"backend/internal/database"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type CrearEmpleadoInput struct {
	Rut        string `json:"rut" binding:"required,rut_valido"`
	Usuario    string `json:"usuario" binding:"required"`
	Contrasena string `json:"contrasena" binding:"required,contrasena_segura"`
	Telefono   string `json:"telefono"`
	Rol        string `json:"rol" binding:"required"`
}

func CrearEmpleado(c *gin.Context) {
	var input CrearEmpleadoInput

	if err := c.ShouldBindJSON(&input); err != nil {

		var errs validator.ValidationErrors
		if errors.As(err, &errs) {

			mensajes := make(map[string]string)

			for _, f := range errs {

				switch f.Tag() {

				case "required":
					mensajes[f.Field()] = "Este campo es obligatorio."

				case "rut_valido":
					mensajes[f.Field()] = "El RUT ingresado no es válido."

				case "contrasena_segura":
					mensajes[f.Field()] = "La contraseña debe tener al menos 8 caracteres, 1 número y 1 caracter especial."

				default:
					mensajes[f.Field()] = "El formato ingresado no es válido."
				}
			}

			c.JSON(http.StatusBadRequest, gin.H{"errores": mensajes})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "El cuerpo de la petición es inválido."})
		return
	}

	hashContrasena, err := bcrypt.GenerateFromPassword([]byte(input.Contrasena), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrió un error interno al procesar la contraseña."})
		return
	}

	nuevoEmpleado := Empleado{
		Rut:        input.Rut,
		Usuario:    input.Usuario,
		Contrasena: string(hashContrasena),
		Telefono:   input.Telefono,
		Rol:        input.Rol,
	}

	err = GuardarEmpleado(database.DB, &nuevoEmpleado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "No se pudo registrar el empleado.",
			"detalle": "El RUT o el Usuario ya se encuentran registrados en el sistema.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":     "Empleado creado exitosamente",
		"id_empleado": nuevoEmpleado.ID,
	})

}
