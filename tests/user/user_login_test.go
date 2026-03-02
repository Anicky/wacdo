package user

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wacdo/tests"

	"github.com/stretchr/testify/assert"
)

func TestUserLoginSuccess(testing *testing.T) {
	router := tests.InitTest()

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
	router := tests.InitTest()

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
	router := tests.InitTest()

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
