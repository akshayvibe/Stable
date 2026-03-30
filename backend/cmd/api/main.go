package main

import (
	"log"
	"os"

	"impact5-backend/internal/api"
	"impact5-backend/internal/database"
	"impact5-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env variables")
	}

	// 1. Connect DB and Migrate schemas
	database.Connect()
	database.Migrate()

	// 2. Setup Supabase Client
	supaURL := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	supaJWT := os.Getenv("SUPABASE_JWT_SECRET")

	client, err := supabase.NewClient(supaURL, supaKey, nil)
	if err != nil {
		log.Fatal("Failed to initialize Supabase client: ", err)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// Health
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	authHandler := &api.AuthHandler{SupaClient: client}
	
	// API Group setup
	v1 := app.Group("/api")

	// --- PUBLIC ROUTES ---
	
	// Webhooks
	v1.Post("/webhooks/stripe", api.StripeWebhook)

	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/forgot-password", func(c *fiber.Ctx) error { return c.SendStatus(200) })

	// --- PROTECTED ROUTES ---
	protected := v1.Group("/", middleware.Protected(supaJWT))
	
	// Auth extensions
	protected.Post("/auth/logout", authHandler.Logout)
	protected.Post("/auth/refresh", authHandler.Refresh)

	// Users
	protected.Get("/users/me", api.GetMe)
	protected.Put("/users/me", api.UpdateMe)
	protected.Put("/users/me/charity", api.SetOwnCharity)
	protected.Get("/users/me/charity", api.GetOwnCharity)

	// Subscriptions
	protected.Post("/subscriptions/checkout", api.Checkout)
	protected.Get("/subscriptions/me", api.GetMeSubscription)
	protected.Post("/subscriptions/cancel", api.CancelSubscription)
	protected.Post("/subscriptions/portal", api.Portal)

	// Scores
	protected.Get("/scores", api.GetOwnScores)
	protected.Post("/scores", api.AddScore)
	protected.Put("/scores/:id", api.EditScore)
	protected.Delete("/scores/:id", api.DeleteScore)

	// Draws (Reads)
	protected.Get("/draws", api.ListDraws)
	protected.Get("/draws/current", api.GetCurrentDrawInfo)
	protected.Get("/draws/:id", api.GetDraw)

	// Charities (Reads)
	protected.Get("/charities", api.ListCharities)
	protected.Get("/charities/:id", api.GetCharity)

	// Winners
	protected.Get("/winners/me", api.GetMeWinners)
	protected.Post("/winners/:id/proof", api.UploadProof)

	// --- ADMIN ROUTES ---
	
	// NOTE: In production we would add an AdminMiddleware checking `role` == 'admin'
	// adminOnly := protected.Group("/admin", middleware.AdminOnly())
	admin := protected.Group("/admin")
	
	// Users Admin
	admin.Get("/users", api.ListUsers)
	admin.Get("/users/:id", api.GetUserAdmin)
	admin.Put("/users/:id", api.EditUserAdmin)

	// Subscriptions Admin
	admin.Get("/subscriptions", api.ListSubscriptions)

	// Scores Admin
	admin.Get("/scores/users/:id", api.GetUserScoresAdmin)
	// admin.Put("/scores/users/:id/:scoreId", api.EditUserScoreAdmin)

	// Draws Admin
	admin.Post("/draws", api.CreateDraw)
	admin.Post("/draws/:id/simulate", api.SimulateDrawAdmin)
	admin.Post("/draws/:id/run", api.ExecuteDrawAdmin)
	admin.Post("/draws/:id/publish", api.PublishDrawAdmin)
	
	// Charities Admin
	admin.Post("/charities", api.AddCharity)
	admin.Put("/charities/:id", api.EditCharity)
	admin.Delete("/charities/:id", api.DeleteCharity)
	
	// Winners Admin
	admin.Get("/winners", api.ListWinners)
	admin.Patch("/winners/:id/verify", api.VerifyWinner)
	admin.Patch("/winners/:id/payout", api.PayoutWinner)

	// Analytics
	admin.Get("/stats", api.AdminStats)
	admin.Get("/reports/subscriptions", api.AdminSubscriptionReports)
	admin.Get("/reports/charity", api.AdminCharityReports)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
