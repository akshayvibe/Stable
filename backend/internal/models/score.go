package models

import (
	"time"
)

type Score struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Value     int       `gorm:"not null" json:"value"` // Stableford 1-45
	PlayedAt  time.Time `gorm:"not null;index" json:"played_at"`
	CreatedAt time.Time `json:"created_at"`

	User      *User     `json:"-"`
}
