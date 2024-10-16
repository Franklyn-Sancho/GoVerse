package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	PostID    uuid.UUID `json:"post_id"`    // Referência ao comentário
	CommentID uuid.UUID `json:"comment_id"` // Referência ao comentário
	UserID    uuid.UUID `json:"user_id"`    // Referência ao usuário que deu o like
	CreatedAt time.Time `json:"created_at"`
}
