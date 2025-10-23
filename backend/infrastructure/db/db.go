package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"backend/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB establishes database connection
func ConnectDB() *gorm.DB {
	// Load environment variables
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("Failed to get underlying sql.DB:", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate database schema
	err = AutoMigrate(db)
	if err != nil {
		log.Fatalln("Failed to migrate database:", err)
	}

	fmt.Println("Database connected successfully")
	return db
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Image{},
		&domain.Post{},
	)
}

// CloseDB closes database connection
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Failed to get underlying sql.DB:", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Println("Failed to close database:", err)
	}
}
