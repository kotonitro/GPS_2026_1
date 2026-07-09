package inventario

import (
	"errors"
	"gorm.io/gorm"
)

// ---- PRODUCTOS ----
// Recibe la BD y un puntero al producto para inyectarle el UUID
func GuardarProducto(db *gorm.DB, producto *Producto) error {
	err := db.Create(producto).Error
	return err
}

// ver todos los productos ACTIVOS de la base de datos.
func ObtenerTodosProductos(db *gorm.DB) ([]Producto, error) {
	var productos []Producto
	// Ocultamos los productos descontinuados
	if err := db.Preload("Categoria").Where("estado = ?", true).Find(&productos).Error; err != nil {
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

func ObtenerProductoPorCodigo(db *gorm.DB, codigo string) (*Producto, error) {
	var producto Producto
	// Usamos Preload para que también nos traiga los datos de su categoría
	if err := db.Preload("Categoria").First(&producto, "codigo_barras = ?", codigo).Error; err != nil {
		return nil, err
	}
	return &producto, nil
}

// ActualizarProducto actualiza los datos de un producto en la base de datos.
func ActualizarProducto(db *gorm.DB, productoExistente *Producto, datosNuevos *Producto) error {
	return db.Model(productoExistente).Updates(map[string]interface{}{
		"nombre":        datosNuevos.Nombre,
		"descripcion":   datosNuevos.Descripcion,
		"stock":         datosNuevos.Stock,
		"stock_minimo":  datosNuevos.StockMinimo,
		"precio":        datosNuevos.Precio,
		"unidad":        datosNuevos.Unidad,
		"marca":         datosNuevos.Marca,
		"codigo_barras": datosNuevos.CodigoBarras,
		"categoria_id":  datosNuevos.CategoriaID,
	}).Error
}

// Desactivamos un producto en la base de datos, marcándolo como descontinuado
func EliminarProducto(db *gorm.DB, id string) error {
	return db.Model(&Producto{}).Where("id = ?", id).Update("estado", false).Error
}

// ---- CATEGORIAS ----

func GuardarCategoria(db *gorm.DB, categoria *Categoria) error {
	return db.Create(categoria).Error
}

// Obtener todas las categorías de la base de datos.
func ObtenerTodasCategorias(db *gorm.DB) ([]Categoria, error) {
	var categorias []Categoria
	if err := db.Find(&categorias).Error; err != nil {
		return nil, err
	}
	return categorias, nil
}

// modifica el nombre de una categoría existente
func ActualizarCategoria(db *gorm.DB, id string, datosNuevos *Categoria) error {
	return db.Model(&Categoria{}).Where("id = ?", id).Updates(datosNuevos).Error
}

// borra una categoría de la base de datos
func EliminarCategoria(db *gorm.DB, id string) error {
	// 1. Verificar si hay productos activos usando esta categoría
	var count int64
	db.Model(&Producto{}).Where("categoria_id = ? AND estado = ?", id, true).Count(&count)
	if count > 0 {
		return errors.New("hay productos activos")
	}

	// 2. Eliminar físicamente los productos descontinuados (estado = false)
	db.Unscoped().Where("categoria_id = ? AND estado = ?", id, false).Delete(&Producto{})

	// 3. Eliminar la categoría
	return db.Where("id = ?", id).Delete(&Categoria{}).Error
}
