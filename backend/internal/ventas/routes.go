package ventas

import "github.com/gin-gonic/gin"

// ConfigurarRutas registra todos los endpoints del módulo de ventas
func ConfigurarRutas(api *gin.RouterGroup) {
	rutasVentas := api.Group("/ventas")
	{
		// Rutas para el CRUD de ventas
		rutasVentas.POST("", CrearVenta)
		rutasVentas.GET("", GetVentas)
		rutasVentas.GET("/:id", GetVentaByID)
		rutasVentas.PUT("/:id", UpdateVenta)
		rutasVentas.DELETE("/:id", DeleteVenta)
	}
}

