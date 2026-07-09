package clientes

import (
	"time"
)

type Cliente struct {
	ID           string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_cliente"`
	Nombre       string     `gorm:"type:varchar(100);not null" json:"nombre"`
	Rut          string     `gorm:"type:varchar(16);unique;not null" json:"rut"`
	Telefono     string     `gorm:"type:varchar(16);unique;not null" json:"telefono"`
	FiadoActual  float64    `gorm:"type:decimal(10,2);default:0.0" json:"fiado_actual"`
	FiadoMaximo  float64    `gorm:"type:decimal(10,2);default:20000.0" json:"fiado_maximo"`
	UltimaCompra *time.Time `json:"ultima_compra"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
