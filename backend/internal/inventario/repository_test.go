package inventario

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=admin123 dbname=pos_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Fallo conexión a BD de pruebas: %v", err)
	}
	db.AutoMigrate(&Categoria{}, &Producto{})
	return db
}

func TestObtenerProductoPorCodigo_NoEncontrado(t *testing.T) {
	db := setupTestDB(t)

	codigoInexistente := "BARCODE-0000X"
	producto, err := ObtenerProductoPorCodigo(db, codigoInexistente)

	// El repositorio debería fallar arrojando ErrRecordNotFound
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, producto)
}
