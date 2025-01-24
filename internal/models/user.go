package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                  uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username            string     `json:"username" gorm:"unique;not null"`
	Email               string     `json:"email" gorm:"unique;not null"`
	Password            string     `json:"password" gorm:"not null"`
	ImageProfile        string     `json:"image_url"`
	IsActive            bool       `json:"is_active"`
	IsPendingDeletion   bool       `json:"is_pending_deletion"`
	DeletionRequestedAt *time.Time `json:"deletion_requested_at,omitempty"`
	IsEmailVerified     bool       `json:"is_email_verified" gorm:"default:false"`
	EmailConfirmToken   string     `json:"email_confirm_token" gorm:"unique;not null"` // Novo campo para o token de confirmação
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	u.EmailConfirmToken = uuid.New().String() // Gera um token único
	return
}
