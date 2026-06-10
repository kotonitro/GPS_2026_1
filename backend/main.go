package main

import (
	"backend/internal/clientes"
	"backend/internal/database"
	"backend/internal/empleados"
	"backend/internal/inventario"
	"backend/internal/validations"
	"backend/internal/ventas"
	"github.com/gin-gonic/gin"
)

func main() {

	// base de datos
	db := database.Init()
	database.Migrations(db)

	// validaciones
	validations.ValidationsConfig()

	// rutas
	r := gin.Default()
	api := r.Group("/api")

	inventario.ConfigurarRutas(api)
	ventas.ConfigurarRutas(api)
	clientes.ConfigurarRutas(api)
	empleados.RoutesConfig(api, db)

	r.Run(":8080")
}
