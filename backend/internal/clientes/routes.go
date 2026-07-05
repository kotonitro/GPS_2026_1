package clientes

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, authMiddleware gin.HandlerFunc) {
	ctrl := NewClienteController(db)

	group := api.Group("clientes")
	group.Use(authMiddleware)
	{
		group.GET("/", ctrl.GetClientesController)
		group.GET("/:id", ctrl.GetClienteByIDController)
		group.GET("/rut/:rut", ctrl.GetClienteByRutController)
		group.GET("/search/nombre", ctrl.GetClientesByNombreController)

		adminGroup := group.Group("")
		adminGroup.Use(auth.RoleMiddleware("Admin"))
		{
			adminGroup.POST("", ctrl.CreateClienteController)
			adminGroup.PATCH("/:id", ctrl.UpdateClienteByIDController)
			adminGroup.DELETE("/:id", ctrl.DeleteClienteByIDController)
		}
	}
}
