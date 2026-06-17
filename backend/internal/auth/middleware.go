package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("auth")
		if err != nil {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado, no ha iniciado sesión."})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo de firma inesperado: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado, token inválido o expirado."})
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok {
			c.Set("id_empleado", claims.ID)
			c.Set("rol", claims.Rol)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado, datos de identidad corruptos."})
			return
		}
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		rolContexto, existe := c.Get("rol")
		if !existe {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No se pudo verificar la identidad del empleado."})
			return
		}

		rolEmpleado, ok := rolContexto.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error interno al procesar el rol del empleado."})
			return
		}

		permitido := false
		for _, rol := range roles {
			if rolEmpleado == rol {
				permitido = true
				break
			}
		}

		if !permitido {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado, no tienes los privilegios para realizar esta acción.",
			})
			return
		}

		c.Next()
	}
}
