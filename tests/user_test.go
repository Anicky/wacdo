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

func TestUserLoginSuccess(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "admin1@example.com",
		"password": "Admin1234!",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	assert.Contains(testing, response.Body.String(), "token")
}

func TestUserLoginInvalidUsername(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "invalid-email@example.com",
		"password": "Admin1234!",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	assert.Contains(testing, response.Body.String(), "Invalid email or password.")
}

func TestUserLoginInvalidPassword(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "admin1@example.com",
		"password": "invalid-password",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	assert.Contains(testing, response.Body.String(), "Invalid email or password.")
}

func TestGetUsersSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/users/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	var results []models.User
	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, 3, len(results))

	assert.Equal(testing, "admin1@example.com", results[0].Email)
	assert.Equal(testing, models.UserRole("admin"), results[0].Role)

	assert.Equal(testing, "greeter1@example.com", results[1].Email)
	assert.Equal(testing, models.UserRole("greeter"), results[1].Role)

	assert.Equal(testing, "greeter2@example.com", results[2].Email)
	assert.Equal(testing, models.UserRole("greeter"), results[2].Role)
}

func TestGetUsersUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/users/", nil)
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

func TestGetUserSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/users/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.User{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "admin1@example.com", result.Email)
	assert.Equal(testing, models.UserRole("admin"), result.Role)
}

func TestGetUserUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/users/1", nil)
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

func TestGetUserNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/users/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "User not found.")
}

func TestGetUserInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/users/a", nil)
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

func TestPostUserSuccess(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "manager1@example.com",
		"password": "Manager1234!",
		"role":     "manager",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusCreated, response.Code)

	result := models.User{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "manager1@example.com", result.Email)
	assert.Equal(testing, models.UserRole("manager"), result.Role)
}

func TestPostUserEmailAlreadyUsed(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "admin1@example.com",
		"password": "Admin1234!",
		"role":     "admin",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	assert.Contains(testing, response.Body.String(), "Email already used.")
}

func TestPostUserInvalidPassword(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "manager1@example.com",
		"password": "a",
		"role":     "manager",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostUserError(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    1234,
		"password": 1234,
		"role":     1234,
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPostUserUnauthorized(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "manager1@example.com",
		"password": "Manager1234!",
		"role":     "manager",
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(data))
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

func TestPutUserSuccess(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "orderpicker1@example.com",
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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	result := models.User{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Unable to decode JSON: ", err)
	}

	assert.Equal(testing, "orderpicker1@example.com", result.Email)
	assert.Equal(testing, models.UserRole("order_picker"), result.Role)
}

func TestPutUserSuccessWithSameEmail(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

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
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	assert.Contains(testing, response.Body.String(), "Email already used.")
}

func TestPutUserInvalidPassword(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutUserError(testing *testing.T) {
	router := InitTest()

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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)
}

func TestPutUserUnauthorized(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "orderpicker1@example.com",
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

func TestPutUserNotFound(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "orderpicker1@example.com",
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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "User not found.")
}

func TestPutUserInvalidId(testing *testing.T) {
	router := InitTest()

	user := map[string]interface{}{
		"email":    "orderpicker1@example.com",
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

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Invalid ID.")
}

func TestDeleteUserSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/users/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)
}

func TestDeleteUserUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/users/1", nil)
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

func TestDeleteUserNotFound(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/users/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "User not found.")
}

func TestDeleteUserInvalidId(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/users/a", nil)
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
