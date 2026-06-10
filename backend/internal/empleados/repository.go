package empleados

import "gorm.io/gorm"

func GetEmpleados(db *gorm.DB) ([]Empleado, error) {
	var empleados []Empleado

	result := db.Find(&empleados)
	return empleados, result.Error
}

func GetEmpleadoByID(db *gorm.DB, id string) (*Empleado, error) {
	var empleado Empleado

	result := db.First(&empleado, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &empleado, nil
}

func CreateEmpleado(db *gorm.DB, empleado *Empleado) error {

	resultado := db.Create(empleado)

	return resultado.Error
}

func DeleteEmpleadoByID(db *gorm.DB, id string) error {
	var empleado Empleado

	if err := db.First(&empleado, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Delete(&empleado).Error; err != nil {
		return err
	}

	return nil
}

func UpdateEmpleadoByID(db *gorm.DB, id string, data UpdateEmpleadoInput) error {
	var empleado Empleado

	if err := db.First(&empleado, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Model(&empleado).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
