package promociones
import ("backend/internal/inventario")

type Promocion struct {
	ID 		  	string `gorm:"primaryKey;type:uuid;default;gen_random_uuid()" json:"id_promocion"`
	Tipo        string `json:"tipo"` //dsp este va a ser NXM,pocentaje o precio fijo 
	Lleva 	    int `json:"lleva"` //N
	Paga        int `json:"paga"`  //M
	Descuento 	float64 `json:"descuento"`
}

type DetallePromocion struct{
	PromocionID string `json:"id_promocion"`
	Promocion 	Promocion`gorm:"foreignKey:PromocionID" json:"promocion,omitzero"`
	EstaActivo 	bool `json:"esta_activa"`
	Producto 	inventario.Producto `gorm:"foreignKey:ProductoID" json:"producto,omitzero"`
}
