package routers

import (
	"GoVersi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupPostRoutes(router *gin.RouterGroup, postHandler *handlers.PostHandler) {
	posts := router.Group("/posts")
	{
		posts.POST("/create", postHandler.CreatePost)
		posts.GET("/:id", postHandler.GetPostById)
		posts.PUT("/:id", postHandler.UpdatePost)
		posts.DELETE("/:id", postHandler.DeletePost)
	}
}
