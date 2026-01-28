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
		log.Fatalf("Unable to find .env file.")
	}

	config.ConnectDB()

	migrateSchema()

	deleteTestData()

	createTestData()
}

func migrateSchema() {
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.ProductCategory{},
	)
	if err != nil {
		log.Fatal("Unable to auto migrate: ", err)
	}
}

func deleteTestData() {
	// Remove users
	var users []models.User
	config.DB.Unscoped().Where("1 = 1").Find(&users)
	config.DB.Unscoped().Delete(&users)

	// Remove products categories
	var productsCategories []models.ProductCategory
	config.DB.Unscoped().Where("1 = 1").Find(&productsCategories)
	config.DB.Unscoped().Delete(&productsCategories)
}

func createTestData() {
	ctx := context.Background()

	// Create users
	err := gorm.G[models.User](config.DB).Create(ctx, &models.User{
		Email:    "admin@wacdo.com",
		Password: hashPassword("Admin1234!"),
		Role:     "admin",
	})
	if err != nil {
		log.Fatal("Unable to create user data: ", err)
	}

	// Create products categories
	err = gorm.G[models.ProductCategory](config.DB).Create(ctx, &models.ProductCategory{
		Name:        "Burgers",
		Description: "Les burgers",
	})
	if err != nil {
		log.Fatal("Unable to create product category data: ", err)
	}
	err = gorm.G[models.ProductCategory](config.DB).Create(ctx, &models.ProductCategory{
		Name:        "Frites",
		Description: "Les frites",
	})
	if err != nil {
		log.Fatal("Unable to create product category data: ", err)
	}
	err = gorm.G[models.ProductCategory](config.DB).Create(ctx, &models.ProductCategory{
		Name:        "Boissons",
		Description: "Les boissons",
	})
	if err != nil {
		log.Fatal("Unable to create product category data: ", err)
	}
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal("Unable to hash password: ", err)
	}

	return string(bytes)
}
