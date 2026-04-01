package services

import (
	"errors"
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"
	"time"
	
	"gorm.io/gorm"
)

func AddScore(userID string, value int, playedAt time.Time) (*models.Score, error) {
	if value < 1 || value > 45 {
		return nil, errors.New("score value must be between 1 and 45")
	}

	var newScore *models.Score

	// Execute within a single transaction to maintain atomic safety (05 Score Bounds)
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Count existing
		var count int64
		if err := tx.Model(&models.Score{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
			return err
		}

		// 2. Drop oldest if 5 exist
		if count >= 5 {
			var oldest models.Score
			if err := tx.Where("user_id = ?", userID).Order("played_at asc").First(&oldest).Error; err != nil {
				return err
			}
			
			if err := tx.Delete(&oldest).Error; err != nil {
				return err
			}
		}

		// 3. Insert new authentic score explicitly retaining valid time
		score := models.Score{
			UserID:   userID,
			Value:    value, 
			PlayedAt: playedAt,
		}

		if err := tx.Create(&score).Error; err != nil {
			return err
		}
		
		newScore = &score
		return nil
	})

	return newScore, err
}
