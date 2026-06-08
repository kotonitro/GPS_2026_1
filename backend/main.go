package main

import (
	"backend/internal/clientes"
	"backend/internal/database"
	"backend/internal/empleados"
	"backend/internal/inventario"
	"backend/internal/validations"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {

	db := database.Init()
	database.Migrations(db)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("contrasena_segura", validations.ValidarContrasena)
		v.RegisterValidation("rut_valido", validations.ValidarRUT)
	}

	r := gin.Default()
	api := r.Group("/api")

	inventario.ConfigurarRutas(api)
	clientes.ConfigurarRutas(api)
	empleados.RoutesConfig(api, db)

	// Encendemos el servidor en el puerto 8080
	r.Run(":8080")
}
