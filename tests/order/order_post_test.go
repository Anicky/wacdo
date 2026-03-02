package order

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

func TestPostOrderSuccess(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "005",
		"items": []map[string]interface{}{
			{
				"quantity": 1,
				"menuID":   1,
			},
			{
				"quantity":  2,
				"productID": 3,
			},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusCreated, response.Code)

	result := models.Order{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "005", result.TicketNumber)
	assert.Equal(testing, models.Created, result.Status)
	assert.Equal(testing, "admin1@example.com", result.User.Email)

	assert.Equal(testing, 2, len(result.Items))

	assert.Equal(testing, 1, result.Items[0].Quantity)
	assert.Equal(testing, "Test menu 1", result.Items[0].OrderContentName)
	assert.Equal(testing, "Test menu description 1", result.Items[0].OrderContentDescription)
	assert.Equal(testing, 8.54, result.Items[0].OrderContentPrice)

	assert.Equal(testing, 2, result.Items[1].Quantity)
	assert.Equal(testing, "Test product 3", result.Items[1].OrderContentName)
	assert.Equal(testing, "Test product description 3", result.Items[1].OrderContentDescription)
	assert.Equal(testing, 3.65, result.Items[1].OrderContentPrice)
}

func TestPostOrderNoItems(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "005",
		"items":        []map[string]interface{}{{}},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "At least one item must be provided.")
}

func TestPostOrderInvalidProduct(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "005",
		"items": []map[string]interface{}{
			{
				"quantity": 1,
				"menuID":   1,
			},
			{
				"quantity":  2,
				"productID": 9999,
			},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Product 9999: item not found.")
}

func TestPostOrderInvalidMenu(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "005",
		"items": []map[string]interface{}{
			{
				"quantity": 1,
				"menuID":   9999,
			},
			{
				"quantity":  2,
				"productID": 1,
			},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Menu 9999: item not found.")
}

func TestPostOrderInvalidQuantity(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "005",
		"items": []map[string]interface{}{
			{
				"quantity": 0,
				"menuID":   1,
			},
			{
				"quantity":  1,
				"productID": 1,
			},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Menu 1: item quantity should be superior than 0.")
}

func TestPostOrderProductNotAvailable(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "005",
		"items": []map[string]interface{}{
			{
				"quantity": 1,
				"menuID":   1,
			},
			{
				"quantity":  1,
				"productID": 2,
			},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Product 2: item is not available.")
}

func TestPostOrderMenuNotAvailable(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "005",
		"items": []map[string]interface{}{
			{
				"quantity": 1,
				"menuID":   2,
			},
			{
				"quantity":  1,
				"productID": 1,
			},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Menu 2: item is not available.")
}

func TestPostOrderError(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": 1234,
		"items":        "abcd",
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostOrderUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "005",
		"items": []map[string]interface{}{
			{
				"quantity": 1,
				"menuID":   1,
			},
			{
				"quantity":  1,
				"productID": 1,
			},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
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

func TestPostOrderAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "005",
		"items": []map[string]interface{}{
			{
				"quantity": 1,
				"menuID":   1,
			},
			{
				"quantity":  1,
				"productID": 1,
			},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/orders/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUser(request, 4)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusForbidden, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Access not allowed.")
}
