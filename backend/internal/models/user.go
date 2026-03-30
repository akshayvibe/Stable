package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"` // Maps directly to Supabase auth.users ID
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Name      string         `json:"name"`
	Role      string         `gorm:"default:'player'" json:"role"` // player, admin
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subscription *Subscription `json:"subscription,omitempty"`
	Scores       []Score      `json:"scores,omitempty"`
	Winnings     []Winner     `json:"winnings,omitempty"`
}
