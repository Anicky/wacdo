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

	body := response.Body.String()

	assert.Contains(testing, body, "Test product 1")
	assert.Contains(testing, body, "Test product description 1")
	assert.Contains(testing, body, "2.5")
	assert.Contains(testing, body, "true")
	assert.Contains(testing, body, "1")
	assert.Contains(testing, body, "Test product 2")
	assert.Contains(testing, body, "Test product description 2")
	assert.Contains(testing, body, "4.99")
	assert.Contains(testing, body, "false")
	assert.Contains(testing, body, "2")
}

func TestGetProducsUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/products/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)
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

	body := response.Body.String()

	assert.Contains(testing, body, "Test product 1")
	assert.Contains(testing, body, "Test product description 1")
	assert.Contains(testing, body, "2.5")
	assert.Contains(testing, body, "true")
	assert.Contains(testing, body, "1")
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
}

func TestPostProductSuccess(testing *testing.T) {
	router := InitTest()

	product := map[string]interface{}{
		"name":        "Test product 3",
		"description": "Test product description 3",
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

	body := response.Body.String()

	assert.Contains(testing, body, "Test product 3")
	assert.Contains(testing, body, "Test product description 3")
	assert.Contains(testing, body, "8.25")
	assert.Contains(testing, body, "true")
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
		"name":        "Test product 3",
		"description": "Test product description 3",
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

	body := response.Body.String()

	assert.Contains(testing, body, "Test product 1b")
	assert.Contains(testing, body, "Test product description 1b")
	assert.Contains(testing, body, "3.29")
	assert.Contains(testing, body, "false")
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
}

func TestDeleteProductSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)
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
}
