package auth

import (
	"github.com/gin-gonic/gin"
)

func RoutesConfig(api *gin.RouterGroup, ctrl *AuthController, authMiddleware gin.HandlerFunc) {

	group := api.Group("auth")
	{
		group.POST("/login", ctrl.Login)

		authGroup := group.Group("")
		authGroup.Use(authMiddleware)
		{
			authGroup.POST("/logout", ctrl.Logout)
			authGroup.GET("/me", ctrl.Me)
		}

	}
}
