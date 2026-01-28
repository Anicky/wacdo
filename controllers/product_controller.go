package controllers

import (
	"net/http"
	"os"
	"wacdo/config"
	"wacdo/models"
	"wacdo/utils"

	"github.com/gin-gonic/gin"
)

// GetProducts godoc
// @Description Récupérer tous les produits
// @Tags Products
// @Produce json
// @Success 200 {array} models.Product
// @Security BearerAuth
// @Router /products [get]
func GetProducts(context *gin.Context) {
	var products []models.Product

	if err := config.DB.Find(&products).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch products."})
		return
	}

	context.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Description Récupérer un produit par son ID
// @Tags Products
// @Produce json
// @Param id path int true "ID du produit"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Produit non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /products/{id} [get]
func GetProduct(context *gin.Context) {
	product, err := models.FindProductById(context)

	if err == nil {
		context.JSON(http.StatusOK, product)
	}
}

// PostProduct godoc
// @Description Créer un nouveau produit
// @Tags PostProducts
// @Accept json
// @Produce json
// @Param product body models.Product true "Données du produit"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /product [post]
func PostProduct(context *gin.Context) {
	var product models.Product

	if err := context.ShouldBindJSON(&product); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	path, err := utils.UploadImage(context)
	if err != nil {
		return
	}

	if path != nil {
		product.Image = *path
	}

	if err := config.DB.Create(&product).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create product."})

		return
	}

	context.JSON(http.StatusCreated, product)
}

// PutProduct godoc
// @Description Mettre à jour un produit existant
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "ID du produit"
// @Param input body models.ProductUpdateInput true "Données de mise à jour"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 404 {object} map[string]string "Produit non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /products/{id} [put]
func PutProduct(context *gin.Context) {
	product, err := models.FindProductById(context)

	if err == nil {
		var input models.ProductUpdateInput
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

		if input.Price != nil {
			updates["price"] = *input.Price
		}

		if input.IsAvailable != nil {
			updates["isAvailable"] = *input.IsAvailable
		}

		path, err := utils.UploadImage(context)
		if err != nil {
			return
		}

		if path != nil {
			if product.Image != "" {
				err = os.Remove(product.Image)

				if err != nil {
					context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete old image."})

					return
				}
			}

			updates["image"] = *path
		}

		if len(updates) == 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "No data to update."})

			return
		}

		if err := config.DB.Model(&product).Updates(updates).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update product."})

			return
		}

		context.JSON(http.StatusOK, product)
	}
}

// DeleteProduct godoc
// @Description Supprimer un produit
// @Tags Products
// @Produce json
// @Param id path int true "ID du produit"
// @Success 200 {object} map[string]string "Message de succès"
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Produit non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /products/{id} [delete]
func DeleteProduct(context *gin.Context) {
	product, err := models.FindProductById(context)

	if err == nil {
		if err = config.DB.Delete(&product).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete product."})

			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully."})
	}
}
