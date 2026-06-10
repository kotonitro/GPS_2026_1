package inventario

import "github.com/gin-gonic/gin"

// ConfigurarRutas registra todos los endpoints del módulo de inventario
func ConfigurarRutas(api *gin.RouterGroup) {
	grupo := api.Group("/inventario")
	{
		// Rutas para categorías
		grupo.POST("/categorias", CrearCategoria)
		grupo.GET("/categorias", GetCategorias)
		grupo.PUT("/categorias/:id", UpdateCategoria)
		grupo.DELETE("/categorias/:id", DeleteCategoria)

		// Rutas para el CRUD de Productos
		grupo.POST("/productos", CrearProducto)
		grupo.GET("/productos", GetProductos)
		grupo.GET("/productos/:id", GetProductoByID)
		grupo.PUT("/productos/:id", UpdateProducto)
		grupo.DELETE("/productos/:id", DeleteProducto)
		// obtener producto por código de barras
		grupo.GET("/productos/codigo/:codigo", GetProductoByCodigo)

	}
}
