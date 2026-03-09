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

func TestPutMenuSuccess(testing *testing.T) {
	router := tests.InitTest()

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

	tests.AuthenticateUserAsAdmin(request)

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
	router := tests.InitTest()

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

func TestPutMenuError(testing *testing.T) {
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

	request, err := http.NewRequest(http.MethodPut, "/menus/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutMenuUnauthorized(testing *testing.T) {
	router := tests.InitTest()

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

	tests.AssertUnauthorized(testing, response)
}

func TestPutMenuAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

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

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertAccessNotAllowed(testing, response)
}

func TestPutMenuNotFound(testing *testing.T) {
	router := tests.InitTest()

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

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 0: item not found.")
}

func TestPutMenuInvalidId(testing *testing.T) {
	router := tests.InitTest()

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

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Invalid ID.")
}
