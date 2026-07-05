package cajas

import "time"

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
