package promociones

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CrearPromocion (c *gin.Context){
	var nuevaPromocion Promocion
	//lee el json
	if err := c.ShouldBindJSON(&nuevaPromocion); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "datos invalidos:" + err.Error()})
		return
	}

	//validacion
	switch nuevaPromocion.Tipo{
		case "NXM":
			if nuevaPromocion.Lleva <=0 || nuevaPromocion.Paga <=0{
					c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Para promociones NXM lleva y paga deben ser mayores a 0"})
					return
			}
		case "porcentaje","precio_fijo":
			if nuevaPromocion.Descuento <=0{
					c.JSON(http.StatusBadRequest,gin.H{"ERROR":"Para la promocion este valor debe ser mayor a 0"})
					return
			}
		default:
				c.JSON(http.StatusBadRequest,gin.H{"ERROR":"Tipo de promocion no valida"})
				return
	}
	//conexion bdd
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	//repositorio
	err := GuardarPromocion(db,&nuevaPromocion)
	if err!= nil{
			c.JSON(http.StatusInternalServerError,gin.H{"ERROR":"No se pudo guardar la promocion en la bdd"})
			return
	}

	//mensaje
	c.JSON(http.StatusCreated, gin.H{
		"mensaje":"promocion creada exitosamente",
		"promocion":nuevaPromocion,
	})
}
func GuardarPromocion(db *gorm.DB, promocion *Promocion) error{
	result :=db.Create(promocion)
	return result.Error
}