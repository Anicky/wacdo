package models

import "time"

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Image       string
	Price       float32
	IsAvailable bool
	CategoryID  uint
	Category    ProductCategory
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
