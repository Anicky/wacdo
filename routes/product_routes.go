package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	routesGroup := router.Group("/products")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", middlewares.Authentication(), controllers.GetProducts)
		routesGroup.GET("/:id", middlewares.Authentication(), controllers.GetProduct)

		// @TODO: add middleware for role (check that only admins can use these routes)
		routesGroup.POST("/", controllers.PostProduct)
		routesGroup.PUT("/:id", controllers.PutProduct)
		routesGroup.DELETE("/:id", controllers.DeleteProduct)
	}
}
