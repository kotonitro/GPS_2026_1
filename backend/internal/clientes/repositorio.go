package clientes

import "gorm.io/gorm"

//cuando creamos clientes, lo guardamos en la db
func GuardarCliente(db *gorm.DB, nuevoCliente *Cliente) error {
	result := db.Create(nuevoCliente)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//verificar si el rut ya existe
func CheckRutExiste(db *gorm.DB, rut string) (bool, error) {
	var count int64

	err := db.Model(&Cliente{}).Where("rut = ?", rut).Count(&count).Error
	if err != nil {

		return false, err
	}

	return count > 0, nil
}

//obtener todos los clientes
func ObtenerTodosLosClientes(db *gorm.DB) ([]Cliente, error) {
	var listaClientes []Cliente
	
	
	result := db.Find(&listaClientes) //busca registros
	if result.Error != nil {
		return nil, result.Error
	}
	
	return listaClientes, nil
}

//obtener por id
func ObtenerClientePorID(db *gorm.DB, id string) (*Cliente, error) {
	var cliente Cliente
	result := db.First(&cliente, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cliente, nil
}

//actualizar 
func ActualizarCliente(db *gorm.DB, clienteExistente *Cliente, datosNuevos *Cliente) error {
	// Updates actualiza los campos modificados basados en el modelo estructurado
	result := db.Model(clienteExistente).Updates(datosNuevos)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//eliminar
func EliminarCliente(db *gorm.DB, id string) error {
	// Pasamos un struct vacío de Cliente con el ID para indicarle a GORM qué borrar
	result := db.Delete(&Cliente{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}