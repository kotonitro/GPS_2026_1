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
		AllowOrigins:     []string{cfg.FrontURL, "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// api
	api := r.Group("/api")

	// auth middleware
	authCtrl := auth.NewAuthController(db, cfg.JWTSecret, cfg.CookieDomain)
	authMiddleware := authCtrl.AuthMiddleware()

	// rutas
	auth.RoutesConfig(api, authCtrl, authMiddleware)
	cajas.RoutesConfig(api, db, authMiddleware)
	clientes.RoutesConfig(api, db, authMiddleware)
	empleados.RoutesConfig(api, db, authMiddleware)
	inventario.RoutesConfig(api, db, authMiddleware)
	ventas.ConfigurarRutas(api, db, authMiddleware)
	promociones.RoutesConfig(api, db, authMiddleware)
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
	
	rutasVentas := r.Group("/ventas")
	{
		rutasVentas.POST("", ventas.CrearVenta)
	}

	r.Run(":8080")
}
