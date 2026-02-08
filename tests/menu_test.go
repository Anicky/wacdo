package tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMenusSuccess(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/menus/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test menu 1")
	assert.Contains(testing, body, "Test menu description 1")
	assert.Contains(testing, body, "8.54")
	assert.Contains(testing, body, "true")
	assert.Contains(testing, body, "Test product 1")
	assert.Contains(testing, body, "Test product description 1")
	assert.Contains(testing, body, "2.5")
	assert.Contains(testing, body, "Test product 2")
	assert.Contains(testing, body, "Test product description 2")
	assert.Contains(testing, body, "4.99")
}

func TestGetMenusUnauthorized(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/menus/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusUnauthorized, response.Code)
}

// @TOOD: get success
// @TOOD: get unauthorized
// @TOOD: get not found
// @TOOD: get invalid id
// @TOOD: post success
// @TOOD: post error
// @TOOD: put success
// @TOOD: put error
// @TOOD: put unauthorized
// @TOOD: put not found
// @TOOD: put invalid id
// @TOOD: delete success
// @TOOD: delete unauthorized
// @TOOD: delete not found
// @TOOD: delete invalid id
