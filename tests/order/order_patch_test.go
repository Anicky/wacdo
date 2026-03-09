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

func TestPatchOrderInPreparationSuccess(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/in-preparation", nil)
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
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/2/in-preparation", nil)
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
	assert.Contains(testing, body, "Order is already in preparation.")
}

func TestPatchOrderInPreparationInvalidStatus2(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/3/in-preparation", nil)
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
	assert.Contains(testing, body, "Order is already prepared.")
}

func TestPatchOrderInPreparationInvalidStatus3(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/4/in-preparation", nil)
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
	assert.Contains(testing, body, "Order is already delivered.")
}

func TestPatchOrderInPreparationUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/in-preparation", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertUnauthorized(testing, response)
}

func TestPatchOrderInPreparationAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/in-preparation", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertAccessNotAllowed(testing, response)
}

func TestPatchOrderInPreparationNotFound(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/0/in-preparation", nil)
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

func TestPatchOrderInPreparationInvalidId(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/a/in-preparation", nil)
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

func TestPatchOrderPreparedSuccess(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/2/prepared", nil)
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
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/prepared", nil)
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
	assert.Contains(testing, body, "Order must be in preparation before it can be prepared.")
}

func TestPatchOrderPreparedInvalidStatus2(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/3/prepared", nil)
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
	assert.Contains(testing, body, "Order is already prepared.")
}

func TestPatchOrderPreparedInvalidStatus3(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/4/prepared", nil)
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
	assert.Contains(testing, body, "Order is already delivered.")
}

func TestPatchOrderPreparedUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/2/prepared", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertUnauthorized(testing, response)
}

func TestPatchOrderPreparedAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/prepared", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertAccessNotAllowed(testing, response)
}

func TestPatchOrderPreparedNotFound(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/0/prepared", nil)
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

func TestPatchOrderPreparedInvalidId(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/a/prepared", nil)
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

func TestPatchOrderDeliveredSuccess(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/3/delivered", nil)
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
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/1/delivered", nil)
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
	assert.Contains(testing, body, "Order must be prepared before it can be delivered.")
}

func TestPatchOrderDeliveredInvalidStatus2(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/2/delivered", nil)
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
	assert.Contains(testing, body, "Order must be prepared before it can be delivered.")
}

func TestPatchOrderDeliveredInvalidStatus3(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/4/delivered", nil)
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
	assert.Contains(testing, body, "Order is already delivered.")
}

func TestPatchOrderDeliveredUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/3/delivered", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertUnauthorized(testing, response)
}

func TestPatchOrderDeliveredAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/3/delivered", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUser(request, 4)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertAccessNotAllowed(testing, response)
}

func TestPatchOrderDeliveredNotFound(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/0/delivered", nil)
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

func TestPatchOrderDeliveredInvalidId(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodPatch, "/orders/a/delivered", nil)
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
