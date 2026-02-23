package models

import (
	"errors"
	"net/http"
	"strconv"
	"wacdo/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductCategory struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Products    []Product `gorm:"foreignKey:CategoryID"`
}

type ProductCategoryInsertInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ProductCategoryUpdateInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func FindProductCategoryByContext(context *gin.Context) (productCategory *ProductCategory, err error) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})

		return nil, err
	}

	return FindProductCategoryById(context, uint(id))
}

func FindProductCategoryById(context *gin.Context, id uint) (productCategory *ProductCategory, err error) {
	if err = config.DB.Preload("Products").First(&productCategory, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "Product category not found."})

			return nil, err
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch product category."})

		return nil, err
	}

	return productCategory, nil
}
