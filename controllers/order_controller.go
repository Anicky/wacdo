package controllers

import (
	"net/http"
	"wacdo/config"
	"wacdo/middlewares"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

// GetOrders godoc
// @Description Récupérer toutes les commandes
// @Tags Orders
// @Produce json
// @Success 200 {array} models.Order
// @Security BearerAuth
// @Router /orders [get]
func GetOrders(context *gin.Context) {
	var orders []models.Order

	if err := config.DB.Preload("User").Find(&orders).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch orders."})
		return
	}

	context.JSON(http.StatusOK, orders)
}

// GetOrder godoc
// @Description Récupérerer une commande par son ID
// @Tags Orders
// @Produce json
// @Param id path int true "ID de la commande"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Commande non trouvée"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /orders/{id} [get]
func GetOrder(context *gin.Context) {
	order, err := models.FindOrderByContext(context)

	if err == nil {
		context.JSON(http.StatusOK, order)
	}
}

// PostOrder godoc
// @Description Créer une nouvelle commande
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Données de la commande"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /order [post]
func PostOrder(context *gin.Context) {
	var input models.OrderInsertInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	userID := *middlewares.GetUserId(context)
	user, err := models.FindUserById(context, userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to find user."})
		return
	}

	orderItems := models.TransformOrderItemInputsToOrderItems(context, input.Items)
	if orderItems == nil {
		return
	}

	order := models.Order{
		TicketNumber: input.TicketNumber,
		User:         *user,
		Status:       models.Created,
		Items:        *orderItems,
	}

	if err := config.DB.Create(&order).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create order."})
		return
	}

	context.JSON(http.StatusCreated, order)
}

// PutOrder godoc
// @Description Mettre à jour une commande existante
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "ID de la commande"
// @Param input body models.OrderUpdateInput true "Données de mise à jour"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 404 {object} map[string]string "Commande non trouvée"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /orders/{id} [put]
func PutOrder(context *gin.Context) {
	order, err := models.FindOrderByContext(context)

	if err == nil {
		if order.Status == models.Prepared || order.Status == models.Delivered {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be modified because it has already been prepared."})

			return
		}

		var input models.OrderUpdateInput
		if err = context.ShouldBindJSON(&input); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

			return
		}

		updates := make(map[string]interface{})

		if input.TicketNumber != nil {
			updates["ticketNumber"] = *input.TicketNumber
		}

		var orderItems *[]models.OrderItem
		if input.Items != nil {
			orderItems = models.TransformOrderItemInputsToOrderItems(context, *input.Items)
			if orderItems == nil {
				return
			}
		}

		if len(updates) == 0 && orderItems == nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "No data to update."})

			return
		}

		if err := config.DB.Model(&order).Updates(updates).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update order."})

			return
		}

		if input.Items != nil {
			if err := config.DB.Model(&order).Association("Items").Unscoped().Replace(orderItems); err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update order items."})

				return
			}
		}

		context.JSON(http.StatusOK, order)
	}
}

// PatchOrderInPreparation godoc
// @Description Indiquer que la commande est en préparation
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "ID de la commande"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 404 {object} map[string]string "Commande non trouvée"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /orders/{id}/in-preparation [patch]
func PatchOrderInPreparation(context *gin.Context) {
	order, err := models.FindOrderByContext(context)

	if err == nil {
		if order.Status == models.InPreparation {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Order is already in preparation."})

			return
		}

		if order.Status == models.Prepared {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Order is already prepared."})

			return
		}

		if order.Status == models.Delivered {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Order is already delivered."})

			return
		}

		updates := map[string]interface{}{
			"status": models.InPreparation,
		}

		if err := config.DB.Model(&order).Updates(updates).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update order."})

			return
		}

		context.JSON(http.StatusOK, order)
	}
}

// PatchOrderPrepared godoc
// @Description Indiquer que la commande a été préparée
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "ID de la commande"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 404 {object} map[string]string "Commande non trouvée"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /orders/{id}/prepared [patch]
func PatchOrderPrepared(context *gin.Context) {
	order, err := models.FindOrderByContext(context)

	if err == nil {
		if order.Status == models.Prepared {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Order is already prepared."})

			return
		}

		if order.Status == models.Created {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Order must be in preparation before it can be prepared."})

			return
		}

		if order.Status == models.Delivered {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Order is already delivered."})

			return
		}

		updates := map[string]interface{}{
			"status": models.Prepared,
		}

		if err := config.DB.Model(&order).Updates(updates).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update order."})

			return
		}

		context.JSON(http.StatusOK, order)
	}
}

// PatchOrderDelivered godoc
// @Description Indiquer que la commande a été livrée
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "ID de la commande"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 404 {object} map[string]string "Commande non trouvée"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /orders/{id}/delivered [patch]
func PatchOrderDelivered(context *gin.Context) {
	order, err := models.FindOrderByContext(context)

	if err == nil {
		if order.Status == models.Delivered {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Order is already delivered."})

			return
		}

		if order.Status != models.Prepared {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Order must be prepared before it can be delivered."})

			return
		}

		updates := map[string]interface{}{
			"status": models.Delivered,
		}

		if err := config.DB.Model(&order).Updates(updates).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update order."})

			return
		}

		context.JSON(http.StatusOK, order)
	}
}
