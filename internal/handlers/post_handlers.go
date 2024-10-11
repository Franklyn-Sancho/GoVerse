package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"GoVersi/internal/models"
	services "GoVersi/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var postService *services.PostService

type PostHandler struct {
	postService *services.PostService
}

func SetPostService(service *services.PostService) {
	postService = service
}

func NewPostHandler(service *services.PostService) *PostHandler {
	return &PostHandler{postService: service}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var request struct {
		Title    string `form:"title" binding:"required"`
		Content  string `form:"content" binding:"required"`
		Topic    string `form:"topic" binding:"required"`
		ImageURL string `form:"image_url"`
	}

	// Pega o user_id do contexto (injetado pelo AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Converte o userID para UUID
	authorID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verifica se a requisição contém uma postagem válida
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Lidar com o upload da imagem
	file, err := c.FormFile("image")
	var imageURL string

	if err == nil {
		// Verificar e criar a pasta uploads se ela não existir
		uploadPath := "uploads"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// Salvar o arquivo
		imagePath := filepath.Join(uploadPath, file.Filename)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload"})
			return
		}
		imageURL = "/" + imagePath // Caminho da imagem
	}

	// Chama o serviço para criar o post
	post, err := h.postService.CreatePost(request.Title, request.Content, request.Topic, imageURL, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	// Chama o serviço para atualizar o post
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
