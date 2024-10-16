package routes

import (
	"GoVersi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupLikeRoutes(router *gin.RouterGroup, likeHandler *handlers.LikeHandler) {
	// Rotas para likes em Posts
	posts := router.Group("/posts/likes/")
	{
		posts.POST("/:post_id", likeHandler.LikePost)     // Curtir post
		posts.DELETE("/:post_id", likeHandler.UnlikePost) // Descurtir post
	}

	// Rotas para likes em Comentários
	comments := router.Group("/comments/likes")
	{
		comments.POST("/:comment_id", likeHandler.LikeComment)     // Curtir comentário
		comments.DELETE("/:comment_id", likeHandler.UnlikeComment) // Descurtir comentário
		comments.GET("/count", likeHandler.GetLikesCount)          // Contar likes no comentário
	}
}
