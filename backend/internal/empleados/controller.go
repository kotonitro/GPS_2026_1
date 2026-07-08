package empleados

import (
	"errors"
	"net/http"
	"strings"

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
	Rut        string `json:"rut" binding:"required,rut_valido"`
	Nombre     string `json:"nombre" binding:"required"`
	Usuario    string `json:"usuario" binding:"required"`
	Contrasena string `json:"contrasena" binding:"required,contrasena_segura"`
	Telefono   string `json:"telefono" binding:"required,telefono_valido"`
	RolID      string `json:"id_rol" binding:"required"`
	Activo     bool   `json:"activo"`
}

func (ctrl *EmpleadoController) CreateEmpleadoController(c *gin.Context) {
	var input CreateEmpleadoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errores := ValidationErrorsFormat(err)
		c.JSON(http.StatusBadRequest, gin.H{"errores": errores})
		return
	}

	idSolCtx, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar al administrador de la sesión."})
		return
	}
	solicitante, _ := GetEmpleadoByID(ctrl.db, idSolCtx.(string))
	esSuperAdmin := solicitante.Rol.Nombre == "Admin"

	rolAsignar, err := GetRolByID(ctrl.db, input.RolID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El rol asignado no existe."})
		return
	}

	if rolAsignar.Nombre == "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "El rol 'Admin' no puede ser asignado a otros empleados."})
		return
	}

	if rolAsignar.EsAdmin && !esSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para crear cuentas de administrador."})
		return
	}

	input.Telefono = normalizarTelefono(input.Telefono)

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
		RolID:      input.RolID,
		Activo:     input.Activo,
	}

	err = CreateEmpleado(ctrl.db, &nuevoEmpleado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el empleado.", "detalle": "RUT o Usuario duplicado."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensaje": "Empleado creado exitosamente.", "id_empleado": nuevoEmpleado.ID})
}

func (ctrl *EmpleadoController) DeleteEmpleadoByIDController(c *gin.Context) {
	idObj := c.Param("id")

	idSolCtx, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar al administrador de la sesión."})
		return
	}

	idSol := idSolCtx.(string)
	solicitante, _ := GetEmpleadoByID(ctrl.db, idSol)
	esSuperAdmin := solicitante.Rol.Nombre == "Admin"

	if idSol == idObj {
		c.JSON(http.StatusForbidden, gin.H{"error": "No puedes eliminar tu propia cuenta."})
		return
	}

	empleadoObj, err := GetEmpleadoByID(ctrl.db, idObj)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El empleado no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al verificar identidad."})
		return
	}

	if empleadoObj.Rol != nil && empleadoObj.Rol.EsAdmin && !esSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para eliminar a un administrador."})
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
	RolID      *string `json:"id_rol"`
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

	if input.Nombre == nil && input.Usuario == nil && input.Contrasena == nil && input.Telefono == nil && input.RolID == nil && input.Activo == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere al menos un campo válido para modificar."})
		return
	}

	idSolCtx, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar al usuario."})
		return
	}
	idSol := idSolCtx.(string)
	solicitante, _ := GetEmpleadoByID(ctrl.db, idSol)
	esSuperAdmin := solicitante.Rol.Nombre == "Admin"

	empleadoObj, err := GetEmpleadoByID(ctrl.db, idObj)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "El empleado a modificar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error interno al verificar el empleado."})
		return
	}

	if empleadoObj.Rol != nil && empleadoObj.Rol.EsAdmin {
		if idSol != idObj {
			if !esSuperAdmin {
				c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para modificar a un administrador."})
				return
			}
		} else {
			if input.Activo != nil && !*input.Activo {
				c.JSON(http.StatusForbidden, gin.H{"error": "No puedes desactivar tu propia cuenta."})
				return
			}
			if input.RolID != nil && *input.RolID != empleadoObj.RolID {
				nuevoRol, errRol := GetRolByID(ctrl.db, *input.RolID)
				if errRol != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "El nuevo rol asignado no existe."})
					return
				}
				if !nuevoRol.EsAdmin {
					c.JSON(http.StatusForbidden, gin.H{"error": "No puedes revocar tus propios privilegios de administrador."})
					return
				}
			}
		}
	}

	if input.RolID != nil && *input.RolID != empleadoObj.RolID {
		nuevoRol, errRol := GetRolByID(ctrl.db, *input.RolID)
		if errRol == nil {
			if nuevoRol.Nombre == "Admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "El rol 'Admin' no puede ser asignado a otros empleados."})
				return
			}

			if nuevoRol.EsAdmin && !esSuperAdmin {
				c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para otorgar el rol de administrador."})
				return
			}
		}
	}

	if input.Telefono != nil {
		telefonoNormalizado := normalizarTelefono(*input.Telefono)
		input.Telefono = &telefonoNormalizado
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

func normalizarTelefono(telefono string) string {
	telefonoLimpio := strings.NewReplacer(
		" ", "",
		"-", "",
		".", "",
		"+", "",
	).Replace(telefono)

	if strings.HasPrefix(telefonoLimpio, "56") {
		telefonoLimpio = telefonoLimpio[2:]
	}

	return telefonoLimpio
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

func (ctrl *EmpleadoController) GetRolesController(c *gin.Context) {
	listaRoles, err := GetRoles(ctrl.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener la lista de roles."})
		return
	}

	c.JSON(http.StatusOK, listaRoles)
}

func (ctrl *EmpleadoController) GetRolByIDController(c *gin.Context) {

	id := c.Param("id")

	rol, err := GetRolByID(ctrl.db, id)
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El rol solicitado no existe."})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar el rol en la base de datos."})
		return
	}

	c.JSON(http.StatusOK, rol)
}

