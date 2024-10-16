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

	// Chama a função para configurar as rotas
	SetupRoutes(r, postHandler, friendshipHandler, commentHandler, likeHandler)

	return r
}

// SetupRoutes agora também recebe um FriendshipHandler
func SetupRoutes(router *gin.Engine, postHandler *handlers.PostHandler, friendshipHandler *handlers.FriendshipHandler, commentHandler *handlers.CommentHandler, likeHandler *handlers.LikeHandler) {
	// Defina a chave secreta do JWT
	secretKey := os.Getenv("JWT_SECRET_KEY")
	log.Printf("SetupRoutes Secret Key: %s", secretKey)

	// Rotas públicas (não requerem autenticação)
	router.POST("/login", handlers.Login)
	router.POST("/register", handlers.RegisterUser)
	router.GET("/confirm-email", handlers.ConfirmEmail)

	// Rotas protegidas (requerem autenticação)
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(secretKey))

	// Configuração das rotas de usuários, postagens e amizade
	SetupUserRoutes(auth)
	SetupPostRoutes(auth, postHandler)             // Aqui passamos o postHandler
	SetupFriendshipRoutes(auth, friendshipHandler) // Registre o FriendshipHandler
	SetupCommentRoutes(auth, commentHandler)
	SetupLikeRoutes(auth, likeHandler)
}
