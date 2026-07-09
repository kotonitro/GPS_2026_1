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

	cfg := config.LoadConfig()

	db := database.Connect(cfg)
	database.Migrations(db)

	config.InitialSetup(db)

	validations.ValidationsConfig()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "database connection failed"})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "database ping failed"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontURL, "http://localhost:5173", "http://127.0.0.1:5173", "http://146.83.198.35:1226", "https://146.83.198.35"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
