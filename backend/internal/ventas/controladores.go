package ventas

import (
	"backend/internal/inventario"
	"backend/internal/promociones"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func calcularDescuentoItem(cantidad int, precio float64, promo promociones.Promocion) float64 {
	switch promo.Tipo {
	case "NXM":
		if promo.Lleva > 0 && promo.Paga > 0 && promo.Lleva > promo.Paga {
			sets := cantidad / promo.Lleva
			descuentoCant := sets * (promo.Lleva - promo.Paga)
			return float64(descuentoCant) * precio
		}
	case "porcentaje":
		if promo.Descuento > 0 {
			return float64(cantidad) * precio * (promo.Descuento / 100.0)
		}
	case "precio_fijo":
		if promo.Descuento > 0 && precio > promo.Descuento {
			unitDiscount := precio - promo.Descuento
			return float64(cantidad) * unitDiscount
		}
	}
	return 0
}

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
		// Cargar todas las promociones activas
		var activePromos []promociones.Promocion
		now := time.Now()
		if err := tx.Where("(fecha_inicio IS NULL OR fecha_inicio <= ?) AND (fecha_fin IS NULL OR fecha_fin >= ?)", now, now).Find(&activePromos).Error; err != nil {
			return err
		}

		var totalDescuento float64 = 0
		var totalVenta float64 = 0

		// 1. Validar, descontar stock y calcular promociones por ítem
		for i, detalle := range nuevaVenta.Detalles {
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

			// Buscar la promoción activa que otorgue el mayor descuento para este ítem
			var bestDiscount float64 = 0
			for _, promo := range activePromos {
				if promo.ProductoID == detalle.ProductoID {
					disc := calcularDescuentoItem(detalle.Cantidad, producto.Precio, promo)
					if disc > bestDiscount {
						bestDiscount = disc
					}
				}
			}

			subtotalItem := producto.Precio * float64(detalle.Cantidad)
			montoFinalItem := subtotalItem - bestDiscount
			nuevaVenta.Detalles[i].MontoFinal = montoFinalItem

			totalDescuento += bestDiscount
			totalVenta += montoFinalItem
		}

		// Sobrescribir los montos totales calculados por seguridad en el backend
		nuevaVenta.MontoDescuento = totalDescuento
		nuevaVenta.MontoTotal = totalVenta

		// 2. Guardar la venta
		return tx.Create(&nuevaVenta).Error
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Preload MetodoPago para consistencia en la respuesta
	ctrl.db.Preload("MetodoPago").First(&nuevaVenta, "id = ?", nuevaVenta.ID)

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Venta creada",
		"venta":   nuevaVenta,
	})
}

func (ctrl *VentasController) GetVentas(c *gin.Context) {
	ventas, err := ObtenerTodasVentas(ctrl.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al obtener las ventas"})
		return
	}

	c.JSON(http.StatusOK, ventas)
}

func (ctrl *VentasController) GetVentaByID(c *gin.Context) {
	id := c.Param("id")

	venta, err := ObtenerVentaPorID(ctrl.db, id)
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

func (ctrl *VentasController) UpdateVenta(c *gin.Context) {
	id := c.Param("id")

	ventaExistente, err := ObtenerVentaPorID(ctrl.db, id)
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

	if err := ActualizarVenta(ctrl.db, ventaExistente, &datosNuevos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo actualizar la venta"})
		return
	}

	ventaActualizada, _ := ObtenerVentaPorID(ctrl.db, id)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Venta actualizada exitosamente",
		"venta":   ventaActualizada,
	})
}

func (ctrl *VentasController) DeleteVenta(c *gin.Context) {
	id := c.Param("id")

	if err := EliminarVenta(ctrl.db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo eliminar la venta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Venta eliminada exitosamente"})
}