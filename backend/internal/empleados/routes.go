package empleados

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB) {

	ctrl := NewEmpleadoController(db)

	grupo := api.Group("empleados")
	{
		grupo.GET("/", ctrl.GetEmpleadosController)
		grupo.GET("/:id", ctrl.GetEmpleadoByIDController)
		grupo.POST("", ctrl.CreateEmpleadoController)
	}
}
