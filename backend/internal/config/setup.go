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

type InitialMetodoPago struct {
	ID           string
	NombreMetodo string
}

func (InitialMetodoPago) TableName() string {
	return "metodo_pagos"
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
	metodosBasicos := []InitialMetodoPago{
		{ID: "11111111-1111-1111-1111-111111111111", NombreMetodo: "Efectivo"},
		{ID: "22222222-2222-2222-2222-222222222222", NombreMetodo: "Tarjeta"},
		{ID: "33333333-3333-3333-3333-333333333333", NombreMetodo: "Fiado"},
	}

	for _, metodo := range metodosBasicos {
		var count int64
		db.Model(&InitialMetodoPago{}).Where("id = ?", metodo.ID).Count(&count)
		if count == 0 {
			if err := db.Create(&metodo).Error; err != nil {
				log.Println("Error al crear método de pago base: ", err)
			} else {
				log.Printf("Método de pago base '%s' creado con éxito.\n", metodo.NombreMetodo)
			}
		}
	}
}
