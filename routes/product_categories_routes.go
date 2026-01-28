package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func ProductCategoryRoutes(router *gin.Engine) {
	routesGroup := router.Group("/products/categories")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", middlewares.Authentication(), controllers.GetProductsCategories)
		routesGroup.GET("/:id", middlewares.Authentication(), controllers.GetProductCategory)
	}
}
