package main

import (
	"context"
	"log"
	"wacdo/config"
	"wacdo/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// If .env file is not found, it is not necessarily an error.
		// With Render, environment variables are injected; there is no need for .env file.
		log.Print("Unable to find .env file.")
	}

	config.ConnectDB()

	migrateSchema()

	deleteTestData()

	createTestData()
}

func migrateSchema() {
	err := config.DB.AutoMigrate(&models.User{}, &models.Product{})
	if err != nil {
		log.Fatal("Unable to auto migrate: ", err)
	}
}

func deleteTestData() {
	// Remove users
	var users []models.User
	config.DB.Unscoped().Where("1 = 1").Find(&users)
	config.DB.Unscoped().Delete(&users)

	// Remove products
	var products []models.Product
	config.DB.Unscoped().Where("1 = 1").Find(&products)
	config.DB.Unscoped().Delete(&products)
}

func createTestData() {
	ctx := context.Background()

	// Create users
	err := gorm.G[models.User](config.DB).Create(ctx, &models.User{
		Email:    "admin@wacdo.com",
		Password: hashPassword("admin"),
	})
	if err != nil {
		log.Fatal("Unable to create user data: ", err)
	}

	// Create products
	err = gorm.G[models.Product](config.DB).Create(ctx, &models.Product{
		Name:        "BigWac",
		Description: "Le burger classique",
		Image:       "https://www.shutterstock.com/image-photo/delicious-double-cheeseburger-big-mac-600nw-2510772013.jpg",
		Price:       4.85,
		IsAvailable: true,
	})
	if err != nil {
		log.Fatal("Unable to create product data: ", err)
	}
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal("Unable to hash password: ", err)
	}

	return string(bytes)
}
