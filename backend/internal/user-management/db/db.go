package db

import (
	"backend/internal/user-management/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// dsn := "host=localhost user=youruser password=yourpassword dbname=yourdbname port=5432 sslmode=disable"

	_, err := os.Stat(".env")
	if os.IsNotExist(err) {
		log.Fatalf(".env file not found in the current working directory")
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Build the DSN (Data Source Name)
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	// fmt.Println(dsn)

	var dsn string
	if password == "" {
		// If no password, omit the password field in DSN
		dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s", host, user, dbname, port, sslmode)
	} else {
		// If password is set, include it
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	}

	// Connect to PostgreSQL database

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema for the User model
	DB.AutoMigrate(&models.User{})
}
