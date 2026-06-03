package inventario

type Categoria struct {
	ID              string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_categoria"`
	NombreCategoria string `json:"nombre_categoria"`
}
type Producto struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_producto"`
	Nombre       string    `json:"nombre"`
	Descripcion  string    `json:"descripcion"`
	Stock        int       `json:"stock"`
	Precio       float64   `json:"precio"`
	Marca        string    `json:"marca"`
	CodigoBarras string    `json:"codigo_barras"`
	CategoriaID  string    `json:"id_categoria"`
	Categoria    Categoria `gorm:"foreignKey:CategoriaID" json:"categoria,omitempty"`
}
