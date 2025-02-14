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
		Username string `json:"username" form:"username" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("Validation Error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user input"})
			return
		}
	}

	imageURL, err := utils.HandleImageUpload(c, "uploads/imageProfile")
	if err != nil {
		log.Printf("Upload image error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Username:     request.Username,
		Email:        request.Email,
		Password:     request.Password,
		ImageProfile: imageURL,
		IsActive:     true,
	}

	if err := userService.RegisterUser(user); err != nil {
		log.Printf("register user error: %v", err)
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	log.Printf("User '%s' registered successfully", user.Username)
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
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	log.Printf("Received token: %s", tokenString)

	secretKey := os.Getenv("JWT_SECRET_KEY")
	log.Printf("Secret Key during logout: %s", secretKey)

	claims, err := utils.ParseTokenClaims(tokenString, secretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Token valid for user ID: %s", claims.UserID)
	expirationTime := claims.ExpiresAt.Time

	err = tokenBlacklistService.AddToTokenBlacklist(tokenString, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")

	user, err := userService.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Handler get user by username
func GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := userService.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Handler get user by email
func GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := userService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Handler delete user
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

// Handler suspend user account
func SuspendUser(c *gin.Context) {
	id := c.Param("id")

	if err := userService.SuspendUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User suspended successfully"})
}

// Handler delete account user solicitation
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

// Handler delete user account permanently
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

func ConfirmEmail(c *gin.Context) {
	token := c.Query("token")

	// get user by token
	user, err := userService.FindByEmailConfirmToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found or invalid token"})
		return
	}

	//
	user.IsEmailVerified = true
	user.EmailConfirmToken = ""
	if err := userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "verify email error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email confirmed successfully"})
}
