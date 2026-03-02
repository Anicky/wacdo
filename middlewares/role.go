package middlewares

import (
	"net/http"
	"slices"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

func CheckRole(roles []models.UserRole) gin.HandlerFunc {
	return func(context *gin.Context) {
		userID := GetUserId(context)

		user, err := models.FindUserById(context, *userID)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to get user."})

			return
		}

		if slices.Contains(roles, user.Role) == false {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access not allowed."})

			return
		}

		context.Next()
	}
}
