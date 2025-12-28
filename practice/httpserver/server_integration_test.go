package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlayerServer(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := &PlayerServer{store}

	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	request := newGetScoreRequest(player)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertResponseBody(t, response.Body.String(), "3")
	assertStatus(t, response.Code, http.StatusOK)
}
