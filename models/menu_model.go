package models

import "time"

type Menu struct {
	ID          uint      `gorm:"primaryKey"`
	Products    []Product `gorm:"many2many:menu_products"`
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Image       string
	Price       float32
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
