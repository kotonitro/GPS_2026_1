package empleados

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, authMiddleware gin.HandlerFunc) {

	ctrl := NewEmpleadoController(db)

	authGroup := api.Group("empleados")
	authGroup.Use(authMiddleware)
	{
		authGroup.GET("/", ctrl.GetEmpleadosController)
		authGroup.GET("/:id", ctrl.GetEmpleadoByIDController)

		adminGroup := authGroup.Group("")
		adminGroup.Use(auth.RoleMiddleware("Admin"))
		{
			adminGroup.POST("", ctrl.CreateEmpleadoController)
			adminGroup.DELETE("/:id", ctrl.DeleteEmpleadoByIDController)
			adminGroup.PATCH("/:id", ctrl.UpdateEmpleadoByIDController)
		}
	}
}
