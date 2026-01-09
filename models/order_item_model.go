package models

type OrderItem struct {
	ID           uint `gorm:"primaryKey"`
	OrderID      uint
	Quantity     int
	OrderContent string
}
