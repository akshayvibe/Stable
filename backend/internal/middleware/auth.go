package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Protected(supabaseJWTSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Native Supabase API Validation bypassing local `.env` secret mismatch
		supaURL := os.Getenv("SUPABASE_URL")
		req, _ := http.NewRequest("GET", supaURL+"/auth/v1/user", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token remotely mapped"})
		}
		defer resp.Body.Close()

		var supaUser struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&supaUser); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Failed to decode robust token response"})
		}

		c.Locals("user_id", supaUser.ID)
		
		return c.Next()
	}
}
