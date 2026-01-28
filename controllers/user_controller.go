package controllers

import (
	"net/http"
	"os"
	"time"
	"wacdo/config"
	"wacdo/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaim struct {
	UserID uint
	jwt.RegisteredClaims
}

// Login godoc
// @Description Se connecter (pour obtenir un token JWT)
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "Identifiants utilisateur (email, password)"
// @Success 200 {object} map[string]string "Token JWT"
// @Failure 400 {object} map[string]string "Identifiants invalides"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Router /users/login [post]
func Login(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password."})

		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password."})

		return
	}

	claim := &CustomClaim{
		UserID: existingUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token."})

		return
	}

	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}

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
	user, err := models.FindUserById(context)

	if err == nil {
		context.JSON(http.StatusOK, models.TransformUserToOutput(user))
	}
}
