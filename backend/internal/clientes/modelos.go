package clientes

import "time"

type Cliente struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Nombre    string     `gorm:"type:varchar(100);not null" json:"nombre" binding:"required"`
	Rut       string     `gorm:"type:varchar(12);unique;not null" json:"rut" binding:"required"`
	Telefono  string     `gorm:"type:varchar(20);not null" json:"telefono" binding:"required"`
	CreatedAt time.Time  `json:"created_at"` //guarda la fecha del primer registro creado
	UpdatedAt time.Time  `json:"updated_at"` //guarda la fecha de la última actualización

	//no son necesarios, pero nos sirven pa despues
}