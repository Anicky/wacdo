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

func TestGetOrdersSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	var results []models.Order
	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, 4, len(results))

	assert.Equal(testing, "001", results[0].TicketNumber)
	assert.Equal(testing, models.Created, results[0].Status)
	assert.Equal(testing, "greeter1@example.com", results[0].User.Email)

	assert.Equal(testing, "002", results[1].TicketNumber)
	assert.Equal(testing, models.InPreparation, results[1].Status)
	assert.Equal(testing, "greeter2@example.com", results[1].User.Email)

	assert.Equal(testing, "003", results[2].TicketNumber)
	assert.Equal(testing, models.Prepared, results[2].Status)
	assert.Equal(testing, "greeter1@example.com", results[2].User.Email)

	assert.Equal(testing, "004", results[3].TicketNumber)
	assert.Equal(testing, models.Delivered, results[3].Status)
	assert.Equal(testing, "greeter2@example.com", results[3].User.Email)
}

func TestGetOrdersUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/", nil)
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

func TestGetOrderSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.Order{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "001", result.TicketNumber)
	assert.Equal(testing, models.Created, result.Status)
	assert.Equal(testing, "greeter1@example.com", result.User.Email)

	assert.Equal(testing, 2, len(result.Items))

	assert.Equal(testing, 2, result.Items[0].Quantity)
	assert.Equal(testing, "Test product 1", result.Items[0].OrderContentName)
	assert.Equal(testing, "Test product description 1", result.Items[0].OrderContentDescription)
	assert.Equal(testing, 2.5, result.Items[0].OrderContentPrice)

	assert.Equal(testing, 1, result.Items[1].Quantity)
	assert.Equal(testing, "Test menu 1", result.Items[1].OrderContentName)
	assert.Equal(testing, "Test menu description 1", result.Items[1].OrderContentDescription)
	assert.Equal(testing, 8.54, result.Items[1].OrderContentPrice)
}

func TestGetOrderUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/1", nil)
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

func TestGetOrderNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Order not found.")
}

func TestGetOrderInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/a", nil)
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

func TestPostOrderSuccess(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

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
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "At least one item must be provided.")
}

func TestPostOrderInvalidProduct(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product 9999: item not found.")
}

func TestPostOrderInvalidMenu(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 9999: item not found.")
}

func TestPostOrderInvalidQuantity(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 1: item quantity should be superior than 0.")
}

func TestPostOrderProductNotAvailable(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product 2: item is not available.")
}

func TestPostOrderMenuNotAvailable(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 2: item is not available.")
}

