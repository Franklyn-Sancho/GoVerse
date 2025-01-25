package services

import (
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
	"GoVersi/internal/service/email"
	"GoVersi/internal/utils"
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
)

type UserService struct {
	UserRepo     repository.UserRepository
	EmailService email.EmailService
}

func NewUserService(repo repository.UserRepository, emailService email.EmailService) *UserService {
	return &UserService{
		UserRepo:     repo,
		EmailService: emailService,
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

	if err := s.EmailService.SendConfirmationEmail(user.Email, user.Username, user.EmailConfirmToken); err != nil {
		log.Printf("Erro ao enviar email: %v", err)
	}

	log.Printf("Usuário registrado com sucesso: %s", user.Username)
	return nil
}

func (s *UserService) LoginUser(email, password string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("something went wrong")
	}

	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")

	return utils.GenerateJWT(user.ID.String(), secretKey)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.UserRepo.UpdateUser(user)
}

func (s *UserService) SuspendUser(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.UserRepo.SuspendUser(userID)
}

func (s *UserService) RequestAccountDeletion(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.UserRepo.RequestAccountDeletion(userID)
}

func (s *UserService) GetUserById(id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.UserRepo.FindByID(userID)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.UserRepo.FindByUsername(username)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.UserRepo.FindByEmail(email)
}

func (s *UserService) DeleteUser(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}
	return s.UserRepo.PermanentlyDeleteUser(userID)
}

func (s *UserService) PermanentlyDeleteUser(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID")
	}
	return s.UserRepo.PermanentlyDeleteUser(userID)
}

func (s *UserService) FindByEmailConfirmToken(token string) (*models.User, error) {
	return s.UserRepo.FindByEmailConfirmToken(token)
}
