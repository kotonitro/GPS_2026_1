package cajas

import "gorm.io/gorm"

func GetCajas(db *gorm.DB) ([]Caja, error) {
	var cajas []Caja

	result := db.Find(&cajas)
	return cajas, result.Error
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

func GetRegistros(db *gorm.DB) ([]RegistroTurno, error) {
	var registros []RegistroTurno

	result := db.Find(&registros)
	return registros, result.Error
}

func GetRegistroByID(db *gorm.DB, id string) (*RegistroTurno, error) {
	var registro RegistroTurno

	result := db.First(&registro, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &registro, nil
}

func CreateRegistro(db *gorm.DB, registro *RegistroTurno) error {

	resultado := db.Create(registro)

	return resultado.Error
}

func DeleteRegistroByID(db *gorm.DB, id string) error {
	var registro RegistroTurno

	if err := db.First(&registro, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Delete(&registro).Error; err != nil {
		return err
	}

	return nil
}

func UpdateRegistroByID(db *gorm.DB, id string, data UpdateRegistroInput) error {
	var registro RegistroTurno

	if err := db.First(&registro, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Model(&registro).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
