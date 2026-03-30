package api

import (
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func Checkout(c *fiber.Ctx) error {
	// Stripe Checkout Session generation goes here
	return c.JSON(fiber.Map{"url": "https://checkout.stripe.com/test"})
}

func GetMeSubscription(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var sub models.Subscription
	if err := database.DB.Where("user_id = ?", userID).First(&sub).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "No active subscription"})
	}
	return c.JSON(sub)
}

func CancelSubscription(c *fiber.Ctx) error {
	// Stripe cancel logic
	return c.JSON(fiber.Map{"status": "canceled at period end"})
}

func Portal(c *fiber.Ctx) error {
	// Stripe customer portal redirect
	return c.JSON(fiber.Map{"url": "https://billing.stripe.com/p/session/test"})
}

// Webhook endpoint (not protected by JWT, but by Stripe Signature)
func StripeWebhook(c *fiber.Ctx) error {
	// Read payload and verify with Stripe secret
	// E.g. Update subscription status upon checkout.session.completed or invoice.payment_failed
	return c.SendStatus(200)
}

// --- ADMIN ---

func ListSubscriptions(c *fiber.Ctx) error {
	var subs []models.Subscription
	database.DB.Preload("User").Find(&subs)
	return c.JSON(subs)
}
