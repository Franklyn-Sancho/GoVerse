package repository

import (
	"GoVersi/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepositoryImpl é uma implementação da interface UserRepository
type UserRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserRepository cria uma nova instância de UserRepositoryImpl
func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{DB: db}
}

// Implementação da interface
func (r *UserRepositoryImpl) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Implementação do método FindByEmailConfirmToken
func (r *UserRepositoryImpl) FindByEmailConfirmToken(token string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email_confirm_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepositoryImpl) DeleteUser(userID uuid.UUID) error {
	return r.DB.Delete(&models.User{}, userID).Error
}

// Implementação do método Create
func (r *UserRepositoryImpl) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// Implementação do método UsernameExists
func (r *UserRepositoryImpl) UsernameExists(username string) (bool, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // Usuário não encontrado
		}
		return false, err // Outro erro
	}
	return true, nil // Usuário encontrado
}

func (r *UserRepositoryImpl) GetUsersWithPendingDeletion() ([]models.User, error) {
	var users []models.User
	// Implemente a lógica para obter usuários com solicitação de exclusão pendente
	if err := r.DB.Where("is_pending_deletion = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Implementação do método FindByEmail
func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Implementação do método FindByID
func (r *UserRepositoryImpl) FindByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Implementação do método FindByUsername
func (r *UserRepositoryImpl) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Implementação do método Delete
func (r *UserRepositoryImpl) Delete(userID uuid.UUID) error {
	var user models.User
	if err := r.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}
	return r.DB.Delete(&user).Error
}

// Implementação do método SuspendUser
func (r *UserRepositoryImpl) SuspendUser(userID uuid.UUID) error {
	user, err := r.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.IsActive = false // Marca o usuário como suspenso
	return r.UpdateUser(user)
}

// Implementação do método RequestAccountDeletion
func (r *UserRepositoryImpl) RequestAccountDeletion(userID uuid.UUID) error {
	user, err := r.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.IsPendingDeletion = true // Marca a conta como pendente de exclusão
	now := time.Now()
	user.DeletionRequestedAt = &now
	return r.UpdateUser(user)
}

// Implementação do método PermanentlyDeleteUser
func (r *UserRepositoryImpl) PermanentlyDeleteUser(userID uuid.UUID) error {
	return r.DeleteUser(userID) // Deleta o usuário permanentemente
}
