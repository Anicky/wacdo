package menu

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

func TestPostMenuSuccess(testing *testing.T) {
	router := tests.InitTest()

	menu := map[string]interface{}{
		"name":        "Test menu 3",
		"description": "Test menu description 3",
		"price":       7.80,
		"isAvailable": true,
		"productsIDs": []int{2, 3},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/menus/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusCreated, response.Code)

	result := models.Menu{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "Test menu 3", result.Name)
	assert.Equal(testing, "Test menu description 3", result.Description)
	assert.Equal(testing, 7.80, result.Price)
	assert.True(testing, result.IsAvailable)
	assert.Equal(testing, 2, len(result.Products))
	assert.Equal(testing, "Test product 2", result.Products[0].Name)
	assert.Equal(testing, "Test product description 2", result.Products[0].Description)
	assert.Equal(testing, 4.99, result.Products[0].Price)
	assert.False(testing, result.Products[0].IsAvailable)
	assert.Equal(testing, "Test product 3", result.Products[1].Name)
	assert.Equal(testing, "Test product description 3", result.Products[1].Description)
	assert.Equal(testing, 3.65, result.Products[1].Price)
	assert.True(testing, result.Products[1].IsAvailable)
}

func TestPostMenuInvalidProduct(testing *testing.T) {
	router := tests.InitTest()

	menu := map[string]interface{}{
		"name":        "Test menu 3",
		"description": "Test menu description 3",
		"price":       7.80,
		"isAvailable": true,
		"productsIDs": []int{0, 1, 9999},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/menus/", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Unable to find products.")
	assert.Contains(testing, body, "missing products")
	assert.Contains(testing, body, "[0,9999]")
}

func TestPostMenuError(testing *testing.T) {
	router := tests.InitTest()

	menu := map[string]interface{}{
		"name":        1234,
		"description": 1234,
		"price":       "abcd",
		"isAvailable": "abcd",
		"productsIDs": "abcd",
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/menus/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostMenuUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	menu := map[string]interface{}{
		"name":        "Test menu 3",
		"description": "Test menu description 3",
		"price":       7.80,
		"isAvailable": true,
		"productsIDs": []int{2, 3},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/menus/", bytes.NewBuffer(data))
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

func TestPostMenuAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	menu := map[string]interface{}{
		"name":        "Test menu 3",
		"description": "Test menu description 3",
		"price":       7.80,
		"isAvailable": true,
		"productsIDs": []int{2, 3},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/menus/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusForbidden, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Access not allowed.")
}
