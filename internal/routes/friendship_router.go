package routes

import (
	"GoVersi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupFriendshipRoutes(r *gin.RouterGroup, handler *handlers.FriendshipHandler) {
	friendshipGroup := r.Group("/friendship")

	friendshipGroup.POST("/send", handler.SendFriendRequest)           // Send Friendship Request
	friendshipGroup.POST("/accept/:id", handler.AcceptFriendRequest)   // Accept Friendship Request
	friendshipGroup.POST("/decline/:id", handler.DeclineFriendRequest) // Decline Friendship Request
	/* friendshipGroup.GET("/friends", handler.ListFriends) */ // Listar amigos
}
