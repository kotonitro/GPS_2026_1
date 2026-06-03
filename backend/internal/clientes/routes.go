package clientes

import "github.com/gin-gonic/gin"

// ConfigurarRutas registra todos los endpoints del módulo de clientes
func ConfigurarRutas(api *gin.RouterGroup) {
	grupo := api.Group("/clientes")
	{
		grupo.POST("", CreateCliente)
		grupo.GET("", GetClientes)
		grupo.GET("/:id", GetClienteByID)
		grupo.GET("/rut/:rut", SearchClienteByRut)
		grupo.GET("/search/nombre", SearchClienteByNombre)
		grupo.PUT("/:id", UpdateCliente)
		grupo.DELETE("/:id", DeleteCliente)
	}
}
