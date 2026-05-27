package ventas

import "gorm.io/gorm"

// Recibe la BD y un puntero a la venta para inyectarle el UUID
func GuardarVenta(db *gorm.DB, venta *Venta) error {
	err := db.Create(venta).Error
	return err
}