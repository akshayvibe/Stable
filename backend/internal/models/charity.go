package models

import (
	"time"
)

type Charity struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	WebsiteURL  string    `json:"website_url"`
	LogoURL     string    `json:"logo_url"`
	Is501c3     bool      `gorm:"default:true" json:"is_501c3"` // Must be verified
	TotalRaised float64   `gorm:"default:0" json:"total_raised"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Subscriptions routing to this charity
	Subscriptions []Subscription `json:"-"`
}
