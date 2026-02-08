package controllers

import (
	"net/http"
	"os"
	"wacdo/config"
	"wacdo/models"
	"wacdo/utils"

	"github.com/gin-gonic/gin"
)

// GetMenus godoc
// @Description Récupérer tous les menus
// @Tags Menus
// @Produce json
// @Success 200 {array} models.Menu
// @Security BearerAuth
// @Router /menus [get]
func GetMenus(context *gin.Context) {
	var menus []models.Menu

	if err := config.DB.Preload("Products").Find(&menus).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch menus."})
		return
	}

	context.JSON(http.StatusOK, menus)
}

// GetMenu godoc
// @Description Récupérer un menu par son ID
// @Tags Menus
// @Produce json
// @Param id path int true "ID du menu"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Menu non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /menus/{id} [get]
func GetMenu(context *gin.Context) {
	menu, err := models.FindMenuByContext(context)

	if err == nil {
		context.JSON(http.StatusOK, menu)
	}
}

// PostMenu godoc
// @Description Créer un nouveau menu
// @Tags Menus
// @Accept json
// @Produce json
// @Param menu body models.Menu true "Données du menu"
// @Success 201 {object} models.Menu
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /menu [post]
func PostMenu(context *gin.Context) {
	var input models.MenuInsertInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	products, _ := models.FindProductsById(context, input.ProductsIDs)
	if products == nil {
		return
	}

	menu := models.Menu{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		IsAvailable: input.IsAvailable,
		Products:    *products,
	}

	// @TODO: image

	if err := config.DB.Create(&menu).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create menu."})
		return
	}

	context.JSON(http.StatusCreated, menu)
}

// PutMenu godoc
// @Description Mettre à jour un menu existant
// @Tags Menus
// @Accept json
// @Produce json
// @Param id path int true "ID du menu"
// @Param input body models.MenuUpdateInput true "Données de mise à jour"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 404 {object} map[string]string "Produit non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /menus/{id} [put]
func PutMenu(context *gin.Context) {
	menu, err := models.FindMenuByContext(context)

	if err == nil {
		var input models.MenuUpdateInput
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
			if menu.Image != "" {
				err = os.Remove(menu.Image)

				if err != nil {
					context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete old image."})

					return
				}
			}

			updates["image"] = *path
		}

		var products *[]models.Product
		if input.ProductsIDs != nil {
			products, _ = models.FindProductsById(context, *input.ProductsIDs)
			if products == nil {
				return
			}
		}

		if len(updates) == 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "No data to update."})

			return
		}

		if err := config.DB.Model(&menu).Updates(updates).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update menu."})

			return
		}

		if input.ProductsIDs != nil {
			if err := config.DB.Model(&menu).Association("Products").Replace(products); err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update menu products."})

				return
			}
		}

		context.JSON(http.StatusOK, menu)
	}
}

// DeleteMenu godoc
// @Description Supprimer un menu
// @Tags Menus
// @Produce json
// @Param id path int true "ID du menu"
// @Success 200 {object} map[string]string "Message de succès"
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Produit non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /menus/{id} [delete]
func DeleteMenu(context *gin.Context) {
	menu, err := models.FindMenuByContext(context)

	if err == nil {
		if err := config.DB.Model(&menu).Association("Products").Replace([]models.Product{}); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete menu products."})

			return
		}

		if err = config.DB.Delete(&menu).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete menu."})

			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully."})
	}
}
