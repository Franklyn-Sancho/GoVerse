package routes

import (
	"GoVersi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupLikeRoutes(router *gin.RouterGroup, likeHandler *handlers.LikeHandler) {
	posts := router.Group("/posts/likes/")
	{
		posts.POST("/:post_id", likeHandler.LikePost)     // Like post
		posts.DELETE("/:post_id", likeHandler.UnlikePost) // Deslike post
	}

	// Rotas para likes em Comentários
	comments := router.Group("/comments/likes")
	{
		comments.POST("/:comment_id", likeHandler.LikeComment)     // Like comment
		comments.DELETE("/:comment_id", likeHandler.UnlikeComment) // Deslike comment
		comments.GET("/count", likeHandler.GetLikesCount)          // count likes on comment
	}
}
