package services

import (
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type FriendshipService struct {
	repo *repository.FriendshipRepository
}

func NewFriendshipService(repo *repository.FriendshipRepository) *FriendshipService {
	return &FriendshipService{repo: repo}
}

func (s *FriendshipService) SendFriendRequest(requesterID, addresseeID uuid.UUID) error {
	existingFriendship, err := s.repo.GetFriendshipBetweenUsers(requesterID, addresseeID)
	if err == nil && existingFriendship != nil {
		return errors.New("friend request already exists or users are already friends")
	}

	return s.repo.SendFriendRequest(requesterID, addresseeID)
}

func (s *FriendshipService) AcceptFriendRequest(id uuid.UUID) error {
	friendship, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if friendship.Status != models.StatusPending {
		return errors.New("only pending requests can be accepted")
	}

	friendship.Status = models.StatusAccepted
	return s.repo.Update(friendship)
}

func (s *FriendshipService) DeclineFriendRequest(id uuid.UUID) error {
	friendship, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if friendship.Status != models.StatusPending {
		return errors.New("only pending requests can be declined")
	}

	friendship.Status = models.StatusDeclined
	return s.repo.Update(friendship)
}

func (s *FriendshipService) GetFriendsForUser(userID uuid.UUID) ([]models.Friendship, error) {
	return s.repo.GetFriendsForUser(userID)
}

func (s *FriendshipService) GetPendingRequestsForUser(userID uuid.UUID) ([]models.Friendship, error) {
	return s.repo.GetPendingRequestsForUser(userID)
}
