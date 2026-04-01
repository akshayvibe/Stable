package api

import (
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"
	"impact5-backend/internal/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// POST /api/scores
func AddScore(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	type input struct {
		Value    int    `json:"value"`
		PlayedAt string `json:"played_at"` // Added explicit date mapping
	}
	var body input
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid parsing", "details": err.Error()})
	}

	// Parse date, fallback to now natively if missing
	var playedAt time.Time
	if body.PlayedAt != "" {
		parsed, err := time.Parse(time.RFC3339, body.PlayedAt)
		if err == nil {
			playedAt = parsed
		} else {
			playedAt = time.Now()
		}
	} else {
		playedAt = time.Now()
	}

	score, err := services.AddScore(userID, body.Value, playedAt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(score)
}

// GET /api/scores
func GetOwnScores(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	
	var scores []models.Score
	if err := database.DB.Where("user_id = ?", userID).Order("played_at desc").Find(&scores).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch scores"})
	}

	return c.JSON(scores)
}

// PUT /api/scores/:id
func EditScore(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	scoreID, _ := strconv.Atoi(c.Params("id"))

	type input struct {
		Value int `json:"value"`
	}
	var body input
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	err := database.DB.Model(&models.Score{}).Where("id = ? AND user_id = ?", scoreID, userID).Update("value", body.Value).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update"})
	}

	return c.JSON(fiber.Map{"message": "updated successfully"})
}

// DELETE /api/scores/:id
func DeleteScore(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	scoreID, _ := strconv.Atoi(c.Params("id"))

	err := database.DB.Where("id = ? AND user_id = ?", scoreID, userID).Delete(&models.Score{}).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete"})
	}

	return c.JSON(fiber.Map{"message": "deleted successfully"})
}

// --- ADMIN ---

// GET /api/scores/users/:id
func GetUserScoresAdmin(c *fiber.Ctx) error {
	targetUserID := c.Params("id")
	var scores []models.Score
	database.DB.Where("user_id = ?", targetUserID).Find(&scores)
	return c.JSON(scores)
}
