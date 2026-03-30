package api

import (
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}

func UpdateMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var updates models.User
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
	}

	database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"name": updates.Name,
		// email usually requires Supabase Auth trigger to change securely.
	})
	
	return c.SendStatus(200)
}

// --- ADMIN ---

func ListUsers(c *fiber.Ctx) error {
	// Optionally parse filters from c.Query()
	var users []models.User
	database.DB.Find(&users)
	return c.JSON(users)
}

func GetUserAdmin(c *fiber.Ctx) error {
	var user models.User
	if err := database.DB.First(&user, "id = ?", c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(user)
}

func EditUserAdmin(c *fiber.Ctx) error {
	var updates models.User
	c.BodyParser(&updates)
	database.DB.Model(&models.User{}).Where("id = ?", c.Params("id")).Updates(updates)
	return c.SendStatus(200)
}
