package database

import (
	"fmt"
	"log"
	"os"

	"backend/internal/clientes"
	"backend/internal/inventario"
	"backend/internal/promociones"
	"backend/internal/ventas"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init inicializa la conexión a PostgreSQL y ejecuta las migraciones automáticas
func Init() *gorm.DB {

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error crítico al conectar con la base de datos: %v", err)
	}

	err = db.AutoMigrate(
		&clientes.Cliente{},
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
		log.Fatalf("Error crítico al ejecutar las migraciones de tablas: %v", err)
	}

	return db
}
