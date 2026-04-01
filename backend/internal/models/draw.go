package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/lib/pq"
)

type Draw struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DrawDate    time.Time      `gorm:"not null" json:"draw_date"`
	TotalPool   float64        `gorm:"not null;default:0" json:"total_pool"`
	
	// Math logic: natively hold drawing array
	WinningNumbers pq.Int64Array `gorm:"type:integer[]" json:"winning_numbers"`
	
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
