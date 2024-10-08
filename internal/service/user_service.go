package services

import (
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
	"GoVersi/internal/utils"
	"errors"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepo repository.UserRepository
}

// NewUserService cria uma nova instância de UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		UserRepo: repository.NewUserRepository(db),
	}
}

func (s *UserService) RegisterUser(user *models.User) error {
	// Verifica se o username já existe
	exists, err := s.UserRepo.UsernameExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists")
	}

	// Hash da senha antes de salvar no banco
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Cria o usuário no banco
	return s.UserRepo.Create(user)
}

// Função para login de usuário
func (s *UserService) LoginUser(email, password string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	// Obtenha a chave secreta das variáveis de ambiente
	secretKey := os.Getenv("JWT_SECRET_KEY") // Certifique-se de que a chave secreta está sendo lida corretamente

	// Converte o UUID para string antes de chamar GenerateJWT
	return utils.GenerateJWT(user.ID.String(), secretKey) // Use a chave secreta aqui
}

// Suspende um usuário
func (s *UserService) SuspendUser(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Chama o repositório para suspender o usuário
	return s.UserRepo.SuspendUser(userID)
}

// Solicita a exclusão da conta
func (s *UserService) RequestAccountDeletion(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Chama o repositório para solicitar a exclusão da conta
	return s.UserRepo.RequestAccountDeletion(userID)
}

// Busca usuário pelo ID
func (s *UserService) GetUserById(id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.UserRepo.FindByID(userID)
}

// Busca usuário pelo username
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.UserRepo.FindByUsername(username)
}

// Busca usuário pelo email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.UserRepo.FindByEmail(email)
}

// Função para deletar um usuário
func (s *UserService) DeleteUser(id string) error {
	userID, err := uuid.Parse(id) // Convertendo de string para uuid.UUID
	if err != nil {
		return errors.New("invalid user ID")
	}
	return s.UserRepo.PermanentlyDeleteUser(userID)
}

// Implementação do método PermanentlyDeleteUser
func (s *UserService) PermanentlyDeleteUser(id string) error {
	userID, err := uuid.Parse(id) // Convertendo de string para uuid.UUID
	if err != nil {
		return errors.New("invalid user ID")
	}
	// Lógica para deletar o usuário permanentemente
	return s.UserRepo.PermanentlyDeleteUser(userID) // Chame o método do repositório
}
