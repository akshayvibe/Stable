package database

import (
	"log"
	"os"

	"impact5-backend/internal/models"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	log.Println("[INFO] Initializing PostgreSQL connection protocol...")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("[ERROR] DATABASE_URL environment variable required but not set")
	}

	var err error
	log.Println("[INFO] Dialing active PostgreSQL Database node...")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("[ERROR] Failed to securely connect to PostgreSQL: %v\n", err)
	}

	log.Println("[INFO] Established secured connection to PostgreSQL Database!")
}

func Migrate() {
	log.Println("[INFO] Executing automated database migrations...")
	err := DB.AutoMigrate(
		&models.User{},
		&models.Charity{},
		&models.Subscription{},
		&models.Score{},
		&models.Draw{},
		&models.Winner{},
	)
	if err != nil {
		log.Fatalf("[ERROR] Failed to execute database schemas migrations: %v\n", err)
	}

	log.Println("[INFO] All database schemas are successfully synchronized and up to date.")
}
