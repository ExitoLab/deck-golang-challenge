package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	controller "golang-execrise/controllers"
)

type Card struct {
	Rank string `json:"rank"`
	Suit string `json:"suit"`
	Code string `json:"code"`
}

// Deck represents a standard 52-card deck of playing cards
type Deck struct {
	Deck_ID   string `json:"deck_id"`
	Cards     []Card `json:"cards"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

// TestNewDeckHandler tests the NewDeckHandler function
func TestNewDeckHandler(t *testing.T) {
	// create a new HTTP request
	req, err := http.NewRequest("POST", "/create/deck", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a new ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// create a new gin engine and router
	engine := gin.Default()
	router := engine.RouterGroup

	// define the routes for the application
	router.POST("/create/deck", controller.CreateNewDeckHandler())

	// call the HTTP handler function with the request and response recorder
	engine.ServeHTTP(rr, req)

	// decode the JSON response into a Deck struct
	var deck Deck
	err = json.Unmarshal(rr.Body.Bytes(), &deck)
	if err != nil {
		t.Fatal(err)
	}

	// check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// check that the deck has the correct properties
	assert.Equal(t, 52, deck.Remaining)
	assert.False(t, deck.Shuffled)

	// Assert that the response is a 200 OK status code and contains the expected message
	assert.Equal(t, http.StatusOK, rr.Code)

}
