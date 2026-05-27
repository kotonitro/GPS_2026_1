package inventario

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CrearProducto(c *gin.Context) {
	var nuevoProducto Producto

	// 1. Lee el JSON
	if err := c.ShouldBindJSON(&nuevoProducto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos"})
		return
	}

	// 2. Extrae la conexión a la base de datos del Middleware
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	// 3. LLamado al repositorio
	err := GuardarProducto(db, &nuevoProducto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo guardar el producto"})
		return
	}

	// 4. Entrega mensaje
	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Producto agregado",
		"producto": nuevoProducto,
	})
}

// GetProductos maneja la solicitud para obtener todos los productos.
func GetProductos(c *gin.Context) {
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	productos, err := ObtenerTodosProductos(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al obtener los productos"})
		return
	}

	c.JSON(http.StatusOK, productos)
}

// GetProductoByID maneja la solicitud para obtener un producto por su ID.
func GetProductoByID(c *gin.Context) {
	id := c.Param("id")

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	producto, err := ObtenerProductoPorID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Producto no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al buscar el producto"})
		return
	}

	c.JSON(http.StatusOK, producto)
}

// UpdateProducto maneja la solicitud para actualizar un producto.
func UpdateProducto(c *gin.Context) {
	id := c.Param("id")

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	// Verificar si el producto existe
	productoExistente, err := ObtenerProductoPorID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "El producto que intenta actualizar no existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al buscar el producto"})
		return
	}

	// Leer los nuevos datos
	var datosNuevos Producto
	if err := c.ShouldBindJSON(&datosNuevos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos"})
		return
	}

	// Llamar al repositorio para actualizar
	if err := ActualizarProducto(db, productoExistente, &datosNuevos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo actualizar el producto"})
		return
	}

	// Volver a buscar el producto para devolver el estado más reciente con la categoría actualizada.
	productoActualizado, _ := ObtenerProductoPorID(db, id)

	c.JSON(http.StatusOK, gin.H{
		"mensaje":  "Producto actualizado exitosamente",
		"producto": productoActualizado,
	})
}

// DeleteProducto maneja la solicitud para eliminar un producto.
func DeleteProducto(c *gin.Context) {
	id := c.Param("id")

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	if err := EliminarProducto(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo eliminar el producto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Producto eliminado exitosamente"})
}
