package cajas

import (
	"backend/internal/empleados"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	EstadoTurnoAbierto   = "abierto"
	EstadoTurnoCerrado   = "cerrado"
	MetodoPagoEfectivoID = "11111111-1111-1111-1111-111111111111"
)

// TurnoCaja representa una sesión de caja asignada a un trabajador.
// El saldo esperado se incrementa automáticamente con cada venta en efectivo
// registrada mientras el turno permanezca abierto.
type TurnoCaja struct {
	ID            string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_turno"`
	UsuarioID     string     `gorm:"type:uuid;not null;index" json:"id_usuario"`
	CajaID        string     `gorm:"type:uuid;not null;index" json:"id_caja"`
	SaldoInicial  float64    `gorm:"type:decimal(12,2);default:0;not null" json:"saldo_inicial"`
	SaldoEsperado float64    `gorm:"type:decimal(12,2);default:0;not null" json:"saldo_esperado"`
	SaldoReal     float64    `gorm:"type:decimal(12,2);default:0;not null" json:"saldo_real"`
	Diferencia    float64    `gorm:"type:decimal(12,2);default:0;not null" json:"diferencia"`
	FechaApertura time.Time  `gorm:"not null;default:now()" json:"fecha_apertura"`
	FechaCierre   *time.Time `json:"fecha_cierre,omitempty"`
	Estado        string     `gorm:"type:varchar(16);default:'abierto';check:estado IN ('abierto','cerrado')" json:"estado"`

	Usuario empleados.Empleado `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
	Caja    Caja               `gorm:"foreignKey:CajaID" json:"caja,omitempty"`
}

type AbrirTurnoRequest struct {
	CajaID       string   `json:"id_caja" binding:"required,uuid4"`
	SaldoInicial *float64 `json:"saldo_inicial" binding:"required,gte=0"`
}

type CerrarTurnoRequest struct {
	SaldoReal *float64 `json:"saldo_real" binding:"required,gte=0"`
}

// CrearTurno persiste un nuevo turno de caja en estado abierto.
func CrearTurno(db *gorm.DB, turno *TurnoCaja) error {
	return db.Create(turno).Error
}

// ObtenerTurnoActivoPorUsuario devuelve el turno abierto actual de un empleado, si existe.
func ObtenerTurnoActivoPorUsuario(db *gorm.DB, usuarioID string) (*TurnoCaja, error) {
	var turno TurnoCaja
	err := db.Where("usuario_id = ? AND estado = ?", usuarioID, EstadoTurnoAbierto).First(&turno).Error
	if err != nil {
		return nil, err
	}
	return &turno, nil
}

// ObtenerTurnoPorID devuelve un turno por su ID.
func ObtenerTurnoPorID(db *gorm.DB, id string) (*TurnoCaja, error) {
	var turno TurnoCaja
	if err := db.First(&turno, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &turno, nil
}

// CerrarTurno finaliza un turno abierto, calculando la diferencia entre saldo real y esperado.
func CerrarTurno(db *gorm.DB, turnoID string, saldoReal float64) (*TurnoCaja, error) {
	var turno TurnoCaja
	if err := db.First(&turno, "id = ?", turnoID).Error; err != nil {
		return nil, err
	}

	if turno.Estado == EstadoTurnoCerrado {
		return nil, errors.New("El turno ya se encuentra cerrado")
	}

	ahora := time.Now()
	diferencia := saldoReal - turno.SaldoEsperado

	turno.SaldoReal = saldoReal
	turno.Diferencia = diferencia
	turno.FechaCierre = &ahora
	turno.Estado = EstadoTurnoCerrado

	if err := db.Save(&turno).Error; err != nil {
		return nil, err
	}
	return &turno, nil
}

// TurnoController expone las operaciones de apertura/cierre de turnos de caja.
type TurnoController struct {
	db *gorm.DB
}

func NewTurnoController(db *gorm.DB) *TurnoController {
	return &TurnoController{db: db}
}

// AbrirTurno crea una nueva sesión de caja para el usuario autenticado.
func (ctrl *TurnoController) AbrirTurno(c *gin.Context) {
	idEmpleado, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró sesión de empleado"})
		return
	}
	usuarioID := idEmpleado.(string)

	var req AbrirTurnoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Verificar que la caja exista y esté activa
	var caja Caja
	if err := ctrl.db.First(&caja, "id = ? AND activo = ?", req.CajaID, true).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La caja seleccionada no existe o está inactiva"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar la caja"})
		return
	}

	// Evitar múltiples turnos abiertos para el mismo usuario
	turnoExistente, err := ObtenerTurnoActivoPorUsuario(ctrl.db, usuarioID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar turnos activos"})
		return
	}
	if turnoExistente != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":    "Ya existe un turno abierto para este usuario",
			"id_turno": turnoExistente.ID,
		})
		return
	}

	nuevoTurno := TurnoCaja{
		UsuarioID:     usuarioID,
		CajaID:        req.CajaID,
		SaldoInicial:  *req.SaldoInicial,
		SaldoEsperado: *req.SaldoInicial,
		SaldoReal:     0,
		Diferencia:    0,
		FechaApertura: time.Now(),
		Estado:        EstadoTurnoAbierto,
	}

	if err := CrearTurno(ctrl.db, &nuevoTurno); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo abrir el turno de caja"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Turno de caja abierto correctamente",
		"id_turno": nuevoTurno.ID,
		"id_caja":  nuevoTurno.CajaID,
	})
}

// CerrarTurno finaliza el turno activo del usuario autenticado.
func (ctrl *TurnoController) CerrarTurno(c *gin.Context) {
	idEmpleado, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró sesión de empleado"})
		return
	}
	usuarioID := idEmpleado.(string)

	var req CerrarTurnoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	turno, err := ObtenerTurnoActivoPorUsuario(ctrl.db, usuarioID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No hay un turno abierto para cerrar"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el turno activo"})
		return
	}

	turnoCerrado, err := CerrarTurno(ctrl.db, turno.ID, *req.SaldoReal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Turno de caja cerrado correctamente",
		"turno":   turnoCerrado,
	})
}

// TurnoActivo devuelve el turno abierto actual del usuario autenticado, si lo tiene.
func (ctrl *TurnoController) TurnoActivo(c *gin.Context) {
	idEmpleado, existe := c.Get("id_empleado")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró sesión de empleado"})
		return
	}

	turno, err := ObtenerTurnoActivoPorUsuario(ctrl.db, idEmpleado.(string))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"activo": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el turno activo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activo": true,
		"turno":  turno,
	})
}

// ActualizarSaldoEsperado suma un monto al saldo esperado del turno activo de un usuario.
// Se utiliza desde el módulo de ventas para reflejar las ventas en efectivo.
func ActualizarSaldoEsperado(db *gorm.DB, usuarioID string, monto float64) error {
	return db.Model(&TurnoCaja{}).
		Where("usuario_id = ? AND estado = ?", usuarioID, EstadoTurnoAbierto).
		Update("saldo_esperado", gorm.Expr("saldo_esperado + ?", monto)).Error
}
