package empleados

import "gorm.io/gorm"

func GetEmpleados(db *gorm.DB) ([]Empleado, error) {
	var empleados []Empleado

	result := db.Preload("Rol").Find(&empleados)
	return empleados, result.Error
}

func GetEmpleadoByID(db *gorm.DB, id string) (*Empleado, error) {
	var empleado Empleado

	result := db.Preload("Rol").First(&empleado, "id = ?", id)
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

func GetRoles(db *gorm.DB) ([]Rol, error) {
	var roles []Rol

	result := db.Find(&roles)
	return roles, result.Error
}

func GetRolByID(db *gorm.DB, id string) (*Rol, error) {
	var rol Rol
	result := db.First(&rol, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &rol, nil
}

func CreateRol(db *gorm.DB, rol *Rol) error {

	resultado := db.Create(rol)

	return resultado.Error
}

func DeleteRolByID(db *gorm.DB, id string) error {
	var rol Rol

	if err := db.First(&rol, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Delete(&rol).Error; err != nil {
		return err
	}

	return nil
}

func UpdateRolByID(db *gorm.DB, id string, data UpdateRolInput) error {
	var rol Rol

	if err := db.First(&rol, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Model(&rol).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
