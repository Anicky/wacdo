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

	if err := config.DB.Find(&productsCategories).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch products categories."})
		return
	}

	context.JSON(http.StatusOK, productsCategories)
}

// GetProduct godoc
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
	productCategory, err := models.FindProductCategoryById(context)

	if err == nil {
		context.JSON(http.StatusOK, productCategory)
	}
}
