package models

import (
	"time"
	"gorm.io/gorm"
)

type Draw struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DrawDate    time.Time      `gorm:"not null" json:"draw_date"`
	TotalPool   float64        `gorm:"not null;default:0" json:"total_pool"`
	// The 5 winning numbers
	Num1        int            `json:"num_1"`
	Num2        int            `json:"num_2"`
	Num3        int            `json:"num_3"`
	Num4        int            `json:"num_4"`
	Num5        int            `json:"num_5"`
	Status      string         `gorm:"default:'pending'" json:"status"` // pending, drawn, published
	JackpotWon  bool           `gorm:"default:false" json:"jackpot_won"`
	
	// Pre-calculated distributions
	JackpotPool float64        `gorm:"default:0" json:"jackpot_pool"`
	Tier2Pool   float64        `gorm:"default:0" json:"tier_2_pool"`
	Tier3Pool   float64        `gorm:"default:0" json:"tier_3_pool"`

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Winners     []Winner       `json:"winners,omitempty"`
}
