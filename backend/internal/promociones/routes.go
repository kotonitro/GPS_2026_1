package promociones

import (
	"backend/internal/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, jwtSecret string) {

	ctrl := NewPromocionController(db)

	group := api.Group("promociones")
	group.Use(auth.AuthMiddleware(jwtSecret))
	{
		group.GET("/", ctrl.GetPromocionesController)
		group.GET("/:id", ctrl.GetPromocionByIDController)
		group.POST("", ctrl.CreatePromocionController)
		group.DELETE("/:id", ctrl.DeletePromocionController)
		group.PATCH("/:id", ctrl.UpdatePromocionByIDController)
	}
}