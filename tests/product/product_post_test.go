package product

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

func TestPostProductSuccess(testing *testing.T) {
	router := tests.InitTest()

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

	tests.AuthenticateUserAsAdmin(request)

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
	router := tests.InitTest()

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

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product category not found.")
}

func TestPostProductError(testing *testing.T) {
	router := tests.InitTest()

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

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostProductUnauthorized(testing *testing.T) {
	router := tests.InitTest()

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

	tests.AssertUnauthorized(testing, response)
}

func TestPostProductAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

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

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertAccessNotAllowed(testing, response)
}
