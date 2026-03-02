package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	routesGroup := router.Group("/products")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", controllers.GetProducts)
		routesGroup.GET("/:id", controllers.GetProduct)
		routesGroup.POST("/", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.PostProduct)
		routesGroup.PUT("/:id", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.PutProduct)
		routesGroup.DELETE("/:id", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.DeleteProduct)
	}
}
