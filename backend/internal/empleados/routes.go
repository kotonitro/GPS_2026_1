package empleados

import "github.com/gin-gonic/gin"

func RoutesConfig(api *gin.RouterGroup) {

	grupo := api.Group("empleados")
	{
		grupo.GET("/", GetEmpleadosController)
		grupo.GET("/:id", GetEmpleadoByIDController)
		grupo.POST("", CreateEmpleadoController)
		// grupo.PUT("/:id", ActualizarEmpleado)
	}
}
