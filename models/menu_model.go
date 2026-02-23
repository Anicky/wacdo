package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"wacdo/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Menu struct {
	ID          uint      `gorm:"primaryKey"`
	Products    []Product `gorm:"many2many:menu_products"`
	Name        string
	Description string
	Image       string
	Price       float64
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MenuInsertInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	IsAvailable bool    `json:"isAvailable" binding:"required"`
	ProductsIDs []uint  `json:"productsIDs" binding:"required"`
	// @TODO: image
}

type MenuUpdateInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	IsAvailable *bool    `json:"isAvailable"`
	ProductsIDs *[]uint  `json:"productsIDs"`
	// @TODO: image
}

func FindMenuByContext(context *gin.Context) (menu *Menu, err error) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})

		return nil, err
	}

	return FindMenuById(context, uint(id))
}

func FindMenuById(context *gin.Context, id uint) (menu *Menu, err error) {
	if err = config.DB.Preload("Products").First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Menu %d: item not found.", id)})

			return nil, err
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Menu %d: unable to fetch item.", id)})

		return nil, err
	}

	return menu, nil
}
