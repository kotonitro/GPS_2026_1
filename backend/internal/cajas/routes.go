package cajas

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, jwtSecret string) {

	ctrl := NewCajaController(db)

	authGroup := api.Group("cajas")
	authGroup.Use(auth.AuthMiddleware(jwtSecret))
	{
		authGroup.GET("/", ctrl.GetCajasController)
		authGroup.GET("/:id", ctrl.GetCajaByIDController)

		adminGroup := authGroup.Group("")
		adminGroup.Use(auth.RoleMiddleware("Admin"))
		{
			adminGroup.POST("", ctrl.CreateCajaController)
			adminGroup.DELETE("/:id", ctrl.DeleteCajaByIDController)
			adminGroup.PATCH("/:id", ctrl.UpdateCajaByIDController)
		}
	}
}
