package inventario

import "github.com/gin-gonic/gin"

// ConfigurarRutas registra todos los endpoints del módulo de inventario
func ConfigurarRutas(api *gin.RouterGroup) {
	grupo := api.Group("/inventario")
	{
		// Rutas para el CRUD de Productos
		grupo.POST("/productos", CrearProducto)
		grupo.GET("/productos", GetProductos)
		grupo.GET("/productos/:id", GetProductoByID)
		grupo.PUT("/productos/:id", UpdateProducto)
		grupo.DELETE("/productos/:id", DeleteProducto)
	}
}
