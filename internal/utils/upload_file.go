package utils

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleImageUpload(c *gin.Context, uploadPath string) (string, error) {
	contentType := c.GetHeader("Content-Type")

	// If it's JSON request, return empty string without error
	if contentType == "application/json" {
		return "", nil
	}

	// Check if it's multipart form data
	if !strings.Contains(contentType, "multipart/form-data") {
		return "", nil
	}

	// Try to get the image file
	file, err := c.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil // No file uploaded, return empty string
		}
		return "", err
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", errors.New("failed to create upload directory")
	}

	// Save the file
	imagePath := filepath.Join(uploadPath, file.Filename)
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		return "", errors.New("failed to upload image")
	}

	return "/" + imagePath, nil
}

// Função para lidar com o upload de vídeo
func HandleVideoUpload(c *gin.Context, uploadPath string) (string, error) {
	contentType := c.GetHeader("Content-Type")

	// If it's JSON request, return empty string without error
	if contentType == "application/json" {
		return "", nil
	}

	// Check if it's multipart form data
	if !strings.Contains(contentType, "multipart/form-data") {
		return "", nil
	}

	// Try to get the video file
	file, err := c.FormFile("video")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil // No file uploaded, return empty string
		}
		return "", err
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", errors.New("failed to create upload directory")
	}

	// Save the file
	videoPath := filepath.Join(uploadPath, file.Filename)
	if err := c.SaveUploadedFile(file, videoPath); err != nil {
		return "", errors.New("failed to upload video")
	}

	return "/" + videoPath, nil
}
