package utils

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database.")
	}

	return db
}
