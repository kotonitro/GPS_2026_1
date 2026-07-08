package ventas

import (
	"backend/internal/cajas"
	"backend/internal/clientes"
	"time"
)

type Venta struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_venta"`
	CajaID         string    `gorm:"type:uuid;not null" json:"id_caja"`
	MetodoID       string    `json:"id_metodo"`
	EmpleadoID     string    `json:"id_empleado"`
	ClienteID      *string   `gorm:"type:uuid" json:"id_cliente,omitempty"`
	FechaEmision   time.Time `json:"fecha_emision"`
	Pago           float64   `json:"pago"`
	Vuelto         float64   `json:"vuelto"`
	MontoTotal     float64   `json:"monto_total"`
	MontoDescuento float64   `json:"monto_descuento"`
	EstadoSync     string    `json:"estado_sync"`

	Detalles   []DetalleVenta `gorm:"foreignKey:VentaID" json:"detalles,omitempty"`
	MetodoPago MetodoPago     `gorm:"foreignKey:MetodoID" json:"metodo_pago,omitempty"`
	Fiado      *Fiado         `gorm:"foreignKey:VentaID" json:"fiado,omitempty"`
	Caja       *cajas.Caja    `gorm:"foreignKey:CajaID;references:ID;constraint:OnDelete:RESTRICT" json:"caja,omitempty"`
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
	ID          string            `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_fiado"`
	ClienteID   string            `gorm:"type:uuid" json:"id_cliente"`
	VentaID     string            `gorm:"type:uuid" json:"id_venta"`
	FechaInicio time.Time         `json:"fecha_inicio"`
	FechaLimite time.Time         `json:"fecha_limite"`
	MontoTotal  float64           `gorm:"type:decimal(10,2);default:0.0" json:"monto_total"`
	Pagado      bool              `gorm:"default:false" json:"pagado"`
	Cliente     *clientes.Cliente `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
}
