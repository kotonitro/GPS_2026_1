package inventario

import "github.com/gin-gonic/gin"

// ConfigurarRutas registra todos los endpoints del módulo de inventario
func ConfigurarRutas(api *gin.RouterGroup) {
	// Creamos un subgrupo específico para inventario
	grupo := api.Group("/inventario")
	{
		// Si la ruta base es /api, esto equivale a POST /api/inventario/
		// Rutas para el CRUD de Productos
		grupo.POST("/productos", CrearProducto)
		grupo.GET("/productos", GetProductos)
		grupo.GET("/productos/:id", GetProductoByID)
		grupo.PUT("/productos/:id", UpdateProducto)
		grupo.DELETE("/productos/:id", DeleteProducto)
	}
}
