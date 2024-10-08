package models

import "time"

type TokenBlacklist struct {
	Token     string    `gorm:"primaryKey"`
	ExpiresAt time.Time `gorm:"not null"`
}
