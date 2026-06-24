package main

import (
	"backend/internal/auth"
	"backend/internal/cajas"
	"backend/internal/clientes"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/empleados"
	"backend/internal/inventario"
	"backend/internal/promociones"
	"backend/internal/validations"
	"backend/internal/ventas"
	"time"

	"github.com/gin-contrib/cors"
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

	//CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// api
	api := r.Group("/api")

	// rutas
	auth.RoutesConfig(api, db, cfg.JWTSecret)
	cajas.RoutesConfig(api, db, cfg.JWTSecret)
	clientes.RoutesConfig(api, db, cfg.JWTSecret)
	empleados.RoutesConfig(api, db, cfg.JWTSecret)
	inventario.RoutesConfig(api, db, cfg.JWTSecret)
	ventas.RoutesConfig(api, db, cfg.JWTSecret)
	promociones.RoutesConfig(api, db, cfg.JWTSecret)

	r.Run(":8080")
}
