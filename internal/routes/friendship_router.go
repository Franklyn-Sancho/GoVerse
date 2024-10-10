package routes

import (
	"GoVersi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupFriendshipRoutes(r *gin.RouterGroup, handler *handlers.FriendshipHandler) {
	// Rotas para a funcionalidade de amizade
	friendshipGroup := r.Group("/friendship")

	friendshipGroup.POST("/send", handler.SendFriendRequest)           // Enviar solicitação de amizade
	friendshipGroup.POST("/accept/:id", handler.AcceptFriendRequest)   // Aceitar solicitação de amizade
	friendshipGroup.POST("/decline/:id", handler.DeclineFriendRequest) // Recusar solicitação de amizade
	/* friendshipGroup.GET("/friends", handler.ListFriends) */ // Listar amigos
}
