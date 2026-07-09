package promociones

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// PromocionController maneja las peticiones web
type PromocionController struct {
	db *gorm.DB
}

// constructor del controlador
func NewPromocionController(db *gorm.DB) *PromocionController {
	return &PromocionController{db: db}
}

type CreatePromocionInput struct {
	ProductoID     *string    `json:"producto_id"`
	ProductosCombo []string   `json:"productos_combo"`
	Tipo           string     `json:"tipo" binding:"required"`
	Lleva          *int       `json:"lleva"`
	Paga           *int       `json:"paga"`
	Descuento      *float64   `json:"descuento"`
	FechaInicio    *time.Time `json:"fecha_inicio"`
	FechaFin       *time.Time `json:"fecha_fin"`
}

type UpdatePromocionInput struct {
	ProductoID     *string    `json:"producto_id"`
	Tipo           *string    `json:"tipo"`
	ProductosCombo []string   `json:"productos_combo"`
	Lleva          *int       `json:"lleva"`
	Paga           *int       `json:"paga"`
	Descuento      *float64   `json:"descuento"`
	FechaInicio    *time.Time `json:"fecha_inicio"`
	FechaFin       *time.Time `json:"fecha_fin"`
}

func (ctrl *PromocionController) CreatePromocionController(c *gin.Context) {
	var input CreatePromocionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error de estructura: " + err.Error()})
		return
	}

	switch input.Tipo {
	case "NXM":
		if input.Lleva == nil || *input.Lleva <= 0 || input.Paga == nil || *input.Paga <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Para NXM, lleva y paga son requeridos y > 0"})
			return
		}
	case "porcentaje", "precio_fijo", "COMBO":
		if input.Descuento == nil || *input.Descuento <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El descuento debe ser mayor a 0"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de promocion no valida"})
		return
	}

	if input.Tipo != "COMBO" {
		if input.ProductoID == nil || *input.ProductoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Debe seleccionar un producto para esta promoción"})
			return
		}
		input.ProductosCombo = []string{}
	}

	// Guardar
	nuevaPromocion := Promocion{
		ProductoID:     input.ProductoID,
		ProductosCombo: input.ProductosCombo,
		Tipo:           input.Tipo,
		Lleva:          input.Lleva,
		Paga:           input.Paga,
		Descuento:      input.Descuento,
		FechaInicio:    input.FechaInicio,
		FechaFin:       input.FechaFin,
	}

	err := GuardarPromocion(ctrl.db, &nuevaPromocion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la promocion."})
		return
	}

	c.JSON(http.StatusCreated, nuevaPromocion)
}

func (ctrl *PromocionController) GetPromocionesController(c *gin.Context) {
	promociones, err := ObtenerTodasLasPromociones(ctrl.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener la lista de promociones."})
		return
	}

	c.JSON(http.StatusOK, promociones)
}

func (ctrl *PromocionController) GetPromocionByIDController(c *gin.Context) {
	id := c.Param("id")

	promocion, err := ObtenerPromocionPorID(ctrl.db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "La promocion solicitada no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar la promocion en la base de datos."})
		return
	}

	c.JSON(http.StatusOK, promocion)
}

func (ctrl *PromocionController) UpdatePromocionByIDController(c *gin.Context) {
	id := c.Param("id")
	var input UpdatePromocionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if input.Tipo != nil {
		switch *input.Tipo {
		case "NXM":
			if (input.Lleva != nil && *input.Lleva <= 0) || (input.Paga != nil && *input.Paga <= 0) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Para promociones NXM lleva y paga deben ser mayores a 0"})
				return
			}
		case "porcentaje", "precio_fijo":
			if input.Descuento != nil && *input.Descuento <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Para la promocion este valor debe ser mayor a 0"})
				return
			}
		case "COMBO":
			if input.Descuento != nil && *input.Descuento <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "El precio final del combo debe ser mayor a 0"})
				return
			}
			if len(input.ProductosCombo) < 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Un combo debe tener al menos 2 productos"})
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de promocion no valida"})
			return
		}

		if *input.Tipo == "COMBO" {
			input.ProductoID = nil
		} else {
			if input.ProductoID == nil || *input.ProductoID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Debe seleccionar un producto para esta promoción"})
				return
			}
			input.ProductosCombo = []string{}
		}
	}
	datosActualizados := make(map[string]interface{})

	if input.Tipo != nil {
		datosActualizados["tipo"] = *input.Tipo
		// Agregamos los campos de producto siempre que se actualice el tipo
		datosActualizados["producto_id"] = input.ProductoID
		datosActualizados["productos_combo"] = input.ProductosCombo
	}
	if input.Lleva != nil {
		datosActualizados["lleva"] = *input.Lleva
	}
	if input.Paga != nil {
		datosActualizados["paga"] = *input.Paga
	}
	if input.Descuento != nil {
		datosActualizados["descuento"] = *input.Descuento
	}
	
	datosActualizados["fecha_inicio"] = input.FechaInicio
	datosActualizados["fecha_fin"] = input.FechaFin

	promocionActualizada, err := ActualizarPromocion(ctrl.db, id, datosActualizados)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "La promocion a modificar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al modificar la promocion."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":   "Promocion modificada exitosamente.",
		"promocion": promocionActualizada,
	})
}

func (ctrl *PromocionController) DeletePromocionController(c *gin.Context) {
	id := c.Param("id")

	err := EliminarPromocion(ctrl.db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "La promocion que intenta eliminar no existe."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al eliminar la promocion."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Promocion eliminada exitosamente."})
}
