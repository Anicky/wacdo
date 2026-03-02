package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	routesGroup := router.Group("/menus")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", controllers.GetMenus)
		routesGroup.GET("/:id", controllers.GetMenu)
		routesGroup.POST("/", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.PostMenu)
		routesGroup.PUT("/:id", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.PutMenu)
		routesGroup.DELETE("/:id", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.DeleteMenu)
	}
}
