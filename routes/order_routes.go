package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	routesGroup := router.Group("/orders")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", middlewares.CheckRole([]models.UserRole{models.Admin, models.OrderPicker, models.Manager}), controllers.GetOrders)
		routesGroup.GET("/:id", middlewares.CheckRole([]models.UserRole{models.Admin, models.OrderPicker, models.Manager}), controllers.GetOrder)
		routesGroup.POST("/", middlewares.CheckRole([]models.UserRole{models.Admin, models.Greeter, models.Manager}), controllers.PostOrder)
		routesGroup.PUT("/:id", middlewares.CheckRole([]models.UserRole{models.Admin, models.Greeter, models.Manager}), controllers.PutOrder)
		routesGroup.PATCH("/:id/in-preparation", middlewares.CheckRole([]models.UserRole{models.Admin, models.OrderPicker}), controllers.PatchOrderInPreparation)
		routesGroup.PATCH("/:id/prepared", middlewares.CheckRole([]models.UserRole{models.Admin, models.OrderPicker}), controllers.PatchOrderPrepared)
		routesGroup.PATCH("/:id/delivered", middlewares.CheckRole([]models.UserRole{models.Admin, models.Manager, models.Greeter}), controllers.PatchOrderDelivered)
	}
}
