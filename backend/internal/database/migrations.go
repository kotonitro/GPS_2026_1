package database

import (
	"backend/internal/cajas"
	"backend/internal/clientes"
	"backend/internal/empleados"
	"backend/internal/inventario"
	"backend/internal/promociones"
	"backend/internal/ventas"
	"log"

	"gorm.io/gorm"
)

// Sincroniza los modelos con las tablas en la base de datos
func Migrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&empleados.Empleado{},
		&clientes.Cliente{},
		&cajas.Caja{},
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
