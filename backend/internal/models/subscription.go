package models

import (
	"time"
)

type Subscription struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	UserID                string    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	StripeCustomerID      string    `gorm:"uniqueIndex" json:"stripe_customer_id"`
	StripeSubscriptionID  string    `gorm:"uniqueIndex" json:"stripe_subscription_id"`
	Status                string    `gorm:"default:'inactive'" json:"status"` // active, past_due, canceled
	CurrentPeriodEnd      time.Time `json:"current_period_end"`
	
	// Charity Preferences
	CharityID             *uint     `json:"charity_id,omitempty"`
	ContributionPercent   float64   `gorm:"default:10.0" json:"contribution_percent"`

	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	User                  *User     `json:"-"`
	Charity               *Charity  `json:"charity,omitempty"`
}
