package cajas

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CajaController struct {
	db *gorm.DB
}

func NewCajaController(db *gorm.DB) *CajaController {
	return &CajaController{db: db}
}

func (ctrl *CajaController) GetCajasController(c *gin.Context) {
	listaCajas, err := GetCajas(ctrl.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener la lista de cajas."})
		return
	}

	c.JSON(http.StatusOK, listaCajas)
}

func (ctrl *CajaController) GetCajaByIDController(c *gin.Context) {

	id := c.Param("id")

	caja, err := GetCajaByID(ctrl.db, id)
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "La caja solicitada no existe."})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar la caja en la base de datos."})
		return
	}

	c.JSON(http.StatusOK, caja)
}

type CreateCajaInput struct {
	Nombre       string `json:"nombre" binding:"required"`
	Ubicacion    string `json:"ubicacion" binding:"required"`
	Activo       bool   `json:"activo"`
	SaldoInicial uint   `json:"saldo_inicial"`
	SaldoFinal   uint   `json:"saldo_final"`
}

func (ctrl *CajaController) CreateCajaController(c *gin.Context) {
	var input CreateCajaInput

	if err := c.ShouldBindJSON(&input); err != nil {

		errores := ValidationErrorsFormat(err)
		c.JSON(http.StatusBadRequest, gin.H{"errores": errores})
		return
	}

	nuevaCaja := Caja{
		Nombre:       input.Nombre,
		Ubicacion:    input.Ubicacion,
		Activo:       input.Activo,
		SaldoInicial: input.SaldoInicial,
		SaldoFinal:   input.SaldoFinal,
	}

	err := CreateCaja(ctrl.db, &nuevaCaja)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "No se pudo registrar la caja.",
			"detalle": "La ubicación ya se encuentra registrada en el sistema.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Caja creada exitosamente.",
		"id_caja": nuevaCaja.ID,
	})

}

func (ctrl *CajaController) DeleteCajaByIDController(c *gin.Context) {
	id := c.Param("id")

	err := DeleteCajaByID(ctrl.db, id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "La caja que intenta eliminar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al eliminar la caja."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Caja eliminada exitosamente."})
}

type UpdateCajaInput struct {
	Nombre       *string `json:"nombre"`
	Ubicacion    *string `json:"ubicacion"`
	SaldoInicial *uint   `json:"saldo_inicial"`
	SaldoFinal   *uint   `json:"saldo_final"`
	Activo       *bool   `json:"activo"`
}

func (ctrl *CajaController) UpdateCajaByIDController(c *gin.Context) {
	id := c.Param("id")

	var input UpdateCajaInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errores := ValidationErrorsFormat(err)
		c.JSON(http.StatusBadRequest, gin.H{"errores": errores})
		return
	}

	if input.Nombre == nil && input.Ubicacion == nil && input.Activo == nil && input.SaldoInicial == nil && input.SaldoFinal == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Se requiere al menos un campo válido para modificar.",
		})
		return
	}

	err := UpdateCajaByID(ctrl.db, id, input)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "La caja a modificar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error interno al modificar la caja."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Caja modificada exitosamente."})
}

func ValidationErrorsFormat(err error) map[string]string {
	var errs validator.ValidationErrors
	mensajes := make(map[string]string)

	if errors.As(err, &errs) {
		for _, f := range errs {
			switch f.Tag() {
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
