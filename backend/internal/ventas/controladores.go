package ventas

import (
	"backend/internal/clientes"
	"backend/internal/inventario"
	"backend/internal/promociones"
	"errors"
	"math"
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
		if err := tx.Create(&nuevaVenta).Error; err != nil {
			return err
		}

		// 3. Si es venta fiada, guardar registro en la tabla fiados
		if nuevaVenta.MetodoID == "33333333-3333-3333-3333-333333333333" && nuevaVenta.ClienteID != nil && *nuevaVenta.ClienteID != "" {
			fiado := Fiado{
				ClienteID:   *nuevaVenta.ClienteID,
				VentaID:     nuevaVenta.ID,
				FechaInicio: time.Now(),
				FechaLimite: time.Now().AddDate(0, 1, 0), // Plazo de 30 días
				MontoTotal:  nuevaVenta.MontoTotal,
				Pagado:      false,
			}
			if err := tx.Create(&fiado).Error; err != nil {
				return err
			}

			// Actualizar la última compra del cliente en la base de datos
			now := time.Now()
			if err := tx.Model(&clientes.Cliente{}).Where("id = ?", *nuevaVenta.ClienteID).Update("ultima_compra", &now).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Preload MetodoPago y Fiado para consistencia en la respuesta
	ctrl.db.Preload("MetodoPago").Preload("Fiado").First(&nuevaVenta, "id = ?", nuevaVenta.ID)

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

type AbonoRequest struct {
	Monto float64 `json:"monto"`
}

func (ctrl *VentasController) RegistrarAbono(c *gin.Context) {
	clienteID := c.Param("id")
	var req AbonoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos: " + err.Error()})
		return
	}

	if req.Monto <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El monto del abono debe ser mayor a cero"})
		return
	}

	err := ctrl.db.Transaction(func(tx *gorm.DB) error {
		// 1. Obtener cliente
		var cliente clientes.Cliente
		if err := tx.First(&cliente, "id = ?", clienteID).Error; err != nil {
			return errors.New("Cliente no encontrado")
		}

		if cliente.FiadoActual <= 0 {
			return errors.New("El cliente no registra deuda pendiente")
		}

		// Asegurar que no abone más de lo que debe
		if req.Monto > cliente.FiadoActual {
			return errors.New("El abono no puede ser mayor que la deuda total del cliente")
		}

		// 2. Descontar del saldo acumulado del cliente
		cliente.FiadoActual = math.Max(cliente.FiadoActual-req.Monto, 0)
		if err := tx.Save(&cliente).Error; err != nil {
			return err
		}

		// 3. Buscar fiados pendientes ordenados por antigüedad
		var fiadosPendientes []Fiado
		if err := tx.Where("cliente_id = ? AND pagado = false", cliente.ID).Order("fecha_inicio ASC").Find(&fiadosPendientes).Error; err != nil {
			return err
		}

		montoRestante := req.Monto
		for i := range fiadosPendientes {
			if montoRestante >= fiadosPendientes[i].MontoTotal {
				fiadosPendientes[i].Pagado = true
				montoRestante -= fiadosPendientes[i].MontoTotal
				if err := tx.Save(&fiadosPendientes[i]).Error; err != nil {
					return err
				}
			} else {
				// Si el abono no cubre por completo esta deuda (o es el remanente),
				// marcamos todos los fiados que quedan activos para renovarles el plazo de 30 días
				for j := i; j < len(fiadosPendientes); j++ {
					fiadosPendientes[j].FechaInicio = time.Now()
					fiadosPendientes[j].FechaLimite = time.Now().AddDate(0, 1, 0) // Nuevos 30 días
					if err := tx.Save(&fiadosPendientes[j]).Error; err != nil {
						return err
					}
				}
				break
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Abono registrado con éxito"})
}
