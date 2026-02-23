package models

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
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
	Price       float64
	IsAvailable bool
	CategoryID  uint
	Category    ProductCategory `gorm:"foreignKey:CategoryID"`
	Menus       []Menu          `gorm:"many2many:menu_products"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductInsertInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	IsAvailable bool    `json:"isAvailable" binding:"required"`
	CategoryID  uint    `json:"categoryID" binding:"required"`
	// @TODO: image
}

type ProductUpdateInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	IsAvailable *bool    `json:"isAvailable"`
	CategoryID  *uint    `json:"categoryID"`
	// @TODO: image
}

func FindProductByContext(context *gin.Context) (product *Product, err error) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})

		return nil, err
	}

	return FindProductById(context, uint(id))
}

func FindProductById(context *gin.Context, id uint) (product *Product, err error) {
	if err = config.DB.Preload("Category").Preload("Menus").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product %d: item not found.", id)})

			return nil, err
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Product %d: unable to fetch item.", id)})

		return nil, err
	}

	return product, nil
}

func FindProductsById(context *gin.Context, productsIDs []uint) (products *[]Product, err error) {
	if err = config.DB.Find(&products, productsIDs).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch products."})

		return nil, err
	}

	if (len(*products)) != len(productsIDs) {
		missingProducts := make([]uint, 0)

		foundProductsIds := make([]uint, 0)
		for _, product := range *products {
			foundProductsIds = append(foundProductsIds, product.ID)
		}

		for _, productID := range productsIDs {
			if slices.Contains(foundProductsIds, productID) == false {
				missingProducts = append(missingProducts, productID)
			}
		}

		context.JSON(http.StatusNotFound, gin.H{
			"error":            "Unable to find products.",
			"missing products": missingProducts,
		})

		return nil, nil
	}

	return products, nil
}
