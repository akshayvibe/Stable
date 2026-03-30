package api

import (
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func GetMeWinners(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var winnings []models.Winner
	database.DB.Where("user_id = ?", userID).Preload("Draw").Find(&winnings)
	return c.JSON(winnings)
}

func UploadProof(c *fiber.Ctx) error {
	// For multipart/form-data image uploads
	// FormFile("screenshot")
	id := c.Params("id")
	// Upload to Supabase Storage -> Get URL
	mockURL := "https://supabase.net/storage/v1/object/public/proofs/" + id + ".png"

	database.DB.Model(&models.Winner{}).Where("id = ?", id).Updates(map[string]interface{}{
		"proof_url": mockURL,
		"status":    "pending", // triggers admin review
	})

	return c.JSON(fiber.Map{"message": "Proof uploaded"})
}

// --- ADMIN ---

func ListWinners(c *fiber.Ctx) error {
	var winners []models.Winner
	database.DB.Preload("User").Preload("Draw").Find(&winners)
	return c.JSON(winners)
}

func VerifyWinner(c *fiber.Ctx) error {
	id := c.Params("id")
	type Body struct {
		Action string `json:"action"` // approve / reject
	}
	var b Body
	c.BodyParser(&b)

	status := "rejected"
	if b.Action == "approve" {
		status = "verified"
	}

	database.DB.Model(&models.Winner{}).Where("id = ?", id).Update("status", status)
	return c.SendStatus(200)
}

func PayoutWinner(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Model(&models.Winner{}).Where("id = ?", id).Update("status", "paid")
	return c.SendStatus(200)
}
