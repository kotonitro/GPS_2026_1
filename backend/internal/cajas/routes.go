package cajas

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, authMiddleware gin.HandlerFunc) {

	ctrl := NewCajaController(db)

	authGroup := api.Group("cajas")
	authGroup.Use(authMiddleware)
	{
		authGroup.GET("", ctrl.GetCajasController)
		authGroup.GET("/:id", ctrl.GetCajaByIDController)

		adminGroup := authGroup.Group("")
		adminGroup.Use(auth.AdminMiddleware())
		{
			adminGroup.POST("", ctrl.CreateCajaController)
			adminGroup.DELETE("/:id", ctrl.DeleteCajaByIDController)
			adminGroup.PATCH("/:id", ctrl.UpdateCajaByIDController)
			adminGroup.GET("/registros", ctrl.GetRegistrosController)
			adminGroup.GET("/registros/:id", ctrl.GetRegistroByIDController)
			adminGroup.POST("/registros", ctrl.CreateRegistroController)
			adminGroup.DELETE("/registros/:id", ctrl.DeleteRegistroByIDController)
			adminGroup.PATCH("/registros/:id", ctrl.UpdateRegistroByIDController)
		}
	}
}
