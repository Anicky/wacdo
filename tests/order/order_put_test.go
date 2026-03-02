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

func TestPutOrderSuccess(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
		"items": []map[string]interface{}{
			{
				"quantity": 2,
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

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.Order{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "006", result.TicketNumber)
	assert.Equal(testing, models.Created, result.Status)
	assert.Equal(testing, "greeter1@example.com", result.User.Email)

	assert.Equal(testing, 2, len(result.Items))

	assert.Equal(testing, 2, result.Items[0].Quantity)
	assert.Equal(testing, "Test menu 1", result.Items[0].OrderContentName)
	assert.Equal(testing, "Test menu description 1", result.Items[0].OrderContentDescription)
	assert.Equal(testing, 8.54, result.Items[0].OrderContentPrice)

	assert.Equal(testing, 1, result.Items[1].Quantity)
	assert.Equal(testing, "Test product 1", result.Items[1].OrderContentName)
	assert.Equal(testing, "Test product description 1", result.Items[1].OrderContentDescription)
	assert.Equal(testing, 2.50, result.Items[1].OrderContentPrice)
}

func TestPutOrderNoItems(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
		"items":        []map[string]interface{}{{}},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
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

func TestPutOrderInvalidProduct(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
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

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
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

func TestPutOrderInvalidMenu(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
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

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
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

func TestPutOrderInvalidQuantity(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
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

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
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

func TestPutOrderProductNotAvailable(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
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

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
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

func TestPutOrderMenuNotAvailable(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
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

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
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

func TestPutOrderInvalidStatus(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
		"items": []map[string]interface{}{
			{
				"quantity": 2,
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

	request, err := http.NewRequest(http.MethodPut, "/orders/3", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Order cannot be modified because it has already been prepared.")
}

func TestPutOrderError(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": 1234,
		"items":        "abcd",
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutOrderUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
		"items": []map[string]interface{}{
			{
				"quantity": 2,
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

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
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

func TestPutOrderAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
		"items": []map[string]interface{}{
			{
				"quantity": 2,
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

	request, err := http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(data))
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

func TestPutOrderNotFound(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
		"items": []map[string]interface{}{
			{
				"quantity": 2,
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

	request, err := http.NewRequest(http.MethodPut, "/orders/0", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "Order not found.")
}

func TestPutOrderInvalidId(testing *testing.T) {
	router := tests.InitTest()

	order := map[string]interface{}{
		"ticketNumber": "006",
		"items": []map[string]interface{}{
			{
				"quantity": 2,
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

	request, err := http.NewRequest(http.MethodPut, "/orders/a", bytes.NewBuffer(data))
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
