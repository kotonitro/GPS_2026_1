package inventario

import "gorm.io/gorm"

// Recibe la BD y un puntero al producto para inyectarle el UUID
func GuardarProducto(db *gorm.DB, producto *Producto) error {
	err := db.Create(producto).Error
	return err
}
