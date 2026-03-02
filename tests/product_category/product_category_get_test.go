package product_category

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wacdo/models"
	"wacdo/tests"

	"github.com/stretchr/testify/assert"
)

func TestGetProductCategoriesSuccess(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	var results []models.ProductCategory
	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, 3, len(results))

	assert.Equal(testing, "Test product category 1", results[0].Name)
	assert.Equal(testing, "Test product category description 1", results[0].Description)

	assert.Equal(testing, "Test product category 2", results[1].Name)
	assert.Equal(testing, "Test product category description 2", results[1].Description)

	assert.Equal(testing, "Test product category 3", results[2].Name)
	assert.Equal(testing, "Test product category description 3", results[2].Description)
}

func TestGetProductCategoriesUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Unauthorized.")
}

func TestGetProductCategorySuccess(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.ProductCategory{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "Test product category 1", result.Name)
	assert.Equal(testing, "Test product category description 1", result.Description)
}

func TestGetProductCategoryUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Unauthorized.")
}

func TestGetProductCategoryNotFound(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product category not found.")
}

func TestGetProductCategoryInvalidId(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/a", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Invalid ID.")
}
