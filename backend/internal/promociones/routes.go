package promociones

import "github.com/gin-gonic/gin"

func ConfigurarRutas(api *gin.RouterGroup){
	grupo := api.Group("/promociones")
	{
		grupo.POST("/promocion",CrearPromocion)
		grupo.GET("/promocion",ObtenerPromociones)
		grupo.GET("/promocion/:id",ObtenerPromocion)
		grupo.DELETE("/promocion/:id",EliminarPromociones)
		grupo.PUT("/promocion/:id",ActualizarPromociones)
	}
}