package handlers

import (
	services "GoVersi/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FriendshipHandler struct {
	friendshipService *services.FriendshipService
}

func NewFriendshipHandler(service *services.FriendshipService) *FriendshipHandler {
	return &FriendshipHandler{friendshipService: service}
}

func (h *FriendshipHandler) SendFriendRequest(c *gin.Context) {
	var request struct {
		AddresseeID string `json:"addressee_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Pega o user_id do contexto (injetado pelo AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Converter o userID (requester) para UUID
	requesterUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requester ID"})
		return
	}

	// Converter o AddresseeID para UUID
	addresseeUUID, err := uuid.Parse(request.AddresseeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid addressee ID"})
		return
	}

	// Chama o serviço para enviar a solicitação de amizade
	if err := h.friendshipService.SendFriendRequest(addresseeUUID, requesterUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}

func (h *FriendshipHandler) AcceptFriendRequest(c *gin.Context) {
	friendshipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friendship ID"})
		return
	}

	if err := h.friendshipService.AcceptFriendRequest(friendshipID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted"})
}

func (h *FriendshipHandler) DeclineFriendRequest(c *gin.Context) {
	friendshipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friendship ID"})
		return
	}

	if err := h.friendshipService.DeclineFriendRequest(friendshipID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request declined"})
}
