package user

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

func TestPutUserSuccess(testing *testing.T) {
	router := tests.InitTest()

	user := map[string]interface{}{
		"email":    "orderpicker2@example.com",
		"password": "OrderPicker1234!",
		"role":     "order_picker",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.User{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "orderpicker2@example.com", result.Email)
	assert.Equal(testing, models.UserRole("order_picker"), result.Role)
}

func TestPutUserSuccessWithSameEmail(testing *testing.T) {
	router := tests.InitTest()

	user := map[string]interface{}{
		"email":    "admin1@example.com",
		"password": "OrderPicker1234!",
		"role":     "order_picker",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.User{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "admin1@example.com", result.Email)
	assert.Equal(testing, models.UserRole("order_picker"), result.Role)
}

func TestPutUserEmailAlreadyUsed(testing *testing.T) {
	router := tests.InitTest()

	user := map[string]interface{}{
		"email":    "greeter1@example.com",
		"password": "Admin1234!",
		"role":     "admin",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	assert.Contains(testing, response.Body.String(), "Email already used.")
}

func TestPutUserInvalidPassword(testing *testing.T) {
	router := tests.InitTest()

	user := map[string]interface{}{
		"email":    "manager1@example.com",
		"password": "a",
		"role":     "manager",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutUserError(testing *testing.T) {
	router := tests.InitTest()

	user := map[string]interface{}{
		"email":    1234,
		"password": 1234,
		"role":     1234,
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutUserUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	user := map[string]interface{}{
		"email":    "orderpicker2@example.com",
		"password": "OrderPicker1234!",
		"role":     "order_picker",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(data))
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

func TestPutUserAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	user := map[string]interface{}{
		"email":    "orderpicker2@example.com",
		"password": "OrderPicker1234!",
		"role":     "order_picker",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(data))
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

func TestPutUserNotFound(testing *testing.T) {
	router := tests.InitTest()

	user := map[string]interface{}{
		"email":    "orderpicker2@example.com",
		"password": "OrderPicker1234!",
		"role":     "order_picker",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/users/0", bytes.NewBuffer(data))
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
	assert.Contains(testing, body, "User not found.")
}

func TestPutUserInvalidId(testing *testing.T) {
	router := tests.InitTest()

	user := map[string]interface{}{
		"email":    "orderpicker2@example.com",
		"password": "OrderPicker1234!",
		"role":     "order_picker",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/users/a", bytes.NewBuffer(data))
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
