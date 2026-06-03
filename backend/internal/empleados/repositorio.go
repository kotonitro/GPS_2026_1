package empleados

import "gorm.io/gorm"

func GuardarEmpleado(db *gorm.DB, empleado *Empleado) error {

	resultado := db.Create(empleado)

	return resultado.Error
}
