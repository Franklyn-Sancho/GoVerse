package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Topic     string    `json:"topic"`
	AuthorID  uuid.UUID `json:"author_id"` // Alterado para UUID
	ImageURL  string    `json:"image_url"` // Novo campo para armazenar o caminho da imagem
	VideoURL  string    `json:"video_url"` // URL do vídeo associado ao post (opcional)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
