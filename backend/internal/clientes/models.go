package clientes

import (
	"time"
)

type Cliente struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_cliente"`
	Nombre    string    `gorm:"type:varchar(100);not null" json:"nombre"`
	Rut       string    `gorm:"type:varchar(16);unique;not null" json:"rut"`
	Telefono  string    `gorm:"type:varchar(16)" json:"telefono"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
