package ventas

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, jwtSecret string) {
	
	ctrl := NewVentasController(db)

	group := api.Group("ventas")
	
	group.Use(auth.AuthMiddleware(jwtSecret)) 
	{
		group.POST("/", ctrl.CrearVenta)
		group.GET("/", ctrl.GetVentas)
		group.GET("/:id", ctrl.GetVentaByID)
		group.PUT("/:id", ctrl.UpdateVenta)
		group.DELETE("/:id", ctrl.DeleteVenta)
	}
}