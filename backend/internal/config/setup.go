package config

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type InitialRol struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Nombre      string
	Descripcion string
	EsAdmin     bool
}

func (InitialRol) TableName() string {
	return "roles"
}

type InitialAdmin struct {
	Rut        string
	Nombre     string
	Usuario    string
	Contrasena string
	RolID      string
	Telefono	string
}

func (InitialAdmin) TableName() string {
	return "empleados"
}

func InitialSetup(db *gorm.DB) {
	var rolAdmin InitialRol

	errRol := db.Where(InitialRol{Nombre: "Admin"}).Attrs(InitialRol{
		Descripcion: "Administrador principal del sistema.",
		EsAdmin:     true,
	}).FirstOrCreate(&rolAdmin).Error

	if errRol != nil {
		log.Fatal("Error al inicializar el rol de Admin: ", errRol)
	}

	var rolEmpleado InitialRol

	errRolEmpleado := db.Where(InitialRol{Nombre: "Empleado"}).Attrs(InitialRol{
		Descripcion: "Empleado regular del sistema.",
		EsAdmin:     false,
	}).FirstOrCreate(&rolEmpleado).Error

	if errRolEmpleado != nil {
		log.Println("Advertencia: No se pudo inicializar el rol de Empleado: ", errRolEmpleado)
	}

	var count int64
	db.Model(&InitialAdmin{}).Where("rol_id = ?", rolAdmin.ID).Count(&count)

	if count == 0 {
		hashContrasena, err := bcrypt.GenerateFromPassword([]byte("Admin123."), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Error al generar la contraseña del admin inicial.")
		}

		admin := InitialAdmin{
			Rut:        "11111111-1",
			Nombre:     "Administrador",
			Usuario:    "admin",
			Contrasena: string(hashContrasena),
			RolID:      rolAdmin.ID,
			Telefono:   "999999999",
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatal("Error al crear el admin inicial: ", err)
		}

		log.Println("Administrador y rol inicial creados con éxito.")
	}
}
