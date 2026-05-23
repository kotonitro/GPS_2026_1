package inventario

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CrearProducto(c *gin.Context) {
	var nuevoProducto Producto

	// 1. Lee el JSON
	if err := c.ShouldBindJSON(&nuevoProducto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos"})
		return
	}

	// 2. Extrae la conexión a la base de datos del Middleware
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	// 3. LLamado al repositorio
	err := GuardarProducto(db, &nuevoProducto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo guardar el producto"})
		return
	}

	// 4. Entrega mensaje
	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Producto agregado",
		"producto": nuevoProducto,
	})
}
