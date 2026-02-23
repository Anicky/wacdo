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
	routes.OrderRoutes(router)

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

	if err = db.Exec("PRAGMA foreign_keys = ON", nil).Error; err != nil {
		log.Fatal("Unable to enforce foreign keys on database: ", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.ProductCategory{},
		&models.Product{},
		&models.Menu{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal("Unable to migrate database: ", err)
	}

	// Products categories
	productCategory1 := &models.ProductCategory{Name: "Test product category 1", Description: "Test product category description 1"}
	productCategory2 := &models.ProductCategory{Name: "Test product category 2", Description: "Test product category description 2"}
	db.Create(productCategory1)
	db.Create(productCategory2)
	db.Create(&models.ProductCategory{Name: "Test product category 3", Description: "Test product category description 3"})

	// Products
	product1 := &models.Product{Name: "Test product 1", Description: "Test product description 1", Price: 2.50, IsAvailable: true, Category: *productCategory1}
	product2 := &models.Product{Name: "Test product 2", Description: "Test product description 2", Price: 4.99, IsAvailable: false, Category: *productCategory2}
	product3 := &models.Product{Name: "Test product 3", Description: "Test product description 3", Price: 3.65, IsAvailable: true, Category: *productCategory1}
	db.Create(product1)
	db.Create(product2)
	db.Create(product3)
	db.Create(&models.Product{Name: "Test product 4", Description: "Test product description 4", Price: 9.10, IsAvailable: true, Category: *productCategory1})

	// Menus
	menu1 := &models.Menu{Name: "Test menu 1", Description: "Test menu description 1", Price: 8.54, IsAvailable: true, Products: []models.Product{*product1, *product2}}
	menu2 := &models.Menu{Name: "Test menu 2", Description: "Test menu description 2", Price: 7.20, IsAvailable: false, Products: []models.Product{*product1, *product3}}
	db.Create(menu1)
	db.Create(menu2)

	// Users
	db.Create(&models.User{Email: "admin1@example.com", Password: utils.HashPassword("Admin1234!"), Role: "admin"})

	userGreeter1 := &models.User{Email: "greeter1@example.com", Password: utils.HashPassword("Greeter1234!"), Role: "greeter"}
	db.Create(userGreeter1)

	userGreeter2 := &models.User{Email: "greeter2@example.com", Password: utils.HashPassword("Greeter5678!"), Role: "greeter"}
	db.Create(userGreeter2)

	// Orders
	orderItem1 := &models.OrderItem{Quantity: 2, OrderContentName: product1.Name, OrderContentDescription: product1.Description, OrderContentImage: product1.Image, OrderContentPrice: product1.Price}
	orderItem2 := &models.OrderItem{Quantity: 1, OrderContentName: menu1.Name, OrderContentDescription: menu1.Description, OrderContentImage: menu1.Image, OrderContentPrice: menu1.Price}
	db.Create(&models.Order{Status: models.Created, TicketNumber: "001", User: *userGreeter1, Items: []models.OrderItem{*orderItem1, *orderItem2}})

	orderItem3 := &models.OrderItem{Quantity: 1, OrderContentName: product2.Name, OrderContentDescription: product2.Description, OrderContentImage: product2.Image, OrderContentPrice: product2.Price}
	db.Create(&models.Order{Status: models.InPreparation, TicketNumber: "002", User: *userGreeter2, Items: []models.OrderItem{*orderItem3}})

	orderItem4 := &models.OrderItem{Quantity: 1, OrderContentName: menu2.Name, OrderContentDescription: menu2.Description, OrderContentImage: menu2.Image, OrderContentPrice: menu2.Price}
	db.Create(&models.Order{Status: models.Prepared, TicketNumber: "003", User: *userGreeter1, Items: []models.OrderItem{*orderItem4}})

	orderItem5 := &models.OrderItem{Quantity: 2, OrderContentName: product3.Name, OrderContentDescription: product3.Description, OrderContentImage: product3.Image, OrderContentPrice: product3.Price}
	db.Create(&models.Order{Status: models.Delivered, TicketNumber: "004", User: *userGreeter2, Items: []models.OrderItem{*orderItem5}})

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
