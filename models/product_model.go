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

type Product struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Image       string
	Price       float32
	IsAvailable bool
	CategoryID  uint
	Category    ProductCategory `gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductInsertInput struct {
	Name        *string  `json:"name" binding:"required"`
	Description *string  `json:"description" binding:"required"`
	Price       *float32 `json:"price" binding:"required"`
	IsAvailable *bool    `json:"isAvailable" binding:"required"`
	CategoryID  *uint    `json:"categoryID" binding:"required"`
	// @TODO: image
}

type ProductUpdateInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float32 `json:"price"`
	IsAvailable *bool    `json:"isAvailable"`
	CategoryID  *uint    `json:"categoryID"`
	// @TODO: image
}

func FindProductById(context *gin.Context) (product *Product, err error) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})

		return nil, err
	}

	if err = config.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "Product not found."})

			return nil, err
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch product."})

		return nil, err
	}

	return product, nil
}
