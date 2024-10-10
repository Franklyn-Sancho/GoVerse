package models

import (
	"time"

	"github.com/google/uuid"
)

type Friendship struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Requester_id uuid.UUID `json:"requester"`
	Addressee_id uuid.UUID `json:"addressee"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
