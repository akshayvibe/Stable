package database

import (
	"fmt"
	"log"
	"os"

	"impact5-backend/internal/models"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	// e.g. "host=localhost user=postgres password=postgres dbname=impact5 port=5432 sslmode=disable"
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable required but not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v\n", err)
	}

	fmt.Println("Connected to Database successfully!")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Charity{},
		&models.Subscription{},
		&models.Score{},
		&models.Draw{},
		&models.Winner{},
	)
	if err != nil {
		log.Fatalf("Failed to execute database migrations: %v\n", err)
	}

	fmt.Println("Database schemas auto-migrated successfully!")
}
