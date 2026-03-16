package routes

import (
	"wacdo/controllers"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoutes(router *gin.Engine) {
	routesGroup := router.Group("/authentication")

	{
		routesGroup.POST("/login", controllers.Authenticate)
	}
}
