package cajas

import "time"

type Caja struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_caja"`
	Ubicacion    string    `gorm:"type:varchar(255);not null" json:"ubicacion" binding:"required"`
	SaldoInicial float64   `gorm:"type:numeric(12,2);not null" json:"saldo_inicial" binding:"required"`
	SaldoFinal   float64   `gorm:"type:numeric(12,2);default:0" json:"saldo_final"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
