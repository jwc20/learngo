package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jwc20/learngowithtests/websockets"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := poker.NewFileSystemPlayerStore(database)
	assertNoError(t, err)

	server, err := poker.NewPlayerServer(store, dummyGame)
	if err != nil {
		t.Fatalf("problem creating player server %v", err)
	}

	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get League", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []poker.Player{
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})
}
