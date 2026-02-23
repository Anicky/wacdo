package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"wacdo/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

type OrderItemInput struct {
	Quantity  int  `json:"quantity" binding:"required,min=1"`
	ProductID uint `json:"productID"`
	MenuID    uint `json:"menuID"`
}

type OrderInsertInput struct {
	TicketNumber string           `json:"ticketNumber" binding:"required"`
	Items        []OrderItemInput `json:"items" binding:"required,min=1"`
}

type OrderUpdateInput struct {
	TicketNumber *string           `json:"ticketNumber"`
	Items        *[]OrderItemInput `json:"items" binding:"min=1"`
}

func FindOrderByContext(context *gin.Context) (order *Order, err error) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})

		return nil, err
	}

	return FindOrderById(context, uint(id))
}

func FindOrderById(context *gin.Context, id uint) (order *Order, err error) {
	if err = config.DB.Preload("User").Preload("Items").First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "Order not found."})

			return nil, err
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch order."})

		return nil, err
	}

	return order, nil
}
