package handlers

import (
	"net/http"
	"time"

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
	var post models.Post

	// Pega o user_id do contexto (injetado pelo AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Converter o userID para UUID
	authorID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verifica se a requisição contém uma postagem válida
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Define o AuthorID como o user_id do usuário autenticado
	post.AuthorID = authorID // Certifique-se de que o `AuthorID` seja do tipo correto

	// Chama o serviço para criar o post
	if err := h.postService.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Responde com o post criado
	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPostById(c *gin.Context) {
	// Converte o ID da URL para UUID
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Chama o serviço com o UUID
	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	// Pega o ID do post do parâmetro da rota
	postIDStr := c.Param("id")

	// Converte o ID do post de string para uuid.UUID
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Busca o post existente no banco de dados
	existingPost, err := h.postService.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var updatedPostData models.Post
	if err := c.ShouldBindJSON(&updatedPostData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualiza apenas os campos que foram modificados
	existingPost.Title = updatedPostData.Title
	existingPost.Content = updatedPostData.Content
	existingPost.Topic = updatedPostData.Topic
	existingPost.UpdatedAt = time.Now() // Atualiza o campo UpdatedAt com a hora atual

	// Chama o serviço para salvar o post atualizado
	if err := h.postService.UpdatePost(existingPost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, existingPost)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	// Converte o ID da URL para UUID
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Chama o serviço para deletar o post usando UUID
	if err := h.postService.DeletePost(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.Status(http.StatusNoContent)
}
