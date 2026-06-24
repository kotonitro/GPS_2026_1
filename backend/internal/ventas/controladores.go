package ventas

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CrearVenta(c *gin.Context) {
	var nuevaVenta Venta

	if err := c.ShouldBindJSON(&nuevaVenta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos"})
		return
	}
	if nuevaVenta.MontoTotal <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El monto total de la venta debe ser mayor a cero"})
		return
	}
	if nuevaVenta.MetodoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El ID del método de pago no puede estar vacío"})
		return
	}
	if nuevaVenta.EmpleadoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El ID del empleado no puede estar vacío"})
		return
	}
	if nuevaVenta.CajaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El ID de la caja no puede estar vacío"})
		return
	}

	if len(nuevaVenta.Detalles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "La venta debe contener al menos un producto en el detalle"})
		return
	}

	for _, detalle := range nuevaVenta.Detalles {
		if detalle.ProductoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Hay un producto en la lista que no tiene un ID válido"})
			return
		}
		if detalle.Cantidad <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "La cantidad de cada producto debe ser mayor a cero"})
			return
		}
		if detalle.MontoFinal <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "El monto final de cada detalle debe ser mayor a cero"})
			return
		}
	}
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	err := GuardarVenta(db, &nuevaVenta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo guardar la venta"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Venta agregada",
		"venta":   nuevaVenta,
	})
}

func GetVentas(c *gin.Context) {
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	ventas, err := ObtenerTodasVentas(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al obtener las ventas"})
		return
	}

	c.JSON(http.StatusOK, ventas)
}

func GetVentaByID(c *gin.Context) {
	id := c.Param("id")

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	venta, err := ObtenerVentaPorID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Venta no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al buscar la venta"})
		return
	}

	c.JSON(http.StatusOK, venta)
}

func UpdateVenta(c *gin.Context) {
	id := c.Param("id")

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	ventaExistente, err := ObtenerVentaPorID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "La venta que intenta actualizar no existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al buscar la venta"})
		return
	}

	var datosNuevos Venta
	if err := c.ShouldBindJSON(&datosNuevos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos"})
		return
	}

	if err := ActualizarVenta(db, ventaExistente, &datosNuevos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo actualizar la venta"})
		return
	}

	ventaActualizada, _ := ObtenerVentaPorID(db, id)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Venta actualizada exitosamente",
		"venta":   ventaActualizada,
	})
}

func DeleteVenta(c *gin.Context) {
	id := c.Param("id")

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	if err := EliminarVenta(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo eliminar la venta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Venta eliminada exitosamente"})
}
