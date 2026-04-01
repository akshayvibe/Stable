package services

import (
	"fmt"
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Prize pool distribution constants (Section 7)
const (
	JackpotSharePct float64 = 0.40 // 5-Match: 40% of pool, rolls over if unclaimed
	Tier2SharePct   float64 = 0.35 // 4-Match: 35% of pool, split equally among winners
	Tier3SharePct   float64 = 0.25 // 3-Match: 25% of pool, split equally among winners
)

// SubscriptionFee is the per-subscriber monthly contribution to the prize pool
const SubscriptionFee float64 = 20.0 // $20/subscriber → pool

// PoolTier represents a single prize tier's computed values
type PoolTier struct {
	Name        string  `json:"name"`
	MatchCount  int     `json:"match_count"`
	PoolShare   float64 `json:"pool_share"`   // absolute $ value
	PctOfPool   float64 `json:"pct_of_pool"`  // 0.40, 0.35, 0.25
	WinnerCount int     `json:"winner_count"` // number of winners in this tier (runtime)
	PrizeEach   float64 `json:"prize_each"`   // pool_share / winner_count
	Rollover    bool    `json:"rollover"`     // only jackpot tier rolls over
}

// DrawPoolPreview returns the three live tier values given active subscriber count
func DrawPoolPreview(activeSubs int) (jackpot, tier2, tier3 float64) {
	gross := float64(activeSubs) * SubscriptionFee
	return gross * JackpotSharePct, gross * Tier2SharePct, gross * Tier3SharePct
}

// ─────────────────────────────────────────────────────────────────────────────

func generateRandomNumbers() pq.Int64Array {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := make(map[int64]bool)
	var arr pq.Int64Array
	for len(arr) < 5 {
		n := int64(r.Intn(45) + 1)
		if !nums[n] {
			nums[n] = true
			arr = append(arr, n)
		}
	}
	return arr
}

func generateAlgorithmicNumbers(tx *gorm.DB) pq.Int64Array {
	type Result struct {
		Value int64
		Count int64
	}
	var results []Result
	tx.Model(&models.Score{}).Select("value, count(value) as count").Group("value").Order("count desc").Limit(5).Scan(&results)

	var arr pq.Int64Array
	for _, res := range results {
		if res.Value >= 1 && res.Value <= 45 {
			arr = append(arr, res.Value)
		}
	}

	if len(arr) < 5 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		nums := make(map[int64]bool)
		for _, v := range arr {
			nums[v] = true
		}
		for len(arr) < 5 {
			n := int64(r.Intn(45) + 1)
			if !nums[n] {
				nums[n] = true
				arr = append(arr, n)
			}
		}
	}
	return arr
}

// SimulateDraw computes parameters without DB persistence ("pre-analysis mode before publish")
func SimulateDraw(totalPool float64, logic string) models.Draw {
	jackpot, tier2, tier3 := totalPool*JackpotSharePct, totalPool*Tier2SharePct, totalPool*Tier3SharePct

	var numbers pq.Int64Array
	if logic == "algorithmic" {
		numbers = generateAlgorithmicNumbers(database.DB)
	} else {
		numbers = generateRandomNumbers()
	}

	return models.Draw{
		DrawDate:       time.Now(),
		TotalPool:      totalPool,
		JackpotPool:    jackpot,
		Tier2Pool:      tier2,
		Tier3Pool:      tier3,
		WinningNumbers: numbers,
		Status:         "pending",
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// matchCount returns how many of the user's top 5 scores match any winning number
func matchCount(scores []models.Score, winningNumbers pq.Int64Array) int {
	winSet := make(map[int64]bool)
	for _, n := range winningNumbers {
		winSet[n] = true
	}
	matches := 0
	for _, s := range scores {
		if winSet[int64(s.Value)] {
			matches++
		}
	}
	return matches
}

// ─────────────────────────────────────────────────────────────────────────────
// ExecuteDraw formally anchors the jackpot pool, rolls over previously empty pools,
// then matches each subscriber's scores against the winning numbers and creates Winner records.
func ExecuteDraw(totalSubscriptions int, subscriptionCost float64, logic string) (*models.Draw, error) {
	var newDraw *models.Draw

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		grossPool := float64(totalSubscriptions) * subscriptionCost

		// Retrieve preceding draw to execute jackpot rollover logic
		var lastDraw models.Draw
		rollover := 0.0
		if err := tx.Where("status = ?", "published").Order("draw_date desc").First(&lastDraw).Error; err == nil {
			if !lastDraw.JackpotWon {
				rollover = lastDraw.JackpotPool
			}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		simulation := SimulateDraw(grossPool, logic)
		simulation.JackpotPool += rollover // jackpot carries forward if previous was unclaimed
		simulation.Status = "drawn"

		if err := tx.Create(&simulation).Error; err != nil {
			return err
		}

		// ── Winner Matching Logic (Section 7) ────────────────────────────────
		// Load all active subscriptions with their user IDs
		var subs []models.Subscription
		if err := tx.Where("status = ?", "active").Find(&subs).Error; err != nil {
			return err
		}

		// Buckets: match tier → list of user IDs
		jackpotWinners := []string{}
		tier2Winners := []string{}
		tier3Winners := []string{}

		for _, sub := range subs {
			// Fetch the user's latest 5 scores (rolling window)
			var scores []models.Score
			tx.Where("user_id = ?", sub.UserID).Order("played_at desc").Limit(5).Find(&scores)

			mc := matchCount(scores, simulation.WinningNumbers)
			switch {
			case mc >= 5:
				jackpotWinners = append(jackpotWinners, sub.UserID)
			case mc == 4:
				tier2Winners = append(tier2Winners, sub.UserID)
			case mc == 3:
				tier3Winners = append(tier3Winners, sub.UserID)
			}
		}

		// Set JackpotWon flag — only true if at least one 5-match winner exists
		didWinJackpot := len(jackpotWinners) > 0
		tx.Model(&simulation).Update("jackpot_won", didWinJackpot)

		// Helper: create winner records, splitting the tier pool equally
		createWinners := func(userIDs []string, matchType string, tierPool float64) error {
			if len(userIDs) == 0 {
				return nil
			}
			prizeEach := tierPool / float64(len(userIDs))
			for _, uid := range userIDs {
				winner := models.Winner{
					DrawID:      simulation.ID,
					UserID:      uid,
					MatchType:   matchType,
					PrizeAmount: prizeEach,
					Status:      "pending",
				}
				if err := tx.Create(&winner).Error; err != nil {
					return err
				}
			}
			return nil
		}

		if err := createWinners(jackpotWinners, "jackpot_5", simulation.JackpotPool); err != nil {
			return err
		}
		if err := createWinners(tier2Winners, "match_4", simulation.Tier2Pool); err != nil {
			return err
		}
		if err := createWinners(tier3Winners, "match_3", simulation.Tier3Pool); err != nil {
			return err
		}

		newDraw = &simulation
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("draw execution failed: %v", err)
	}

	return newDraw, nil
}
