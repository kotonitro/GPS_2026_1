package cajas

import "gorm.io/gorm"

func GetCajas(db *gorm.DB) ([]Caja, error) {
	var cajas []Caja

	result := db.Find(&cajas)
	return cajas, result.Error
}

// GetCajasConTurnoActivo devuelve la lista de cajas enriquecida con el turno abierto más reciente
// de cada una, incluyendo el responsable y el saldo esperado actual.
func GetCajasConTurnoActivo(db *gorm.DB) ([]Caja, error) {
	var cajas []Caja
	if err := db.Find(&cajas).Error; err != nil {
		return nil, err
	}

	for i := range cajas {
		var turno TurnoCaja
		err := db.
			Preload("Usuario").
			Where("caja_id = ? AND estado = ?", cajas[i].ID, EstadoTurnoAbierto).
			Order("fecha_apertura DESC").
			First(&turno).Error

		if err == nil {
			cajas[i].TurnoActivo = &turno
		}
	}

	return cajas, nil
}

func GetCajaByID(db *gorm.DB, id string) (*Caja, error) {
	var caja Caja

	result := db.First(&caja, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &caja, nil
}

func CreateCaja(db *gorm.DB, caja *Caja) error {

	resultado := db.Create(caja)

	return resultado.Error
}

func DeleteCajaByID(db *gorm.DB, id string) error {
	var caja Caja

	if err := db.First(&caja, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Delete(&caja).Error; err != nil {
		return err
	}

	return nil
}

func UpdateCajaByID(db *gorm.DB, id string, data UpdateCajaInput) error {
	var caja Caja

	if err := db.First(&caja, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Model(&caja).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
