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

func TestGetMenusSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/menus/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	var results []models.Menu
	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, 2, len(results))

	assert.Equal(testing, "Test menu 1", results[0].Name)
	assert.Equal(testing, "Test menu description 1", results[0].Description)
	assert.Equal(testing, 8.54, results[0].Price)
	assert.True(testing, results[0].IsAvailable)
	assert.Equal(testing, 2, len(results[0].Products))
	assert.Equal(testing, "Test product 1", results[0].Products[0].Name)
	assert.Equal(testing, "Test product description 1", results[0].Products[0].Description)
	assert.Equal(testing, 2.5, results[0].Products[0].Price)
	assert.True(testing, results[0].Products[0].IsAvailable)
	assert.Equal(testing, "Test product 2", results[0].Products[1].Name)
	assert.Equal(testing, "Test product description 2", results[0].Products[1].Description)
	assert.Equal(testing, 4.99, results[0].Products[1].Price)
	assert.False(testing, results[0].Products[1].IsAvailable)

	assert.Equal(testing, "Test menu 2", results[1].Name)
	assert.Equal(testing, "Test menu description 2", results[1].Description)
	assert.Equal(testing, 7.20, results[1].Price)
	assert.False(testing, results[1].IsAvailable)
	assert.Equal(testing, 2, len(results[1].Products))
	assert.Equal(testing, "Test product 1", results[1].Products[0].Name)
	assert.Equal(testing, "Test product description 1", results[1].Products[0].Description)
	assert.Equal(testing, 2.5, results[1].Products[0].Price)
	assert.True(testing, results[1].Products[0].IsAvailable)
	assert.Equal(testing, "Test product 3", results[1].Products[1].Name)
	assert.Equal(testing, "Test product description 3", results[1].Products[1].Description)
	assert.Equal(testing, 3.65, results[1].Products[1].Price)
	assert.True(testing, results[1].Products[1].IsAvailable)
}

func TestGetMenusUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/menus/", nil)
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

func TestGetMenuSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/menus/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.Menu{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "Test menu 1", result.Name)
	assert.Equal(testing, "Test menu description 1", result.Description)
	assert.Equal(testing, 8.54, result.Price)
	assert.True(testing, result.IsAvailable)
	assert.Equal(testing, 2, len(result.Products))
	assert.Equal(testing, "Test product 1", result.Products[0].Name)
	assert.Equal(testing, "Test product description 1", result.Products[0].Description)
	assert.Equal(testing, 2.5, result.Products[0].Price)
	assert.True(testing, result.Products[0].IsAvailable)
	assert.Equal(testing, "Test product 2", result.Products[1].Name)
	assert.Equal(testing, "Test product description 2", result.Products[1].Description)
	assert.Equal(testing, 4.99, result.Products[1].Price)
	assert.False(testing, result.Products[1].IsAvailable)
}

func TestGetMenuUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/menus/1", nil)
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

func TestGetMenuNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/menus/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 0: item not found.")
}

func TestGetMenuInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/menus/a", nil)
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

func TestPostMenuSuccess(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

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
	router := InitTest()

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

	AuthenticateUser(request)

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
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostMenuUnauthorized(testing *testing.T) {
	router := InitTest()

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

func TestPutMenuSuccess(testing *testing.T) {
	router := InitTest()

	menu := map[string]interface{}{
		"name":        "Test menu 1b",
		"description": "Test menu description 1b",
		"price":       8.21,
		"isAvailable": false,
		"productsIDs": []int{3},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/menus/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.Menu{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "Test menu 1b", result.Name)
	assert.Equal(testing, "Test menu description 1b", result.Description)
	assert.Equal(testing, 8.21, result.Price)
	assert.False(testing, result.IsAvailable)
	assert.Equal(testing, 1, len(result.Products))
	assert.Equal(testing, "Test product 3", result.Products[0].Name)
	assert.Equal(testing, "Test product description 3", result.Products[0].Description)
	assert.Equal(testing, 3.65, result.Products[0].Price)
	assert.True(testing, result.Products[0].IsAvailable)
}

func TestPutMenuInvalidCategory(testing *testing.T) {
	router := InitTest()

	menu := map[string]interface{}{
		"name":        "Test menu 1b",
		"description": "Test menu description 1b",
		"price":       8.21,
		"isAvailable": false,
		"productsIDs": []int{0, 1, 9999},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/menus/1", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Unable to find products.")
	assert.Contains(testing, body, "missing products")
	assert.Contains(testing, body, "[0,9999]")
}

func TestPutMenuError(testing *testing.T) {
	router := InitTest()

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

	request, err := http.NewRequest(http.MethodPut, "/menus/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutMenuUnauthorized(testing *testing.T) {
	router := InitTest()

	menu := map[string]interface{}{
		"name":        "Test menu 1b",
		"description": "Test menu description 1b",
		"price":       8.21,
		"isAvailable": false,
		"productsIDs": []int{3},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/menus/1", bytes.NewBuffer(data))
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

func TestPutMenuNotFound(testing *testing.T) {
	router := InitTest()

	menu := map[string]interface{}{
		"name":        "Test menu 1b",
		"description": "Test menu description 1b",
		"price":       8.21,
		"isAvailable": false,
		"productsIDs": []int{3},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/menus/0", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Menu 0: item not found.")
}

func TestPutMenuInvalidId(testing *testing.T) {
	router := InitTest()

	menu := map[string]interface{}{
		"name":        "Test menu 1b",
		"description": "Test menu description 1b",
		"price":       8.21,
		"isAvailable": false,
		"productsIDs": []int{3},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/menus/a", bytes.NewBuffer(data))
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

func TestDeleteMenuSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/menus/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)
}

func TestDeleteMenuUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/menus/1", nil)
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

func TestDeleteMenuNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/menus/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 0: item not found.")
}

func TestDeleteMenuInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/menus/a", nil)
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
