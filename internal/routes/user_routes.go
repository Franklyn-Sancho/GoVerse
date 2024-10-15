package routes

import (
	"GoVersi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/:id", handlers.GetUserById)
		users.GET("/username/:username", handlers.GetUserByUsername)
		users.GET("/email/:email", handlers.GetUserByEmail)
		users.DELETE("/:id", handlers.DeleteUser)

		// Outras rotas...
		router.GET("/confirm-email", handlers.ConfirmEmail)

		// Novas funcionalidades relacionadas ao usuário
		users.PATCH("/:id/suspend", handlers.SuspendUser)
		users.POST("/:id/request-deletion", handlers.RequestAccountDeletion)
		users.DELETE("/:id/permanently-delete", handlers.PermanentlyDeleteUser)

		// Logout
		users.POST("/logout", handlers.Logout)
	}
}
