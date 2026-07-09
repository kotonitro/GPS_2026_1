package promociones

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, authMiddleware gin.HandlerFunc) {

	ctrl := NewPromocionController(db)

	group := api.Group("promociones")
	group.Use(authMiddleware)
	{
		group.GET("/", ctrl.GetPromocionesController)
		group.GET("/:id", ctrl.GetPromocionByIDController)
		group.POST("", ctrl.CreatePromocionController)
		group.DELETE("/:id", ctrl.DeletePromocionController)
		group.PATCH("/:id", ctrl.UpdatePromocionByIDController)
	}
}
