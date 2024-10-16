package handlers

import (
	"GoVersi/internal/models"
	services "GoVersi/internal/service"
	"GoVersi/internal/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Dependência de serviço de usuário
var userService *services.UserService
var tokenBlacklistService *services.TokenBlacklistService

// Função para configurar o serviço de usuário
func SetUserService(svc *services.UserService) {
	userService = svc
}

func SetTokenBlacklistService(s *services.TokenBlacklistService) {
	tokenBlacklistService = s
}

func RegisterUser(c *gin.Context) {
	var request struct {
		Username string `form:"username" binding:"required"`
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	// Verifica se a requisição contém dados válidos
	if err := c.ShouldBind(&request); err != nil {
		log.Printf("Erro de validação: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user input"})
		return
	}

	// Chama a função de upload de imagem
	imageURL, err := utils.HandleImageUpload(c, "uploads/imageProfile")
	if err != nil {
		log.Printf("Erro ao fazer upload da imagem: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cria o usuário com os dados capturados
	user := &models.User{
		Username:     request.Username,
		Email:        request.Email,
		Password:     request.Password,
		ImageProfile: imageURL,
		IsActive:     true,
	}

	// Chama o serviço para registrar o usuário
	if err := userService.RegisterUser(user); err != nil {
		log.Printf("Erro ao registrar usuário: %v", err)
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	log.Printf("Usuário '%s' registrado com sucesso", user.Username)
	c.JSON(http.StatusCreated, user)
}

func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := userService.LoginUser(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	// Pega o token do cabeçalho da requisição
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	// Remove "Bearer " do tokenString se presente
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	log.Printf("Received token: %s", tokenString) // Log do token recebido

	// Carrega a chave secreta das variáveis de ambiente
	secretKey := os.Getenv("JWT_SECRET_KEY")
	log.Printf("Secret Key during logout: %s", secretKey)

	// Tenta analisar e verificar as claims do token
	claims, err := utils.ParseTokenClaims(tokenString, secretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Token valid for user ID: %s", claims.UserID)
	expirationTime := claims.ExpiresAt.Time

	// Adiciona o token à blacklist com a data de expiração
	err = tokenBlacklistService.AddToTokenBlacklist(tokenString, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// Handler para buscar usuário pelo ID
func GetUserById(c *gin.Context) {
	id := c.Param("id")

	user, err := userService.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Handler para buscar usuário pelo username
func GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := userService.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Handler para buscar usuário pelo email
func GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := userService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Handler para deletar usuário
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := userService.DeleteUser(id); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Handler para suspender usuário
func SuspendUser(c *gin.Context) {
	id := c.Param("id")

	if err := userService.SuspendUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User suspended successfully"})
}

// Handler para solicitar a exclusão da conta
func RequestAccountDeletion(c *gin.Context) {
	id := c.Param("id")

	if err := userService.RequestAccountDeletion(id); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to request account deletion"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deletion requested successfully"})
}

// Handler para excluir usuário permanentemente
func PermanentlyDeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := userService.PermanentlyDeleteUser(id); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User will be deleted in 30 days"})
}

// ConfirmEmail manipula a confirmação de e-mail
func ConfirmEmail(c *gin.Context) {
	token := c.Query("token")

	// Encontre o usuário pelo token
	user, err := userService.FindByEmailConfirmToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado ou token inválido"})
		return
	}

	// Atualize o status de verificação do e-mail
	user.IsEmailVerified = true
	user.EmailConfirmToken = "" // Limpa o token após a confirmação
	if err := userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao verificar e-mail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "E-mail confirmado com sucesso"})
}
