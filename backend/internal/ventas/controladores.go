package ventas

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CrearVenta(c *gin.Context) {
	// conexion a la base de datos
	db, existe := c.Get("db")
	if !existe {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sin conexión a la base de datos"})
		return
	}

	
	var nuevaVenta Venta

	if err := c.ShouldBindJSON(&nuevaVenta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if err := GuardarVenta(db.(*gorm.DB), &nuevaVenta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar en la BD"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Venta creada con éxito",
		"venta":   nuevaVenta,
	})
}