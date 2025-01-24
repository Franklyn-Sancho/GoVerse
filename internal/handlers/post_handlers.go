package handlers

import (
	"net/http"

	"GoVersi/internal/models"
	services "GoVersi/internal/service"
	"GoVersi/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(service *services.PostService) *PostHandler {
	return &PostHandler{postService: service}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var request struct {
		Title   string `form:"title" json:"title" binding:"required"`
		Content string `form:"content" json:"content" binding:"required"`
		Topic   string `form:"topic" json:"topic" binding:"required"`
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authorID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	contentType := c.ContentType()
	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
			return
		}
	} else if contentType == "multipart/form-data" {
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form-data input"})
			return
		}
	} else {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported content type"})
		return
	}

	imageURL, err := utils.HandleImageUpload(c, "uploads/images")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading image: " + err.Error()})
		return
	}

	videoURL, err := utils.HandleVideoUpload(c, "uploads/videos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading video: " + err.Error()})
		return
	}

	post, err := h.postService.CreatePost(request.Title, request.Content, request.Topic, imageURL, videoURL, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating post: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPostById(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var updatedPostData models.Post
	if err := c.ShouldBindJSON(&updatedPostData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPost, err := h.postService.UpdatePost(postID, &updatedPostData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	if err := h.postService.DeletePost(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.Status(http.StatusNoContent)
}
