package empleados

import "github.com/gin-gonic/gin"

func ConfigurarRutas(api *gin.RouterGroup) {

	grupo := api.Group("empleados")
	{
		grupo.POST("", CrearEmpleado)
		// grupo.GET("/", ListarEmpleados)
		// grupo.GET("/:id", ObtenerEmpleado)
		// grupo.PUT("/:id", ActualizarEmpleado)
	}
}