type CreateRolInput struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion" binding:"required"`
	EsAdmin     bool   `json:"es_admin"`
}

func (ctrl *EmpleadoController) CreateRolController(c *gin.Context) {
	var input CreateRolInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errores := ValidationErrorsFormat(err)
		c.JSON(http.StatusBadRequest, gin.H{"errores": errores})
		return
	}

	idSolCtx, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar al administrador de la sesión."})
		return
	}
	solicitante, err := GetEmpleadoByID(ctrl.db, idSolCtx.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar permisos del solicitante."})
		return
	}
	esSuperAdmin := solicitante.Rol.Nombre == "Admin"

	if input.EsAdmin && !esSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para crear roles administrativos."})
		return
	}

	nuevoRol := Rol{
		Nombre:      input.Nombre,
		Descripcion: input.Descripcion,
		EsAdmin:     input.EsAdmin,
	}

	err = CreateRol(ctrl.db, &nuevoRol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "No se pudo registrar el rol.",
			"detalle": "El nombre ya se encuentra registrado en el sistema.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Rol creado exitosamente.",
		"id_rol":  nuevoRol.ID,
	})
}

func (ctrl *EmpleadoController) DeleteRolByIDController(c *gin.Context) {
	id := c.Param("id")

	idSolCtx, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar al administrador de la sesión."})
		return
	}
	solicitante, _ := GetEmpleadoByID(ctrl.db, idSolCtx.(string))
	esSuperAdmin := solicitante.Rol.Nombre == "Admin"

	rolObj, err := GetRolByID(ctrl.db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El rol que intenta eliminar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al verificar el rol."})
		return
	}

	if rolObj.Nombre == "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No se puede eliminar el rol principal del sistema."})
		return
	}

	if rolObj.EsAdmin && !esSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para eliminar roles administrativos."})
		return
	}

	err = DeleteRolByID(ctrl.db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al eliminar el rol."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Rol eliminado exitosamente."})
}

type UpdateRolInput struct {
	Nombre      *string `json:"nombre"`
	Descripcion *string `json:"descripcion"`
	EsAdmin     *bool   `json:"es_admin"`
}

func (ctrl *EmpleadoController) UpdateRolByIDController(c *gin.Context) {
	id := c.Param("id")
	var input UpdateRolInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errores := ValidationErrorsFormat(err)
		c.JSON(http.StatusBadRequest, gin.H{"errores": errores})
		return
	}

	if input.Nombre == nil && input.Descripcion == nil && input.EsAdmin == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere al menos un campo válido para modificar."})
		return
	}

	idSolCtx, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar al administrador de la sesión."})
		return
	}
	solicitante, _ := GetEmpleadoByID(ctrl.db, idSolCtx.(string))
	esSuperAdmin := solicitante.Rol.Nombre == "Admin"

	rolObj, err := GetRolByID(ctrl.db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "El rol a modificar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error interno al verificar el rol."})
		return
	}

	if rolObj.Nombre == "Admin" {
		if input.EsAdmin != nil && !*input.EsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "No puedes revocar los privilegios del rol principal."})
			return
		}
		if input.Nombre != nil && *input.Nombre != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "No puedes cambiar el nombre del rol principal."})
			return
		}
	}

	if !esSuperAdmin {
		if rolObj.EsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para modificar otros roles administrativos."})
			return
		}
		if input.EsAdmin != nil && *input.EsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para otorgar privilegios administrativos a un rol."})
			return
		}
	}

	err = UpdateRolByID(ctrl.db, id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error interno al modificar el rol."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Rol modificado exitosamente."})
}
