package main

import (
	"context"
	"log"
	"wacdo/models"
	"wacdo/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	utils.InitEnvVars()

	db := utils.InitDatabase()

	migrateSchema(db)

	deleteTestData(db)

	createTestData(db)
}

func migrateSchema(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("Failed to migrate schema.")
	}
}

func deleteTestData(db *gorm.DB) {
	// Remove users
	var users []models.User
	db.Unscoped().Where("1 = 1").Find(&users)
	db.Unscoped().Delete(&users)
}

func createTestData(db *gorm.DB) {
	ctx := context.Background()

	// Create users
	err := gorm.G[models.User](db).Create(ctx, &models.User{Username: "admin", Password: hashPassword("admin")})

	if err != nil {
		log.Fatal("Failed to create user data.")
	}
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal("Failed to hash password.")
	}

	return string(bytes)
}
