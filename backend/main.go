package main

import (
	"backend/internal/auth"
	"backend/internal/cajas"
	"backend/internal/clientes"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/empleados"
	"backend/internal/inventario"
	"backend/internal/validations"
	"backend/internal/ventas"

	"github.com/gin-gonic/gin"
)

func main() {

	// variables de entorno
	cfg := config.LoadConfig()

	// base de datos
	db := database.Connect(cfg)
	database.Migrations(db)

	// setup inicial
	config.InitialSetup(db)

	// validaciones
	validations.ValidationsConfig()

	// router
	r := gin.Default()

	// api
	api := r.Group("/api")

	// rutas
	auth.RoutesConfig(api, db, cfg.JWTSecret)
	cajas.RoutesConfig(api, db, cfg.JWTSecret)
	clientes.ConfigurarRutas(api)
	empleados.RoutesConfig(api, db, cfg.JWTSecret)
	inventario.ConfigurarRutas(api)
	ventas.ConfigurarRutas(api)

	r.Run(":8080")
}
