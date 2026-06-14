package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
}

func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Aviso: No se encontro el archivo .env")
	}

	config := &AppConfig{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	if config.DBHost == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" || config.DBPort == "" {
		log.Fatal("Error: Faltan variables de entorno para la base de datos")
	}

	if config.JWTSecret == "" {
		log.Fatal("Error: JWT_SECRET no está configurado")
	}

	return config
}
