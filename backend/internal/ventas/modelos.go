package ventas
import "time"

type Venta struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_venta"`
	CajaID         string    `json:"id_caja"`
	MetodoID       string    `json:"id_metodo"`
	EmpleadoID     string    `json:"id_empleado"`
	FechaEmision   time.Time `json:"fecha_emision"`
	Pago           float64   `json:"pago"`
	Vuelto         float64   `json:"vuelto"`
	Precio         float64   `json:"precio"`
	MontoDescuento float64   `json:"monto_descuento"`
	EstadoSync     string    `json:"estado_sync"`

	Detalles   []DetalleVenta `gorm:"foreignKey:VentaID" json:"detalles,omitempty"`
	MetodoPago MetodoPago     `gorm:"foreignKey:MetodoID" json:"metodo_pago,omitempty"`
}


type DetalleVenta struct {
		ID         string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_detalle"`
		VentaID    string  `json:"id_venta"`
		ProductoID string  `json:"id_producto"`
		Cantidad   int     `json:"cantidad"`
		MontoFinal float64 `json:"monto_final"`
}

type MetodoPago struct {
		ID           string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_metodo"`
		NombreMetodo string `json:"nombre_metodo"`
}

type Fiado struct {
		ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_fiado"`
		ClienteID   string    `json:"id_cliente"`
		VentaID     string    `json:"id_venta"`
		FechaInicio time.Time `json:"fecha_inicio"`
		FechaLimite time.Time `json:"fecha_limite"`
}
