package routes

import (
	"GoVersi/internal/handlers"
	"GoVersi/internal/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Defina a chave secreta como uma variável
	secretKey := os.Getenv("JWT_SECRET_KEY")
	log.Printf("SetupRoutes Secret Key: %s", secretKey) // Adicione este log

	// Rotas públicas (não requerem autenticação)
	router.POST("/login", handlers.Login)
	router.POST("/users", handlers.RegisterUser)

	// Rotas protegidas (requerem autenticação)
	auth := router.Group("/users")
	auth.Use(middleware.AuthMiddleware(secretKey))
	{
		// Rotas de usuários
		auth.GET("/:id", handlers.GetUserById)
		auth.GET("/username/:username", handlers.GetUserByUsername)
		auth.GET("/email/:email", handlers.GetUserByEmail)
		auth.DELETE("/:id", handlers.DeleteUser)

		// Novas funcionalidades relacionadas ao usuário
		auth.PATCH("/:id/suspend", handlers.SuspendUser)                       // Suspender usuário
		auth.POST("/:id/request-deletion", handlers.RequestAccountDeletion)    // Solicitar exclusão da conta
		auth.DELETE("/:id/permanently-delete", handlers.PermanentlyDeleteUser) // Excluir permanentemente

		// Logout
		auth.POST("/logout", handlers.Logout)
	}
}
