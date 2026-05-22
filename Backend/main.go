package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Crea el enrutador de Gin con la configuración por defecto
	r := gin.Default()

	// Creamos una ruta de prueba
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"mensaje": "¡El backend con Gin y Docker está funcionando",
		})
	})

	// Encendemos el servidor en el puerto 8080
	r.Run(":8080")
}
