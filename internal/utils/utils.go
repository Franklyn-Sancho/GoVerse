package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword gera um hash da senha
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifica se a senha corresponde ao hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HandleImageUpload(c *gin.Context, uploadPath string) (string, error) {
	// Tenta pegar o arquivo de imagem (caso exista)
	file, err := c.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		return "", errors.New("Error handling image upload")
	}

	if err == nil {
		// Verificar e criar a pasta de upload se não existir
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			return "", errors.New("Failed to create upload directory")
		}

		// Salvar o arquivo no diretório
		imagePath := filepath.Join(uploadPath, file.Filename)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			return "", errors.New("Failed to upload image")
		}

		return "/" + imagePath, nil // Retorna o caminho da imagem
	}

	return "", nil // Retorna vazio se não houver imagem
}

// Função para lidar com o upload de vídeo
func HandleVideoUpload(c *gin.Context, uploadPath string) (string, error) {
	// Tenta pegar o arquivo de vídeo (caso exista)
	file, err := c.FormFile("video")
	if err != nil && err != http.ErrMissingFile {
		return "", errors.New("Error handling video upload")
	}

	if err == nil {
		// Verificar e criar a pasta de upload se não existir
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			return "", errors.New("Failed to create upload directory")
		}

		// Salvar o arquivo no diretório
		videoPath := filepath.Join(uploadPath, file.Filename)
		if err := c.SaveUploadedFile(file, videoPath); err != nil {
			return "", errors.New("Failed to upload video")
		}

		return "/" + videoPath, nil // Retorna o caminho do vídeo
	}

	return "", nil // Retorna vazio se não houver vídeo
}

func SendEmail(to string, subject string, body string) error {
	from := "your_email@example.com" // pode ser qualquer endereço de e-mail

	// Configura o servidor SMTP do MailHog
	smtpHost := "localhost" // MailHog normalmente está rodando no localhost
	smtpPort := "1025"      // porta do MailHog

	// Cria a mensagem
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	// Envia o e-mail sem autenticação
	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, msg)
}
