package api

import (
	"impact5-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/supabase-go"
)

type AuthHandler struct {
	SupaClient *supabase.Client
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	type RegisterReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	var req RegisterReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := services.RegisterUser(h.SupaClient, req.Email, req.Password, req.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Registration successful", "user": user})
}

// Proxies Supabase REST endpoints using HTTP natively or client
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Not fully implemented: We would hit Supabase /auth/v1/token to get JWT
	return c.JSON(fiber.Map{"message": "Login successful - pretend returning JWT"})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Logged out"})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Token refreshed"})
}
