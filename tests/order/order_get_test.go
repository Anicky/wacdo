package order

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wacdo/models"
	"wacdo/tests"

	"github.com/stretchr/testify/assert"
)

func TestGetOrdersSuccess(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

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
	router := tests.InitTest()

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

func TestGetOrdersAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusForbidden, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Access not allowed.")
}

func TestGetOrderSuccess(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

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
	router := tests.InitTest()

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

func TestGetOrderAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusForbidden, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Access not allowed.")
}

func TestGetOrderNotFound(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Order not found.")
}

func TestGetOrderInvalidId(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodGet, "/orders/a", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Invalid ID.")
}
