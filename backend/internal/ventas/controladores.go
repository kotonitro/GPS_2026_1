package ventas

import (
	"errors"
	"net/http"
	"time"
	"backend/internal/inventario"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VentasController struct {
	db *gorm.DB
}

func NewVentasController(db *gorm.DB) *VentasController {
	return &VentasController{db: db}
}

func (ctrl *VentasController) CrearVenta(c *gin.Context) {
	var nuevaVenta Venta

	if err := c.ShouldBindJSON(&nuevaVenta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos: " + err.Error()})
		return
	}

	// Extraer empleado del contexto
	idEmpleado, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "No se encontró sesión de empleado"})
		return
	}
	nuevaVenta.EmpleadoID = idEmpleado.(string)
	nuevaVenta.FechaEmision = time.Now()

	err := ctrl.db.Transaction(func(tx *gorm.DB) error {
		// 1. Validar y descontar stock
		for _, detalle := range nuevaVenta.Detalles {
			var producto inventario.Producto
			if err := tx.First(&producto, "id = ?", detalle.ProductoID).Error; err != nil {
				return errors.New("El producto " + detalle.ProductoID + " no existe")
			}
			
			if producto.Stock < detalle.Cantidad {
				return errors.New("Stock insuficiente para: " + producto.Nombre)
			}
			
			producto.Stock -= detalle.Cantidad
			if err := tx.Save(&producto).Error; err != nil {
				return err
			}
		}

		// 2. Guardar la venta
		return tx.Create(&nuevaVenta).Error
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	ctrl.db.Preload("MetodoPago").First(&nuevaVenta, "id = ?", nuevaVenta.ID)
	c.JSON(http.StatusCreated, gin.H{
    "mensaje": "Venta creada", 
    "venta":   nuevaVenta,
	})
}

func (ctrl *VentasController) GetVentas(c *gin.Context) {
	ventas, err := ObtenerTodasVentas(ctrl.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al obtener ventas"})
		return
	}
	c.JSON(http.StatusOK, ventas)
}

func (ctrl *VentasController) GetVentaByID(c *gin.Context) {
	id := c.Param("id")
	venta, err := ObtenerVentaPorID(ctrl.db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Venta no encontrada"})
		return
	}
	c.JSON(http.StatusOK, venta)
}

func (ctrl *VentasController) UpdateVenta(c *gin.Context) {
	id := c.Param("id")
	ventaExistente, _ := ObtenerVentaPorID(ctrl.db, id)
	var datosNuevos Venta
	if err := c.ShouldBindJSON(&datosNuevos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos"})
		return
	}
	ActualizarVenta(ctrl.db, ventaExistente, &datosNuevos)
	c.JSON(http.StatusOK, gin.H{"mensaje": "Venta actualizada"})
}

func (ctrl *VentasController) DeleteVenta(c *gin.Context) {
	id := c.Param("id")
	if err := EliminarVenta(ctrl.db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo eliminar: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Venta eliminada"})
}