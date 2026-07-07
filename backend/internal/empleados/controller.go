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

	empleado, err := GetEmpleadoByID(ctrl.db, id)
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El empleado solicitado no existe."})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar el empleado en la base de datos."})
		return
	}

	c.JSON(http.StatusOK, empleado)
}

type CreateEmpleadoInput struct {
	Rut        string  `json:"rut" binding:"required,rut_valido"`
	Nombre     string  `json:"nombre" binding:"required"`
	Usuario    string  `json:"usuario" binding:"required"`
	Contrasena string  `json:"contrasena" binding:"required,contrasena_segura"`
	Telefono   *string `json:"telefono"`
	Rol        string  `json:"rol" binding:"required"`
	Activo     bool    `json:"activo"`
}

func (ctrl *EmpleadoController) CreateEmpleadoController(c *gin.Context) {
	var input CreateEmpleadoInput

	if err := c.ShouldBindJSON(&input); err != nil {

		errores := ValidationErrorsFormat(err)
		c.JSON(http.StatusBadRequest, gin.H{"errores": errores})
		return
	}

	hashContrasena, err := bcrypt.GenerateFromPassword([]byte(input.Contrasena), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al procesar la contraseña."})
		return
	}

	nuevoEmpleado := Empleado{
		Rut:        input.Rut,
		Nombre:     input.Nombre,
		Usuario:    input.Usuario,
		Contrasena: string(hashContrasena),
		Telefono:   input.Telefono,
		Rol:        input.Rol,
		Activo:     input.Activo,
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

func (ctrl *EmpleadoController) DeleteEmpleadoByIDController(c *gin.Context) {
	idObj := c.Param("id")

	idSolCtx, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar al administrador de la sesión."})
		return
	}
	idSol := idSolCtx.(string)

	if idSol == idObj {
		c.JSON(http.StatusForbidden, gin.H{"error": "No puedes eliminar tu propia cuenta de administrador."})
		return
	}

	empleadoObj, err := GetEmpleadoByID(ctrl.db, idObj)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El empleado que intenta eliminar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al verificar la identidad del empleado."})
		return
	}

	if empleadoObj.Rol == "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No puedes eliminar a otro administrador."})
		return
	}

	err = DeleteEmpleadoByID(ctrl.db, idObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al eliminar el empleado."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Empleado eliminado exitosamente."})
}

type UpdateEmpleadoInput struct {
	Nombre     *string `json:"nombre"`
	Usuario    *string `json:"usuario"`
	Contrasena *string `json:"contrasena" binding:"omitempty,contrasena_segura"`
	Telefono   *string `json:"telefono"`
	Rol        *string `json:"rol"`
	Activo     *bool   `json:"activo"`
}

func (ctrl *EmpleadoController) UpdateEmpleadoByIDController(c *gin.Context) {
	idObj := c.Param("id")

	var input UpdateEmpleadoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errores := ValidationErrorsFormat(err)
		c.JSON(http.StatusBadRequest, gin.H{"errores": errores})
		return
	}

	if input.Nombre == nil && input.Usuario == nil && input.Contrasena == nil && input.Telefono == nil && input.Rol == nil && input.Activo == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Se requiere al menos un campo válido para modificar.",
		})
		return
	}

	idSolCtx, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar al administrador de la sesión."})
		return
	}
	idSol := idSolCtx.(string)

	empleadoObj, err := GetEmpleadoByID(ctrl.db, idObj)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "El empleado a modificar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error interno al verificar el empleado."})
		return
	}

	if empleadoObj.Rol == "Admin" {
		if idSol != idObj {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para modificar los datos de otro administrador."})
			return
		} else {
			if input.Activo != nil && !*input.Activo {
				c.JSON(http.StatusForbidden, gin.H{"error": "No puedes desactivar tu propia cuenta."})
				return
			}

			if input.Rol != nil && *input.Rol != "Admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "No puedes revocar tus propios privilegios de administrador."})
				return
			}
		}
	}

	if input.Contrasena != nil {
		hashContrasena, err := bcrypt.GenerateFromPassword([]byte(*input.Contrasena), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al procesar la contraseña."})
			return
		}

		hashString := string(hashContrasena)
		input.Contrasena = &hashString
	}

	err = UpdateEmpleadoByID(ctrl.db, idObj, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error interno al modificar el empleado."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Empleado modificado exitosamente."})
}

func ValidationErrorsFormat(err error) map[string]string {
	var errs validator.ValidationErrors
	mensajes := make(map[string]string)

	if errors.As(err, &errs) {
		for _, f := range errs {
			switch f.Tag() {
			case "rut_valido":
				mensajes[f.Field()] = "El RUT ingresado no es válido."
			case "contrasena_segura":
				mensajes[f.Field()] = "La contraseña debe tener al menos 8 caracteres, 1 número y 1 caracter especial."
			case "required":
				mensajes[f.Field()] = "Este campo es obligatorio."
			default:
				mensajes[f.Field()] = "El formato ingresado no es válido."
			}
		}
		return mensajes
	}

	mensajes["error"] = "El cuerpo de la petición es inválido."
	return mensajes
}
