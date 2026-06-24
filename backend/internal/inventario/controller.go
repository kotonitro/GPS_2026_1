package inventario

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 1. Creamos el Struct y el Constructor
type InventarioController struct {
	db *gorm.DB
}

func NewInventarioController(db *gorm.DB) *InventarioController {
	return &InventarioController{db: db}
}

// 2. Convertimos tus funciones en métodos (ctrl *InventarioController)
func (ctrl *InventarioController) CrearProducto(c *gin.Context) {
	var nuevoProducto Producto

	if err := c.ShouldBindJSON(&nuevoProducto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos"})
		return
	}

	if nuevoProducto.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El stock no puede ser negativo"})
		return
	}
	if nuevoProducto.Precio <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El precio debe ser mayor a cero"})
		return
	}

	// ¡Usamos ctrl.db en lugar de c.Get("db")!
	err := GuardarProducto(ctrl.db, &nuevoProducto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo guardar el producto"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Producto agregado",
		"producto": nuevoProducto,
	})
}

func (ctrl *InventarioController) GetProductos(c *gin.Context) {
	productos, err := ObtenerTodosProductos(ctrl.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al obtener los productos"})
		return
	}
	c.JSON(http.StatusOK, productos)
}

func (ctrl *InventarioController) GetProductoByID(c *gin.Context) {
	id := c.Param("id")

	producto, err := ObtenerProductoPorID(ctrl.db, id)
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

func (ctrl *InventarioController) GetProductoByCodigo(c *gin.Context) {
	codigo := c.Param("codigo")

	producto, err := ObtenerProductoPorCodigo(ctrl.db, codigo)
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

func (ctrl *InventarioController) UpdateProducto(c *gin.Context) {
	id := c.Param("id")

	var productoExistente Producto
	if err := ctrl.db.First(&productoExistente, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "El producto que intentas actualizar no existe"})
		return
	}

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

	if err := ActualizarProducto(ctrl.db, &productoExistente, &datosNuevos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron guardar los cambios en la base de datos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":  "Producto actualizado correctamente",
		"producto": productoExistente,
	})
}

func (ctrl *InventarioController) DeleteProducto(c *gin.Context) {
	id := c.Param("id")

	if err := EliminarProducto(ctrl.db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo eliminar el producto"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "Producto eliminado exitosamente"})
}

// ---- CATEGORIAS ----

func (ctrl *InventarioController) CrearCategoria(c *gin.Context) {
	var nuevaCategoria Categoria

	if err := c.ShouldBindJSON(&nuevaCategoria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos"})
		return
	}

	err := GuardarCategoria(ctrl.db, &nuevaCategoria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo guardar la categoría"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":   "Categoría agregada",
		"categoria": nuevaCategoria,
	})
}

func (ctrl *InventarioController) GetCategorias(c *gin.Context) {
	categorias, err := ObtenerTodasCategorias(ctrl.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al obtener las categorías"})
		return
	}
	c.JSON(http.StatusOK, categorias)
}

func (ctrl *InventarioController) UpdateCategoria(c *gin.Context) {
	id := c.Param("id")

	var datosNuevos Categoria
	if err := c.ShouldBindJSON(&datosNuevos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos de categoría inválidos"})
		return
	}

	if err := ActualizarCategoria(ctrl.db, id, &datosNuevos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo actualizar la categoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Categoría actualizada exitosamente"})
}

func (ctrl *InventarioController) DeleteCategoria(c *gin.Context) {
	id := c.Param("id")

	err := EliminarCategoria(ctrl.db, id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"Error": "No puedes eliminar esta categoría porque hay productos que la están usando. Reasigna los productos primero.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Categoría eliminada exitosamente"})
}
