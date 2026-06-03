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
//obtiene todas las promociones
func ObtenerPromociones(c *gin.Context){
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)
	promociones, err := ObtenerTodasLasPromociones(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR":"Error al consultar bdd"})
		return
	}
	c.JSON(http.StatusOK, promociones)
}

//Obtiene una promocion por ID
func ObtenerPromocion(c *gin.Context){
	dbInstance, _:= c.Get("db")
	db := dbInstance.(*gorm.DB)
	id:= c.Param("id")
	promocion, err := ObtenerPromocionPorID(db,id)
	if err != nil{
		c.JSON(http.StatusNotFound,gin.H{"ERROR":"Promocion no encontrada"})
		return
	}
	c.JSON(http.StatusOK,promocion)
}
// Actualiza una promocion
func ActualizarPromociones(c *gin.Context){
	dbInstance, _:= c.Get("db")
	db := dbInstance.(*gorm.DB)
	id := c.Param("id")

	var datosNuevos Promocion
	if err := c.ShouldBindJSON(&datosNuevos); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"ERROR":"datos no validos" + err.Error()})
		return
	}
	promocionActualizada, err:= ActualizarPromocion(db,id,&datosNuevos)
	if err!= nil{
		c.JSON(http.StatusNotFound,gin.H{"ERROR":"No se pudo actualizar la promocion"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"mensaje":"Promocion actualizada","promocion":promocionActualizada})
}
//elimina una promocion
func EliminarPromociones(c *gin.Context){
	dbInstance, _:= c.Get("db")
	db := dbInstance.(*gorm.DB)
	id:= c.Param("id")

	err:=EliminarPromocion(db,id)
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"ERROR":"No se pudo eliminar la promocion"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje":"Promocion eliminada exitosamente"})
}
