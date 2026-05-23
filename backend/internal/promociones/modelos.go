package promociones
//import ("GPS_2026_1\backend\internal\inventario")

type Promocion struct {
	ID 		  	string `gorm:"primaryKey;type:uuid;default;gen_random_uuid()" json:"id_promocion"`
	Mayorista 	int `json:"mayorista"`
	Descuento 	float64 `json:"descuento"`
}

type DetalleVenta struct{
	PromocionID string `json:"id_promocion"`
	Promocion 	Promocion`gorm:"foreignKey:PromocionID" json:"promocion,omitzero"`
	EstaActivo 	bool `json:"esta_activa"`
//	Producto 	inventario.Producto `gorm:"foreignKey:ProductoID" json:"producto,omitzero"`
}
