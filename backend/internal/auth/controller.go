package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type EmpleadoAuth struct {
	ID         string
	Nombre     string
	Usuario    string
	Contrasena string
	Rol        string
	Activo     bool
}

func (EmpleadoAuth) TableName() string {
	return "empleados"
}

type AuthController struct {
	db        *gorm.DB
	jwtSecret string
	domain    string
}

func NewAuthController(db *gorm.DB, secret string, domain string) *AuthController {
	return &AuthController{
		db:        db,
		jwtSecret: secret,
		domain:    domain,
	}
}

type LoginInput struct {
	Usuario    string `json:"usuario" binding:"required"`
	Contrasena string `json:"contrasena" binding:"required"`
}

type JWTClaims struct {
	ID     string `json:"id_empleado"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
	jwt.RegisteredClaims
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe enviar usuario y contraseña."})
		return
	}

	var empleado EmpleadoAuth

	if err := ctrl.db.Where("usuario = ?", input.Usuario).First(&empleado).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas."})
		return
	}

	if !empleado.Activo {
		c.JSON(http.StatusForbidden, gin.H{"error": "Esta cuenta de empleado esta desactivada."})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(empleado.Contrasena), []byte(input.Contrasena))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas."})
		return
	}

	tiempoExpiracion := time.Now().Add(24 * time.Hour)
	claims := JWTClaims{
		ID:     empleado.ID,
		Nombre: empleado.Nombre,
		Rol:    empleado.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tiempoExpiracion),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(ctrl.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al generar el token de autorización."})
		return
	}

	c.SetCookie(
		"auth",
		tokenString,
		int(24*time.Hour.Seconds()),
		"/",
		ctrl.domain,
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Sesión iniciada exitosamente.",
		"empleado": gin.H{
			"id_empleado": empleado.ID,
			"usuario":     empleado.Usuario,
			"rol":         empleado.Rol,
		},
	})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	c.SetCookie(
		"auth",
		"",
		-1,
		"/",
		ctrl.domain,
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Sesión cerrada exitosamente.",
	})
}

type EmpleadoResponse struct {
	ID      string `json:"id_empleado"`
	Usuario string `json:"usuario"`
	Rol     string `json:"rol"`
}

func (ctrl *AuthController) Me(c *gin.Context) {
	ID, exists := c.Get("id_empleado")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesión no encontrada."})
		return
	}

	var sesion EmpleadoResponse

	err := ctrl.db.Model(&EmpleadoAuth{}).
		Select("id, usuario, rol").
		Where("id = ?", ID).
		First(&sesion).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "El empleado ya no existe en el sistema."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener los datos del empleado."})
		return
	}

	c.JSON(http.StatusOK, sesion)
}
