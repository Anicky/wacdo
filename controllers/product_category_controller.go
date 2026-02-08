package controllers

import (
	"net/http"
	"wacdo/config"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

// GetProductsCategories godoc
// @Description Récupérer toutes les catégories de produits
// @Tags ProductsCategories
// @Produce json
// @Success 200 {array} models.ProductCategory
// @Security BearerAuth
// @Router /products/categories [get]
func GetProductsCategories(context *gin.Context) {
	var productsCategories []models.ProductCategory

	if err := config.DB.Preload("Products").Find(&productsCategories).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch products categories."})
		return
	}

	context.JSON(http.StatusOK, productsCategories)
}

// GetProductCategory GetProduct godoc
// @Description Récupérer une catégorie de produit par son ID
// @Tags ProductsCategories
// @Produce json
// @Param id path int true "ID de la catégorie du produit"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Catégorie de produit non trouvée"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /products/categories/{id} [get]
func GetProductCategory(context *gin.Context) {
	productCategory, err := models.FindProductCategoryByContext(context)

	if err == nil {
		context.JSON(http.StatusOK, productCategory)
	}
}

// PostProductCategory godoc
// @Description Créer une nouvelle catégorie de produit
// @Tags ProductsCategories
// @Accept json
// @Produce json
// @Param productCategory body models.ProductCategory true "Données de la catégorie du produit"
// @Success 201 {object} models.ProductCategory
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /product/categories [post]
func PostProductCategory(context *gin.Context) {
	var input models.ProductCategoryInsertInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	productCategory := models.ProductCategory{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := config.DB.Create(&productCategory).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create product category."})

		return
	}

	context.JSON(http.StatusCreated, productCategory)
}

// PutProductCategory godoc
// @Description Mettre à jour une catégorie de produit existante
// @Tags ProductsCategories
// @Accept json
// @Produce json
// @Param id path int true "ID de la catégorie de produit"
// @Param input body models.ProductCategoryUpdateInput true "Données de mise à jour"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 404 {object} map[string]string "Catégorie de produit non trouvée"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /products/categories/{id} [put]
func PutProductCategory(context *gin.Context) {
	productCategory, err := models.FindProductCategoryByContext(context)

	if err == nil {
		var input models.ProductCategoryUpdateInput
		if err = context.ShouldBindJSON(&input); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

			return
		}

		updates := make(map[string]interface{})

		if input.Name != nil {
			updates["name"] = *input.Name
		}

		if input.Description != nil {
			updates["description"] = *input.Description
		}

		if len(updates) == 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "No data to update."})

			return
		}

		if err := config.DB.Model(&productCategory).Updates(updates).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update product category."})

			return
		}

		context.JSON(http.StatusOK, productCategory)
	}
}

// DeleteProductCategory godoc
// @Description Supprimer une catégorie de produit
// @Tags ProductsCategories
// @Produce json
// @Param id path int true "ID de la catégorie de produit"
// @Success 200 {object} map[string]string "Message de succès"
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Catégorie de produit non trouvée"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /products/categories/{id} [delete]
func DeleteProductCategory(context *gin.Context) {
	productCategory, err := models.FindProductCategoryByContext(context)

	if err == nil {
		if err = config.DB.Delete(&productCategory).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete product category."})

			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Product category deleted successfully."})
	}
}
