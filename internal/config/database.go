package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB configura e retorna uma conexão ao banco de dados
func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=123456789 dbname=goversidb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}
	log.Println("Database connection established")
	return db, nil
}
