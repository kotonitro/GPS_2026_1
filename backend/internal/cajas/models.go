package cajas

import (
	"backend/internal/empleados"
	"time"
)

type Caja struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_caja"`
	Nombre       string    `gorm:"type:varchar(64);unique;not null" json:"nombre"`
	Ubicacion    string    `gorm:"type:varchar(64);unique;not null" json:"ubicacion"`
	Activo       bool      `gorm:"type:bool;default:true" json:"activo"`
	SaldoInicial uint      `gorm:"type:integer;check:saldo_inicial >= 0;default:0" json:"saldo_inicial"`
	SaldoFinal   uint      `gorm:"type:integer;check:saldo_final >= 0;default:0" json:"saldo_final"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegistroTurno struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_registro"`
	CajaID     string `gorm:"type:uuid;not null" json:"id_caja"`
	EmpleadoID string `gorm:"type:uuid;not null" json:"id_empleado"`

	Caja     *Caja               `gorm:"foreignKey:CajaID;references:ID" json:"caja,omitempty"`
	Empleado *empleados.Empleado `gorm:"foreignKey:EmpleadoID;references:ID" json:"empleado,omitempty"`

	FechaInicio time.Time `json:"fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
