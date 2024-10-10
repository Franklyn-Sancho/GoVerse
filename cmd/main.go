package main

import (
	"GoVersi/internal/config"
	"GoVersi/internal/handlers"
	"GoVersi/internal/middleware"
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
	routers "GoVersi/internal/routes"
	services "GoVersi/internal/service"
	"log"
	"os"

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

	// Inicializar os repositórios
	postRepository := repository.NewPostRepository(db)

	// Inicializar o serviço de postagens com o repositório
	postService := services.NewPostService(postRepository)

	// Inicializar o serviço de Token Blacklist
	tokenBlacklistService := services.NewTokenBlacklistService(db)

	// Configurar os handlers com os serviços
	handlers.SetUserService(userService)
	handlers.SetPostService(postService)                     // Configure o serviço de postagens
	handlers.SetTokenBlacklistService(tokenBlacklistService) // Configure o serviço de Token Blacklist

	// Criar uma instância do PostHandler
	postHandler := handlers.NewPostHandler(postService)

	// Inicializar o router
	r := setupRouter(postHandler) // Passa o postHandler

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

	return db
}

// SetupRoutes agora também recebe um PostHandler
func SetupRoutes(router *gin.Engine, postHandler *handlers.PostHandler) {
	// Defina a chave secreta do JWT
	secretKey := os.Getenv("JWT_SECRET_KEY")
	log.Printf("SetupRoutes Secret Key: %s", secretKey)

	// Rotas públicas (não requerem autenticação)
	router.POST("/login", handlers.Login)
	router.POST("/register", handlers.RegisterUser)

	// Rotas protegidas (requerem autenticação)
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(secretKey))

	// Configuração das rotas de usuários e postagens
	routers.SetupUserRoutes(auth)
	routers.SetupPostRoutes(auth, postHandler) // Aqui passamos o postHandler
}

// setupRouter configura e retorna o router com as rotas definidas
func setupRouter(postHandler *handlers.PostHandler) *gin.Engine {
	r := gin.Default()
	SetupRoutes(r, postHandler) // Passa o postHandler para SetupRoutes
	return r
}

// startServer inicia o servidor na porta especificada
func startServer(r *gin.Engine) {
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
