package ventas

import "gorm.io/gorm"

// Recibe la BD y un puntero a la venta para inyectarle el UUID
func GuardarVenta(db *gorm.DB, venta *Venta) error {
	err := db.Create(venta).Error
	return err
}

// ObtenerTodasVentas recupera todo el historial de ventas
func ObtenerTodasVentas(db *gorm.DB) ([]Venta, error) {
	var ventas []Venta
	// Si a futuro quieres traer el detalle de los productos vendidos, puedes agregar .Preload("DetalleVentas")
	if err := db.Find(&ventas).Error; err != nil {
		return nil, err
	}
	return ventas, nil
}

// ObtenerVentaPorID busca una venta específica por su ID
func ObtenerVentaPorID(db *gorm.DB, id string) (*Venta, error) {
	var venta Venta
	if err := db.First(&venta, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &venta, nil
}

// ActualizarVenta actualiza los datos de una transacción (copiando el estilo de inventario)
func ActualizarVenta(db *gorm.DB, ventaExistente *Venta, datosNuevos *Venta) error {
	return db.Model(ventaExistente).Updates(datosNuevos).Error
}

// EliminarVenta elimina una transacción de la base de datos por su ID
func EliminarVenta(db *gorm.DB, id string) error {
	return db.Delete(&Venta{}, "id = ?", id).Error
}