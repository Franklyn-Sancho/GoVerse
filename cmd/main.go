package main

import (
	"GoVersi/internal/config"
	"GoVersi/internal/handlers"
	"GoVersi/internal/models"
	"GoVersi/internal/routes"
	services "GoVersi/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Carregar variáveis de ambiente
	loadEnv()

	// Conectar ao banco de dados
	db := connectDatabase()

	// Inicializar o serviço de usuários
	userService := services.NewUserService(db)

	// Configurar os handlers com o serviço de usuários
	handlers.SetUserService(userService)

	tokenBlacklistService := services.NewTokenBlacklistService(db)

	// Configurar os handlers com os serviços
	handlers.SetUserService(userService)
	handlers.SetTokenBlacklistService(tokenBlacklistService) // Adicione essa linha

	// Inicializar o router
	r := setupRouter()

	// Iniciar o servidor
	startServer(r)
}

// loadEnv carrega as variáveis de ambiente do arquivo .env
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using environment variables.")
	}
}

func connectDatabase() *gorm.DB {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Faz a migração automática para criar as tabelas necessárias
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = db.AutoMigrate(&models.TokenBlacklist{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// setupRouter configura e retorna o router com as rotas definidas
func setupRouter() *gin.Engine {
	r := gin.Default()
	routes.SetupRoutes(r)
	return r
}

// startServer inicia o servidor na porta especificada
func startServer(r *gin.Engine) {
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
