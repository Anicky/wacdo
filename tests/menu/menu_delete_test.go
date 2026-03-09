package menu

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wacdo/tests"

	"github.com/stretchr/testify/assert"
)

func TestDeleteMenuSuccess(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/menus/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)
}

func TestDeleteMenuUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/menus/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertUnauthorized(testing, response)
}

func TestDeleteMenuAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/menus/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertAccessNotAllowed(testing, response)
}

func TestDeleteMenuNotFound(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/menus/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Menu 0: item not found.")
}

func TestDeleteMenuInvalidId(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/menus/a", nil)
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
