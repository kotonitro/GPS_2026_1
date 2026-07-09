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

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	db := database.Connect(cfg)
	database.Migrations(db)

	config.InitialSetup(db)

	validations.ValidationsConfig()

	r := gin.Default()

	config := cors.DefaultConfig()
	// 1. Origen exacto de tu frontend (sin asteriscos)
	config.AllowOrigins = []string{"http://146.83.198.35:1226"}
	// 2. Permitir envío de credenciales (soluciona tu error rojo en consola)
	config.AllowCredentials = true
	// 3. Métodos permitidos
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	// 4. Headers permitidos (asegúrate de incluir los que envíe tu frontend)
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(config))
	api := r.Group("/api")

	authCtrl := auth.NewAuthController(db, cfg.JWTSecret, cfg.CookieDomain)
	authMiddleware := authCtrl.AuthMiddleware()

	auth.RoutesConfig(api, authCtrl, authMiddleware)
	cajas.RoutesConfig(api, db, authMiddleware)
	clientes.RoutesConfig(api, db, authMiddleware)
	empleados.RoutesConfig(api, db, authMiddleware)
	inventario.RoutesConfig(api, db, authMiddleware)
	ventas.ConfigurarRutas(api, db, authMiddleware)
	promociones.RoutesConfig(api, db, authMiddleware)

	r.Run(":8080")
}
