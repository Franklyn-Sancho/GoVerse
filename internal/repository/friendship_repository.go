package repository

import (
	"GoVersi/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/* type FriendshipRepository interface {
	SendFriendRequest(requesterID, addresseeID uuid.UUID) error
	AcceptFriendRequest(id uuid.UUID) error
	DeclineFriendRequest(id uuid.UUID) error
	GetPendingRequestsForUser(userID uuid.UUID) ([]models.Friendship, error)
	GetFriendsForUser(userID uuid.UUID) ([]models.Friendship, error)
	GetFriendshipBetweenUsers(requesterID, addresseeID uuid.UUID) (*models.Friendship, error)
	FindByID(id uuid.UUID) (*models.Friendship, error)
	Update(friendship *models.Friendship) error
} */

type FriendshipRepository struct {
	db *gorm.DB
}

func NewFriendshipRepository(db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{db: db}
}

func (r *FriendshipRepository) SendFriendRequest(requesterID, addresseeID uuid.UUID) error {
	friendship := &models.Friendship{
		ID:          uuid.New(),
		RequesterID: requesterID,
		AddresseeID: addresseeID,
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Validação antes de criar
	if err := friendship.Validate(); err != nil {
		return err
	}

	return r.db.Create(friendship).Error
}

// get friendship between users
func (r *FriendshipRepository) GetFriendshipBetweenUsers(requesterID, addresseeID uuid.UUID) (*models.Friendship, error) {
	var friendship models.Friendship
	err := r.db.Where("(requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)",
		requesterID, addresseeID, addresseeID, requesterID).
		First(&friendship).Error

	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

// get friendship by id
func (r *FriendshipRepository) FindByID(id uuid.UUID) (*models.Friendship, error) {
	var friendship models.Friendship
	err := r.db.First(&friendship, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

// decline friendship request
func (r *FriendshipRepository) DeclineFriendRequest(id uuid.UUID) error {
	return r.db.Model(&models.Friendship{}).Where("id = ?", id).Update("status", "declined").Error
}

// get pendent requests for user
func (r *FriendshipRepository) GetPendingRequestsForUser(userID uuid.UUID) ([]models.Friendship, error) {
	var requests []models.Friendship
	err := r.db.Where("addressee_id = ? AND status = ?", userID, "pending").Find(&requests).Error
	return requests, err
}

// update friendship (accept or decline)
func (r *FriendshipRepository) Update(friendship *models.Friendship) error {
	return r.db.Save(friendship).Error
}

// get friends by user (accepted friendships)
func (r *FriendshipRepository) GetFriendsForUser(userID uuid.UUID) ([]models.Friendship, error) {
	var friends []models.Friendship
	err := r.db.Where("(requester_id = ? OR addressee_id = ?) AND status = ?", userID, userID, "accepted").Find(&friends).Error
	return friends, err
}
