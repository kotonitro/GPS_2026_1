package empleados

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, jwtSecret string) {

	ctrl := NewEmpleadoController(db)

	group := api.Group("empleados")
	group.Use(auth.AuthMiddleware(jwtSecret))
	{
		group.GET("/", ctrl.GetEmpleadosController)
		group.GET("/:id", ctrl.GetEmpleadoByIDController)

		adminGroup := group.Group("")
		adminGroup.Use(auth.RoleMiddleware("Admin"))
		{
			adminGroup.POST("", ctrl.CreateEmpleadoController)
			adminGroup.DELETE("/:id", ctrl.DeleteEmpleadoByIDController)
			adminGroup.PATCH("/:id", ctrl.UpdateEmpleadoByIDController)
		}
	}
}
