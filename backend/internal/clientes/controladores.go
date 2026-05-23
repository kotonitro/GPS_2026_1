package clientes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gonic/gin"
	"gorm.io/gorm"
)

func CreateCliente(c *gin.Context) {
	var nuevoCliente Cliente

	//ve los atributos
	if err := c.ShouldBindJSON(&nuevoCliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos o campos obligatorios faltantes"})
		return
	}

	//extraemos la conexión a la base de datos 
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	//lo importante es validar el formato del rut(chileno)
	nuevoCliente.Rut = strings.TrimSpace(nuevoCliente.Rut)
	if len(nuevoCliente.Rut) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El formato del RUT no es válido"})
		return
	}

	//verificamos que el rut no esté en la bd
	existe, err := CheckRutExiste(db, nuevoCliente.Rut)
	if err != nil {
		//error general
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al verificar la duplicidad del RUT"})
		return
	}
	if existe {
		c.JSON(http.StatusConflict, gin.H{"Error": "El RUT ingresado ya se encuentra registrado"})
		return
	}

	//guardar
	err = GuardarCliente(db, &nuevoCliente)
	if err != nil {
		//error general
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo guardar el cliente en la base de datos"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Cliente agregado exitosamente",
		"cliente": nuevoCliente,
	})
}


func GetClientes(c *gin.Context) {
	
	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	listaClientes, err := ObtenerTodosLosClientes(db)
	if err != nil {
		//error general
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al obtener la lista de clientes"})
		return
	}

	c.JSON(http.StatusOK, listaClientes)
}


func GetClienteByID(c *gin.Context) {
	//primeramente vemos el id
	id := c.Param("id")

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	
	cliente, err := ObtenerClientePorID(db, id)
	if err != nil {
		
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "El cliente solicitado no existe"})
			return
		}
		//error general
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al buscar el cliente en la base de datos"})
		return
	}

	c.JSON(http.StatusOK, cliente)
}


func UpdateCliente(c *gin.Context) {
	//ver el id
	id := c.Param("id")

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	//verificamos que exista realmente
	clienteExistente, err := ObtenerClientePorID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "El cliente que intenta actualizar no existe"})
			return
		}
		//error general
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al buscar el cliente en la base de datos"})
		return
	}

	//leemos y validamos los nuevos datos 
	var datosNuevos Cliente
	if err := c.ShouldBindJSON(&datosNuevos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos inválidos o campos obligatorios faltantes"})
		return
	}

	//formato del rut
	if datosNuevos.Rut != "" && len(datosNuevos.Rut) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "El formato del nuevo RUT no es válido"})
		return
	}

	//unico con su rut
	if datosNuevos.Rut != "" && datosNuevos.Rut != clienteExistente.Rut {
		existe, err := CheckRutExiste(db, datosNuevos.Rut)
		if err != nil {
			//error general
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al verificar la duplicidad del RUT"})
			return
		}
		if existe {
			c.JSON(http.StatusConflict, gin.H{"Error": "El nuevo RUT ingresado ya pertenece a otro cliente"})
			return
		}
	}

	err = ActualizarCliente(db, clienteExistente, &datosNuevos)
	if err != nil {
		//error general
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudieron actualizar los datos del cliente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Cliente actualizado exitosamente",
		"cliente": clienteExistente,
	})
}


func DeleteCliente(c *gin.Context) {
	//vemos el id
	id := c.Param("id")

	dbInstance, _ := c.Get("db")
	db := dbInstance.(*gorm.DB)

	//verificamos que exista realmente
	_, err := ObtenerClientePorID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Error": "El cliente que intenta eliminar no existe"})
			return
		}
		//error general
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error al verificar el cliente en la base de datos"})
		return
	}

	err = EliminarCliente(db, id)
	if err != nil {
		//error general
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "No se pudo eliminar el cliente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Cliente eliminado exitosamente",
	})
}