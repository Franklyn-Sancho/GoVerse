package handlers

import (
	"GoVersi/internal/models"
	services "GoVersi/internal/service"
	"GoVersi/internal/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	var user models.User

	// Verifica se os dados do usuário foram enviados corretamente
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user input"})
		return
	}

	// Tenta pegar o arquivo de imagem (caso exista)
	file, err := c.FormFile("image")
	if err == nil {
		// Verificar e criar a pasta uploads/imageProfile se ela não existir
		uploadPath := "uploads/imageProfile"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// Salvar o arquivo no diretório
		imagePath := filepath.Join(uploadPath, file.Filename)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}
		user.ImageProfile = "/" + imagePath // Caminho da imagem salvo no campo ImageProfile
	} else if err != http.ErrMissingFile {
		// Se houve um erro diferente de arquivo ausente, retorna erro
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error handling image upload"})
		return
	}

	// Chama o serviço para registrar o usuário
	if err := userService.RegisterUser(&user); err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	// Retorna o usuário criado
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
