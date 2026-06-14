package cajas

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB) {

	ctrl := NewCajaController(db)

	grupo := api.Group("cajas")
	{
		grupo.GET("/", ctrl.GetCajasController)
		grupo.GET("/:id", ctrl.GetCajaByIDController)
		grupo.POST("", ctrl.CreateCajaController)
		grupo.DELETE("/:id", ctrl.DeleteCajaByIDController)
		grupo.PATCH("/:id", ctrl.UpdateCajaByIDController)
	}
}
