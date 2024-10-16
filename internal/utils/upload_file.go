package utils

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

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
