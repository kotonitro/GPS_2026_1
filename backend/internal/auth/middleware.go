package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (ctrl *AuthController) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado, debes iniciar sesión."})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo de firma inesperado: %v", token.Header["alg"])
			}
			return []byte(ctrl.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado, token inválido o expirado."})
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok {
			c.Set("id_empleado", claims.ID)

			var estado struct {
				Activo  bool
				Nombre  string
				EsAdmin bool
			}

			if err := ctrl.db.Table("empleados").
				Select("empleados.activo, roles.nombre, roles.es_admin").
				Joins("JOIN roles ON empleados.rol_id = roles.id").
				Where("empleados.id = ?", claims.ID).
				First(&estado).Error; err != nil {

				c.SetCookie("auth", "", -1, "/", ctrl.domain, false, true)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Esta cuenta de empleado ya no existe en el sistema."})
				return
			}

			if !estado.Activo {
				c.SetCookie("auth", "", -1, "/", ctrl.domain, false, true)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Esta cuenta de empleado está desactivada."})
				return
			}

			c.Set("rol", estado.Nombre)
			c.Set("es_admin", estado.EsAdmin)

		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado, datos de identidad corruptos."})
			return
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		esAdminCtx, existe := c.Get("es_admin")
		if !existe {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No se pudo verificar el nivel de privilegios."})
			return
		}

		esAdmin, ok := esAdminCtx.(bool)
		if !ok || !esAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado, se requieren privilegios de administrador para realizar esta acción.",
			})
			return
		}

		c.Next()
	}
}
