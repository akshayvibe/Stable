package models

import (
	"time"
)

type Winner struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DrawID      uint      `gorm:"not null;index" json:"draw_id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	MatchType   string    `gorm:"not null" json:"match_type"` // "jackpot_5", "match_4", "match_3"
	PrizeAmount float64   `gorm:"not null" json:"prize_amount"`
	
	// Verification
	ProofURL    *string   `json:"proof_url,omitempty"`   // S3/Supabase Storage link to screenshot
	Status      string    `gorm:"default:'pending'" json:"status"` // pending, verified, rejected, paid

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Draw        *Draw     `json:"-"`
	User        *User     `json:"user,omitempty"`
}
