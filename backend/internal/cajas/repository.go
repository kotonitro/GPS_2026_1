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