func TestPostOrderError(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostOrderUnauthorized(testing *testing.T) {
	router := InitTest()

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

func TestPutOrderSuccess(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

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
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "At least one item must be provided.")
}

func TestPutOrderInvalidProduct(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product 9999: item not found.")
}

func TestPutOrderInvalidMenu(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 9999: item not found.")
}

func TestPutOrderInvalidQuantity(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 1: item quantity should be superior than 0.")
}

func TestPutOrderProductNotAvailable(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product 2: item is not available.")
}

func TestPutOrderMenuNotAvailable(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 2: item is not available.")
}

func TestPutOrderInvalidStatus(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Order cannot be modified because it has already been prepared.")
}

func TestPutOrderError(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutOrderUnauthorized(testing *testing.T) {
	router := InitTest()

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

func TestPutOrderNotFound(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Order not found.")
}

func TestPutOrderInvalidId(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Invalid ID.")
}

func TestPatchOrderInPreparationSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/in-preparation", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.Order{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "001", result.TicketNumber)
	assert.Equal(testing, models.InPreparation, result.Status)
	assert.Equal(testing, "greeter1@example.com", result.User.Email)

	assert.Equal(testing, 2, len(result.Items))

	assert.Equal(testing, 2, result.Items[0].Quantity)
	assert.Equal(testing, "Test product 1", result.Items[0].OrderContentName)
	assert.Equal(testing, "Test product description 1", result.Items[0].OrderContentDescription)
	assert.Equal(testing, 2.5, result.Items[0].OrderContentPrice)

	assert.Equal(testing, 1, result.Items[1].Quantity)
	assert.Equal(testing, "Test menu 1", result.Items[1].OrderContentName)
	assert.Equal(testing, "Test menu description 1", result.Items[1].OrderContentDescription)
	assert.Equal(testing, 8.54, result.Items[1].OrderContentPrice)
}

func TestPatchOrderInPreparationInvalidStatus1(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/2/in-preparation", nil)
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
	assert.Contains(testing, body, "Order is already in preparation.")
}

func TestPatchOrderInPreparationInvalidStatus2(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/3/in-preparation", nil)
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
	assert.Contains(testing, body, "Order is already prepared.")
}

func TestPatchOrderInPreparationInvalidStatus3(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/4/in-preparation", nil)
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
	assert.Contains(testing, body, "Order is already delivered.")
}

func TestPatchOrderInPreparationUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/in-preparation", nil)
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

func TestPatchOrderInPreparationNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/0/in-preparation", nil)
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
	assert.Contains(testing, body, "Order not found.")
}

func TestPatchOrderInPreparationInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/a/in-preparation", nil)
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

func TestPatchOrderPreparedSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/2/prepared", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.Order{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "002", result.TicketNumber)
	assert.Equal(testing, models.Prepared, result.Status)
	assert.Equal(testing, "greeter2@example.com", result.User.Email)

	assert.Equal(testing, 1, len(result.Items))

	assert.Equal(testing, 1, result.Items[0].Quantity)
	assert.Equal(testing, "Test product 2", result.Items[0].OrderContentName)
	assert.Equal(testing, "Test product description 2", result.Items[0].OrderContentDescription)
	assert.Equal(testing, 4.99, result.Items[0].OrderContentPrice)
}

func TestPatchOrderPreparedInvalidStatus1(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/prepared", nil)
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
	assert.Contains(testing, body, "Order must be in preparation before it can be prepared.")
}

func TestPatchOrderPreparedInvalidStatus2(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/3/prepared", nil)
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
	assert.Contains(testing, body, "Order is already prepared.")
}

func TestPatchOrderPreparedInvalidStatus3(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/4/prepared", nil)
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
	assert.Contains(testing, body, "Order is already delivered.")
}

func TestPatchOrderPreparedUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/prepared", nil)
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

func TestPatchOrderPreparedNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/0/prepared", nil)
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
	assert.Contains(testing, body, "Order not found.")
}

func TestPatchOrderPreparedInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/a/prepared", nil)
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

func TestPatchOrderDeliveredSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/3/delivered", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.Order{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "003", result.TicketNumber)
	assert.Equal(testing, models.Delivered, result.Status)
	assert.Equal(testing, "greeter1@example.com", result.User.Email)

	assert.Equal(testing, 1, len(result.Items))

	assert.Equal(testing, 1, result.Items[0].Quantity)
	assert.Equal(testing, "Test menu 2", result.Items[0].OrderContentName)
	assert.Equal(testing, "Test menu description 2", result.Items[0].OrderContentDescription)
	assert.Equal(testing, 7.2, result.Items[0].OrderContentPrice)
}

func TestPatchOrderDeliveredInvalidStatus1(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/delivered", nil)
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
	assert.Contains(testing, body, "Order must be prepared before it can be delivered.")
}

func TestPatchOrderDeliveredInvalidStatus2(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/2/delivered", nil)
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
	assert.Contains(testing, body, "Order must be prepared before it can be delivered.")
}

func TestPatchOrderDeliveredInvalidStatus3(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/4/delivered", nil)
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
	assert.Contains(testing, body, "Order is already delivered.")
}

func TestPatchOrderDeliveredUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/delivered", nil)
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

func TestPatchOrderDeliveredNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/0/delivered", nil)
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
	assert.Contains(testing, body, "Order not found.")
}

func TestPatchOrderDeliveredInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/a/delivered", nil)
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
