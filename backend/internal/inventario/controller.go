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

	//------------- validaciones--------------
	if nuevoProducto.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El stock no puede ser negativo"})
		return
	}
	if nuevoProducto.Precio <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El precio debe ser mayor a cero"})
		return
	}
	// ---------------------------------------

	// 2. Extrae la conexión a la base de datos
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	// 3. Llamado al repositorio
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

// UpdateProducto maneja la petición PUT para modificar el inventario
func UpdateProducto(c *gin.Context) {
	// 1. Obtener el UUID
	id := c.Param("id")

	// 2. Extraer la base de datos del contexto
	dbInstance, existe := c.Get("db")
	if !existe {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	db := dbInstance.(*gorm.DB)

	// 3.Buscar el producto existente en la BD primero
	var productoExistente Producto
	if err := db.First(&productoExistente, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "El producto que intentas actualizar no existe"})
		return
	}

	// 4. Leer el JSON del frontend con los datos nuevos
	var datosNuevos Producto
	if err := c.ShouldBindJSON(&datosNuevos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El formato de los datos es incorrecto"})
		return
	}

	if datosNuevos.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Operación rechazada: El stock no puede ser negativo"})
		return
	}
	if datosNuevos.Precio <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Operación rechazada: El precio debe ser mayor a cero"})
		return
	}

	// 5. Llamar al Repositorio
	if err := ActualizarProducto(db, &productoExistente, &datosNuevos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron guardar los cambios en la base de datos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":  "Producto actualizado correctamente",
		"producto": productoExistente,
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
