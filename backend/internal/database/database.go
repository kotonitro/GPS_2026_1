package database

import (
	"backend/internal/cajas"
	"backend/internal/clientes"
	"backend/internal/config"
	"backend/internal/empleados"
	"backend/internal/inventario"
	"backend/internal/promociones"
	"backend/internal/ventas"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Inicializa la conexión a PostgreSQL
func Connect(cfg *config.AppConfig) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error: No se pudo conectar a la base de datos \n", err)
	}

	fmt.Println("Conexión con la base de datos establecida")

	return db
}

// Migra automaticamente las tablas
func Migrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&empleados.Empleado{},
		&clientes.Cliente{},
		&cajas.Caja{},
		&cajas.RegistroTurnos{},
		&inventario.Categoria{},
		&inventario.Producto{},
		&promociones.Promocion{},
		&promociones.DetallePromocion{},
		&ventas.MetodoPago{},
		&ventas.Venta{},
		&ventas.DetalleVenta{},
		&ventas.Fiado{},
	)

	if err != nil {
		log.Fatal("Error fatal al ejecutar las migraciones: ", err)
	}

	log.Println("Migraciones ejecutadas exitosamente")
}
