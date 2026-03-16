package controllers

import (
	"net/http"
	"wacdo/config"
	"wacdo/models"
	"wacdo/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers godoc
// @Description Récupérer tous les utilisateurs
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Security BearerAuth
// @Router /users [get]
func GetUsers(context *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch users."})
		return
	}

	context.JSON(http.StatusOK, models.TransformUsersToOutput(users))
}

// GetUser godoc
// @Description Récupérer un utilisateur par son ID
// @Tags Users
// @Produce json
// @Param id path int true "ID de l'utilisateur"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Utilisateur non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /users/{id} [get]
func GetUser(context *gin.Context) {
	user, err := models.FindUserByContext(context)

	if err == nil {
		context.JSON(http.StatusOK, models.TransformUserToOutput(user))
	}
}

// PostUser godoc
// @Description Créer un nouvel utilisateur
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserInsertInput true "Données de l'utilisateur"
// @Success 201 {object} models.UserOutput
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /users [post]
func PostUser(context *gin.Context) {
	var input models.UserInsertInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", input.Email).Count(&count)

	if count > 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email already used."})

		return
	}

	if err := utils.ValidatePassword(input.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password."})

		return
	}

	if !input.Role.IsValid() {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role."})

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password."})

		return
	}

	user := models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user."})
		return
	}

	context.JSON(http.StatusCreated, models.TransformUserToOutput(&user))
}

// PutUser godoc
// @Description Mettre à jour un utilisateur
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID de l'utilisateur"
// @Param user body models.UserUpdateInput true "Données de l'utilisateur à mettre à jour"
// @Success 200 {object} models.UserOutput
// @Failure 400 {object} map[string]string "Données ou ID invalides"
// @Failure 404 {object} map[string]string "Utilisateur non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /users/{id} [put]
func PutUser(context *gin.Context) {
	user, err := models.FindUserByContext(context)
	if err == nil {
		var input models.UserUpdateInput
		if err := context.ShouldBindJSON(&input); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

			return
		}

		updates := make(map[string]interface{})

		if input.Email != nil {
			var count int64
			config.DB.Model(&models.User{}).Where("id <> ? AND email = ?", user.ID, *input.Email).Count(&count)

			if count > 0 {
				context.JSON(http.StatusBadRequest, gin.H{"error": "Email already used."})

				return
			}

			updates["email"] = *input.Email
		}

		if input.Role != nil {
			if !input.Role.IsValid() {
				context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role."})

				return
			}

			updates["role"] = *input.Role
		}

		if input.Password != nil {
			if err := utils.ValidatePassword(*input.Password); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password."})

				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password."})

				return
			}

			updates["password"] = string(hashedPassword)
		}

		if len(updates) == 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "No data to update."})

			return
		}

		if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update user."})

			return
		}

		context.JSON(http.StatusOK, models.TransformUserToOutput(user))
	}
}

// DeleteUser godoc
// @Description Supprimer un utilisateur
// @Tags Users
// @Param id path int true "ID de l'utilisateur"
// @Success 204 "Pas de contenu"
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Utilisateur non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /users/{id} [delete]
func DeleteUser(context *gin.Context) {
	user, err := models.FindUserByContext(context)
	if err == nil {
		if err := config.DB.Delete(&user).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete user."})
			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "User deleted successfully."})
	}
}
