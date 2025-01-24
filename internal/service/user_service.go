package services

import (
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
	"GoVersi/internal/utils"
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
)

type UserService struct {
	UserRepo repository.UserRepository
}

// NewUserService cria uma nova instância de UserService
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

func (s *UserService) RegisterUser(user *models.User) error {
	log.Printf("Iniciando registro do usuário: %s", user.Username)

	exists, err := s.UserRepo.UsernameExists(user.Username)
	if err != nil {
		log.Printf("Erro ao verificar username: %v", err)
		return err
	}
	if exists {
		return errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if err := s.UserRepo.Create(user); err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		return err
	}

	if err := utils.SendConfirmationEmail(user.Email, user.Username, user.EmailConfirmToken); err != nil {
		log.Printf("Erro ao enviar email: %v", err)
	}

	log.Printf("Usuário registrado com sucesso: %s", user.Username)
	return nil
}

// Função para login de usuário
func (s *UserService) LoginUser(email, password string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("something went wrong")
	}

	if user == nil {
		return "", errors.New("invalid credentials") // Usuário não encontrado
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials") // Senha incorreta
	}

	// Obtenha a chave secreta das variáveis de ambiente
	secretKey := os.Getenv("JWT_SECRET_KEY")

	// Converte o UUID para string antes de chamar GenerateJWT
	return utils.GenerateJWT(user.ID.String(), secretKey)
}

// UpdateUser atualiza os dados do usuário no banco de dados
func (s *UserService) UpdateUser(user *models.User) error {
	return s.UserRepo.UpdateUser(user)
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

// Implementação do método FindByEmailConfirmToken no serviço
func (s *UserService) FindByEmailConfirmToken(token string) (*models.User, error) {
	return s.UserRepo.FindByEmailConfirmToken(token)
}
