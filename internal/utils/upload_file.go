package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func isValidImageExt(ext string) bool {
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return validExts[ext]
}

func isValidVideoExt(ext string) bool {
	validExts := map[string]bool{
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".wmv":  true,
		".mkv":  true,
		".webm": true,
	}
	return validExts[ext]
}

func HandleImageUpload(c *gin.Context, uploadPath string) (string, error) {
	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		return "", nil
	}

	if !strings.Contains(contentType, "multipart/form-data") {
		return "", nil
	}

	file, err := c.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}

	if file.Size > 10*1024*1024 {
		return "", errors.New("image file too large (max 10MB)")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExt(ext) {
		return "", errors.New("invalid image format")
	}

	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102-150405")
	safeFilename := fmt.Sprintf("%s_%s%s", timestamp, uniqueID, ext)

	// Create absolute path for storage
	absPath, err := filepath.Abs(uploadPath)
	if err != nil {
		return "", errors.New("failed to get absolute path")
	}

	if err := os.MkdirAll(absPath, os.ModePerm); err != nil {
		return "", errors.New("failed to create upload directory")
	}

	// Save file
	filePath := filepath.Join(absPath, safeFilename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", errors.New("failed to upload image")
	}

	return "/" + filepath.Join(uploadPath, safeFilename), nil
}

func HandleVideoUpload(c *gin.Context, uploadPath string) (string, error) {
	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		return "", nil
	}

	if !strings.Contains(contentType, "multipart/form-data") {
		return "", nil
	}

	file, err := c.FormFile("video")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}

	// Validate file size (100MB limit)
	if file.Size > 100*1024*1024 {
		return "", errors.New("video file too large (max 100MB)")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidVideoExt(ext) {
		return "", errors.New("invalid video format")
	}

	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102-150405")
	safeFilename := fmt.Sprintf("%s_%s%s", timestamp, uniqueID, ext)

	// Create absolute path for storage
	absPath, err := filepath.Abs(uploadPath)
	if err != nil {
		return "", errors.New("failed to get absolute path")
	}

	if err := os.MkdirAll(absPath, os.ModePerm); err != nil {
		return "", errors.New("failed to create upload directory")
	}

	// Save file
	filePath := filepath.Join(absPath, safeFilename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", errors.New("failed to upload video")
	}

	return "/" + filepath.Join("uploads/videos", safeFilename), nil
}
