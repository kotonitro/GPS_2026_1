package promociones

import "gorm.io/gorm"

func GuardarPromocion(db *gorm.DB, promocion *Promocion) error{
	result :=db.Create(promocion)
	return result.Error
}

//Muestra todas las promociones de la bdd
func ObtenerTodasLasPromociones(db *gorm.DB)([]Promocion,error){
	var promociones []Promocion
	result :=db.Find(&promociones)
	return promociones, result.Error
}

//Muestra una promocion seleccionada por ID
func ObtenerPromocionPorID(db *gorm.DB, id string) (Promocion,error){
	var promocion Promocion
	result := db.First(&promocion, id)
	return promocion, result.Error
}

//Actualiza una promocion mediante su id
func ActualizarPromocion(db *gorm.DB, id string, datosActualizados *Promocion) (Promocion,error){
	var promocion Promocion
	if err:= db.First(&promocion, id).Error; err != nil{
		return promocion,err
	}
	result := db.Model(&promocion).Updates(datosActualizados)
	return promocion, result.Error
}

//elimina promocion
func EliminarPromocion(db *gorm.DB, id string)error{
	result := db.Delete(&Promocion{},id)
	return result.Error
}