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

// Enviar solicitação de amizade
func (s *FriendshipService) SendFriendRequest(requesterID, addresseeID uuid.UUID) error {
	// Verificar se já existe uma amizade ou solicitação entre os usuários
	existingFriendship, err := s.repo.GetFriendshipBetweenUsers(requesterID, addresseeID)
	if err == nil && existingFriendship != nil {
		return errors.New("friend request already exists or users are already friends")
	}

	// Criar nova solicitação de amizade com status "pending"
	return s.repo.SendFriendRequest(requesterID, addresseeID)
}

// Aceitar solicitação de amizade
func (s *FriendshipService) AcceptFriendRequest(id uuid.UUID) error {
	friendship, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Garantir que a solicitação esteja "pending" antes de aceitar
	if friendship.Status != models.StatusPending {
		return errors.New("only pending requests can be accepted")
	}

	friendship.Status = models.StatusAccepted
	return s.repo.Update(friendship)
}

// Recusar solicitação de amizade
func (s *FriendshipService) DeclineFriendRequest(id uuid.UUID) error {
	friendship, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Garantir que a solicitação esteja "pending" antes de recusar
	if friendship.Status != models.StatusPending {
		return errors.New("only pending requests can be declined")
	}

	friendship.Status = models.StatusDeclined
	return s.repo.Update(friendship)
}

// Buscar amigos aceitos de um usuário
func (s *FriendshipService) GetFriendsForUser(userID uuid.UUID) ([]models.Friendship, error) {
	return s.repo.GetFriendsForUser(userID)
}

// Buscar solicitações pendentes de amizade para um usuário
func (s *FriendshipService) GetPendingRequestsForUser(userID uuid.UUID) ([]models.Friendship, error) {
	return s.repo.GetPendingRequestsForUser(userID)
}
