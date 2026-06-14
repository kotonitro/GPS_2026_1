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
			group.POST("", ctrl.CreateEmpleadoController)
			group.DELETE("/:id", ctrl.DeleteEmpleadoByIDController)
			group.PATCH("/:id", ctrl.UpdateEmpleadoByIDController)
		}
	}
}
