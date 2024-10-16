package repository

import (
	"GoVersi/internal/models"

	"github.com/google/uuid"
)

// UserRepository define a interface para operações de banco de dados relacionadas a usuários
type UserRepository interface {
	GetUserByID(userID uuid.UUID) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID uuid.UUID) error
	GetUsersWithPendingDeletion() ([]models.User, error)

	// Métodos adicionais
	SuspendUser(userID uuid.UUID) error
	PermanentlyDeleteUser(userID uuid.UUID) error
	Create(user *models.User) error
	UsernameExists(username string) (bool, error)
	FindByEmail(email string) (*models.User, error)

	// Novos métodos
	FindByID(userID uuid.UUID) (*models.User, error)      // Para buscar usuário por ID
	FindByUsername(username string) (*models.User, error) // Para buscar usuário por nome de usuário
	RequestAccountDeletion(userID uuid.UUID) error        // Para solicitar a exclusão de conta
	FindByEmailConfirmToken(token string) (*models.User, error)
}
