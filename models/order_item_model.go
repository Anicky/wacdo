package models

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderItem struct {
	ID                      uint `gorm:"primaryKey"`
	OrderID                 uint
	Quantity                int
	OrderContentName        string
	OrderContentDescription string
	OrderContentImage       string
	OrderContentPrice       float64
}

func TransformOrderItemInputsToOrderItems(context *gin.Context, items []OrderItemInput) *[]OrderItem {
	var orderItems []OrderItem

	for _, item := range items {
		if item.ProductID != 0 {
			product, _ := FindProductById(context, item.ProductID)
			if product == nil {
				return nil
			}

			if product.IsAvailable == false {
				context.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Product %d: item is not available.", item.ProductID)})

				return nil
			}

			if item.Quantity <= 0 {
				context.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Product %d: item quantity should be superior than 0.", item.ProductID)})

				return nil
			}

			orderItems = append(orderItems, OrderItem{
				Quantity:                item.Quantity,
				OrderContentName:        product.Name,
				OrderContentDescription: product.Description,
				OrderContentImage:       product.Image,
				OrderContentPrice:       product.Price,
			})
		} else if item.MenuID != 0 {
			menu, _ := FindMenuById(context, item.MenuID)
			if menu == nil {
				return nil
			}

			if menu.IsAvailable == false {
				context.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Menu %d: item is not available.", item.MenuID)})

				return nil
			}

			if item.Quantity <= 0 {
				context.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Menu %d: item quantity should be superior than 0.", item.MenuID)})

				return nil
			}

			orderItems = append(orderItems, OrderItem{
				Quantity:                item.Quantity,
				OrderContentName:        menu.Name,
				OrderContentDescription: menu.Description,
				OrderContentImage:       menu.Image,
				OrderContentPrice:       menu.Price,
			})
		}
	}

	if orderItems == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "At least one item must be provided."})

		return nil
	}

	return &orderItems
}
