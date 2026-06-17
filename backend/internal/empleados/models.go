package empleados

import (
	"time"
)

type Empleado struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_empleado"`
	Rut        string    `gorm:"type:varchar(16);unique;not null" json:"rut"`
	Usuario    string    `gorm:"type:varchar(16);unique;not null" json:"usuario"`
	Contrasena string    `gorm:"type:varchar(255);not null" json:"-"`
	Telefono   *string   `gorm:"type:varchar(16);unique" json:"telefono"`
	Rol        string    `gorm:"type:varchar(16);not null" json:"rol"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
