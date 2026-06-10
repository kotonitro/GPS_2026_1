package promociones

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// verifica el JSON antes de que llegue al controlador
func ValidarPromocionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var promo Promocion
		
		// Lee el JSON
		if err := c.ShouldBindJSON(&promo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Datos inválidos: " + err.Error()})
			c.Abort() //  detiene la ejecución 
			return
		}

		// validacion 
		switch promo.Tipo {
		case "NXM":
			if promo.Lleva <= 0 || promo.Paga <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Para promociones NXM, 'lleva' y 'paga' deben ser mayores a 0"})
				c.Abort()
				return
			}
		case "porcentaje", "precio_fijo":
			if promo.Descuento <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Para esta promoción el descuento debe ser mayor a 0"})
				c.Abort()
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Tipo de promoción no válida"})
			c.Abort()
			return
		}

		c.Set("promocionValidada", promo)
		c.Next()
	}
}

func ValidarRutas(api *gin.RouterGroup) {
	grupo := api.Group("/promociones")
	{
		grupo.POST("/promocion", ValidarPromocionMiddleware(), CrearPromocion)
		grupo.PUT("/promocion/:id", ValidarPromocionMiddleware(), ActualizarPromociones)
	}
}

func CrearPromocionMiddleware(c *gin.Context) {
	// promoción que el middleware ya validó y guardó
	promoData, _ := c.Get("promocionValidada")
	nuevaPromocion := promoData.(Promocion)

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	err := GuardarPromocion(db, &nuevaPromocion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "No se pudo guardar la promoción en la BDD"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":   "Promoción creada exitosamente",
		"promocion": nuevaPromocion,
	})
}

func ActualizarPromocionesMiddleware(c *gin.Context) {
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)
	id := c.Param("id")

	// promoción validada del middleware
	promoData, _ := c.Get("promocionValidada")
	datosNuevos := promoData.(Promocion)

	promocionActualizada, err := ActualizarPromocion(db, id, &datosNuevos)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"ERROR": "No se pudo actualizar la promoción"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"mensaje":   "Promoción actualizada",
		"promocion": promocionActualizada,
	})
}