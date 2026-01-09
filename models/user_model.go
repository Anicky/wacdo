package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique" binding:"required,email"`
	Password  string `binding:"required,min=8"`
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}
