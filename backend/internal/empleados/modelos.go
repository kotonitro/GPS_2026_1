package empleados

import "time"

type Empleado struct {
	ID         string    `gorm:"primaryKey;type:uuid;default;gen_random_uuid()" json:"id_empleado"`
	Rut        string    `gorm:"type:varchar(12);unique;not null" json:"rut" binding:"required"`
	Usuario    string    `gorm:"type:varchar(16);unique;not null" json:"usuario" binding:"required"`
	Contraseña int       `gorm:"type:varchar(255);not null" json:"-"`
	Telefono   string    `gorm:"type:varchar(20)" json:"telefono"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
