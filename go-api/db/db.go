package db

import (
	"fmt"
	"log"
	"os"
	"proevilz/api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file couln't be loaded!")
	}
	username := os.Getenv("username")
	password := os.Getenv("password")
	host := os.Getenv("host")
	dbName := os.Getenv("dbName")

	dsn := username + ":" + password + "@tcp(" + host + ":3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to configure database connection pool: %v", err)
	}

	// Set up connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Set up table relationships and constraints
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		return fmt.Errorf("failed to configure database table relationships: %v", err)
	}

	// Set the global DB variable to the database instance

	log.Println("Connected to the database")
	DB = db
	return nil
}
