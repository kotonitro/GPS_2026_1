package empleados

import (
	"time"
)

type Empleado struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_empleado"`
	Rut        string    `gorm:"type:varchar(16);unique;not null" json:"rut"`
	Nombre     string    `gorm:"type:varchar(64);unique;not null" json:"nombre"`
	Usuario    string    `gorm:"type:varchar(16);unique;not null" json:"usuario"`
	Contrasena string    `gorm:"type:varchar(255);not null" json:"-"`
	Telefono   *string   `gorm:"type:varchar(16);unique" json:"telefono"`
	Rol        string    `gorm:"type:varchar(16);not null" json:"rol"`
	Activo     bool      `gorm:"type:bool;default:true" json:"activo"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
