package ventas

import "gorm.io/gorm"

func GuardarVenta(db *gorm.DB, venta *Venta) error {
	err := db.Create(venta).Error
	return err
}

func ObtenerTodasVentas(db *gorm.DB) ([]Venta, error) {
    var ventas []Venta
    if err := db.Preload("MetodoPago").Preload("Detalles").Find(&ventas).Error; err != nil {
        return nil, err
    }
    return ventas, nil
}

func ObtenerVentaPorID(db *gorm.DB, id string) (*Venta, error) {
    var venta Venta
    if err := db.Preload("MetodoPago").Preload("Detalles").First(&venta, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &venta, nil
}

func ActualizarVenta(db *gorm.DB, ventaExistente *Venta, datosNuevos *Venta) error {
	return db.Model(ventaExistente).Updates(datosNuevos).Error
}

func EliminarVenta(db *gorm.DB, id string) error {
	return db.Delete(&Venta{}, "id = ?", id).Error
}
