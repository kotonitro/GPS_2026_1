package clientes

import "gorm.io/gorm"
// obetener todos los clientes de la DB
func GetClientes(db *gorm.DB) ([]Cliente, error) {
	var clientes []Cliente

	result := db.Find(&clientes)
	return clientes, result.Error
}
// Obtener un cliente por su ID
func GetClienteByID(db *gorm.DB, id string) (*Cliente, error) {
	var cliente Cliente

	result := db.First(&cliente, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &cliente, nil
}
//crear cliente
func CreateCliente(db *gorm.DB, cliente *Cliente) error {
	resultado := db.Create(cliente)
	return resultado.Error
}
//borrar cliente usando ID
func DeleteClienteByID(db *gorm.DB, id string) error {
	var cliente Cliente

	if err := db.First(&cliente, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Delete(&cliente).Error; err != nil {
		return err
	}

	return nil
}
//modicifar cliente
func UpdateClienteByID(db *gorm.DB, id string, data UpdateClienteInput) error {
	var cliente Cliente

	if err := db.First(&cliente, "id = ?", id).Error; err != nil {
		return err
	}

	if err := db.Model(&cliente).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
//obtener cliente por rut
func GetClienteByRut(db *gorm.DB, rut string) (*Cliente, error) {
	var cliente Cliente

	result := db.Where("rut = ?", rut).First(&cliente)
	if result.Error != nil {
		return nil, result.Error
	}

	return &cliente, nil
}
//obtener cliente por nombre
func GetClientesByNombre(db *gorm.DB, nombre string) ([]Cliente, error) {
	var clientes []Cliente

	result := db.Where("nombre ILIKE ?", "%"+nombre+"%").Find(&clientes)
	if result.Error != nil {
		return nil, result.Error
	}

	return clientes, nil
}
