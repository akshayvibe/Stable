package api

import (
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func ListCharities(c *fiber.Ctx) error {
	var charities []models.Charity
	database.DB.Find(&charities)
	return c.JSON(charities)
}

func GetCharity(c *fiber.Ctx) error {
	id := c.Params("id")
	var charity models.Charity
	if err := database.DB.First(&charity, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(charity)
}

func GetOwnCharity(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var sub models.Subscription
	if err := database.DB.Preload("Charity").Where("user_id = ?", userID).First(&sub).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "no subscription found"})
	}
	return c.JSON(sub.Charity)
}

func SetOwnCharity(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	
	type Input struct {
		CharityID           uint    `json:"charity_id"`
		ContributionPercent float64 `json:"contribution_percent"`
	}
	var body Input
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad payload"})
	}

	if body.ContributionPercent < 10.0 {
		return c.Status(400).JSON(fiber.Map{"error": "Minimum contribution is 10%"})
	}

	err := database.DB.Model(&models.Subscription{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"charity_id":           body.CharityID,
		"contribution_percent": body.ContributionPercent,
	}).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update selection"})
	}

	return c.JSON(fiber.Map{"message": "charity updated"})
}

// --- ADMIN --- 

func AddCharity(c *fiber.Ctx) error {
	var charity models.Charity
	c.BodyParser(&charity)
	database.DB.Create(&charity)
	return c.Status(201).JSON(charity)
}

func EditCharity(c *fiber.Ctx) error {
	id := c.Params("id")
	var updates models.Charity
	c.BodyParser(&updates)
	database.DB.Model(&models.Charity{}).Where("id = ?", id).Updates(updates)
	return c.SendStatus(200)
}

func DeleteCharity(c *fiber.Ctx) error {
	database.DB.Delete(&models.Charity{}, c.Params("id"))
	return c.SendStatus(200)
}
