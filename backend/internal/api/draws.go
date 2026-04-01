package api

import (
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"
	"impact5-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

func ListDraws(c *fiber.Ctx) error {
	var draws []models.Draw
	database.DB.Where("status = ?", "published").Order("draw_date desc").Find(&draws)
	return c.JSON(draws)
}

func GetDraw(c *fiber.Ctx) error {
	id := c.Params("id")
	var draw models.Draw
	if err := database.DB.Preload("Winners").First(&draw, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(draw)
}

func GetCurrentDrawInfo(c *fiber.Ctx) error {
	var count int64
	database.DB.Model(&models.Subscription{}).Where("status = ?", "active").Count(&count)

	jackpot, tier2, tier3 := services.DrawPoolPreview(int(count))

	sim := services.SimulateDraw(float64(count)*services.SubscriptionFee, "random")
	sim.JackpotPool = jackpot
	sim.Tier2Pool = tier2
	sim.Tier3Pool = tier3

	return c.JSON(fiber.Map{
		"draw":            sim,
		"active_subs":     count,
		"jackpot_pool":    jackpot,
		"tier2_pool":      tier2,
		"tier3_pool":      tier3,
		"total_pool":      float64(count) * services.SubscriptionFee,
		"winning_numbers": sim.WinningNumbers,
		"draw_date":       sim.DrawDate,
	})
}


// --- ADMIN ROUTES ---

func CreateDraw(c *fiber.Ctx) error {
	draw := models.Draw{}
	database.DB.Create(&draw)
	return c.JSON(draw)
}

func SimulateDrawAdmin(c *fiber.Ctx) error {
	// Parse logic flag e.g. /api/admin/draws/simulate?logic=algorithmic
	logic := c.Query("logic", "random")
	sim := services.SimulateDraw(1000.0, logic)
	return c.JSON(sim)
}

func ExecuteDrawAdmin(c *fiber.Ctx) error {
	var count int64
	database.DB.Model(&models.Subscription{}).Where("status = ?", "active").Count(&count)
	
	logic := c.Query("logic", "random")

	draw, err := services.ExecuteDraw(int(count), 20.0, logic)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(draw)
}

func PublishDrawAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Model(&models.Draw{}).Where("id = ?", id).Update("status", "published")
	return c.JSON(fiber.Map{"status": "published"})
}
