package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	routesGroup := router.Group("/menus")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", middlewares.Authentication(), controllers.GetMenus)
		routesGroup.GET("/:id", middlewares.Authentication(), controllers.GetMenu)

		// @TODO: add middleware for role (check that only admins can use these routes)
		routesGroup.POST("/", controllers.PostMenu)
		routesGroup.PUT("/:id", controllers.PutMenu)
		routesGroup.DELETE("/:id", controllers.DeleteMenu)
	}
}
