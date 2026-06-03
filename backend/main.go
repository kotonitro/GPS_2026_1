package main

import (
	"backend/internal/clientes"
	"backend/internal/database"
	"backend/internal/empleados"
	"backend/internal/inventario"
	"backend/internal/validations"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {

	database.Init()

	err := database.DB.AutoMigrate(&empleados.Empleado{})
	if err != nil {
		log.Fatal("Error al ejecutar la migración de la base de datos:", err)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("contrasena_segura", validations.ValidarContrasena)
		v.RegisterValidation("rut_valido", validations.ValidarRUT)
	}

	r := gin.Default()

	// Rutas de la api
	api := r.Group("/api")

	inventario.ConfigurarRutas(api)
	clientes.ConfigurarRutas(api)

	// Encendemos el servidor en el puerto 8080
	r.Run(":8080")
}
