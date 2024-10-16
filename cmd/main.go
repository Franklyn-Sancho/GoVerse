package main

import (
	"GoVersi/internal/config"
	"GoVersi/internal/handlers"
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
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

	// Inicializar os repositórios
	userRepository := repository.NewUserRepository(db)     // Crie uma instância do repositório
	userService := services.NewUserService(userRepository) // Passe o repositório para o UserService

	postRepository := repository.NewPostRepository(db)
	tokenBlacklistService := services.NewTokenBlacklistService(db)
	friendshipRepository := repository.NewFriendshipRepository(db)
	commentRepository := repository.NewCommentRepository(db)

	postService := services.NewPostService(postRepository)
	friendshipService := services.NewFriendshipService(friendshipRepository)
	commentService := services.NewCommentService(commentRepository)

	// Configurar os handlers com os serviços
	handlers.SetUserService(userService)
	handlers.SetTokenBlacklistService(tokenBlacklistService) // Configure o serviço de Token Blacklist

	// Criar uma instância do PostHandler
	postHandler := handlers.NewPostHandler(postService)

	// Inicializar o FriendshipHandler
	friendshipHandler := handlers.NewFriendshipHandler(friendshipService)

	commentHandler := handlers.NewCommentHandler(commentService)

	// Inicializar o router
	r := routes.SetupRouter(postHandler, friendshipHandler, commentHandler)

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

	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = db.AutoMigrate(&models.Friendship{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = db.AutoMigrate(&models.Comment{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// startServer inicia o servidor na porta especificada
func startServer(r *gin.Engine) {
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
