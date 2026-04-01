package api

import (
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

// AdminStats — live aggregated analytics from PostgreSQL
func AdminStats(c *fiber.Ctx) error {
	var userCount, subCount int64
	var totalCharity, totalPool float64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Subscription{}).Where("status = ?", "active").Count(&subCount)

	// Sum all active subscription amounts => prize pool (80% of revenue)
	database.DB.Model(&models.Subscription{}).
		Where("status = ?", "active").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPool)

	// Sum charity contributions (contribution_percent * amount per sub)
	type CharitySum struct{ Total float64 }
	var cs CharitySum
	database.DB.Raw(`SELECT COALESCE(SUM(amount * contribution_percent / 100), 0) as total FROM subscriptions WHERE status = 'active'`).Scan(&cs)
	totalCharity = cs.Total

	// Count draws this month
	var drawCount int64
	database.DB.Model(&models.Draw{}).Where("draw_date >= date_trunc('month', NOW())").Count(&drawCount)

	// Count pending winners (awaiting verification)
	var pendingWinners int64
	database.DB.Model(&models.Winner{}).Where("status = ?", "pending").Count(&pendingWinners)

	return c.JSON(fiber.Map{
		"total_users":      userCount,
		"active_subs":      subCount,
		"total_pool":       totalPool * 0.80,
		"total_charity":    totalCharity,
		"draws_this_month": drawCount,
		"pending_winners":  pendingWinners,
	})
}

// AdminSubscriptionReports — per-subscription breakdown
func AdminSubscriptionReports(c *fiber.Ctx) error {
	type SubReport struct {
		Status string  `json:"status"`
		Count  int64   `json:"count"`
		Total  float64 `json:"total_revenue"`
	}
	var reports []SubReport
	database.DB.Raw(`
		SELECT status, COUNT(*) as count, COALESCE(SUM(amount), 0) as total_revenue
		FROM subscriptions
		GROUP BY status
	`).Scan(&reports)
	return c.JSON(reports)
}

// AdminCharityReports — charity-level contribution breakdown
func AdminCharityReports(c *fiber.Ctx) error {
	type CharityReport struct {
		CharityID   uint    `json:"charity_id"`
		CharityName string  `json:"charity_name"`
		SubCount    int64   `json:"subscriber_count"`
		TotalRouted float64 `json:"total_routed"`
	}
	var reports []CharityReport
	database.DB.Raw(`
		SELECT s.charity_id, c.name as charity_name,
		       COUNT(s.id) as sub_count,
		       COALESCE(SUM(s.amount * s.contribution_percent / 100), 0) as total_routed
		FROM subscriptions s
		JOIN charities c ON c.id = s.charity_id
		WHERE s.status = 'active'
		GROUP BY s.charity_id, c.name
		ORDER BY total_routed DESC
	`).Scan(&reports)
	return c.JSON(reports)
}


// AdminDrawStats — draw-specific analytics
func AdminDrawStats(c *fiber.Ctx) error {
	type DrawStat struct {
		Status     string  `json:"status"`
		Count      int64   `json:"count"`
		TotalPrize float64 `json:"total_prize"`
	}
	var stats []DrawStat
	database.DB.Raw(`
		SELECT status, COUNT(*) as count, COALESCE(SUM(jackpot_pool), 0) as total_prize
		FROM draws
		GROUP BY status
	`).Scan(&stats)
	return c.JSON(stats)
}

