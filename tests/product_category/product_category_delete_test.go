package product_category

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wacdo/tests"

	"github.com/stretchr/testify/assert"
)

func TestDeleteProductCategorySuccess(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/3", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)
}

func TestDeleteProductCategoryErrorAssociatedProducts(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusBadRequest, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Cannot delete product category: there are products associated with it.")
}

func TestDeleteProductCategoryUnauthorized(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertUnauthorized(testing, response)
}

func TestDeleteProductCategoryAccessNotAllowed(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUser(request, 2)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	tests.AssertAccessNotAllowed(testing, response)
}

func TestDeleteProductCategoryNotFound(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/0", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	tests.AuthenticateUserAsAdmin(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusNotFound, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "error")
	assert.Contains(testing, body, "Product category not found.")
}

func TestDeleteProductCategoryInvalidId(testing *testing.T) {
	router := tests.InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/products/categories/a", nil)
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
