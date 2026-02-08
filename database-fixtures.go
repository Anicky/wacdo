package main

import (
	"context"
	"log"
	"wacdo/config"
	"wacdo/models"
	"wacdo/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to find .env file.")
	}

	config.ConnectDB()

	migrateSchema()

	deleteTestData()

	createTestData()
}

func migrateSchema() {
	err := config.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Unable to auto migrate: ", err)
	}
}

func deleteTestData() {
	// Remove users
	var users []models.User
	config.DB.Unscoped().Where("1 = 1").Find(&users)
	config.DB.Unscoped().Delete(&users)
}

func createTestData() {
	ctx := context.Background()

	// Create users
	err := gorm.G[models.User](config.DB).Create(ctx, &models.User{
		Email:    "admin@wacdo.com",
		Password: utils.HashPassword("Admin1234!"),
		Role:     "admin",
	})
	if err != nil {
		log.Fatal("Unable to create user data: ", err)
	}
}
