package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"library-api/router"

	"github.com/stretchr/testify/assert"
)

func TestGetBooks(t *testing.T) {
	testRouter := router.SetupRouter("test")
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	expectedBody := `[{"id":"2","title":"Chronicles of Narnia","author":"CS Lewis","pages":562,"checkedOut":false,"patron":""},{"id":"1","title":"Lord of the Rings","author":"JRR Tolkien","pages":1077,"checkedOut":false,"patron":""},{"id":"3","title":"Cat in the Hat","author":"Dr. Seuss","pages":32,"checkedOut":false,"patron":""},{"id":"4","title":"Super Powereds","author":"Drew Hayes","pages":398,"checkedOut":false,"patron":""},{"id":"5","title":"Tale of Two Cities","author":"Charles Dickens","pages":1200,"checkedOut":false,"patron":""}]`
	assert.JSONEq(t, expectedBody, resp.Body.String())

}

func TestGetBook(t *testing.T) {
	testRouter := router.SetupRouter("test")

	req, err := http.NewRequest("GET", "/books/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	expectedBody := `{"id": "1", "title": "Lord of the Rings", "author": "JRR Tolkien", "pages": 1077, "checkedOut": false, "patron":""}`
	assert.JSONEq(t, expectedBody, resp.Body.String())

}
