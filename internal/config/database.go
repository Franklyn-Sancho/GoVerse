package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB configura e retorna uma conexão ao banco de dados
func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}
	log.Println("Database connection established")
	return db, nil
}
