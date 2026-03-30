package api

import (
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func AdminStats(c *fiber.Ctx) error {
	var userCount, subCount int64
	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Subscription{}).Where("status = ?", "active").Count(&subCount)

	return c.JSON(fiber.Map{
		"total_users": userCount,
		"active_subs": subCount,
		"total_pools": 150000.0,
		"total_charity": 50000.0,
	})
}

func AdminSubscriptionReports(c *fiber.Ctx) error {
	// Aggregate mock revenue over time query
	return c.JSON(fiber.Map{"data": "subscription reports graph data..."})
}

func AdminCharityReports(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"data": "charity distribution sums..."})
}
