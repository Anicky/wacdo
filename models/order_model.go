package models

import (
	"time"
)

type Order struct {
	ID           uint        `gorm:"primaryKey"`
	Status       OrderStatus `binding:"required"`
	TicketNumber string
	Items        []OrderItem
	UserID       uint
	User         User `binding:"required"`
	CreatedAt    time.Time
	PreparedAt   time.Time
	DeliveredAt  time.Time
}
