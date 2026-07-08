package ventas

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ConfigurarRutas(api *gin.RouterGroup, db *gorm.DB, authMiddleware gin.HandlerFunc) {

	ctrl := NewVentasController(db)

	rutasVentas := api.Group("/ventas")

	rutasVentas.Use(authMiddleware)
	{
		rutasVentas.POST("", ctrl.CrearVenta)
		rutasVentas.GET("", ctrl.GetVentas)
		rutasVentas.GET("/:id", ctrl.GetVentaByID)

		rutasAdmin := rutasVentas.Group("")
		rutasAdmin.Use(auth.AdminMiddleware())
		{
			rutasAdmin.PUT("/:id", ctrl.UpdateVenta)
			rutasAdmin.DELETE("/:id", ctrl.DeleteVenta)
		}
	}
}