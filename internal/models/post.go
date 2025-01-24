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
	AuthorID  uuid.UUID `json:"author_id"`
	ImageURL  string    `json:"image_url"`
	VideoURL  string    `json:"video_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
