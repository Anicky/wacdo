package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wacdo/models"

	"github.com/stretchr/testify/assert"
)

func TestGetProductsSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	var results []models.Product
	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, 4, len(results))

	assert.Equal(testing, "Test product 1", results[0].Name)
	assert.Equal(testing, "Test product description 1", results[0].Description)
	assert.Equal(testing, 2.5, results[0].Price)
	assert.True(testing, results[0].IsAvailable)

	assert.Equal(testing, "Test product 2", results[1].Name)
	assert.Equal(testing, "Test product description 2", results[1].Description)
	assert.Equal(testing, 4.99, results[1].Price)
	assert.False(testing, results[1].IsAvailable)

	assert.Equal(testing, "Test product 3", results[2].Name)
	assert.Equal(testing, "Test product description 3", results[2].Description)
	assert.Equal(testing, 3.65, results[2].Price)
	assert.True(testing, results[2].IsAvailable)

	assert.Equal(testing, "Test product 4", results[3].Name)
	assert.Equal(testing, "Test product description 4", results[3].Description)
	assert.Equal(testing, 9.10, results[3].Price)
	assert.True(testing, results[3].IsAvailable)
}

func TestGetProductsUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/", nil)
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

func TestGetProductSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.Product{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "Test product 1", result.Name)
	assert.Equal(testing, "Test product description 1", result.Description)
	assert.Equal(testing, 2.5, result.Price)
	assert.True(testing, result.IsAvailable)
}

func TestGetProductUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/1", nil)
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

func TestGetProductNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product 0: item not found.")
}

func TestGetProductInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/a", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Invalid ID.")
}

func TestPostProductSuccess(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product 5",
		"description": "Test product description 5",
		"price":       8.25,
		"isAvailable": true,
		"categoryId":  1,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusCreated, response.Code)

	result := models.Product{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "Test product 5", result.Name)
	assert.Equal(testing, "Test product description 5", result.Description)
	assert.Equal(testing, 8.25, result.Price)
	assert.True(testing, result.IsAvailable)
}

func TestPostProductInvalidCategory(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product 5",
		"description": "Test product description 5",
		"price":       8.25,
		"isAvailable": true,
		"categoryId":  9999,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product category not found.")
}

func TestPostProductError(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        1234,
		"description": 1234,
		"price":       "abcd",
		"isAvailable": "abcd",
		"categoryId":  "abcd",
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostProductUnauthorized(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product 5",
		"description": "Test product description 5",
		"price":       8.25,
		"isAvailable": true,
		"categoryId":  1,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Unauthorized.")
}

func TestPutProductSuccess(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product 1b",
		"description": "Test product description 1b",
		"price":       3.29,
		"isAvailable": false,
		"categoryId":  1,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.Product{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "Test product 1b", result.Name)
	assert.Equal(testing, "Test product description 1b", result.Description)
	assert.Equal(testing, 3.29, result.Price)
	assert.False(testing, result.IsAvailable)
}

func TestPutProductInvalidCategory(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product 1b",
		"description": "Test product description 1b",
		"price":       3.29,
		"isAvailable": false,
		"categoryId":  0,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product category not found.")
}

func TestPutProductError(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        1234,
		"description": 1234,
		"price":       "abcd",
		"isAvailable": "abcd",
		"categoryId":  "abcd",
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutProductUnauthorized(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product 1b",
		"description": "Test product description 1b",
		"price":       3.29,
		"isAvailable": false,
		"categoryId":  1,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Unauthorized.")
}

func TestPutProductNotFound(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product 1b",
		"description": "Test product description 1b",
		"price":       3.29,
		"isAvailable": false,
		"categoryId":  1,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/0", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product 0: item not found.")
}

func TestPutProductInvalidId(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product 1b",
		"description": "Test product description 1b",
		"price":       3.29,
		"isAvailable": false,
		"categoryId":  1,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/a", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Invalid ID.")
}

func TestDeleteProductSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/4", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)
}

func TestDeleteProductErrorAssociatedMenus(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Cannot delete product: there are menus associated with it.")
}

func TestDeleteProductUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/1", nil)
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

func TestDeleteProductNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product 0: item not found.")
}

func TestDeleteProductInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/a", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Invalid ID.")
}
