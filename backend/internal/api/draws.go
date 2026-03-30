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
	// E.g., simulate based on active subscriptions right now
	var count int64
	database.DB.Model(&models.Subscription{}).Where("status = ?", "active").Count(&count)

	sim := services.SimulateDraw(float64(count) * 20.0) // Assume $20 sub
	return c.JSON(sim)
}

// --- ADMIN ROUTES ---

func CreateDraw(c *fiber.Ctx) error {
	// Standard empty creation
	draw := models.Draw{}
	database.DB.Create(&draw)
	return c.JSON(draw)
}

func SimulateDrawAdmin(c *fiber.Ctx) error {
	// ... retrieve active subscriptions
	sim := services.SimulateDraw(1000.0) // mockup pool
	return c.JSON(sim)
}

func ExecuteDrawAdmin(c *fiber.Ctx) error {
	var count int64
	database.DB.Model(&models.Subscription{}).Where("status = ?", "active").Count(&count)

	draw, err := services.ExecuteDraw(int(count), 20.0) // $20 sub rate assumption
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
