package models

import (
	"time"
)

type Order struct {
	ID           uint `gorm:"primaryKey"`
	Status       OrderStatus
	TicketNumber string
	Items        []OrderItem
	UserID       uint
	User         User
	CreatedAt    time.Time
	PreparedAt   time.Time
	DeliveredAt  time.Time
}
