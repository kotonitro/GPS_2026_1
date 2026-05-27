package promociones

import "gorm.io/gorm"

func GuardarPromocion(db *gorm.DB, promocion *Promocion) error{
	result :=db.Create(promocion)
	return result.Error
}