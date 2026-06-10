package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Inicializa la conexión a PostgreSQL
func Init() *gorm.DB {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Aviso: No se encontro el archivo .env")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" {
		log.Fatal("Error: Las variables de entorno de la base de datos no están configuradas")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error: No se pudo conectar a la base de datos \n", err)
	}

	fmt.Println("Conexión con la base de datos establecida")

	return db
}
