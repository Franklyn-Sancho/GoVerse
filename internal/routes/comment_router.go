package routes

import (
	"GoVersi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupCommentRoutes(router *gin.RouterGroup, commentHandler *handlers.CommentHandler) {
	posts := router.Group("/posts/comments")
	{
		posts.POST("/:post_id/create", commentHandler.CreateComment)
		posts.GET("/:comment_id", commentHandler.GetCommentById)
		posts.PUT("/:comment_id", commentHandler.UpdateComment)
		posts.DELETE("/:comment_id", commentHandler.DeleteComment)
	}
}
