package main

import (
	"backend/internal/clientes"
	"backend/internal/database"
	"backend/internal/inventario"
	"backend/internal/ventas"

	"github.com/gin-gonic/gin"
)

func main() {

	db := database.Init()

	// Crea el enrutador de Gin con la configuración por defecto
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Rutas de la api
	api := r.Group("/api")

	inventario.ConfigurarRutas(api)

	rutasClientes := r.Group("/clientes")
	{
		rutasClientes.POST("", clientes.CreateCliente)
		rutasClientes.GET("", clientes.GetClientes)
		rutasClientes.GET("/:id", clientes.GetClienteByID)
		rutasClientes.PUT("/:id", clientes.UpdateCliente)
		rutasClientes.DELETE("/:id", clientes.DeleteCliente)
	}

	rutasVentas := r.Group("/ventas")
	{
		rutasVentas.POST("", ventas.CrearVenta)
	}

	// Encendemos el servidor en el puerto 8080
	r.Run(":8080")
}
