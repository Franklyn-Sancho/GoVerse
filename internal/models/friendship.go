package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type FriendshipStatus string

const (
	StatusPending  FriendshipStatus = "pending"
	StatusAccepted FriendshipStatus = "accepted"
	StatusDeclined FriendshipStatus = "declined"
)

// Friendship model
type Friendship struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RequesterID uuid.UUID        `json:"requester_id"`
	AddresseeID uuid.UUID        `json:"addressee_id"`
	Status      FriendshipStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// Validar o status antes de salvar
func (f *Friendship) Validate() error {
	switch f.Status {
	case StatusPending, StatusAccepted, StatusDeclined:
		return nil
	default:
		return errors.New("invalid status: must be 'pending', 'accepted', or 'declined'")
	}
}
