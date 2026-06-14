package config

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type InitialAdmin struct {
	Rut        string
	Usuario    string
	Contrasena string
	Rol        string
}

func (InitialAdmin) TableName() string {
	return "empleados"
}

func InitialSetup(db *gorm.DB) {
	var count int64

	db.Model(&InitialAdmin{}).Where("rol = ?", "Admin").Count(&count)

	if count == 0 {

		hashContrasena, err := bcrypt.GenerateFromPassword([]byte("Admin123."), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Error al generar la contraseña del admin inicial.")
		}

		admin := InitialAdmin{
			Rut:        "11111111-1",
			Usuario:    "admin",
			Contrasena: string(hashContrasena),
			Rol:        "Admin",
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatal("Error al crear el admin inicial: ", err)
		}

		log.Println("Administrador inicial creado con éxito.")
	}
}
