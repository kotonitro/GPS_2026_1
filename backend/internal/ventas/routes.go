package ventas

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, authMiddleware gin.HandlerFunc) {
	
	ctrl := NewVentasController(db)

	group := api.Group("ventas")
	
	group.Use(authMiddleware) 
	{
		group.POST("/", ctrl.CrearVenta)
		group.GET("/", ctrl.GetVentas)
		group.GET("/:id", ctrl.GetVentaByID)
		group.PUT("/:id", ctrl.UpdateVenta)
		group.DELETE("/:id", ctrl.DeleteVenta)
	}
}