package routes

import (
	"GoVersi/internal/handlers"
	"GoVersi/internal/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// setupRouter inicializa as rotas da aplicação
func SetupRouter(postHandler *handlers.PostHandler, friendshipHandler *handlers.FriendshipHandler, commentHandler *handlers.CommentHandler, likeHandler *handlers.LikeHandler) *gin.Engine {
	r := gin.Default()

	SetupRoutes(r, postHandler, friendshipHandler, commentHandler, likeHandler)

	return r
}

// SetupRoutes agora também recebe um FriendshipHandler
func SetupRoutes(router *gin.Engine, postHandler *handlers.PostHandler, friendshipHandler *handlers.FriendshipHandler, commentHandler *handlers.CommentHandler, likeHandler *handlers.LikeHandler) {
	// secret key
	secretKey := os.Getenv("JWT_SECRET_KEY")
	log.Printf("SetupRoutes Secret Key: %s", secretKey)

	// public routes (authentication not required)
	router.POST("/login", handlers.Login)
	router.POST("/register", handlers.RegisterUser)
	router.GET("/confirm-email", handlers.ConfirmEmail)

	// protected routes (authentication required)
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(secretKey))

	// user routes
	SetupUserRoutes(auth)
	SetupPostRoutes(auth, postHandler)
	SetupFriendshipRoutes(auth, friendshipHandler)
	SetupCommentRoutes(auth, commentHandler)
	SetupLikeRoutes(auth, likeHandler)
}
