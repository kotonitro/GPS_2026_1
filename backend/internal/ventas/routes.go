package ventas

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ConfigurarRutas(api *gin.RouterGroup, db *gorm.DB, authMiddleware gin.HandlerFunc) {

	rutasVentas := api.Group("/ventas")

	rutasVentas.Use(authMiddleware)
	{
		rutasVentas.POST("", CrearVenta)
		rutasVentas.GET("", GetVentas)
		rutasVentas.GET("/:id", GetVentaByID)

		rutasAdmin := rutasVentas.Group("")
		rutasAdmin.Use(auth.AdminMiddleware())
		{
			rutasAdmin.PUT("/:id", UpdateVenta)
			rutasAdmin.DELETE("/:id", DeleteVenta)
		}
	}
}
