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

func TestPostProductCategorySuccess(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        "Test product category 4",
		"description": "Test product category description 4",
	}

	data, err := json.Marshal(productCategory)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/categories/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusCreated, response.Code)

	result := models.ProductCategory{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "Test product category 4", result.Name)
	assert.Equal(testing, "Test product category description 4", result.Description)
}

func TestPostProductCategoryError(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        1234,
		"description": 1234,
	}

	data, err := json.Marshal(productCategory)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/categories/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostProductCategoryUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        "Test product category 4",
		"description": "Test product category description 4",
	}

	data, err := json.Marshal(productCategory)
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

	tests.AssertUnauthorized(testing, response)
}

func TestPostProductCategoryAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	productCategory := map[string]interface{}{
		"name":        "Test product category 4",
		"description": "Test product category description 4",
	}

	data, err := json.Marshal(productCategory)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/products/categories/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertAccessNotAllowed(testing, response)
}
