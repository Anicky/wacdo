package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

func ProductCategoryRoutes(router *gin.Engine) {
	routesGroup := router.Group("/products/categories")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", controllers.GetProductsCategories)
		routesGroup.GET("/:id", controllers.GetProductCategory)
		routesGroup.POST("/", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.PostProductCategory)
		routesGroup.PUT("/:id", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.PutProductCategory)
		routesGroup.DELETE("/:id", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.DeleteProductCategory)
	}
}
