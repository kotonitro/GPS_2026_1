package empleados

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type EmpleadoController struct {
	db *gorm.DB
}

func NewEmpleadoController(db *gorm.DB) *EmpleadoController {
	return &EmpleadoController{db: db}
}

func (ctrl *EmpleadoController) GetEmpleadosController(c *gin.Context) {
	listaEmpleados, err := GetEmpleados(ctrl.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener la lista de empleados."})
		return
	}

	c.JSON(http.StatusOK, listaEmpleados)
}

func (ctrl *EmpleadoController) GetEmpleadoByIDController(c *gin.Context) {

	id := c.Param("id")

	empleado, err := GetEmpleadosByID(ctrl.db, id)
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El empleado no existe."})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar el empleado en la base de datos."})
		return
	}

	c.JSON(http.StatusOK, empleado)
}

type CreateEmpleadoInput struct {
	Rut        string `json:"rut" binding:"required,rut_valido"`
	Usuario    string `json:"usuario" binding:"required"`
	Contrasena string `json:"contrasena" binding:"required,contrasena_segura"`
	Telefono   string `json:"telefono"`
	Rol        string `json:"rol" binding:"required"`
}

func (ctrl *EmpleadoController) CreateEmpleadoController(c *gin.Context) {
	var input CreateEmpleadoInput

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al procesar la contraseña."})
		return
	}

	nuevoEmpleado := Empleado{
		Rut:        input.Rut,
		Usuario:    input.Usuario,
		Contrasena: string(hashContrasena),
		Telefono:   input.Telefono,
		Rol:        input.Rol,
	}

	err = CreateEmpleado(ctrl.db, &nuevoEmpleado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "No se pudo registrar el empleado.",
			"detalle": "El RUT o el Usuario ya se encuentran registrados en el sistema.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":     "Empleado creado exitosamente.",
		"id_empleado": nuevoEmpleado.ID,
	})

}
