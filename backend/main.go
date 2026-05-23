package main

import (
	"backend/internal/clientes"
	"backend/internal/database"
	"backend/internal/inventario"

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

	rutasClientes := r.Group("/clientes")
	{
		rutasClientes.POST("", clientes.CreateCliente)
		rutasClientes.GET("", clientes.GetClientes)
		rutasClientes.GET("/:id", clientes.GetClienteByID)
		rutasClientes.PUT("/:id", clientes.UpdateCliente)
		rutasClientes.DELETE("/:id", clientes.DeleteCliente)
	}

	rutasInventario := r.Group("/inventario")
	{
		rutasInventario.POST("/productos", inventario.CrearProducto)
	}

	// Encendemos el servidor en el puerto 8080
	r.Run(":8080")
}
