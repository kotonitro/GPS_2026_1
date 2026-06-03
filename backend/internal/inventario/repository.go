package inventario

import "gorm.io/gorm"

// Recibe la BD y un puntero al producto para inyectarle el UUID
func GuardarProducto(db *gorm.DB, producto *Producto) error {
	err := db.Create(producto).Error
	return err
}

// ObtenerTodosProductos recupera todos los productos de la base de datos, incluyendo su categoría.
func ObtenerTodosProductos(db *gorm.DB) ([]Producto, error) {
	var productos []Producto
	// Preload("Categoria") carga la información de la categoría asociada a cada producto.
	if err := db.Preload("Categoria").Find(&productos).Error; err != nil {
		return nil, err
	}
	return productos, nil
}

// ObtenerProductoPorID busca un producto por su ID, incluyendo su categoría.
func ObtenerProductoPorID(db *gorm.DB, id string) (*Producto, error) {
	var producto Producto
	if err := db.Preload("Categoria").First(&producto, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &producto, nil
}

// ActualizarProducto actualiza los datos de un producto en la base de datos.
func ActualizarProducto(db *gorm.DB, productoExistente *Producto, datosNuevos *Producto) error {
	return db.Model(productoExistente).Updates(datosNuevos).Error
}

// EliminarProducto elimina un producto de la base de datos por su ID.
func EliminarProducto(db *gorm.DB, id string) error {
	return db.Delete(&Producto{}, "id = ?", id).Error
}
