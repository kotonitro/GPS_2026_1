package promociones

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB) {

	ctrl := NewPromocionController(db)

	grupo := api.Group("promociones")
	{
		grupo.GET("/", ctrl.GetPromocionesController)
		grupo.GET("/:id", ctrl.GetPromocionByIDController)
		grupo.POST("", ctrl.CreatePromocionController)
		grupo.DELETE("/:id", ctrl.DeletePromocionController)
		grupo.PATCH("/:id", ctrl.UpdatePromocionByIDController)
	}
}