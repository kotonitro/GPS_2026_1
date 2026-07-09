package cajas

import (
	"backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, authMiddleware gin.HandlerFunc) {

	ctrl := NewCajaController(db)
	turnoCtrl := NewTurnoController(db)

	authGroup := api.Group("cajas")
	authGroup.Use(authMiddleware)
	{
		authGroup.GET("", ctrl.GetCajasController)
		authGroup.GET("/:id", ctrl.GetCajaByIDController)

		// Turnos de caja (apertura/cierre por trabajador)
		authGroup.GET("/turnos/activo", turnoCtrl.TurnoActivo)
		authGroup.POST("/turnos/apertura", turnoCtrl.AbrirTurno)
		authGroup.POST("/turnos/cierre", turnoCtrl.CerrarTurno)

		adminGroup := authGroup.Group("")
		adminGroup.Use(auth.AdminMiddleware())
		{
			adminGroup.POST("", ctrl.CreateCajaController)
			adminGroup.DELETE("/:id", ctrl.DeleteCajaByIDController)
			adminGroup.PATCH("/:id", ctrl.UpdateCajaByIDController)
		}
	}
}
