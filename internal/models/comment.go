package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Content   string    `json:"content"`
	PostID    string    `json:"post_id"`
	AuthorID  uuid.UUID `json:"author_id"` // Alterado para UUID
	ImageURL  string    `json:"image_url"` // Novo campo para armazenar o caminho da imagem
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
