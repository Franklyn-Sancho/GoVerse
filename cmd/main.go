package main

import (
	"GoVersi/internal/config"
	"GoVersi/internal/handlers"
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
	"GoVersi/internal/routes"
	services "GoVersi/internal/service"
	"os"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Carregar variáveis de ambiente
	loadEnv()

	// database connect
	db := connectDatabase()

	// start repositories and services
	userRepository := repository.NewUserRepository(db)     // create a nre instance of repository
	userService := services.NewUserService(userRepository) // import repository to service

	postRepository := repository.NewPostRepository(db)
	tokenBlacklistService := services.NewTokenBlacklistService(db)
	friendshipRepository := repository.NewFriendshipRepository(db)
	commentRepository := repository.NewCommentRepository(db)
	likeRepository := repository.NewLikeRepository(db)

	postService := services.NewPostService(postRepository)
	friendshipService := services.NewFriendshipService(friendshipRepository)
	commentService := services.NewCommentService(commentRepository)
	likeService := services.NewLikeService(likeRepository)

	// Configure the handlers with the services
	handlers.SetUserService(userService)
	handlers.SetTokenBlacklistService(tokenBlacklistService) // Configure the Token Blacklist service

	// Create an instance of PostHandler
	postHandler := handlers.NewPostHandler(postService)

	// Initialize the FriendshipHandler
	friendshipHandler := handlers.NewFriendshipHandler(friendshipService)

	commentHandler := handlers.NewCommentHandler(commentService)

	likeHandler := handlers.NewLikeHandler(likeService)

	// Initialize the router
	r := routes.SetupRouter(postHandler, friendshipHandler, commentHandler, likeHandler)

	r.Static("/uploads", "./uploads")

	// Create the uploads directory if it doesn't exist
	os.MkdirAll("uploads/imageProfile", 0755)
	os.MkdirAll("uploads/images", 0755)
	os.MkdirAll("uploads/videos", 0755)

	// Start the server
	startServer(r)

}

// loadEnv load .env
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

	// automatic migrations
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

	err = db.AutoMigrate(&models.Like{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// start server on :8080 port
func startServer(r *gin.Engine) {
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
