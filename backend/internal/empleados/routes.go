package empleados

import "github.com/gin-gonic/gin"

func ConfigurarRutas(api *gin.RouterGroup) {

	grupo := api.Group("empleados")
	{
		grupo.POST("", CrearEmpleado)
		// grupoEmpleados.GET("/", ListarEmpleados)
		// grupoEmpleados.GET("/:id", ObtenerEmpleado)
		// grupoEmpleados.PUT("/:id", ActualizarEmpleado)
	}
}
