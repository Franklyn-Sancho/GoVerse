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
	IsActive            bool       `json:"is_active"`           // true se o usuário está ativo, false se suspenso
	IsPendingDeletion   bool       `json:"is_pending_deletion"` // true se a exclusão foi solicitada
	DeletionRequestedAt *time.Time `json:"deletion_requested_at,omitempty"`
}

// Método para gerar um novo UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
