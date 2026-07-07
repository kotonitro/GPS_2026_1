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
	ProductoID string  `json:"producto_id" binding:"required,uuid"`
	Tipo      string  `json:"tipo" binding:"required"`
	Lleva     int     `json:"lleva"`
	Paga      int     `json:"paga"`
	Descuento float64 `json:"descuento"`
	FechaInicio *time.Time `json:"fecha_inicio"`
	FechaFin    *time.Time `json:"fecha_fin"`
}

type UpdatePromocionInput struct {
	Tipo      *string  `json:"tipo"`
	Lleva     *int     `json:"lleva"`
	Paga      *int     `json:"paga"`
	Descuento *float64 `json:"descuento"`
	FechaInicio *time.Time `json:"fecha_inicio"`
	FechaFin    *time.Time `json:"fecha_fin"`
}

func (ctrl *PromocionController) CreatePromocionController(c *gin.Context) {
	var input CreatePromocionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	switch input.Tipo {
	case "NXM":
		if input.Lleva <= 0 || input.Paga <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Para promociones NXM lleva y paga deben ser mayores a 0"})
			return
		}
	case "porcentaje", "precio_fijo":
		if input.Descuento <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Para la promocion este valor debe ser mayor a 0"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de promocion no valida"})
		return
	}

	nuevaPromocion := Promocion{
		ProductoID: input.ProductoID,
		Tipo:      input.Tipo,
		Lleva:     input.Lleva,
		Paga:      input.Paga,
		Descuento: input.Descuento,
		FechaInicio: input.FechaInicio, 
		FechaFin:    input.FechaFin,
	}

	err := GuardarPromocion(ctrl.db, &nuevaPromocion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la promocion en la base de datos."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":   "Promocion creada exitosamente.",
		"promocion": nuevaPromocion,
	})
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
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de promocion no valida"})
			return
		}
	}

	var datosActualizados Promocion
	if input.Tipo != nil {
		datosActualizados.Tipo = *input.Tipo
	}
	if input.Lleva != nil {
		datosActualizados.Lleva = *input.Lleva
	}
	if input.Paga != nil {
		datosActualizados.Paga = *input.Paga
	}
	if input.Descuento != nil {
		datosActualizados.Descuento = *input.Descuento
	}

	promocionActualizada, err := ActualizarPromocion(ctrl.db, id, &datosActualizados)
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
