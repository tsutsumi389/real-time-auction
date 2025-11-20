package repository

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates a new database connection
func NewDatabase() (*gorm.DB, error) {
	// Get database connection details from environment variables
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "auction_user")
	dbPassword := getEnv("DB_PASSWORD", "auction_pass")
	dbName := getEnv("DB_NAME", "auction_db")

	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Tokyo",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	// Configure GORM logger
	var gormLogger logger.Interface
	if getEnv("ENV", "development") == "production" {
		gormLogger = logger.Default.LogMode(logger.Error)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully")

	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
