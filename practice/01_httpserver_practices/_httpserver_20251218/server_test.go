package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	score    map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(player string) int {
	return s.score[player]
}

func (s *StubPlayerStore) RecordWin(player string) {
	s.winCalls = append(s.winCalls, player)
}

// ------------------------------------------------------------------------

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{}
	server := &PlayerServer{&store}
	player := "Pepper"

	request := newPostWinRequest(player)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertStatus(t, response.Code, http.StatusAccepted)

	if len(store.winCalls) != 1 {
		t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != player {
		t.Errorf("got %q want %q", store.winCalls[0], player)
	}
}

func newPostWinRequest(player string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	return req
}

// ------------------------------------------------------------------------

func TestServer(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
	}
	server := &PlayerServer{&store}
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})
	t.Run("returns 404 if no players", func(t *testing.T) {
		request := newGetScoreRequest("Asss")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

// --helpers-------------------------------------------------------------------

func newGetScoreRequest(player string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
	return req
}

// --asserts-------------------------------------------------------------------

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

// ------------------------------------------------------------------------
