package clientes

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ClienteController struct {
	db *gorm.DB
}

func NewClienteController(db *gorm.DB) *ClienteController {
	return &ClienteController{db: db}
}

func (ctrl *ClienteController) GetClientesController(c *gin.Context) {
	listaClientes, err := GetClientes(ctrl.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener la lista de clientes."})
		return
	}

	c.JSON(http.StatusOK, listaClientes)
}

func (ctrl *ClienteController) GetClienteByIDController(c *gin.Context) {
	id := c.Param("id")

	cliente, err := GetClienteByID(ctrl.db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El cliente solicitado no existe."})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar el cliente en la base de datos."})
		return
	}

	c.JSON(http.StatusOK, cliente)
}

type CreateClienteInput struct {
	Nombre      string   `json:"nombre" binding:"required"`
	Rut         string   `json:"rut" binding:"required,rut_valido"`
	Telefono    string   `json:"telefono" binding:"required,telefono_valido"`
	FiadoActual *float64 `json:"fiado_actual"`
	FiadoMaximo *float64 `json:"fiado_maximo"`
}

func (ctrl *ClienteController) CreateClienteController(c *gin.Context) {
	var input CreateClienteInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errores := ValidationErrorsFormat(err)
		c.JSON(http.StatusBadRequest, gin.H{"errores": errores})
		return
	}

	input.Telefono = normalizarTelefono(input.Telefono)

	nuevoCliente := Cliente{
		Nombre:   input.Nombre,
		Rut:      input.Rut,
		Telefono: input.Telefono,
	}

	if input.FiadoMaximo != nil {
		nuevoCliente.FiadoMaximo = *input.FiadoMaximo
	} else {
		nuevoCliente.FiadoMaximo = 20000.0
	}

	if input.FiadoActual != nil {
		if *input.FiadoActual > nuevoCliente.FiadoMaximo {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "No se pudo registrar el cliente.",
				"detalle": "El saldo de fiado no puede superar el límite máximo de $20.000.",
			})
			return
		}
		nuevoCliente.FiadoActual = *input.FiadoActual
	}

	err := CreateCliente(ctrl.db, &nuevoCliente)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "rut") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo registrar el cliente.",
				"detalle": "El RUT ya se encuentra registrado en el sistema.",
			})
			return
		}

		if strings.Contains(errMsg, "telefono") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo registrar el cliente.",
				"detalle": "El teléfono ya esta registrado en el sistema.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "No se pudo registrar el cliente.",
			"detalle": "Error interno al crear el cliente.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":    "Cliente creado exitosamente.",
		"id_cliente": nuevoCliente.ID,
	})
}

func (ctrl *ClienteController) DeleteClienteByIDController(c *gin.Context) {
	id := c.Param("id")

	err := DeleteClienteByID(ctrl.db, id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El cliente que intenta eliminar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al eliminar el cliente."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Cliente eliminado exitosamente."})
}

type UpdateClienteInput struct {
	Nombre      *string  `json:"nombre"`
	Rut         *string  `json:"rut" binding:"omitempty,rut_valido"`
	Telefono    *string  `json:"telefono" binding:"omitempty,telefono_valido"`
	FiadoActual *float64 `json:"fiado_actual"`
	FiadoMaximo *float64 `json:"fiado_maximo"`
}

func (ctrl *ClienteController) UpdateClienteByIDController(c *gin.Context) {
	id := c.Param("id")

	var input UpdateClienteInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errores := ValidationErrorsFormat(err)
		c.JSON(http.StatusBadRequest, gin.H{"errores": errores})
		return
	}

	if input.Nombre == nil && input.Rut == nil && input.Telefono == nil && input.FiadoActual == nil && input.FiadoMaximo == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Se requiere al menos un campo válido para modificar.",
		})
		return
	}

	clienteActual, err := GetClienteByID(ctrl.db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El cliente a modificar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar el cliente."})
		return
	}

	fiadoMaximo := clienteActual.FiadoMaximo
	if input.FiadoMaximo != nil {
		fiadoMaximo = *input.FiadoMaximo
	}

	if input.FiadoActual != nil {
		if *input.FiadoActual > fiadoMaximo {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "No se pudo actualizar el cliente.",
				"detalle": "El saldo de fiado no puede superar el límite máximo.",
			})
			return
		}
	}

	if input.Telefono != nil {
		telefonoNormalizado := normalizarTelefono(*input.Telefono)
		input.Telefono = &telefonoNormalizado
	}

	err = UpdateClienteByID(ctrl.db, id, input)

	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "rut") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo actualizar el cliente.",
				"detalle": "El RUT ya se encuentra registrado en el sistema.",
			})
			return
		}

		if strings.Contains(errMsg, "telefono") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo actualizar el cliente.",
				"detalle": "El teléfono ya se encuentra registrado en el sistema.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al modificar el cliente."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Cliente modificado exitosamente."})
}

func (ctrl *ClienteController) GetClienteByRutController(c *gin.Context) {
	rut := c.Param("rut")

	cliente, err := GetClienteByRut(ctrl.db, rut)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró cliente con ese RUT."})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar el cliente."})
		return
	}

	c.JSON(http.StatusOK, cliente)
}

func (ctrl *ClienteController) GetClientesByNombreController(c *gin.Context) {
	nombre := c.Query("nombre")

	if nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre es requerido."})
		return
	}

	clientes, err := GetClientesByNombre(ctrl.db, nombre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar clientes."})
		return
	}

	if len(clientes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron clientes con ese nombre.", "resultados": []Cliente{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resultados": clientes,
		"cantidad":   len(clientes),
	})
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
//ValidationErrorsFormat formatea los errores de validación de Gin y devuelve los mensajes de error correspondientes.
func ValidationErrorsFormat(err error) map[string]string {
	var errs validator.ValidationErrors
	mensajes := make(map[string]string)

	if errors.As(err, &errs) {
		for _, f := range errs {
			switch f.Tag() {
			case "rut_valido":
				mensajes[f.Field()] = "El RUT ingresado no es válido."
			case "telefono_valido":
				mensajes[f.Field()] = "El formato del teléfono no es válido"
			case "required":
				fieldName := f.Field()
				if fieldName == "Telefono" || fieldName == "telefono" {
					mensajes[fieldName] = "Se requiere un numero de teléfono para llamar."
				} else {
					mensajes[fieldName] = "Este campo es obligatorio."
				}
			default:
				mensajes[f.Field()] = "El formato ingresado no es válido."
			}
		}
		return mensajes
	}

	mensajes["error"] = "El cuerpo de la petición es inválido."
	return mensajes
}
