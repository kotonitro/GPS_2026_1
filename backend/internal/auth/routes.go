package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, jwtSecret string) {

	ctrl := NewAuthController(db, jwtSecret)

	group := api.Group("auth")
	{
		group.POST("/login", ctrl.Login)

		authGroup := group.Group("")
		authGroup.Use(AuthMiddleware(jwtSecret))
		{
			authGroup.POST("/logout", ctrl.Logout)
			authGroup.GET("/me", ctrl.Me)
		}

	}
}
