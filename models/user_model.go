package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"wacdo/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique" binding:"required,email"`
	Password  string `binding:"required,min=8"`
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserUpdateInput struct {
	Email    *string   `json:"email"`
	Password *string   `json:"password"`
	Role     *UserRole `json:"role"`
}

type UserOutput struct {
	ID        uint
	Email     string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindUserByContext(context *gin.Context) (user *User, err error) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})

		return nil, err
	}

	return FindUserById(context, uint(id))
}

func FindUserById(context *gin.Context, id uint) (user *User, err error) {
	if err = config.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "User not found."})

			return nil, err
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user."})

		return nil, err
	}

	return user, nil
}

func TransformUsersToOutput(users []User) []UserOutput {
	var outputUsers []UserOutput

	for _, user := range users {
		outputUsers = append(outputUsers, TransformUserToOutput(&user))
	}

	return outputUsers
}

func TransformUserToOutput(user *User) UserOutput {
	return UserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
