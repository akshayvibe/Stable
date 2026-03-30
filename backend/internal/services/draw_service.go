package services

import (
	"fmt"
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// SimulateDraw calculates distributions without saving
func SimulateDraw(totalPool float64) models.Draw {
	// Requirements: 40% jackpot, 35% tier2, 25% tier3
	jackpot := totalPool * 0.40
	tier2 := totalPool * 0.35
	tier3 := totalPool * 0.25

	return models.Draw{
		DrawDate:    time.Now(),
		TotalPool:   totalPool,
		JackpotPool: jackpot,
		Tier2Pool:   tier2,
		Tier3Pool:   tier3,
		Status:      "pending",
	}
}

// ExecuteDraw formally sets the jackpot pool, potentially resolving previous rollover.
func ExecuteDraw(totalSubscriptions int, subscriptionCost float64) (*models.Draw, error) {
	var newDraw *models.Draw

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Calculate pool (assume pool represents subset of revenue for prizes, e.g., 50% going to pool)
		// For MVP, letting total = subscriptions * cost
		grossPool := float64(totalSubscriptions) * subscriptionCost

		// Check if previous draw had NO jackpot winner
		var lastDraw models.Draw
		rollover := 0.0

		// Find the most recently completed draw
		if err := tx.Where("status = ?", "published").Order("draw_date desc").First(&lastDraw).Error; err == nil {
			if !lastDraw.JackpotWon {
				rollover = lastDraw.JackpotPool
			}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		simulation := SimulateDraw(grossPool)
		// Add the rollover directly to the new Jackpot Pool
		simulation.JackpotPool += rollover

		// In a real drawing, setting actual winning balls here.
		simulation.Status = "drawn"

		if err := tx.Create(&simulation).Error; err != nil {
			return err
		}

		newDraw = &simulation
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("draw execution failed: %v", err)
	}

	return newDraw, nil
}
