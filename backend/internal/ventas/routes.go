package ventas

import (
	"backend/internal/auth"
	"github.com/gin-gonic/gin"
)
func ConfigurarRutas(api *gin.RouterGroup, jwtSecret string) {
	
	rutasVentas := api.Group("/ventas")
	
	rutasVentas.Use(auth.AuthMiddleware(jwtSecret))
	{
		rutasVentas.POST("", CrearVenta)
		rutasVentas.GET("", GetVentas)
		rutasVentas.GET("/:id", GetVentaByID)

		rutasAdmin := rutasVentas.Group("")
		rutasAdmin.Use(auth.RoleMiddleware("Admin"))
		{
			rutasAdmin.PUT("/:id", UpdateVenta)
			rutasAdmin.DELETE("/:id", DeleteVenta)
		}
	}
}