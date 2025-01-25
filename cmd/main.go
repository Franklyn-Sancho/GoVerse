package main

import (
	"GoVersi/internal/config"
	"GoVersi/internal/handlers"
	"GoVersi/internal/infrastrucuture/queue"
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
	"GoVersi/internal/routes"
	services "GoVersi/internal/service"
	"GoVersi/internal/service/email"
	"encoding/json"

	/* "encoding/json" */
	"os"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	rabbitMQ, err := queue.NewRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQ.Close()

	// Start processing RabbitMQ messages in a separate goroutine
	go processRabbitMQMessages(rabbitMQ)

	// Load environment variables
	loadEnv()

	// Database connection
	db := connectDatabase()

	// Initialize repositories and services
	queueService := email.NewEmailQueueService(rabbitMQ)
	mailService := email.NewEmailService(queueService)

	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository, mailService)

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
	handlers.SetTokenBlacklistService(tokenBlacklistService)

	postHandler := handlers.NewPostHandler(postService)
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

// Function to process RabbitMQ messages
func processRabbitMQMessages(rabbitMQ *queue.RabbitMQ) {
	msgs, err := rabbitMQ.Consume("email_queue")
	if err != nil {
		log.Printf("Error consuming RabbitMQ messages: %v", err)
		return
	}

	emailService := email.NewEmailService(email.NewEmailQueueService(rabbitMQ))

	for d := range msgs {
		var msg email.EmailMessage
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		if err := emailService.SendEmail(msg.To, msg.Subject, msg.Body); err != nil {
			log.Printf("Error sending email: %v", err)
		} else {
			log.Printf("Email sent to %s successfully", msg.To)
		}
	}
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

func startServer(r *gin.Engine) {
	log.Println("Starting server on port :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
