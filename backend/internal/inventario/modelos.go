package inventario

// Modelo de Categoría
type Categoria struct {
	ID              string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_categoria"`
	NombreCategoria string `json:"nombre_categoria"`
}

// Modelo de Producto
type Producto struct {
	ID           string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_producto"`
	Nombre       string  `json:"nombre"`
	Descripcion  string  `json:"descripcion"`
	Stock        int     `json:"stock"`
	Precio       float64 `json:"precio"`
	Marca        string  `json:"marca"`
	CodigoBarras string  `json:"codigo_barras"`

	// La Llave Foránea
	CategoriaID string `json:"id_categoria"`

	// GORM usa la palabra reservada `foreignKey` para entender la relación
	Categoria Categoria `gorm:"foreignKey:CategoriaID" json:"categoria,omitempty"`
}
