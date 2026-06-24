package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesConfig(api *gin.RouterGroup, db *gorm.DB, jwtSecret string) {

	ctrl := NewAuthController(db, jwtSecret)

	grupo := api.Group("auth")
	{
		grupo.POST("/login", ctrl.Login)
		grupo.POST("/logout", ctrl.Logout)
	}
}
