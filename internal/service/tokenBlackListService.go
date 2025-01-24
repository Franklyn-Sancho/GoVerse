package services

import (
	"GoVersi/internal/models"
	"log"
	"time"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

type TokenBlacklistService struct {
	DB *gorm.DB
}

func NewTokenBlacklistService(db *gorm.DB) *TokenBlacklistService {
	return &TokenBlacklistService{DB: db}
}

func (s *TokenBlacklistService) AddToTokenBlacklist(token string, expiresAt time.Time) error {
	tokenBlacklistEntry := models.TokenBlacklist{
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return s.DB.Create(&tokenBlacklistEntry).Error
}

func (s *TokenBlacklistService) IsTokenBlacklisted(token string) (bool, error) {
	var count int64
	err := s.DB.Model(&models.TokenBlacklist{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *TokenBlacklistService) RemoveExpiredTokens() error {
	return s.DB.Where("expires_at < ?", time.Now()).Delete(&models.TokenBlacklist{}).Error
}

func (s *TokenBlacklistService) StartCronJob() {
	c := cron.New()
	c.AddFunc("@daily", func() {
		err := s.RemoveExpiredTokens()
		if err != nil {
			log.Printf("Failed to remove expired tokens: %v", err)
		} else {
			log.Println("Expired tokens removed from blacklist")
		}
	})
	c.Start()
}
