package repository

import (
	"GoVersi/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByID(userID uuid.UUID) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID uuid.UUID) error
	GetUsersWithPendingDeletion() ([]models.User, error)

	SuspendUser(userID uuid.UUID) error
	PermanentlyDeleteUser(userID uuid.UUID) error
	Create(user *models.User) error
	UsernameExists(username string) (bool, error)
	FindByEmail(email string) (*models.User, error)

	FindByID(userID uuid.UUID) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	RequestAccountDeletion(userID uuid.UUID) error
	FindByEmailConfirmToken(token string) (*models.User, error)
}
