package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProductCategoriesSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test product category 1")
	assert.Contains(testing, body, "Test product category description 1")
	assert.Contains(testing, body, "Test product category 2")
	assert.Contains(testing, body, "Test product category description 2")
}

func TestGetProductCategoriesUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)
}

func TestGetProductCategorySuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test product category 1")
	assert.Contains(testing, body, "Test product category description 1")
}

func TestGetProductCategoryUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)
}

func TestGetProductCategoryNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)
}

func TestGetProductCategoryInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/categories/a", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostProductCategorySuccess(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product category 3",
		"description": "Test product category description 3",
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/categories/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusCreated, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test product category 3")
	assert.Contains(testing, body, "Test product category description 3")
}

func TestPostProductCategoryError(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        1234,
		"description": 1234,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/categories/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostProductCategoryUnauthorized(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product category 3",
		"description": "Test product category description 3",
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/categories/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)
}

func TestPutProductCategorySuccess(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product category 1b",
		"description": "Test product category description 1b",
		"price":       3.29,
		"isAvailable": false,
		"categoryId":  1,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test product category 1b")
	assert.Contains(testing, body, "Test product category description 1b")
}

func TestPutProductCategoryError(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        1234,
		"description": 1234,
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutProductCategoryUnauthorized(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product category 1b",
		"description": "Test product category description 1b",
	}

	data, err := json.Marshal(product)
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

	assert.Equal(testing, http.StatusUnauthorized, response.Code)
}

func TestPutProductCategoryNotFound(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product category 1b",
		"description": "Test product category description 1b",
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/0", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)
}

func TestPutProductCategoryInvalidId(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product category 1b",
		"description": "Test product category description 1b",
	}

	data, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/products/categories/a", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestDeleteProductCategorySuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)
}

func TestDeleteProductCategoryUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)
}

func TestDeleteProductCategoryNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)
}

func TestDeleteProductCategoryInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/a", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}
