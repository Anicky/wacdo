package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	routesGroup := router.Group("/users")

	{
		routesGroup.POST("/login", controllers.Login)

		// @TODO: add middleware for role (check that only admins can use these routes)
		routesGroup.GET("/", middlewares.Authentication(), controllers.GetUsers)
		routesGroup.GET("/:id", middlewares.Authentication(), controllers.GetUser)
		routesGroup.POST("/", middlewares.Authentication(), controllers.PostUser)
		routesGroup.PUT("/:id", middlewares.Authentication(), controllers.PutUser)
		routesGroup.DELETE("/:id", middlewares.Authentication(), controllers.DeleteUser)
	}
}
