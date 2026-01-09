package models

type ProductCategory struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `binding:"required"`
	Description string `binding:"required"`
}
