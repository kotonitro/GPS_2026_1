package ventas

import (
	"backend/internal/cajas"
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
		if promo.Lleva != nil && promo.Paga != nil {
			lleva := *promo.Lleva
			paga := *promo.Paga
			if lleva > 0 && paga > 0 && lleva > paga {
				sets := cantidad / lleva
				descuentoCant := sets * (lleva - paga)
				return float64(descuentoCant) * precio
			}
		}
	case "porcentaje":
		if promo.Descuento != nil {
			descuento := *promo.Descuento
			if descuento > 0 {
				return float64(cantidad) * precio * (descuento / 100.0)
			}
		}
	case "precio_fijo":
		if promo.Descuento != nil {
			descuento := *promo.Descuento
			if descuento > 0 && precio > descuento {
				unitDiscount := precio - descuento
				return float64(cantidad) * unitDiscount
			}
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

type SyncVentasRequest struct {
	Ventas []Venta `json:"ventas" binding:"required"`
}

type SyncResultado struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// procesarVentaIndividual centraliza la lógica de creación de una venta.
// Si respectOriginalDate es true y la venta trae una FechaEmision válida, se respeta;
// de lo contrario se asigna la fecha actual. Esto permite sincronizar ventas offline.
func (ctrl *VentasController) procesarVentaIndividual(tx *gorm.DB, venta *Venta, empleadoID string, respectOriginalDate bool) error {
	venta.EmpleadoID = empleadoID

	if !respectOriginalDate || venta.FechaEmision.IsZero() {
		venta.FechaEmision = time.Now()
	}

	// Cargar todas las promociones activas
	var activePromos []promociones.Promocion
	now := time.Now()
	if err := tx.Where("(fecha_inicio IS NULL OR fecha_inicio <= ?) AND (fecha_fin IS NULL OR fecha_fin >= ?)", now, now).Find(&activePromos).Error; err != nil {
		return err
	}

	var totalDescuento float64 = 0
	var totalVenta float64 = 0

	// 1. Validar, descontar stock y calcular promociones por ítem
	for i, detalle := range venta.Detalles {
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
		var totalDescuento float64 = 0
		var subtotalTotal float64 = 0

		cantidades := make(map[string]int)
		precios := make(map[string]float64)

		// 1. Validar, descontar stock y preparar datos
		for i, detalle := range nuevaVenta.Detalles {
			var producto inventario.Producto
			if err := tx.First(&producto, "id = ?", detalle.ProductoID).Error; err != nil {
				return errors.New("El producto " + detalle.ProductoID + " no existe")
			}

		// Buscar la promoción activa que otorgue el mayor descuento para este ítem
		var bestDiscount float64 = 0
		for _, promo := range activePromos {
			// Validar que ProductoID no sea nil antes de desreferenciarlo
			if promo.ProductoID != nil && *promo.ProductoID == detalle.ProductoID {
				disc := calcularDescuentoItem(detalle.Cantidad, producto.Precio, promo)
				if disc > bestDiscount {
					bestDiscount = disc
				}
			}
		}

		subtotalItem := producto.Precio * float64(detalle.Cantidad)
		montoFinalItem := subtotalItem - bestDiscount
		venta.Detalles[i].MontoFinal = montoFinalItem

		totalDescuento += bestDiscount
		totalVenta += montoFinalItem
	}

	// Sobrescribir los montos totales calculados por seguridad en el backend
	venta.MontoDescuento = totalDescuento
	venta.MontoTotal = totalVenta

	// 2. Guardar la venta
	if err := tx.Create(&venta).Error; err != nil {
		return err
	}

	// 3. Si es venta en efectivo, actualizar el saldo esperado del turno activo del usuario
	if venta.MetodoID == cajas.MetodoPagoEfectivoID {
		var turnoActivo cajas.TurnoCaja
		if err := tx.Where("usuario_id = ? AND estado = ?", empleadoID, cajas.EstadoTurnoAbierto).First(&turnoActivo).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("No hay un turno de caja abierto para registrar ventas en efectivo")
			}
			return err
		}

		if err := cajas.ActualizarSaldoEsperado(tx, empleadoID, venta.MontoTotal); err != nil {
			return err
		}
	}

	// 4. Si es venta fiada, guardar registro en la tabla fiados
	if venta.MetodoID == "33333333-3333-3333-3333-333333333333" && venta.ClienteID != nil && *venta.ClienteID != "" {
		fiadoFechaInicio := venta.FechaEmision
		if fiadoFechaInicio.IsZero() {
			fiadoFechaInicio = time.Now()
		}

		fiado := Fiado{
			ClienteID:   *venta.ClienteID,
			VentaID:     venta.ID,
			FechaInicio: fiadoFechaInicio,
			FechaLimite: fiadoFechaInicio.AddDate(0, 1, 0), // Plazo de 30 días
			MontoTotal:  venta.MontoTotal,
			Pagado:      false,
		}
		if err := tx.Create(&fiado).Error; err != nil {
			return err
		}
			cantidades[detalle.ProductoID] += detalle.Cantidad
			precios[detalle.ProductoID] = producto.Precio

			subtotalItem := producto.Precio * float64(detalle.Cantidad)
			subtotalTotal += subtotalItem
			nuevaVenta.Detalles[i].MontoFinal = subtotalItem // Monto base, el descuento se calcula globalmente
		}

		// 2. Calcular Descuentos de Combos
		for _, promo := range activePromos {
			if promo.Tipo == "COMBO" && len(promo.ProductosCombo) > 0 && promo.Descuento != nil {
				sets := 999999
				var costoNormalCombo float64 = 0

				reqMap := make(map[string]int)
				for _, pid := range promo.ProductosCombo {
					reqMap[pid]++
				}

				for pid, reqQty := range reqMap {
					if reqQty > 0 {
						avail := cantidades[pid]
						possible := avail / reqQty
						if possible < sets {
							sets = possible
						}
						costoNormalCombo += precios[pid] * float64(reqQty)
					}
				}

				if sets > 0 && sets != 999999 {
					// El descuento es la diferencia entre el precio normal y el precio final del combo
					ahorroPorCombo := costoNormalCombo - *promo.Descuento
					if ahorroPorCombo > 0 {
						totalDescuento += ahorroPorCombo * float64(sets)
						for pid, reqQty := range reqMap {
							cantidades[pid] -= reqQty * sets
						}
					}
				}
			}
		}

		// 3. Calcular Descuentos Individuales para el remanente
		for pid, remanente := range cantidades {
			if remanente > 0 {
				var bestDiscount float64 = 0
				for _, promo := range activePromos {
					if promo.Tipo != "COMBO" && promo.ProductoID != nil && *promo.ProductoID == pid {
						disc := calcularDescuentoItem(remanente, precios[pid], promo)
						if disc > bestDiscount {
							bestDiscount = disc
						}
					}
				}
				totalDescuento += bestDiscount
			}
		}

		totalVenta := subtotalTotal - totalDescuento

		// Sobrescribir los montos totales calculados por seguridad en el backend
		nuevaVenta.MontoDescuento = totalDescuento
		nuevaVenta.MontoTotal = totalVenta

		// Actualizar la última compra del cliente en la base de datos
		if err := tx.Model(&clientes.Cliente{}).Where("id = ?", *venta.ClienteID).Update("ultima_compra", &fiadoFechaInicio).Error; err != nil {
			return err
		}
	}
	return nil
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

	err := ctrl.db.Transaction(func(tx *gorm.DB) error {
		return ctrl.procesarVentaIndividual(tx, &nuevaVenta, idEmpleado.(string), false)
			// Actualizar la última compra del cliente en la base de datos
			now := time.Now()
			if err := tx.Model(&clientes.Cliente{}).Where("id = ?", *nuevaVenta.ClienteID).Update("ultima_compra", &now).Error; err != nil {
				return err
			}
		}

		// 4. Actualizar el saldo final de la caja
		// Aumentamos el saldo de la caja (asumiendo que las ventas al contado o tarjeta suman al cuadre)
		// Si es Fiado, el dinero no entra a la caja en este momento.
		if nuevaVenta.MetodoID != "33333333-3333-3333-3333-333333333333" {
			var caja cajas.Caja
			if err := tx.First(&caja, "id = ?", nuevaVenta.CajaID).Error; err != nil {
				return errors.New("Error al buscar la caja para actualizar el saldo: " + err.Error())
			}
			caja.SaldoFinal += uint(nuevaVenta.MontoTotal)
			if err := tx.Save(&caja).Error; err != nil {
				return errors.New("Error al actualizar el saldo de la caja: " + err.Error())
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Preload MetodoPago, Detalles, Fiado y Caja para consistencia en la respuesta
	ctrl.db.Preload("MetodoPago").Preload("Detalles").Preload("Fiado").Preload("Caja").First(&nuevaVenta, "id = ?", nuevaVenta.ID)

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Venta creada",
		"venta":   nuevaVenta,
	})
}

// SyncVentasOffline recibe un lote de ventas creadas offline, valida duplicados por UUID
// y las procesa respetando las fechas originales enviadas por el cliente.
func (ctrl *VentasController) SyncVentasOffline(c *gin.Context) {
	var req SyncVentasRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos: " + err.Error()})
		return
	}

	idEmpleado, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "No se encontró sesión de empleado"})
		return
	}

	var resultados []SyncResultado

	for _, venta := range req.Ventas {
		if venta.ID == "" {
			resultados = append(resultados, SyncResultado{ID: "", Status: "error", Error: "El UUID de la venta es requerido"})
			continue
		}

		// Validar duplicado por UUID
		var existing Venta
		if err := ctrl.db.First(&existing, "id = ?", venta.ID).Error; err == nil {
			resultados = append(resultados, SyncResultado{ID: venta.ID, Status: "duplicado", Error: "La venta ya fue sincronizada"})
			continue
		}

		err := ctrl.db.Transaction(func(tx *gorm.DB) error {
			return ctrl.procesarVentaIndividual(tx, &venta, idEmpleado.(string), true)
		})

		if err != nil {
			resultados = append(resultados, SyncResultado{ID: venta.ID, Status: "error", Error: err.Error()})
		} else {
			resultados = append(resultados, SyncResultado{ID: venta.ID, Status: "ok"})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":    "Sincronización completada",
		"resultados": resultados,
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
