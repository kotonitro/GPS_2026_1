package empleados

import (
	"time"
)

type Empleado struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_empleado"`
	Rut        string `gorm:"type:varchar(16);unique;not null" json:"rut"`
	Nombre     string `gorm:"type:varchar(64);not null" json:"nombre"`
	Usuario    string `gorm:"type:varchar(16);unique;not null" json:"usuario"`
	Contrasena string `gorm:"type:varchar(255);not null" json:"-"`
	Telefono   string `gorm:"type:varchar(16);unique;not null" json:"telefono"`
	Activo     bool   `gorm:"type:bool;default:true" json:"activo"`

	RolID string `gorm:"type:uuid;not null" json:"id_rol"`
	Rol   *Rol   `gorm:"foreignKey:RolID;references:ID" json:"rol,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Rol struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_rol"`
	Nombre      string    `gorm:"type:varchar(16);not null;unique" json:"nombre"`
	Descripcion string    `gorm:"type:varchar(255);not null" json:"descripcion"`
	EsAdmin     bool      `gorm:"default:false;not null" json:"es_admin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Rol) TableName() string {
	return "roles"
}
