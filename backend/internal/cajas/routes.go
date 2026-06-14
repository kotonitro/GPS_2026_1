package cajas

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, jwtSecret string) {

	ctrl := NewCajaController(db)

	group := api.Group("cajas")
	group.Use(auth.AuthMiddleware(jwtSecret))
	{
		group.GET("/", ctrl.GetCajasController)
		group.GET("/:id", ctrl.GetCajaByIDController)

		adminGroup := group.Group("")
		adminGroup.Use(auth.RoleMiddleware("Admin"))
		{

			group.POST("", ctrl.CreateCajaController)
			group.DELETE("/:id", ctrl.DeleteCajaByIDController)
			group.PATCH("/:id", ctrl.UpdateCajaByIDController)

		}
	}
}
