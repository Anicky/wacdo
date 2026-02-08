package tests

import (
	"log"
	"net/http"
	"os"
	"time"
	"wacdo/config"
	"wacdo/controllers"
	"wacdo/models"
	"wacdo/routes"
	"wacdo/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitTest() *gin.Engine {
	gin.SetMode(gin.TestMode)

	err := os.Setenv("JWT_SECRET", "test_secret")
	if err != nil {
		log.Fatal("Unable to set environment variable: ", err)
	}

	config.DB = setupTestDatabase()

	router := gin.Default()

	routes.UserRoutes(router)
	routes.ProductCategoryRoutes(router)
	routes.ProductRoutes(router)
	routes.MenuRoutes(router)

	return router
}

func AuthenticateUser(request *http.Request) {
	token := generateTestToken(1)

	request.Header.Set("Authorization", "Bearer "+token)
}

func setupTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to setup database: ", err)
	}

	err = db.AutoMigrate(
		&models.ProductCategory{},
		&models.Product{},
		&models.Menu{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("Unable to migrate database: ", err)
	}

	// Products categories
	db.Create(&models.ProductCategory{Name: "Test product category 1", Description: "Test product category description 1"})
	db.Create(&models.ProductCategory{Name: "Test product category 2", Description: "Test product category description 2"})

	// Products
	product1 := &models.Product{Name: "Test product 1", Description: "Test product description 1", Price: 2.50, IsAvailable: true, CategoryID: 1}
	product2 := &models.Product{Name: "Test product 2", Description: "Test product description 2", Price: 4.99, IsAvailable: false, CategoryID: 2}
	product3 := &models.Product{Name: "Test product 3", Description: "Test product description 3", Price: 3.65, IsAvailable: true, CategoryID: 1}
	db.Create(product1)
	db.Create(product2)
	db.Create(product3)

	// Menus
	db.Create(&models.Menu{Name: "Test menu 1", Description: "Test menu description 1", Price: 8.54, IsAvailable: true, Products: []models.Product{*product1, *product2}})
	db.Create(&models.Menu{Name: "Test menu 2", Description: "Test menu description 2", Price: 7.20, IsAvailable: false, Products: []models.Product{*product1, *product3}})

	// Users
	db.Create(&models.User{Email: "admin1@example.com", Password: utils.HashPassword("Admin1234!"), Role: "admin"})
	db.Create(&models.User{Email: "greeter1@example.com", Password: utils.HashPassword("Greeter1234!"), Role: "greeter"})

	return db
}

func generateTestToken(userID uint) string {
	claim := &controllers.CustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString
}
