package clientes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ClienteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verifica que la conexión a la base de datos esté disponible.
		dbInstance, ok := c.Get("db")
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Error": "Conexión a la base de datos no disponible",
			})
			return
		}

		if _, ok := dbInstance.(*gorm.DB); !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Error": "Conexión a la base de datos inválida",
			})
			return
		}

		// el request pasó por el middleware de clientes.
		c.Set("clienteMiddlewareActivo", true)

		c.Next()
	}
}
