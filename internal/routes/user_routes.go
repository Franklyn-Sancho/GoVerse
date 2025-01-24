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

		users.PATCH("/:id/suspend", handlers.SuspendUser)
		users.POST("/:id/request-deletion", handlers.RequestAccountDeletion)
		users.DELETE("/:id/permanently-delete", handlers.PermanentlyDeleteUser)

		users.POST("/logout", handlers.Logout)
	}
}
