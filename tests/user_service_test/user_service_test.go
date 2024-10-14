package user_service_test

import (
	"GoVersi/internal/models"
	services "GoVersi/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository é um mock da interface UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) UsernameExists(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(userID uuid.UUID) (*models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) SuspendUser(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) RequestAccountDeletion(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) PermanentlyDeleteUser(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersWithPendingDeletion() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestLoginUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	// Prepara um usuário simulado
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	mockRepo.On("FindByEmail", "test@example.com").Return(&models.User{Password: string(hashedPassword)}, nil)

	// Teste de sucesso
	token, err := userService.LoginUser("test@example.com", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Teste de falha ao tentar logar com credenciais inválidas
	_, err = userService.LoginUser("test@example.com", "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	// Teste de sucesso
	mockRepo.On("UsernameExists", "testuser").Return(false, nil)
	mockRepo.On("Create", mock.Anything).Return(nil)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}

	err := userService.RegisterUser(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Teste de falha ao registrar um usuário com username existente
	mockRepo.On("UsernameExists", "existinguser").Return(true, nil)

	user2 := &models.User{
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "password",
	}

	err = userService.RegisterUser(user2)

	assert.Error(t, err)
	assert.Equal(t, "username already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestSuspendUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	userID := uuid.New()

	// Teste de sucesso
	mockRepo.On("SuspendUser", userID).Return(nil)

	err := userService.SuspendUser(userID.String())
	assert.NoError(t, err)

	// Teste de falha ao suspender usuário inválido
	invalidUserID := "invalid-uuid"
	err = userService.SuspendUser(invalidUserID)
	assert.Error(t, err)
	assert.Equal(t, "invalid user ID", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetUserById(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	userID := uuid.New()
	expectedUser := &models.User{ID: userID, Username: "testuser"}

	// Teste de sucesso
	mockRepo.On("FindByID", userID).Return(expectedUser, nil)

	user, err := userService.GetUserById(userID.String())
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Teste de falha com ID inválido
	invalidUserID := "invalid-uuid"
	_, err = userService.GetUserById(invalidUserID)
	assert.Error(t, err)
	assert.Equal(t, "invalid user ID", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestRequestAccountDeletion(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	userID := uuid.New()

	// Teste de sucesso
	mockRepo.On("RequestAccountDeletion", userID).Return(nil)

	err := userService.RequestAccountDeletion(userID.String())
	assert.NoError(t, err)

	// Teste de falha ao passar um ID inválido
	invalidUserID := "invalid-uuid"
	err = userService.RequestAccountDeletion(invalidUserID)
	assert.Error(t, err)
	assert.Equal(t, "invalid user ID", err.Error())

	mockRepo.AssertExpectations(t)
}
