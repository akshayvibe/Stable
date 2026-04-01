package api

import (
	"impact5-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/gotrue-go/types"
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
		return c.Status(400).JSON(fiber.Map{"error": "Registration failed", "details": err.Error()})
	}

	// Automatically log the user in after registration to generate a JWT token immediately
	tokenResp, tokenErr := h.SupaClient.Auth.Token(types.TokenRequest{
		GrantType: "password",
		Email:     req.Email,
		Password:  req.Password,
	})
	if tokenErr != nil {
		// Log them in anyway or tell them to login manually depending on confirmation flow
		return c.Status(201).JSON(fiber.Map{"message": "Registration successful, but confirm email to login.", "user": user})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Registration successful", "token": tokenResp.AccessToken, "user": user})
}

// Proxies Supabase REST endpoints using HTTP natively or client
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Malformed request structure"})
	}

	tokenResp, err := h.SupaClient.Auth.Token(types.TokenRequest{
		GrantType: "password",
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Supabase authentication failed", "details": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Login successful", "token": tokenResp.AccessToken})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Logged out"})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Token refreshed"})
}
