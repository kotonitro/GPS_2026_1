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
		authGroup.GET("", ctrl.GetEmpleadosController)
		authGroup.GET("/:id", ctrl.GetEmpleadoByIDController)
		authGroup.GET("/roles", ctrl.GetRolesController)
		authGroup.GET("/roles/:id", ctrl.GetRolByIDController)

		adminGroup := authGroup.Group("")
		adminGroup.Use(auth.AdminMiddleware())
		{
			adminGroup.POST("", ctrl.CreateEmpleadoController)
			adminGroup.DELETE("/:id", ctrl.DeleteEmpleadoByIDController)
			adminGroup.PATCH("/:id", ctrl.UpdateEmpleadoByIDController)
			adminGroup.POST("/roles", ctrl.CreateRolController)
			adminGroup.DELETE("/roles/:id", ctrl.DeleteRolByIDController)
			adminGroup.PATCH("/roles/:id", ctrl.UpdateRolByIDController)
		}
	}
}
