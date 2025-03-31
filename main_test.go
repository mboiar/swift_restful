package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

// func TestPostUser(t *testing.T) {
// 	router := setupRouter()
// 	// router = postUser(router)

// 	w := httptest.NewRecorder()

// 	// Create an example user for testing
// 	// exampleUser := User{
// 	// 	Username: "test_name",
// 	// 	Gender:   "male",
// 	// }
// 	// userJson, _ := json.Marshal(exampleUser)
// 	// req, _ := http.NewRequest("POST", "/user/add", strings.NewReader(string(userJson)))
// 	// router.ServeHTTP(w, req)

// 	// assert.Equal(t, 200, w.Code)
// 	// // Compare the response body with the json data of exampleUser
// 	// assert.Equal(t, string(userJson), w.Body.String())
// }
