package inventario

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoutesConfig reemplaza a la antigua función ConfigurarRutas
func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, authMiddleware gin.HandlerFunc) {

	// 1. Instanciamos tu nuevo controlador pasándole la base de datos
	ctrl := NewInventarioController(db)

	grupo := api.Group("/inventario")

	// 2. Middleware: Exige que el usuario haya iniciado sesión (Token JWT válido)
	grupo.Use(authMiddleware)
	{
		// Rutas de lectura (Cualquier empleado cajero o admin puede ver productos)
		grupo.GET("/categorias", ctrl.GetCategorias)
		grupo.GET("/productos", ctrl.GetProductos)
		grupo.GET("/productos/:id", ctrl.GetProductoByID)
		grupo.GET("/productos/codigo/:codigo", ctrl.GetProductoByCodigo)

		// 3. Subgrupo de Administrador: Solo los 'Admin' pueden crear, editar o borrar
		adminGroup := grupo.Group("")
		adminGroup.Use(auth.RoleMiddleware("Admin"))
		{
			// CRUD Categorías
			adminGroup.POST("/categorias", ctrl.CrearCategoria)
			adminGroup.PUT("/categorias/:id", ctrl.UpdateCategoria)
			adminGroup.DELETE("/categorias/:id", ctrl.DeleteCategoria)

			// CRUD Productos
			adminGroup.POST("/productos", ctrl.CrearProducto)
			adminGroup.PUT("/productos/:id", ctrl.UpdateProducto)
			adminGroup.DELETE("/productos/:id", ctrl.DeleteProducto)
		}
	}
}
