package product_category

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wacdo/models"
	"wacdo/tests"

	"github.com/stretchr/testify/assert"
)

func TestPutProductCategorySuccess(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        "Test product category 1b",
		"description": "Test product category description 1b",
		"price":       3.29,
		"isAvailable": false,
		"categoryId":  1,
	}

	data, err := json.Marshal(productCategory)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.ProductCategory{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "Test product category 1b", result.Name)
	assert.Equal(testing, "Test product category description 1b", result.Description)
}

func TestPutProductCategoryError(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        1234,
		"description": 1234,
	}

	data, err := json.Marshal(productCategory)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutProductCategoryUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        "Test product category 1b",
		"description": "Test product category description 1b",
	}

	data, err := json.Marshal(productCategory)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertUnauthorized(testing, response)
}

func TestPutProductCategoryAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        "Test product category 1b",
		"description": "Test product category description 1b",
	}

	data, err := json.Marshal(productCategory)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertAccessNotAllowed(testing, response)
}

func TestPutProductCategoryNotFound(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        "Test product category 1b",
		"description": "Test product category description 1b",
	}

	data, err := json.Marshal(productCategory)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/0", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product category not found.")
}

func TestPutProductCategoryInvalidId(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        "Test product category 1b",
		"description": "Test product category description 1b",
	}

	data, err := json.Marshal(productCategory)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/a", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Invalid ID.")
}
