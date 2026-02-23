package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	routesGroup := router.Group("/orders")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", middlewares.Authentication(), controllers.GetOrders)
		routesGroup.GET("/:id", middlewares.Authentication(), controllers.GetOrder)
		routesGroup.POST("/", controllers.PostOrder)
		routesGroup.PUT("/:id", controllers.PutOrder)
		routesGroup.PATCH("/:id/in-preparation", controllers.PatchOrderInPreparation)
		routesGroup.PATCH("/:id/prepared", controllers.PatchOrderPrepared)
		routesGroup.PATCH("/:id/delivered", controllers.PatchOrderDelivered)
	}
}
