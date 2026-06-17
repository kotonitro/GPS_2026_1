package inventario

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCrearProducto_ErrorPrecioNegativo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	ctrl := NewInventarioController(nil) // No tocamos la base de datos
	router.POST("/inventario/productos", ctrl.CrearProducto)

	// JSON con precio inválido
	cuerpoJSON := []byte(`{
		"nombre": "Monitor LED",
		"precio": -15000,
		"stock": 10
	}`)

	req, _ := http.NewRequest(http.MethodPost, "/inventario/productos", bytes.NewBuffer(cuerpoJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verificamos que el servidor detecte la regla de negocio y rechace (Error 400)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "mayor a cero")
}
