package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	routesGroup := router.Group("/users")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.GetUsers)
		routesGroup.GET("/:id", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.GetUser)
		routesGroup.POST("/", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.PostUser)
		routesGroup.PUT("/:id", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.PutUser)
		routesGroup.DELETE("/:id", middlewares.CheckRole([]models.UserRole{models.Admin}), controllers.DeleteUser)
	}
}
